package cs

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_production_produce(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// build 1 mine
	planet.ProductionQueue = []ProductionQueueItem{{Type: QueueItemTypeMine, Quantity: 1}}
	planet.Cargo = Cargo{10, 20, 30, 2500}
	planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 100, MaxPossibleMines: 100}
	planet.Mines = 0

	// should build 1 mine, leaving empty queue
	producer := newProducer(planet, player)
	producer.produce()
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
	producer = newProducer(planet, player)
	producer.produce()
	assert.Equal(t, 5, planet.Mines)
	assert.Equal(t, 1, len(planet.ProductionQueue))
	assert.Equal(t, QueueItemTypeAutoMines, planet.ProductionQueue[0].Type)
	assert.Equal(t, 5, planet.ProductionQueue[0].Quantity)
	assert.Equal(t, 1, len(player.Messages))
}

func Test_production_produce2(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// build 5 auto factories, leaving them in the queue
	planet.ProductionQueue = []ProductionQueueItem{{Type: QueueItemTypeAutoFactories, Quantity: 5}}
	planet.Cargo = Cargo{10, 20, 30, 2500}
	planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 100, MaxFactories: 100}
	planet.Factories = 0
	player.Messages = []PlayerMessage{}

	// should build 5 mine, leaving the auto build in the queu
	producer := newProducer(planet, player)
	producer.produce()
	assert.Equal(t, 5, planet.Factories)
	assert.Equal(t, 1, len(planet.ProductionQueue))
	assert.Equal(t, QueueItemTypeAutoFactories, planet.ProductionQueue[0].Type)
	assert.Equal(t, 5, planet.ProductionQueue[0].Quantity)
	assert.Equal(t, 1, len(player.Messages))

}

func Test_production_produce3(t *testing.T) {
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
	producer := newProducer(planet, player)
	producer.produce()
	assert.Equal(t, 2, planet.Factories)
	assert.Equal(t, 5, planet.Mines)
	assert.Equal(t, 2, len(planet.ProductionQueue))
	assert.Equal(t, QueueItemTypeAutoFactories, planet.ProductionQueue[0].Type)
	assert.Equal(t, 5, planet.ProductionQueue[0].Quantity)
	assert.Equal(t, QueueItemTypeAutoMines, planet.ProductionQueue[1].Type)
	assert.Equal(t, 5, planet.ProductionQueue[1].Quantity)
}

func Test_production_produce4(t *testing.T) {
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
	producer := newProducer(planet, player)
	producer.produce()
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

func Test_production_produceTerraform(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// build 2/5 auto factories and one partial mine
	planet.ProductionQueue = []ProductionQueueItem{
		{Type: QueueItemTypeTerraformEnvironment, Quantity: 5},
	}
	planet.Cargo = Cargo{0, 0, 8, 2500}
	planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 1000}
	planet.BaseHab = Hab{40, 40, 40}
	planet.Hab = Hab{40, 40, 40}

	player.Spec.Terraform[TerraformHabTypeAll] = &TotalTerraform3
	player.Messages = []PlayerMessage{}

	// build 5 terraform steps, make planet better, log some messages
	producer := newProducer(planet, player)
	producer.produce()
	assert.Equal(t, Hab{42, 42, 41}, planet.Hab)
	assert.Equal(t, 0, len(planet.ProductionQueue))
	assert.Equal(t, 5, len(player.Messages))

}

func Test_production_allocatePartialBuild(t *testing.T) {
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

			production := production{planet: planet}

			if got := production.allocatePartialBuild(tt.args.costPerItem, tt.args.allocated); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Planet.allocatePartialBuild() = %v, want %v", got, tt.want)
			}
		})
	}
}
