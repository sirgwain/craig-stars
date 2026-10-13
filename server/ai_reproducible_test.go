package server

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/sirgwain/craig-stars/ai"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/stretchr/testify/require"
)

// Play the same AI game twice, a year at a time, and stop at the first year
// the two differ. The regression recordings only work if a game always plays
// out the same way, and when one doesn't this shows where it went wrong.
//
// AI_REPRODUCIBLE_RACE picks the race by its primary trait, HE by default.
func TestAIGamesAreReproducible(t *testing.T) {
	race := ai.Races[0]
	if prt := os.Getenv("AI_REPRODUCIBLE_RACE"); prt != "" {
		index := slices.IndexFunc(ai.Races, func(r cs.Race) bool { return string(r.PRT) == prt })
		require.NotEqual(t, -1, index, "no AI race for %s", prt)
		race = ai.Races[index]
	}

	host := func() (*gameRunner, int64) {
		gr := createTestGameRunner(t.Context()).(*gameRunner)
		t.Cleanup(func() { gr.dbConn.Close() })
		game, err := gr.hostAIGame(race, aiGameOpponent(), aiGameSeed, cs.SizeSmall)
		require.NoError(t, err)
		return gr, game.ID
	}
	first, firstID := host()
	second, secondID := host()

	for turn := 0; turn <= aiGameTurns; turn++ {
		a, b := describeAIGame(t, first, firstID), describeAIGame(t, second, secondID)
		differences := []string{}
		for i := 0; i < max(len(a), len(b)) && len(differences) < 12; i++ {
			var x, y string
			if i < len(a) {
				x = a[i]
			}
			if i < len(b) {
				y = b[i]
			}
			if x != y {
				differences = append(differences, "first:  "+x, "second: "+y)
			}
		}
		require.Empty(t, differences, "the two games differ after %d turns:\n%s", turn, strings.Join(differences, "\n"))

		_, err := first.GenerateTurn(firstID)
		require.NoError(t, err)
		_, err = second.GenerateTurn(secondID)
		require.NoError(t, err)
	}
}

// describeAIGame lists everything in a game that should match between two
// plays of it, one line per object, in a fixed order
func describeAIGame(t *testing.T, gr *gameRunner, gameID int64) []string {
	full, err := gr.loadGame(gr.dbConn.NewReadClient(), gameID)
	require.NoError(t, err)

	lines := []string{}
	for _, p := range full.Players {
		lines = append(lines, fmt.Sprintf("player %d researching %v tech %+v spent %+v designs %d", p.Num, p.Researching, p.TechLevels, p.TechLevelsSpent, len(p.Designs)))
		for _, intel := range p.PlanetIntels {
			lines = append(lines, fmt.Sprintf("player %d intel planet %d age %d owner %d colonists %d", p.Num, intel.Num, intel.ReportAge, intel.PlayerNum, intel.Cargo.Colonists))
		}
	}
	planets := []string{}
	for _, planet := range full.Planets {
		queue := []string{}
		for _, item := range planet.ProductionQueue {
			queue = append(queue, fmt.Sprintf("%s:%d:%d:%+v", item.Type, item.DesignNum, item.Quantity, item.Allocated))
		}
		planets = append(planets, fmt.Sprintf("planet %04d owner %d cargo %+v mines %d factories %d defenses %d hab %+v concentration %+v mine years %+v queue %s", planet.Num, planet.PlayerNum, planet.Cargo, planet.Mines, planet.Factories, planet.Defenses, planet.Hab, planet.MineralConcentration, planet.MineYears, strings.Join(queue, " ")))
	}
	slices.Sort(planets)
	fleets := []string{}
	for _, fleet := range full.Fleets {
		tokens, waypoints := []string{}, []string{}
		for _, token := range fleet.Tokens {
			tokens = append(tokens, fmt.Sprintf("%dx%d damage %v/%d", token.DesignNum, token.Quantity, token.Damage, token.QuantityDamaged))
		}
		for _, wp := range fleet.Waypoints {
			waypoints = append(waypoints, fmt.Sprintf("%v/%d/%s/w%d", wp.Position, wp.TargetNum, wp.Task, wp.WarpSpeed))
		}
		fleets = append(fleets, fmt.Sprintf("fleet %d/%04d at %v cargo %+v fuel %d tags %v ships %s orders %s", fleet.PlayerNum, fleet.Num, fleet.Position, fleet.Cargo, fleet.Fuel, fleet.Tags, strings.Join(tokens, ";"), strings.Join(waypoints, " ")))
	}
	slices.Sort(fleets)
	others := []string{}
	for _, minefield := range full.Minefields {
		others = append(others, fmt.Sprintf("minefield %d/%d at %v mines %d", minefield.PlayerNum, minefield.Num, minefield.Position, minefield.NumMines))
	}
	for _, packet := range full.MineralPackets {
		others = append(others, fmt.Sprintf("packet %d/%d at %v cargo %+v", packet.PlayerNum, packet.Num, packet.Position, packet.Cargo))
	}
	for _, salvage := range full.Salvages {
		others = append(others, fmt.Sprintf("salvage %d/%d at %v cargo %+v", salvage.PlayerNum, salvage.Num, salvage.Position, salvage.Cargo))
	}
	slices.Sort(others)
	return slices.Concat(lines, planets, fleets, others)
}
