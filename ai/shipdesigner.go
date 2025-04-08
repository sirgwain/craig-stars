package ai

import (
	"fmt"

	"github.com/sirgwain/craig-stars/cs"
)

// get a ship design for a purpose
func (ai *aiPlayer) designShip(name string, purpose cs.ShipDesignPurpose, fleetPurpose cs.FleetPurpose) (updated *cs.ShipDesign, err error) {

	existing, found := ai.designsByPurpose[purpose]
	if found && existing.ID == 0 {
		// we already created a new design for this purpose, just use it
		return existing, nil
	}

	var hull *cs.TechHull
	switch purpose {
	case cs.ShipDesignPurposeScout:
		hull = ai.getBestHull(ai.techStore.GetHullsByType(cs.TechHullTypeScout))
	case cs.ShipDesignPurposeColonizer:
		hull = ai.getBestHull(ai.techStore.GetHullsByType(cs.TechHullTypeColonizer))
	case cs.ShipDesignPurposeFuelFreighter:
		hull = ai.getBestHull(ai.techStore.GetHullsByType(cs.TechHullTypeFuelTransport))
		if hull != nil {
			break
		}
		fallthrough
	case cs.ShipDesignPurposeColonistFreighter, cs.ShipDesignPurposeFreighter:
		hull = ai.getBestHull(ai.techStore.GetHullsByType(cs.TechHullTypeFreighter))
	case cs.ShipDesignPurposeArmedFreighter:
		hull = ai.getBestHull(ai.techStore.GetHullsByType(cs.TechHullTypeMultiPurposeFreighter))
	case cs.ShipDesignPurposeBeamFighter, cs.ShipDesignPurposeTorpedoFighter:
		hull = ai.getBestHull(ai.techStore.GetHullsByType(cs.TechHullTypeCapitalShip))
		if hull != nil {
			break
		}
		fallthrough
	case cs.ShipDesignPurposeFighterScout, cs.ShipDesignPurposeStartingFighter:
		hull = ai.getBestHull(ai.techStore.GetHullsByType(cs.TechHullTypeFighter))
	case cs.ShipDesignPurposeBomber:
		hull = ai.getBestHull(ai.techStore.GetHullsByType(cs.TechHullTypeBomber))
	case cs.ShipDesignPurposeFuelDepot:
		// TODO: Add support for searching for cheap fuel docks
		fallthrough
	case cs.ShipDesignPurposeStarterColony, cs.ShipDesignPurposeStarbase, cs.ShipDesignPurposeStarbaseQuarter,
		cs.ShipDesignPurposeStarbaseHalf, cs.ShipDesignPurposeStarbaseUnarmed:
		hull = ai.getBestHull(ai.techStore.GetHullsByType(cs.TechHullTypeStarbase))
		if hull != nil {
			break
		}
		fallthrough
	case cs.ShipDesignPurposePacketThrower, cs.ShipDesignPurposeStargater, cs.ShipDesignPurposeFort:
		hull = ai.getBestHull(ai.techStore.GetHullsByType(cs.TechHullTypeOrbitalFort))
	}

	if hull == nil {
		return existing, nil
	}

	updated, err = cs.DesignShip(&ai.game.Rules, ai.Player, hull, name, ai.GetNextDesignNum(ai.Designs), ai.DefaultHullSet, purpose, fleetPurpose)
	if err != nil {
		return existing, fmt.Errorf("cs.DesignShip returned error: %w", err)
	}
	// no need to compute ship design specs as functions already compute it before returning

	// if we tried to build a bomber with no bombs, ignore it
	if purpose == cs.ShipDesignPurposeBomber && !updated.Spec.Bomber {
		return nil, nil
	}

	if existing != nil {
		updated.Version = existing.Version + 1
	} else {
		updated.Version = 1
	}
	updated.Name = fmt.Sprintf("%s v%d", name, updated.Version+1)

	// if our existing design is equivalent or higher rated, return it
	if found && existing.SlotsEqual(updated.Slots) || (existing != nil && existing.Spec.PowerRating != 0 && existing.Spec.PowerRating >= updated.Spec.PowerRating) {
		return existing, nil
	}

	if err := updated.Validate(&ai.game.Rules, ai.Player); err != nil {
		return nil, fmt.Errorf("invalid updated design: %w", err)
	}

	ai.Designs = append(ai.Designs, updated)
	ai.designsByPurpose[updated.Purpose] = updated

	// otherwise return the new design
	return updated, nil
}

// design various types of starbases we may build on planets
func (ai *aiPlayer) designStarbases() error {
	starbasePurposes := []cs.ShipDesignPurpose{
		cs.ShipDesignPurposeFuelDepot,
		cs.ShipDesignPurposeFort,
		cs.ShipDesignPurposeStarbaseUnarmed,
		cs.ShipDesignPurposeStarbaseQuarter,
		cs.ShipDesignPurposeStarbaseHalf,
		cs.ShipDesignPurposeStarbase,
	}

	for _, purpose := range starbasePurposes {
		design, err := ai.designShip(ai.config.namesByPurpose[purpose], purpose, cs.FleetPurposeFromShipDesignPurpose(purpose))
		if err != nil {
			return fmt.Errorf("unable to design ship %v: %w", purpose, err)
		}
		ai.designsByPurpose[purpose] = design
	}

	return nil
}

func (ai *aiPlayer) removeUnusedDesigns() {
	var unusedDesigns map[int]bool = map[int]bool{}

	// find any designs with no instances
	for _, design := range ai.Designs {
		if design.Spec.NumInstances == 0 && !design.CannotDelete {
			// ai.log.Debug().
			// 	Int64("GameID", ai.GameID).
			// 	Int("PlayerNum", ai.Num).
			// 	Msgf("marking %s for deletion, unused", design.Name)

			unusedDesigns[design.Num] = true
		}
	}

	// find any designs that may be queued on planets, don't delete those
	for _, planet := range ai.Planets {
		for _, item := range planet.ProductionQueue {
			if item.DesignNum != 0 {
				delete(unusedDesigns, item.DesignNum)

				// ai.log.Debug().
				// 	Int64("GameID", ai.GameID).
				// 	Int("PlayerNum", ai.Num).
				// 	Msgf("design %d still used, not marking for deletion", item.DesignNum)
			}
		}
	}

	for _, design := range ai.Designs {
		if found, found2 := unusedDesigns[design.Num]; found && found2 {
			// log a message if we're deleting an existing design
			if design.ID != 0 {
				ai.log.Debug().
					Int64("GameID", ai.GameID).
					Int("PlayerNum", ai.Num).
					Msgf("marking %s, design %d for deletion, unused", design.Name, design.Num)
			}
			design.Delete = true
		}
	}
}

// assign a purpose to any fleets that don't have it (for starter fleets)
func (ai *aiPlayer) assignPurpose() {
	for _, fleet := range ai.Fleets {
		if fleet.Purpose == cs.FleetPurposeNone {
			design := ai.GetDesign(fleet.Tokens[0].DesignNum)
			fleet.Purpose = cs.FleetPurposeFromShipDesignPurpose(design.Purpose)
		}
	}
}

// get the best hull we can build by iterating through the list backwards
// TODO: Add ability to use different search criteria (cheapness, rating, etc.) & not be dependent
// on the best techs being at the back of the list
func (ai *aiPlayer) getBestHull(hulls []*cs.TechHull) *cs.TechHull {
	// iterate over hulls backwards. better hulls are later so start there
	for i := len(hulls) - 1; i >= 0; i-- {
		hull := hulls[i]
		if ai.HasTech(&hull.Tech) {
			return hull
		}
	}
	return nil
}
