//go:build !wasi && !wasm

package cs

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func newTestPlayerPlanet() (player *Player, planet *Planet) {
	player = NewPlayer(1, NewRace()).WithNum(1)
	player.Race.Spec = ComputeRaceSpec(&player.Race, &rules)
	planet = &Planet{}
	planet.PlayerNum = player.Num
	planet.BaseHab = player.Race.HabCenter()
	planet.Hab = planet.BaseHab

	player.Spec = ComputePlayerSpec(player, &rules, []*Planet{planet})
	planet.Spec = ComputePlanetSpec(&rules, player, planet)

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
				design: NewShipDesign(player.Num, 1).
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
				design: NewShipDesign(player.Num, 1).
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
		{
			name:   "empty planet",
			fields: fields{Hab{50, 50, 50}, 0},
			args:   args{NewPlayer(1, NewRace()), 1_000_000},
			want:   0,
		},
		{
			name:   "low cap, full growth rate",
			fields: fields{Hab{50, 50, 50}, 100_000},
			args:   args{NewPlayer(1, NewRace()), 1_200_000},
			want:   10_000,
		},
		{
			name:   "fractional pop doesn't add growth",
			fields: fields{Hab{50, 50, 50}, 111_199},
			args:   args{NewPlayer(1, NewRace().WithGrowthRate(9)), 1_200_000},
			want:   9999, // would be 10K had the last 100 fractional pop counted for growth
		},
		{
			name:   "50% cap; growth slows down",
			fields: fields{Hab{50, 50, 50}, 600_000},
			args:   args{NewPlayer(1, NewRace()), 1_200_000},
			want:   26_666, // 0.444x growth multi
		},
		{
			name:   "98% capacity, next to no growth",
			fields: fields{Hab{50, 50, 50}, 1_180_000},
			args:   args{NewPlayer(1, NewRace()), 1_200_000}, // 98.3% cap
			want:   58,                                       // 0.0004938x growth multi
		},
		{
			name:   "slightly hostile planet; low deaths override overpop calcs",
			fields: fields{Hab{10, 15, 15}, 2500},
			args:   args{NewPlayer(1, NewRace()), 100},
			want:   -13,
		},
		{
			name:   "hostile planet; high deaths",
			fields: fields{Hab{0, 0, 0}, 100_000},
			args:   args{NewPlayer(1, NewRace()), 0},
			want:   -4500,
		},
		{
			name:   "overcap planet; lose pop",
			fields: fields{Hab{50, 50, 50}, 2_400_000},
			args:   args{NewPlayer(1, NewRace()), 1_200_000},
			want:   -96_000,
		},
		{
			name:   "5x cap; pop loss capped",
			fields: fields{Hab{50, 50, 50}, 6_000_000},
			args:   args{NewPlayer(1, NewRace()), 1_200_000},
			want:   -720_000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPlanet().WithHab(tt.fields.Hab).WithPopulation(tt.fields.Population)
			// If at default, set growth rate to 10% for easier math
			if tt.args.player.Race.GrowthRate == 15 {
				tt.args.player.Race.GrowthRate = 10
			}
			tt.args.player.Race.Spec = ComputeRaceSpec(&tt.args.player.Race, &rules)
			if got := p.GetGrowthAmount(tt.args.player, tt.args.maxPopulation, rules.PopulationOvercrowdDieoffRate, rules.PopulationOvercrowdDieoffRateMax); got != tt.want {
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
	tests := []struct {
		name   string
		planet *Planet
		want   *Planet
	}{
		{
			name: "reduces mineral conc",
			// 1.5M / 100 / 100 = 150 mine-years to reduce
			planet: NewPlanet().
				WithMineralConcentration(Mineral{100, 100, 100}).
				WithMineYears(Mineral{151, 151, 151}),
			want: NewPlanet().
				WithMineralConcentration(Mineral{99, 99, 99}).
				WithMineYears(Mineral{1, 1, 1}),
		},
		{
			name: "Homeworld can go below 30 conc",
			// 1.5M / 30 / 30 = 1,666 mine-years to reduce
			planet: NewPlanet().WithHomeworld(true).
				WithMineralConcentration(Mineral{30, 30, 30}).
				WithMineYears(Mineral{1667, 1667, 1667}),
			want: NewPlanet().WithHomeworld(true).
				WithMineralConcentration(Mineral{29, 29, 29}).
				WithMineYears(Mineral{1, 1, 1}),
		},
		{
			name: "Cannot lower below 1",
			// 1.5M / 2 / 2 = 375,000 mine-years to reduce
			planet: NewPlanet().
				WithMineralConcentration(Mineral{2, 2, 2}).
				WithMineYears(Mineral{5e7, 5e7, 5e7}),
			want: NewPlanet().
				WithMineralConcentration(Mineral{1, 1, 1}).
				WithMineYears(Mineral{0, 0, 0}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.planet.reduceMineralConcentration(&rules)
			test.CompareAsJSON(t, tt.planet, tt.want)
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
	planet.Spec = ComputePlanetSpec(&rules, player, planet)

	assert.Equal(t, planet.Spec.ScanRange, 82)
	assert.Equal(t, planet.Spec.ScanRangePen, 0)

	// now use a death start
	planet.Starbase = testDeathStar(player, planet)
	planet.Spec = ComputePlanetSpec(&rules, player, planet)

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
			r := NewRules()
			r.HabDropoffRange = tt.fields.habDropoff
			r.MinHab = tt.fields.minHab
			r.MaxHab = tt.fields.maxHab
			r.random = tt.rng

			got.randomize(&r, false)

			test.CompareAsJSON(t, got, tt.want)

		})
	}
}

func TestPlanet_grow(t *testing.T) {
	type fields struct {
		hab         Hab
		population  int
		turnsToGrow int
	}
	tests := []struct {
		name   string
		fields fields
		race   *Race
		want   int
	}{
		{
			name:   "standard humanoid starter world",
			fields: fields{hab: Hab{50, 50, 50}, population: 25000, turnsToGrow: 1},
			race:   NewRace().WithSpec(&rules),
			want:   28_750,
		},
		{
			name:   "partially full world; reduced growth",
			fields: fields{hab: Hab{50, 50, 50}, population: 500_000, turnsToGrow: 1},
			race:   NewRace().WithSpec(&rules),
			want:   545_370, // 60.4% GR multi due to crowding
		},
		{
			name:   "pas de population, pas de croissance",
			fields: fields{hab: Hab{50, 50, 50}, population: 0, turnsToGrow: 1},
			race:   NewRace().WithSpec(&rules),
			want:   0,
		},
		// TODO: Check in OG game if 0% worlds actually grow pop or not
		/* {
			name:   "0% value world, don't make colonists",
			fields: fields{hab: Hab{15, 15, 15}, population: 100, turnsToGrow: 1},
			race:   NewRace().WithSpec(&rules),
			want:   100,
		}, */
		{
			name:   "40% value world, 20% HE, 3 years",
			fields: fields{hab: Hab{35, 21, 45}, population: 600, turnsToGrow: 3},
			race:   NewRace().WithPRT(HE).WithGrowthRate(20).WithSpec(&rules),
			want:   904, // 600 --> 696 --> 792 --> 904
		},
		{
			name:   "hostile world",
			fields: fields{hab: Hab{1, 1, 1}, population: 25000, turnsToGrow: 1},
			race:   NewRace().WithSpec(&rules),
			want:   23_950, // -42% value = -4.2% pop per year
		},
		{
			name:   "hostile world, pop rounding",
			fields: fields{hab: Hab{1, 1, 1}, population: 200, turnsToGrow: 1},
			race:   NewRace().WithSpec(&rules),
			want:   192,
		},
		{
			name:   "hostile world, pop floor",
			fields: fields{hab: Hab{1, 1, 1}, population: 100, turnsToGrow: 1},
			race:   NewRace().WithSpec(&rules),
			want:   100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(0, tt.race).WithNum(1)
			planet := NewPlanet().WithPlayerNum(player.Num).
				WithHab(tt.fields.hab).WithPopulation(tt.fields.population)
			for range tt.fields.turnsToGrow {
				planet.Spec = ComputePlanetSpec(&rules, player, planet) // only really needed for growth amount
				planet.grow(player)
			}

			if exactPop := planet.exactPopulation(); exactPop != tt.want {
				t.Errorf("planet.grow() gave %v pop, want %v", exactPop, tt.want)
			}

		})
	}
}

func TestPlanet_getMineralOutput(t *testing.T) {
	tests := []struct {
		name       string
		planet     *Planet
		numMines   int
		mineOutput int
		want       Mineral
	}{
		// TODO: Change tests 2 and 4 after frac mineral output bug is fixed
		{
			name:       "whole number outputs",
			planet:     NewPlanet().WithMineralConcentration(Mineral{100, 100, 100}),
			numMines:   10,
			mineOutput: 10,
			want:       Mineral{10, 10, 10},
		},
		{
			name:       "mixed conc; truncates",
			planet:     NewPlanet().WithMineralConcentration(Mineral{25, 45, 65}),
			numMines:   22,
			mineOutput: 10,
			want:       Mineral{5, 9, 14}, // I/B get truncated
		},
		{
			name:       "Unowned homeworld; no conc floor",
			planet:     NewPlanet().WithMineralConcentration(Mineral{1, 1, 1}).WithHomeworld(true),
			numMines:   100,
			mineOutput: 10,
			want:       Mineral{1, 1, 1},
		},
		{
			name:       "Owned homeworld min conc floor",
			planet:     NewPlanet().WithMineralConcentration(Mineral{1, 1, 1}).WithHomeworld(true).WithPlayerNum(1),
			numMines:   22, // 6.6 minerals
			mineOutput: 10,
			want:       Mineral{6, 6, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.planet.getMineralOutput(&rules, tt.numMines, tt.mineOutput); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Planet.getMineralOutput() returned output \n%+v, want \n%+v", got, tt.want)
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
			name: "Negative Hab AR planet",
			args: args{
				player: NewPlayer(1, NewRace().WithPRT(AR).WithSpec(&rules)).
					WithTechLevels(TechLevel{10, 0, 0, 0, 0, 0}).withSpec(&rules), // makes calcs easier
				numFacts:        0,
				productivePop:   1_000_000,
				installationPop: 1_000_000,
			},
			spec: PlanetSpec{
				Habitability:  -45, // floored to 25 for AR
				MaxPopulation: 500_000,
			},
			want: PlanetSpec{
				Habitability:     -45,
				MaxPopulation:    500_000,
				ResourcesPerYear: 250,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.spec.ComputeResourcesPerYear(tt.args.player, tt.args.numFacts, tt.args.productivePop, tt.args.installationPop)
			test.CompareAsJSON(t, tt.spec, tt.want)
		})
	}
}
