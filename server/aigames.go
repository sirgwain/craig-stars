package server

import (
	"fmt"
	"log/slog"

	"github.com/sirgwain/craig-stars/ai"
	configpkg "github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/cs"
)

// The AI games played by the regression test and by `serve --ai-mode`: each
// standard AI race against a single JoaT opponent, on a fixed map.
const (
	aiGameSeed  = 20261012
	aiGameTurns = 80
)

// the JoaT every other AI race is matched against
func aiGameOpponent() cs.Race {
	return ai.Races[len(ai.Races)-1]
}

// hostAIGame hosts a two player game of one AI race against the opponent
func (gr *gameRunner) hostAIGame(race, opponent cs.Race, seed int64, size cs.Size) (*cs.FullGame, error) {
	rules := cs.NewRulesWithSeed(seed)
	settings := cs.NewGameSettings().WithName(fmt.Sprintf("AI %s vs %s", race.PRT, opponent.PRT)).WithSize(size)
	settings.RandomEvents = false
	settings.Rules = &rules
	settings.WithAIPlayerRace(race, cs.AIDifficultyNormal, 0)
	settings.WithAIPlayerRace(opponent, cs.AIDifficultyNormal, 0)
	// hosted by the admin user every database starts with
	game, err := gr.HostGame(1, settings)
	if err != nil {
		return nil, err
	}
	game.Seed = seed
	if err := gr.dbConn.NewReadWriteClient().SaveGame(gr.ctx, game.Game); err != nil {
		return nil, err
	}
	return game, nil
}

// createAIGames plays every AI game to its last turn, leaving the games in
// the database for an admin to open and view as either player
func createAIGames(dbConn DBConnection, config configpkg.Config) error {
	gr := NewGameRunner(dbConn, config).(*gameRunner)
	opponent := aiGameOpponent()
	for _, race := range ai.Races {
		game, err := gr.hostAIGame(race, opponent, aiGameSeed, cs.SizeSmall)
		if err != nil {
			return fmt.Errorf("host ai game for %s: %w", race.PRT, err)
		}
		for range aiGameTurns {
			if _, err := gr.GenerateTurn(game.ID); err != nil {
				return fmt.Errorf("generate turn for %s: %w", game.Name, err)
			}
		}
		slog.Info("Created AI game", slog.String("game", game.Name), slog.Int64("GameID", game.ID))
	}
	return nil
}
