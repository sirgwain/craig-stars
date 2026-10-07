//go:build !wasi && !wasm

package cs

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"maps"
	"slices"
	"strconv"
	"strings"
)

// TestScenario describes a deterministic universe. References to designs and
// planets use names; BuildScenario resolves them before normal finalization.
type TestScenario struct {
	Name           string           `json:"name"`
	Players        []ScenarioPlayer `json:"players"`
	Planets        []ScenarioPlanet `json:"planets"`
	Wormholes      []Wormhole
	MysteryTraders []MysteryTrader
}

// ScenarioPlayer defaults to a humanoid with tech 3 in every field.
// Relations, when present, contain one entry per player, including this player.
type ScenarioPlayer struct {
	Player         *Player                 `json:"player,omitempty"`
	Relations      []PlayerRelationship    `json:"relations,omitempty"`
	Designs        []ShipDesign            `json:"designs,omitempty"`
	Fleets         []ScenarioFleet         `json:"fleets,omitempty"`
	Salvages       []Salvage               `json:"salvages,omitempty"`
	MineralPackets []ScenarioMineralPacket `json:"mineralPackets,omitempty"`
	Minefields     []Minefield             `json:"minefields,omitempty"`
}

// ScenarioPlanet defaults to hab 50/50/50 and concentration 100/100/100.
// Pointers distinguish omitted values from an explicitly hostile or depleted planet.
type ScenarioPlanet struct {
	Name                              string
	Owner                             int
	Position                          Vector
	Hab                               *Hab
	BaseHab                           *Hab
	Concentration                     *Mineral
	Homeworld                         bool
	Cargo                             Cargo
	PartialPopulation                 int
	Mines, Factories, Defenses        int
	Scanner                           bool
	Starbase                          string
	StarbaseDamage                    float64
	ProductionQueue                   []ScenarioProductionQueueItem
	RouteTo                           string
	PacketTo                          string
	PacketSpeed                       int
	ContributesOnlyLeftoverToResearch bool
}

// ScenarioValue makes an explicit optional value, including zero.
func ScenarioValue[T any](value T) *T { return &value }

// ScenarioFleet contains either one Design or mixed Tokens. Quantity defaults
// to one; Fuel defaults to capacity unless EmptyFuel is true. Waypoints describe
// the full route, including waypoint zero; omission creates a stationary route.
type ScenarioFleet struct {
	Name          string
	Design        string
	Quantity      int
	Tokens        []ScenarioShipToken
	At            string
	Position      Vector
	Fuel          int
	EmptyFuel     bool
	Cargo         Cargo
	Age           int
	Waypoints     []ScenarioWaypoint
	RepeatOrders  bool
	BattlePlanNum int
	Purpose       FleetPurpose
}

type ScenarioShipToken struct {
	Design          string
	Quantity        int
	QuantityDamaged int
	Damage          float64
}

type ScenarioWaypoint struct {
	To                           string
	Position                     Vector
	Warp                         int
	Task                         WaypointTask
	TransportTasks               WaypointTransportTasks
	WaitAtWaypoint               bool
	LayMinefieldDuration         int
	PatrolRange, PatrolWarpSpeed int
	TransferToPlayer             int
}

type ScenarioProductionQueueItem struct {
	Type      QueueItemType
	Design    string
	Quantity  int
	Allocated Cost
	Tags      Tags
}

type ScenarioMineralPacket struct {
	Name                     string
	Position                 Vector
	Cargo                    Cargo
	WarpSpeed, SafeWarpSpeed int
	To                       string
}

// Homeworld returns the standard homeworld used by browser and logic fixtures.
func Homeworld(name string, owner int) ScenarioPlanet {
	return ScenarioPlanet{Name: name, Owner: owner, Homeworld: true,
		Cargo: Cargo{Ironium: 1000, Boranium: 1000, Germanium: 1000, Colonists: 2500}}
}

func AIPlayer(name string) ScenarioPlayer {
	return ScenarioPlayer{Player: &Player{Name: name, AIControlled: true, Race: *NewRace()}}
}

// BuildScenario constructs fresh objects on every call, validates references,
// and initializes specs, intel and maps through the universe finalizer.
func BuildScenario(scenario TestScenario) *FullGame {
	b := &scenarioBuilder{scenario: scenario, game: newScenarioGame(scenario.Name), planets: map[string]*Planet{}}
	for _, sp := range scenario.Planets {
		b.addPlanet(sp)
	}
	for _, wormhole := range scenario.Wormholes {
		addScenarioWormhole(b.game, &wormhole)
	}
	for _, trader := range scenario.MysteryTraders {
		addScenarioMysteryTrader(b.game, &trader)
	}
	for index, sp := range scenario.Players {
		b.addPlayer(index, sp)
	}
	b.applyRelations()

	ug := NewUniverseGenerator(b.game.Game, b.game.Players)
	if err := ug.GenerateWithUniverse(b.game.Universe); err != nil {
		panic(fmt.Errorf("scenario %q: %w", scenario.Name, err))
	}

	b.fillFuel()
	return b.game
}

// scenarioBuilder holds the game being built and the name lookups that
// scenario references resolve against.
type scenarioBuilder struct {
	scenario TestScenario
	game     *FullGame
	planets  map[string]*Planet
}

// scenarioPlayerBuilder resolves design names for one player.
type scenarioPlayerBuilder struct {
	*scenarioBuilder
	player  *Player
	designs map[string]*ShipDesign
}

func (b *scenarioBuilder) planet(context, name string) *Planet {
	planet := b.planets[name]
	if planet == nil {
		panic(fmt.Sprintf("%s: unknown planet %q", context, name))
	}
	return planet
}

func (b *scenarioBuilder) addPlanet(sp ScenarioPlanet) {
	if sp.Name == "" {
		panic("scenario planet: name is required")
	}
	if b.planets[sp.Name] != nil {
		panic(fmt.Sprintf("duplicate planet %q", sp.Name))
	}
	if sp.Owner < 0 || sp.Owner > len(b.scenario.Players) {
		panic(fmt.Sprintf("planet %q: unknown owner %d", sp.Name, sp.Owner))
	}
	if sp.Starbase != "" && sp.Owner == 0 {
		panic(fmt.Sprintf("planet %q: starbase requires an owner", sp.Name))
	}

	hab := Hab{50, 50, 50}
	if sp.Hab != nil {
		hab = *sp.Hab
	}
	baseHab := hab
	if sp.BaseHab != nil {
		baseHab = *sp.BaseHab
	}
	concentration := NewMineral(100, 100, 100)
	if sp.Concentration != nil {
		concentration = *sp.Concentration
	}

	planet := &Planet{
		MapObject:            MapObject{Type: MapObjectTypePlanet, Num: len(b.game.Planets) + 1, Name: sp.Name, PlayerNum: sp.Owner, Position: sp.Position},
		Hab:                  hab,
		BaseHab:              baseHab,
		MineralConcentration: concentration,
		Homeworld:            sp.Homeworld,
		Cargo:                sp.Cargo,
		PartialPopulation:    sp.PartialPopulation,
		Mines:                sp.Mines,
		Factories:            sp.Factories,
		Defenses:             sp.Defenses,
		Scanner:              sp.Scanner,
		PlanetOrders:         PlanetOrders{PacketSpeed: sp.PacketSpeed, ContributesOnlyLeftoverToResearch: sp.ContributesOnlyLeftoverToResearch},
	}
	b.game.Planets = append(b.game.Planets, planet)
	b.planets[sp.Name] = planet
}

func (b *scenarioBuilder) addPlayer(index int, sp ScenarioPlayer) {
	template := defaultScenarioPlayer()
	if sp.Player != nil {
		// Clone player configuration; specs are recomputed by the finalizer.
		// Mutable maps and slices must not be shared between builds.
		template = Player{}
		data, err := json.Marshal(sp.Player)
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(data, &template); err != nil {
			panic(err)
		}
	}
	if template.Stats == nil {
		template.Stats = &PlayerStats{}
	}
	player := addScenarioPlayer(b.game, template.WithNum(index+1))
	if player.AIControlled {
		player.SubmittedTurn = true
	}

	pb := &scenarioPlayerBuilder{scenarioBuilder: b, player: player, designs: map[string]*ShipDesign{}}
	for _, sd := range sp.Designs {
		pb.addDesign(sd)
	}
	for fi, sf := range sp.Fleets {
		pb.addFleet(fi+1, sf)
	}
	for pi, planet := range b.game.Planets {
		if planet.PlayerNum == player.Num {
			pb.addPlanetOrders(planet, b.scenario.Planets[pi])
		}
	}
	for packetIndex, packet := range sp.MineralPackets {
		pb.addMineralPacket(packetIndex+1, packet)
	}
	for mi, field := range sp.Minefields {
		addScenarioMinefield(b.game, player, &field, mi+1)
	}
	for _, salvage := range sp.Salvages {
		addScenarioSalvage(b.game, player, &salvage)
	}
}

func (pb *scenarioPlayerBuilder) design(context, name string) *ShipDesign {
	design := pb.designs[name]
	if design == nil {
		panic(fmt.Sprintf("%s: unknown design %q for player %d", context, name, pb.player.Num))
	}
	return design
}

func (pb *scenarioPlayerBuilder) addDesign(sd ShipDesign) {
	if sd.Name == "" {
		panic(fmt.Sprintf("player %d: design name is required", pb.player.Num))
	}
	if pb.designs[sd.Name] != nil {
		panic(fmt.Sprintf("player %d: duplicate design %q", pb.player.Num, sd.Name))
	}
	design := addScenarioDesign(pb.player, &ShipDesign{Name: sd.Name, Hull: sd.Hull, Purpose: sd.Purpose, Slots: slices.Clone(sd.Slots)})
	pb.designs[sd.Name] = design
}

// fleetBaseName strips a trailing " #<number>" from a fleet name, so a fleet
// named "MT Meet-er #2" has the base name "MT Meet-er" wherever it is listed.
func fleetBaseName(name string) string {
	i := strings.LastIndex(name, " #")
	if i < 0 {
		return name
	}
	if _, err := strconv.Atoi(name[i+2:]); err != nil {
		return name
	}
	return name[:i]
}

func (pb *scenarioPlayerBuilder) addFleet(num int, sf ScenarioFleet) {
	player := pb.player

	tokens := sf.Tokens
	if sf.Design != "" {
		if len(tokens) > 0 {
			panic(fmt.Sprintf("fleet %q: use Design or Tokens, not both", sf.Name))
		}
		tokens = []ScenarioShipToken{{Design: sf.Design, Quantity: sf.Quantity}}
	}
	if len(tokens) == 0 {
		panic(fmt.Sprintf("player %d fleet %d: a design is required", player.Num, num))
	}

	name := sf.Name
	if name == "" {
		name = fmt.Sprintf("%s #%d", tokens[0].Design, num)
	}
	for _, f := range pb.game.Fleets {
		if f.PlayerNum == player.Num && f.Name == name {
			panic(fmt.Sprintf("player %d: duplicate fleet %q", player.Num, name))
		}
	}
	context := fmt.Sprintf("fleet %q", name)

	fleet := &Fleet{
		MapObject:   MapObject{Type: MapObjectTypeFleet, Num: num, PlayerNum: player.Num, Name: name, Position: sf.Position},
		BaseName:    fleetBaseName(name),
		Cargo:       sf.Cargo,
		Fuel:        sf.Fuel,
		Age:         sf.Age,
		FleetOrders: FleetOrders{RepeatOrders: sf.RepeatOrders, BattlePlanNum: sf.BattlePlanNum, Purpose: sf.Purpose},
	}
	for _, st := range tokens {
		quantity := st.Quantity
		if quantity == 0 {
			quantity = 1
		}
		if quantity < 0 {
			panic(fmt.Sprintf("%s: negative quantity", context))
		}
		design := pb.design(context, st.Design)
		fleet.Tokens = append(fleet.Tokens, ShipToken{DesignNum: design.Num, Quantity: quantity, QuantityDamaged: st.QuantityDamaged, Damage: st.Damage})
	}

	var orbiting *Planet
	if sf.At != "" {
		orbiting = pb.planet(context, sf.At)
		fleet.Position = orbiting.Position
		fleet.OrbitingPlanetNum = orbiting.Num
	}

	if len(sf.Waypoints) == 0 {
		// stationary: a single waypoint where the fleet is
		wp := NewPositionWaypoint(fleet.Position, 0)
		if orbiting != nil {
			wp = NewPlanetWaypoint(orbiting.Position, orbiting.Num, orbiting.Name, 0)
		}
		fleet.Waypoints = []Waypoint{wp}
	}
	for _, sw := range sf.Waypoints {
		fleet.Waypoints = append(fleet.Waypoints, pb.waypoint(context, sw))
	}

	pb.game.Fleets = append(pb.game.Fleets, fleet)
}

// waypoint resolves a waypoint target by planet name, then mystery trader name,
// falling back to a position when To is empty.
func (b *scenarioBuilder) waypoint(context string, sw ScenarioWaypoint) Waypoint {
	wp := NewPositionWaypoint(sw.Position, sw.Warp)
	if sw.To != "" {
		if planet := b.planets[sw.To]; planet != nil {
			wp = NewPlanetWaypoint(planet.Position, planet.Num, planet.Name, sw.Warp)
		} else if trader := b.mysteryTrader(sw.To); trader != nil {
			wp = NewMysteryTraderWaypoint(trader, sw.Warp)
		} else {
			panic(fmt.Sprintf("%s: unknown planet or trader %q", context, sw.To))
		}
	}
	wp.Task = sw.Task
	wp.TransportTasks = sw.TransportTasks
	wp.WaitAtWaypoint = sw.WaitAtWaypoint
	wp.LayMinefieldDuration = sw.LayMinefieldDuration
	wp.PatrolRange = sw.PatrolRange
	wp.PatrolWarpSpeed = sw.PatrolWarpSpeed
	wp.TransferToPlayer = sw.TransferToPlayer
	return wp
}

func (b *scenarioBuilder) mysteryTrader(name string) *MysteryTrader {
	for _, trader := range b.game.MysteryTraders {
		if trader.Name == name {
			return trader
		}
	}
	return nil
}

// addPlanetOrders adds the parts of an owned planet that refer to the owner's
// designs or to other planets: its starbase, production queue and targets.
func (pb *scenarioPlayerBuilder) addPlanetOrders(planet *Planet, source ScenarioPlanet) {
	player := pb.player
	context := fmt.Sprintf("planet %q", planet.Name)

	if source.Starbase != "" {
		design := pb.design(context, source.Starbase)
		starbase := &Fleet{
			MapObject:   MapObject{Type: MapObjectTypeFleet, Name: design.Name, PlayerNum: player.Num, Position: planet.Position},
			BaseName:    design.Name,
			PlanetNum:   planet.Num,
			Starbase:    true,
			Tokens:      []ShipToken{{DesignNum: design.Num, Quantity: 1}},
			FleetOrders: FleetOrders{Waypoints: []Waypoint{NewPositionWaypoint(planet.Position, 0)}},
		}
		if source.StarbaseDamage > 0 {
			starbase.Tokens[0].QuantityDamaged = 1
			starbase.Tokens[0].Damage = source.StarbaseDamage
		}
		planet.Starbase = starbase
		pb.game.Starbases = append(pb.game.Starbases, starbase)
	}

	for _, sq := range source.ProductionQueue {
		q := ProductionQueueItem{Type: sq.Type, Quantity: sq.Quantity, Allocated: sq.Allocated, Tags: maps.Clone(sq.Tags)}
		if sq.Design != "" {
			q.DesignNum = pb.design(context, sq.Design).Num
		}
		planet.ProductionQueue = append(planet.ProductionQueue, q)
	}

	if source.RouteTo != "" {
		planet.RouteTargetType = MapObjectTypePlanet
		planet.RouteTargetNum = pb.planet(context, source.RouteTo).Num
	}
	if source.PacketTo != "" {
		planet.PacketTargetNum = pb.planet(context, source.PacketTo).Num
	}
}

func (pb *scenarioPlayerBuilder) addMineralPacket(num int, packet ScenarioMineralPacket) {
	target := pb.planet(fmt.Sprintf("packet %q", packet.Name), packet.To)
	mp := &MineralPacket{
		MapObject:       MapObject{Position: packet.Position},
		Cargo:           packet.Cargo,
		WarpSpeed:       packet.WarpSpeed,
		SafeWarpSpeed:   packet.SafeWarpSpeed,
		TargetPlanetNum: target.Num,
	}
	addScenarioMineralPacket(pb.game, pb.player, mp, num)
	if packet.Name != "" {
		mp.Name = packet.Name
	}
}

// applyRelations sets explicit relations before GenerateWithUniverse
// initializes discoverers and scans, so alliances affect initial intel too.
// Players without explicit relations get the defaults from the finalizer.
func (b *scenarioBuilder) applyRelations() {
	for index, sp := range b.scenario.Players {
		player := b.game.Players[index]
		relations := player.Relations
		if sp.Relations != nil {
			relations = sp.Relations
		}
		if relations == nil {
			continue
		}
		if len(relations) != len(b.scenario.Players) {
			panic(fmt.Sprintf("player %d: relations must contain %d entries", index+1, len(b.scenario.Players)))
		}
		player.Relations = slices.Clone(relations)
	}
}

// fillFuel gives fleets a full tank unless the scenario set Fuel or EmptyFuel.
// It runs after the finalizer because fuel capacity comes from the fleet spec.
func (b *scenarioBuilder) fillFuel() {
	fleets := b.game.Fleets
	for _, sp := range b.scenario.Players {
		for _, sf := range sp.Fleets {
			fleet := fleets[0]
			fleets = fleets[1:]
			if sf.Fuel == 0 && !sf.EmptyFuel {
				fleet.Fuel = fleet.Spec.FuelCapacity
			}
		}
	}
}

var scenarioColors = []string{
	"#0000FF",
	"#C33232",
	"#1F8BA7",
	"#43A43E",
	"#8D29CB",
	"#B88628",
	"#FF4500",
	"#FF8C00",
	"#008000",
	"#00FA9A",
	"#7FFFD4",
	"#8A2BE2",
	"#FF1493",
	"#D2691E",
	"#F0FFF0",
}

func defaultScenarioPlayer() Player {
	return Player{
		UserID:        1,
		Race:          *NewRace(),
		TechLevels:    TechLevel{Energy: 3, Weapons: 3, Propulsion: 3, Construction: 3, Electronics: 3, Biotechnology: 3},
		Stats:         &PlayerStats{},
		AcquiredTechs: map[string]bool{},
		PlayerOrders: PlayerOrders{
			Researching:       Energy,
			ResearchAmount:    15,
			NextResearchField: NextResearchFieldSameField,
			CargoTransfers:    CargoTransfers{},
		},
	}
}

func newScenarioGame(name string) *FullGame {
	client := NewGamer()
	game := client.CreateGame(1, *NewGameSettings().WithName(name))
	game.RandomEvents = false
	game.Area = Vector{X: 200, Y: 200}
	game.Seed = 0
	game.Rules.ResetSeed(0)
	game.State = GameStateWaitingForPlayers
	universe := NewUniverse(slog.Default(), &game.Rules)

	return &FullGame{
		Game:      game,
		Universe:  &universe,
		TechStore: &StaticTechStore,
		Players:   []*Player{},
	}
}

func addScenarioPlayer(game *FullGame, player *Player) *Player {
	player.Num = len(game.Players) + 1
	player.Color = scenarioColors[player.Num-1]

	if player.Name == "" {
		player.Name = fmt.Sprintf("Player #%d", player.Num)
	}

	game.Players = append(game.Players, player)
	return player
}

func addScenarioDesign(player *Player, design *ShipDesign) *ShipDesign {
	design.Num = len(player.Designs) + 1
	design.PlayerNum = player.Num
	player.Designs = append(player.Designs, design)
	return design
}

func addScenarioMineralPacket(game *FullGame, player *Player, mineralPacket *MineralPacket, num int) *MineralPacket {
	mineralPacket.Type = MapObjectTypeMineralPacket
	mineralPacket.PlayerNum = player.Num
	mineralPacket.Num = num
	mineralPacket.Name = fmt.Sprintf("%s Mineral Packet #%d", player.Race.PluralName, mineralPacket.Num)
	mineralPacket.Heading = (game.Planets[mineralPacket.TargetPlanetNum-1].Position.Subtract(mineralPacket.Position)).Normalized()

	game.MineralPackets = append(game.MineralPackets, mineralPacket)
	return mineralPacket
}

func addScenarioMinefield(game *FullGame, player *Player, minefield *Minefield, num int) *Minefield {
	minefield.Tags = maps.Clone(minefield.Tags)
	minefield.Type = MapObjectTypeMinefield
	minefield.PlayerNum = player.Num
	minefield.Num = len(game.Minefields) + 1
	minefield.Name = fmt.Sprintf("%s %s Minefield #%d", player.Race.PluralName, minefield.MinefieldType.String(), num)
	game.Minefields = append(game.Minefields, minefield)
	return minefield
}

func addScenarioSalvage(game *FullGame, player *Player, salvage *Salvage) *Salvage {
	salvage.Tags = maps.Clone(salvage.Tags)
	salvage.Type = MapObjectTypeSalvage
	salvage.PlayerNum = player.Num
	salvage.Num = len(game.Salvages) + 1
	salvage.Name = fmt.Sprintf("Salvage #%d", salvage.Num)
	game.Salvages = append(game.Salvages, salvage)
	return salvage
}

func addScenarioWormhole(game *FullGame, wormhole *Wormhole) *Wormhole {
	wormhole.Tags = maps.Clone(wormhole.Tags)
	wormhole.Type = MapObjectTypeWormhole
	wormhole.Num = len(game.Wormholes) + 1
	game.Wormholes = append(game.Wormholes, wormhole)
	return wormhole
}

func addScenarioMysteryTrader(game *FullGame, mysteryTrader *MysteryTrader) *MysteryTrader {
	mysteryTrader.Tags = maps.Clone(mysteryTrader.Tags)
	mysteryTrader.Type = MapObjectTypeMysteryTrader
	mysteryTrader.Num = len(game.MysteryTraders) + 1
	mysteryTrader.Name = fmt.Sprintf("Mystery Trader #%d", mysteryTrader.Num)
	mysteryTrader.Heading = (mysteryTrader.Destination.Subtract(mysteryTrader.Position)).Normalized()
	mysteryTrader.PlayersRewarded = maps.Clone(mysteryTrader.PlayersRewarded)
	if mysteryTrader.PlayersRewarded == nil {
		mysteryTrader.PlayersRewarded = map[int]bool{}
	}
	game.MysteryTraders = append(game.MysteryTraders, mysteryTrader)
	return mysteryTrader
}
