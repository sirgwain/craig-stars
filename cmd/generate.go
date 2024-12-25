//go:build !wasi && !wasm

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/spf13/cobra"
)

func newGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate ",
		Long:  `Update a record in the database.`,
	}

	cmd.AddCommand(newGenerateTechsJson())
	return cmd
}

func newGenerateTechsJson() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "techsjson",
		Short: "Generate the techs.json content",
		Long:  `Generate the techs.json content. During build time we need to update this to the latest techs.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			techs := cs.StaticTechStore

			techsJson, err := json.Marshal(techs)
			if err != nil {
				return fmt.Errorf("failed to marshal StaticTechStore to json")
			}

			fmt.Println(string(techsJson))
			return nil
		},
	}

	return cmd
}

func init() {
	rootCmd.AddCommand(newGenerateCmd())
}
