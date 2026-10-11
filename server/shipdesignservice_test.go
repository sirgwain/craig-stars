package server

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
	craig_starsv1 "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_shipDesignService_DeleteShipDesign_removesQueuedItems(t *testing.T) {
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
	fullGame, err := gr.HostGame(user.ID, cs.NewGameSettings().WithHost(cs.Humanoids()))
	require.NoError(t, err)

	player := fullGame.Players[0]
	var deleted, kept *cs.ShipDesign
	for _, design := range player.Designs {
		if design.Spec.Starbase || design.CannotDelete {
			continue
		}
		if deleted == nil {
			deleted = design
		} else if kept == nil {
			kept = design
		}
	}
	require.NotNil(t, deleted)
	require.NotNil(t, kept)

	planets, err := client.GetPlanetsForPlayer(ctx, fullGame.ID, player.Num)
	require.NoError(t, err)
	require.NotEmpty(t, planets)
	planet := planets[0]
	planet.ProductionQueue = []cs.ProductionQueueItem{
		{Type: cs.QueueItemTypeShipToken, Quantity: 2, DesignNum: deleted.Num},
		{Type: cs.QueueItemTypeMine, Quantity: 5},
		{Type: cs.QueueItemTypeShipToken, Quantity: 1, DesignNum: kept.Num},
		{Type: cs.QueueItemTypeShipToken, Quantity: 3, DesignNum: deleted.Num},
	}
	require.NoError(t, client.SavePlanet(ctx, planet))

	game, err := client.GetGame(ctx, fullGame.ID)
	require.NoError(t, err)
	ctx = context.WithValue(ctx, keyDbRead, dbConn.NewReadClient())
	ctx = context.WithValue(ctx, keyGame, game)
	ctx = context.WithValue(ctx, keyGamePlayer, &game.Players[0])
	response, err := NewshipDesignServiceHandler(dbConn).DeleteShipDesign(ctx, connect.NewRequest(&craig_starsv1.DeleteShipDesignRequest{
		GameId: fullGame.ID,
		Num:    int32(deleted.Num),
	}))
	require.NoError(t, err)

	// the deleted design is gone from the saved queue, everything else is left in order
	saved, err := client.GetPlanetByNum(ctx, fullGame.ID, planet.Num)
	require.NoError(t, err)
	require.Len(t, saved.ProductionQueue, 2)
	assert.Equal(t, cs.QueueItemTypeMine, saved.ProductionQueue[0].Type)
	assert.Equal(t, kept.Num, saved.ProductionQueue[1].DesignNum)

	// the response has the updated queue with estimates for the remaining items
	var responseQueue []*craig_starsv1.ProductionQueueItem
	for _, p := range response.Msg.Planets {
		if int(p.MapObject.Num) == planet.Num {
			responseQueue = p.PlanetOrders.ProductionQueue
		}
	}
	require.Len(t, responseQueue, 2)
	assert.Equal(t, int32(kept.Num), responseQueue[1].DesignNum)
	assert.NotZero(t, responseQueue[1].QueueItemCompletionEstimate.GetYearsToBuildAll())
}
