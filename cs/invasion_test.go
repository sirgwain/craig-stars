package cs

import (
	"reflect"
	"testing"
)

func Test_invasion_resolve(t *testing.T) {
	tests := []struct {
		name     string
		invasion invasion
		want     invasionResult
	}{
		{
			name: "10000 attackers 10000 defenders, attacker wins",
			invasion: invasion{
				planet: &Planet{
					Cargo: Cargo{}.WithPopulation(10_000),
				},
				defender:  NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules),
				attacker:  NewPlayer(2, NewRace().WithSpec(&rules)).WithNum(2).withSpec(&rules),
				attackers: 10_000,
			},
			want: invasionResult{
				defenders:          10_000,
				attackersKilled:    9_100,
				defendersKilled:    10_000,
				remainingAttackers: 900,
				remainingDefenders: 0,
				successful:         true,
			},
		},
		{
			name: "5000 attackers for 10000 undefended defenders, defenders win",
			invasion: invasion{
				planet: &Planet{
					Cargo: Cargo{}.WithPopulation(10_000),
				},
				defender:  NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules),
				attacker:  NewPlayer(2, NewRace().WithSpec(&rules)).WithNum(2).withSpec(&rules),
				attackers: 5000,
			},
			want: invasionResult{
				defenders:          10_000,
				attackersKilled:    5000,
				defendersKilled:    5500,
				remainingAttackers: 0,
				remainingDefenders: 4500,
				successful:         false,
			},
		},
		{
			name: "100,000 attackers for 100,000 well defended defenders, defenders win",
			invasion: invasion{
				planet: &Planet{
					Cargo: Cargo{}.WithPopulation(100_000),
					Spec:  PlanetSpec{DefenseCoverage: .9},
				},
				defender:  NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules),
				attacker:  NewPlayer(2, NewRace().WithSpec(&rules)).WithNum(2).withSpec(&rules),
				attackers: 100_000,
			},
			want: invasionResult{
				defenders:          100_000,
				attackersKilled:    100_000,
				defendersKilled:    35700,
				remainingAttackers: 0,
				remainingDefenders: 64_300,
				successful:         false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.invasion.resolve(&rules)
			// zero out the invasion itself in the result, we only care about the result numbers
			got.invasion = invasion{}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("invasion.resolve() = %v, want %v", got, tt.want)
			}
		})
	}
}
