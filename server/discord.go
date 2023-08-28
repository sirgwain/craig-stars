package server

import (
	"context"
	"fmt"
	"net/http"
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
	user := s.contextUser(r)
	game := s.contextGame(r)
	db := s.contextDb(r)

	if user.ID != game.HostID {
		log.Error().Int64("GameID", game.ID).Str("User", user.Username).Msg("access denied for sending discord updates")
		render.Render(w, r, ErrForbidden)
		return
	}

	s.sendNewTurnNotification(db, game.ID)
}

// send a notification about a new turn
// this will not send for single players games or games with the admin (my tests)
func (s *server) sendNewTurnNotification(db TXClient, gameID int64) {
	if !s.config.Discord.WebhookNotify {
		// no webhook notifications for this server
		return
	}

	go func() {
		game, err := db.GetGame(gameID)
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

		users, err := db.GetUsersForGame(gameID)
		if err != nil {
			log.Error().Err(err).Msg("get users for game for discord notification")
			return
		}

		userAts := make([]string, 0, len(users))
		for _, user := range users {
			if user.DiscordID != nil && *user.DiscordID != "" {
				userAts = append(userAts, fmt.Sprintf("<@%s>", *user.DiscordID))
			}
		}

		// construct new webhook client
		// https://discord.com/api/webhooks/<id>/<token>
		id, err := snowflake.Parse(s.config.Discord.WebhookID)
		if err != nil {
			log.Error().Err(err).Msg("parse discord webhook id")
			return
		}
		client := webhook.New(id, s.config.Discord.WebhookToken)

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
	}()

}
