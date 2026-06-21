package server

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"log/slog"

	"github.com/disgoorg/snowflake/v2"
	"github.com/go-chi/render"
	"github.com/sirgwain/craig-stars/config"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/webhook"
)

type discordWebhook struct {
	id    string
	token string
}

// discordNotifier handles Discord webhook notifications
type discordNotifier struct {
	db     DBConnection
	config config.Config
}

// newDiscordNotifier creates a new Discord notifier
func newDiscordNotifier(db DBConnection, cfg config.Config) *discordNotifier {
	return &discordNotifier{
		db:     db,
		config: cfg,
	}
}

// SendNewTurnNotification sends a notification about a new turn
// this will not send for single players games or games with the admin (my tests)
func (d *discordNotifier) SendNewTurnNotification(gameID int64) {
	if !d.config.Discord.WebhookNotify {
		// no webhook notifications for this server
		return
	}

	readClient := d.db.NewReadClient()
	ctx := context.Background()
	go func() {
		game, err := readClient.GetGame(ctx, gameID)
		if err != nil {
			slog.Error("get game for discord notification", slog.Any("error", err))
			return
		}

		if game.IsSinglePlayer() {
			return
		}

		// don't notify for the admin
		// don't want to spam the real discord when I'm testing
		if !d.config.Discord.WebhookNotifyForAdmin {
			for _, player := range game.Players {
				if player.Name == "admin" {
					return
				}
			}
		}

		webhooks := []discordWebhook{}

		host, err := readClient.GetUser(ctx, game.HostID)
		if err != nil {
			slog.Error("get host user for game for discord notification", slog.Any("error", err))
			return
		}

		users, err := readClient.GetUsersForGame(ctx, gameID)
		if err != nil {
			slog.Error("get users for game for discord notification", slog.Any("error", err))
			return
		}

		// if the host has their own webhook url, they probably have their own server. Don't notify the main server
		if host.DiscordWebhookURL == "" && d.config.Discord.WebhookID != "" && d.config.Discord.WebhookToken != "" {
			// no host webhook, so setup the main game webhook
			webhooks = append(webhooks, discordWebhook{
				id:    d.config.Discord.WebhookID,
				token: d.config.Discord.WebhookToken,
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
					slog.Error("unable to parse user webhook url", slog.Any("error", err), slog.String("username", user.Username), slog.String("webhookURL", user.DiscordWebhookURL))
					continue
				}
				webhooks = append(webhooks, discordWebhook{id: id, token: token})
			}
		}

		slog.Debug("notifying players of new turn",
			slog.Int64("gameID", gameID),
			slog.Int("webhooks", len(webhooks)),
			slog.Int("userAts", len(userAts)),
		)
		for _, hook := range webhooks {
			d.sendWebhookMessage(ctx, hook, discord.NewWebhookMessageCreate().
				WithContentf("**%s** has a new turn. \n%s", game.Name, strings.Join(userAts, ", ")).
				WithEmbeds(discord.NewEmbed().
					WithTitlef("%s - %d", game.Name, game.Year).
					WithURLf("%s/games/%d", d.config.Auth.URL, game.ID)))
		}
	}()
}

// SendTestWebhook sends a test message to a user's webhook
func (d *discordNotifier) SendTestWebhook(ctx context.Context, webhookURL, discordID string) error {
	webhookID, token, err := parseDiscordWebhookUrl(webhookURL)
	if err != nil {
		return err
	}

	hook := discordWebhook{id: webhookID, token: token}

	go func() {
		d.sendWebhookMessage(ctx, hook, discord.NewWebhookMessageCreate().
			WithContentf("This is a test of your discord webhook.\n<@%s>", discordID).
			WithEmbeds(discord.NewEmbed().
				WithTitlef("craig-stars").
				WithURLf("%s", d.config.Auth.URL)))
	}()

	return nil
}

// sendWebhookMessage sends a message to a Discord webhook
func (d *discordNotifier) sendWebhookMessage(_ context.Context, hook discordWebhook, message discord.WebhookMessageCreate) {
	// construct new webhook client
	// https://discord.com/api/webhooks/<id>/<token>
	id, err := snowflake.Parse(hook.id)
	if err != nil {
		slog.Error("parse discord webhook id", slog.Any("error", err))
		return
	}
	client := webhook.New(id, hook.token)

	defer client.Close(context.TODO())

	if _, err := client.CreateMessage(message,
		rest.CreateWebhookMessageParams{},
		// delay each request by 2 seconds
		rest.WithDelay(2*time.Second),
	); err != nil {
		slog.Error("sending discord message", slog.Any("error", err))
	}
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

func (s *server) pingDiscordForGameUpdate(w http.ResponseWriter, r *http.Request) {
	user := s.contextUserSession(r)
	game := s.contextGame(r)

	if user.ID != game.HostID {
		slog.Error("access denied for sending discord updates", slog.Int64("GameID", game.ID), slog.String("User", user.Username))
		render.Render(w, r, ErrForbidden)
		return
	}

	s.discordNotifier.SendNewTurnNotification(game.ID)
}
