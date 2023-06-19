package game

import (
	"testing"
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
		{"Terrible Hab", args{Hab{1, 1, 1}, &Race{HabLow: Hab{99, 99, 99}, HabHigh: Hab{100, 100, 100}}}, -45},
		// {"Moderate Hab", args{Hab{50, 50, 30}, NewRace()}, 75}, TODO: make sure this works with regular stars data
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.race.GetPlanetHabitability(tt.args.hab); got != tt.want {
				t.Errorf("Race.GetPlanetHabitability() = %v, want %v", got, tt.want)
			}
		})
	}
}
