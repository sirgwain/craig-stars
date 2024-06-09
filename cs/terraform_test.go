package cs

import (
	"reflect"
	"testing"
)

func Test_getTerraformAbility(t *testing.T) {
	type args struct {
		player *Player
	}
	tests := []struct {
		name string
		args args
		want Hab
	}{
		{
			name: "No ability",
			args: args{
				player: NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules),
			},
			want: Hab{},
		},
		{
			name: "Humanoid starting tech",
			args: args{
				player: NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).WithTechLevels(TechLevel{3, 3, 3, 3, 3, 3}).withSpec(&rules),
			},
			want: Hab{3, 3, 3},
		},
		{
			name: "Race with TT",
			args: args{
				player: NewPlayer(1, NewRace().WithLRT(TT).WithSpec(&rules)).WithNum(1).withSpec(&rules),
			},
			want: Hab{3, 3, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			terraformer := terraform{}
			if got := terraformer.getTerraformAbility(tt.args.player); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTerraformAbility() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTerraformAmount(t *testing.T) {
	planet := &Planet{
		BaseHab:           Hab{Grav: 50, Temp: 50, Rad: 50},
		Hab:               Hab{Grav: 47, Temp: 53, Rad: 50},
		TerraformedAmount: Hab{Grav: 0, Temp: 0, Rad: 0},
	}

	tests := []struct {
		name        string
		planet      *Planet
		player      *Player
		terraformer *Player
		expected    Hab
	}{
		{
			name:        "player is nil",
			planet:      planet,
			player:      nil,
			terraformer: NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules),
			expected:    Hab{},
		},
		{
			name:        "player and terraformer are nil",
			planet:      planet,
			player:      nil,
			terraformer: nil,
			expected:    Hab{},
		},
		{
			name:        "terraformer is nil, but player can terraform",
			planet:      planet,
			player:      NewPlayer(1, NewRace().WithLRT(TT).WithSpec(&rules)).WithNum(1).withSpec(&rules),
			terraformer: nil,
			expected:    Hab{Grav: 3, Temp: -3, Rad: 0},
		},
		{
			name:     "can terraform all habs up to ability",
			planet:   planet,
			player:   NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).WithTechLevels(TechLevel{3, 3, 3, 3, 3, 3}).withSpec(&rules),
			expected: Hab{Grav: 3, Temp: -3, Rad: 0},
		},
		{
			name: "can terraform some habs up to ability",
			planet: &Planet{
				BaseHab: Hab{Grav: 45, Temp: 45, Rad: 45},
				Hab:     Hab{Grav: 46, Temp: 48, Rad: 47},
			},
			player:   NewPlayer(1, NewRace().WithLRT(TT).WithSpec(&rules)).WithNum(1).withSpec(&rules),
			expected: Hab{Grav: 2, Temp: 0, Rad: 1},
		},
		{
			name: "can terraform some habs up to ability both directions",
			planet: &Planet{
				BaseHab: Hab{Grav: 47, Temp: 54, Rad: 50},
				Hab:     Hab{Grav: 48, Temp: 53, Rad: 50},
			},
			player:   NewPlayer(1, NewRace().WithLRT(TT).WithSpec(&rules)).WithNum(1).withSpec(&rules),
			expected: Hab{Grav: 2, Temp: -2, Rad: 0},
		},
		{
			name: "high ability, mostly terraformed",
			planet: &Planet{
				BaseHab: Hab{Grav: 50, Temp: 33, Rad: 23},
				Hab:     Hab{Grav: 50, Temp: 41, Rad: 34},
			},
			player:   NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).WithTechLevels(TechLevel{11, 14, 11, 15, 10, 10}).withSpec(&rules),
			expected: Hab{Grav: 0, Temp: 3, Rad: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			terraformer := terraform{}
			terraformAmount := terraformer.getTerraformAmount(tt.planet.Hab, tt.planet.BaseHab, tt.player, tt.terraformer)
			if !reflect.DeepEqual(terraformAmount, tt.expected) {
				t.Errorf("Got terraform amount %v, expected %v", terraformAmount, tt.expected)
			}
		})
	}
}

func Test_getMinTerraformAmount(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).WithTechLevels(TechLevel{3, 3, 3, 3, 3, 3}).withSpec(&rules)

	type args struct {
		planet      *Planet
		player      *Player
		terraformer *Player
	}
	tests := []struct {
		name string
		args args
		want Hab
	}{
		{
			name: "nil player",
			args: args{
				planet: &Planet{},
			},
			want: Hab{},
		},
		{
			name: "off by one, still in range, no terraform",
			args: args{
				player: player,
				planet: &Planet{
					BaseHab: Hab{Grav: 50, Temp: 50, Rad: 50},
					Hab:     Hab{Grav: 49, Temp: 51, Rad: 50},
				},
			},
			want: Hab{},
		},
		{
			name: "out of range by one/two, terraform",
			args: args{
				player: player,
				planet: &Planet{
					BaseHab: Hab{Grav: 14, Temp: 87, Rad: 50},
					Hab:     Hab{Grav: 14, Temp: 87, Rad: 50},
				},
			},
			want: Hab{Grav: 1, Temp: 2, Rad: 0},
		},
		{
			name: "out of range by two, but we've already terraformed one",
			args: args{
				player: player,
				planet: &Planet{
					BaseHab: Hab{Grav: 13, Temp: 50, Rad: 50},
					Hab:     Hab{Grav: 14, Temp: 50, Rad: 50},
				},
			},
			want: Hab{Grav: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			terraformer := terraform{}
			if got := terraformer.getMinTerraformAmount(tt.args.planet.Hab, tt.args.planet.BaseHab, tt.args.player, tt.args.terraformer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMinTerraformAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBestTerraform(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).WithTechLevels(TechLevel{3, 3, 3, 3, 3, 3}).withSpec(&rules)

	grav := int(Grav)
	rad := int(Rad)

	type args struct {
		planet      *Planet
		player      *Player
		terraformer *Player
	}
	tests := []struct {
		name string
		args args
		want *HabType
	}{
		{
			name: "nil player",
			args: args{
				planet: &Planet{},
			},
			want: nil,
		},
		{
			name: "best hab to terraform is gravity",
			args: args{
				player: player,
				planet: &Planet{
					BaseHab: Hab{Grav: 49, Temp: 50, Rad: 50},
					Hab:     Hab{Grav: 49, Temp: 50, Rad: 50},
				},
			},
			want: (*HabType)(&grav),
		},
		{
			name: "best hab to terraform is rad",
			args: args{
				player: player,
				planet: &Planet{
					BaseHab: Hab{Grav: 49, Temp: 50, Rad: 52}, // terraforming grav/rad are equivalent in hab bonus, so pick rad because it's farther away
					Hab:     Hab{Grav: 49, Temp: 50, Rad: 52},
				},
			},
			want: (*HabType)(&rad),
		},
		{
			name: "best hab to terraform is rad because it's less red",
			args: args{
				player: player,
				planet: &Planet{
					BaseHab: Hab{Grav: 13, Temp: 50, Rad: 14},
					Hab:     Hab{Grav: 13, Temp: 50, Rad: 14},
				},
			},
			want: (*HabType)(&rad),
		},
		{
			name: "best hab to terraform is rad because it's less red other direction",
			args: args{
				player: player,
				planet: &Planet{
					BaseHab: Hab{Grav: 13, Temp: 50, Rad: 86},
					Hab:     Hab{Grav: 13, Temp: 50, Rad: 86},
				},
			},
			want: (*HabType)(&rad),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			terraformer := terraform{}
			got := terraformer.GetBestTerraform(tt.args.planet, tt.args.player, tt.args.terraformer)
			if got == nil && tt.want != nil {
				t.Errorf("GetBestTerraform() = %v, want %v", got, *tt.want)
			} else if got != nil && tt.want != nil && *got != *tt.want {
				t.Errorf("GetBestTerraform() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func Test_getBestUnterraform(t *testing.T) {
	player := NewPlayer(1, NewRace().WithLRT(TT).WithSpec(&rules)).WithNum(1).withSpec(&rules)

	temp := int(Temp)
	grav := int(Grav)

	type args struct {
		planet      *Planet
		player      *Player
		terraformer *Player
	}
	tests := []struct {
		name string
		args args
		want *HabType
	}{
		{
			name: "nil player",
			args: args{
				planet: &Planet{},
			},
			want: nil,
		},
		{
			name: "can't unterraform past our abilities",
			args: args{
				player: player,
				planet: &Planet{
					Hab:               Hab{Grav: 47, Temp: 47, Rad: 47},
					TerraformedAmount: Hab{Grav: -3, Temp: 3, Rad: -3},
				},
			},
			want: nil,
		},
		{
			name: "best hab to unterraform is temp",
			args: args{
				player: player,
				planet: &Planet{
					Hab:               Hab{Grav: 47, Temp: 48, Rad: 47},
					TerraformedAmount: Hab{Grav: -3, Temp: -2, Rad: -3},
				},
			},
			want: (*HabType)(&temp),
		},
		{
			name: "he farthest from ideal habType so we push this planet into the red",
			args: args{
				player: player,
				planet: &Planet{
					Hab:               Hab{Grav: 50, Temp: 52, Rad: 50},
					TerraformedAmount: Hab{Grav: 0, Temp: 2, Rad: 0},
				},
			},
			want: (*HabType)(&temp),
		},
		{
			name: "perfect planet, pick the first",
			args: args{
				player: player,
				planet: &Planet{
					BaseHab:           Hab{Grav: 50, Temp: 50, Rad: 50},
					Hab:               Hab{Grav: 50, Temp: 50, Rad: 50},
					TerraformedAmount: Hab{Grav: 0, Temp: 0, Rad: 0},
				},
			},
			want: (*HabType)(&grav),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			terraformer := terraform{}
			got := terraformer.getBestUnterraform(tt.args.planet, tt.args.player, tt.args.terraformer)
			if got == nil && tt.want != nil {
				t.Errorf("getWorstTerraform() = %v, want %v", got, *tt.want)
			} else if got != nil && tt.want != nil && *got != *tt.want {
				t.Errorf("getWorstTerraform() = %v, want %v", *got, *tt.want)
			}
		})
	}
}
