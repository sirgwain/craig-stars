package server

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/sirgwain/craig-stars/ai"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/stretchr/testify/require"
)

// Exercise the persisted game path with real fog of war. The default run is a
// short smoke test. AI_STRATEGY_ALL_MODES or AI_STRATEGY_REPORT runs the full
// audit: one race per PRT for thirty complete turns, optionally saving the
// observations. AI_STRATEGY_SEED overrides the map seed.
func TestAIStrategyThirtyTurns(t *testing.T) {
	seed, races, turns := int64(8032026), ai.Races[:10], 30
	if os.Getenv("AI_STRATEGY_ALL_MODES") == "" && os.Getenv("AI_STRATEGY_REPORT") == "" {
		races, turns = ai.Races[:4], 12
	}
	if s := os.Getenv("AI_STRATEGY_SEED"); s != "" {
		var err error
		seed, err = strconv.ParseInt(s, 10, 64)
		require.NoError(t, err)
	}
	var report strings.Builder
	fmt.Fprintln(&report, "difficulty,year,prt,planets,resources,tech,idle_colony_ships,loaded_idle_colony_fleets,queued_colony_ships,unexplored,armed_ships,miners,mine_layers,population,factories,mines")
	difficulties := []cs.AIDifficulty{cs.AIDifficultyNormal}
	if os.Getenv("AI_STRATEGY_ALL_MODES") != "" {
		difficulties = []cs.AIDifficulty{cs.AIDifficultyEasy, cs.AIDifficultyNormal, cs.AIDifficultyHard, cs.AIDifficultyCheater}
	}
	for _, difficulty := range difficulties {
		t.Run(string(difficulty), func(t *testing.T) {
			gr := createTestGameRunner(t.Context()).(*gameRunner)
			t.Cleanup(func() { gr.dbConn.Close() })
			rules := cs.NewRulesWithSeed(seed)
			settings := cs.NewGameSettings().WithName("AI strategy audit").WithSize(cs.SizeMedium)
			settings.RandomEvents = false
			settings.Rules = &rules
			for i, race := range races {
				settings.WithAIPlayerRace(race, difficulty, i%4)
			}
			game, err := gr.HostGame(1, settings)
			require.NoError(t, err)
			game.Seed = seed
			require.NoError(t, gr.dbConn.NewReadWriteClient().SaveGame(t.Context(), game.Game))
			for turn := 0; turn <= turns; turn++ {
				full, err := gr.loadGame(gr.dbConn.NewReadClient(), game.ID)
				require.NoError(t, err)
				require.NoError(t, gr.client.ComputeSpecs(full))
				for _, p := range full.Players {
					pmo := full.GetPlayerMapObjects(p.Num)
					// Surface AI errors rather than relying on the runner's logging.
					require.NoError(t, ai.NewAIPlayer(full.Game, full.TechStore, p, pmo).ProcessTurn())
					if turn%10 != 0 && turn != turns {
						continue
					}
					idle, loaded, queued, unknown := 0, 0, 0, 0
					armed, miners, layers, population, factories, mines := 0, 0, 0, 0, 0, 0
					for _, f := range pmo.Fleets {
						if f.Delete {
							continue
						}
						for _, token := range f.Tokens {
							if d := p.GetDesign(token.DesignNum); d != nil {
								if d.Spec.PowerRating > 0 {
									armed += token.Quantity
								}
								if d.Spec.MiningRate > 0 {
									miners += token.Quantity
								}
								if len(d.Spec.MineLayingRateByMineType) > 0 {
									layers += token.Quantity
								}
							}
						}
						if !f.Idle() || !f.Spec.Colonizer {
							continue
						}
						if f.Cargo.Colonists > 0 {
							loaded++
						}
						for _, token := range f.Tokens {
							if d := p.GetDesign(token.DesignNum); d != nil && d.Spec.Colonizer {
								idle += token.Quantity
							}
						}
					}
					for _, planet := range pmo.Planets {
						population += planet.GetPopulation()
						factories += planet.Factories
						mines += planet.Mines
						for _, item := range planet.ProductionQueue {
							if item.Type == cs.QueueItemTypeShipToken {
								if d := p.GetDesign(item.DesignNum); d != nil && d.Spec.Colonizer {
									queued += item.Quantity
								}
							}
						}
					}
					for _, intel := range p.PlanetIntels {
						if intel.Unexplored() {
							unknown++
						}
					}
					score := p.GetScore()
					fmt.Fprintf(&report, "%s,%d,%s,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d\n", difficulty, full.Year, p.Race.PRT, len(pmo.Planets), score.Resources, score.TechLevels, idle, loaded, queued, unknown, armed, miners, layers, population, factories, mines)
					if turn == turns {
						if turns >= 30 {
							require.Greater(t, len(pmo.Planets), 1, "%s should expand", p.Race.PRT)
						}
						require.LessOrEqual(t, idle+queued, 11, "%s must bound colony ship commitments", p.Race.PRT)
					}
				}
				if turn < turns {
					_, err = gr.GenerateTurn(game.ID)
					require.NoError(t, err)
				}
			}
		})
	}
	t.Log("\n" + report.String())
	if path := os.Getenv("AI_STRATEGY_REPORT"); path != "" {
		require.NoError(t, os.WriteFile(path, []byte(report.String()), 0644))
	}
}

func TestAIDifficultyRacePoolsAndOpeningOrders(t *testing.T) {
	gr := createTestGameRunner(t.Context()).(*gameRunner)
	t.Cleanup(func() { gr.dbConn.Close() })
	settings := cs.NewGameSettings()
	for _, difficulty := range []cs.AIDifficulty{cs.AIDifficultyEasy, cs.AIDifficultyNormal, cs.AIDifficultyHard, cs.AIDifficultyCheater} {
		settings.WithAIPlayer(difficulty, 0)
	}
	game, err := gr.HostGame(1, settings)
	require.NoError(t, err)
	for _, player := range game.Players {
		pool := ai.Races
		if player.AIDifficulty == cs.AIDifficultyCheater {
			pool = ai.CheaterRaces
		}
		found := false
		for _, race := range pool {
			if reflect.DeepEqual(raceSettings(player.Race), raceSettings(race)) {
				found = true
				break
			}
		}
		require.True(t, found, "%s must select its own race pool", player.AIDifficulty)
		// A race that starts without a scout, such as WM, orders one instead.
		scouting := false
		for _, f := range game.GetPlayerMapObjects(player.Num).Fleets {
			if f.GetTag(cs.TagPurpose) == string(cs.FleetPurposeScout) && len(f.Waypoints) > 1 {
				scouting = true
			}
		}
		for _, planet := range game.GetPlayerMapObjects(player.Num).Planets {
			for _, item := range planet.ProductionQueue {
				if item.GetTag(cs.TagPurpose) == string(cs.FleetPurposeScout) {
					scouting = true
				}
			}
		}
		require.True(t, scouting, "%s %s should scout or order a scout before the first turn", player.AIDifficulty, player.Race.PRT)
	}
}

func raceSettings(race cs.Race) cs.Race {
	race.DBObject = cs.DBObject{}
	race.UserID = 0
	race.Name, race.PluralName = "", ""
	race.Spec = cs.RaceSpec{}
	return race
}
