//go:build !wasi && !wasm

package server

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/disgoorg/snowflake/v2"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/webhook"
)

func (s *server) pingDiscordForGameUpdate(w http.ResponseWriter, r *http.Request) {
	user := s.contextUserSession(r)
	game := s.contextGame(r)

	if user.ID != game.HostID {
		log.Error().Int64("GameID", game.ID).Str("User", user.Username).Msg("access denied for sending discord updates")
		render.Render(w, r, ErrForbidden)
		return
	}

	s.sendNewTurnNotification(r.Context(), game.ID)
}

// send a notification about a new turn
// this will not send for single players games or games with the admin (my tests)
func (s *server) sendNewTurnNotification(ctx context.Context, gameID int64) {
	if !s.config.Discord.WebhookNotify {
		// no webhook notifications for this server
		return
	}

	readClient := s.db.NewReadClient()

	type discordWebhook struct {
		id    string
		token string
	}

	go func() {
		game, err := readClient.GetGame(ctx, gameID)
		if err != nil {
			log.Error().Err(err).Msg("get game for discord notification")
			return
		}

		if game.IsSinglePlayer() {
			return
		}

		// don't notify for the admin
		// don't want to spam the real discord when I'm testing
		if !s.config.Discord.WebhookNotifyForAdmin {
			for _, player := range game.Players {
				if player.Name == "admin" {
					return
				}
			}
		}

		webhooks := []discordWebhook{}

		host, err := readClient.GetUser(ctx, game.HostID)
		if err != nil {
			log.Error().Err(err).Msg("get host user for game for discord notification")
			return
		}

		users, err := readClient.GetUsersForGame(ctx, gameID)
		if err != nil {
			log.Error().Err(err).Msg("get users for game for discord notification")
			return
		}

		// if the host has their own webhook url, they probably have their own server. Don't notify the main server
		if host.DiscordWebhookURL == "" && s.config.Discord.WebhookID != "" && s.config.Discord.WebhookToken != "" {
			// no host webhook, so setup the main game webhook
			webhooks = append(webhooks, discordWebhook{
				id:    s.config.Discord.WebhookID,
				token: s.config.Discord.WebhookToken,
			})
		}

		userAts := make([]string, 0, len(users))
		for _, user := range users {
			if user.DiscordID != "" {
				userAts = append(userAts, fmt.Sprintf("<@%s>", user.DiscordID))
			}

			// if this user has their own webhook, notify them of a new game
			if user.DiscordWebhookURL != "" {
				// this user has their own webhook, add it to the list of webhooks to call
				id, token, err := parseDiscordWebhookUrl(user.DiscordWebhookURL)
				if err != nil {
					// don't fail on bad user data, just log it and move on
					log.Error().
						Err(err).
						Msgf("unable to parse user %s webhook url: %s", user.Username, user.DiscordWebhookURL)
					continue
				}
				webhooks = append(webhooks, discordWebhook{id: id, token: token})
			}
		}

		log.Debug().Msgf("notifying players of game %d of new turn at %d webhooks", gameID, len(webhooks))
		for _, hook := range webhooks {

			// construct new webhook client
			// https://discord.com/api/webhooks/<id>/<token>
			id, err := snowflake.Parse(hook.id)
			if err != nil {
				log.Error().Err(err).Msg("parse discord webhook id")
				continue
			}
			client := webhook.New(id, hook.token)

			defer client.Close(context.TODO())

			if _, err := client.CreateMessage(discord.NewWebhookMessageCreateBuilder().
				SetContentf("**%s** has a new turn. \n%s", game.Name, strings.Join(userAts, ", ")).
				SetEmbeds(discord.NewEmbedBuilder().
					SetTitlef("%s - %d", game.Name, game.Year).
					SetURLf("%s/games/%d", s.config.Auth.URL, game.ID).
					Build()).
				Build(),
				// delay each request by 2 seconds
				rest.WithDelay(2*time.Second),
			); err != nil {
				log.Error().Err(err).Msgf("sending discord message")
			}
		}

	}()
}

var discordWebhookRegex = regexp.MustCompile(`^https://discord\.com/api/webhooks/([^/\s]+)/([^/\s]+)$`)

// parseDiscordWebhookUrl parses a url for a webhook id and token
func parseDiscordWebhookUrl(url string) (id, token string, err error) {
	matches := discordWebhookRegex.FindStringSubmatch(url)
	if matches == nil || len(matches) != 3 {
		return "", "", fmt.Errorf("webhook is unsupported format")
	}

	id = matches[1]
	token = matches[2]
	return id, token, nil
}

// testDiscordWebhook test a user's webhook
func (s *server) testDiscordWebhook(w http.ResponseWriter, r *http.Request) {
	user := s.contextUser(r)

	if user.DiscordWebhookURL == "" {
		render.Render(w, r, ErrBadRequest(fmt.Errorf("no webhook url for user")))
		return
	}

	if user.DiscordID == "" {
		render.Render(w, r, ErrBadRequest(fmt.Errorf("no discord id user")))
		return
	}

	webhookID, token, err := parseDiscordWebhookUrl(user.DiscordWebhookURL)
	if err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	userAt := user.DiscordID
	log.Info().Msgf("sending test discord message for %s", user.Username)
	go func() {

		// construct new webhook client
		// https://discord.com/api/webhooks/<id>/<token>
		id, err := snowflake.Parse(webhookID)
		if err != nil {
			log.Error().Err(err).Msg("parse discord webhook id")
			return
		}
		client := webhook.New(id, token)

		defer client.Close(context.TODO())

		if _, err := client.CreateMessage(discord.NewWebhookMessageCreateBuilder().
			SetContentf("This is a test of your discord webhook.\n<@%s>", userAt).
			SetEmbeds(discord.NewEmbedBuilder().
				SetTitlef("craig-stars").
				SetURLf("%s", s.config.Auth.URL).
				Build()).
			Build(),
			// delay each request by 2 seconds
			rest.WithDelay(2*time.Second),
		); err != nil {
			log.Error().Err(err).Msgf("sending test discord message")
		}
	}()
}
