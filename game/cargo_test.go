package game

import "testing"

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
