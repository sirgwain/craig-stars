package ai

import (
	"fmt"

	"github.com/sirgwain/craig-stars/cs"
)

func (ai *aiPlayer) produce() error {

	// check for builds of whole fleets
	for fleetPurpose, quantity := range ai.requests.fleetBuilds {
		if quantity <= 0 {
			continue
		}

		for _, planet := range ai.Planets {
			fleetMakeup := ai.fleetsByPurpose[fleetPurpose].clone()

			// make sure this planet is ready to build this type of fleet
			if !ai.isPlanetReadyToBuildFleet(planet, fleetMakeup.purpose) {
				continue
			}

			for shipIndex, ship := range fleetMakeup.ships {
				idleShips := ai.getIdleShipCount(planet, fleetPurpose, ship.purpose)
				if idleShips > 0 {
					fleetMakeup.ships[shipIndex].quantity = max(0, ship.quantity-idleShips)
				}
			}

			for _, ship := range fleetMakeup.ships {
				if ship.quantity <= 0 {
					continue
				}

				// design and upgrade this ship
				design, err := ai.designShip(ai.config.namesByPurpose[ship.purpose], ship.purpose, fleetMakeup.purpose)
				if err != nil {
					return fmt.Errorf("unable to design ship %v: %w", ship.purpose, err)
				}
				if design == nil {
					ai.log.Debug().
						Int64("GameID", ai.GameID).
						Int("PlayerNum", ai.Num).
						Msgf("unable to design ship %v", ship.purpose)
					continue
				}

				if !planet.CanBuild(design.Spec.Mass) {
					continue
				}

				if !ai.isShipInQueue(planet, fleetMakeup.purpose, ship.purpose, ship.quantity) {
					ai.log.Debug().
						Int64("GameID", ai.GameID).
						Int("PlayerNum", ai.Num).
						Str("FleetPurpose", string(fleetMakeup.purpose)).
						Str("Purpose", string(ship.purpose)).
						Int("PlayerNum", ai.Num).
						Msgf("adding %d %s to %s queue", ship.quantity, design.Name, planet.Name)

					item := cs.NewProductionQueueItemShip(design, ship.quantity).WithTag(cs.TagPurpose, string(fleetMakeup.purpose))
					ai.prependToQueue(planet, *item)
					if err := ai.client.UpdatePlanetOrders(&ai.game.Rules, ai.Player, planet, planet.PlanetOrders, ai.Planets); err != nil {
						return err
					}
				}
			}
		}
	}

	// build scanners and starbases on each planet where applicable
	for _, planet := range ai.Planets {
		if !planet.Scanner && !ai.isItemInQueue(planet, cs.QueueItemTypePlanetaryScanner) {
			item := cs.ProductionQueueItem{Type: cs.QueueItemTypePlanetaryScanner, Quantity: 1}
			if err := ai.checkAndBuildItem(planet, item, ai.config.minYearsToBuildScanner); err != nil {
				return err
			}
		}

		// see if we should build a starbase or upgrade an existing one
		if err := ai.buildOrUpgradeStarbase(planet); err != nil {
			return err
		}

	}
	return nil
}

// check if a planet is in a state where it will build a fleet
func (ai *aiPlayer) isPlanetReadyToBuildFleet(planet *cs.Planet, purpose cs.FleetPurpose) bool {
	if !planet.Spec.HasStarbase || planet.Spec.DockCapacity == 0 {
		return false
	}

	// don't worry about production, let this planet build scouts
	if purpose == cs.FleetPurposeScout {
		return true
	}

	// don't build colonizers or freighters on this planet if it doesn't have the pop to move them
	if purpose == cs.FleetPurposeColonistFreighter &&
		planet.Spec.PopulationDensity >= ai.config.colonistTransportDensity {
		return true
	}
	if purpose == cs.FleetPurposeColonizer &&
		planet.Spec.PopulationDensity >= ai.config.colonizerPopulationDensity {
		return true
	}

	// don't build certain fleets or bombers unless we have enough installations
	// TODO: Change this for -f races and stuff
	planetaryStructuresBuilt := min(float64(planet.Mines)/float64(planet.Spec.MaxMines),
		float64(planet.Factories)/float64(planet.Spec.MaxFactories))

	cutoff := ai.config.fleetProductionCutoff
	if purpose == cs.FleetPurposeBomber {
		cutoff = ai.config.bomberProductionCutoff
	}

	if planetaryStructuresBuilt < cutoff {
		return false
	}

	return true
}

// add a new ship build request
func (ai *aiPlayer) addFleetBuildRequest(purpose cs.FleetPurpose, count int) {
	current := ai.requests.fleetBuilds[purpose]
	ai.requests.fleetBuilds[purpose] = current + count
}

// Build a new or upgrade an existing starbase on this planet.
func (ai *aiPlayer) buildOrUpgradeStarbase(planet *cs.Planet) error {
	// if we're already building a starbase, don't do anything
	if ai.isStarbaseInQueue(planet) {
		return nil
	}

	// check if we are targeted or being bombed by bombers
	enemyOrbitingFleets := ai.enemyShipsAbovePlanet(planet)
	attackShipsInOrbit := ai.hasAttackShips(enemyOrbitingFleets)
	_, targeted := ai.targetedPlanets[planet.Num]

	threatened := targeted || attackShipsInOrbit

	// don't build starbases if this planet has not moved forward enough economically
	// and we don't need them (ie not being attacked)
	planetaryStructuresBuilt := min(float64(planet.Mines)/float64(planet.Spec.MaxMines),
		float64(planet.Factories)/float64(planet.Spec.MaxFactories))
	if !threatened && planetaryStructuresBuilt < ai.config.fleetProductionCutoff {
		return nil
	}

	timeToWait := ai.config.minYearsToQueueStarbasePeaceTime
	if attackShipsInOrbit {
		timeToWait = ai.config.minYearsToQueueStarbaseWarTime
	}

	return ai.addStarbaseToQueue(planet, timeToWait, threatened)
}

// upgrade an existing starbase to a better or newer model, or build a new one if none are present
func (ai *aiPlayer) addStarbaseToQueue(planet *cs.Planet, timeToWait int, threatened bool) error {
	existingDesign := ai.GetDesign(planet.Spec.StarbaseDesignNum)
	if (existingDesign == nil) != planet.Spec.HasStarbase {
		// we either got a starbase design despite not expecting one or vice versa; funky stuff happened
		var err error
		var msg string
		if planet.Spec.HasStarbase {
			err = fmt.Errorf("no existing base design found despite planet.Spec.HasStarbase being true")
			msg = "design not found despite expecting one"
		} else {
			err = fmt.Errorf("existing base design found despite planet.Spec.HasStarbase being false")
			msg = "design found despite not expecting one"
		}

		ai.log.Err(err).
			Int64("GameID", ai.GameID).
			Int("PlayerNum", ai.Num).
			Int("PlanetNum", planet.Num).
			Str("PlanetName", planet.Name).
			Int("DesignNum", planet.Spec.StarbaseDesignNum).
			Str("DesignName", planet.Spec.StarbaseDesignName).
			Bool("HasStarbase", planet.Spec.HasStarbase).
			Msg(msg)

		return err
	}

	var purpose cs.ShipDesignPurpose
	if existingDesign == nil {
		if threatened {
			// we don't have a starbase, but we need one (invaded, bombed, etc)
			purpose = cs.ShipDesignPurposeFort
		} else {
			// we don't need an armed starbase yet, so build a fuel depot
			purpose = cs.ShipDesignPurposeFuelDepot
		}
	} else {
		// we have a starbase, so we need to upgrade it
		switch existingDesign.Purpose {
		case cs.ShipDesignPurposeFort, cs.ShipDesignPurposeFuelDepot, cs.ShipDesignPurposeStarbaseUnarmed:
			// forts & fuel depots --> 1/4 starbase
			purpose = cs.ShipDesignPurposeStarbaseQuarter
		case cs.ShipDesignPurposeStarbaseQuarter:
			// 1/4 starbase --> 1/2 starbase
			purpose = cs.ShipDesignPurposeStarbaseHalf
		case cs.ShipDesignPurposeStarbaseHalf, cs.ShipDesignPurposeStarbase:
			// 1/2 starbase --> full (or update existing ones)
			purpose = cs.ShipDesignPurposeStarbase
		}
	}

	design := ai.designsByPurpose[purpose]
	if design == nil || (existingDesign != nil && design.Num == existingDesign.Num) {
		// we either lack a new design to upgrade to or want to upgrade to the
		// exact same base as before; skip
		return nil
	}

	item := cs.ProductionQueueItem{
		Type:      cs.QueueItemTypeStarbase,
		DesignNum: design.Num,
		Quantity:  1,
	}

	return ai.checkAndBuildItem(planet, item, timeToWait)
}

// Compute years to build a given ProductionQueueItem, prepending it to the queue if time allows.
func (ai *aiPlayer) checkAndBuildItem(planet *cs.Planet, item cs.ProductionQueueItem, yearsCutoff int) error {
	yearsToBuild, err := ai.getYearsToBuild(planet, item)
	if err != nil {
		return err
	}

	if yearsToBuild != cs.Infinite && yearsToBuild <= yearsCutoff {
		ai.prependToQueue(planet, item)
	}
	return nil
}

func (ai *aiPlayer) getYearsToBuild(planet *cs.Planet, item cs.ProductionQueueItem) (yearsToBuild int, err error) {
	yearlyAvailableToSpend := cs.NewCostFromMineralAndResources(planet.Spec.MiningOutput, planet.Spec.ResourcesPerYearAvailable)
	costCalculator := cs.NewCostCalculator(&ai.game.Rules, ai.TechLevels, &ai.Race.Spec)
	completionEstimator := cs.NewCompletionEstimator()

	var oldStarbase *cs.ShipDesign
	if planet.Spec.HasStarbase {
		// leave design nil if planet lacks a starbase (since GetItemCost)
		oldStarbase = ai.GetDesign(planet.Starbase.Tokens[0].DesignNum)
	}
	cost, err := costCalculator.GetItemCost(item, oldStarbase)
	if err != nil {
		return 0, err
	}

	yearsToBuild = completionEstimator.GetYearsToBuild(item, cost, planet.Spec.MiningOutput, yearlyAvailableToSpend)
	return yearsToBuild, nil
}
