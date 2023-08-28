package cmd

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/db"
	"github.com/spf13/cobra"
)

var client db.Client
var dbClient db.DBClient

func dbPreRun(cmd *cobra.Command, args []string) error {
	debugPreRun(cmd, args)
	dbClient = db.NewClient()
	cfg := config.GetConfig()
	dbClient.Connect(cfg)

	var err error
	client, err = dbClient.BeginTransaction()
	if err != nil {
		return err
	}

	return nil
}

func dbPostRun(cmd *cobra.Command, args []string) error {
	if err := dbClient.Commit(client); err != nil {
		return err
	}
	return nil
}

func dbRollback() error {
	if client != nil && dbClient != nil {
		log.Info().Msgf("rolling back transaction")
		if err := dbClient.Rollback(client); err != nil {
			panic(fmt.Errorf("failed to rollback transaction %v", err))
		}
	}
	return nil
}
