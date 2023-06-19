package cs

import (
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

type Universe struct {
	Planets              []*Planet                           `json:"planets,omitempty"`
	Fleets               []*Fleet                            `json:"fleets,omitempty"`
	Starbases            []*Fleet                            `json:"starbases,omitempty"`
	Wormholes            []*Wormhole                         `json:"wormholes,omitempty"`
	MineralPackets       []*MineralPacket                    `json:"mineralPackets,omitempty"`
	MineFields           []*MineField                        `json:"mineFields,omitempty"`
	MysteryTraders       []*MysteryTrader                    `json:"mysteryTraders,omitempty"`
	Salvages             []*Salvage                          `json:"salvage,omitempty"`
	rules                *Rules                              `json:"-"`
	battlePlansByNum     map[playerBattlePlanNum]*BattlePlan `json:"-"`
	mapObjectsByPosition map[Vector][]interface{}            `json:"-"`
	fleetsByPosition     map[Vector][]*Fleet                 `json:"-"`
	fleetsByNum          map[playerObject]*Fleet             `json:"-"`
	designsByNum         map[playerObject]*ShipDesign        `json:"-"`
	mineFieldsByNum      map[playerObject]*MineField         `json:"-"`
	mineralPacketsByNum  map[playerObject]*MineralPacket     `json:"-"`
	salvagesByNum        map[int]*Salvage                    `json:"-"`
	salvagesByPosition   map[Vector]*Salvage                 `json:"-"`
	mysteryTradersByNum  map[int]*MysteryTrader              `json:"-"`
	wormholesByNum       map[int]*Wormhole                   `json:"-"`
}

func NewUniverse(rules *Rules) Universe {
	return Universe{
		rules:                rules,
		battlePlansByNum:     make(map[playerBattlePlanNum]*BattlePlan),
		mapObjectsByPosition: make(map[Vector][]interface{}),
		fleetsByPosition:     make(map[Vector][]*Fleet),
		designsByNum:         make(map[playerObject]*ShipDesign),
		fleetsByNum:          make(map[playerObject]*Fleet),
		mineFieldsByNum:      make(map[playerObject]*MineField),
		mineralPacketsByNum:  make(map[playerObject]*MineralPacket),
		salvagesByNum:        make(map[int]*Salvage),
		salvagesByPosition:   make(map[Vector]*Salvage),
		wormholesByNum:       make(map[int]*Wormhole),
		mysteryTradersByNum:  make(map[int]*MysteryTrader),
	}
}

type mapObjectGetter interface {
	getShipDesign(playerNum int, num int) *ShipDesign
	getMapObject(mapObjectType MapObjectType, num int, playerNum int) *MapObject
	getPlanet(num int) *Planet
	getFleet(playerNum int, num int) *Fleet
	getMineField(playerNum int, num int) *MineField
	getMysteryTrader(num int) *MysteryTrader
	getWormhole(num int) *Wormhole
	getSalvage(num int) *Salvage
	getCargoHolder(mapObjectType MapObjectType, num int, playerNum int) cargoHolder
	getMapObjectsAtPosition(position Vector) []interface{}
	isPositionValid(pos Vector, occupiedLocations *[]Vector, minDistance float64) bool
	updateMapObjectAtPosition(mo interface{}, originalPosition, newPosition Vector)
}

type playerObject struct {
	PlayerNum int
	Num       int
}

func playerObjectKey(playerNum int, num int) playerObject { return playerObject{playerNum, num} }

type playerBattlePlanNum struct {
	PlayerNum int
	Num       int
}

// build the maps used for the Get functions
func (u *Universe) buildMaps(players []*Player) {

	// make a big map to hold all of our universe objects by position
	u.mapObjectsByPosition = make(map[Vector][]interface{}, len(u.Planets))

	// build a map of designs by num
	// so we can inject the design into each token
	numDesigns := 0
	numBattlePlans := 0
	for _, p := range players {
		numDesigns += len(p.Designs)
		numBattlePlans += len(p.BattlePlans)
	}
	u.designsByNum = make(map[playerObject]*ShipDesign, numDesigns)
	u.battlePlansByNum = make(map[playerBattlePlanNum]*BattlePlan, numBattlePlans)

	for _, p := range players {
		for i := range p.Designs {
			design := p.Designs[i]
			u.designsByNum[playerObjectKey(design.PlayerNum, design.Num)] = design
		}

		for i := range p.BattlePlans {
			plan := &p.BattlePlans[i]
			u.battlePlansByNum[playerBattlePlanNum{PlayerNum: p.Num, Num: plan.Num}] = plan
		}
	}

	u.fleetsByPosition = make(map[Vector][]*Fleet)
	u.fleetsByNum = make(map[playerObject]*Fleet, len(u.Fleets))
	for _, fleet := range u.Fleets {
		u.addMapObjectByPosition(fleet, fleet.Position)
		fleets, found := u.fleetsByPosition[fleet.Position]
		if !found {
			fleets = []*Fleet{fleet}
		}
		fleets = append(fleets, fleet)
		u.fleetsByPosition[fleet.Position] = fleets

		u.fleetsByNum[playerObjectKey(fleet.PlayerNum, fleet.Num)] = fleet

		fleet.battlePlan = u.battlePlansByNum[playerBattlePlanNum{fleet.PlayerNum, fleet.BattlePlanNum}]

		// inject the design into this
		for i := range fleet.Tokens {
			token := &fleet.Tokens[i]
			token.design = u.designsByNum[playerObjectKey(fleet.PlayerNum, token.DesignNum)]
		}
	}

	for _, starbase := range u.Starbases {
		u.Planets[starbase.PlanetNum-1].starbase = starbase
		for i := range starbase.Tokens {
			token := &starbase.Tokens[i]
			token.design = u.designsByNum[playerObjectKey(starbase.PlayerNum, token.DesignNum)]
		}
		starbase.battlePlan = u.battlePlansByNum[playerBattlePlanNum{starbase.PlayerNum, starbase.BattlePlanNum}]
	}

	for _, planet := range u.Planets {
		u.addMapObjectByPosition(planet, planet.Position)
	}

	for _, mineralPacket := range u.MineralPackets {
		u.mineralPacketsByNum[playerObjectKey(mineralPacket.PlayerNum, mineralPacket.Num)] = mineralPacket
		u.addMapObjectByPosition(mineralPacket, mineralPacket.Position)
	}
	for _, mineField := range u.MineFields {
		u.mineFieldsByNum[playerObjectKey(mineField.PlayerNum, mineField.Num)] = mineField
		u.addMapObjectByPosition(mineField, mineField.Position)
	}

	u.salvagesByPosition = make(map[Vector]*Salvage, len(u.Salvages))
	u.salvagesByNum = make(map[int]*Salvage, len(u.Salvages))
	for _, salvage := range u.Salvages {
		u.salvagesByNum[salvage.Num] = salvage
		u.addMapObjectByPosition(salvage, salvage.Position)
		u.salvagesByPosition[salvage.Position] = salvage
	}

	u.wormholesByNum = make(map[int]*Wormhole, len(u.Wormholes))
	for _, wormhole := range u.Wormholes {
		u.wormholesByNum[wormhole.Num] = wormhole
		u.addMapObjectByPosition(wormhole, wormhole.Position)
	}

	u.mysteryTradersByNum = make(map[int]*MysteryTrader, len(u.MysteryTraders))
	for _, mysteryTrader := range u.MysteryTraders {
		u.mysteryTradersByNum[mysteryTrader.Num] = mysteryTrader
		u.addMapObjectByPosition(mysteryTrader, mysteryTrader.Position)
	}

}

func (u *Universe) addMapObjectByPosition(mo interface{}, position Vector) {
	mos, found := u.mapObjectsByPosition[position]
	if !found {
		mos = []interface{}{}
		u.mapObjectsByPosition[position] = mos
	}
	mos = append(mos, mo)
	u.mapObjectsByPosition[position] = mos
}

// Check if a position vector is a mininum distance away from all other points
func (u *Universe) isPositionValid(pos Vector, occupiedLocations *[]Vector, minDistance float64) bool {
	minDistanceSquared := minDistance * minDistance

	for _, to := range *occupiedLocations {
		if pos.DistanceSquaredTo(to) <= minDistanceSquared {
			return false
		}
	}
	return true
}

// get all commandable map objects for a player
func (u *Universe) GetPlayerMapObjects(playerNum int) PlayerMapObjects {
	pmo := PlayerMapObjects{}

	pmo.Fleets = u.getFleets(playerNum)
	pmo.Planets = u.getPlanets(playerNum)
	pmo.MineFields = u.getMineFields(playerNum)

	return pmo
}

func (u *Universe) getMapObject(mapObjectType MapObjectType, num int, playerNum int) *MapObject {
	switch mapObjectType {
	case MapObjectTypePlanet:
		planet := u.getPlanet(num)
		if planet != nil {
			return &planet.MapObject
		}
	case MapObjectTypeFleet:
		fleet := u.getFleet(playerNum, num)
		if fleet != nil {
			return &fleet.MapObject
		}
	case MapObjectTypeWormhole:
		wormhole := u.getWormhole(num)
		if wormhole != nil {
			return &wormhole.MapObject
		}
	case MapObjectTypeMineField:
		mineField := u.getMineField(playerNum, num)
		if mineField != nil {
			return &mineField.MapObject
		}
	case MapObjectTypeMysteryTrader:
		mysteryTrader := u.getMysteryTrader(num)
		if mysteryTrader != nil {
			return &mysteryTrader.MapObject
		}
	case MapObjectTypeSalvage:
		salvage := u.getSalvage(num)
		if salvage != nil {
			return &salvage.MapObject
		}
	case MapObjectTypeMineralPacket:
		mineralPacket := u.getMineralPacket(playerNum, num)
		if mineralPacket != nil {
			return &mineralPacket.MapObject
		}
	}
	return nil
}

// get a ship design by num
func (u *Universe) getShipDesign(playerNum int, num int) *ShipDesign {
	return u.designsByNum[playerObjectKey(playerNum, num)]
}

// Get a planet by num
func (u *Universe) getPlanet(num int) *Planet {
	return u.Planets[num-1]
}

// Get a fleet by player num and fleet num
func (u *Universe) getFleet(playerNum int, num int) *Fleet {
	return u.fleetsByNum[playerObjectKey(playerNum, num)]
}

// Get a planet by num
func (u *Universe) getWormhole(num int) *Wormhole {
	return u.wormholesByNum[num]
}

// Get a salvage by num
func (u *Universe) getSalvage(num int) *Salvage {
	return u.salvagesByNum[num]
}

func (u *Universe) getMineField(playerNum int, num int) *MineField {
	return u.mineFieldsByNum[playerObjectKey(playerNum, num)]
}

// get a minefield that is close to a position
func (u *Universe) getMineFieldNearPosition(playerNum int, position Vector, mineFieldType MineFieldType) *MineField {
	for _, mineField := range u.MineFields {
		if mineField.PlayerNum == playerNum && mineField.MineFieldType == mineFieldType && isPointInCircle(position, mineField.Position, mineField.Spec.Radius) {
			return mineField
		}
	}

	return nil
}

func (u *Universe) getMineralPacket(playerNum int, num int) *MineralPacket {
	return u.mineralPacketsByNum[playerObjectKey(playerNum, num)]
}

func (u *Universe) getMysteryTrader(num int) *MysteryTrader {
	return u.mysteryTradersByNum[num]
}

// get a cargo holder by natural key (num, playerNum, etc)
func (u *Universe) getCargoHolder(mapObjectType MapObjectType, num int, playerNum int) cargoHolder {
	switch mapObjectType {
	case MapObjectTypePlanet:
		return u.getPlanet(num)
	case MapObjectTypeFleet:
		return u.getFleet(playerNum, num)
	case MapObjectTypeSalvage:
		return u.getSalvage(num)
	}
	return nil
}

// mark a fleet as deleted and remove it from the universe
func (u *Universe) deleteFleet(fleet *Fleet) {
	fleet.Delete = true

	// decrease token count
	for _, token := range fleet.Tokens {
		token.design.Spec.NumInstances -= token.Quantity
		token.design.MarkDirty()
	}

	index := slices.Index(u.Fleets, fleet)
	slices.Delete(u.Fleets, index, index)

	delete(u.fleetsByNum, playerObjectKey(fleet.PlayerNum, fleet.Num))
	delete(u.fleetsByPosition, fleet.Position)
	u.removeMapObjectAtPosition(fleet, fleet.Position)
}

// move a fleet from one position to another
func (u *Universe) moveFleet(fleet *Fleet, originalPosition Vector) {
	fleet.MarkDirty()
	originalPositionFleets := u.fleetsByPosition[originalPosition]
	index := slices.Index(originalPositionFleets, fleet)
	if index != -1 {
		slices.Delete(originalPositionFleets, index, index)
		if len(originalPositionFleets) == 0 {
			delete(u.fleetsByPosition, originalPosition)
		}
	}

	// move to new location
	fleets, found := u.fleetsByPosition[fleet.Position]
	if !found {
		fleets = []*Fleet{}
	}
	fleets = append(fleets, fleet)
	u.fleetsByPosition[fleet.Position] = fleets

	// upadte mapobjects position
	u.updateMapObjectAtPosition(fleet, originalPosition, fleet.Position)
}

// move a wormhole from one position to another
func (u *Universe) moveWormhole(wormhole *Wormhole, originalPosition Vector) {
	wormhole.MarkDirty()
	u.updateMapObjectAtPosition(wormhole, originalPosition, wormhole.Position)
}

// delete a wormhole from the universe
func (u *Universe) deleteWormhole(wormhole *Wormhole) {
	wormhole.Delete = true

	index := slices.Index(u.Wormholes, wormhole)
	slices.Delete(u.Wormholes, index, index)

	delete(u.wormholesByNum, wormhole.Num)
	u.removeMapObjectAtPosition(wormhole, wormhole.Position)
}

// create a new wormhole in the universe
func (u *Universe) createWormhole(position Vector, stability WormholeStability, companion *Wormhole) *Wormhole {
	num := 1
	if len(u.Wormholes) > 0 {
		num = u.Wormholes[len(u.Wormholes)-1].Num + 1
	}

	wormhole := newWormhole(position, num, stability)

	if companion != nil {
		companion.DestinationNum = wormhole.Num
		wormhole.DestinationNum = companion.Num
	}

	// compute the spec for this wormhole
	wormhole.Spec = computeWormholeSpec(wormhole, u.rules)

	u.Wormholes = append(u.Wormholes, wormhole)
	u.wormholesByNum[wormhole.Num] = wormhole
	u.addMapObjectByPosition(wormhole, wormhole.Position)

	return wormhole
}

// delete a wormhole from the universe
func (u *Universe) createSalvage(position Vector, playerNum int, cargo Cargo) *Salvage {
	num := 1
	if len(u.Salvages) > 0 {
		num = u.Salvages[len(u.Salvages)-1].Num + 1
	}
	salvage := newSalvage(position, num, playerNum, cargo)
	u.Salvages = append(u.Salvages, salvage)
	u.salvagesByNum[num] = salvage
	u.addMapObjectByPosition(salvage, salvage.Position)

	return salvage
}

// delete a salvage from the universe
func (u *Universe) deleteSalvage(salvage *Salvage) {
	salvage.Delete = true

	index := slices.Index(u.Salvages, salvage)
	slices.Delete(u.Salvages, index, index)

	delete(u.salvagesByNum, salvage.Num)
	u.removeMapObjectAtPosition(salvage, salvage.Position)
}

func (u *Universe) getPlanets(playerNum int) []*Planet {
	planets := []*Planet{}
	for _, planet := range u.Planets {
		if planet.PlayerNum == playerNum {
			planets = append(planets, planet)
		}
	}
	return planets
}

func (u *Universe) getFleets(playerNum int) []*Fleet {
	fleets := []*Fleet{}
	for _, fleet := range u.Fleets {
		if fleet.PlayerNum == playerNum {
			fleets = append(fleets, fleet)
		}
	}
	return fleets
}

func (u *Universe) getMineFields(playerNum int) []*MineField {
	mineFields := []*MineField{}
	for _, mineField := range u.MineFields {
		if mineField.PlayerNum == playerNum {
			mineFields = append(mineFields, mineField)
		}
	}
	return mineFields
}

// get a slice of mapobjects at a position, or nil if none
func (u *Universe) getMapObjectsAtPosition(position Vector) []interface{} {
	return u.mapObjectsByPosition[position]
}

// get a slice of mapobjects at a position, or nil if none
func (u *Universe) updateMapObjectAtPosition(mo interface{}, originalPosition, newPosition Vector) {
	mos := u.mapObjectsByPosition[originalPosition]
	if mos != nil {
		index := slices.IndexFunc(mos, func(item interface{}) bool { return item == mo })
		if index >= 0 && index < len(mos) {
			slices.Delete(mos, index, index)
		} else {
			log.Warn().Msgf("tried to update position of %s from %v to %v, but index %d of original position out of range", mo, originalPosition, newPosition, index)
		}
	} else {
		log.Warn().Msgf("tried to update position of %s from %v to %v, no mapobjects were found at %v", mo, originalPosition, newPosition, originalPosition)
	}

	// add the new object to the list
	u.addMapObjectByPosition(mo, newPosition)
}

// get a slice of mapobjects at a position, or nil if none
func (u *Universe) removeMapObjectAtPosition(mo interface{}, position Vector) {
	mos := u.mapObjectsByPosition[position]
	if mos != nil {
		index := slices.IndexFunc(mos, func(item interface{}) bool { return item == mo })
		if index >= 0 && index < len(mos) {
			slices.Delete(mos, index, index)
		} else {
			log.Warn().Msgf("tried to remove mapobject %s at position %v but index %d of position out of range", mo, position, index)
		}
	} else {
		log.Warn().Msgf("tried to to remove mapobject %s at position %v, no mapobjects were found at %v", mo, position, position)
	}
}

// get the next design number to use
func (u *Universe) getNextMineFieldNum() int {
	num := 0
	for _, mineField := range u.MineFields {
		num = maxInt(num, mineField.Num)
	}
	return num + 1
}

func (u *Universe) addMineField(mineField *MineField) {
	u.MineFields = append(u.MineFields, mineField)
	u.mineFieldsByNum[playerObjectKey(mineField.PlayerNum, mineField.Num)] = mineField
	u.addMapObjectByPosition(mineField, mineField.Position)
}
