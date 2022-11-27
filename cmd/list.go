package cmd

import (
	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/db"
	"github.com/sirgwain/craig-stars/game"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List various resources",
	Long:  `List various resources`,
}

// listUsersCmd represents the listUsers command
var listUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "List users",
	Long:  `List users in the database`,
	RunE: func(cmd *cobra.Command, args []string) error {
		db := db.NewClient()
		cfg := config.GetConfig()
		db.Connect(cfg)

		users, err := db.GetUsers()
		if err != nil {
			return err
		}
		PrintTable("Users", users)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(listUsersCmd)

	addListGamesCmd()
}

func addListGamesCmd() {
	var userID int64

	// listUsersCmd represents the listUsers command
	var listGamesCmd = &cobra.Command{
		Use:   "games",
		Short: "List games",
		Long:  `List games in the database`,
		RunE: func(cmd *cobra.Command, args []string) error {
			db := db.NewClient()
			cfg := config.GetConfig()
			db.Connect(cfg)

			var games []game.Game
			var err error
			if userID != 0 {
				games, err = db.GetGamesForUser(userID)
			} else {
				games, err = db.GetGames()
			}

			if err != nil {
				return err
			}

			PrintTable("Games", games)
			return nil
		},
	}

	listGamesCmd.Flags().Int64VarP(&userID, "user-id", "u", 0, "List games for user id")
	listCmd.AddCommand(listGamesCmd)

}
