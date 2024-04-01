package ai

import (
	"fmt"
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

// make a clone of this fleet makeup so we modify it
func (f fleet) clone() fleet {
	c := f
	c.ships = make([]fleetShip, len(f.ships))
	copy(c.ships, f.ships)
	return c
}

// merge idle fleets matching the purposes we require into a single fleet
func (f *fleet) mergeFromIdleFleets(ai *aiPlayer, fleets []*cs.Fleet) (fleet *cs.Fleet, remainingFleets []*cs.Fleet, err error) {
	required := make(map[cs.ShipDesignPurpose]int, len(f.ships))
	for _, ship := range f.ships {
		current := required[ship.purpose]
		required[ship.purpose] = current + ship.quantity
	}

	remainingFleets = []*cs.Fleet{}

	// log.Debug().
	// 	Int64("GameID", ai.GameID).
	// 	Int("PlayerNum", ai.Num).
	// 	Str("Purpose", string(f.purpose)).
	// 	Msgf("%d fleets at location", len(fleets))

	// see if we have enough of what we need
	fleetsToMerge := []*cs.Fleet{}
	for i, fleet := range fleets {
		if fleet.GetTag("purpose") != string(f.purpose) {
			continue
		}
		foundShip := false
		if len(fleet.Tokens) == 1 {
			design := ai.GetDesign(fleet.Tokens[0].DesignNum)
			if design != nil {
				// if we need dthis design, add it to our fleets to merge
				if requiredQuantity, found := required[design.Purpose]; found {
					required[design.Purpose] = requiredQuantity - fleet.Tokens[0].Quantity
					fleetsToMerge = append(fleetsToMerge, fleet)
					foundShip = true
					// log.Debug().
					// 	Int64("GameID", ai.GameID).
					// 	Int("PlayerNum", ai.Num).
					// 	Str("Purpose", string(f.purpose)).
					// 	Msgf("tapping %s at planet %d for %s", fleet.Name, fleet.OrbitingPlanetNum, design.Purpose)

					// we've found all the ships we need for this requirement, remove it
					if required[design.Purpose] <= 0 {
						delete(required, design.Purpose)
					}
				}
			}
		}
		if !foundShip {
			remainingFleets = append(remainingFleets, fleet)
			// log.Debug().
			// 	Int64("GameID", ai.GameID).
			// 	Int("PlayerNum", ai.Num).
			// 	Str("Purpose", string(f.purpose)).
			// 	Msgf("skipping %s", fleet.Name)
		}
		// we're done, we have a fleet!
		if len(required) == 0 {
			// add any fleets we skipped to the remaining list and break out, we're done
			remainingFleets = append(remainingFleets, fleets[i+1:]...)
			// log.Debug().
			// 	Int64("GameID", ai.GameID).
			// 	Int("PlayerNum", ai.Num).
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
		fleet.Rename(ai.fleetName(fleet, f.purpose))

		return fleet, remainingFleets, nil
	}

	return nil, remainingFleets, nil
}

// get any fleets matching this fleetmakeup
func (f *fleet) getFleetsMatchingMakeup(ai *aiPlayer, fleets []*cs.Fleet) []*cs.Fleet {
	matchingFleets := []*cs.Fleet{}
	required := make(map[cs.ShipDesignPurpose]int, len(f.ships))
	for _, ship := range f.ships {
		current := required[ship.purpose]
		required[ship.purpose] = current + ship.quantity
	}

	for _, fleet := range fleets {
		if fleet.GetTag(cs.TagPurpose) != string(f.purpose) {
			continue
		}

		matches := 0
		for _, token := range fleet.Tokens {
			design := ai.GetDesign(token.DesignNum)
			if design != nil {
				// if we need dthis design, add it to our fleets to merge
				if requiredQuantity, found := required[design.Purpose]; found {
					required[design.Purpose] = requiredQuantity - token.Quantity
					if required[design.Purpose] <= 0 {
						matches++
					}
				}
			}
		}

		// if we match all ships in the fleet
		if matches == len(f.ships) {
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
