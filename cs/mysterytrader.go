package cs

import (
	"fmt"
	"math"
)

// The mystery trader travels through space and gives a boon to any player that gives it a fleet
// full of minerals
// TODO: not yet implemented
type MysteryTrader struct {
	MapObject
	WarpSpeed     int                     `json:"warpSpeed,omitempty"`
	Destination   Vector                  `json:"destination"`
	RequestedBoon int                     `json:"requestedBoon,omitempty"`
	RewardType    MysteryTraderRewardType `json:"rewardType"`
	Heading       Vector                  `json:"heading,omitempty"`
	Spec          MysteryTraderSpec       `json:"spec,omitempty"`
}

type MysteryTraderSpec struct {
}

type MysteryTraderRewardType string

const (
	MysteryTraderRewardNone       MysteryTraderRewardType = ""
	MysteryTraderRewardResearch   MysteryTraderRewardType = "Research"
	MysteryTraderRewardEngine     MysteryTraderRewardType = "Engine"
	MysteryTraderRewardBomb       MysteryTraderRewardType = "Bomb"
	MysteryTraderRewardArmor      MysteryTraderRewardType = "Armor"
	MysteryTraderRewardShield     MysteryTraderRewardType = "Shield"
	MysteryTraderRewardElectrical MysteryTraderRewardType = "Electrical"
	MysteryTraderRewardMechanical MysteryTraderRewardType = "Mechanical"
	MysteryTraderRewardTorpedo    MysteryTraderRewardType = "Torpedo"
	MysteryTraderRewardMineRobot  MysteryTraderRewardType = "MineRobot"
	MysteryTraderRewardShipHull   MysteryTraderRewardType = "ShipHull"
	MysteryTraderRewardBeamWeapon MysteryTraderRewardType = "BeamWeapon"
	MysteryTraderRewardGenesis    MysteryTraderRewardType = "Genesis"
	MysteryTraderRewardJumpgate   MysteryTraderRewardType = "Jumpgate"
	MysteryTraderRewardLifeboat   MysteryTraderRewardType = "Lifeboat"
)

var Rewards = [15]MysteryTraderRewardType{
	MysteryTraderRewardNone,
	MysteryTraderRewardResearch,
	MysteryTraderRewardEngine,
	MysteryTraderRewardBomb,
	MysteryTraderRewardArmor,
	MysteryTraderRewardShield,
	MysteryTraderRewardElectrical,
	MysteryTraderRewardMechanical,
	MysteryTraderRewardTorpedo,
	MysteryTraderRewardBeamWeapon,
	MysteryTraderRewardMineRobot,
	MysteryTraderRewardShipHull,
	MysteryTraderRewardGenesis,
	MysteryTraderRewardJumpgate,
	MysteryTraderRewardLifeboat,
}

type MysteryTraderReward struct {
	Type       MysteryTraderRewardType `json:"type"`
	TechLevels TechLevel               `json:"techLevels"`
	Tech       string                  `json:"tech,omitempty"`
}

// create a new mysterytrader object
func newMysteryTrader(position Vector, num int, warpSpeed int, destination Vector, requestedBoon int, reward MysteryTraderRewardType) *MysteryTrader {
	return &MysteryTrader{
		MapObject: MapObject{
			Type:     MapObjectTypeMysteryTrader,
			Position: position,
			Num:      num,
			Name:     fmt.Sprintf("Mystery Trader #%d", num),
		},
		WarpSpeed:     warpSpeed,
		Destination:   destination,
		Heading:       (destination.Subtract(position)).Normalized(),
		RequestedBoon: requestedBoon,
		RewardType:    reward,
	}
}

func generateMysteryTrader(rules *Rules, game *Game, num int) *MysteryTrader {

	mtRules := rules.MysteryTraderRules

	turn := game.Year - rules.StartingYear

	if turn < mtRules.MinYear {
		return nil
	}

	// only generate mystery traders on even years
	if mtRules.EvenYearOnly && turn&1 > 0 {
		return nil
	}

	// every turn has a different chance of generating a random number generator
	chance := mtRules.Chances[rules.random.Intn(len(mtRules.Chances))]
	if rules.random.Intn(chance) != 0 {
		return nil
	}

	// We made it this far, WE HAVE A MYSTERY TRADER!

	// determine speed
	warp := mtRules.MinWarp + rules.random.Intn(mtRules.MaxWarp-mtRules.MinWarp)

	// we want two sets of coords, a random set somewhere on the x/y range
	// and whatever our edge coords are
	var randomXYCoords, yEdgeCoords, xEdgeCoords [2]int

	randomXYCoords = [2]int{
		20 + rules.random.Intn(int(game.Area.X)-39),
		20 + rules.random.Intn(int(game.Area.Y)-39),
	}

	xEdgeCoords = [2]int{
		-20,
		int(game.Area.X) + 20,
	}

	yEdgeCoords = [2]int{
		-20,
		int(game.Area.Y) + 20,
	}

	// our position/dest always uses one edge coord, one random coord
	// the edge coord is either 0 or the max size, based on some randomness
	var position, destination Vector
	if rules.random.Intn(2) == 0 {
		position.X = float64(randomXYCoords[0])
		position.Y = float64(yEdgeCoords[rules.random.Intn(2)])
	} else {
		position.X = float64(xEdgeCoords[rules.random.Intn(2)])
		position.Y = float64(randomXYCoords[1])
	}

	// set our destination to maximize the chances it goes through the universe
	if position.X > game.Area.X/2 {
		destination.X = -20
	} else {
		destination.X = game.Area.X + 20
	}

	if position.Y > game.Area.Y/2 {
		destination.Y = -20
	} else {
		destination.Y = game.Area.Y + 20
	}

	// swap position/destination at random
	if rules.random.Intn(2) == 0 {
		position, destination = destination, position
	}

	// TODO: populate reward
	reward := MysteryTraderRewardResearch

	return newMysteryTrader(position, num, warp, destination, mtRules.RequestedBoon, reward)
}

// move a mystery trader
func (mt *MysteryTrader) move() {
	totalDist := mt.Position.DistanceTo(mt.Destination)
	dist := float64(mt.WarpSpeed * mt.WarpSpeed)

	// make sure we end up at a whole number
	vectorTravelled := mt.Destination.Subtract(mt.Position).Normalized().Scale(dist)
	dist = vectorTravelled.Length()
	// don't overshoot
	dist = math.Min(totalDist, dist)

	if totalDist == dist {
		mt.Position = mt.Destination
		mt.Delete = true
	} else {
		// move along the heading...
		mt.Heading = (mt.Destination.Subtract(mt.Position)).Normalized()
		mt.Position = mt.Position.Add(mt.Heading.Scale(dist))
		mt.Position = mt.Position.Round()
	}
}

// meet a mystery trader and recieve a reward!
func (mt *MysteryTrader) meet(rules *Rules, fleet *Fleet, player *Player) MysteryTraderReward {
	gift := fleet.Cargo.ToMineral().Total()
	if fleet.Cargo.ToMineral().Total() >= mt.RequestedBoon {
		// it's a major award!

		switch mt.RewardType {
		// give the player tech levels
		case MysteryTraderRewardResearch:
			numLevels := player.TechLevels.Sum()
			levels := TechLevel{}

			for _, techLevelReward := range rules.MysteryTraderRules.TechBoon {
				if numLevels <= techLevelReward.TechLevels {
					// player tech level count is below threshold for this reward, grant some random levels
					for _, reward := range techLevelReward.Rewards {
						if gift <= reward.MineralsGiven {
							// minerals given is below this threshold, give this number of levels
							for i := 0; i < reward.Reward; i++ {
								availableFields := levels.LearnableTechFields(rules)
								field := availableFields[rules.random.Intn(len(availableFields))]
								levels.Set(field, levels.Get(field)+1)
							}
							break
						}
					}
					break
				}
			}
			return MysteryTraderReward{
				Type:       MysteryTraderRewardResearch,
				TechLevels: levels,
			}
		}
	}
	return MysteryTraderReward{}
}

var HushABoom = TechHullComponent{Tech: NewTechWithOrigin("Hush-a-Boom", NewCost(1, 2, 0, 2), TechRequirements{Acquirable: true, TechLevel: TechLevel{Weapons: 12, Electronics: 12, Biotechnology: 12}}, 75, TechCategoryBomb, OriginMysteryTrader),
	Mass:                 45,
	KillRate:             3,
	StructureDestroyRate: 2.0,
	HullSlotType:         HullSlotTypeBomb,
}

var EnigmaPulsar = TechEngine{
	TechHullComponent: TechHullComponent{Tech: NewTechWithOrigin("Enigma Pulsar", NewCost(12, 15, 11, 40), TechRequirements{Acquirable: true, TechLevel: TechLevel{Energy: 7, Propulsion: 13, Construction: 5, Electronics: 9}}, 85, TechCategoryEngine, OriginMysteryTrader),
		Mass:          20,
		HullSlotType:  HullSlotTypeEngine,
		MovementBonus: 1,
		CloakUnits:    20,
	},
	Engine: Engine{
		IdealSpeed:   10,
		FreeSpeed:    1,
		MaxSafeSpeed: 10,
		FuelUsage: [11]int{
			0,
			0,
			0,
			0,
			0,
			0,
			60,
			70,
			80,
			90,
			100,
		},
	},
}
var MegaPolyShell = TechHullComponent{Tech: NewTechWithOrigin("Mega Poly Shell", NewCost(14, 5, 5, 52), TechRequirements{Acquirable: true, TechLevel: TechLevel{Energy: 14, Construction: 14, Electronics: 14, Biotechnology: 6}}, 95, TechCategoryArmor, OriginMysteryTrader),

	Mass:           20,
	Shield:         100,
	Armor:          400,
	CloakUnits:     40,
	TorpedoJamming: .2,
	ScanRange:      80,
	ScanRangePen:   40,
	HullSlotType:   HullSlotTypeShield,
}
var LangstonShell = TechHullComponent{Tech: NewTechWithOrigin("Langston Shell", NewCost(6, 1, 4, 12), TechRequirements{Acquirable: true, TechLevel: TechLevel{Energy: 12, Propulsion: 9, Electronics: 9}}, 65, TechCategoryShield, OriginMysteryTrader),

	Mass:           10,
	Shield:         125,
	Armor:          65,
	CloakUnits:     20,
	TorpedoJamming: .05,
	ScanRange:      50,
	ScanRangePen:   25,
	HullSlotType:   HullSlotTypeShield,
}
var MultiFunctionPod = TechHullComponent{Tech: NewTechWithOrigin("Multi-Function Pod", NewCost(5, 0, 5, 15), TechRequirements{Acquirable: true, TechLevel: TechLevel{Energy: 11, Propulsion: 11, Electronics: 11}}, 35, TechCategoryElectrical, OriginMysteryTrader),

	Mass:           2,
	CloakUnits:     60,
	TorpedoJamming: .1,
	MovementBonus:  1,
	HullSlotType:   HullSlotTypeElectrical,
}
var AntiMatterTorpedo = TechHullComponent{Tech: NewTechWithOrigin("Anti Matter Torpedo", NewCost(3, 8, 1, 50), TechRequirements{Acquirable: true, TechLevel: TechLevel{Weapons: 11, Propulsion: 12, Biotechnology: 21}}, 65, TechCategoryTorpedo, OriginMysteryTrader),

	Mass:         8,
	Initiative:   0,
	Accuracy:     85,
	Power:        60,
	Range:        6,
	HullSlotType: HullSlotTypeWeapon,
}
var JumpGate = TechHullComponent{Tech: NewTechWithOrigin("Jump Gate", NewCost(0, 0, 38, 30), TechRequirements{Acquirable: true, TechLevel: TechLevel{Energy: 16, Propulsion: 20, Construction: 20, Electronics: 16}}, 75, TechCategoryMechanical, OriginMysteryTrader),
	Mass:         10,
	CanJump:      true, // TODO: add support for this
	HullSlotType: HullSlotTypeMechanical,
}
var MultiContainedMunition = TechHullComponent{Tech: NewTechWithOrigin("Multi Contained Munition", NewCost(5, 32, 5, 32), TechRequirements{Acquirable: true, TechLevel: TechLevel{Energy: 21, Weapons: 21, Electronics: 16, Biotechnology: 12}}, 175, TechCategoryBeamWeapon, OriginMysteryTrader),
	Mass:                 8,
	Initiative:           6,
	Power:                140,
	Range:                3,
	CloakUnits:           20,
	TorpedoBonus:         .1,
	ScanRange:            150,
	ScanRangePen:         75,
	KillRate:             2,
	StructureDestroyRate: 5,
	MineLayingRate:       40,
	HullSlotType:         HullSlotTypeWeapon,
}
var AlienMiner = TechHullComponent{Tech: NewTechWithOrigin("Alien Miner", NewCost(4, 0, 1, 10), TechRequirements{Acquirable: true, TechLevel: TechLevel{Energy: 5, Construction: 10, Electronics: 5, Biotechnology: 5}}, 55, TechCategoryMineRobot, OriginMysteryTrader),

	Mass:           20,
	MiningRate:     10,
	CloakUnits:     60,
	TorpedoJamming: .3,
	MovementBonus:  1,
	HullSlotType:   HullSlotTypeMining,
}
var MultiCargoPod = TechHullComponent{Tech: NewTechWithOrigin("Multi Cargo Pod", NewCost(12, 0, 3, 25), TechRequirements{Acquirable: true, TechLevel: TechLevel{Energy: 5, Construction: 11, Electronics: 5}}, 35, TechCategoryMechanical, OriginMysteryTrader),
	Mass:         9,
	CargoBonus:   250,
	Armor:        50,
	CloakUnits:   20,
	HullSlotType: HullSlotTypeMechanical,
}
var MiniMorph = TechHull{Tech: NewTechWithOrigin("Mini Morph", NewCost(30, 8, 8, 100), TechRequirements{Acquirable: true, TechLevel: TechLevel{Construction: 8}}, 305, TechCategoryShipHull, OriginMysteryTrader),
	Type:              TechHullTypeMultiPurposeFreighter,
	Mass:              70,
	Armor:             250,
	Initiative:        2,
	FuelCapacity:      400,
	CargoCapacity:     150,
	CargoSlotPosition: Vector{0, -0.5},
	CargoSlotSize:     Vector{1, 2},
	Slots: []TechHullSlot{
		{Position: Vector{-2, 0}, Type: HullSlotTypeEngine, Capacity: 2, Required: true},
		{Position: Vector{-1, 0}, Type: HullSlotTypeGeneral, Capacity: 3},
		{Position: Vector{1, -0.5}, Type: HullSlotTypeGeneral, Capacity: 1},
		{Position: Vector{1, 0.5}, Type: HullSlotTypeGeneral, Capacity: 1},
		{Position: Vector{2, 0}, Type: HullSlotTypeGeneral, Capacity: 1},
		{Position: Vector{-1, -1}, Type: HullSlotTypeGeneral, Capacity: 2},
		{Position: Vector{-1, 1}, Type: HullSlotTypeGeneral, Capacity: 2},
	},
}
var GenesisDevice = TechPlanetary{Tech: NewTechWithOrigin("Genesis Device", NewCost(10, 10, 70, 100), TechRequirements{Acquirable: true, TechLevel: TechLevel{Energy: 20, Weapons: 10, Propulsion: 10, Construction: 20, Electronics: 10, Biotechnology: 20}}, 45, TechCategoryPlanetaryDefense, OriginMysteryTrader),
	ResetPlanet: true,
}
