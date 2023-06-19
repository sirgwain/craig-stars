package cmd

import (
	"github.com/sirgwain/craig-stars/appcontext"
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
	Run: func(cmd *cobra.Command, args []string) {
		ctx := appcontext.Initialize()
		PrintTable("Users", *ctx.DB.GetUsers())
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(listUsersCmd)

	addListGamesCmd()
}

func addListGamesCmd() {
	var userID uint

	// listUsersCmd represents the listUsers command
	var listGamesCmd = &cobra.Command{
		Use:   "games",
		Short: "List games",
		Long:  `List games in the database`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := appcontext.Initialize()
			var games []game.Game
			if userID != 0 {
				games = ctx.DB.GetGamesByUser(userID)
			} else {
				games = ctx.DB.GetGames()
			}
			PrintTable("Games", games)
		},
	}

	listGamesCmd.Flags().UintVarP(&userID, "user-id", "u", 0, "List games for user id")
	listCmd.AddCommand(listGamesCmd)

}
