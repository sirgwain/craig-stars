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
		Short: "Generate",
		Long:  `Update a record in the database.`,
	}

	cmd.AddCommand(newGenerateTechsJson())
	cmd.AddCommand(newGenerateRulesJson())
	return cmd
}

func newGenerateTechsJson() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "techsjson",
		Short: "Generate the techs.json content",
		Long:  `Generate the techs.json content. During build time we need to update this to the latest techs.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			techs := cs.StaticTechStore

			techsJson, err := json.MarshalIndent(techs, "", "\t")
			if err != nil {
				return fmt.Errorf("failed to marshal StaticTechStore to json: \n%w", err)
			}

			fmt.Println(string(techsJson))
			return nil
		},
	}

	return cmd
}

func newGenerateRulesJson() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rulesjson",
		Short: "Generate the rules.json content",
		Long:  `Generate the rules.json content. During build time we need to update this to the latest default rules.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			rules := cs.NewRules()

			rulesJson, err := json.MarshalIndent(rules, "", "\t")
			if err != nil {
				return fmt.Errorf("failed to marshal rules to json: \n%w", err)
			}

			fmt.Println(string(rulesJson))
			return nil
		},
	}

	return cmd
}

func init() {
	rootCmd.AddCommand(newGenerateCmd())
}
