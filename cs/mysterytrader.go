package cs

import (
	"fmt"
	"math"
)

// The mystery trader travels through space and gives a boon to any player that gives it a fleet full of minerals
type MysteryTrader struct {
	MapObject
	WarpSpeed       int                     `json:"warpSpeed,omitempty"`
	Destination     Vector                  `json:"destination"`
	RequestedBoon   int                     `json:"requestedBoon,omitempty"`
	RewardType      MysteryTraderRewardType `json:"rewardType"`
	Heading         Vector                  `json:"heading,omitempty"`
	PlayersRewarded map[int]bool            `json:"playersRewarded"`
	Spec            MysteryTraderSpec       `json:"spec,omitempty"`
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
	MysteryTraderRewardJumpGate   MysteryTraderRewardType = "JumpGate"
	MysteryTraderRewardLifeboat   MysteryTraderRewardType = "Lifeboat"
)

var MysteryTraderRewards = []MysteryTraderRewardType{
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
	MysteryTraderRewardJumpGate,
	MysteryTraderRewardLifeboat,
}

var MysteryTraderRewardParts = []MysteryTraderRewardType{
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
}

func (t MysteryTraderRewardType) IsPart() bool {
	switch t {
	case MysteryTraderRewardEngine:
		fallthrough
	case MysteryTraderRewardBomb:
		fallthrough
	case MysteryTraderRewardArmor:
		fallthrough
	case MysteryTraderRewardShield:
		fallthrough
	case MysteryTraderRewardElectrical:
		fallthrough
	case MysteryTraderRewardMechanical:
		fallthrough
	case MysteryTraderRewardTorpedo:
		fallthrough
	case MysteryTraderRewardMineRobot:
		fallthrough
	case MysteryTraderRewardShipHull:
		fallthrough
	case MysteryTraderRewardBeamWeapon:
		fallthrough
	case MysteryTraderRewardGenesis:
		fallthrough
	case MysteryTraderRewardJumpGate:
		return true
	}
	return false
}

func (t MysteryTraderRewardType) Category() TechCategory {
	switch t {
	case MysteryTraderRewardEngine:
		return TechCategoryEngine
	case MysteryTraderRewardBomb:
		return TechCategoryBomb
	case MysteryTraderRewardArmor:
		return TechCategoryArmor
	case MysteryTraderRewardShield:
		return TechCategoryShield
	case MysteryTraderRewardElectrical:
		return TechCategoryElectrical
	case MysteryTraderRewardMechanical:
		return TechCategoryMechanical
	case MysteryTraderRewardTorpedo:
		return TechCategoryTorpedo
	case MysteryTraderRewardMineRobot:
		return TechCategoryMineRobot
	case MysteryTraderRewardShipHull:
		return TechCategoryShipHull
	case MysteryTraderRewardBeamWeapon:
		return TechCategoryBeamWeapon
	case MysteryTraderRewardGenesis:
		return TechCategoryPlanetary
	case MysteryTraderRewardJumpGate:
		return TechCategoryMechanical
	}
	return TechCategoryNone
}

// MysteryTraderRewardTypeForTech converts a mystery trader tech into a reward type
// we use this if we give the player a random part, we want them to know what type of
// reward it is
func MysteryTraderRewardTypeForTech(tech *Tech) MysteryTraderRewardType {
	if tech.Origin != OriginMysteryTrader {
		// not a MT part
		return MysteryTraderRewardNone
	}

	// jump gate is mechanical, but it's special
	if tech.Name == JumpGate.Name {
		return MysteryTraderRewardJumpGate
	}

	switch tech.Category {
	case TechCategoryEngine:
		return MysteryTraderRewardEngine
	case TechCategoryBomb:
		return MysteryTraderRewardBomb
	case TechCategoryArmor:
		return MysteryTraderRewardArmor
	case TechCategoryShield:
		return MysteryTraderRewardShield
	case TechCategoryElectrical:
		return MysteryTraderRewardElectrical
	case TechCategoryMechanical:
		return MysteryTraderRewardMechanical
	case TechCategoryTorpedo:
		return MysteryTraderRewardTorpedo
	case TechCategoryMineRobot:
		return MysteryTraderRewardMineRobot
	case TechCategoryShipHull:
		return MysteryTraderRewardShipHull
	case TechCategoryBeamWeapon:
		return MysteryTraderRewardBeamWeapon
	case TechCategoryPlanetary:
		return MysteryTraderRewardGenesis
	}
	return MysteryTraderRewardNone
}

type MysteryTraderReward struct {
	Type       MysteryTraderRewardType `json:"type"`
	TechLevels TechLevel               `json:"techLevels"`
	Tech       string                  `json:"tech,omitempty"`
	Ship       ShipDesign              `json:"ship,omitempty"`
	ShipCount  int                     `json:"shipCount,omitempty"`
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
		WarpSpeed:       warpSpeed,
		Destination:     destination,
		Heading:         (destination.Subtract(position)).Normalized(),
		PlayersRewarded: map[int]bool{},
		RequestedBoon:   requestedBoon,
		RewardType:      reward,
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
	chance := mtRules.ChanceSpawn[rules.random.Intn(len(mtRules.ChanceSpawn))]
	if rules.random.Intn(chance) != 0 {
		return nil
	}

	// We made it this far, WE HAVE A MYSTERY TRADER!

	// determine speed
	warpSpeed := mtRules.MinWarp + rules.random.Intn(mtRules.MaxWarp-mtRules.MinWarp)

	// determine where on the edge of the map we start
	position := generateRandomMysteryTraderCoords(rules, game)
	destination := generateRandomMysteryTraderDestination(rules, game, position)

	// swap position/destination at random
	if rules.random.Intn(2) == 0 {
		position, destination = destination, position
	}

	// generate a random reward type
	reward := generateMysteryTraderReward(rules, game.Year, warpSpeed)

	return newMysteryTrader(position, num, warpSpeed, destination, mtRules.RequestedBoon, reward)
}

// generateRandomMysteryTraderCoords gets random coords on the edge of the universe
func generateRandomMysteryTraderCoords(rules *Rules, game *Game) (coords Vector) {
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
	if rules.random.Intn(2) == 0 {
		coords.X = float64(randomXYCoords[0])
		coords.Y = float64(yEdgeCoords[rules.random.Intn(2)])
	} else {
		coords.X = float64(xEdgeCoords[rules.random.Intn(2)])
		coords.Y = float64(randomXYCoords[1])
	}
	return coords
}

// generateRandomMysteryDestination gets random coords on the edge of the universe
func generateRandomMysteryTraderDestination(rules *Rules, game *Game, position Vector) (destination Vector) {

	// start with a random dest
	randCoords := Vector{
		float64(20 + rules.random.Intn(int(game.Area.X)-39)),
		float64(20 + rules.random.Intn(int(game.Area.Y)-39)),
	}

	var yEdgeCoords, xEdgeCoords [2]int
	xEdgeCoords = [2]int{
		-20,
		int(game.Area.X) + 20,
	}

	yEdgeCoords = [2]int{
		-20,
		int(game.Area.Y) + 20,
	}

	// make destinations on each side of the map, with the other coord being random on the other axis
	// i.e. right is the far right edge of the map, and some random point on the Y
	left := Vector{float64(xEdgeCoords[0]), randCoords.Y}
	right := Vector{float64(xEdgeCoords[1]), randCoords.Y}
	top := Vector{randCoords.X, float64(yEdgeCoords[0])}
	bottom := Vector{randCoords.X, float64(yEdgeCoords[1])}

	// find the edge that is farthest from our current point and go there
	maxDist := -1.
	for _, dest := range [4]Vector{left, right, top, bottom} {
		dist := position.DistanceSquaredTo(dest)
		if dist > maxDist {
			maxDist = dist
			destination = dest
		}
	}

	return destination
}

// generate a random mystery trader reward based on the year of the game and speed of the MT
// early game is more likely to be research (or 1/6th chance of a ship)
func generateMysteryTraderReward(rules *Rules, year int, warpSpeed int) MysteryTraderRewardType {
	turn := year - rules.StartingYear

	var chance int
	if turn < 100 {
		// before year 100, we weight towards research rewards
		chance = 5
	} else if turn < 250 {
		// year 100 to 250, weight a little more towards random components
		chance = 3
	} else {
		// after year 250, it's a toss up if we go the research path vs random component
		chance = 2
	}

	// faster mystery traders have more of a chance for good rewards
	if warpSpeed <= 9 {
		chance++
	} else if warpSpeed >= 11 {
		chance--
	}

	// default to research
	reward := MysteryTraderRewardResearch

	if rules.random.Intn(10) < chance {
		// roll a 10 sided die, if it's less than our starting chance
		// 5 out of 6 traders will have research, the others will give a ship
		if rules.random.Intn(6) == 5 {
			reward = MysteryTraderRewardLifeboat
		} else {
			reward = MysteryTraderRewardResearch
		}
	} else {
		// we made it to the good stuff!
		// grab a random reward from the list, skipping the first one "None"
		randomReward := MysteryTraderRewards[rules.random.Intn(len(MysteryTraderRewards)-1)+1]
		if randomReward == MysteryTraderRewardTorpedo ||
			randomReward == MysteryTraderRewardBeamWeapon ||
			randomReward == MysteryTraderRewardGenesis ||
			randomReward == MysteryTraderRewardJumpGate {

			// the first random reward was something cool, make them roll again
			randomReward = MysteryTraderRewards[rules.random.Intn(len(MysteryTraderRewards)-1)+1]

			if turn < 120 && randomReward == MysteryTraderRewardBeamWeapon ||
				turn < 150 && randomReward == MysteryTraderRewardGenesis ||
				turn < 180 && randomReward == MysteryTraderRewardJumpGate {
				// we are too early for one of these good rewards. Half the time we just give research
				if rules.random.Intn(2) > 0 {
					randomReward = MysteryTraderRewardResearch
				}
			}
		}
		reward = randomReward
	}

	return reward
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
	} else {
		// move along the heading...
		mt.Heading = (mt.Destination.Subtract(mt.Position)).Normalized()
		mt.Position = mt.Position.Add(mt.Heading.Scale(dist))
		mt.Position = mt.Position.Round()
	}
}

func (mt *MysteryTrader) change(rules *Rules, game *Game) bool {
	if mt.WarpSpeed > rules.MysteryTraderRules.MaxWarp {
		return false
	}

	if rules.random.Intn(rules.MysteryTraderRules.ChanceCourseChange) > 0 {
		return false
	}

	// course change!
	mt.WarpSpeed++
	if rules.random.Intn(rules.MysteryTraderRules.ChanceSpeedUpOnly) == 0 {
		// speed up only
		return true
	}

	// point the MT at a new location
	mt.Destination = generateRandomMysteryTraderDestination(rules, game, mt.Position)

	return true
}

func (mt *MysteryTrader) again(rules *Rules, game *Game, numHumanPlayers int) bool {
	if len(mt.PlayersRewarded) == numHumanPlayers || rules.random.Intn(rules.MysteryTraderRules.ChanceAgain) > 0 {
		// we've given to all human players or we aren't going again
		return false
	}

	if mt.WarpSpeed-2 >= rules.MysteryTraderRules.MinWarp {
		// slow it down
		mt.WarpSpeed -= 2
	}

	// point the MT at a new location
	mt.Destination = generateRandomMysteryTraderDestination(rules, game, mt.Position)
	mt.Heading = (mt.Destination.Subtract(mt.Position)).Normalized()
	return true
}

func (mt *MysteryTrader) rewardedPlayer(num int) bool {
	return mt.PlayersRewarded[num]
}

// meet a mystery trader and recieve a reward!
func (mt *MysteryTrader) meet(rules *Rules, game *Game, fleet *Fleet, player *Player) MysteryTraderReward {
	gift := fleet.Cargo.ToMineral().Total()
	if fleet.Cargo.ToMineral().Total() >= mt.RequestedBoon {
		// it's a major award!
		rewardType := mt.RewardType

		if rewardType == MysteryTraderRewardResearch && player.TechLevels.Min() == rules.MaxTechLevel {
			if rules.random.Intn(rules.MysteryTraderRules.ChanceMaxTechGetsPart) > 0 {
				// player gets nothing
				return MysteryTraderReward{}
			}
			// player is maxed on tech, give them a part
			rewardType = MysteryTraderRewardParts[rules.random.Intn(len(MysteryTraderRewardParts))]
		}
		if rewardType.IsPart() {
			// the player gets a part, make sure they don't have them all
			tech, reward := getMysteryTraderPart(rules.random, player, rewardType)
			if tech != nil {
				// we found a tech for the player
				return MysteryTraderReward{
					Type: reward,
					Tech: tech.Name,
				}
			}
			// couldn't find a tech, so we get a different reward
			rewardType = reward
		}

		switch rewardType {
		case MysteryTraderRewardResearch:
			// give the player tech levels
			return mt.getTechLevelReward(rules, player, gift)
		case MysteryTraderRewardJumpGate:
			// player gets a jump gate
			return MysteryTraderReward{
				Type: MysteryTraderRewardJumpGate,
				Tech: JumpGate.Name,
			}
		case MysteryTraderRewardGenesis:
			// player gets a genesis device
			return MysteryTraderReward{
				Type: MysteryTraderRewardGenesis,
				Tech: GenesisDevice.Name,
			}
		case MysteryTraderRewardLifeboat:
			ship, count := getRandomLifeboat(rules.random, game.Year-rules.StartingYear)
			return MysteryTraderReward{
				Type:      MysteryTraderRewardLifeboat,
				Ship:      ship,
				ShipCount: count,
			}
		}

	}
	return MysteryTraderReward{}
}

// getTechLevelReward returns MysteryTraderReward awarding tech levels for a player
func (mt *MysteryTrader) getTechLevelReward(rules *Rules, player *Player, gift int) MysteryTraderReward {
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

// get a random lifeboat design the player will be given
func getRandomLifeboat(rng rng, turn int) (ShipDesign, int) {

	// figure out which ships the player doesn't already have
	ships := []ShipDesign{
		MysteryTraderScout,
		MysteryTraderProbe,
		MysteryTraderLifeboat,
	}

	count := 1
	if rng.Intn(3) == 0 {
		count = 2
	}

	if turn > 100 && rng.Intn(2) > 0 {
		// after turn 100, return probes and lifeboats
		return ships[1+rng.Intn(len(ships)-1)], count
	} else {
		return ships[rng.Intn(len(ships))], count
	}
}

// get a random mystery trader tech
// currently there is no "random" to it, because we only have one tech of each category
// but there *could be*
func getMysteryTraderPart(rng rng, player *Player, reward MysteryTraderRewardType) (*Tech, MysteryTraderRewardType) {

	// jump gate and genesis device are special, not categories
	// if the player doesn't already have them and the MT is offering it, give it
	if reward == MysteryTraderRewardJumpGate && !player.HasAcquiredTech(&JumpGate.Tech) {
		return &JumpGate.Tech, reward
	}

	if reward == MysteryTraderRewardGenesis && !player.HasAcquiredTech(&GenesisDevice.Tech) {
		return &GenesisDevice.Tech, reward
	}

	techsByCategory := make(map[TechCategory][]Tech, len(TechCategories))

	// for the initial check, don't include the jump gate and genesis device
	// in the techs the player can get. Those are more rare
	for _, tech := range MysteryTraderRandomTechs {
		techsByCategory[tech.Category] = append(techsByCategory[tech.Category], tech)
	}

	category := reward.Category()
	if category == TechCategoryNone || len(techsByCategory[category]) == 0 {
		// no techs for this category, they get ships
		return nil, MysteryTraderRewardLifeboat
	}

	tech := &techsByCategory[category][rng.Intn(len(techsByCategory[category]))]
	if player.HasAcquiredTech(tech) {
		// the player already has this tech, get a new random one from all techs (including jump gate and genesis device)
		for i := 0; i < 25; i++ {
			tech = &MysteryTraderTechs[rng.Intn(len(MysteryTraderTechs))]
			if !player.HasAcquiredTech(tech) {
				return tech, MysteryTraderRewardTypeForTech(tech)
			}
		}

		// couldn't find a tech the player didn't already have, they get ships
		return nil, MysteryTraderRewardLifeboat
	}

	// pick a random tech to aware for this category
	return tech, reward
}

// a list of all Mystery Trader Techs
var MysteryTraderTechs = []Tech{
	HushABoom.Tech,
	EnigmaPulsar.Tech,
	MegaPolyShell.Tech,
	LangstonShell.Tech,
	MultiFunctionPod.Tech,
	AntiMatterTorpedo.Tech,
	MultiContainedMunition.Tech,
	AlienMiner.Tech,
	MultiCargoPod.Tech,
	MiniMorph.Tech,
	JumpGate.Tech,
	GenesisDevice.Tech,
}

// a list of mystery trader techs that can be randomly
var MysteryTraderRandomTechs = []Tech{
	HushABoom.Tech,
	EnigmaPulsar.Tech,
	MegaPolyShell.Tech,
	LangstonShell.Tech,
	MultiFunctionPod.Tech,
	AntiMatterTorpedo.Tech,
	MultiContainedMunition.Tech,
	AlienMiner.Tech,
	MultiCargoPod.Tech,
	MiniMorph.Tech,
}

var HushABoom = TechHullComponent{Tech: NewTechWithOrigin("Hush-a-Boom", NewCost(1, 2, 0, 2), TechRequirements{Acquirable: true, TechLevel: TechLevel{Weapons: 12, Electronics: 12, Biotechnology: 12}}, 75, TechCategoryBomb, OriginMysteryTrader),
	Mass:                 5,
	KillRate:             3,
	StructureDestroyRate: 2,
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
	Scanner:        true,
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
	Scanner:        true,
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
	Scanner:              true,
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
var GenesisDevice = TechPlanetary{Tech: NewTechWithOrigin("Genesis Device", NewCost(0, 0, 0, 5000), TechRequirements{Acquirable: true, TechLevel: TechLevel{Energy: 20, Weapons: 10, Propulsion: 10, Construction: 20, Electronics: 10, Biotechnology: 20}}, 45, TechCategoryPlanetary, OriginMysteryTrader),
	ResetPlanet: true,
}

var MysteryTraderScout = ShipDesign{
	Hull:          MiniMorph.Name,
	Name:          "M.T. Scout",
	MysteryTrader: true,
	Slots: []ShipDesignSlot{
		{HullComponent: EnigmaPulsar.Name, HullSlotIndex: 1, Quantity: 2},
		{HullComponent: LangstonShell.Name, HullSlotIndex: 2, Quantity: 3},
		{HullComponent: MultiCargoPod.Name, HullSlotIndex: 3, Quantity: 1},
		{HullComponent: MultiFunctionPod.Name, HullSlotIndex: 4, Quantity: 1},
		{HullComponent: JumpGate.Name, HullSlotIndex: 5, Quantity: 1},
		{HullComponent: AntiMatterTorpedo.Name, HullSlotIndex: 6, Quantity: 2},
		{HullComponent: AntiMatterTorpedo.Name, HullSlotIndex: 7, Quantity: 2},
	},
}

var MysteryTraderProbe = ShipDesign{
	Hull:          MiniMorph.Name,
	Name:          "M.T. Probe",
	MysteryTrader: true,
	Slots: []ShipDesignSlot{
		{HullComponent: EnigmaPulsar.Name, HullSlotIndex: 1, Quantity: 2},
		{HullComponent: MegaPolyShell.Name, HullSlotIndex: 2, Quantity: 3},
		{HullComponent: MultiCargoPod.Name, HullSlotIndex: 3, Quantity: 1},
		{HullComponent: MultiFunctionPod.Name, HullSlotIndex: 4, Quantity: 1},
		{HullComponent: JumpGate.Name, HullSlotIndex: 5, Quantity: 1},
		{HullComponent: AntiMatterTorpedo.Name, HullSlotIndex: 6, Quantity: 2},
		{HullComponent: AntiMatterTorpedo.Name, HullSlotIndex: 7, Quantity: 2},
	},
}

var MysteryTraderLifeboat = ShipDesign{
	Hull:          Nubian.Name,
	Name:          "M.T. Lifeboat",
	MysteryTrader: true,
	Slots: []ShipDesignSlot{
		{HullComponent: EnigmaPulsar.Name, HullSlotIndex: 1, Quantity: 3},
		{HullComponent: MegaPolyShell.Name, HullSlotIndex: 2, Quantity: 3},
		{HullComponent: MegaPolyShell.Name, HullSlotIndex: 3, Quantity: 3},
		{HullComponent: LangstonShell.Name, HullSlotIndex: 4, Quantity: 3},
		{HullComponent: LangstonShell.Name, HullSlotIndex: 5, Quantity: 3},
		{HullComponent: MultiContainedMunition.Name, HullSlotIndex: 6, Quantity: 3},
		{HullComponent: MultiContainedMunition.Name, HullSlotIndex: 7, Quantity: 3},
		{HullComponent: MultiContainedMunition.Name, HullSlotIndex: 8, Quantity: 3},
		{HullComponent: MultiCargoPod.Name, HullSlotIndex: 9, Quantity: 3},
		{HullComponent: MultiFunctionPod.Name, HullSlotIndex: 10, Quantity: 3},
		{HullComponent: MultiFunctionPod.Name, HullSlotIndex: 11, Quantity: 3},
	},
}
