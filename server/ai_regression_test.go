package server

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/sirgwain/craig-stars/ai"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/stretchr/testify/require"
)

const (
	aiRegressionEvery = 10
	aiRegressionDir   = "ai-regression"
)

// ship categories recorded for each snapshot
var aiRegressionShips = []string{"scout", "colonizer", "freighter", "miner", "mine_layer", "bomber", "warship", "other"}

// Play each PRT's AI race against a single JoaT opponent and compare a set of
// metrics, recorded every ten years, against the file committed for that PRT.
// A difference means AI behavior changed. If the change is intended, record it
// with AI_REGRESSION_UPDATE=1 and review the diff.
//
// `serve --ai-mode` plays the same games into a database to view in the app.
func TestAIRegression(t *testing.T) {
	opponent := aiGameOpponent()
	// AI_REGRESSION_WITHOUT plays without some behaviors, to see what they change in the recordings
	if without := strings.Split(os.Getenv("AI_REGRESSION_WITHOUT"), ","); without[0] != "" {
		ai.Experiment = func(_ int, behavior string) bool { return slices.Contains(without, behavior) }
		t.Cleanup(func() { ai.Experiment = nil })
	}
	require.Equal(t, cs.JoaT, opponent.PRT)

	for _, race := range ai.Races {
		t.Run(string(race.PRT), func(t *testing.T) {
			actual := playAIRegressionGame(t, race, opponent)
			path := filepath.Join(aiRegressionDir, fmt.Sprintf("%s.csv", race.PRT))
			if os.Getenv("AI_REGRESSION_UPDATE") != "" {
				require.NoError(t, os.MkdirAll(aiRegressionDir, 0755))
				require.NoError(t, os.WriteFile(path, []byte(actual.String()), 0644))
				return
			}
			recorded, err := os.ReadFile(path)
			require.NoError(t, err, "no recording, run with AI_REGRESSION_UPDATE=1")

			// report each metric that changed, rather than two whole files
			expected := map[string]string{}
			for _, line := range strings.Split(strings.TrimSpace(string(recorded)), "\n") {
				metric, values, _ := strings.Cut(line, ",")
				expected[metric] = values
			}
			changes := []string{}
			for _, metric := range actual.metrics {
				values := strings.Join(actual.values[metric], ",")
				if recordedValues, found := expected[metric]; !found {
					changes = append(changes, fmt.Sprintf("%s: not recorded, now %s", metric, values))
				} else if recordedValues != values {
					changes = append(changes, fmt.Sprintf("%s: recorded %s, now %s", metric, recordedValues, values))
				}
				delete(expected, metric)
			}
			delete(expected, "metric")
			for metric := range expected {
				changes = append(changes, fmt.Sprintf("%s: recorded but no longer measured", metric))
			}
			if len(changes) > 0 {
				t.Errorf("%s AI behavior changed from %s (AI_REGRESSION_UPDATE=1 records the new behavior):\n  %s", race.PRT, path, strings.Join(changes, "\n  "))
			}
		})
	}
}

// metrics in the order they were first recorded, with a value per snapshot
type aiRegressionRecording struct {
	years   []string
	metrics []string
	values  map[string][]string
}

func (r *aiRegressionRecording) set(metric string, value int) {
	if _, found := r.values[metric]; !found {
		r.metrics = append(r.metrics, metric)
	}
	r.values[metric] = append(r.values[metric], strconv.Itoa(value))
}

// one row per metric and one column per snapshot, so a diff shows what changed and when
func (r *aiRegressionRecording) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "metric,%s\n", strings.Join(r.years, ","))
	for _, metric := range r.metrics {
		fmt.Fprintf(&b, "%s,%s\n", metric, strings.Join(r.values[metric], ","))
	}
	return b.String()
}

func playAIRegressionGame(t *testing.T, race, opponent cs.Race) *aiRegressionRecording {
	gr := createTestGameRunner(t.Context()).(*gameRunner)
	t.Cleanup(func() { gr.dbConn.Close() })
	game, err := gr.hostAIGame(race, opponent, aiGameSeed, cs.SizeSmall)
	require.NoError(t, err)

	recording := &aiRegressionRecording{values: map[string][]string{}}
	// the first year each of these happened, or 0 if it never did
	milestones := map[string]int{}
	milestone := func(name string, year int, reached bool) {
		if reached && milestones[name] == 0 {
			milestones[name] = year
		}
	}

	startingWarships := 0
	// who owned each planet last year, to count worlds changing hands
	owners := map[int]int{}
	taken, lost, emptied, battles := 0, 0, 0, 0
	for turn := 0; turn <= aiGameTurns; turn++ {
		full, err := gr.loadGame(gr.dbConn.NewReadClient(), game.ID)
		require.NoError(t, err)
		require.NoError(t, gr.client.ComputeSpecs(full))
		for _, planet := range full.Planets {
			if before, seen := owners[planet.Num]; seen && before != planet.PlayerNum {
				switch {
				case before == 2 && planet.PlayerNum == 1:
					taken++
				case before == 1 && planet.PlayerNum == 2:
					lost++
				case before == 2 && !planet.Owned():
					emptied++
				}
			}
			owners[planet.Num] = planet.PlayerNum
		}
		battles += len(full.Players[0].BattleRecords)
		for _, p := range full.Players {
			pmo := full.GetPlayerMapObjects(p.Num)
			// Surface AI errors rather than relying on the runner's logging.
			require.NoError(t, ai.NewAIPlayer(full.Game, full.TechStore, p, pmo).ProcessTurn(), "%s in %d", p.Race.PRT, full.Year)
			if p.Num != 1 {
				continue
			}

			ships, idle, queued := map[string]int{}, map[string]int{}, map[string]int{}
			idleLoadedColonizers := 0
			for _, f := range pmo.Fleets {
				if f.Delete {
					continue
				}
				if f.Idle() && f.Spec.Colonizer && f.Cargo.Colonists > 0 {
					idleLoadedColonizers++
				}
				for _, token := range f.Tokens {
					category := aiRegressionShipCategory(p.GetDesign(token.DesignNum))
					ships[category] += token.Quantity
					if f.Idle() {
						idle[category] += token.Quantity
					}
				}
			}
			starbases, population, factories, mines, minerals := 0, 0, 0, 0, 0
			stargates, massDrivers, colonyPopulation, coloniesAtReserve := 0, 0, 0, 0
			defenses, homeworldResources, homeworldPopulation, homeworldOrdered, homeworldMinerals := 0, 0, 0, 0, 0
			for _, planet := range pmo.Planets {
				defenses += planet.Defenses
				if !planet.Homeworld {
					colonyPopulation += planet.GetPopulation()
					if planet.Spec.PopulationDensity >= .25 {
						coloniesAtReserve++
					}
				}
				if planet.Spec.HasStargate {
					stargates++
				}
				if planet.Spec.HasMassDriver {
					massDrivers++
				}
				if planet.Homeworld {
					// resources still to be paid on the ships and starbases it has on order
					homeworldResources += planet.Spec.ResourcesPerYear
					homeworldPopulation += planet.GetPopulation()
					homeworldMinerals += planet.Cargo.ToMineral().Total()
					for _, item := range planet.ProductionQueue {
						if design := p.GetDesign(item.DesignNum); design != nil && (item.Type == cs.QueueItemTypeShipToken || item.Type == cs.QueueItemTypeStarbase) {
							homeworldOrdered += item.Quantity*design.Spec.Cost.Resources - item.Allocated.Resources
						}
					}
				}
				if planet.Spec.HasStarbase {
					starbases++
				}
				population += planet.GetPopulation()
				factories += planet.Factories
				mines += planet.Mines
				minerals += planet.Cargo.ToMineral().Total()
				for _, item := range planet.ProductionQueue {
					if item.Type == cs.QueueItemTypeShipToken {
						queued[aiRegressionShipCategory(p.GetDesign(item.DesignNum))] += item.Quantity
					}
				}
			}
			unexplored, enemyPlanetsKnown := 0, 0
			for _, intel := range p.PlanetIntels {
				if intel.Unexplored() {
					unexplored++
				} else if intel.Owned() && intel.PlayerNum != p.Num {
					enemyPlanetsKnown++
				}
			}

			if turn == 0 {
				startingWarships = ships["warship"] + ships["bomber"]
			}
			milestone("first_colony_year", full.Year, len(pmo.Planets) > 1)
			milestone("second_starbase_year", full.Year, starbases > 1)
			milestone("first_warship_built_year", full.Year, ships["warship"]+ships["bomber"] > startingWarships)
			milestone("met_opponent_year", full.Year, enemyPlanetsKnown > 0)
			attackingFleets := 0
			for _, f := range pmo.Fleets {
				if len(f.Waypoints) < 2 || f.GetTag(cs.TagPurpose) != string(cs.FleetPurposeBomber) {
					continue
				}
				if target := p.GetPlanetIntel(f.Waypoints[len(f.Waypoints)-1].TargetNum); target != nil && target.Owned() && target.PlayerNum != p.Num {
					attackingFleets++
				}
			}
			milestone("first_attack_year", full.Year, attackingFleets > 0)
			// a loaded colony ship with nowhere to go is always a bug
			require.Zero(t, idleLoadedColonizers, "%s has an idle loaded colony ship in %d", p.Race.PRT, full.Year)

			if turn == 0 || turn%aiRegressionEvery != 0 {
				continue
			}
			recording.years = append(recording.years, strconv.Itoa(full.Year))
			score := p.GetScore()
			recording.set("score", score.Score)
			recording.set("opponent_score", full.Players[1].GetScore().Score)
			recording.set("opponent_planets", len(full.GetPlayerMapObjects(2).Planets))
			recording.set("planets", len(pmo.Planets))
			recording.set("starbases", starbases)
			recording.set("population", population)
			recording.set("resources", score.Resources)
			recording.set("factories", factories)
			recording.set("mines", mines)
			recording.set("surface_minerals", minerals)
			recording.set("tech_levels", score.TechLevels)
			recording.set("tech_energy", p.TechLevels.Energy)
			recording.set("tech_weapons", p.TechLevels.Weapons)
			recording.set("tech_propulsion", p.TechLevels.Propulsion)
			recording.set("tech_construction", p.TechLevels.Construction)
			recording.set("tech_electronics", p.TechLevels.Electronics)
			recording.set("tech_biotechnology", p.TechLevels.Biotechnology)
			recording.set("colony_population", colonyPopulation)
			recording.set("colonies_at_growth_reserve", coloniesAtReserve)
			recording.set("homeworld_population", homeworldPopulation)
			recording.set("homeworld_resources", homeworldResources)
			recording.set("homeworld_resources_on_order", homeworldOrdered)
			recording.set("homeworld_surface_minerals", homeworldMinerals)
			recording.set("defenses", defenses)
			recording.set("minefields", len(pmo.Minefields))
			recording.set("stargates", stargates)
			recording.set("mass_drivers", massDrivers)
			recording.set("packets_in_flight", len(pmo.MineralPackets))
			recording.set("unexplored_planets", unexplored)
			recording.set("enemy_planets_known", enemyPlanetsKnown)
			recording.set("attacking_fleets", attackingFleets)
			recording.set("battles", battles)
			recording.set("enemy_planets_taken", taken)
			recording.set("enemy_planets_emptied", emptied)
			recording.set("planets_lost", lost)
			for _, category := range aiRegressionShips {
				recording.set("ships_"+category, ships[category])
			}
			for _, category := range aiRegressionShips {
				recording.set("idle_"+category, idle[category])
			}
			for _, category := range aiRegressionShips {
				recording.set("queued_"+category, queued[category])
			}
		}
		if turn < aiGameTurns {
			_, err = gr.GenerateTurn(game.ID)
			require.NoError(t, err)
		}
	}

	// an AI that never expands has stopped working, whatever was recorded
	require.NotZero(t, milestones["first_colony_year"], "%s never colonized a planet", race.PRT)

	for _, name := range []string{"first_colony_year", "second_starbase_year", "first_warship_built_year", "met_opponent_year", "first_attack_year"} {
		recording.metrics = append(recording.metrics, name)
		recording.values[name] = []string{strconv.Itoa(milestones[name])}
	}
	return recording
}

func aiRegressionShipCategory(design *cs.ShipDesign) string {
	if design == nil {
		return "other"
	}
	switch {
	case design.Spec.Colonizer:
		return "colonizer"
	case design.Spec.Bomber:
		return "bomber"
	case design.Spec.PowerRating > 0:
		return "warship"
	case design.Spec.MiningRate > 0:
		return "miner"
	case len(design.Spec.MineLayingRateByMineType) > 0:
		return "mine_layer"
	case design.Purpose == cs.ShipDesignPurposeScout:
		return "scout"
	case design.Spec.CargoCapacity > 0 || design.Purpose == cs.ShipDesignPurposeFuelFreighter:
		return "freighter"
	}
	return "other"
}
