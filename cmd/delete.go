package cmd

import (
	"github.com/sirgwain/craig-stars/appcontext"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
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
	var id uint

	// deleteUserCmd represents the deleteUsers command
	var deleteUserCmd = &cobra.Command{
		Use:   "user",
		Short: "Delete user",
		Long:  `Delete user from the database`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := appcontext.Initialize()
			ctx.DB.DeleteUserById(id)
			PrintTable("Remaining Users", *ctx.DB.GetUsers())
		},
	}
	deleteUserCmd.Flags().UintVarP(&id, "user-id", "u", 0, "Delete games for user id")
	deleteUserCmd.MarkFlagRequired("user-id")

	deleteCmd.AddCommand(deleteUserCmd)

}
func addDeleteGameCmd() {
	var id uint

	// deleteUsersCmd represents the deleteUsers command
	var deleteGameCmd = &cobra.Command{
		Use:   "game",
		Short: "Delete game",
		Long:  `Delete game from the database`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := appcontext.Initialize()
			ctx.DB.DeleteGameById(id)
			PrintTable("Remaining Games", ctx.DB.GetGames())
		},
	}

	deleteGameCmd.Flags().UintVarP(&id, "game-id", "g", 0, "Delete game by id")
	deleteGameCmd.MarkFlagRequired("game-id")
	deleteCmd.AddCommand(deleteGameCmd)
}
