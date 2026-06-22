package cmd

import (
	"fmt"
	"io"
	"strings"

	"connectrpc.com/connect"
	"github.com/lensesio/tableprinter"
	craig_starsv1 "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1"
	"github.com/spf13/cobra"
)

func newRemoteRacesCmd() *cobra.Command {
	var jsonOut bool
	var name string
	cmd := &cobra.Command{
		Use:   "races",
		Short: "List races from a craig-stars server",
		Long:  `List races for the authenticated user from a craig-stars server.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := resolveCLIClientConfig(true)
			if err != nil {
				return err
			}
			clients := newRemoteClients(cfg)
			resp, err := clients.Races.GetRaces(cmd.Context(), connect.NewRequest(&craig_starsv1.GetRacesRequest{}))
			if err != nil {
				return err
			}

			if name != "" {
				race := findRaceByName(resp.Msg.GetRaces(), name)
				if race == nil {
					return fmt.Errorf("race not found: %s", name)
				}
				return printProtoJSON(cmd.OutOrStdout(), race)
			}
			if jsonOut {
				return printProtoJSON(cmd.OutOrStdout(), resp.Msg)
			}
			return printRacesTable(cmd.OutOrStdout(), resp.Msg.GetRaces())
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "Print JSON output")
	cmd.Flags().StringVar(&name, "name", "", "Show details for a race by exact name")
	return cmd
}

func findRaceByName(races []*craig_starsv1.Race, name string) *craig_starsv1.Race {
	for _, race := range races {
		if strings.EqualFold(race.GetName(), name) {
			return race
		}
	}
	return nil
}

func printRacesTable(out io.Writer, races []*craig_starsv1.Race) error {
	type raceRow struct {
		ID        int64  `header:"ID"`
		Name      string `header:"Name"`
		Plural    string `header:"Plural"`
		PRT       string `header:"PRT"`
		Growth    int32  `header:"Growth"`
		Factories string `header:"Factories"`
		Mines     string `header:"Mines"`
	}

	rows := make([]raceRow, 0, len(races))
	for _, race := range races {
		rows = append(rows, raceRow{
			ID:        race.GetId(),
			Name:      race.GetName(),
			Plural:    race.GetPluralName(),
			PRT:       race.GetPrt().String(),
			Growth:    race.GetGrowthRate(),
			Factories: fmt.Sprintf("%d/%d/%d", race.GetFactoryOutput(), race.GetFactoryCost(), race.GetNumFactories()),
			Mines:     fmt.Sprintf("%d/%d/%d", race.GetMineOutput(), race.GetMineCost(), race.GetNumMines()),
		})
	}

	tableprinter.New(out).Print(rows)
	return nil
}

func init() {
	rootCmd.AddCommand(newRemoteRacesCmd())
}
