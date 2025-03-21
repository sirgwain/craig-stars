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

func enableDebug(args []js.Value) any {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime, NoColor: true})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	debug = true
	log.Debug().Msg("enabled debug mode")
	return js.Undefined()
}

// set the rules used by this wasm instance
// rules default to a standard ruleset, but are overloaded during game load
func setRules(args []js.Value) any {
	if len(args) != 1 {
		return wasm.NewError(fmt.Errorf("setRules: number of arguments doesn't match"))
	}

	ctx.rules = wasm.GetRules(args[0])
	// TODO: support loaded tech stores eventually
	ctx.rules.SetTechStore(&cs.StaticTechStore)

	return js.Undefined()
}

// setPlayer sets or updates the current player for this wasm instance
func setPlayer(args []js.Value) any {
	if len(args) != 1 {
		return wasm.NewError(fmt.Errorf("setPlayer: number of arguments doesn't match"))
	}

	player := wasm.GetPlayer(args[0])
	player.Designs = ctx.player.Designs
	ctx.player = player

	log.Debug().Msgf("setting active player with %d designs", len(player.Designs))
	return js.Undefined()
}

// setDesigns sets or updates the current player's designs for this wasm instance
func setDesigns(args []js.Value) any {
	if len(args) != 1 {
		return wasm.NewError(fmt.Errorf("setDesigns: number of arguments doesn't match"))
	}

	designs := wasm.GetSlice(args[0], wasm.GetShipDesign)
	ctx.player.Designs = make([]*cs.ShipDesign, len(designs))
	for i := range designs {
		ctx.player.Designs[i] = &designs[i]
	}

	log.Debug().Msgf("setting %d player designs", len(designs))
	return js.Undefined()
}

// wasm wrapper for calculating race points
// takes one argument, the race
func calculateRacePoints(args []js.Value) any {
	if len(args) != 1 {
		return wasm.NewError(fmt.Errorf("number of arguments doesn't match"))
	}

	race := wasm.GetRace(args[0])
	points := race.ComputeRacePoints(ctx.rules.RaceStartingPoints)
	log.Debug().Msgf("calculated points for race %s: %d", race.PluralName, points)

	return js.ValueOf(points)
}

// wasm wrapper for calculating resources required to reach a given tech level
// takes one argument, the tech level
func getResearchCost(args []js.Value) any {
	if len(args) != 1 {
		return wasm.NewError(fmt.Errorf("number of arguments doesn't match"))
	}

	techLevel := wasm.GetTechLevel(args[0])

	resources := ctx.player.GetResearchCost(&ctx.rules, techLevel)
	log.Debug().Msgf("calculated research cost for %v: %d", techLevel, resources)

	return js.ValueOf(resources)
}

// wasm wrapper for computing a ship's ShipDesignSpec
// takes one argument, the design
func computeShipDesignSpec(args []js.Value) any {
	if len(args) != 1 {
		return wasm.NewError(fmt.Errorf("number of arguments doesn't match"))
	}

	design := wasm.GetShipDesign(args[0])
	spec, err := cs.ComputeShipDesignSpec(&ctx.rules, ctx.player.TechLevels, ctx.player.Race.Spec, &design)
	if err != nil {
		return wasm.NewError(fmt.Errorf("failed to compute spec for design %s: %v", design.Name, err))
	}
	log.Debug().Msgf("computed spec for design %s", design.Name)

	o := js.ValueOf(map[string]any{})
	wasm.SetShipDesignSpec(o, &spec)

	return o
}

// wasm wrapper for calculating race points
// takes one argument, the race
func starbaseUpgradeCost(args []js.Value) any {
	if len(args) != 2 {
		return wasm.NewError(fmt.Errorf("number of arguments doesn't match"))
	}

	design := wasm.GetShipDesign(args[0])
	newDesign := wasm.GetShipDesign(args[1])

	costCalculator := cs.NewCostCalculator()
	cost, err := costCalculator.StarbaseUpgradeCost(&ctx.rules, ctx.player.TechLevels, ctx.player.Race.Spec, &design, &newDesign)
	if err != nil {
		return wasm.NewError(fmt.Errorf("failed to calculate starbase upgrade cost: %v", err))
	}

	log.Debug().Msgf("computed starbase upgrade cost for design %q -> %q: %v", design.Name, newDesign.Name, cost)

	o := js.ValueOf(map[string]any{})
	wasm.SetCost(o, &cost)

	return o
}

// wasm wrapper for calculating race points
// takes one argument, the race
func techCost(args []js.Value) any {
	if len(args) != 1 {
		return wasm.NewError(fmt.Errorf("number of arguments doesn't match"))
	}

	tech := wasm.GetTech(args[0])
	costCalculatoor := cs.NewCostCalculator()
	cost := costCalculatoor.GetTechCost(&ctx.rules, ctx.player.TechLevels, ctx.player.Race.Spec, tech)

	log.Debug().Msgf("computed tech cost %s: %v", tech.Name, cost)

	o := js.ValueOf(map[string]any{})
	wasm.SetCost(o, &cost)

	return o
}

// wasm wrapper for estimating planet production
// takes 1 argument: the planet
func estimateProduction(args []js.Value) any {
	if len(args) != 1 {
		return wasm.NewError(fmt.Errorf("number of arguments doesn't match"))
	}

	planet := wasm.GetPlanet(args[0])

	if len(planet.ProductionQueue) == 0 {
		log.Debug().Msgf("Empty production queue; no estimates made")
		return args[0]
	}

	// setup the starbase
	if planet.Spec.HasStarbase {
		planet.Starbase = &cs.Fleet{
			Tokens: []cs.ShipToken{
				{Quantity: 1, DesignNum: planet.Spec.StarbaseDesignNum},
			},
		}
	}

	// populate starbase & production queue designs for starbase upgrades and packet cancels
	if err := planet.PopulateStarbaseDesign(&ctx.player); err != nil {
		return wasm.NewError(fmt.Errorf("failed to populate starbase with player design: %v", err))
	}

	planet.Starbase.Spec = cs.ComputeFleetSpec(&ctx.rules, &ctx.player, planet.Starbase)
	if err := planet.PopulateProductionQueueDesigns(&ctx.player); err != nil {
		return wasm.NewError(fmt.Errorf("failed to populate production queue designs: %v", err))
	}

	if err := planet.PopulateProductionQueueEstimates(&ctx.rules, &ctx.player); err != nil {
		return wasm.NewError(fmt.Errorf("failed to compute production queue estimates: %v", err))
	}

	log.Debug().Msgf("estimated production of planet %s\n", planet.Name)
	o := js.ValueOf(map[string]any{})
	wasm.SetPlanet(o, &planet)
	return o
}

// wasm wrapper for calculating planet maximium buildable installations
// takes 2 arguments: the planet and item type
func maxBuildable(args []js.Value) any {
	if len(args) != 2 {
		return wasm.NewError(fmt.Errorf("number of arguments doesn't match"))
	}

	planet := wasm.GetPlanet(args[0])
	itemType := wasm.GetQueueItemType(args[1])

	// auto items and ships have infinite cap
	maxBuild := 5000
	if !itemType.IsAuto() /* || itemType == cs.QueueItemTypeAutoMineralAlchemy */ {
		maxBuild = planet.MaxBuildable(&ctx.player, itemType)
		// Infinite is the constant integer of -1, but we want very big number
		if maxBuild == cs.Infinite {
			maxBuild = 5000
		}
	}

	log.Debug().Msgf("calculated planet max buildable for itemType %s: %d\n", itemType, maxBuild)
	return js.ValueOf(maxBuild)
}

// wasm wrapper for updating planet yearly resource production
// takes 1 argument: the planet
func updateResourcesAvailable(args []js.Value) any {
	if len(args) != 1 {
		return wasm.NewError(fmt.Errorf("number of arguments doesn't match"))
	}

	planet := wasm.GetPlanet(args[0])

	planet.Spec.ComputeResourcesPerYear(&ctx.player, planet.Factories,
		cs.ProductivePopulation(planet.GetPopulation(), planet.Spec.MaxPopulation,
			ctx.rules.PopulationOvercrowdResourcePenalty, ctx.rules.PopulationOvercrowdResourceMax),
		min(planet.GetPopulation(), planet.Spec.MaxPopulation))

	log.Debug().Msgf("calculated resource stats for planet %s\n", planet.Name)
	o := js.ValueOf(map[string]any{})
	wasm.SetPlanet(o, &planet)
	return o
}

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	wasm.ExposeFunction("setRules", setRules)
	wasm.ExposeFunction("setPlayer", setPlayer)
	wasm.ExposeFunction("setDesigns", setDesigns)
	wasm.ExposeFunction("enableDebug", enableDebug)
	wasm.ExposeFunction("calculateRacePoints", calculateRacePoints)
	wasm.ExposeFunction("getResearchCost", getResearchCost)
	wasm.ExposeFunction("computeShipDesignSpec", computeShipDesignSpec)
	wasm.ExposeFunction("starbaseUpgradeCost", starbaseUpgradeCost)
	wasm.ExposeFunction("techCost", techCost)
	wasm.ExposeFunction("estimateProduction", estimateProduction)
	wasm.ExposeFunction("maxBuildable", maxBuildable)
	wasm.ExposeFunction("resourcesAvailable", updateResourcesAvailable)
	wasm.Ready()

	// fmt.Println("wasm initialized")
	<-make(chan bool) // To use anything from Go WASM, the program may not exit.
}
