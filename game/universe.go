package game

type Universe struct {
	Area             Vector
	Planets          []Planet                  `json:"planets,omitempty"`
	Fleets           []Fleet                   `json:"fleets,omitempty"`
	Wormholes        []Wormohole               `json:"wormholes,omitempty"`
	MineralPackets   []MineralPacket           `json:"mineralPackets,omitempty"`
	Salvage          []Salvage                 `json:"salvage,omitempty"`
	FleetsByPosition map[Vector]*Fleet         `json:"-"`
	FleetsByNum      map[playerFleetNum]*Fleet `json:"-"`
}

type MapObjectGetter interface {
	GetPlanet(num int) *Planet
	GetFleet(playerNum int, num int) *Fleet
	GetWormhole(num int) *Wormohole
	GetSalvage(num int) *Salvage
	GetCargoHolder(mapObjectType MapObjectType, num int, playerNum int) CargoHolder
}

type playerFleetNum struct {
	PlayerNum int
	Num       int
}

// build the maps used for the Get functions
func (u *Universe) buildMaps() {
	u.FleetsByPosition = make(map[Vector]*Fleet, len(u.Fleets))
	u.FleetsByNum	 = make(map[playerFleetNum]*Fleet, len(u.Fleets))
	for i := range u.Fleets {
		fleet := &u.Fleets[i]
		u.FleetsByPosition[fleet.Position] = fleet
		u.FleetsByNum[playerFleetNum{fleet.PlayerNum, fleet.Num}] = fleet
	}
}

// Get a planet by num
func (u *Universe) GetPlanet(num int) *Planet {
	return &u.Planets[num]
}

// Get a fleet by player num and fleet num
func (u *Universe) GetFleet(playerNum int, num int) *Fleet {
	return u.FleetsByNum[playerFleetNum{playerNum, num}]
}

// Get a planet by num
func (u *Universe) GetWormhole(num int) *Wormohole {
	return &u.Wormholes[num]
}

// Get a salvage by num
func (u *Universe) GetSalvage(num int) *Salvage {
	return &u.Salvage[num]
}

// Get a mineralpacket by num
func (u *Universe) GetMineralPacket(num int) *MineralPacket {
	return &u.MineralPackets[num]
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
	delete(u.FleetsByNum, playerFleetNum{fleet.PlayerNum, fleet.Num})
	delete(u.FleetsByPosition, fleet.Position)
}
