package cs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_transferPlanetCargo(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules))
	// scout := testLongRangeScout(player)
	freighter := testSmallFreighter(player)
	type args struct {
		planet         *Planet
		transferAmount Cargo
	}
	tests := []struct {
		name    string
		fleet   *Fleet
		args    args
		wantErr bool
	}{
		{"Should transfer from planet", freighter, args{NewPlanet().WithCargo(Cargo{1, 2, 3, 4}), Cargo{1, 0, 0, 0}}, false},
		{"Should fail to transfer from planet", freighter, args{NewPlanet().WithCargo(Cargo{1, 2, 3, 4}), Cargo{2, 0, 0, 0}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourceCargo := tt.fleet.Cargo
			destCargo := tt.args.planet.Cargo
			o := orders{}
			if err := o.TransferPlanetCargo(tt.fleet, tt.args.planet, tt.args.transferAmount); (err != nil) != tt.wantErr {
				t.Errorf("Fleet.TransferPlanetCargo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				// if we didn't want an error, ensure we transferred from the planet to the fleet
				assert.Equal(t, tt.fleet.Cargo, sourceCargo.Add(tt.args.transferAmount))
				assert.Equal(t, tt.args.planet.Cargo, destCargo.Subtract(tt.args.transferAmount))
			}

		})
	}
}

func Test_transferFleetCargo(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules))
	// scout := testLongRangeScout(player)
	freighter := testSmallFreighter(player)
	type args struct {
		fleet          *Fleet
		transferAmount Cargo
	}
	tests := []struct {
		name    string
		fleet   *Fleet
		args    args
		wantErr bool
	}{
		{"Should transfer from fleet", freighter, args{testSmallFreighter(player).withCargo(Cargo{1, 2, 3, 4}), Cargo{1, 0, 0, 0}}, false},
		{"Should fail to transfer from fleet", freighter, args{testSmallFreighter(player).withCargo(Cargo{1, 2, 3, 4}), Cargo{2, 0, 0, 0}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourceCargo := tt.fleet.Cargo
			destCargo := tt.args.fleet.Cargo
			o := orders{}

			if err := o.TransferFleetCargo(tt.fleet, tt.args.fleet, tt.args.transferAmount); (err != nil) != tt.wantErr {
				t.Errorf("Fleet.TransferFleetCargo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				// if we didn't want an error, ensure we transferred from the fleet to the fleet
				assert.Equal(t, tt.fleet.Cargo, sourceCargo.Add(tt.args.transferAmount))
				assert.Equal(t, tt.args.fleet.Cargo, destCargo.Subtract(tt.args.transferAmount))
			}

		})
	}
}
