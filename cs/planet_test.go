package cs

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func newTestPlayerPlanet() (player *Player, planet *Planet) {
	player = NewPlayer(1, NewRace())
	player.Race.Spec = computeRaceSpec(&player.Race, &rules)
	planet = &Planet{}
	planet.PlayerNum = player.Num
	planet.BaseHab = player.Race.HabCenter()
	planet.Hab = planet.BaseHab

	player.Spec = computePlayerSpec(player, &rules, []*Planet{planet})

	return player, planet
}

func testSpaceStation(player *Player, planet *Planet) *Fleet {
	fleet := &Fleet{
		MapObject: MapObject{
			PlayerNum: player.Num,
		},
		BaseName: "Starbase",
		Tokens: []ShipToken{
			{
				DesignNum: 1,
				Quantity:  1,
				design: NewShipDesign(player, 1).
					WithHull(SpaceStation.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: Laser.Name, HullSlotIndex: 2, Quantity: 8},
						{HullComponent: MoleSkinShield.Name, HullSlotIndex: 3, Quantity: 8},
						{HullComponent: Laser.Name, HullSlotIndex: 4, Quantity: 8},
						{HullComponent: MoleSkinShield.Name, HullSlotIndex: 6, Quantity: 8},
						{HullComponent: Laser.Name, HullSlotIndex: 8, Quantity: 8},
						{HullComponent: Laser.Name, HullSlotIndex: 10, Quantity: 8},
					}).
					WithSpec(&rules, player)},
		},
		battlePlan: &player.BattlePlans[0],
		FleetOrders: FleetOrders{
			Waypoints: []Waypoint{
				NewPositionWaypoint(planet.Position, 0),
			},
		},
	}
	fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
	return fleet
}

func testDeathStar(player *Player, planet *Planet) *Fleet {
	fleet := &Fleet{
		MapObject: MapObject{
			PlayerNum: player.Num,
		},
		BaseName: "Starbase",
		Tokens: []ShipToken{
			{
				DesignNum: 1,
				Quantity:  1,
				design: NewShipDesign(player, 1).
					WithHull(DeathStar.Name).
					WithSlots([]ShipDesignSlot{}).
					WithSpec(&rules, player)},
		},
		battlePlan: &player.BattlePlans[0],
		FleetOrders: FleetOrders{
			Waypoints: []Waypoint{
				NewPositionWaypoint(planet.Position, 0),
			},
		},
	}
	fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
	return fleet
}

func TestPlanet_String(t *testing.T) {

	tests := []struct {
		name string
		p    *Planet
		want string
	}{
		{"MapObject String()", &Planet{MapObject: MapObject{GameDBObject: GameDBObject{GameID: 1, ID: 2}, Num: 3, Name: "Bob's Revenge"}},
			"Planet GameID:     1, ID:     2, Num:   3 Bob's Revenge"},
		{"MapObject String()", &Planet{MapObject: MapObject{GameDBObject: GameDBObject{GameID: 12345, ID: 23456}, Num: 120, Name: "Craig's Planet"}},
			"Planet GameID: 12345, ID: 23456, Num: 120 Craig's Planet"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.String(); got != tt.want {
				t.Errorf("MapObject.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanet_innateMines(t *testing.T) {
	player := NewPlayer(1, &Race{Spec: RaceSpec{InnateMining: false}})
	planet := Planet{}
	planet.setPopulation(16000)

	if got := planet.innateMines(player, planet.population()); got != 0 {
		t.Errorf("Planet.GetInnateMines() = %v, want %v", got, 0)
	}

	// should get 40 mines for 16k pop when the player has innate mining
	player.Race.Spec.InnateMining = true
	player.Race.Spec.InnatePopulationFactor = .1
	if got := planet.innateMines(player, planet.population()); got != 12 {
		t.Errorf("Planet.GetInnateMines() = %v, want %v", got, 12)
	}

}

func TestPlanet_innateScanner(t *testing.T) {
	player := NewPlayer(1, &Race{Spec: RaceSpec{InnateMining: false}})
	planet := Planet{}
	planet.setPopulation(67300)

	if got := planet.innateScanner(player, planet.population()); got != 0 {
		t.Errorf("Planet.GetInnateMines() = %v, want %v", got, 0)
	}

	// should get 40 mines for 16k pop when the player has innate mining
	player.Race.Spec.InnateScanner = true
	player.Race.Spec.InnatePopulationFactor = .1
	if got := planet.innateScanner(player, planet.population()); got != 82 {
		t.Errorf("Planet.GetInnateMines() = %v, want %v", got, 82)
	}

}

func TestPlanet_getGrowthAmount(t *testing.T) {
	type fields struct {
		Hab        Hab
		Population int
	}
	type args struct {
		player        *Player
		maxPopulation int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{name: "empty planet", fields: fields{Hab{50, 50, 50}, 0}, args: args{NewPlayer(1, NewRace()), 1_000_000}, want: 0},
		{name: "less than 25% cap, grows at full 10% growth rate", fields: fields{Hab{50, 50, 50}, 100_000}, args: args{NewPlayer(1, NewRace()), 1_200_000}, want: 10_000},
		{name: "at 50% cap, it slows down in growth", fields: fields{Hab{50, 50, 50}, 600_000}, args: args{NewPlayer(1, NewRace()), 1_200_000}, want: 26_700},
		{name: "we are basicallly at capacity, we only grow a tiny amount", fields: fields{Hab{50, 50, 50}, 1_180_000}, args: args{NewPlayer(1, NewRace()), 1_200_000}, want: 100},
		{name: "no more growth past a certain capacity", fields: fields{Hab{50, 50, 50}, 1_190_000}, args: args{NewPlayer(1, NewRace()), 1_200_000}, want: 0},
		{name: "hostile planets kill off colonists", fields: fields{Hab{10, 15, 15}, 2500}, args: args{NewPlayer(1, NewRace()), 0}, want: -100},
		{name: "super hostile planet with 100k people, should be -45% habitable, so should kill off -4.5% of the pop", fields: fields{Hab{}, 100_000}, args: args{NewPlayer(1, NewRace()), 0}, want: -4500},
		{name: "double cap planet should kill off 4% of the pop", fields: fields{Hab{50, 50, 50}, 2_400_000}, args: args{NewPlayer(1, NewRace()), 1_200_000}, want: -96_000},
		{name: "5x cap planet should kill off max of 12% of the pop", fields: fields{Hab{50, 50, 50}, 6_000_000}, args: args{NewPlayer(1, NewRace()), 1_200_000}, want: -720_000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Planet{
				Cargo: Cargo{Colonists: tt.fields.Population / 100},
				Hab:   tt.fields.Hab,
			}
			// 10% growth for easier math
			tt.args.player.Race.GrowthRate = 10
			tt.args.player.Race.Spec = computeRaceSpec(&tt.args.player.Race, &rules)
			if got := p.getGrowthAmount(tt.args.player, tt.args.maxPopulation, rules.PopulationOvercrowdDieoffRate, rules.PopulationOvercrowdDieoffRateMax); got != tt.want {
				t.Errorf("Planet.getGrowthAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueueItemType_IsAuto(t *testing.T) {
	tests := []struct {
		name string
		tr   QueueItemType
		want bool
	}{
		{"Not Auto", QueueItemTypeFactory, false},
		{"Not Auto", QueueItemTypeAutoMaxTerraform, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.IsAuto(); got != tt.want {
				t.Errorf("QueueItemType.IsAuto() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanet_reduceMineralConcentration(t *testing.T) {
	type args struct {
		rules Rules
	}
	tests := []struct {
		name   string
		planet *Planet
		args   args
		want   Mineral
	}{
		{"Redcue empty planet min conc", NewPlanet(), args{rules}, Mineral{1, 1, 1}},
		{
			"150 mines should reduce 100% conc by 1 if we have 151 mineyears",
			NewPlanet().
				WithMineralConcentration(Mineral{100, 100, 100}).
				WithMines(150).
				WithMineYears(Mineral{151, 151, 151}),
			args{rules},
			Mineral{99, 99, 99},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.planet.reduceMineralConcentration(&tt.args.rules)

			if got := tt.planet.MineralConcentration; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Planet.reduceMineralConcentration() = %v, want %v", got, tt.want)
			}

		})
	}
}

func Test_getMaxPopulation(t *testing.T) {
	type args struct {
		hab int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"joat homeworld", args{100}, 1_200_000},
		{"low hab world", args{1}, 60_000},
		{"bad hab world", args{-45}, 60_000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planet := NewPlanet()
			player := NewPlayer(0, NewRace().WithSpec(&rules)).withSpec(&rules)
			if got := planet.getMaxPopulation(&rules, player, tt.args.hab); got != tt.want {
				t.Errorf("getMaxPopulation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_computePlanetSpec(t *testing.T) {
	player, planet := newTestPlayerPlanet()
	planet.Starbase = testSpaceStation(player, planet)

	player.Race.Spec.InnateScanner = true
	player.Race.Spec.InnatePopulationFactor = .1
	planet.setPopulation(67300)
	planet.Spec = computePlanetSpec(&rules, player, planet)

	assert.Equal(t, planet.Spec.ScanRange, 82)
	assert.Equal(t, planet.Spec.ScanRangePen, 0)

	// now use a death start
	planet.Starbase = testDeathStar(player, planet)
	planet.Spec = computePlanetSpec(&rules, player, planet)

	assert.Equal(t, planet.Spec.ScanRange, 82)
	assert.Equal(t, planet.Spec.ScanRangePen, 41)

}

func TestPlanet_randomize(t *testing.T) {

	type args struct {
		rng rng
	}
	tests := []struct {
		name string
		args args
		want Planet
	}{
		{
			name: "planet gen with all 0 rng",
			args: args{newIntRandom().addFloats(.50)},
			want: Planet{
				MapObject:            MapObject{Type: MapObjectTypePlanet, PlayerNum: Unowned},
				Dirty:                true,
				Hab:                  Hab{1, 1, 1},
				BaseHab:              Hab{1, 1, 1},
				MineralConcentration: Mineral{1, 31, 31},
				PlanetOrders: PlanetOrders{
					ProductionQueue: []ProductionQueueItem{},
				},
			},
		},
		{
			name: "planet with random artifact",
			args: args{newIntRandom().addFloats(.33)},
			want: Planet{
				MapObject:            MapObject{Type: MapObjectTypePlanet, PlayerNum: Unowned},
				Dirty:                true,
				Hab:                  Hab{1, 1, 1},
				BaseHab:              Hab{1, 1, 1},
				MineralConcentration: Mineral{1, 31, 31},
				RandomArtifact:       true,
				PlanetOrders: PlanetOrders{
					ProductionQueue: []ProductionQueueItem{},
				},
			},
		},
		{
			name: "planet with no random artifact",
			args: args{newIntRandom().addFloats(.50)},
			want: Planet{
				MapObject:            MapObject{Type: MapObjectTypePlanet, PlayerNum: Unowned},
				Dirty:                true,
				Hab:                  Hab{1, 1, 1},
				BaseHab:              Hab{1, 1, 1},
				MineralConcentration: Mineral{1, 31, 31},
				RandomArtifact:       false,
				PlanetOrders: PlanetOrders{
					ProductionQueue: []ProductionQueueItem{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPlanet()

			rules := NewRules()
			rules.random = tt.args.rng
			got.randomize(&rules)

			if !reflect.DeepEqual(got, &tt.want) {
				// dump json, but this won't include some fields
				test.CompareAsJSON(t, got, tt.want)
				t.Errorf("randomize() = %#v, want %#v", got, tt.want)
			}

		})
	}
}

func TestPlanet_grow(t *testing.T) {
	type fields struct {
		hab        Hab
		population int
	}
	type args struct {
		race *Race
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantPopulation int
	}{
		{"standard humanoid starter world", fields{hab: Hab{50, 50, 50}, population: 25000}, args{NewRace().WithSpec(&rules)}, 28800},
		{"full world", fields{hab: Hab{50, 50, 50}, population: 500_000}, args{NewRace().WithSpec(&rules)}, 545_400},
		{"hostile world", fields{hab: Hab{1, 1, 1}, population: 25000}, args{NewRace().WithSpec(&rules)}, 23900},
		{"hostile world, low pop", fields{hab: Hab{1, 1, 1}, population: 200}, args{NewRace().WithSpec(&rules)}, 100},
		{"hostile world, low pop 2", fields{hab: Hab{1, 1, 1}, population: 100}, args{NewRace().WithSpec(&rules)}, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(0, tt.args.race).WithNum(1)
			planet := NewPlanet().WithPlayerNum(player.Num)
			planet.Hab = tt.fields.hab
			planet.BaseHab = tt.fields.hab
			planet.setPopulation(tt.fields.population)
			planet.Spec = computePlanetSpec(&rules, player, planet)

			planet.grow(player)

			if planet.population() != tt.wantPopulation {
				t.Errorf("grow() = %v, want %v", planet.population(), tt.wantPopulation)
			}

		})
	}
}
