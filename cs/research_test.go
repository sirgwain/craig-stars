package cs

import (
	"testing"
)

func Test_research_getNextResearchField(t *testing.T) {
	rules := NewRules()
	type args struct {
		player *Player
	}
	tests := []struct {
		name          string
		args          args
		wantNextField TechField
	}{
		{"Next Weapons", args{testPlayer().WithNextResearchField(NextResearchFieldWeapons)}, Weapons},
		{"Lowest", args{testPlayer().WithTechLevels(TechLevel{3, 3, 3, 3, 3, 2}).WithNextResearchField(NextResearchFieldLowestField)}, Biotechnology},
		{"Same", args{testPlayer().WithResearching(Construction).WithNextResearchField(NextResearchFieldSameField)}, Construction},
		{"At Max, do lowest", args{testPlayer().WithTechLevels(TechLevel{26, 2, 2, 2, 2, 1}).WithNextResearchField(NextResearchFieldEnergy)}, Biotechnology},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &research{
				rules: &rules,
			}
			if gotNextField := r.getNextResearchField(tt.args.player); gotNextField != tt.wantNextField {
				t.Errorf("research.getNextResearchField() = %v, want %v", gotNextField, tt.wantNextField)
			}
		})
	}
}

func Test_research_isAtMaxLevel(t *testing.T) {
	rules := NewRules()
	type args struct {
		player *Player
		field  TechField
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"no levels", args{testPlayer(), Energy}, false},
		{"max weapons", args{testPlayer().WithTechLevels(TechLevel{Weapons: 26}), Weapons}, true},
		{"max weapons, check energy", args{testPlayer().WithTechLevels(TechLevel{Weapons: 26}), Energy}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &research{
				rules: &rules,
			}
			if got := r.isAtMaxLevel(tt.args.player, tt.args.field); got != tt.want {
				t.Errorf("research.isAtMaxLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_research_getTotalCost(t *testing.T) {
	rules := NewRules()
	type args struct {
		player *Player
		field  TechField
		level  int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"level 0", args{testPlayer(), Energy, 0}, 50},
		{"level 26, no more research possible", args{testPlayer().WithTechLevels(TechLevel{Energy: 26}), Energy, 26}, 0},
		{"level 0 prop with 10 energy", args{testPlayer().WithTechLevels(TechLevel{Energy: 10}), Propulsion, 0}, 150},
		{"level 11 energy with 10 energy", args{testPlayer().WithTechLevels(TechLevel{Energy: 10}), Energy, 11}, 9970},
		{
			name: "level 1 energy with extra cost",
			args: args{
				player: NewPlayer(1, NewRace().withResearchCost(ResearchCost{Energy: ResearchCostExtra}).WithSpec(&rules)).withSpec(&rules),
				field:  Energy,
				level:  0,
			},
			want: 87,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &research{
				rules: &rules,
			}
			if got := r.getTotalCost(tt.args.player, tt.args.field, tt.args.level); got != tt.want {
				t.Errorf("research.getTotalCost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_research_researchFieldOnce(t *testing.T) {
	rules := NewRules()
	playerWithResearchSpent := testPlayer()
	playerWithResearchSpent.TechLevelsSpent = TechLevel{Energy: 25}
	type args struct {
		player           *Player
		field            TechField
		resourcesToSpend int
	}
	tests := []struct {
		name                      string
		args                      args
		wantLevelGained           bool
		wantResourcesLeftover     int
		wantPlayerTechLevels      TechLevel
		wantPlayerTechLevelsSpent TechLevel
	}{
		{
			name: "1st level energy, no gain",
			args: args{
				player:           testPlayer(),
				field:            Energy,
				resourcesToSpend: 25,
			},
			wantPlayerTechLevelsSpent: TechLevel{Energy: 25},
		},
		{
			name: "1st level energy, gain level",
			args: args{
				player:           testPlayer(),
				field:            Energy,
				resourcesToSpend: 51,
			},
			wantLevelGained:           true,
			wantResourcesLeftover:     1,
			wantPlayerTechLevels:      TechLevel{Energy: 1},
			wantPlayerTechLevelsSpent: TechLevel{Energy: 0},
		},
		{
			name: "spent 25 energy already",
			args: args{
				player:           playerWithResearchSpent,
				field:            Energy,
				resourcesToSpend: 26,
			},
			wantLevelGained:           true,
			wantResourcesLeftover:     1,
			wantPlayerTechLevels:      TechLevel{Energy: 1},
			wantPlayerTechLevelsSpent: TechLevel{Energy: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &research{
				rules: &rules,
			}
			gotLevelGained, gotResourcesLeftover := r.researchFieldOnce(tt.args.player, tt.args.field, tt.args.resourcesToSpend)
			if gotLevelGained != tt.wantLevelGained {
				t.Errorf("research.research() gotLevelGained = %v, want %v", gotLevelGained, tt.wantLevelGained)
			}
			if gotResourcesLeftover != tt.wantResourcesLeftover {
				t.Errorf("research.research() gotResourcesLeftover = %v, want %v", gotResourcesLeftover, tt.wantResourcesLeftover)
			}
			if tt.args.player.TechLevels != tt.wantPlayerTechLevels {
				t.Errorf("research.research() gotPlayerTechLevels = %v, want %v", tt.args.player.TechLevels, tt.wantPlayerTechLevels)
			}
			if tt.args.player.TechLevelsSpent != tt.wantPlayerTechLevelsSpent {
				t.Errorf("research.research() gotPlayerTechLevelsSpent = %v, want %v", tt.args.player.TechLevelsSpent, tt.wantPlayerTechLevelsSpent)
			}
		})
	}
}

func Test_research_researchField(t *testing.T) {
	rules := NewRules()
	type args struct {
		player           *Player
		field            TechField
		resourcesToSpend int
	}
	tests := []struct {
		name                      string
		args                      args
		wantPlayerTechLevels      TechLevel
		wantPlayerTechLevelsSpent TechLevel
	}{
		{"research energy partially", args{testPlayer(), Energy, 25}, TechLevel{}, TechLevel{Energy: 25}},
		{"research energy gain a couple levels, leave some", args{testPlayer(), Energy, 150}, TechLevel{Energy: 2}, TechLevel{Energy: 10}},
		{"research max", args{testPlayer().WithTechLevels(TechLevel{Energy: 26}), Energy, 10}, TechLevel{Energy: 26}, TechLevel{Energy: 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &research{
				rules: &rules,
			}
			r.researchField(tt.args.player, tt.args.field, tt.args.resourcesToSpend, func(player *Player, field TechField) {})
			if tt.args.player.TechLevels != tt.wantPlayerTechLevels {
				t.Errorf("research.research() gotPlayerTechLevels = %v, want %v", tt.args.player.TechLevels, tt.wantPlayerTechLevels)
			}
			if tt.args.player.TechLevelsSpent != tt.wantPlayerTechLevelsSpent {
				t.Errorf("research.research() gotPlayerTechLevelsSpent = %v, want %v", tt.args.player.TechLevelsSpent, tt.wantPlayerTechLevelsSpent)
			}

		})
	}
}

func Test_research_research(t *testing.T) {
	rules := NewRules()
	type args struct {
		player           *Player
		resourcesToSpend int
	}
	tests := []struct {
		name                      string
		args                      args
		wantSpent                 TechLevel
		wantPlayerTechLevels      TechLevel
		wantPlayerTechLevelsSpent TechLevel
	}{
		{
			name: "research none",
			args: args{testPlayer().WithResearching(TechFieldNone), 100},
		},
		{
			name:                      "spend 25 on energy, don't gain level",
			args:                      args{testPlayer(), 25},
			wantSpent:                 TechLevel{Energy: 25},
			wantPlayerTechLevels:      TechLevel{},
			wantPlayerTechLevelsSpent: TechLevel{Energy: 25},
		},
		{
			name:                      "spend 150 on energy, research lowest, gain level",
			args:                      args{testPlayer().WithNextResearchField(NextResearchFieldLowestField), 150},
			wantSpent:                 TechLevel{Energy: 50, Weapons: 60, Propulsion: 40},
			wantPlayerTechLevels:      TechLevel{Energy: 1, Weapons: 1},
			wantPlayerTechLevelsSpent: TechLevel{Energy: 0, Weapons: 0, Propulsion: 40},
		},
		{
			name: "spend 25 on energy with previous spent, gain level",
			args: args{
				player: testPlayer().
					WithTechLevelsSpent(TechLevel{Energy: 25}),
				resourcesToSpend: 25,
			},
			wantSpent:                 TechLevel{Energy: 25},
			wantPlayerTechLevels:      TechLevel{Energy: 1},
			wantPlayerTechLevelsSpent: TechLevel{Energy: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &research{
				rules: &rules,
			}
			gotSpent := r.research(tt.args.player, tt.args.resourcesToSpend, func(player *Player, field TechField) {})
			if gotSpent != tt.wantSpent {
				t.Errorf("research.research() gotSpent = %v, want %v", gotSpent, tt.wantSpent)
			}
			if tt.args.player.TechLevels != tt.wantPlayerTechLevels {
				t.Errorf("research.research() gotPlayerTechLevels = %v, want %v", tt.args.player.TechLevels, tt.wantPlayerTechLevels)
			}
			if tt.args.player.TechLevelsSpent != tt.wantPlayerTechLevelsSpent {
				t.Errorf("research.research() gotPlayerTechLevelsSpent = %v, want %v", tt.args.player.TechLevelsSpent, tt.wantPlayerTechLevelsSpent)
			}

		})
	}
}
