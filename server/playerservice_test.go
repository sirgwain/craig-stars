package server

import (
	"context"
	"path/filepath"
	"sync"
	"testing"

	"connectrpc.com/connect"
	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
	craig_starsv1 "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1"
	"github.com/stretchr/testify/require"
)

// connect to a fresh file database (WAL, separate read and write pools, like the real server)
func connectSubmitTestDB(t *testing.T) (db.DBConn, config.Config) {
	t.Helper()
	cfg := config.Config{}
	cfg.Database.Filename = filepath.Join(t.TempDir(), "data.db")
	cfg.Database.ReadConnectionParams = "?_txlock=deferred"
	cfg.Database.WriteConnectionParams = "?_txlock=immediate&_busy_timeout=1200"
	cfg.Database.SkipUpgrade = true
	dbConn := db.NewConn()
	require.NoError(t, dbConn.Connect(context.Background(), &cfg))
	t.Cleanup(func() { dbConn.Close() })
	return dbConn, cfg
}

// host a started game with two human players (users 1 and 2) and one AI player
func hostTwoHumanGame(t *testing.T, dbConn db.DBConn, cfg config.Config) int64 {
	t.Helper()
	ctx := context.Background()
	c := dbConn.NewReadWriteClient()
	host, err := c.CreateUser(ctx, cs.NewUser("host", "pw", "host@example.com", cs.RoleUser))
	require.NoError(t, err)
	other, err := c.CreateUser(ctx, cs.NewUser("other", "pw", "other@example.com", cs.RoleUser))
	require.NoError(t, err)

	gr := NewGameRunner(dbConn, cfg)
	fullGame, err := gr.HostGame(host.ID, cs.NewGameSettings().
		WithSize(cs.SizeTiny).
		WithHost(cs.Humanoids()).
		WithOpenPlayerSlot().
		WithAIPlayer(cs.AIDifficultyNormal, 0))
	require.NoError(t, err)
	require.NoError(t, gr.JoinGame(fullGame.ID, other.ID, "other", cs.Humanoids()))
	require.NoError(t, gr.StartGame(fullGame.Game))
	return fullGame.ID
}

// call the SubmitTurn handler as the given player, with the context the interceptors would build
func submitAs(t *testing.T, svc *playerService, dbConn db.DBConn, gameID int64, playerNum int) error {
	ctx := context.Background()
	game, err := dbConn.NewReadClient().GetGame(ctx, gameID)
	if err != nil {
		return err
	}
	var gamePlayer *cs.GamePlayer
	for i := range game.Players {
		if game.Players[i].Num == playerNum {
			gamePlayer = &game.Players[i]
		}
	}
	require.NotNil(t, gamePlayer)
	ctx = context.WithValue(ctx, keyDbRead, dbConn.NewReadClient())
	ctx = context.WithValue(ctx, keyDbWrite, dbConn.NewReadWriteClient())
	ctx = context.WithValue(ctx, keyGame, game)
	ctx = context.WithValue(ctx, keyGamePlayer, gamePlayer)
	_, err = svc.SubmitTurn(ctx, connect.NewRequest(&craig_starsv1.SubmitTurnRequest{GameId: gameID}))
	return err
}

// When the last two human players submit at the same moment, exactly one turn must be generated:
// not zero (the game would sit with everyone submitted) and not two.
func Test_playerService_SubmitTurn_concurrentLastSubmits(t *testing.T) {
	dbConn, cfg := connectSubmitTestDB(t)
	gameID := hostTwoHumanGame(t, dbConn, cfg)
	svc := NewPlayerServiceHandler(dbConn, cfg, newDiscordNotifier(dbConn, cfg)).(*playerService)
	read := dbConn.NewReadClient()

	rounds := 25
	if testing.Short() {
		rounds = 5
	}
	for round := 0; round < rounds; round++ {
		before, err := read.GetGame(context.Background(), gameID)
		require.NoError(t, err)

		start := make(chan struct{})
		errs := make([]error, 2)
		var wg sync.WaitGroup
		for i, num := range []int{1, 2} {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				errs[i] = submitAs(t, svc, dbConn, gameID, num)
			}()
		}
		close(start)
		wg.Wait()
		require.NoError(t, errs[0])
		require.NoError(t, errs[1])

		after, err := read.GetGame(context.Background(), gameID)
		require.NoError(t, err)
		require.Equalf(t, before.Year+1, after.Year, "round %d: both players submitted at once, want exactly one new turn (state %s)", round, after.State)
		require.Equal(t, cs.GameStateWaitingForPlayers, after.State)
		for _, p := range after.Players {
			if !p.AIControlled {
				require.Falsef(t, p.SubmittedTurn, "round %d: player %d still submitted after the new turn", round, p.Num)
			}
		}
	}
}
