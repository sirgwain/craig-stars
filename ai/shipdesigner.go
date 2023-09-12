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
	case cs.ShipDesignPurposeFreighter:
		fallthrough
	case cs.ShipDesignPurposeColonistFreighter:
		fallthrough
	case cs.ShipDesignPurposeFuelFreighter:
		hull = ai.getBestHull(ai.techStore.GetHullsByType(cs.TechHullTypeFreighter))
	case cs.ShipDesignPurposeFighter:
		hull = ai.getBestHull(ai.techStore.GetHullsByType(cs.TechHullTypeFighter))
	case cs.ShipDesignPurposeBomber:
		hull = ai.getBestHull(ai.techStore.GetHullsByType(cs.TechHullTypeBomber))
	case cs.ShipDesignPurposeStarbase:
		fallthrough
	case cs.ShipDesignPurposePacketThrower:
		fallthrough
	case cs.ShipDesignPurposeStargater:
		fallthrough
	case cs.ShipDesignPurposeFuelDepot:
		hull = ai.getBestHull(ai.techStore.GetHullsByType(cs.TechHullTypeStarbase))
	case cs.ShipDesignPurposeFort:
		fallthrough
	case cs.ShipDesignPurposeStarterColony:
		hull = ai.getBestHull(ai.techStore.GetHullsByType(cs.TechHullTypeStarbase))
	}

	if hull == nil {
		return existing, nil
	}

	updated = cs.DesignShip(ai.techStore, hull, name, ai.Player, ai.GetNextDesignNum(ai.Designs), ai.DefaultHullSet, purpose, fleetPurpose)
	updated.HullSetNumber = ai.DefaultHullSet
	updated.Purpose = purpose
	updated.Spec = cs.ComputeShipDesignSpec(&ai.game.Rules, ai.TechLevels, ai.Race.Spec, updated)

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

	// if our existing design is equivalent, return it
	if found && existing.SlotsEqual(updated) {
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

// get the best hull we can build by iterating through the list backwards
func (ai *aiPlayer) getBestHull(hulls []*cs.TechHull) *cs.TechHull {
	for i := len(hulls) - 1; i >= 0; i-- {
		hull := hulls[i]
		if ai.HasTech(&hull.Tech) {
			return hull
		}
	}
	return nil
}
