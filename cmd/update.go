//go:build !wasi && !wasm

package cmd

import (
	"github.com/spf13/cobra"
)

func newUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a record",
		Long:  `Update a record in the database.`,
	}

	cmd.AddCommand(newUpdateGameCmd())
	return cmd
}

func init() {
	rootCmd.AddCommand(newUpdateCmd())
}
