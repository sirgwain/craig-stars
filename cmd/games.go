package cmd

import (
	"github.com/spf13/cobra"
)

func newGamesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "games",
		Short: "Interact with games",
		Long:  `Root command for interacting with games.`,
	}

	cmd.AddCommand(newGamesListCmd())
	cmd.AddCommand(newGamesUpdateCmd())
	return cmd
}

func init() {
	rootCmd.AddCommand(newGamesCmd())
}
