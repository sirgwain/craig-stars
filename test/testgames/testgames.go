//go:build !wasi && !wasm

package testgames

import (
	"context"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
)

var TestGames = []cs.TestScenario{
	cs.SingleUnitScenario(),
	cs.TwoPlayerScenario(),
	cs.ScenarioScoutTest(),
	cs.ScenarioColonizerTest(),
	cs.ScenarioColonizerTestAR(),
	cs.ScenarioKitchenSink(),
	cs.ScenarioSDMinefield(),
	cs.ScenarioCargoTransferInvasion(),
	cs.ScenarioCargoTransferInvasionStarbase(),
	cs.ScenarioCargoTransferPlanetOwned(),
	cs.ScenarioCargoTransferPlanetSteal(),
	cs.ScenarioCargoTransferFleetSteal(),
	cs.ScenarioCargoTransferJettison(),
	cs.ScenarioCargoTransferSalvage(),
	cs.ScenarioCargoTransferFleets(),
	cs.ScenarioCargoTransferSplit(),
	cs.ScenarioCargoTransferMineralPacket(),
	cs.ScenarioMysteryTrader(),
	cs.ScenarioBattle1(),
	cs.ScenarioProductionStarbases(),
}

// CreateTestGames creates one of each test game for manual testing
// for automated testing they should be created new each time
func CreateTestGames(db db.Client) error {
	ctx := context.Background()
	for _, testGame := range TestGames {
		var err error
		game := cs.BuildScenario(testGame)
		err = db.SaveGame(ctx, game.Game)
		if err != nil {
			return err
		}
		for _, player := range game.Players {
			player.GameID = game.ID

			if err := db.SavePlayer(ctx, player); err != nil {
				return err
			}
			for _, design := range player.Designs {
				design.GameID = game.ID
			}
		}

		// save to db
		if err := db.UpdateFullGame(ctx, game); err != nil {
			return err
		}
	}

	return nil
}
