package cmd

import (
	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/db"

	"github.com/spf13/cobra"
)

// deleteCmd is the root delete command for delete operations
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete various resources",
	Long:  `Delete various resources`,
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	addDeleteUserCmd()
	addDeleteGameCmd()
}

func addDeleteUserCmd() {
	var id int64

	// deleteUserCmd deletes a user from the database
	var deleteUserCmd = &cobra.Command{
		Use:   "user",
		Short: "Delete user",
		Long:  `Delete user from the database`,
		RunE: func(cmd *cobra.Command, args []string) error {
			db := db.NewClient()
			cfg := config.GetConfig()
			db.Connect(cfg)

			db.DeleteUser(id)
			users, err := db.GetUsers()
			if err != nil {
				return err
			}

			PrintTable("Remaining Users", users)
			return nil
		},
	}
	deleteUserCmd.Flags().Int64VarP(&id, "user-id", "u", 0, "Delete games for user id")
	deleteUserCmd.MarkFlagRequired("user-id")

	deleteCmd.AddCommand(deleteUserCmd)

}
func addDeleteGameCmd() {
	var id int64

	// deleteGameCmd deletes a game from the database
	var deleteGameCmd = &cobra.Command{
		Use:   "game",
		Short: "Delete game",
		Long:  `Delete game from the database`,
		RunE: func(cmd *cobra.Command, args []string) error {
			db := db.NewClient()
			cfg := config.GetConfig()
			db.Connect(cfg)

			db.DeleteGame(id)
			games, err := db.GetGames()
			if err != nil {
				return err
			}

			PrintTable("Remaining Games", games)
			return nil
		},
	}

	deleteGameCmd.Flags().Int64VarP(&id, "game-id", "g", 0, "Delete game by id")
	deleteGameCmd.MarkFlagRequired("game-id")
	deleteCmd.AddCommand(deleteGameCmd)
}
