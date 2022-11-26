package game

import (
	"sort"

	"github.com/google/uuid"
)

type Universe struct {
	Planets           []*Planet                            `json:"planets,omitempty" gorm:"foreignKey:GameID;references:ID"`
	Fleets            []*Fleet                             `json:"fleets,omitempty" gorm:"foreignKey:GameID;references:ID"`
	Starbases         []*Fleet                             `json:"starbases,omitempty" gorm:"foreignKey:GameID;references:ID"`
	Wormholes         []*Wormohole                         `json:"wormholes,omitempty" gorm:"foreignKey:GameID;references:ID"`
	MineralPackets    []*MineralPacket                     `json:"mineralPackets,omitempty" gorm:"foreignKey:GameID;references:ID"`
	MineFields        []*MineField                         `json:"mineFields,omitempty" gorm:"foreignKey:GameID;references:ID"`
	Salvage           []*Salvage                           `json:"salvage,omitempty" gorm:"foreignKey:GameID;references:ID"`
	fleetsByPosition  map[Vector]*Fleet                    `json:"-"`
	fleetsByNum       map[playerFleetNum]*Fleet            `json:"-"`
	designsByUUID     map[uuid.UUID]*ShipDesign            `json:"-"`
	battlePlansByName map[playerBattlePlanName]*BattlePlan `json:"-"`
}

type MapObjectGetter interface {
	GetShipDesign(uuid uuid.UUID) *ShipDesign
	GetPlanet(num int) *Planet
	GetFleet(playerNum int, num int) *Fleet
	GetWormhole(num int) *Wormohole
	GetSalvage(num int) *Salvage
	GetCargoHolder(mapObjectType MapObjectType, num int, playerNum int) CargoHolder
}

// gets lists of player commandable mapobjects
type PlayerMapObjectGetter interface {
	GetPlanets(playerNum int) []*Planet
	GetFleets(playerNum int) []*Fleet
	GetMineFields(playerNum int) []*MineField
}

type playerFleetNum struct {
	PlayerNum int
	Num       int
}

type playerBattlePlanName struct {
	PlayerNum int
	Name      string
}

// build the maps used for the Get functions
func (u *Universe) buildMaps(players []*Player) {

	// build a map of designs by uuid
	// so we can inject the design into each token
	numDesigns := 0
	numBattlePlans := 0
	for _, p := range players {
		numDesigns += len(p.Designs)
		numBattlePlans += len(p.BattlePlans)
	}
	u.designsByUUID = make(map[uuid.UUID]*ShipDesign, numDesigns)
	u.battlePlansByName = make(map[playerBattlePlanName]*BattlePlan, numBattlePlans)

	for _, p := range players {
		for i := range p.Designs {
			design := &p.Designs[i]
			u.designsByUUID[design.UUID] = design
		}

		for i := range p.BattlePlans {
			plan := &p.BattlePlans[i]
			u.battlePlansByName[playerBattlePlanName{PlayerNum: p.Num, Name: plan.Name}] = plan
		}
	}

	u.fleetsByPosition = make(map[Vector]*Fleet, len(u.Fleets))
	u.fleetsByNum = make(map[playerFleetNum]*Fleet, len(u.Fleets))
	for _, fleet := range u.Fleets {
		u.fleetsByPosition[fleet.Position] = fleet
		u.fleetsByNum[playerFleetNum{fleet.PlayerNum, fleet.Num}] = fleet

		fleet.battlePlan = u.battlePlansByName[playerBattlePlanName{fleet.PlayerNum, fleet.BattlePlanName}]
		// inject the design into this
		for i := range fleet.Tokens {
			token := &fleet.Tokens[i]
			token.design = u.designsByUUID[token.DesignUUID]
		}
	}

	for _, starbase := range u.Starbases {
		u.Planets[starbase.PlanetNum-1].starbase = starbase
	}

}

// get all commandable map objects for a player
func (u *Universe) GetPlayerMapObjects(playerNum int) PlayerMapObjects {
	pmo := PlayerMapObjects{}

	pmo.Fleets = u.GetFleets(playerNum)
	pmo.Planets = u.GetPlanets(playerNum)
	pmo.MineFields = u.GetMineFields(playerNum)

	return pmo
}

// get a ship design by uuid
func (u *Universe) GetShipDesign(uuid uuid.UUID) *ShipDesign {
	return u.designsByUUID[uuid]
}

// Get a planet by num
func (u *Universe) GetPlanet(num int) *Planet {
	return u.Planets[num-1]
}

// Get a fleet by player num and fleet num
func (u *Universe) GetFleet(playerNum int, num int) *Fleet {
	return u.fleetsByNum[playerFleetNum{playerNum, num}]
}

// Get a planet by num
func (u *Universe) GetWormhole(num int) *Wormohole {
	return u.Wormholes[num]
}

// Get a salvage by num
func (u *Universe) GetSalvage(num int) *Salvage {
	return u.Salvage[num]
}

// Get a mineralpacket by num
func (u *Universe) GetMineralPacket(num int) *MineralPacket {
	return u.MineralPackets[num]
}

// get a cargo holder by natural key (num, playerNum, etc)
func (u *Universe) GetCargoHolder(mapObjectType MapObjectType, num int, playerNum int) CargoHolder {
	switch mapObjectType {
	case MapObjectTypePlanet:
		return u.GetPlanet(num)
	case MapObjectTypeFleet:
		return u.GetFleet(playerNum, num)
	}
	return nil
}

// mark a fleet as deleted and remove it from the universe
func (u *Universe) DeleteFleet(fleet *Fleet) {
	fleet.Dirty = true
	fleet.Delete = true
	delete(u.fleetsByNum, playerFleetNum{fleet.PlayerNum, fleet.Num})
	delete(u.fleetsByPosition, fleet.Position)
}

// move a fleet from one position to another
func (u *Universe) MoveFleet(fleet *Fleet, originalPosition Vector) {
	fleet.Dirty = true
	delete(u.fleetsByPosition, originalPosition)
	u.fleetsByPosition[originalPosition] = fleet
}

func (u *Universe) GetPlanets(playerNum int) []*Planet {
	planets := []*Planet{}
	for _, planet := range u.Planets {
		if planet.PlayerNum == playerNum {
			planets = append(planets, planet)
		}
	}
	return planets
}

func (u *Universe) GetFleets(playerNum int) []*Fleet {
	fleets := []*Fleet{}
	for _, fleet := range u.Fleets {
		if fleet.PlayerNum == playerNum {
			fleets = append(fleets, fleet)
		}
	}
	return fleets
}

func (u *Universe) GetMineFields(playerNum int) []*MineField {
	mineFields := []*MineField{}
	for _, mineField := range u.MineFields {
		if mineField.PlayerNum == playerNum {
			mineFields = append(mineFields, mineField)
		}
	}
	return mineFields
}

func (u *Universe) getNextFleetNum(playerNum int) int {
	num := 1

	playerFleets := u.GetFleets(playerNum)
	orderedFleets := make([]*Fleet, len(playerFleets))
	copy(orderedFleets, playerFleets)
	sort.Slice(orderedFleets, func(i, j int) bool { return orderedFleets[i].Num < orderedFleets[j].Num })

	for i := 0; i < len(orderedFleets); i++ {
		// todo figure out starbasees
		fleet := orderedFleets[i]
		if i > 0 {
			// if we are past fleet #1 and we skipped a number, used the skipped number
			if fleet.Num > 1 && fleet.Num != orderedFleets[i-1].Num+1 {
				return orderedFleets[i-1].Num + 1
			}
		}
		// we are the next num...
		num = fleet.Num + 1
	}

	return num
}
