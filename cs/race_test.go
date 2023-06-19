package cs

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
		{"Terrible Hab", args{Hab{0, 0, 0}, &Race{HabLow: Hab{99, 99, 99}, HabHigh: Hab{100, 100, 100}}}, -45},
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

func TestRace_ComputeRacePoints(t *testing.T) {
	startingPoints := 1650 // defind in rules

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
