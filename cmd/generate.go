//go:build !wasi && !wasm

package cmd

import (
	"fmt"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/proto/converter"
	v1 "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
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
			// Load current in-memory tech store
			techs := cs.StaticTechStore

			var resp = &v1.GetTechsResponse{
				PlanetaryScanners: converter.C.ConvertCSTechPlanetaryScanners(techs.PlanetaryScanners),
				Terraforms:        converter.C.ConvertCSTechTerraforms(techs.Terraforms),
				Defenses:          converter.C.ConvertCSTechDefenses(techs.Defenses),
				Planetaries:       converter.C.ConvertCSTechPlanetaries(techs.Planetaries),
				HullComponents:    converter.C.ConvertCSTechHullComponents(techs.HullComponents),
				Hulls:             converter.C.ConvertCSTechHulls(techs.Hulls),
			}

			// Marshal with protobuf JSON, using proto field names and indentation
			marshaler := protojson.MarshalOptions{
				Multiline: true,
			}
			techsJSON, err := marshaler.Marshal(resp)
			if err != nil {
				return fmt.Errorf("failed to marshal techs to protobuf json: %w", err)
			}

			fmt.Println(string(techsJSON))
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

			// Marshal with protobuf JSON, using proto field names and indentation
			marshaler := protojson.MarshalOptions{
				Multiline: true,
			}
			rulesJson, err := marshaler.Marshal(converter.C.ConvertCSRules(&rules))
			if err != nil {
				return fmt.Errorf("failed to marshal techs to protobuf json: %w", err)
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
