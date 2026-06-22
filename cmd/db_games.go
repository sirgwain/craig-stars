package cmd

import (
	"github.com/spf13/cobra"
)

func newDBCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db",
		Short: "Interact with the local sqlite database",
		Long:  `Root command for interacting with the local database.`,
	}

	cmd.AddCommand(newDBGamesCmd())
	return cmd
}

func newDBGamesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "games",
		Short: "Interact with games",
		Long:  `Root command for interacting with games.`,
	}

	cmd.AddCommand(newDBGamesListCmd())
	cmd.AddCommand(newDBGamesUpdateCmd())
	return cmd
}

func init() {
	rootCmd.AddCommand(newDBCmd())
}
