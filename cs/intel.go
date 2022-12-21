package cs

import (
	"fmt"

	"github.com/google/uuid"
)

const ReportAgeUnexplored = -1
const Unowned = 0

type discover struct {
	player              *Player
	fleetIntelsByKey    map[string]*FleetIntel
	wormholeIntelsByKey map[int]*WormholeIntel
	designIntelsByKey   map[uuid.UUID]*ShipDesignIntel
}

type discoverer interface {
	playerInfoDiscover(player *Player)
	clearTransientReports()
	discoverPlanet(rules *Rules, player *Player, planet *Planet, penScanned bool) error
	discoverPlanetCargo(player *Player, planet *Planet) error
	discoverFleet(player *Player, fleet *Fleet)
	discoverFleetCargo(player *Player, fleet *Fleet)
	discoverWormhole(player *Player, wormhole *Wormhole)
	discoverDesign(player *Player, design *ShipDesign, discoverSlots bool)

	getWormholeIntel(num int) *WormholeIntel
}

func newDiscoverer(player *Player) discoverer {
	d := &discover{
		player: player,
	}
	d.fleetIntelsByKey = make(map[string]*FleetIntel, len(player.FleetIntels))
	for i := range player.FleetIntels {
		intel := &player.FleetIntels[i]
		d.fleetIntelsByKey[intel.String()] = intel
	}

	d.wormholeIntelsByKey = make(map[int]*WormholeIntel, len(player.WormholeIntels))
	for i := range player.WormholeIntels {
		intel := &player.WormholeIntels[i]
		d.wormholeIntelsByKey[intel.Num] = intel
	}

	d.designIntelsByKey = make(map[uuid.UUID]*ShipDesignIntel, len(player.ShipDesignIntels))
	for i := range player.ShipDesignIntels {
		intel := &player.ShipDesignIntels[i]
		d.designIntelsByKey[intel.UUID] = intel
	}

	return d
}

type Intel struct {
	Name      string `json:"name"`
	Num       int    `json:"num"`
	PlayerNum int    `json:"playerNum"`
	ReportAge int    `json:"reportAge"`
}

type MapObjectIntel struct {
	Intel
	Type     MapObjectType `json:"type"`
	Position Vector        `json:"position"`
}

func (intel *Intel) String() string {
	return fmt.Sprintf("Num: %3d %s", intel.Num, intel.Name)
}

func (intel *Intel) Owned() bool {
	return intel.PlayerNum != Unowned
}

type PlanetIntel struct {
	MapObjectIntel
	Hab                  Hab         `json:"hab,omitempty"`
	MineralConcentration Mineral     `json:"mineralConcentration,omitempty"`
	Population           uint        `json:"population,omitempty"`
	Starbase             *FleetIntel `json:"starbase,omitempty"`
	Cargo                Cargo       `json:"cargo,omitempty"`
	CargoDiscovered      bool        `json:"cargoDiscovered,omitempty"`
}

type ShipDesignIntel struct {
	Intel
	UUID          uuid.UUID        `json:"uuid,omitempty"`
	Name          string           `json:"name,omitempty"`
	Hull          string           `json:"hull,omitempty"`
	HullSetNumber int              `json:"hullSetNumber,omitempty"`
	Version       int              `json:"version,omitempty"`
	Armor         int              `json:"armor,omitempty"`
	Shields       int              `json:"shields,omitempty"`
	Slots         []ShipDesignSlot `json:"slots,omitempty"`
}

type FleetIntel struct {
	MapObjectIntel
	PlanetIntelID   int64 `json:"-"` // for starbase fleets that are owned by a planet
	Cargo           Cargo `json:"cargo,omitempty"`
	CargoDiscovered bool  `json:"cargoDiscovered,omitempty"`
}

type MineralPacketIntel struct {
	MapObjectIntel
	WarpFactor uint   `json:"warpFactor,omitempty"`
	Heading    Vector `json:"position"`
	Cargo      Cargo  `json:"cargo,omitempty"`
}

type SalvageIntel struct {
	MapObjectIntel
	Cargo Cargo `json:"cargo,omitempty"`
}

type MineFieldIntel struct {
	MapObjectIntel
	NumMines uint          `json:"numMines,omitempty"`
	Type     MineFieldType `json:"type,omitempty"`
}

type WormholeIntel struct {
	MapObjectIntel
	DestinationNum int               `json:"destinationNum,omitempty"`
	Stability      WormholeStability `json:"stability,omitempty"`
}

type PlayerIntel struct {
	Name           string `json:"name,omitempty"`
	Num            int    `json:"num,omitempty"`
	Color          string `json:"color,omitempty"`
	Seen           bool   `json:"seen,omitempty"`
	RaceName       string `json:"raceName,omitempty"`
	RacePluralName string `json:"racePluralName,omitempty"`
}

func (p *PlanetIntel) String() string {
	return fmt.Sprintf("Planet %s", &p.MapObjectIntel)
}

func (f *FleetIntel) String() string {
	return fmt.Sprintf("Player: %d, Fleet: %s", f.PlayerNum, f.Name)
}

func (d *ShipDesignIntel) String() string {
	return fmt.Sprintf("Player: %d, Fleet: %s", d.PlayerNum, d.Name)
}

// create a new FleetIntel object by key
func newFleetIntel(playerNum int, name string) FleetIntel {
	return FleetIntel{
		MapObjectIntel: MapObjectIntel{
			Type: MapObjectTypeFleet,
			Intel: Intel{
				Name:      name,
				PlayerNum: playerNum,
			},
		},
	}
}

// create a new WormholeIntel object by key
func newWormholeIntel(num int) *WormholeIntel {
	return &WormholeIntel{
		MapObjectIntel: MapObjectIntel{
			Type: MapObjectTypeWormhole,
			Intel: Intel{
				Num: num,
			},
		},
	}
}

// true if we haven't explored this planet
func (intel *PlanetIntel) Unexplored() bool {
	return intel.ReportAge == ReportAgeUnexplored
}

// true if we have explored this planet
func (intel *PlanetIntel) Explored() bool {
	return intel.ReportAge != ReportAgeUnexplored
}

// clear any transient player reports that are refreshed each turn
func (d *discover) clearTransientReports() {
	d.player.FleetIntels = []FleetIntel{}
}

// discover a planet and add it to the player's intel
func (d *discover) discoverPlanet(rules *Rules, player *Player, planet *Planet, penScanned bool) error {

	var intel *PlanetIntel
	planetIndex := planet.Num - 1

	if planetIndex < 0 || planetIndex >= len(player.PlanetIntels) {
		return fmt.Errorf("player %s cannot discover planet %s, planetIndex %d out of range", player, planet, planetIndex)
	}

	intel = &player.PlanetIntels[planetIndex]

	// everyone knows these about planets
	intel.Position = planet.Position
	intel.Name = planet.Name
	intel.Num = planet.Num

	ownedByPlayer := planet.PlayerNum != Unowned && player.Num == planet.PlayerNum

	if penScanned || ownedByPlayer {
		intel.PlayerNum = planet.PlayerNum

		// if we pen scanned the planet, we learn some things
		intel.Hab = planet.Hab
		intel.MineralConcentration = planet.MineralConcentration
		intel.ReportAge = 0

		// players know their planet pops, but other planets are slightly off
		if ownedByPlayer {
			intel.Population = uint(planet.population())
		} else {
			var randomPopulationError = rules.random.Float64()*(rules.PopulationScannerError-(-rules.PopulationScannerError)) - rules.PopulationScannerError
			intel.Population = uint(float64(planet.population()) * (1 - randomPopulationError))
		}
	}
	return nil
}

// discover the cargo of a planet
func (d *discover) discoverPlanetCargo(player *Player, planet *Planet) error {

	var intel *PlanetIntel
	planetIndex := planet.Num - 1

	if planetIndex < 0 || planetIndex >= len(player.PlanetIntels) {
		return fmt.Errorf("player %s cannot discover planet %s, planetIndex %d out of range", player, planet, planetIndex)
	}

	intel = &player.PlanetIntels[planetIndex]

	intel.CargoDiscovered = true
	intel.Cargo = Cargo{
		Ironium:   planet.Cargo.Ironium,
		Boranium:  planet.Cargo.Boranium,
		Germanium: planet.Cargo.Germanium,
	}

	return nil

}

// discover a fleet and add it to the player's fleet intel
func (d *discover) discoverFleet(player *Player, fleet *Fleet) {
	intel := newFleetIntel(fleet.PlayerNum, fleet.Name)

	intel.Name = fleet.Name
	intel.PlayerNum = fleet.PlayerNum
	intel.Position = fleet.Position

	player.FleetIntels = append(player.FleetIntels, intel)
	d.fleetIntelsByKey[intel.String()] = &intel
}

// discover cargo for an existing fleet
func (d *discover) discoverFleetCargo(player *Player, fleet *Fleet) {
	key := newFleetIntel(fleet.PlayerNum, fleet.Name)

	existingIntel, found := d.fleetIntelsByKey[key.String()]
	if found {
		existingIntel.Cargo = fleet.Cargo
		existingIntel.CargoDiscovered = true
	}

}

// discover a wormhole and add it to the player's wormhole intel
func (d *discover) discoverWormhole(player *Player, wormhole *Wormhole) {
	intel, found := d.wormholeIntelsByKey[wormhole.Num]
	if !found {
		// discover this new wormhole
		player.WormholeIntels = append(player.WormholeIntels, *newWormholeIntel(wormhole.Num))
		intel = &player.WormholeIntels[len(player.WormholeIntels)-1]
		d.wormholeIntelsByKey[intel.Num] = intel
	}

	intel.Name = wormhole.Name
	intel.Position = wormhole.Position
	intel.Stability = wormhole.Stability
}

// discover a player's design. This is a noop if we already know about
// the design and aren't discovering slots
func (d *discover) discoverDesign(player *Player, design *ShipDesign, discoverSlots bool) {
	intel, found := d.designIntelsByKey[design.UUID]
	if !found {
		// create a new intel for this design
		intel = &ShipDesignIntel{
			Intel: Intel{
				Name:      design.Name,
				PlayerNum: design.PlayerNum,
			},
			UUID:          design.UUID,
			Hull:          design.Hull,
			HullSetNumber: design.HullSetNumber,
		}

		// save this new design to our intel
		player.ShipDesignIntels = append(player.ShipDesignIntels, *intel)
		intel = &player.ShipDesignIntels[len(player.ShipDesignIntels)-1]
		d.designIntelsByKey[intel.UUID] = intel
	}

	// discover slots if we haven't already and this scanner discovers them
	if discoverSlots && len(intel.Slots) == 0 {
		intel.Slots = make([]ShipDesignSlot, len(design.Slots))
		copy(intel.Slots, design.Slots)
		intel.Armor = design.Spec.Armor
		intel.Shields = design.Spec.Shield
	}
}

func (d *discover) playerInfoDiscover(player *Player) {
	// d.game <- players to discover
	// discover info about other players
}

func (d *discover) getWormholeIntel(num int) *WormholeIntel {
	return d.wormholeIntelsByKey[num]
}
