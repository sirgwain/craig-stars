package cs

import (
	"reflect"
	"testing"
)

func Test_startingFleeter_getStartingHull(t *testing.T) {
	tests := []struct {
		name          string
		race          *Race
		techLevels    TechLevel
		startingFleet *StartingFleet
		wantHull      *TechHull
	}{
		{
			name: "JoaT Scout",
			race: NewRace().WithPRT(JoaT).WithSpec(&rules),
			techLevels: TechLevel{
				Energy:        3,
				Weapons:       3,
				Propulsion:    3,
				Construction:  3,
				Electronics:   3,
				Biotechnology: 3,
			},
			startingFleet: &StartingFleet{
				HullType: TechHullTypeScout,
				Type:     StartingFleetTypeScout,
			},
			wantHull: &Scout,
		},
		{
			name: "JoaT Destroyer",
			race: NewRace().WithPRT(JoaT).WithSpec(&rules),
			techLevels: TechLevel{
				Energy:        3,
				Weapons:       3,
				Propulsion:    3,
				Construction:  3,
				Electronics:   3,
				Biotechnology: 3,
			},
			startingFleet: &StartingFleet{
				HullType: TechHullTypeFighter,
				Type:     StartingFleetTypeFighter,
			},
			wantHull: &Destroyer,
		},
		{
			name: "JoaT Remote Miner",
			race: NewRace().WithPRT(JoaT),
			techLevels: TechLevel{
				Energy:        3,
				Weapons:       3,
				Propulsion:    3,
				Construction:  3,
				Electronics:   3,
				Biotechnology: 3,
			},
			startingFleet: &StartingFleet{
				HullType: TechHullTypeMiner,
				Type:     StartingFleetTypeMiner,
			},
			wantHull: &MiniMiner,
		},
		{
			name: "JoaT Colony Ship",
			race: NewRace().WithPRT(JoaT).WithSpec(&rules),
			techLevels: TechLevel{
				Energy:        3,
				Weapons:       3,
				Propulsion:    3,
				Construction:  3,
				Electronics:   3,
				Biotechnology: 3,
			},
			startingFleet: &StartingFleet{
				HullType: TechHullTypeColonizer,
				Type:     StartingFleetTypeColonizer,
			},
			wantHull: &ColonyShip,
		},
		{
			name: "JoaT Freighter",
			race: NewRace().WithPRT(JoaT).WithSpec(&rules),
			techLevels: TechLevel{
				Energy:        3,
				Weapons:       3,
				Propulsion:    3,
				Construction:  3,
				Electronics:   3,
				Biotechnology: 3,
			},
			startingFleet: &StartingFleet{
				HullType: TechHullTypeFreighter,
				Type:     StartingFleetTypeFighter,
			},
			wantHull: &MediumFreighter,
		},
		{
			name: "JoaT Con 4 Privateer",
			race: NewRace().WithPRT(JoaT).WithSpec(&rules),
			techLevels: TechLevel{
				Energy:        4,
				Weapons:       4,
				Propulsion:    4,
				Construction:  4,
				Electronics:   4,
				Biotechnology: 4,
			},
			startingFleet: &StartingFleet{
				HullType: TechHullTypeFreighter,
				Type:     StartingFleetTypeFighter,
			},
			wantHull: &Privateer,
		},
		{
			name:       "SS Freighter",
			race:       NewRace().WithPRT(SS).WithSpec(&rules),
			techLevels: TechLevel{Electronics: 5},
			startingFleet: &StartingFleet{
				HullType: TechHullTypeFreighter,
				Type:     StartingFleetTypeCloakedFreighter,
			},
			wantHull: &SmallFreighter,
		},
		{
			name: "WM Bomber - doesn't exist",
			race: NewRace().WithPRT(WM).WithSpec(&rules),
			techLevels: TechLevel{
				Energy:     1,
				Weapons:    6,
				Propulsion: 2,
			},
			startingFleet: &StartingFleet{
				HullType: TechHullTypeBomber,
				Type:     StartingFleetTypeBomber,
			},
			wantHull: nil,
		},
		{
			name: "WM Bomber - Con 3",
			race: NewRace().WithPRT(WM).WithSpec(&rules),
			techLevels: TechLevel{
				Energy:        3,
				Weapons:       3,
				Propulsion:    3,
				Construction:  3,
				Electronics:   3,
				Biotechnology: 3,
			},
			startingFleet: &StartingFleet{
				HullType: TechHullTypeBomber,
				Type:     StartingFleetTypeBomber,
			},
			wantHull: &MiniBomber,
		},
		{
			name: "HE Minicol",
			race: NewRace().WithPRT(HE).WithSpec(&rules),
			techLevels: TechLevel{
				Energy:        3,
				Weapons:       3,
				Propulsion:    3,
				Construction:  3,
				Electronics:   3,
				Biotechnology: 3,
			},
			startingFleet: &StartingFleet{
				HullType:         TechHullTypeColonizer,
				Type:             StartingFleetTypeColonizer,
				UsesCheapestHull: true,
			},
			wantHull: &MiniColonyShip,
		},
		{
			name: "ARM Midget Miner",
			race: NewRace().WithPRT(JoaT).WithLRT(ARM).WithSpec(&rules),
			techLevels: TechLevel{
				Energy:        3,
				Weapons:       3,
				Propulsion:    3,
				Construction:  3,
				Electronics:   3,
				Biotechnology: 3,
			},
			startingFleet: &StartingFleet{
				HullType:         TechHullTypeMiner,
				Type:             StartingFleetTypeMiner,
				UsesCheapestHull: true,
			},
			wantHull: &MidgetMiner,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(0, tt.race.WithSpec(&rules)).WithTechLevels(tt.techLevels)
			player.Name = tt.name
			sf := newStartingFleeter(&rules, player)
			if gotHull := sf.getStartingHull(tt.startingFleet); !reflect.DeepEqual(gotHull, tt.wantHull) {
				t.Errorf("startingFleeter.getStartingHull() returned starting hull \n%s, want \n%s", gotHull, tt.wantHull)
			}
		})
	}
}

func Test_startingFleeter_getStartingEngine(t *testing.T) {
	tests := []struct {
		name              string
		race              *Race
		techLevels        TechLevel
		startingFleetType StartingFleetType
		wantEngine        *TechEngine
	}{
		{
			name: "chooses DLL7 over mizer",
			race: NewRace().WithPRT(JoaT).WithLRT(IFE).WithSpec(&rules),
			techLevels: TechLevel{
				Energy:        4,
				Weapons:       4,
				Propulsion:    5,
				Construction:  4,
				Electronics:   4,
				Biotechnology: 4,
			},
			startingFleetType: StartingFleetTypeFighter,
			wantEngine:        &DaddyLongLegs7,
		},
		{
			name: "uses radram over AD8",
			race: NewRace().WithPRT(IT).WithLRT(IFE).WithLRT(CE).WithSpec(&rules),
			techLevels: TechLevel{
				Energy:        3,
				Weapons:       3,
				Propulsion:    7,
				Construction:  5,
				Electronics:   3,
				Biotechnology: 3,
			},
			startingFleetType: StartingFleetTypeFighter,
			wantEngine:        &RadiatingHydroRamScoop,
		},
		{
			name: "avoids radram for colony ships",
			race: NewRace().WithPRT(IT).WithLRT(IFE).WithLRT(CE).WithSpec(&rules),
			techLevels: TechLevel{
				Energy:        3,
				Weapons:       3,
				Propulsion:    7,
				Construction:  5,
				Electronics:   3,
				Biotechnology: 3,
			},
			startingFleetType: StartingFleetTypeColonizer,
			wantEngine:        &AlphaDrive8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(0, tt.race.WithSpec(&rules)).WithTechLevels(tt.techLevels)
			player.Name = tt.name
			hullName := Scout.Name
			sf := newStartingFleeter(&rules, player)
			if gotEngine := sf.getStartingEngine(tt.startingFleetType, hullName); !reflect.DeepEqual(gotEngine, tt.wantEngine) {
				t.Errorf("startingFleeter.getStartingEngine() returned starting engine \n%s, want \n%s", gotEngine.Name, tt.wantEngine)
			}
		})
	}
}

func Test_startingFleeter_createStartingDesign(t *testing.T) {
	tests := []struct {
		name          string
		race          *Race
		techLevels    TechLevel
		startingFleet *StartingFleet
		wantSlots     map[string]int
	}{
		{
			name: "Standard scout",
			race: NewRace().WithPRT(JoaT).WithSpec(&rules),
			techLevels: TechLevel{
				Energy:        3,
				Weapons:       3,
				Propulsion:    3,
				Construction:  3,
				Electronics:   3,
				Biotechnology: 3,
			},
			startingFleet: &StartingFleet{
				HullType: TechHullTypeScout,
				Type:     StartingFleetTypeScout,
			},
			wantSlots: map[string]int{
				LongHump6.Name:    1,
				RhinoScanner.Name: 1,
				FuelTank.Name:     1,
			},
		},
		{
			name: "JoaT ARM Midget Miner",
			race: NewRace().WithPRT(JoaT).WithLRT(IFE).WithLRT(ARM).WithSpec(&rules),
			techLevels: TechLevel{
				Energy:        4,
				Weapons:       4,
				Propulsion:    5,
				Construction:  4,
				Electronics:   4,
				Biotechnology: 4,
			},
			startingFleet: &StartingFleet{
				HullType:         TechHullTypeMiner,
				Type:             StartingFleetTypeMiner,
				UsesCheapestHull: true,
			},
			wantSlots: map[string]int{
				DaddyLongLegs7.Name: 1,
				RoboMiner.Name:      2,
			},
		},
		{
			name: "JoaT Prop 6 Con 4 Privateer",
			race: NewRace().WithPRT(JoaT).WithLRT(IFE).WithLRT(CE).WithSpec(&rules),
			techLevels: TechLevel{
				Energy:        4,
				Weapons:       4,
				Propulsion:    6,
				Construction:  4,
				Electronics:   4,
				Biotechnology: 4,
			},
			startingFleet: &StartingFleet{
				HullType: TechHullTypeFreighter,
				Type:     StartingFleetTypeFighter,
			},
			wantSlots: map[string]int{
				RadiatingHydroRamScoop.Name: 1,
				XRayLaser.Name:              1,
				AlphaTorpedo.Name:           1,
				MoleScanner.Name:            1,
				Carbonic.Name:               2,
			},
		},
		{
			name: "WM Bomber - Tech 3",
			race: NewRace().WithPRT(WM).WithLRT(IFE).WithSpec(&rules),
			techLevels: TechLevel{
				Energy:        3,
				Weapons:       6,
				Propulsion:    3,
				Construction:  3,
				Electronics:   3,
				Biotechnology: 3,
			},
			startingFleet: &StartingFleet{
				HullType: TechHullTypeBomber,
				Type:     StartingFleetTypeBomber,
			},
			wantSlots: map[string]int{
				LongHump6.Name:    1,
				BlackCatBomb.Name: 2,
			},
		},
		{
			name: "WM Bomber doesn't exist",
			race: NewRace().WithPRT(WM).WithLRT(IFE).WithSpec(&rules),
			techLevels: TechLevel{
				Energy:     1,
				Weapons:    6,
				Propulsion: 2,
			},
			startingFleet: &StartingFleet{
				HullType: TechHullTypeBomber,
				Type:     StartingFleetTypeBomber,
			},
			wantSlots: map[string]int{},
		},
		{
			name: "HE Minicol - Tech 3",
			race: NewRace().WithPRT(HE).WithSpec(&rules),
			techLevels: TechLevel{
				Energy:        3,
				Weapons:       3,
				Propulsion:    3,
				Construction:  3,
				Electronics:   3,
				Biotechnology: 3,
			},
			startingFleet: &StartingFleet{
				HullType:         TechHullTypeColonizer,
				Type:             StartingFleetTypeColonizer,
				UsesCheapestHull: true,
			},
			wantSlots: map[string]int{
				SettlersDelight.Name:    1,
				ColonizationModule.Name: 1,
			},
		},
		{
			name: "SD Minelayer",
			race: NewRace().WithPRT(SD).WithSpec(&rules),
			techLevels: TechLevel{
				Propulsion:    2,
				Biotechnology: 2,
			},
			startingFleet: &StartingFleet{
				HullType: TechHullTypeMineLayer,
				Type:     StartingFleetTypeMineLayer,
			},
			wantSlots: map[string]int{
				QuickJump5.Name:      1,
				MineDispenser40.Name: 4,
				BatScanner.Name:      1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(0, tt.race).WithTechLevels(tt.techLevels)
			player.Name = tt.name
			sf := newStartingFleeter(&rules, player)
			got, err := sf.createStartingDesign(tt.startingFleet, 1)
			if err != nil {
				t.Errorf("startingFleeter.createStartingDesign() returned error \n%v", err)
			}

			if (got == nil) != (len(tt.wantSlots) == 0) {
				t.Errorf("startingFleeter.createStartingDesign() returned nil design; wanted slots \n%v", tt.wantSlots)
			}

			if got == nil {
				return
			}

			tallyMap := map[string]int{}
			for _, slot := range got.Slots {
				tallyMap[slot.HullComponent] += slot.Quantity
			}

			if !reflect.DeepEqual(tallyMap, tt.wantSlots) {
				t.Errorf("ship design from startingFleeter.createStartingDesign() had parts \n%+v, want \n%+v", tallyMap, tt.wantSlots)
			}
		})
	}
}
