package cmd

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/db"
	"github.com/spf13/cobra"
)

var client db.Client
var dbConn db.DBConn

func dbPreRun(cmd *cobra.Command, args []string) error {
	debugPreRun(cmd, args)
	dbConn = db.NewConn()
	cfg := config.GetConfig()
	dbConn.Connect(cfg)

	var err error
	client, err = dbConn.BeginTransaction()
	if err != nil {
		return err
	}

	return nil
}

func dbPostRun(cmd *cobra.Command, args []string) error {
	if err := dbConn.Commit(client); err != nil {
		return err
	}
	return nil
}

func dbRollback() error {
	if client != nil && dbConn != nil {
		log.Info().Msgf("rolling back transaction")
		if err := dbConn.Rollback(client); err != nil {
			panic(fmt.Errorf("failed to rollback transaction %v", err))
		}
	}
	return nil
}
