package cs

import (
	"reflect"
	"testing"
)

func TestCargo_CanTransfer(t *testing.T) {
	type args struct {
		transferAmount Cargo
	}
	tests := []struct {
		name  string
		cargo Cargo
		args  args
		want  bool
	}{
		{"Can transfer", Cargo{1, 2, 3, 4}, args{Cargo{1, 2, 3, 4}}, true},
		{"Cannot transfer", Cargo{1, 2, 3, 4}, args{Cargo{1, 2, 3, 5}}, false},
		{"Cannot transfer", Cargo{1, 2, 3, 4}, args{Cargo{0, 0, 0, 5}}, false},
		{"Can transfer", Cargo{1, 2, 3, 4}, args{Cargo{1, 0, 0, 0}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cargo.CanTransfer(tt.args.transferAmount); got != tt.want {
				t.Errorf("Cargo.CanTransfer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCargo_Split(t *testing.T) {
	type args struct {
		sourceCapacity int
		capacity1      int
		capacity2      int
	}
	tests := []struct {
		name       string
		source     Cargo
		args       args
		wantCargo1 Cargo
		wantCargo2 Cargo
		wantErr    bool
	}{
		{
			name:    "invalid args",
			source:  Cargo{1, 1, 1, 1},
			args:    args{sourceCapacity: 4, capacity1: 5, capacity2: 0},
			wantErr: true,
		},
		{
			name:       "no split",
			source:     Cargo{1, 1, 1, 1},
			args:       args{sourceCapacity: 4, capacity1: 4, capacity2: 0},
			wantCargo1: Cargo{1, 1, 1, 1},
			wantCargo2: Cargo{},
		},
		{
			name:       "split all",
			source:     Cargo{1, 1, 1, 1},
			args:       args{sourceCapacity: 4, capacity1: 0, capacity2: 4},
			wantCargo1: Cargo{},
			wantCargo2: Cargo{1, 1, 1, 1},
		},
		{
			name:       "split even",
			source:     Cargo{2, 2, 2, 2},
			args:       args{sourceCapacity: 8, capacity1: 4, capacity2: 4},
			wantCargo1: Cargo{1, 1, 1, 1},
			wantCargo2: Cargo{1, 1, 1, 1},
		},
		{
			name:       "split 1",
			source:     Cargo{2, 2, 2, 2},
			args:       args{sourceCapacity: 8, capacity1: 7, capacity2: 1},
			wantCargo1: Cargo{2, 2, 2, 1},
			wantCargo2: Cargo{0, 0, 0, 1},
		},
		{
			name:       "split 2",
			source:     Cargo{2, 2, 2, 2},
			args:       args{sourceCapacity: 8, capacity1: 6, capacity2: 2},
			wantCargo1: Cargo{2, 2, 2, 0},
			wantCargo2: Cargo{0, 0, 0, 2},
		},
		{
			name:       "split 3",
			source:     Cargo{2, 2, 2, 2},
			args:       args{sourceCapacity: 8, capacity1: 5, capacity2: 3},
			wantCargo1: Cargo{1, 1, 1, 2},
			wantCargo2: Cargo{1, 1, 1, 0},
		},
		{
			name:       "split many",
			source:     Cargo{100, 200, 300, 400},
			args:       args{sourceCapacity: 2000, capacity1: 1500, capacity2: 500},
			wantCargo1: Cargo{75, 150, 225, 300},
			wantCargo2: Cargo{25, 50, 75, 100},
		},
		{
			name:       "split 33%",
			source:     Cargo{100, 0, 0, 0},
			args:       args{sourceCapacity: 300, capacity1: 200, capacity2: 100},
			wantCargo1: Cargo{67, 0, 0, 0},
			wantCargo2: Cargo{33, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.source.Split(tt.args.sourceCapacity, tt.args.capacity1, tt.args.capacity2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Cargo.Split() errored unexpectedly; err = %v", err)
			}

			if !reflect.DeepEqual(got, tt.wantCargo1) {
				t.Errorf("Cargo.Split() got = %v, want %v", got, tt.wantCargo1)
			}
			if !reflect.DeepEqual(got1, tt.wantCargo2) {
				t.Errorf("Cargo.Split() got1 = %v, want %v", got1, tt.wantCargo2)
			}
		})
	}
}
