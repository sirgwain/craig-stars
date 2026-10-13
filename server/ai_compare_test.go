package server

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/sirgwain/craig-stars/ai"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/stretchr/testify/require"
)

// Compare the AI with and without some behaviors. For every map and AI race,
// two games are played against the JoaT opponent: one where only the race
// has the behaviors and one where only the opponent does. The side with them
// should score higher if they help. Score rewards spreading population over
// many planets, so income is reported as well.
//
//	AI_COMPARE=behavior[,behavior]   the ai.Experiment behaviors to compare (required)
//	AI_COMPARE_SEEDS=1,2,3           map seeds, four by default
//	AI_COMPARE_TURNS=80              years to play
//	AI_COMPARE_SIZE=Small            universe size
//	AI_COMPARE_RACES=PP,IT           only these races, all ten by default
func TestAICompare(t *testing.T) {
	behaviors := strings.Split(os.Getenv("AI_COMPARE"), ",")
	if behaviors[0] == "" {
		t.Skip("set AI_COMPARE to the behaviors to compare")
	}
	seeds := []int64{aiGameSeed, 777, 31337, 9001}
	if v := os.Getenv("AI_COMPARE_SEEDS"); v != "" {
		seeds = nil
		for _, s := range strings.Split(v, ",") {
			seed, err := strconv.ParseInt(s, 10, 64)
			require.NoError(t, err)
			seeds = append(seeds, seed)
		}
	}
	turns := 80
	if v := os.Getenv("AI_COMPARE_TURNS"); v != "" {
		var err error
		turns, err = strconv.Atoi(v)
		require.NoError(t, err)
	}
	size := cs.SizeSmall
	if v := os.Getenv("AI_COMPARE_SIZE"); v != "" {
		size = cs.Size(v)
	}

	// play a game where one player goes without the behaviors, returning both scores
	without := 0
	ai.Experiment = func(playerNum int, behavior string) bool {
		return playerNum == without && slices.Contains(behaviors, behavior)
	}
	t.Cleanup(func() { ai.Experiment = nil })
	play := func(race cs.Race, seed int64, playerWithout int) (raceScore, opponentScore cs.PlayerScore) {
		without = playerWithout
		gr := createTestGameRunner(t.Context()).(*gameRunner)
		defer gr.dbConn.Close()
		game, err := gr.hostAIGame(race, aiGameOpponent(), seed, size)
		require.NoError(t, err)
		for range turns {
			_, err := gr.GenerateTurn(game.ID)
			require.NoError(t, err)
		}
		full, err := gr.loadGame(gr.dbConn.NewReadClient(), game.ID)
		require.NoError(t, err)
		return full.Players[0].GetScore(), full.Players[1].GetScore()
	}

	var report strings.Builder
	fmt.Fprintf(&report, "\nscore with / without %s after %d turns (%s)\n", strings.Join(behaviors, ","), turns, size)
	totalScore, totalIncome := 0.0, 0.0
	for _, seed := range seeds {
		with, wo, incomeWith, incomeWithout := 0, 0, 0, 0
		wins := []string{}
		for _, race := range ai.Races {
			if only := os.Getenv("AI_COMPARE_RACES"); only != "" && !slices.Contains(strings.Split(only, ","), string(race.PRT)) {
				continue
			}
			raceWith, opponentWithout := play(race, seed, 2)
			raceWithout, opponentWith := play(race, seed, 1)
			with += raceWith.Score + opponentWith.Score
			wo += raceWithout.Score + opponentWithout.Score
			incomeWith += raceWith.Resources + opponentWith.Resources
			incomeWithout += raceWithout.Resources + opponentWithout.Resources
			wins = append(wins, fmt.Sprintf("%s %+d%%", race.PRT, 100*(raceWith.Score-raceWithout.Score)/max(1, raceWithout.Score)))
		}
		change := 100*float64(with)/float64(wo) - 100
		income := 100*float64(incomeWith)/float64(max(1, incomeWithout)) - 100
		totalScore += change
		totalIncome += income
		fmt.Fprintf(&report, "  seed %-9d score %6d / %6d %+5.1f%%  income %+5.1f%%   %s\n", seed, with, wo, change, income, strings.Join(wins, " "))
	}
	fmt.Fprintf(&report, "  average score %+.2f%%, income %+.2f%% over %d maps\n", totalScore/float64(len(seeds)), totalIncome/float64(len(seeds)), len(seeds))
	t.Log(report.String())
	if path := os.Getenv("AI_COMPARE_REPORT"); path != "" {
		require.NoError(t, os.WriteFile(path, []byte(report.String()), 0644))
	}
}
