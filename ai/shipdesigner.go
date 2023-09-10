package ai

import (
	"fmt"

	"github.com/sirgwain/craig-stars/cs"
)

// get a ship design for a purpose
func (ai *aiPlayer) designShip(name string, purpose cs.ShipDesignPurpose) (updated *cs.ShipDesign, err error) {

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
	case cs.ShipDesignPurposeColonistFreighter:
		hull = ai.getBestHull(ai.techStore.GetHullsByType(cs.TechHullTypeFreighter))
	case cs.ShipDesignPurposeFighter:
		hull = ai.getBestHull(ai.techStore.GetHullsByType(cs.TechHullTypeFighter))
	}

	if hull == nil {
		return existing, nil
	}

	updated = cs.DesignShip(ai.techStore, hull, name, ai.Player, ai.GetNextDesignNum(ai.Designs), ai.DefaultHullSet, purpose)
	updated.HullSetNumber = ai.DefaultHullSet
	updated.Purpose = purpose
	updated.Spec = cs.ComputeShipDesignSpec(&ai.game.Rules, ai.TechLevels, ai.Race.Spec, updated)

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

func (ai *aiPlayer) getHull(purpose fleetPurpose) *cs.TechHull {
	switch purpose {
	case scout:
		return ai.getBestHull(ai.techStore.GetHullsByType(cs.TechHullTypeScout))
	case colonizer:
		return ai.getBestHull(ai.techStore.GetHullsByType(cs.TechHullTypeColonizer))
	case fighter:
		return ai.getBestHull(ai.techStore.GetHullsByType(cs.TechHullTypeFighter))
	case bomber:
		return ai.getBestHull(ai.techStore.GetHullsByType(cs.TechHullTypeBomber))
	case transport:
		return ai.getBestHull(ai.techStore.GetHullsByType(cs.TechHullTypeFreighter))
	}
	return nil
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
