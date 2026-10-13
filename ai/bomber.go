package ai

import (
	"log/slog"
	"math"
	"slices"

	"github.com/sirgwain/craig-stars/cs"
)

func (ai *aiPlayer) bomb() error {
	// catalogue all explored planets owned by our enemies
	bombablePlanets := map[int]*cs.Planet{}
	for _, planet := range ai.Player.PlanetIntels {
		if ai.IsEnemy(planet.PlayerNum) {
			bombablePlanets[planet.Num] = planet
		}
	}

	fleetMakeup := ai.fleetsByPurpose[cs.FleetPurposeBomber]

	// find all idle ships that can be part of our colonizer fleets
	// and merge them into a single fleet with purpose
	bomberFleets, err := ai.assembleFromIdleFleets(fleetMakeup)
	if err != nil {
		return err
	}

	// go through fleets in space that may have been misassigned
	// The recipe grows over the years. A fleet assembled to an earlier recipe is
	// still a bomber fleet, so don't require it to match today's.
	for _, fleet := range ai.Fleets {
		if fleet.GetTag(cs.TagPurpose) == string(cs.FleetPurposeBomber) && fleet.Spec.Bomber && len(fleet.Tokens) > 1 {
			if len(fleet.Waypoints) > 1 {
				// this fleet is already targeting a planet, remove the target from the bombable planets list
				for _, wp := range fleet.Waypoints[1:] {
					target := ai.GetPlanetIntel(wp.TargetNum)
					if target == nil || !target.Owned() || !ai.IsEnemy(target.PlayerNum) {
						// this fleet is targeting a planet that is no longer bombable, clear its waypoints so it can be used
						// to bomb something else
						fleet.Waypoints = fleet.Waypoints[:1]
						bomberFleets = append(bomberFleets, fleet)
						ai.log.Debug("Fleet was going to bomb target, but it's no longer bombable", slog.String("fleet", fleet.Name))

					} else {
						delete(bombablePlanets, wp.TargetNum)
					}
				}
			} else {
				// if we are already orbiting a bombable planet, we're good
				if _, found := bombablePlanets[fleet.OrbitingPlanetNum]; found {
					delete(bombablePlanets, fleet.OrbitingPlanetNum)
				} else {
					// we aren't orbiting a bombable planet yet, go find one!
					bomberFleets = append(bomberFleets, fleet)
				}
			}
		}
	}

	// after colonizing, we may have idle fleets leftover
	idleFleets := len(bomberFleets)
	ai.log.Debug("Found bombers and bombable planets", slog.Int("bombers", idleFleets), slog.Int("bombablePlanets", len(bombablePlanets)))

	// fleets at one of our worlds with nothing they can take on
	waiting := map[int][]*cs.Fleet{}
	for _, fleet := range bomberFleets {
		if ai.enabled("war-retreat") && ai.lowOnFuel(fleet) {
			// Go home to refuel before the next attack, rather than crawl
			// onwards or sit stranded where the last fight ended.
			if starbase := ai.getClosestStarbasePlanet(fleet); starbase != nil && starbase.Num != fleet.OrbitingPlanetNum {
				if index := fleet.AddWaypoint(ai.Player, cs.WaypointDest{MO: starbase.MapObject}, len(fleet.Waypoints)-1, false); index > 0 {
					ai.client.UpdateFleetOrders(ai.Player, fleet, fleet.FleetOrders)
					idleFleets--
				}
				continue
			}
			if ai.lowOnFuel(fleet) {
				// at a starbase, the tanks are full again next year
				continue
			}
		}
		bestPlanet := ai.getBestPlanetToBomb(fleet, bombablePlanets)
		if bestPlanet == nil && ai.getPlanet(fleet.OrbitingPlanetNum) != nil {
			waiting[fleet.OrbitingPlanetNum] = append(waiting[fleet.OrbitingPlanetNum], fleet)
		}
		if bestPlanet != nil {
			newWpIndex := fleet.AddWaypoint(ai.Player, cs.WaypointDest{MO: bestPlanet.MapObject}, len(fleet.Waypoints)-1, true)
			if newWpIndex == 0 {
				ai.log.Warn("Fleet tried to target planet for bombing but did not add the waypoint", slog.String("fleet", fleet.Name), slog.String("planet", bestPlanet.Name))
				continue
			}

			ai.client.UpdateFleetOrders(ai.Player, fleet, fleet.FleetOrders)
			delete(bombablePlanets, bestPlanet.Num)
			idleFleets--

			ai.log.Debug("Fleet targeting planet for bombing",
				slog.String("fleet", fleet.Name),
				slog.String("planet", bestPlanet.Name),
				slog.Int("Fuel", fleet.Fuel),
				slog.Int("WarpSpeed", fleet.Waypoints[newWpIndex].WarpSpeed),
				slog.Int("EstFuelUsage", fleet.Waypoints[newWpIndex].EstFuelUsage),
				slog.Int("Population", bestPlanet.GetPopulation()),
				slog.Bool("HasStarbase", bestPlanet.Spec.HasStarbase))

		}
	}

	// Fleets too weak for anything on their own are stronger together.
	if ai.enabled("mass-fleets") {
		planetNums := make([]int, 0, len(waiting))
		for num := range waiting {
			planetNums = append(planetNums, num)
		}
		slices.Sort(planetNums)
		for _, num := range planetNums {
			if len(waiting[num]) < 2 {
				continue
			}
			fleet, err := ai.merge(waiting[num])
			if err != nil {
				return err
			}
			fleet.Rename(ai.fleetName(cs.FleetPurposeBomber))
			fleet.SetTag(cs.TagPurpose, string(cs.FleetPurposeBomber))
			fleet.Purpose = cs.FleetPurposeBomber
		}
	}

	// build bomber fleets where necessary
	if len(bombablePlanets)-idleFleets > 0 {
		ai.addFleetBuildRequest(cs.FleetPurposeBomber, min(ai.warLimit(), len(bombablePlanets)-idleFleets))
	}
	return nil
}

// a war fleet with less than this share of its fuel goes home before its next attack
const warFleetLowFuel = .25

func (ai *aiPlayer) lowOnFuel(fleet *cs.Fleet) bool {
	return float64(fleet.Fuel) < warFleetLowFuel*float64(fleet.Spec.FuelCapacity)
}

// a target further than this many years away isn't worth setting out for
const maxAttackYears = 10

// get the planet with the best distance to hab ratio
func (ai *aiPlayer) getBestPlanetToBomb(fleet *cs.Fleet, planets map[int]*cs.Planet) *cs.Planet {
	var best *cs.Planet = nil

	// lowest weight wins
	bestWeight := 0.0
	yearlyTravelDistance := float64(fleet.Spec.Engine.IdealSpeed * fleet.Spec.Engine.IdealSpeed)

	for num := range planets {
		intel := planets[num]
		if !ai.canAttackPlanet(fleet, intel) {
			continue
		}
		// go for low pop planets to clean them up quickly
		// like the snowball method of paying off debt
		pop := intel.GetPopulation()

		// avoid starbases if there are better targets
		starbaseFactor := 1.0
		if intel.Spec.HasStarbase {
			starbaseFactor = 2.0
		}
		dist := intel.Position.DistanceTo(fleet.Position)
		yearsToTravel := max(1, math.Ceil(dist/yearlyTravelDistance))
		if ai.enabled("war-range") {
			// what matters is how fast the fleet can really get there on its fuel
			warp := ai.getWarpSpeed(fleet, intel.ToTarget(), true)
			yearsToTravel = max(1, math.Ceil(dist/float64(max(1, warp*warp))))
			if yearsToTravel > maxAttackYears {
				continue
			}
		}

		// weight is based on distance (in years) and pop
		// closer is better, lower pop is better, starbases are discouraged
		weight := (1 / float64(max(1, pop))) / (2 * yearsToTravel) / starbaseFactor

		// ai.log.Debug().
		// 	Float64("dist", dist).
		// 	Int("pop", pop).
		// 	Float64("yearsToTravel", yearsToTravel).
		// 	Bool("hasStarbase", intel.Spec.HasStarbase).
		// 	Float64("weight", weight*1000).
		// 	Msgf("getBestPlanetToBomb %s", intel.Name)

		if weight > bestWeight || (weight == bestWeight && best != nil && intel.Num < best.Num) {
			bestWeight = weight
			best = intel
		}
	}

	// best = ai.getClosestPlanetIntel(fleet.Position, colonizablePlanets)

	return best
}

// canAttackPlanet compares a fleet with the defenses we have actually
// discovered, without looking into the full universe. Power rating is a coarse
// estimate, not a battle prediction.
func (ai *aiPlayer) canAttackPlanet(fleet *cs.Fleet, planet *cs.Planet) bool {
	// how much stronger than the known defenders the fleet must be: Easy
	// doesn't check, Normal wants an even fight, Hard wants an edge
	margin := 1.25
	switch ai.AIDifficulty {
	case cs.AIDifficultyEasy:
		return true
	case cs.AIDifficultyNormal:
		if !ai.enabled("attack-judgment") {
			return true
		}
		margin = 1
	}
	defending := 0
	if planet.Spec.HasStarbase {
		if d := ai.GetForeignDesign(planet.PlayerNum, planet.Spec.StarbaseDesignNum); d != nil {
			defending += d.Spec.PowerRating
		}
	}
	for _, enemy := range ai.enemyShipsAbovePlanet(planet) {
		for _, token := range enemy.Tokens {
			if d := ai.GetForeignDesign(enemy.PlayerNum, token.DesignNum); d != nil {
				defending += token.Quantity * d.Spec.PowerRating
			}
		}
	}
	return float64(ai.fleetPower(fleet)) >= float64(defending)*margin
}
