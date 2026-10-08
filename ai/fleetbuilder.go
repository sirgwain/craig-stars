package ai

import (
	"fmt"
	"math"
	"slices"

	"github.com/sirgwain/craig-stars/cs"
)

type fleet struct {
	purpose cs.FleetPurpose
	ships   []fleetShip
}

type fleetShip struct {
	purpose  cs.ShipDesignPurpose
	quantity int
}

func (f fleet) requiredShips() map[cs.ShipDesignPurpose]int {
	required := make(map[cs.ShipDesignPurpose]int, len(f.ships))
	for _, ship := range f.ships {
		if ship.quantity > 0 {
			required[ship.purpose] += ship.quantity
		}
	}
	return required
}

// countCompleteFleets counts whole fleets, combining designs with the same purpose.
func (f fleet) countCompleteFleets(ships map[cs.ShipDesignPurpose]int) int {
	required := f.requiredShips()
	if len(required) == 0 {
		return 0
	}
	count := math.MaxInt
	for purpose, quantity := range required {
		count = min(count, ships[purpose]/quantity)
	}
	return count
}

// merge idle fleets matching the purposes we require into a single fleet
func (f *fleet) mergeFromIdleFleets(ai *aiPlayer, fleets []*cs.Fleet) (fleet *cs.Fleet, remainingFleets []*cs.Fleet, err error) {
	required := make(map[cs.ShipDesignPurpose]int, len(f.ships))
	for _, ship := range f.ships {
		current := required[ship.purpose]
		required[ship.purpose] = current + ship.quantity
	}

	remainingFleets = []*cs.Fleet{}

	// ai.log.Debug().
	// 	Str("Purpose", string(f.purpose)).
	// 	Msgf("%d fleets at location", len(fleets))

	// check existing idle ships to see if we have enough already lying around
	fleetsToMerge := []*cs.Fleet{}
	for i, fleet := range fleets {
		if fleet.GetTag("purpose") != string(f.purpose) {
			// fleet's current purpose doesn't match our desired purpose;
			// assume it's doing something else
			continue
		}

		if len(fleet.Tokens) != 1 {
			// fleet has multiple token types, likely indicating
			// an already partially built fleet
			continue
		}

		foundShip := false
		design := ai.GetDesign(fleet.Tokens[0].DesignNum)
		if design == nil {
			// no design; skip
			continue
		}

		// if we need this design, add it to our fleets to merge slice
		// TODO: Consider imposing a distance/ETA requirement to this?
		if requiredQuantity, found := required[design.Purpose]; found {
			required[design.Purpose] = requiredQuantity - fleet.Tokens[0].Quantity
			fleetsToMerge = append(fleetsToMerge, fleet)
			foundShip = true

			// ai.log.Debug().
			// 	Str("Purpose", string(f.purpose)).
			// 	Msgf("tapping %s at planet %d for %s", fleet.Name, fleet.OrbitingPlanetNum, design.Purpose)

			// we've found all the ships we need for this requirement, remove it
			if required[design.Purpose] <= 0 {
				delete(required, design.Purpose)
			}
		}
		if !foundShip {
			remainingFleets = append(remainingFleets, fleet)
			// ai.log.Debug().
			// 	Str("Purpose", string(f.purpose)).
			// 	Msgf("skipping %s", fleet.Name)
		}
		// we're done, we have a fleet!
		if len(required) == 0 {
			// add any fleets we skipped to the remaining list and break out, we're done
			remainingFleets = append(remainingFleets, fleets[i+1:]...)
			// ai.log.Debug().
			// 	Msgf("found ships for %s, %d fleets remaining", string(f.purpose), len(remainingFleets))

			break
		}

	}

	if len(required) == 0 {
		// we only needed one fleet, return
		if len(fleetsToMerge) == 1 {
			return fleetsToMerge[0], remainingFleets, nil
		}
		fleet, err := ai.merge(fleetsToMerge)
		if err != nil {
			return nil, remainingFleets, err
		}

		// rename this fleet based on our purpose
		fleet.Rename(ai.fleetName(f.purpose))

		return fleet, remainingFleets, nil
	}

	return nil, remainingFleets, nil
}

// get any fleets matching this fleetmakeup
func (f *fleet) getFleetsMatchingMakeup(ai *aiPlayer, fleets []*cs.Fleet) []*cs.Fleet {
	matchingFleets := []*cs.Fleet{}

	for _, fleet := range fleets {
		if fleet.GetTag(cs.TagPurpose) != string(f.purpose) {
			continue
		}

		ships := make(map[cs.ShipDesignPurpose]int, len(fleet.Tokens))
		for _, token := range fleet.Tokens {
			design := ai.GetDesign(token.DesignNum)
			if design != nil {
				ships[design.Purpose] += token.Quantity
			}
		}

		if f.countCompleteFleets(ships) > 0 {
			matchingFleets = append(matchingFleets, fleet)
		}
	}

	return matchingFleets
}

// assemble fleets of a given purpose from all idle fleets over planets (i.e. fleets that were just built)
func (ai *aiPlayer) assembleFromIdleFleets(fleetMakeup fleet) ([]*cs.Fleet, error) {
	// find all idle fleets that are colonizers
	// and merge them into a single fleet with purpose
	fleets := []*cs.Fleet{}
	for planetNum, fleets := range ai.fleetsByPlanetNum {

		fleetsToCheck := make([]*cs.Fleet, len(fleets))
		copy(fleetsToCheck, fleets)
		if planetNum != cs.None {
			for {
				fleet, remainingFleets, err := fleetMakeup.mergeFromIdleFleets(ai, fleetsToCheck)
				fleetsToCheck = remainingFleets
				if err != nil {
					return nil, fmt.Errorf("failed to merge %s fleet %v", fleetMakeup.purpose, err)
				}
				if fleet == nil {
					break
				}
				// we made a fleet, hoorah
				fleet.Purpose = fleetMakeup.purpose
			}
		}
	}

	return fleets, nil
}

// merge a fleet and update the ai maps
func (ai *aiPlayer) merge(fleets []*cs.Fleet) (*cs.Fleet, error) {
	fleet, err := ai.client.Merge(&ai.game.Rules, ai.Player, fleets)
	if err != nil {
		return nil, err
	}

	// update the AI with new fleet info
	updatedFleets := make([]*cs.Fleet, 0, len(ai.Fleets)-len(fleets)+1)
	for _, existingFleet := range ai.Fleets {
		if slices.Index(fleets, existingFleet) == -1 {
			updatedFleets = append(updatedFleets, existingFleet)
		}
	}
	updatedFleets = append(updatedFleets, fleet)
	ai.Fleets = updatedFleets

	updatedFleetsAtPlanetNum := make([]*cs.Fleet, 0, len(ai.fleetsByPlanetNum[fleet.OrbitingPlanetNum])-len(fleets)+1)
	for _, existingFleet := range ai.fleetsByPlanetNum[fleet.OrbitingPlanetNum] {
		if slices.Index(fleets, existingFleet) == -1 {
			updatedFleetsAtPlanetNum = append(updatedFleetsAtPlanetNum, existingFleet)
		}
	}
	updatedFleetsAtPlanetNum = append(updatedFleetsAtPlanetNum, fleet)
	ai.fleetsByPlanetNum[fleet.OrbitingPlanetNum] = updatedFleetsAtPlanetNum

	// remove from fleetsByNum
	for i := 1; i < len(fleets); i++ {
		delete(ai.fleetsByNum, fleets[i].Num)
	}

	return fleet, nil
}
