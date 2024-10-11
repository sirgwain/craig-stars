package cs

import (
	"testing"

	"github.com/sirgwain/craig-stars/test"
)

func TestRace_GetPlanetHabitability(t *testing.T) {
	type args struct {
		hab  Hab
		race *Race
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Perfect Hab", args{Hab{50, 50, 50}, NewRace()}, 100},
		{"Terrible Hab", args{Hab{0, 0, 0}, &Race{HabLow: Hab{99, 99, 99}, HabHigh: Hab{100, 100, 100}}}, -45},
		{"1% away", args{Hab{48, 50, 50}, NewRace()}, 99},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.race.GetPlanetHabitability(tt.args.hab); got != tt.want {
				t.Errorf("Race.GetPlanetHabitability() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRace_ComputeRacePoints(t *testing.T) {
	startingPoints := 1650 // defined in rules

	immuneInsectoids := Insectoids()
	tests := []struct {
		name string
		race Race
		want int
	}{
		{"Humanoids", Humanoids(), 25},
		{"all immune", *NewRace().withImmuneGrav(true).withImmuneTemp(true).withImmuneRad(true), -3900},
		{"Rabbitoids", Rabbitoids(), 32},
		{"Insectoids", Insectoids(), 43},
		{"All Immune Insectoid", *immuneInsectoids.withImmuneRad(true).withImmuneTemp(true), -2112},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.race.ComputeRacePoints(startingPoints); got != tt.want {
				t.Errorf("Race.ComputeRacePoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_computeRaceSpec(t *testing.T) {
	tests := []struct {
		name string
		race *Race
		want RaceSpec
	}{
		{
			name: "humanoids w/arm",
			race: NewRace().WithLRT(ARM),
			want: RaceSpec{
				Costs: map[QueueItemType]Cost{
					QueueItemTypeAutoDefenses:       {5, 5, 5, 15},
					QueueItemTypeAutoFactories:      {Germanium: 4, Resources: 10},
					QueueItemTypeAutoMaxTerraform:   {Resources: 100},
					QueueItemTypeAutoMineralAlchemy: {Resources: 100},
					QueueItemTypeAutoMineralPacket: {
						Resources: 10,
						Ironium:   44,
						Boranium:  44,
						Germanium: 44,
					},
					QueueItemTypeAutoMines:              {Resources: 5},
					QueueItemTypeAutoMinTerraform:       {Resources: 100},
					QueueItemTypeBoraniumMineralPacket:  {Resources: 10, Boranium: 110},
					QueueItemTypeDefenses:               {5, 5, 5, 15},
					QueueItemTypeFactory:                {Germanium: 4, Resources: 10},
					QueueItemTypeGenesisDevice:          rules.MysteryTraderRules.GenesisDeviceCost,
					QueueItemTypeGermaniumMineralPacket: {Resources: 10, Germanium: 110},
					QueueItemTypeIroniumMineralPacket:   {Resources: 10, Ironium: 110},
					QueueItemTypeMine:                   {Resources: 5},
					QueueItemTypeMineralAlchemy:         {Resources: 100},
					QueueItemTypeMixedMineralPacket: {
						Resources: 10,
						Ironium:   44,
						Boranium:  44,
						Germanium: 44,
					},
					QueueItemTypePlanetaryScanner:     rules.PlanetaryScannerCost,
					QueueItemTypeTerraformEnvironment: {Resources: 100},
				},
				MiniaturizationSpec: MiniaturizationSpec{
					MiniaturizationMax:      0.75,
					MiniaturizationPerLevel: 0.04,
					NewTechCostFactor:       1,
				},
				ScannerSpec: ScannerSpec{
					BuiltInScannerMultiplier: 20,
					ScanRangeFactor:          1,
				},
				StartingPlanets: []StartingPlanet{{
					Defenses:           10,
					Factories:          10,
					Homeworld:          true,
					Mines:              10,
					Population:         25000,
					StarbaseDesignName: "Starbase",
					StarbaseHull:       SpaceStation.Name,
					StartingFleets: []StartingFleet{
						{"Long Range Scout", StartingFleetHullScout, 0, ShipDesignPurposeScout},
						{"Santa Maria", StartingFleetHullColonyShip, 0, ShipDesignPurposeColonizer},
						{"Teamster", StartingFleetHullMediumFreighter, 0, ShipDesignPurposeFreighter},
						{"Cotton Picker", StartingFleetHullMiniMiner, 0, ShipDesignPurposeMiner},
						{"Armed Probe", StartingFleetHullScout, 1, ShipDesignPurposeFighterScout},
						{"Stalwart Defender", StartingFleetHullDestroyer, 0, ShipDesignPurposeFighter},
						{"Potato Bug", StartingFleetHullMidgetMiner, 0, ShipDesignPurposeMiner},
						{"Potato Bug", StartingFleetHullMidgetMiner, 0, ShipDesignPurposeMiner},
					},
				}},
				ArmorStrengthFactor:            1,
				CanBuildDefenses:               true,
				EngineReliableSpeed:            10,
				GrowthFactor:                   1,
				HabCenter:                      Hab{50, 50, 50},
				InnatePopulationFactor:         1,
				InvasionAttackBonus:            1.1,
				InvasionDefendBonus:            1,
				MaxPopulationOffset:            .2,
				MineFieldBaseDecayRate:         .02,
				MineFieldDetonateDecayRate:     .25,
				MineFieldMaxDecayRate:          .5,
				MineFieldMinDecayFactor:        1,
				MineFieldPlanetDecayRate:       .04,
				MineralsPerMixedMineralPacket:  40,
				MineralsPerSingleMineralPacket: 100,
				PacketDecayFactor:              1,
				PacketMineralCostFactor:        1.1, // 10% overhead
				PacketPermaTerraformSizeUnit:   100,
				PacketReceiverFactor:           1,
				PacketResourceCost:             10,
				RepairFactor:                   1,
				ResearchFactor:                 1,
				ShieldStrengthFactor:           1,
				ShipsVanishInVoid:              true,
				StarbaseCostFactor:             1,
				StarbaseRepairFactor:           1,
				StartingPopulationFactor:       1,
				StartingTechLevels:             TechLevel{3, 3, 3, 3, 3, 3},
				TechsCostExtraLevel:            4,
			},
		},
		{
			name: "humanoid",
			race: NewRace(),
			want: RaceSpec{
				Costs: map[QueueItemType]Cost{
					QueueItemTypeAutoDefenses:       {5, 5, 5, 15},
					QueueItemTypeAutoFactories:      {Germanium: 4, Resources: 10},
					QueueItemTypeAutoMaxTerraform:   {Resources: 100},
					QueueItemTypeAutoMineralAlchemy: {Resources: 100},
					QueueItemTypeAutoMineralPacket: {
						Resources: 10,
						Ironium:   44,
						Boranium:  44,
						Germanium: 44,
					},
					QueueItemTypeAutoMines:              {Resources: 5},
					QueueItemTypeAutoMinTerraform:       {Resources: 100},
					QueueItemTypeBoraniumMineralPacket:  {Resources: 10, Boranium: 110},
					QueueItemTypeDefenses:               {5, 5, 5, 15},
					QueueItemTypeFactory:                {Germanium: 4, Resources: 10},
					QueueItemTypeGenesisDevice:          rules.MysteryTraderRules.GenesisDeviceCost,
					QueueItemTypeGermaniumMineralPacket: {Resources: 10, Germanium: 110},
					QueueItemTypeIroniumMineralPacket:   {Resources: 10, Ironium: 110},
					QueueItemTypeMine:                   {Resources: 5},
					QueueItemTypeMineralAlchemy:         {Resources: 100},
					QueueItemTypeMixedMineralPacket: {
						Resources: 10,
						Ironium:   44,
						Boranium:  44,
						Germanium: 44,
					},
					QueueItemTypePlanetaryScanner:     rules.PlanetaryScannerCost,
					QueueItemTypeTerraformEnvironment: {Resources: 100},
				},
				MiniaturizationSpec: MiniaturizationSpec{
					MiniaturizationMax:      0.75,
					MiniaturizationPerLevel: 0.04,
					NewTechCostFactor:       1,
				},
				ScannerSpec: ScannerSpec{
					BuiltInScannerMultiplier: 20,
					ScanRangeFactor:          1,
				},
				StartingPlanets: []StartingPlanet{{
					Defenses:           10,
					Factories:          10,
					Homeworld:          true,
					Mines:              10,
					Population:         25000,
					StarbaseDesignName: "Starbase",
					StarbaseHull:       SpaceStation.Name,
					StartingFleets: []StartingFleet{
						{"Long Range Scout", StartingFleetHullScout, 0, ShipDesignPurposeScout},
						{"Santa Maria", StartingFleetHullColonyShip, 0, ShipDesignPurposeColonizer},
						{"Teamster", StartingFleetHullMediumFreighter, 0, ShipDesignPurposeFreighter},
						{"Cotton Picker", StartingFleetHullMiniMiner, 0, ShipDesignPurposeMiner},
						{"Armed Probe", StartingFleetHullScout, 1, ShipDesignPurposeFighterScout},
						{"Stalwart Defender", StartingFleetHullDestroyer, 0, ShipDesignPurposeFighter},
					},
				}},
				ArmorStrengthFactor:            1,
				CanBuildDefenses:               true,
				EngineReliableSpeed:            10,
				GrowthFactor:                   1,
				HabCenter:                      Hab{50, 50, 50},
				InnatePopulationFactor:         1,
				InvasionAttackBonus:            1.1,
				InvasionDefendBonus:            1,
				MaxPopulationOffset:            .2,
				MineFieldBaseDecayRate:         .02,
				MineFieldDetonateDecayRate:     .25,
				MineFieldMaxDecayRate:          .5,
				MineFieldMinDecayFactor:        1,
				MineFieldPlanetDecayRate:       .04,
				MineralsPerMixedMineralPacket:  40,
				MineralsPerSingleMineralPacket: 100,
				PacketDecayFactor:              1,
				PacketMineralCostFactor:        1.1, // 10% overhead
				PacketPermaTerraformSizeUnit:   100,
				PacketReceiverFactor:           1,
				PacketResourceCost:             10,
				RepairFactor:                   1,
				ResearchFactor:                 1,
				ShieldStrengthFactor:           1,
				ShipsVanishInVoid:              true,
				StarbaseCostFactor:             1,
				StarbaseRepairFactor:           1,
				StartingPopulationFactor:       1,
				StartingTechLevels:             TechLevel{3, 3, 3, 3, 3, 3},
				TechsCostExtraLevel:            4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := computeRaceSpec(tt.race, &rules); !test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("computeRaceSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}
