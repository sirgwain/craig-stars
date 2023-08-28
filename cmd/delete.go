package cmd

import (
	"github.com/spf13/cobra"
)

// deleteCmd is the root delete command for delete operations
var deleteCmd = &cobra.Command{
	Use:                "delete",
	Short:              "Delete various resources",
	Long:               `Delete various resources`,
	PersistentPreRunE:  dbPreRun,
	PersistentPostRunE: dbPostRun,
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
			client.DeleteUser(id)
			users, err := client.GetUsers()
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
			client.DeleteGame(id)
			games, err := client.GetGames()
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
