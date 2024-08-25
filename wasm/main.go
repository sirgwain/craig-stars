//go:build wasi || wasm

package main

import (
	"encoding/json"
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

var rules = cs.NewRules()
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
func setRules(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return wasm.NewError(fmt.Errorf("setRules: number of arguments doesn't match"))
	}

	log.Debug().Msgf("setting rules override")

	rulesJson := args[0].String()
	instanceRules := &cs.Rules{}
	json.Unmarshal([]byte(rulesJson), &instanceRules)

	rules = *instanceRules

	return js.Undefined()
}

// wasm wrapper for calculating race points
// takes one argument, the race
func calculateRacePoints(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return wasm.NewError(fmt.Errorf("number of arguments doesn't match"))
	}

	race := wasm.GetRace(args[0])
	points := race.ComputeRacePoints(rules.RaceStartingPoints)
	log.Debug().Msgf("calculated points for race %s: %d", race.PluralName, points)

	return js.ValueOf(points)
}

// wasm wrapper for estimating planet production
// takes 1 arguments: planet, player (with designs)
func estimateProduction(this js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return wasm.NewError(fmt.Errorf("number of arguments doesn't match"))
	}

	planetJson := args[0].String()
	planet := &cs.Planet{}
	json.Unmarshal([]byte(planetJson), &planet)

	playerJson := args[1].String()
	player := &cs.Player{}
	json.Unmarshal([]byte(playerJson), &player)

	// make sure if we have a starbase, it has a design so we can compute
	// upgrade costs
	if err := planet.PopulateStarbaseDesign(player); err != nil {
		return wasm.NewError(fmt.Errorf("failed to populate starbase with player design. %v", err))
	}

	if err := planet.PopulateProductionQueueDesigns(player); err != nil {
		return wasm.NewError(fmt.Errorf("failed to populate production queue designs. %v", err))
	}

	planet.PopulateProductionQueueEstimates(&rules, player)

	result, err := json.MarshalIndent(planet, "", "  ")
	if err != nil {
		return wasm.NewError(fmt.Errorf("failed to serialize result to json %v\n", err))
	}
	log.Debug().Msgf("estimatied production of %s for player %s\n%s", planet.Name, player.Name, result)

	return js.ValueOf(string(result))
}

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	wasm.ExposeFunction("setRules", js.FuncOf(setRules))
	wasm.ExposeFunction("enableDebug", js.FuncOf(enableDebug))
	wasm.ExposeFunction("calculateRacePoints", js.FuncOf(calculateRacePoints))
	wasm.ExposeFunction("estimateProduction", js.FuncOf(estimateProduction))
	wasm.Ready()

	// fmt.Println("wasm initialized")
	<-make(chan bool) // To use anything from Go WASM, the program may not exit.
}
