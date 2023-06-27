package main

import (
	"embed"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/sirgwain/craig-stars/cmd"
	"github.com/sirgwain/craig-stars/server"
)

// must use all: so we include _ files from sveltekit
//go:embed all:frontend/build
var assets embed.FS

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	server.SetAssets(assets)
	cmd.Execute()
}
