package server

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
	"github.com/sirgwain/craig-stars/proto/converter"
	craig_starsv1 "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func Test_fleetService_SplitFleet(t *testing.T) {
	tests := []struct {
		name         string
		existingDest bool
		allRight     bool
		partial      bool
		destBaseName string
	}{
		{name: "new destination with no ships moved"},
		{name: "new destination after moving all ships back", destBaseName: "Another design"},
		{name: "existing destination with all ships moved left", existingDest: true},
		{name: "existing destination with all ships moved right", existingDest: true, allRight: true},
		{name: "new destination with all ships moved right", allRight: true},
		{name: "new destination with some ships moved right", partial: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			require.GreaterOrEqual(t, len(fullGame.Fleets), 2)

			source := fullGame.Fleets[0]
			source.Tokens[0].Quantity = 2
			require.NoError(t, client.SaveFleet(ctx, source))
			fleetsBefore, err := client.GetFleetsForPlayer(ctx, fullGame.ID, source.PlayerNum)
			require.NoError(t, err)
			shipsBefore := map[int]int{}
			cargoBefore := cs.Cargo{}
			fuelBefore := 0
			for _, fleet := range fleetsBefore {
				for _, token := range fleet.Tokens {
					shipsBefore[token.DesignNum] += token.Quantity
				}
				cargoBefore = cargoBefore.Add(fleet.Cargo)
				fuelBefore += fleet.Fuel
			}

			sourceTokens := converter.C.ConvertCSFleet(source).Tokens
			var destNum int32
			if tt.existingDest {
				dest := fullGame.Fleets[1]
				destNum = int32(dest.Num)
				sourceTokens = append(sourceTokens, converter.C.ConvertCSFleet(dest).Tokens...)
			}
			destTokens := make([]*craig_starsv1.ShipToken, len(sourceTokens))
			for i, token := range sourceTokens {
				destTokens[i] = proto.Clone(token).(*craig_starsv1.ShipToken)
				destTokens[i].Quantity = 0
				destTokens[i].QuantityDamaged = 0
				destTokens[i].Damage = 0
			}
			if tt.allRight {
				sourceTokens, destTokens = destTokens, sourceTokens
			} else if tt.partial {
				sourceTokens[0].Quantity--
				destTokens[0].Quantity++
			}

			game, err := client.GetGame(ctx, fullGame.ID)
			require.NoError(t, err)
			ctx = context.WithValue(ctx, keyDbRead, dbConn.NewReadClient())
			ctx = context.WithValue(ctx, keyGame, game)
			ctx = context.WithValue(ctx, keyGamePlayer, &game.Players[0])
			response, err := NewFleetServiceHandler(dbConn).SplitFleet(ctx, connect.NewRequest(&craig_starsv1.SplitFleetRequest{
				GameId:         fullGame.ID,
				SourceFleetNum: int32(source.Num),
				DestFleetNum:   destNum,
				SourceTokens:   sourceTokens,
				DestTokens:     destTokens,
				DestBaseName:   tt.destBaseName,
			}))
			require.NoError(t, err)
			require.NotNil(t, response.Msg.Source)
			if tt.partial {
				require.NotNil(t, response.Msg.Dest)
			} else {
				assert.Nil(t, response.Msg.Dest)
			}

			fleetsAfter, err := client.GetFleetsForPlayer(ctx, fullGame.ID, source.PlayerNum)
			require.NoError(t, err)
			wantFleetCount := len(fleetsBefore)
			if tt.existingDest {
				wantFleetCount--
			} else if tt.partial {
				wantFleetCount++
			}
			assert.Len(t, fleetsAfter, wantFleetCount)
			shipsAfter := map[int]int{}
			cargoAfter := cs.Cargo{}
			fuelAfter := 0
			for _, fleet := range fleetsAfter {
				assert.NotEmpty(t, fleet.Tokens, "fleet %d must have ships", fleet.Num)
				for _, token := range fleet.Tokens {
					assert.Positive(t, token.Quantity)
					shipsAfter[token.DesignNum] += token.Quantity
				}
				cargoAfter = cargoAfter.Add(fleet.Cargo)
				fuelAfter += fleet.Fuel
			}
			assert.Equal(t, shipsBefore, shipsAfter)
			assert.Equal(t, cargoBefore, cargoAfter)
			assert.Equal(t, fuelBefore, fuelAfter)
			storedSource, err := client.GetFleetByNum(ctx, fullGame.ID, source.PlayerNum, int(response.Msg.Source.MapObject.Num))
			require.NoError(t, err)
			require.NotNil(t, storedSource)
			assert.Equal(t, response.Msg.Source.Tokens, converter.C.ConvertCSFleet(storedSource).Tokens)
		})
	}
}
