//go:build !wasi && !wasm

package cmd

import (
	"github.com/sirgwain/craig-stars/update"
	"github.com/spf13/cobra"
)

var gameID int64

func newUpdateGameCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "game",
		Short: "Update a game",
		Long:  `Update a game.`,
	}
	cmd.PersistentFlags().Int64Var(&gameID, "game-id", 0, "The game to update")
	cmd.MarkFlagRequired("game-id")

	cmd.AddCommand(newUpdateGameHostCmd())
	cmd.AddCommand(newUpdateGamePlayerCmd())
	return cmd
}

func newUpdateGameHostCmd() *cobra.Command {
	var userID int64
	cmd := &cobra.Command{
		Use:   "host",
		Short: "Update a game's host",
		Long:  `Update a game's host to a new player.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			return update.UpdateHost(gameID, userID)
		},
	}
	cmd.Flags().Int64Var(&userID, "user-id", 0, "The id of the user to make the new host")
	cmd.MarkFlagRequired("user-id")
	return cmd
}

func newUpdateGamePlayerCmd() *cobra.Command {
	var userID int64
	var playerNum int
	cmd := &cobra.Command{
		Use:   "player",
		Short: "Update a game's player",
		Long:  `Update a game's player to a new user.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			return update.UpdatePlayer(gameID, playerNum, userID)
		},
	}
	cmd.Flags().Int64Var(&userID, "user-id", 0, "The id of the user to take over the player")
	cmd.Flags().IntVar(&playerNum, "player-num", 0, "The number of the player to update")
	cmd.MarkFlagRequired("user-id")
	cmd.MarkFlagRequired("player-num")
	return cmd
}
