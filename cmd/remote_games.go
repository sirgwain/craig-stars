package cmd

import (
	"fmt"
	"io"

	"connectrpc.com/connect"
	"github.com/lensesio/tableprinter"
	craig_starsv1 "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func newRemoteGamesCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "games",
		Short: "List games from a craig-stars server",
		Long:  `List games for the authenticated user from a craig-stars server.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := resolveCLIClientConfig(true)
			if err != nil {
				return err
			}
			clients := newRemoteClients(cfg)
			resp, err := clients.Games.GetGames(cmd.Context(), connect.NewRequest(&craig_starsv1.GetGamesRequest{}))
			if err != nil {
				return err
			}
			if jsonOut {
				return printProtoJSON(cmd.OutOrStdout(), resp.Msg)
			}
			return printGamesTable(cmd.OutOrStdout(), resp.Msg.GetGames())
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "Print JSON output")
	return cmd
}

func printGamesTable(out io.Writer, games []*craig_starsv1.GameWithPlayers) error {
	type gameRow struct {
		ID      int64  `header:"ID"`
		Name    string `header:"Name"`
		State   string `header:"State"`
		Year    string `header:"Year"`
		Players int    `header:"Players"`
		Open    int32  `header:"Open"`
		Public  bool   `header:"Public"`
	}

	rows := make([]gameRow, 0, len(games))
	for _, item := range games {
		game := item.GetGame()
		rows = append(rows, gameRow{
			ID:      game.GetId(),
			Name:    game.GetName(),
			State:   game.GetState().String(),
			Year:    fmt.Sprintf("%d", game.GetYear()),
			Players: len(item.GetPlayers()),
			Open:    game.GetOpenPlayerSlots(),
			Public:  game.GetPublic(),
		})
	}

	tableprinter.New(out).Print(rows)
	return nil
}

func printProtoJSON(out interface{ Write([]byte) (int, error) }, msg proto.Message) error {
	marshaler := protojson.MarshalOptions{
		Multiline:       true,
		Indent:          "  ",
		EmitUnpopulated: true,
	}
	b, err := marshaler.Marshal(msg)
	if err != nil {
		return err
	}
	b = append(b, '\n')
	_, err = out.Write(b)
	return err
}

func init() {
	rootCmd.AddCommand(newRemoteGamesCmd())
}
