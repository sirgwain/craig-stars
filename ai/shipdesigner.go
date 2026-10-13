package ai

import (
	"fmt"
	"log/slog"

	"github.com/sirgwain/craig-stars/cs"
)

// get a ship design for a purpose
func (ai *aiPlayer) designShip(name string, purpose cs.ShipDesignPurpose, fleetPurpose cs.FleetPurpose) (updated *cs.ShipDesign, err error) {

	existing := ai.designsByPurpose[purpose]
	found := existing != nil
	if found && existing.ID == 0 {
		// we already created a new design for this purpose, just use it
		return existing, nil
	}

	var hull *cs.TechHull
	switch purpose {
	case cs.ShipDesignPurposeScout:
		hull = ai.getBestHull(ai.techStore.GetHullsByType(cs.TechHullTypeScout))
	case cs.ShipDesignPurposeMiner:
		hull = ai.getBestHull(ai.techStore.GetHullsByType(cs.TechHullTypeMiner))
	case cs.ShipDesignPurposeDamageMineLayer:
		hull = ai.getBestHull(ai.techStore.GetHullsByType(cs.TechHullTypeMineLayer))
		if hull == nil && ai.HasTech(&cs.Frigate.Tech) {
			hull = &cs.Frigate
		}
	case cs.ShipDesignPurposeColonizer:
		hull = ai.getBestHull(ai.techStore.GetHullsByType(cs.TechHullTypeColonizer))
		// HE and AR colony modules only fit their own colony hulls
		if ai.Race.PRT != cs.HE && ai.Race.PRT != cs.AR {
			if privateer := ai.techStore.GetHull(cs.Privateer.Name); privateer != nil && ai.HasTech(&privateer.Tech) {
				hull = privateer
			}
		}
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

	layout := purpose
	if purpose == cs.ShipDesignPurposeScout && ai.doctrine.cloakedScouts {
		layout = cs.ShipDesignPurposeCloakedScout
	}
	updated, err = cs.DesignShip(&ai.game.Rules, ai.Player, hull, name, ai.GetNextDesignNum(ai.Designs), ai.DefaultHullSet, layout, fleetPurpose)
	if err != nil {
		return existing, fmt.Errorf("cs.DesignShip returned error: %w", err)
	}
	// a cloaked scout still does a scout's job
	updated.Purpose = purpose
	// no need to compute ship design specs as functions already compute it before returning

	// if we tried to build a bomber with no bombs, ignore it
	if (purpose == cs.ShipDesignPurposeBomber && !updated.Spec.Bomber) || (purpose == cs.ShipDesignPurposeColonizer && !updated.Spec.Colonizer) || (purpose == cs.ShipDesignPurposeMiner && updated.Spec.MiningRate == 0) || (purpose == cs.ShipDesignPurposeDamageMineLayer && len(updated.Spec.MineLayingRateByMineType) == 0) {
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
	var err error

	purpose := cs.ShipDesignPurposeFuelDepot
	ai.fuelDepotDesign, err = ai.designShip(ai.config.namesByPurpose[purpose], purpose, cs.FleetPurposeFromShipDesignPurpose(purpose))
	if err != nil {
		return fmt.Errorf("unable to design ship %v: %w", purpose, err)
	}

	purpose = cs.ShipDesignPurposeFort
	ai.fortDesign, err = ai.designShip(ai.config.namesByPurpose[purpose], purpose, cs.FleetPurposeFromShipDesignPurpose(purpose))
	if err != nil {
		return fmt.Errorf("unable to design ship %v: %w", purpose, err)
	}

	purpose = cs.ShipDesignPurposeStarbaseUnarmed
	ai.starbaseUnarmedDesign, err = ai.designShip(ai.config.namesByPurpose[purpose], purpose, cs.FleetPurposeFromShipDesignPurpose(purpose))
	if err != nil {
		return fmt.Errorf("unable to design ship %v: %w", purpose, err)
	}

	purpose = cs.ShipDesignPurposeStarbaseQuarter
	ai.starbaseQuarterDesign, err = ai.designShip(ai.config.namesByPurpose[purpose], purpose, cs.FleetPurposeFromShipDesignPurpose(purpose))
	if err != nil {
		return fmt.Errorf("unable to design ship %v: %w", purpose, err)
	}

	purpose = cs.ShipDesignPurposeStarbaseHalf
	ai.starbaseHalfDesign, err = ai.designShip(ai.config.namesByPurpose[purpose], purpose, cs.FleetPurposeFromShipDesignPurpose(purpose))
	if err != nil {
		return fmt.Errorf("unable to design ship %v: %w", purpose, err)
	}

	purpose = cs.ShipDesignPurposeStarbase
	ai.starbaseDesign, err = ai.designShip(ai.config.namesByPurpose[purpose], purpose, cs.FleetPurposeFromShipDesignPurpose(purpose))
	if err != nil {
		return fmt.Errorf("unable to design ship %v: %w", purpose, err)
	}

	if ai.doctrine.outposts {
		// the cheap orbital fort hull can carry the one part an outpost is for
		purpose = cs.ShipDesignPurposeStargater
		if ai.Race.PRT == cs.PP {
			purpose = cs.ShipDesignPurposePacketThrower
		}
		ai.outpostDesign, err = ai.designShip(ai.config.namesByPurpose[purpose], purpose, cs.FleetPurposeFromShipDesignPurpose(purpose))
		if err != nil {
			return fmt.Errorf("unable to design ship %v: %w", purpose, err)
		}
		if ai.outpostDesign != nil && ai.outpostDesign.Spec.Stargate == "" && ai.outpostDesign.Spec.MassDriver == "" {
			// no gate or driver to put on it yet
			ai.outpostDesign = nil
		}
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
				ai.log.Debug("marking design for deletion, unused",
					slog.String("design", design.Name),
					slog.Int("designNum", design.Num),
					slog.Int64("GameID", ai.GameID),
					slog.Int("PlayerNum", ai.Num))
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
