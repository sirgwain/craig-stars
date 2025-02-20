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
		{"MapObject String()", &Planet{GameDBObject: GameDBObject{GameID: 1, ID: 2}, MapObject: MapObject{Num: 3, Name: "Bob's Revenge"}},
			"Planet GameID:     1, ID:     2, Num:   3 Bob's Revenge"},
		{"MapObject String()", &Planet{GameDBObject: GameDBObject{GameID: 12345, ID: 23456}, MapObject: MapObject{Num: 120, Name: "Craig's Planet"}},
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

func Test_innateMines(t *testing.T) {
	tests := []struct {
		name      string
		pop       int
		popFactor float64
		want      int
	}{
		{name: "100 pop", pop: 100, popFactor: 0.1, want: 1},
		{name: "10K pop", pop: 10_000, popFactor: 0.1, want: 10},
		{name: "40K pop", pop: 40_000, popFactor: 0.1, want: 20},
		{name: "1M pop", pop: 1_000_000, popFactor: 0.1, want: 100},
	}
	for _, tt := range tests {
		planet := NewPlanet()
		planet.Name = tt.name
		if got := innateMines(tt.popFactor, tt.pop); got != tt.want {
			t.Errorf("planet.GetInnateMines() = %v, want %v", got, tt.want)
		}
	}
}

func Test_innateScanner(t *testing.T) {
	tests := []struct {
		name      string
		pop       int
		popFactor float64
		want      int
	}{
		{name: "100 pop", pop: 100, popFactor: 0.1, want: 3},       // sqrt(10)
		{name: "10K pop", pop: 10_000, popFactor: 0.1, want: 31},   // sqrt(1000)
		{name: "50.5K pop", pop: 50_500, popFactor: 0.1, want: 71}, // sqrt(5050) ≈ 71
		{name: "144K pop", pop: 144_000, popFactor: 0.1, want: 120},
		{name: "6.4M pop", pop: 6_400_000, popFactor: 0.1, want: 800},
	}
	for _, tt := range tests {
		planet := NewPlanet()
		planet.Name = tt.name
		if got := innateScanner(tt.popFactor, tt.pop); got != tt.want {
			t.Errorf("innateScanner() = %v, want %v", got, tt.want)
		}
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
	tests := []struct {
		name string
		hab  int
		want int
	}{
		{"joat homeworld", 100, 1_200_000},
		{"low hab world", 1, 60_000},
		{"bad hab world", -45, 60_000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planet := NewPlanet()
			player := NewPlayer(0, NewRace().WithSpec(&rules)).withSpec(&rules)
			if got := planet.getMaxPopulation(&rules, player, tt.hab); got != tt.want {
				t.Errorf("getMaxPopulation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_computePlanetSpec(t *testing.T) {
	player, planet := newTestPlayerPlanet()
	planet.Starbase = testSpaceStation(player, planet)

	player.Race.Spec.InnateScanner = true
	player.Race.Spec.InnateMinesFactor = 0.1
	player.Race.Spec.InnateScannerFactor = 0.1
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
	type fields struct {
		habDropoff Hab
		minHab     int
		maxHab     int
	}
	tests := []struct {
		name   string
		fields fields
		planet *Planet
		rng    rng
		want   *Planet
	}{
		{
			name:   "normal w/ all 0 rng; shouldn't reset production queue",
			fields: fields{rules.HabDropoffRange, rules.MinHab, rules.MaxHab},
			planet: NewPlanet().WithOrders(PlanetOrders{
				ProductionQueue: []ProductionQueueItem{
					{Type: QueueItemTypeFactory, Quantity: 1, Allocated: Cost{0, 0, 2, 6}},
					{Type: QueueItemTypeAutoDefenses, Quantity: 100},
					{Type: QueueItemTypeAutoFactories, Quantity: 10},
				},
			}),
			rng: newIntRandom(),
			want: &Planet{
				MapObject:            MapObject{Type: MapObjectTypePlanet, PlayerNum: Unowned},
				Dirty:                true,
				Hab:                  Hab{1, 1, 1},
				BaseHab:              Hab{1, 1, 1},
				MineralConcentration: Mineral{1, 1, 1},
				MineYears:            Mineral{},
				PlanetOrders: PlanetOrders{ProductionQueue: []ProductionQueueItem{
					{Type: QueueItemTypeFactory, Quantity: 1, Allocated: Cost{0, 0, 2, 6}},
					{Type: QueueItemTypeAutoDefenses, Quantity: 100},
					{Type: QueueItemTypeAutoFactories, Quantity: 10},
				}},
			},
		},
		{
			name:   "custom rules/rng seed",
			fields: fields{Hab{20, 20, 20}, 10, 90}, // 10 + rand[0,61) + rand[0,21)
			planet: NewPlanet(),
			rng:    newIntRandom(50, 43, 11, 2, 3, 5),
			want: &Planet{
				MapObject:            MapObject{Type: MapObjectTypePlanet, PlayerNum: Unowned},
				Dirty:                true,
				Hab:                  Hab{62, 56, 26},
				BaseHab:              Hab{62, 56, 26},
				MineralConcentration: Mineral{1, 1, 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.planet
			r := &rules
			r.HabDropoffRange = tt.fields.habDropoff
			r.MinHab = tt.fields.minHab
			r.MaxHab = tt.fields.maxHab
			r.random = tt.rng

			got.randomize(r)

			if !reflect.DeepEqual(got, tt.want) {
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

func TestPlanet_getMineralOutput(t *testing.T) {
	type fields struct {
		MineralConcentration Mineral
	}
	type args struct {
		numMines   int
		mineOutput int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Mineral
	}{
		{
			name:   "100 conc, 10 mines, 10 output",
			fields: fields{MineralConcentration: Mineral{100, 100, 100}},
			args:   args{numMines: 10, mineOutput: 10},
			want:   Mineral{10, 10, 10},
		},
		{
			name:   "100 conc, 10 mines, 8 output",
			fields: fields{MineralConcentration: Mineral{100, 100, 100}},
			args:   args{numMines: 10, mineOutput: 8},
			want:   Mineral{8, 8, 8},
		},
		{
			name:   "mixed conc, 100 mines, 10 output",
			fields: fields{MineralConcentration: Mineral{25, 45, 65}},
			args:   args{numMines: 100, mineOutput: 10},
			want:   Mineral{25, 45, 65},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Planet{
				MineralConcentration: tt.fields.MineralConcentration,
			}
			if got := p.getMineralOutput(tt.args.numMines, tt.args.mineOutput); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Planet.getMineralOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanetSpec_computeResourcesPerYear(t *testing.T) {
	type args struct {
		player          *Player
		numFacts        int
		productivePop   int
		installationPop int
	}
	tests := []struct {
		name string
		args args
		spec PlanetSpec
		want PlanetSpec
	}{
		{
			name: "Normal JoaT HW",
			args: args{
				player:          testPlayer(),
				numFacts:        10,
				productivePop:   25_000,
				installationPop: 25_000,
			},
			spec: PlanetSpec{
				MaxFactories:  100,
				MaxPopulation: 1_000_000,
			},
			want: PlanetSpec{
				MaxFactories:         25,
				MaxPossibleFactories: 1000,
				MaxPopulation:        1_000_000,
				ResourcesPerYear:     35,
			},
		},
		{
			name: "HE tiny planet overcapped on facts",
			args: args{
				player:          NewPlayer(1, NewRace().WithPRT(HE).WithSpec(&rules)).withSpec(&rules),
				numFacts:        99999,
				productivePop:   27_500,
				installationPop: 27_500,
			},
			spec: PlanetSpec{
				MaxFactories:  100,
				MaxPopulation: 27_500,
			},
			want: PlanetSpec{
				MaxFactories:         27,
				MaxPossibleFactories: 27,
				MaxPopulation:        27_500,
				ResourcesPerYear:     54,
			},
		},
		{
			name: "AR HW",
			args: args{
				player: NewPlayer(1, NewRace().WithPRT(AR).WithSpec(&rules)).
					WithTechLevels(TechLevel{1, 0, 0, 0, 0, 0}).
					withSpec(&rules),
				numFacts:        0,
				productivePop:   25_000,
				installationPop: 25_000,
			},
			spec: PlanetSpec{
				Habitability:  100,
				MaxPopulation: 1_000_000,
			},
			want: PlanetSpec{
				Habitability:     100,
				MaxPopulation:    1_000_000,
				ResourcesPerYear: 50,
			},
		},
		{
			name: "Crappy AR starter colony with lots of pop, En 10",
			args: args{
				player: NewPlayer(1, NewRace().WithPRT(AR).WithSpec(&rules)).
					WithTechLevels(TechLevel{10, 0, 0, 0, 0, 0}).withSpec(&rules), // makes calcs easier
				numFacts:        0,
				productivePop:   1_000_000,
				installationPop: 1_000_000,
			},
			spec: PlanetSpec{
				Habitability:  25, // min hab floor for AR
				MaxPopulation: 500_000,
			},
			want: PlanetSpec{
				Habitability:     25,
				MaxPopulation:    500_000,
				ResourcesPerYear: 250,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.spec.computeResourcesPerYear(tt.args.player, tt.args.numFacts, tt.args.productivePop, tt.args.installationPop)
			if !test.CompareAsJSON(t, tt.spec, tt.want) {
				// TODO: refactor after CompareAsJSON PR gets mergeed
				t.Errorf("computeResourcesPerYear resulted in bad specs")
			}
		})
	}
}
