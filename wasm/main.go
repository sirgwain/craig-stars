//go:build wasi || wasm

package main

import (
	"fmt"
	"os"
	"syscall/js"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/wasm/wasm"
)

// Each wasm instannce is unique to a browser session, so keep track of state so we don't have to
// send it and serialize it for each call
type state struct {
	rules   cs.Rules
	player  cs.Player
	designs []cs.ShipDesign
}

var ctx = state{
	rules: cs.NewRules(),
}
var debug = false

func enableDebug(this js.Value, args []js.Value) interface{} {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime, NoColor: true})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	debug = true
	log.Debug().Msg("enabled debug mode")
	return js.Undefined()
}

// set the rules used by this wasm instance
// rules default to a standard ruleset, but are overloaded during game load
func setRules(_ js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return wasm.NewError(fmt.Errorf("setRules: number of arguments doesn't match"))
	}

	ctx.rules = wasm.GetRules(args[0])

	return js.Undefined()
}

// setPlayer sets or updates the current player for this wasm instance
func setPlayer(_ js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return wasm.NewError(fmt.Errorf("setPlayer: number of arguments doesn't match"))
	}

	player := wasm.GetPlayer(args[0])
	player.Designs = ctx.player.Designs
	ctx.player = player

	log.Debug().Msgf("setting active player")
	return js.Undefined()
}

// setDesigns sets or updates the current player's designs for this wasm instance
func setDesigns(_ js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return wasm.NewError(fmt.Errorf("setDesigns: number of arguments doesn't match"))
	}

	designs := wasm.GetSlice(args[0], wasm.GetShipDesign)
	ctx.player.Designs = make([]*cs.ShipDesign, len(designs))
	for i := range designs {
		ctx.player.Designs[i] = &designs[i]
	}

	log.Debug().Msgf("setting player designs")
	return js.Undefined()
}

// wasm wrapper for calculating race points
// takes one argument, the race
func calculateRacePoints(_ js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return wasm.NewError(fmt.Errorf("number of arguments doesn't match"))
	}

	race := wasm.GetRace(args[0])
	points := race.ComputeRacePoints(ctx.rules.RaceStartingPoints)
	log.Debug().Msgf("calculated points for race %s: %d", race.PluralName, points)

	return js.ValueOf(points)
}

// wasm wrapper for estimating planet production
// takes 1 arguments: planet, player (with designs)
func estimateProduction(_ js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return wasm.NewError(fmt.Errorf("number of arguments doesn't match"))
	}

	planet := wasm.GetPlanet(args[0])

	// setup the starbase
	if planet.Spec.HasStarbase {
		planet.Starbase = &cs.Fleet{
			Tokens: []cs.ShipToken{
				{Quantity: 1, DesignNum: planet.Spec.StarbaseDesignNum},
			},
		}
	}

	// make sure if we have a starbase, it has a design so we can compute
	// upgrade costs
	if err := planet.PopulateStarbaseDesign(&ctx.player); err != nil {
		return wasm.NewError(fmt.Errorf("failed to populate starbase with player design. %v", err))
	}

	if err := planet.PopulateProductionQueueDesigns(&ctx.player); err != nil {
		return wasm.NewError(fmt.Errorf("failed to populate production queue designs. %v", err))
	}

	planet.PopulateProductionQueueEstimates(&ctx.rules, &ctx.player)

	log.Debug().Msgf("estimatied production of %s\n", planet.Name)
	o := js.ValueOf(map[string]any{})
	wasm.SetPlanet(o, &planet)
	return o
}

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	wasm.ExposeFunction("setRules", js.FuncOf(setRules))
	wasm.ExposeFunction("setPlayer", js.FuncOf(setPlayer))
	wasm.ExposeFunction("setDesigns", js.FuncOf(setDesigns))
	wasm.ExposeFunction("enableDebug", js.FuncOf(enableDebug))
	wasm.ExposeFunction("calculateRacePoints", js.FuncOf(calculateRacePoints))
	wasm.ExposeFunction("estimateProduction", js.FuncOf(estimateProduction))
	wasm.Ready()

	// fmt.Println("wasm initialized")
	<-make(chan bool) // To use anything from Go WASM, the program may not exit.
}
