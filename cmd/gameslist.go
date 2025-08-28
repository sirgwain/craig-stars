package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/lensesio/tableprinter"
	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
	"github.com/spf13/cobra"
	"log/slog"
)

func newGamesListCmd() *cobra.Command {
	var full bool
	var save bool
	var generate bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list games",
		Long:  `List games in database.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			cfg := config.GetConfig()

			// create a new connection to the database
			dbConn := db.NewConn()
			if err := dbConn.Connect(cfg); err != nil {
				return err
			}
			defer func() { dbConn.Close() }()
			db := dbConn.NewReadWriteClient()

			games, err := db.GetGamesWithPlayers(ctx)
			if err != nil {
				return fmt.Errorf("failed to load games %w", err)
			}

			type game struct {
				Name    string `header:"Name"`
				Players int    `header:"Players"`
			}

			printer := tableprinter.New(os.Stdout)
			tableGames := make([]game, 0, len(games))

			if !full {
				for _, g := range games {
					tableGames = append(tableGames, game{Name: g.Name, Players: len(g.Players)})
				}
				printer.Print(tableGames)
				return nil
			}

			fmt.Printf("Loading full games and players\n\n")
			for _, g := range games {
				fullGame, err := db.GetFullGame(ctx, g.ID)
				if err != nil {
					return fmt.Errorf("failed to load full game %s (%d): %w", g.Name, g.ID, err)
				}
				fmt.Printf("loaded %s\n", fullGame.Name)
				tableGames = append(tableGames, game{Name: fullGame.Name, Players: len(fullGame.Players)})

				if save {
					for _, p := range fullGame.Planets {
						p.MarkDirty()
					}
					db.UpdateFullGame(ctx, fullGame)
					fmt.Printf("saved %s\n", fullGame.Name)
				}
				if generate {
					// Set log level to warn for turn generation
					warnLogger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn}))
					slog.SetDefault(warnLogger)
					client := cs.NewGamer()
					if err := client.GenerateTurn(fullGame.Game, fullGame.Universe, fullGame.Players); err != nil {
						return fmt.Errorf("failed to generate turn for game %s (%d): %w", fullGame.Name, fullGame.ID, err)
					}
					fmt.Printf("generated turn %s\n", fullGame.Name)
				}
			}
			printer.Print(tableGames)
			return nil
		},
	}
	cmd.Flags().BoolVar(&full, "full", false, "Load the full game")
	cmd.Flags().BoolVar(&save, "save", false, "Save the full game")
	cmd.Flags().BoolVar(&generate, "generate", false, "Generate a turn (without saving)")

	return cmd
}
