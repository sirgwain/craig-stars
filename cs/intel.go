package cs

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

const ReportAgeUnexplored = -1
const Unowned = 0

type discover struct {
	player                   *Player
	fleetIntelsByKey         map[playerObject]*FleetIntel
	salvageIntelsByKey       map[int]*SalvageIntel
	mineFieldIntelsByKey     map[playerObject]*MineFieldIntel
	mineralPacketIntelsByKey map[playerObject]*MineralPacketIntel
	designIntelsByKey        map[playerObject]*ShipDesignIntel
	wormholeIntelsByKey      map[int]*WormholeIntel
	mysteryTraderIntelsByKey map[int]*MysteryTraderIntel
}

type discoverer interface {
	clearTransientReports()
	discoverPlayer(player *Player)
	discoverPlayerScores(player *Player)
	discoverPlanet(rules *Rules, player *Player, planet *Planet, penScanned bool) error
	discoverPlanetStarbase(player *Player, planet *Planet) error
	discoverPlanetCargo(player *Player, planet *Planet) error
	discoverPlanetTerraformability(player *Player, planetNum int) error
	discoverFleet(player *Player, fleet *Fleet)
	discoverFleetCargo(player *Player, fleet *Fleet)
	discoverMineField(player *Player, mineField *MineField)
	discoverMineralPacket(rules *Rules, player *Player, mineralPacket *MineralPacket, packetPlayer *Player, target *Planet)
	discoverSalvage(salvage *Salvage)
	discoverWormhole(wormhole *Wormhole)
	discoverWormholeLink(wormhole1, wormhole2 *Wormhole)
	discoverMysteryTrader(mineField *MysteryTrader)
	discoverDesign(player *Player, design *ShipDesign, discoverSlots bool)

	getPlanetIntel(num int) *PlanetIntel
	getWormholeIntel(num int) *WormholeIntel
	getMysteryTraderIntel(num int) *MysteryTraderIntel
	getMineFieldIntel(playerNum, num int) *MineFieldIntel
	getMineralPacketIntel(playerNum, num int) *MineralPacketIntel
	getFleetIntel(playerNum, num int) *FleetIntel
	getSalvageIntel(num int) *SalvageIntel
}

func newDiscoverer(player *Player) discoverer {
	d := &discover{
		player: player,
	}
	d.fleetIntelsByKey = make(map[playerObject]*FleetIntel, len(player.FleetIntels))
	for i := range player.FleetIntels {
		intel := &player.FleetIntels[i]
		d.fleetIntelsByKey[intel.playerObjectKey()] = intel
	}

	d.mineFieldIntelsByKey = make(map[playerObject]*MineFieldIntel, len(player.MineFieldIntels))
	for i := range player.MineFieldIntels {
		intel := &player.MineFieldIntels[i]
		d.mineFieldIntelsByKey[intel.playerObjectKey()] = intel
	}

	d.mineralPacketIntelsByKey = make(map[playerObject]*MineralPacketIntel, len(player.MineralPacketIntels))
	for i := range player.MineralPacketIntels {
		intel := &player.MineralPacketIntels[i]
		d.mineralPacketIntelsByKey[intel.playerObjectKey()] = intel
	}

	d.salvageIntelsByKey = make(map[int]*SalvageIntel, len(player.SalvageIntels))
	for i := range player.SalvageIntels {
		intel := &player.SalvageIntels[i]
		d.salvageIntelsByKey[intel.Num] = intel
	}

	d.wormholeIntelsByKey = make(map[int]*WormholeIntel, len(player.WormholeIntels))
	for i := range player.WormholeIntels {
		intel := &player.WormholeIntels[i]
		d.wormholeIntelsByKey[intel.Num] = intel
	}

	d.mysteryTraderIntelsByKey = make(map[int]*MysteryTraderIntel, len(player.MysteryTraderIntels))
	for i := range player.MysteryTraderIntels {
		intel := &player.MysteryTraderIntels[i]
		d.mysteryTraderIntelsByKey[intel.Num] = intel
	}

	d.designIntelsByKey = make(map[playerObject]*ShipDesignIntel, len(player.ShipDesignIntels))
	for i := range player.ShipDesignIntels {
		intel := &player.ShipDesignIntels[i]
		d.designIntelsByKey[intel.playerObjectKey()] = intel
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

func (intel *Intel) playerObjectKey() playerObject {
	return playerObjectKey(intel.PlayerNum, intel.Num)
}

func (intel *Intel) Owned() bool {
	return intel.PlayerNum != Unowned
}

type PlanetIntel struct {
	MapObjectIntel
	Hab                           Hab         `json:"hab,omitempty"`
	BaseHab                       Hab         `json:"baseHab,omitempty"`
	MineralConcentration          Mineral     `json:"mineralConcentration,omitempty"`
	Starbase                      *FleetIntel `json:"starbase,omitempty"`
	Cargo                         Cargo       `json:"cargo,omitempty"`
	CargoDiscovered               bool        `json:"cargoDiscovered,omitempty"`
	PlanetHabitability            int         `json:"planetHabitability,omitempty"`
	PlanetHabitabilityTerraformed int         `json:"planetHabitabilityTerraformed,omitempty"`
	Spec                          PlanetSpec  `json:"spec,omitempty"`
}

type ShipDesignIntel struct {
	Intel
	Hull          string           `json:"hull,omitempty"`
	HullSetNumber int              `json:"hullSetNumber,omitempty"`
	Version       int              `json:"version,omitempty"`
	Slots         []ShipDesignSlot `json:"slots,omitempty"`
	Spec          ShipDesignSpec   `json:"spec,omitempty"`
}

type FleetIntel struct {
	MapObjectIntel
	PlanetIntelID     int64       `json:"-,omitempty"` // for starbase fleets that are owned by a planet
	Heading           Vector      `json:"heading,omitempty"`
	OrbitingPlanetNum int         `json:"orbitingPlanetNum,omitempty"`
	WarpSpeed         int         `json:"warpSpeed,omitempty"`
	Mass              int         `json:"mass,omitempty"`
	Cargo             Cargo       `json:"cargo,omitempty"`
	CargoDiscovered   bool        `json:"cargoDiscovered,omitempty"`
	Tokens            []ShipToken `json:"tokens,omitempty"`
}

type MineralPacketIntel struct {
	MapObjectIntel
	WarpSpeed       int    `json:"warpSpeed,omitempty"`
	Heading         Vector `json:"heading"`
	Cargo           Cargo  `json:"cargo,omitempty"`
	TargetPlanetNum int    `json:"targetPlanetNum,omitempty"`
}

type SalvageIntel struct {
	MapObjectIntel
	Cargo Cargo `json:"cargo,omitempty"`
}

type MineFieldIntel struct {
	MapObjectIntel
	NumMines      int           `json:"numMines"`
	MineFieldType MineFieldType `json:"mineFieldType"`
	Spec          MineFieldSpec `json:"spec"`
}

type WormholeIntel struct {
	MapObjectIntel
	DestinationNum int               `json:"destinationNum,omitempty"`
	Stability      WormholeStability `json:"stability,omitempty"`
}

type MysteryTraderIntel struct {
	MapObjectIntel
	WarpSpeed int    `json:"warpSpeed,omitempty"`
	Heading   Vector `json:"heading"`
}

type PlayerIntel struct {
	Name           string `json:"name,omitempty"`
	Num            int    `json:"num,omitempty"`
	Color          string `json:"color,omitempty"`
	Seen           bool   `json:"seen,omitempty"`
	RaceName       string `json:"raceName,omitempty"`
	RacePluralName string `json:"racePluralName,omitempty"`
}

type ScoreIntel struct {
	ScoreHistory []PlayerScore `json:"scoreHistory"`
}

func (p *PlanetIntel) String() string {
	return fmt.Sprintf("Planet %s", &p.MapObjectIntel)
}

func (f *FleetIntel) String() string {
	return fmt.Sprintf("Player: %d, Fleet: %s", f.PlayerNum, f.Name)
}

func (f *SalvageIntel) String() string {
	return fmt.Sprintf("Player: %d, Salvage: %s", f.PlayerNum, f.Name)
}

func (f *MineFieldIntel) String() string {
	return fmt.Sprintf("Player: %d, MineField: %s", f.PlayerNum, f.Name)
}

func (f *MineralPacketIntel) String() string {
	return fmt.Sprintf("Player: %d, MineralPacket: %s", f.PlayerNum, f.Name)
}

func (d *ShipDesignIntel) String() string {
	return fmt.Sprintf("Player: %d, Fleet: %s", d.PlayerNum, d.Name)
}

// create a new FleetIntel object by key
func newFleetIntel(playerNum int, num int) FleetIntel {
	return FleetIntel{
		MapObjectIntel: MapObjectIntel{
			Type: MapObjectTypeFleet,
			Intel: Intel{
				PlayerNum: playerNum,
				Num:       num,
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

// create a new SalvageIntel object by key
func newSalvageIntel(playerNum int, num int) *SalvageIntel {
	return &SalvageIntel{
		MapObjectIntel: MapObjectIntel{
			Type: MapObjectTypeSalvage,
			Intel: Intel{
				PlayerNum: playerNum,
				Num:       num,
			},
		},
	}
}

// create a new MineFieldIntel object by key
func newMineFieldIntel(playerNum int, num int) *MineFieldIntel {
	return &MineFieldIntel{
		MapObjectIntel: MapObjectIntel{
			Type: MapObjectTypeMineField,
			Intel: Intel{
				PlayerNum: playerNum,
				Num:       num,
			},
		},
	}
}

// create a new MineralPacketIntel object by key
func newMineralPacketIntel(playerNum int, num int) *MineralPacketIntel {
	return &MineralPacketIntel{
		MapObjectIntel: MapObjectIntel{
			Type: MapObjectTypeMineralPacket,
			Intel: Intel{
				PlayerNum: playerNum,
				Num:       num,
			},
		},
	}
}

// create a new MysteryTraderIntel object by key
func newMysteryTraderIntel(num int) *MysteryTraderIntel {
	return &MysteryTraderIntel{
		MapObjectIntel: MapObjectIntel{
			Type: MapObjectTypeMysteryTrader,
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
	d.fleetIntelsByKey = make(map[playerObject]*FleetIntel)
	d.player.SalvageIntels = []SalvageIntel{}
	d.salvageIntelsByKey = make(map[int]*SalvageIntel)
	d.player.MineralPacketIntels = []MineralPacketIntel{}
	d.mineralPacketIntelsByKey = make(map[playerObject]*MineralPacketIntel)
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

		if !ownedByPlayer && intel.ReportAge == ReportAgeUnexplored {
			// let the player know we discovered a new planet
			messager.planetDiscovered(player, planet)
		}

		// if we pen scanned the planet, we learn some things
		intel.ReportAge = 0
		intel.Hab = planet.Hab
		intel.BaseHab = planet.BaseHab
		intel.MineralConcentration = planet.MineralConcentration
		intel.Spec.Habitability = player.Race.GetPlanetHabitability(intel.Hab)

		// terraforming
		terraformer := NewTerraformer()
		intel.Spec.TerraformAmount = terraformer.getTerraformAmount(intel.Hab, intel.BaseHab, player, player)
		intel.Spec.MinTerraformAmount = terraformer.getMinTerraformAmount(intel.Hab, intel.BaseHab, player, player)
		intel.Spec.CanTerraform = intel.Spec.TerraformAmount.absSum() > 0
		intel.Spec.TerraformedHabitability = player.Race.GetPlanetHabitability(planet.Hab.Add(intel.Spec.TerraformAmount))
		intel.Spec.MaxPopulation = planet.getMaxPopulation(rules, player, intel.Spec.Habitability)

		// discover starbases on scan, but don't discover designs
		intel.Spec.HasStarbase = planet.Spec.HasStarbase
		intel.Spec.HasMassDriver = planet.Spec.HasMassDriver
		intel.Spec.HasStargate = planet.Spec.HasStargate

		// these should never be nil...
		if !ownedByPlayer && planet.Spec.HasStarbase && planet.Starbase != nil && planet.Starbase.Tokens[0].design != nil {
			design := planet.Starbase.Tokens[0].design
			intel.Spec.StarbaseDesignName = design.Name
			intel.Spec.StarbaseDesignNum = design.Num
			d.discoverDesign(player, design, false)
		}

		// players know their planet pops, but other planets are slightly off
		if ownedByPlayer {
			intel.Spec.Population = planet.population()
		} else {
			var randomPopulationError = rules.random.Float64()*(rules.PopulationScannerError*2) - rules.PopulationScannerError
			intel.Spec.Population = maxInt(0, int(float64(planet.population())*(1-randomPopulationError)))
		}
	}
	return nil
}

// discover a planet's starbase specs after a battle
func (d *discover) discoverPlanetStarbase(player *Player, planet *Planet) error {
	var intel *PlanetIntel
	planetIndex := planet.Num - 1

	if planetIndex < 0 || planetIndex >= len(player.PlanetIntels) {
		return fmt.Errorf("player %s cannot discover planet %s, planetIndex %d out of range", player, planet, planetIndex)
	}

	intel = &player.PlanetIntels[planetIndex]

	// discover starbases on scan, but don't discover designs
	intel.Spec.HasStarbase = planet.Spec.HasStarbase
	intel.Spec.HasStargate = planet.Spec.HasStargate
	intel.Spec.HasMassDriver = planet.Spec.HasMassDriver

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

func (d *discover) discoverPlanetTerraformability(player *Player, planetNum int) error {
	var intel *PlanetIntel
	planetIndex := planetNum - 1

	if planetIndex < 0 || planetIndex >= len(player.PlanetIntels) {
		return fmt.Errorf("planetIndex %d out of range", planetIndex)
	}

	intel = &player.PlanetIntels[planetIndex]

	// if we've discovered this planet before, update the terraform stats
	if intel.ReportAge != ReportAgeUnexplored {
		// terraforming
		terraformer := NewTerraformer()
		intel.Spec.TerraformAmount = terraformer.getTerraformAmount(intel.Hab, intel.BaseHab, player, player)
		intel.Spec.MinTerraformAmount = terraformer.getMinTerraformAmount(intel.Hab, intel.BaseHab, player, player)
		intel.Spec.CanTerraform = intel.Spec.TerraformAmount.absSum() > 0
		intel.Spec.TerraformedHabitability = player.Race.GetPlanetHabitability(intel.Hab.Add(intel.Spec.TerraformAmount))
	}
	return nil
}

// discover a fleet and add it to the player's fleet intel
func (d *discover) discoverFleet(player *Player, fleet *Fleet) {
	intel := newFleetIntel(fleet.PlayerNum, fleet.Num)

	intel.Name = fleet.Name
	intel.Position = fleet.Position
	intel.OrbitingPlanetNum = fleet.OrbitingPlanetNum
	intel.Heading = fleet.Heading
	intel.WarpSpeed = fleet.WarpSpeed
	intel.Mass = fleet.Spec.Mass
	intel.Tokens = fleet.Tokens

	player.FleetIntels = append(player.FleetIntels, intel)
	d.fleetIntelsByKey[playerObjectKey(intel.PlayerNum, intel.Num)] = &intel
}

// discover cargo for an existing fleet
func (d *discover) discoverFleetCargo(player *Player, fleet *Fleet) {
	existingIntel, found := d.fleetIntelsByKey[playerObjectKey(fleet.PlayerNum, fleet.Num)]
	if found {
		existingIntel.Cargo = fleet.Cargo
		existingIntel.CargoDiscovered = true
	}

}

// discover a salvage and add it to the player's salvage intel
func (d *discover) discoverSalvage(salvage *Salvage) {
	player := d.player
	intel, found := d.salvageIntelsByKey[salvage.Num]
	if !found {
		// discover this new wormhole
		player.SalvageIntels = append(player.SalvageIntels, *newSalvageIntel(salvage.PlayerNum, salvage.Num))
		intel = &player.SalvageIntels[len(player.SalvageIntels)-1]
		d.salvageIntelsByKey[intel.Num] = intel
	}

	intel.Name = salvage.Name
	intel.PlayerNum = salvage.PlayerNum
	intel.Position = salvage.Position
	intel.Cargo = salvage.Cargo

}

// discover a mineField and add it to the player's mineField intel
func (d *discover) discoverMineField(player *Player, mineField *MineField) {
	key := playerObjectKey(mineField.PlayerNum, mineField.Num)
	intel, found := d.mineFieldIntelsByKey[key]
	if !found {
		// discover this new mineField
		intel = newMineFieldIntel(mineField.PlayerNum, mineField.Num)
		player.MineFieldIntels = append(player.MineFieldIntels, *intel)
		intel = &player.MineFieldIntels[len(player.MineFieldIntels)-1]
		d.mineFieldIntelsByKey[intel.playerObjectKey()] = intel
	}

	intel.Name = mineField.Name
	intel.Position = mineField.Position
	intel.MineFieldType = mineField.MineFieldType
	intel.NumMines = mineField.NumMines
	intel.Spec.Radius = mineField.Spec.Radius
}

// discover a mineralPacket and add it to the player's mineralPacket intel
func (d *discover) discoverMineralPacket(rules *Rules, player *Player, mineralPacket *MineralPacket, packetPlayer *Player, target *Planet) {
	key := playerObjectKey(mineralPacket.PlayerNum, mineralPacket.Num)
	intel, found := d.mineralPacketIntelsByKey[key]
	if !found {
		// discover this new mineralPacket
		intel = newMineralPacketIntel(mineralPacket.PlayerNum, mineralPacket.Num)
		player.MineralPacketIntels = append(player.MineralPacketIntels, *intel)
		intel = &player.MineralPacketIntels[len(player.MineralPacketIntels)-1]
		d.mineralPacketIntelsByKey[intel.playerObjectKey()] = intel
	}

	if player.Num != mineralPacket.PlayerNum {
		if target.PlayerNum == player.Num {
			damage := mineralPacket.getDamage(rules, target)
			messager.mineralPacketDiscoveredTargettingPlayer(player, mineralPacket, packetPlayer, target, damage)
		} else {
			messager.mineralPacketDiscovered(player, mineralPacket, packetPlayer, target)
		}
	}

	intel.Name = mineralPacket.Name
	intel.Position = mineralPacket.Position
	intel.Heading = mineralPacket.Heading
	intel.WarpSpeed = mineralPacket.WarpSpeed
	intel.Cargo = mineralPacket.Cargo
	intel.TargetPlanetNum = mineralPacket.TargetPlanetNum
}

// discover a player's design. This is a noop if we already know about
// the design and aren't discovering slots
func (d *discover) discoverDesign(player *Player, design *ShipDesign, discoverSlots bool) {
	key := playerObjectKey(design.PlayerNum, design.Num)
	intel, found := d.designIntelsByKey[key]
	if !found {
		// create a new intel for this design
		intel = &ShipDesignIntel{
			Intel: Intel{
				Name:      design.Name,
				PlayerNum: design.PlayerNum,
				Num:       design.Num,
			},
			Hull:          design.Hull,
			HullSetNumber: design.HullSetNumber,
		}

		// save this new design to our intel
		player.ShipDesignIntels = append(player.ShipDesignIntels, *intel)
		intel = &player.ShipDesignIntels[len(player.ShipDesignIntels)-1]
		d.designIntelsByKey[key] = intel
	}

	// discover slots if we haven't already and this scanner discovers them
	if discoverSlots && len(intel.Slots) == 0 {
		intel.Slots = make([]ShipDesignSlot, len(design.Slots))
		copy(intel.Slots, design.Slots)
		intel.Spec.Mass = design.Spec.Mass
		intel.Spec.Armor = design.Spec.Armor
		intel.Spec.Shields = design.Spec.Shields
		intel.Spec.PowerRating = design.Spec.PowerRating
		intel.Spec.FuelCapacity = design.Spec.FuelCapacity
		intel.Spec.Movement = design.Spec.Movement
		intel.Spec.Initiative = design.Spec.Initiative
		intel.Spec.CloakPercent = design.Spec.CloakPercent
	}
}

// discover a wormhole and add it to the player's wormhole intel
func (d *discover) discoverWormhole(wormhole *Wormhole) {
	player := d.player
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

func (d *discover) discoverWormholeLink(wormhole1, wormhole2 *Wormhole) {
	player := d.player
	intel1, found := d.wormholeIntelsByKey[wormhole1.Num]
	if !found {
		// discover this new wormhole
		player.WormholeIntels = append(player.WormholeIntels, *newWormholeIntel(wormhole1.Num))
		intel1 = &player.WormholeIntels[len(player.WormholeIntels)-1]
		d.wormholeIntelsByKey[intel1.Num] = intel1
	}

	intel2, found := d.wormholeIntelsByKey[wormhole2.Num]
	if !found {
		// discover this new wormhole
		player.WormholeIntels = append(player.WormholeIntels, *newWormholeIntel(wormhole2.Num))
		intel2 = &player.WormholeIntels[len(player.WormholeIntels)-1]
		d.wormholeIntelsByKey[intel2.Num] = intel2
	}

	intel1.Name = wormhole1.Name
	intel1.Position = wormhole1.Position
	intel1.Stability = wormhole1.Stability
	intel1.DestinationNum = wormhole1.DestinationNum

	intel2.Name = wormhole2.Name
	intel2.Position = wormhole2.Position
	intel2.Stability = wormhole2.Stability
	intel2.DestinationNum = wormhole2.DestinationNum
}

// discover a mysteryTrader and add it to the player's mysteryTrader intel
func (d *discover) discoverMysteryTrader(mysteryTrader *MysteryTrader) {
	player := d.player
	intel, found := d.mysteryTraderIntelsByKey[mysteryTrader.Num]
	if !found {
		// discover this new mysteryTrader
		player.MysteryTraderIntels = append(player.MysteryTraderIntels, *newMysteryTraderIntel(mysteryTrader.Num))
		intel = &player.MysteryTraderIntels[len(player.MysteryTraderIntels)-1]
		d.mysteryTraderIntelsByKey[intel.Num] = intel
	}

	intel.Name = mysteryTrader.Name
	intel.Position = mysteryTrader.Position
	intel.Heading = mysteryTrader.Heading
	intel.WarpSpeed = mysteryTrader.WarpSpeed
}

// discover a player's race
func (d *discover) discoverPlayer(player *Player) {
	intel := &d.player.PlayerIntels.PlayerIntels[player.Num-1]

	if !intel.Seen {
		log.Debug().Msgf("player %s discovered %s", d.player.Name, player.Name)
		messager.playerDiscovered(d.player, player)
		intel.Seen = true
		intel.Name = player.Name
		intel.Color = player.Color
		intel.RaceName = player.Race.Name
		intel.RacePluralName = player.Race.PluralName
	}
}

// discover a player's score
func (d *discover) discoverPlayerScores(player *Player) {
	intel := &d.player.PlayerIntels.ScoreIntels[player.Num-1]

	intel.ScoreHistory = make([]PlayerScore, len(player.ScoreHistory))
	copy(intel.ScoreHistory, player.ScoreHistory)
}

func (d *discover) getPlanetIntel(num int) *PlanetIntel {
	return &d.player.PlanetIntels[num-1]
}

func (d *discover) getWormholeIntel(num int) *WormholeIntel {
	return d.wormholeIntelsByKey[num]
}

func (d *discover) getMysteryTraderIntel(num int) *MysteryTraderIntel {
	return d.mysteryTraderIntelsByKey[num]
}

func (d *discover) getMineFieldIntel(playerNum, num int) *MineFieldIntel {
	key := playerObjectKey(playerNum, num)
	return d.mineFieldIntelsByKey[key]
}

func (d *discover) getMineralPacketIntel(playerNum, num int) *MineralPacketIntel {
	key := playerObjectKey(playerNum, num)
	return d.mineralPacketIntelsByKey[key]
}

func (d *discover) getFleetIntel(playerNum, num int) *FleetIntel {
	key := playerObjectKey(playerNum, num)
	return d.fleetIntelsByKey[key]
}

func (d *discover) getSalvageIntel(num int) *SalvageIntel {
	return d.salvageIntelsByKey[num]
}
