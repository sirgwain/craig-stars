package cs

import (
	"fmt"
	"slices"
)

// A StartingFleet represents a "blueprint" for a race's starting ship designs,
// used to create player starting designs and fleets during universe generation.
type StartingFleet struct {
	HullType         TechHullType      `json:"hullType,omitempty"`         // hull type of designed ship
	Type             StartingFleetType `json:"type,omitempty"`             // fleet type; denotes kind and order of parts to place
	UsesCheapestHull bool              `json:"usesCheapestHull,omitempty"` // whether the design uses the cheapest or highest tech hull; used for ARM/HE
	name             string            // used to store design name so we can find it later
}

// Fleet type used during starting fleet design to fill slots
type StartingFleetType string

const (
	StartingFleetTypeBomber           StartingFleetType = "Bomber"           // Bombers
	StartingFleetTypeCloakedFreighter StartingFleetType = "CloakedFreighter" // SS cloaked freighters
	StartingFleetTypeColonizer        StartingFleetType = "Colonizer"        // (Mini-)Colony ships
	StartingFleetTypeFighter          StartingFleetType = "Fighter"          // Destroyers, armed scouts and normal freighters
	StartingFleetTypeMineLayer        StartingFleetType = "MineLayer"        // SD standard minelayer
	StartingFleetTypeMiner            StartingFleetType = "Miner"            // Remote miners
	StartingFleetTypeSpeedMineLayer   StartingFleetType = "SpeedMineLayer"   // SD Speed Trap minelayer
	StartingFleetTypeScout            StartingFleetType = "Scout"            // Unarmed scouts
	StartingFleetTypeTerraformer      StartingFleetType = "Terraformer"      // CA terraformers
)

// Return the default name and hullSetNumber for this StartingFleet, using its Spec to resolve naming disputes.
// Spec is assumed to be non-nil.
func (sf *StartingFleet) getDefaultName(spec *ShipDesignSpec) (name string, hullSetNumber int) {
	switch sf.Type {
	case StartingFleetTypeBomber:
		return "Gadfly", 1
	case StartingFleetTypeCloakedFreighter:
		return "Shadow Transport", 1
	case StartingFleetTypeMineLayer:
		return "Little Hen", 0
	case StartingFleetTypeSpeedMineLayer:
		return "Speed Turtle", 0
	case StartingFleetTypeTerraformer:
		return "Change of Heart", 2
	case StartingFleetTypeColonizer:
		if sf.UsesCheapestHull {
			return "Spore Cloud", 0
		}
		if spec.OrbitalConstructionModule {
			return "Pinta", 1
		}
		return "Santa Maria", 0
	case StartingFleetTypeFighter:
		switch sf.HullType {
		case TechHullTypeScout:
			return "Armed Probe", 1
		case TechHullTypeFighter:
			return "Stalwart Defender", 0
		default:
			if spec.HasWeapons {
				return "Swashbuckler", 0
			}
			return "Teamster", 0
		}
	case StartingFleetTypeMiner:
		if sf.UsesCheapestHull {
			return "Potato Bug", 2
		}
		return "Cotton Picker", 1
	case StartingFleetTypeScout:
		if spec.CloakUnits > 0 {
			return "Shadow Sleuth", 3
		}
		return "Long Range Scout", 0
	}
	return "BRRR SKIBIDI DOP DOP DOP DOP YES YES YES", 3
}

// the startingFleeter struct creates starting designs for starting fleets.
type startingFleeter struct {
	rules  *Rules
	player *Player
	// cache containing all usable tech items of all types;
	// speeds up execution by avoiding looping over unusable techs
	// and is saved per instance
	usableTechs map[TechItemType][]TechItem
}

func newStartingFleeter(rules *Rules, player *Player) (sf *startingFleeter) {
	return &startingFleeter{
		rules:       rules,
		player:      player,
		usableTechs: rules.techs.GetUsableTechs(player.TechLevels),
	}
}

// createStartingDesign creates starting ship designs for a player.
// This modifies the StartingFleet to contain the newly created design name for storage.
func (sf *startingFleeter) createStartingDesign(startingFleet *StartingFleet, num int) (design *ShipDesign, err error) {
	// Basic setup
	player := sf.player

	hull := sf.getStartingHull(startingFleet)
	if hull == nil {
		// no hull means we don't create this design
		// TODO: Should we log this?
		return nil, nil
	}

	engine := sf.getStartingEngine(startingFleet.Type, hull.Name)
	if engine == nil {
		// no engine; report error
		engines := sf.usableTechs[TechItemTypeEngine]
		return nil, fmt.Errorf("no usable engines found when creating starting design; Engines: %v", engines)
	}

	design = NewShipDesign(player.Num, num).WithHull(hull.Name)

	numBeamWeapons := 0
	numTorpedoes := 0
	var hasScanner bool

	// ordered list of fleet parts to use
	var fleetParts []TechTag
	switch startingFleet.Type {
	case StartingFleetTypeBomber:
		fleetParts = []TechTag{TechTagBomb, TechTagShield}
	case StartingFleetTypeCloakedFreighter:
		fleetParts = []TechTag{TechTagCloak, TechTagShield}
	case StartingFleetTypeColonizer:
		fleetParts = []TechTag{TechTagColonyModule}
	case StartingFleetTypeFighter:
		fleetParts = []TechTag{TechTagScanner, TechTagBeamWeapon, TechTagTorpedo, TechTagArmor, TechTagTorpedoBonus}
	case StartingFleetTypeMineLayer:
		fleetParts = []TechTag{TechTagMineLayer, TechTagScanner}
	case StartingFleetTypeMiner:
		fleetParts = []TechTag{TechTagMiningRobot, TechTagScanner}
	case StartingFleetTypeScout:
		fleetParts = []TechTag{TechTagScanner, TechTagCloak}
	case StartingFleetTypeSpeedMineLayer:
		fleetParts = []TechTag{TechTagSpeedMineLayer, TechTagScanner}
	case StartingFleetTypeTerraformer:
		fleetParts = []TechTag{TechTagTerraforming, TechTagScanner}
	}

	partCache := newPartCache(func(hst HullSlotType, tag TechTag) *TechHullComponent {
		return sf.getStartingPart(hull.Name, hst, tag)
	})

	// fill slots consecutively
	for i, hullSlot := range hull.Slots {
		slot := ShipDesignSlot{HullSlotIndex: i + 1, Quantity: hullSlot.Capacity}
		hst := hullSlot.Type
		// Handle engine slots first
		if hst&HullSlotTypeEngine != 0 {
			slot.HullComponent = engine.Name
			design.Slots = append(design.Slots, slot)
			continue
		}

		// check tags one by one in order to see what we can add
		for _, tag := range fleetParts {
			// cover edge cases
			switch tag {
			case TechTagScanner:
				if hasScanner {
					// already have scanner; skip
					continue
				}

				if scanner := partCache.get(hst, TechTagScanner); scanner != nil {
					slot.HullComponent = scanner.Name
					hasScanner = true
				}
			case TechTagTorpedo, TechTagBeamWeapon:
				// check guns
				beamWeapon := partCache.get(hst, TechTagBeamWeapon)
				torpedo := partCache.get(hst, TechTagTorpedo)
				switch {
				case beamWeapon != nil && torpedo != nil:
					// if both exist, try to balance between beams and torps
					if numTorpedoes > numBeamWeapons {
						slot.HullComponent = beamWeapon.Name
						numBeamWeapons += slot.Quantity
					} else {
						slot.HullComponent = torpedo.Name
						numTorpedoes += slot.Quantity
					}
				case beamWeapon != nil:
					// torpedo unavailable; use beam weapon
					slot.HullComponent = beamWeapon.Name
					numBeamWeapons += slot.Quantity
				case torpedo != nil:
					// beam unavailable; use torpedo
					slot.HullComponent = torpedo.Name
					numTorpedoes += slot.Quantity
				default:
					// No weapons available for this slot; move on
				}
			default:
				if part := partCache.get(hst, tag); part != nil {
					slot.HullComponent = part.Name
				}
			}

			if slot.HullComponent != "" {
				// we found a component for this slot
				break
			}
		}

		if slot.HullComponent == "" {
			// check for fuel tanks as a last resort (and to fill in empty mech slots)
			if part := partCache.get(hst, TechTagFuelTank); part != nil {
				slot.HullComponent = part.Name
			}
		}

		if slot.HullComponent != "" {
			// tack on new slot to design and move on
			design.Slots = append(design.Slots, slot)
		}

	}

	// compute design name/hull set if not previously modified
	if startingFleet.name == "" {
		defaultName, defaultHullSet := startingFleet.getDefaultName(&design.Spec)
		startingFleet.name = defaultName
		design.Name = defaultName
		design.HullSetNumber = defaultHullSet
	}

	// final wrapup
	design.Spec, err = ComputeShipDesignSpec(sf.rules, player.TechLevels, player.Race.Spec, design)
	if err != nil {
		return nil, fmt.Errorf("computeShipDesignSpec() errored during createStartingDesigns: \n%w", err)
	}
	return design, nil
}

// getStartingHull returns the hull to use for the given StartingFleet.
func (sf *startingFleeter) getStartingHull(startingFleet *StartingFleet) (hull *TechHull) {
	player := sf.player

	usableHulls := sf.usableTechs[TechItemTypeHull]
	var cheapestHull *TechHull

	// loop over hulls in reverse and grab the first one we find that works
	for i := len(usableHulls) - 1; i >= 0; i-- {
		hull, ok := usableHulls[i].(*TechHull)
		if !ok {
			panic(fmt.Sprintf("startingFleeter part cache contained non-hull inside hull cache; \n%v", usableHulls[i]))
		}

		if !player.CanLearnTech(&hull.Tech) || !player.HasAcquiredTech(&hull.Tech) {
			// tech is incompatible or unacquired; skip
			continue
		}

		var correctHullType bool
		// cover edge special cases for certain hull types
		switch startingFleet.HullType {
		case TechHullTypeFighter:
			correctHullType = hull.Type == TechHullTypeFighter || hull.Type == TechHullTypeCapitalShip
		case TechHullTypeFreighter:
			correctHullType = hull.Type == TechHullTypeFreighter || hull.Type == TechHullTypeMultiPurposeFreighter
		default:
			correctHullType = hull.Type == startingFleet.HullType
		}

		if !correctHullType {
			// incorrect hull type; skip
			continue
		}

		if !startingFleet.UsesCheapestHull {
			return hull // normally, return the first valid hull found
		}

		// for cheapest mode, track the hull with the lowest cost
		if cheapestHull == nil || hull.Cost.Total() < cheapestHull.Cost.Total() {
			cheapestHull = hull
		}
	}

	return cheapestHull
}

// getStartingEngine returns the proper engine for a starting fleet.
// It may not (and in fact, often isn't) the best choice,
// but at least it's accurate dammit!
func (sf *startingFleeter) getStartingEngine(startingFleetType StartingFleetType, hullName string) (engine *TechEngine) {
	player := sf.player
	usableEngines := sf.usableTechs[TechItemTypeEngine]
	var firstEngine *TechEngine
	// iterate through engines backwards
	for i := len(usableEngines) - 1; i >= 0; i-- {
		eng, ok := usableEngines[i].(*TechEngine)
		if !ok {
			panic(fmt.Sprintf("startingFleeter part cache contained non-engine inside engine cache at position %d: \n%v", len(usableEngines)-i-1, usableEngines[i]))
		}

		if !player.CanLearnTech(&eng.Tech) || !player.HasAcquiredTech(&eng.Tech) || // tech is incompatible or unacquired; skip
			!eng.allowedOnHull(hullName) || // engine is incompatible with hull
			// engine would kill colonists for colonizers (and *ONLY* colonizers I may add)
			(eng.Radiating && startingFleetType == StartingFleetTypeColonizer &&
				// TODO: Refactor this once Radiating becomes a property of the TechHullComponent
				!(player.Race.ImmuneRad || player.Race.Spec.HabCenter.Rad >= sf.rules.RadiatingImmune)) {
			continue
		}

		// If engine has hull restrictions and this hull is allowed, use it
		// This gets around settler's delight for HE minicols
		if len(eng.Requirements.HullsAllowed) > 0 &&
			slices.Contains(eng.Requirements.HullsAllowed, hullName) {
			return eng
		}

		// otherwise, keep track of the first one we found and default to that
		// if no exclusive engines are found
		if firstEngine == nil {
			firstEngine = eng
		}
	}
	return firstEngine
}

// getStartingPart returns the best starting part for a given type and hull slot, or nil if none are found.
//
// Unlike most starting design functions, this one actually chooses the best part available regardless of order.
func (sf *startingFleeter) getStartingPart(hullName string, hst HullSlotType, tag TechTag) (bestTech *TechHullComponent) {
	player := sf.player
	usableComps := sf.usableTechs[TechItemTypeHullComponent]
	tc := newTechComparerWithoutCost(sf.rules, player)

	for i := range usableComps {
		hc, ok := usableComps[i].(*TechHullComponent)
		if !ok {
			panic(fmt.Sprintf("startingFleeter part cache contained non-hull component inside hull component cache at position %d: \n%s", i, usableComps[i]))
		}

		if !player.CanLearnTech(&hc.Tech) || // tech is unavailable
			!player.HasAcquiredTech(&hc.Tech) || // tech has not been acquired
			hst&hc.HullSlotType == 0 || // wrong slot type
			!hc.Tags.HasTag(tag) || // wrong type of part
			!hc.allowedOnHull(hullName) { // part is incompatible with hull
			continue
		}

		if tc.compareFieldsByTag(bestTech, hc, tag) {
			// this is better than our prior part; tack it on
			bestTech = hc
		}
	}

	return bestTech
}
