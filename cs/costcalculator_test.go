package cs

import (
	"testing"
)

func TestCostCalculator_StarbaseUpgradeCost(t *testing.T) {
	p := NewCostCalculator()
	type args struct {
		techLevels          TechLevel
		miniaturizationSpec MiniaturizationSpec
		techCostOffset      TechCostOffset
		oldDesignHull       string
		newDesignHull       string
		oldDesignSlots      []ShipDesignSlot
		newDesignSlots      []ShipDesignSlot
		starbaseCostFactor  float64
	}
	tests := []struct {
		name     string
		args     args
		want     Cost
		wanterr  bool
	}{
		{
			name: "Invalid station",
			args: args{
				techLevels:          TechLevel{0, 0, 0, 0, 0, 0},
				miniaturizationSpec: MiniaturizationSpec{1.0, 0.75, 0.04},
				techCostOffset:      TechCostOffset{0, 0, 0, 0, 0, 0, 0},
				oldDesignHull:       "NotAStation",
				newDesignHull:       SpaceStation.Name,
				oldDesignSlots:      []ShipDesignSlot{},
				newDesignSlots:      []ShipDesignSlot{},
				starbaseCostFactor:  1,
			},
			want: Cost{}, wanterr: true,
		},
		{
			name: "Invalid parts",
			args: args{
				techLevels:          TechLevel{0, 0, 0, 0, 0, 0},
				miniaturizationSpec: MiniaturizationSpec{1.0, 0.75, 0.04},
				techCostOffset:      TechCostOffset{0, 0, 0, 0, 0, 0, 0},
				oldDesignHull:       SpaceStation.Name,
				newDesignHull:       SpaceStation.Name,
				oldDesignSlots:      []ShipDesignSlot{
					{HullComponent: "pew pew laser gun", HullSlotIndex: 2, Quantity: 1},
				},
				newDesignSlots:      []ShipDesignSlot{},
				starbaseCostFactor:  1,
			},
			want: Cost{}, wanterr: true,
		},
		{
			name: "Identical Bases",
			args: args{
				techLevels:          TechLevel{0, 0, 0, 0, 0, 0},
				miniaturizationSpec: MiniaturizationSpec{1.0, 0.75, 0.04},
				techCostOffset:      TechCostOffset{0, 0, 0, 0, 0, 0, 0},
				oldDesignHull:       SpaceStation.Name,
				newDesignHull:       SpaceStation.Name,
				oldDesignSlots:      []ShipDesignSlot{},
				newDesignSlots:      []ShipDesignSlot{},
				starbaseCostFactor:  1,
			},
			want: Cost{}, wanterr: false,
		},
		{
			name: "Items on former base not on latter",
			args: args{
				techLevels:          TechLevel{0, 13, 0, 0, 0, 0},
				miniaturizationSpec: MiniaturizationSpec{1.0, 0.75, 0.04},
				techCostOffset:      TechCostOffset{0, 0, 0, 0, 0, 0, 0},
				oldDesignHull:       SpaceStation.Name,
				newDesignHull:       SpaceStation.Name,
				oldDesignSlots:      []ShipDesignSlot{
					{HullComponent: Bludgeon.Name, HullSlotIndex: 2, Quantity: 1},
				},
				newDesignSlots:      []ShipDesignSlot{},
				starbaseCostFactor:  1,
			},
			want: Cost{}, wanterr: false,
		},
		{
			name: "Adding weapons",
			args: args{
				techLevels:          TechLevel{0, 20, 0, 0, 0, 0},
				miniaturizationSpec: MiniaturizationSpec{1.0, 0.75, 0.04},
				techCostOffset:      TechCostOffset{0, 0, 0, 0, 0, 0, 0},
				oldDesignHull:       SpaceStation.Name,
				newDesignHull:       SpaceStation.Name,
				oldDesignSlots:      []ShipDesignSlot{},
				newDesignSlots:      []ShipDesignSlot{
					{HullComponent: Disruptor.Name, HullSlotIndex: 2, Quantity: 1},
					{HullComponent: Disruptor.Name, HullSlotIndex: 4, Quantity: 1},
				},
				starbaseCostFactor:  1,
			},
			want: Cost{
				Ironium:   0,
				Boranium:  16,
				Germanium: 0,
				Resources: 20,
			}, wanterr: false,
		},
		{
			name: "Adding single orbital",
			args: args{
				techLevels:          TechLevel{0, 0, 5, 5, 0, 0},
				miniaturizationSpec: MiniaturizationSpec{1.0, 0.75, 0.04},
				techCostOffset:      TechCostOffset{0, 0, 0, 0, 0, 0, 0},
				oldDesignHull:       SpaceStation.Name,
				newDesignHull:       SpaceStation.Name,
				oldDesignSlots:      []ShipDesignSlot{},
				newDesignSlots: []ShipDesignSlot{
					{HullComponent: Stargate100_250.Name, HullSlotIndex: 1, Quantity: 1},
				},
				starbaseCostFactor:  1,
			},
			want: Cost{
				Ironium:   50,
				Boranium:  20,
				Germanium: 20,
				Resources: 200,
			}, wanterr: false,
		},
		{
			name: "Hull swap",
			args: args{
				techLevels:          TechLevel{0, 20, 0, 0, 0, 0},
				miniaturizationSpec: MiniaturizationSpec{1.0, 0.75, 0.04},
				techCostOffset:      TechCostOffset{0, 0, 0, 0, 0, 0, 0},
				oldDesignHull:       OrbitalFort.Name,
				newDesignHull:       SpaceStation.Name,
				oldDesignSlots:      []ShipDesignSlot{},
				newDesignSlots:      []ShipDesignSlot{},
				starbaseCostFactor:  1,
			},
			want: Cost{
				Ironium:   114,
				Boranium:  80,
				Germanium: 242,
				Resources: 580,
			}, wanterr: false,
		},
		{
			name: "Hull Swap + added components",
			args: args{
				techLevels:          TechLevel{0, 20, 0, 0, 0, 0},
				miniaturizationSpec: MiniaturizationSpec{1.0, 0.75, 0.04},
				techCostOffset:      TechCostOffset{0, 0, 0, 0, 0, 0, 0},
				oldDesignHull:       OrbitalFort.Name,
				newDesignHull:       SpaceStation.Name,
				oldDesignSlots:      []ShipDesignSlot{},
				newDesignSlots: []ShipDesignSlot{
					{HullComponent: Disruptor.Name, HullSlotIndex: 2, Quantity: 10},
					// 80B, 100R
				},
				starbaseCostFactor:  1,
			},
			want: Cost{
				Ironium:   114,
				Boranium:  160,
				Germanium: 242,
				Resources: 680,
			}, wanterr: false,
		},
		{
			name: "Min Price Floor - same category",
			args: args{
				techLevels:          TechLevel{0, 22, 0, 0, 0, 0},
				miniaturizationSpec: MiniaturizationSpec{1.0, 0.75, 0.04},
				techCostOffset:      TechCostOffset{0, 0, 0, 0, 0, 0, 0},
				oldDesignHull:       SpaceStation.Name,
				newDesignHull:       SpaceStation.Name,
				oldDesignSlots:      []ShipDesignSlot{
					{HullComponent: SyncroSapper.Name, HullSlotIndex: 2, Quantity: 16},
					{HullComponent: SyncroSapper.Name, HullSlotIndex: 4, Quantity: 16},
					// Some amount of G, much much resources
				},
				newDesignSlots: []ShipDesignSlot{
					{HullComponent: MegaDisruptor.Name, HullSlotIndex: 2, Quantity: 10},
					// 150B, 165R
					// Boranium should stay same since not being reduced
				},
				starbaseCostFactor:  1,
			},
			want: Cost{
				Ironium:   0,
				Boranium:  150,
				Germanium: 0,
				Resources: 33,
			}, wanterr: false,
		},
		{
			name: "Component Swap (different categories)",
			args: args{
				techLevels:          TechLevel{22, 22, 22, 10, 0, 0},
				miniaturizationSpec: MiniaturizationSpec{1.0, 0.75, 0.04},
				techCostOffset:      TechCostOffset{0, 0, 0, 0, 0, 0, 0},
				oldDesignHull:       SpaceStation.Name,
				newDesignHull:       SpaceStation.Name,
				oldDesignSlots: []ShipDesignSlot{
					{HullComponent: UpsilonTorpedo.Name, HullSlotIndex: 2, Quantity: 16},
					// 320I, 112B, 72G, 120R
				},
				newDesignSlots: []ShipDesignSlot{
					{HullComponent: SyncroSapper.Name, HullSlotIndex: 2, Quantity: 12},
					// 48G, 120R
				},
				starbaseCostFactor:  1,
			},
			want: Cost{
				Ironium:   0,
				Boranium:  0,
				Germanium: 15,
				Resources: 36,
			}, wanterr: false,
		},
		{
			name: "Component Swap (same categories)",
			args: args{
				techLevels:          TechLevel{22, 22, 22, 22, 21, 22},
				miniaturizationSpec: MiniaturizationSpec{1.0, 0.75, 0.04},
				techCostOffset:      TechCostOffset{0, 0, 0, 0, 0, 0, 0},
				oldDesignHull:       SpaceStation.Name,
				newDesignHull:       SpaceStation.Name,
				oldDesignSlots: []ShipDesignSlot{
					{HullComponent: SyncroSapper.Name, HullSlotIndex: 2, Quantity: 12},
					// 48G, 120R
				},
				newDesignSlots: []ShipDesignSlot{
					{HullComponent: MegaDisruptor.Name, HullSlotIndex: 2, Quantity: 10},
					// 150B, 165R
					{HullComponent: BattleNexus.Name, HullSlotIndex: 1, Quantity: 2},
					// 28G, 14R
				},
				starbaseCostFactor:  1,
			},
			want: Cost{
				Ironium:   0,
				Boranium:  150,
				Germanium: 9, // technically 8.4 but gets rounded up to 9
				Resources: 83,
			}, wanterr: false,
		},
		{
			name: "ISB Component Swap",
			args: args{
				techLevels:          TechLevel{22, 22, 22, 22, 21, 22},
				miniaturizationSpec: MiniaturizationSpec{1.0, 0.75, 0.04},
				techCostOffset:      TechCostOffset{0, 0, 0, 0, 0, 0, 0},
				oldDesignHull:       SpaceStation.Name,
				newDesignHull:       SpaceStation.Name,
				oldDesignSlots: []ShipDesignSlot{
					{HullComponent: SyncroSapper.Name, HullSlotIndex: 2, Quantity: 12},
					// 48G, 120R
				},
				newDesignSlots: []ShipDesignSlot{
					{HullComponent: MegaDisruptor.Name, HullSlotIndex: 2, Quantity: 10},
					// 150B, 165R
					{HullComponent: BattleNexus.Name, HullSlotIndex: 1, Quantity: 2},
					// 28G, 14R
				},
				starbaseCostFactor:  0.8,
			},
			want: Cost{
				Ironium:   0,
				Boranium:  120,
				Germanium: 7,
				Resources: 67,
			}, wanterr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(1, NewRace().WithSpec(&rules))
			player.TechLevels = tt.args.techLevels
			player.Race.Spec.MiniaturizationSpec = tt.args.miniaturizationSpec
			player.Race.Spec.StarbaseCostFactor = tt.args.starbaseCostFactor
			player.Race.Spec.TechCostOffset = tt.args.techCostOffset
			design := NewShipDesign(player, 1).
				WithHull(tt.args.oldDesignHull).
				WithSlots(tt.args.oldDesignSlots)
			newDesign := NewShipDesign(player, 1).
				WithHull(tt.args.newDesignHull).
				WithSlots(tt.args.newDesignSlots)
			got, err := p.StarbaseUpgradeCost(&rules, tt.args.techLevels, player.Race.Spec, design, newDesign)
			if (err != nil) != tt.wanterr {
				t.Errorf("costCalculate.StarbaseUpgradeCost() errored unexpectedly; err = %v", err)
			} else if got != tt.want {
				t.Errorf("costCalculate.StarbaseUpgradeCost() returned incorrect cost %v, want %v", got, tt.want)
			}
		})
	}
}
