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
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := appcontext.Initialize()
		users, err := ctx.DB.GetUsers()
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
	var userID uint64

	// listUsersCmd represents the listUsers command
	var listGamesCmd = &cobra.Command{
		Use:   "games",
		Short: "List games",
		Long:  `List games in the database`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := appcontext.Initialize()
			var games []*game.Game
			var err error
			if userID != 0 {
				games, err = ctx.DB.GetGamesByUser(userID)
			} else {
				games, err = ctx.DB.GetGames()
			}

			if err != nil {
				return err
			}

			PrintTable("Games", games)
			return nil
		},
	}

	listGamesCmd.Flags().Uint64VarP(&userID, "user-id", "u", 0, "List games for user id")
	listCmd.AddCommand(listGamesCmd)

}
