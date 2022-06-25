package cmd

import (
	"github.com/sirgwain/craig-stars/appcontext"
	"github.com/sirgwain/craig-stars/game"

	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a game or turn",
	Long:  `Generate a game or turn.`,
}

func init() {
	rootCmd.AddCommand(generateCmd)
	addGenerateGameCommand()
}

func addGenerateGameCommand() {

	var name string
	var userID uint

	// generateGameCmd represents the generateGame command
	generateGameCmd := &cobra.Command{
		Use:   "game",
		Short: "A brief description of your command",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := appcontext.Initialize()

			g := game.NewGame()
			g.Name = name

			// save the game to the db
			if userID != 0 {
				g.AddPlayer(game.NewPlayer(userID, game.NewRace()))
			}

			if err := ctx.DB.CreateGame(g); err != nil {
				return err
			}

			if err := g.GenerateUniverse(); err != nil {
				return err
			}

			if err := ctx.DB.SaveGame(g); err != nil {
				return err
			}
			return nil
		},
	}

	generateGameCmd.Flags().StringVarP(&name, "name", "n", "A Barefoot Jaywalk", "The name of the game")
	generateGameCmd.Flags().UintVarP(&userID, "user-id", "u", 0, "A user to create a player for")
	generateGameCmd.MarkFlagRequired("gamename")
	generateGameCmd.MarkFlagRequired("password")

	generateCmd.AddCommand(generateGameCmd)
}
