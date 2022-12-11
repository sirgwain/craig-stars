package cs

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestPlayerPlanet() (player *Player, planet *Planet) {
	rules := NewRules()
	player = NewPlayer(1, NewRace())
	player.Race.Spec = computeRaceSpec(&player.Race, &rules)
	planet = &Planet{}
	planet.PlayerID = player.ID
	planet.PlayerNum = player.Num

	return player, planet
}

func TestPlanet_String(t *testing.T) {

	tests := []struct {
		name string
		p    *Planet
		want string
	}{
		{"MapObject String()", &Planet{MapObject: MapObject{GameID: 1, ID: 2, Num: 3, Name: "Bob's Revenge"}},
			"Planet GameID:     1, ID:     2, Num:   3 Bob's Revenge"},
		{"MapObject String()", &Planet{MapObject: MapObject{GameID: 12345, ID: 23456, Num: 120, Name: "Craig's Planet"}},
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

func TestPlanet_GetInnateMines(t *testing.T) {
	player := NewPlayer(1, &Race{Spec: RaceSpec{InnateMining: false}})
	planet := Planet{}
	planet.setPopulation(16000)

	if got := planet.innateMines(player); got != 0 {
		t.Errorf("Planet.GetInnateMines() = %v, want %v", got, 0)
	}

	// should get 40 mines for 16k pop when the player has innate mining
	player.Race.Spec.InnateMining = true
	if got := planet.innateMines(player); got != 40 {
		t.Errorf("Planet.GetInnateMines() = %v, want %v", got, 40)
	}

}

func TestPlanet_getGrowthAmount(t *testing.T) {
	rules := NewRules()
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
			if got := p.getGrowthAmount(tt.args.player, tt.args.maxPopulation); got != tt.want {
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

func TestPlanet_produce(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// build 1 mine
	planet.ProductionQueue = []ProductionQueueItem{{Type: QueueItemTypeMine, Quantity: 1}}
	planet.Cargo = Cargo{10, 20, 30, 2500}
	planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 100, MaxPossibleMines: 100}
	planet.Mines = 0

	// should build 1 mine, leaving empty queue
	planet.produce(player)
	assert.Equal(t, 1, planet.Mines)
	assert.Equal(t, 0, len(planet.ProductionQueue))
	assert.Equal(t, 1, len(player.Messages))

	// build 5 auto mines, leaving them in the queue
	planet.ProductionQueue = []ProductionQueueItem{{Type: QueueItemTypeAutoMines, Quantity: 5}}
	planet.Cargo = Cargo{10, 20, 30, 2500}
	planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 100, MaxMines: 100}
	planet.Mines = 0
	player.Messages = []PlayerMessage{}

	// should build 5 mine, leaving the auto build in the queu
	planet.produce(player)
	assert.Equal(t, 5, planet.Mines)
	assert.Equal(t, 1, len(planet.ProductionQueue))
	assert.Equal(t, QueueItemTypeAutoMines, planet.ProductionQueue[0].Type)
	assert.Equal(t, 5, planet.ProductionQueue[0].Quantity)
	assert.Equal(t, 1, len(player.Messages))
}

func TestPlanet_produce2(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// build 5 auto factories, leaving them in the queue
	planet.ProductionQueue = []ProductionQueueItem{{Type: QueueItemTypeAutoFactories, Quantity: 5}}
	planet.Cargo = Cargo{10, 20, 30, 2500}
	planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 100, MaxFactories: 100}
	planet.Factories = 0
	player.Messages = []PlayerMessage{}

	// should build 5 mine, leaving the auto build in the queu
	planet.produce(player)
	assert.Equal(t, 5, planet.Factories)
	assert.Equal(t, 1, len(planet.ProductionQueue))
	assert.Equal(t, QueueItemTypeAutoFactories, planet.ProductionQueue[0].Type)
	assert.Equal(t, 5, planet.ProductionQueue[0].Quantity)
	assert.Equal(t, 1, len(player.Messages))

}

func TestPlanet_produce3(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// build 2/5 auto factories, leaving them in the queue
	planet.ProductionQueue = []ProductionQueueItem{
		{Type: QueueItemTypeAutoFactories, Quantity: 5},
		{Type: QueueItemTypeAutoMines, Quantity: 5},
	}
	planet.Cargo = Cargo{0, 0, 8, 2500}
	planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 100, MaxFactories: 100, MaxMines: 100}
	planet.Factories = 0
	player.Messages = []PlayerMessage{}

	// should build 2 factories and 5 mines, leaving the auto builds in the queue
	planet.produce(player)
	assert.Equal(t, 2, planet.Factories)
	assert.Equal(t, 5, planet.Mines)
	assert.Equal(t, 2, len(planet.ProductionQueue))
	assert.Equal(t, QueueItemTypeAutoFactories, planet.ProductionQueue[0].Type)
	assert.Equal(t, 5, planet.ProductionQueue[0].Quantity)
	assert.Equal(t, QueueItemTypeAutoMines, planet.ProductionQueue[1].Type)
	assert.Equal(t, 5, planet.ProductionQueue[1].Quantity)
}

func TestPlanet_produce4(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// build 2/5 auto factories and one partial mine
	planet.ProductionQueue = []ProductionQueueItem{
		{Type: QueueItemTypeAutoFactories, Quantity: 5},
		{Type: QueueItemTypeAutoMines, Quantity: 5},
	}
	planet.Cargo = Cargo{0, 0, 8, 2500}
	planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 10*2 + 8, MaxFactories: 100, MaxMines: 100}
	planet.Factories = 0
	player.Messages = []PlayerMessage{}

	// should build 2 factories, 1 mine and one partial mine, leaving the auto builds in the queue
	planet.produce(player)
	assert.Equal(t, 2, planet.Factories)
	assert.Equal(t, 1, planet.Mines)
	assert.Equal(t, 3, len(planet.ProductionQueue))
	assert.Equal(t, QueueItemTypeMine, planet.ProductionQueue[0].Type)
	assert.Equal(t, 1, planet.ProductionQueue[0].Quantity)
	assert.Equal(t, Cost{0, 0, 0, 3}, planet.ProductionQueue[0].Allocated)
	assert.Equal(t, QueueItemTypeAutoFactories, planet.ProductionQueue[1].Type)
	assert.Equal(t, 5, planet.ProductionQueue[1].Quantity)
	assert.Equal(t, QueueItemTypeAutoMines, planet.ProductionQueue[2].Type)
	assert.Equal(t, 5, planet.ProductionQueue[2].Quantity)

}

func TestPlanet_allocatePartialBuild(t *testing.T) {
	_, planet := newTestPlayerPlanet()

	type args struct {
		costPerItem Cost
		allocated   Cost
	}
	tests := []struct {
		name string
		args args
		want Cost
	}{
		{"Allocate Partial", args{costPerItem: Cost{10, 20, 30, 40}, allocated: Cost{5, 100, 100, 100}}, Cost{5, 10, 15, 20}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := planet.allocatePartialBuild(tt.args.costPerItem, tt.args.allocated); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Planet.allocatePartialBuild() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanet_reduceMineralConcentration(t *testing.T) {
	rules := NewRules()
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
