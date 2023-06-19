package cs

import (
	"testing"

	"github.com/sirgwain/craig-stars/test"
)

func Test_invadePlanet(t *testing.T) {
	rules := NewRules()
	type args struct {
		planet           *Planet
		fleet            *Fleet
		defender         *Player
		attacker         *Player
		colonistsDropped int
	}
	tests := []struct {
		name string
		args args
		want Planet
	}{
		{
			name: "10000 attackers 10000 defenders, attacker wins",
			args: args{
				planet: &Planet{
					MapObject: MapObject{
						PlayerNum: 1,
						Name:      "Brin",
					},
					Cargo:     Cargo{}.WithPopulation(10_000),
					Mines:     100,
					Factories: 100,
					Defenses:  0,
				},
				fleet: &Fleet{
					MapObject: MapObject{
						PlayerNum: 2,
						Name:      "Teamster #1",
					},
				},
				defender:         NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules),
				attacker:         NewPlayer(2, NewRace().WithSpec(&rules)).WithNum(2).withSpec(&rules),
				colonistsDropped: 10_000,
			},
			want: Planet{MapObject: MapObject{Name: "Brin", PlayerNum: 2}, Cargo: Cargo{}.WithPopulation(900), Mines: 100, Factories: 100},
		},
		{
			name: "5000 attackers for 10000 undefended defenders, defenders win",
			args: args{
				planet: &Planet{
					MapObject: MapObject{
						PlayerNum: 1,
						Name:      "Brin",
					},
					Cargo:     Cargo{}.WithPopulation(10_000),
					Mines:     100,
					Factories: 100,
					Defenses:  0,
				},
				fleet: &Fleet{
					MapObject: MapObject{
						PlayerNum: 2,
						Name:      "Teamster #1",
					},
				},
				defender:         NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules),
				attacker:         NewPlayer(2, NewRace().WithSpec(&rules)).WithNum(2).withSpec(&rules),
				colonistsDropped: 5000,
			},
			want: Planet{MapObject: MapObject{Name: "Brin", PlayerNum: 1}, Cargo: Cargo{}.WithPopulation(4500), Mines: 100, Factories: 100},
		},
		{
			name: "100,000 attackers for 100,000 well defended defenders, defenders win",
			args: args{
				planet: &Planet{
					MapObject: MapObject{
						PlayerNum: 1,
						Name:      "Brin",
					},
					Cargo:     Cargo{}.WithPopulation(100_000),
					Mines:     100,
					Factories: 100,
					Defenses:  1000,
				},
				fleet: &Fleet{
					MapObject: MapObject{
						PlayerNum: 2,
						Name:      "Teamster #1",
					},
				},
				defender:         NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules),
				attacker:         NewPlayer(2, NewRace().WithSpec(&rules)).WithNum(2).withSpec(&rules),
				colonistsDropped: 100_000,
			},
			want: Planet{MapObject: MapObject{Name: "Brin", PlayerNum: 1}, Cargo: Cargo{}.WithPopulation(42_000), Mines: 100, Factories: 100, Defenses: 1000},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.planet.Spec = computePlanetSpec(&rules, tt.args.attacker, tt.args.planet)
			tt.want.Spec = computePlanetSpec(&rules, tt.args.attacker, &tt.want)
			invadePlanet(tt.args.planet, tt.args.fleet, tt.args.defender, tt.args.attacker, tt.args.colonistsDropped, rules.InvasionDefenseCoverageFactor)

			// recompute planet spec after invasion
			tt.args.planet.Spec = computePlanetSpec(&rules, tt.args.attacker, tt.args.planet)

			if got := *tt.args.planet; !test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("Planet = %v, want %v", got, tt.want)
			}

		})
	}
}
