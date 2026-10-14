package server

import (
	"context"
	"strings"
	"testing"

	"connectrpc.com/connect"
	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
	craig_starsv1 "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_playerMailService(t *testing.T) {
	ctx := t.Context()
	dbConn := db.NewConn()
	cfg := config.Config{}
	cfg.Database.Filename = ":memory:"
	require.NoError(t, dbConn.Connect(ctx, &cfg))
	t.Cleanup(func() { dbConn.Close() })
	client := dbConn.NewReadWriteClient()
	user, err := client.CreateUser(ctx, cs.NewUser("host", "", "", cs.RoleUser))
	require.NoError(t, err)
	gr := NewGameRunner(dbConn, cfg)
	fullGame, err := gr.HostGame(user.ID, cs.NewGameSettings().
		WithHost(cs.Humanoids()).
		WithAIPlayer(cs.AIDifficultyNormal, 0).
		WithAIPlayer(cs.AIDifficultyNormal, 0))
	require.NoError(t, err)

	baseCtx := context.WithValue(ctx, keyDbRead, dbConn.NewReadClient())
	baseCtx = context.WithValue(baseCtx, keyDbWrite, client)
	// a context for a player's request, with the game as it is now
	playerCtx := func(num int) context.Context {
		game, err := client.GetGame(ctx, fullGame.ID)
		require.NoError(t, err)
		require.Len(t, game.Players, 3)
		c := context.WithValue(baseCtx, keyGame, game)
		return context.WithValue(c, keyGamePlayer, &game.Players[num-1])
	}

	service := NewPlayerMailServiceHandler()
	save := func(ctx context.Context, mail *craig_starsv1.PlayerMail) (*craig_starsv1.PlayerMail, error) {
		response, err := service.SavePlayerMail(ctx, connect.NewRequest(&craig_starsv1.SavePlayerMailRequest{GameId: fullGame.ID, Mail: mail}))
		if err != nil {
			return nil, err
		}
		return response.Msg.Mail, nil
	}
	remove := func(ctx context.Context, id int64) error {
		_, err := service.DeletePlayerMail(ctx, connect.NewRequest(&craig_starsv1.DeletePlayerMailRequest{GameId: fullGame.ID, MailId: id}))
		return err
	}
	mail := func(ctx context.Context) []*craig_starsv1.PlayerMail {
		response, err := service.GetPlayerMail(ctx, connect.NewRequest(&craig_starsv1.GetPlayerMailRequest{GameId: fullGame.ID}))
		require.NoError(t, err)
		return response.Msg.Mail
	}
	delivered := func(playerNum int) []cs.PlayerMessage {
		player, err := client.GetPlayerForGame(ctx, fullGame.ID, playerNum)
		require.NoError(t, err)
		messages := []cs.PlayerMessage{}
		for _, message := range player.Messages {
			if message.Type == cs.PlayerMessagePlayerMail {
				messages = append(messages, message)
			}
		}
		return messages
	}

	p1, p2, p3 := playerCtx(1), playerCtx(2), playerCtx(3)

	// the sender and year come from the request's player and game, not the mail
	direct, err := save(p1, &craig_starsv1.PlayerMail{FromPlayerNum: 3, Year: 2500, ToPlayerNum: 2, Text: " truce "})
	require.NoError(t, err)
	assert.Equal(t, "truce", direct.Text)
	assert.Equal(t, int32(1), direct.FromPlayerNum)
	assert.Equal(t, int32(fullGame.Year), direct.Year)
	assert.NotZero(t, direct.GameDbObject.Id)

	// mail can be changed until it is delivered
	direct.Text = "truce?"
	direct, err = save(p1, direct)
	require.NoError(t, err)

	everybody, err := save(p1, &craig_starsv1.PlayerMail{Text: "hi all"})
	require.NoError(t, err)
	deleted, err := save(p1, &craig_starsv1.PlayerMail{ToPlayerNum: 3, Text: "never mind"})
	require.NoError(t, err)
	require.NoError(t, remove(p1, deleted.GameDbObject.Id))

	_, err = save(p1, &craig_starsv1.PlayerMail{ToPlayerNum: 2, Text: " "})
	assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
	_, err = save(p1, &craig_starsv1.PlayerMail{ToPlayerNum: 4, Text: "hi"})
	assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
	_, err = save(p1, &craig_starsv1.PlayerMail{ToPlayerNum: 2, Text: strings.Repeat("a", cs.PlayerMailMaxLength+1)})
	assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))

	// only the sender sees mail before it is delivered, and only the sender can change it
	assert.Len(t, mail(p1), 2)
	assert.Empty(t, mail(p2))
	assert.Empty(t, mail(p3))
	_, err = save(p2, direct)
	assert.Equal(t, connect.CodeNotFound, connect.CodeOf(err))
	assert.Equal(t, connect.CodeNotFound, connect.CodeOf(remove(p2, direct.GameDbObject.Id)))

	_, err = gr.GenerateTurn(fullGame.ID)
	require.NoError(t, err)
	p1, p2, p3 = playerCtx(1), playerCtx(2), playerCtx(3)

	// the turn delivers the mail as messages
	assert.Empty(t, delivered(1))
	got := delivered(2)
	require.Len(t, got, 2)
	assert.Equal(t, "truce?", got[0].Spec.Note)
	assert.Equal(t, 1, got[0].Spec.SourcePlayerNum)
	assert.Equal(t, "hi all", got[1].Spec.Note)
	got = delivered(3)
	require.Len(t, got, 1)
	assert.Equal(t, "hi all", got[0].Spec.Note)

	// and the recipients keep it
	assert.Len(t, mail(p1), 2)
	assert.Len(t, mail(p2), 2)
	require.Len(t, mail(p3), 1)
	assert.Equal(t, everybody.GameDbObject.Id, mail(p3)[0].GameDbObject.Id)

	// delivered mail can't be changed
	_, err = save(p1, direct)
	assert.Equal(t, connect.CodeFailedPrecondition, connect.CodeOf(err))
	assert.Equal(t, connect.CodeFailedPrecondition, connect.CodeOf(remove(p1, direct.GameDbObject.Id)))

	// mail is only delivered once
	_, err = gr.GenerateTurn(fullGame.ID)
	require.NoError(t, err)
	assert.Empty(t, delivered(2))
	assert.Len(t, mail(playerCtx(2)), 2)

	// a player can only write so much mail a year
	p1 = playerCtx(1)
	for range cs.PlayerMailMaxPerYear {
		_, err = save(p1, &craig_starsv1.PlayerMail{ToPlayerNum: 2, Text: "spam"})
		require.NoError(t, err)
	}
	_, err = save(p1, &craig_starsv1.PlayerMail{ToPlayerNum: 2, Text: "spam"})
	assert.Equal(t, connect.CodeResourceExhausted, connect.CodeOf(err))
}
