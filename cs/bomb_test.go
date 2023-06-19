package cs

import (
	"testing"
)

func Test_bomb_getColonistsKilledForBombs(t *testing.T) {
	rules := NewRules()

	type args struct {
		population      int
		defenseCoverage float64
		bombs           []Bomb
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "No bombs",
			args: args{
				population:      1000,
				defenseCoverage: 0.5,
				bombs:           []Bomb{},
			},
			want: 0.0,
		},
		{
			name: "One bomb, no defense",
			args: args{
				population:      10000,
				defenseCoverage: 0.0,
				bombs: []Bomb{
					{Quantity: 10, KillRate: 2.5},
				},
			},
			want: 2500,
		},
		{
			name: "One bomb, partial defense",
			args: args{
				population:      10000,
				defenseCoverage: 0.9792,
				bombs: []Bomb{
					{Quantity: 10, KillRate: 2.5},
				},
			},
			want: 52.0,
		},
		{
			name: "Multiple bombs, no defense",
			args: args{
				population:      10000,
				defenseCoverage: 0.0,
				bombs: []Bomb{
					{Quantity: 10, KillRate: 2.5},
					{Quantity: 5, KillRate: 1.2},
				},
			},
			want: 3100.0,
		},
		{
			name: "Multiple bombs, partial defense",
			args: args{
				population:      10000,
				defenseCoverage: 0.9792,
				bombs: []Bomb{
					{Quantity: 10, KillRate: 2.5},
					{Quantity: 5, KillRate: 1.2},
				},
			},
			want: 64.48,
		},
		{
			name: "One bomb, full defense",
			args: args{
				population:      1000,
				defenseCoverage: 1.0,
				bombs: []Bomb{
					{Quantity: 1, KillRate: 50.0},
				},
			},
			want: 0.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bomb{
				rules: &rules,
			}
			if got := b.getColonistsKilledForBombs(tt.args.population, tt.args.defenseCoverage, tt.args.bombs); int(got) != int(tt.want) {
				t.Errorf("bomb.getColonistsKilledForBombs() = %v, want %v", int(got), int(tt.want))
			}
		})
	}
}

func Test_bomb_getStructuresDestroyed(t *testing.T) {
	rules := NewRules()
	type args struct {
		defenseCoverage float64
		bombs           []Bomb
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "No bombs",
			args: args{
				defenseCoverage: 0.5,
				bombs:           []Bomb{},
			},
			want: 0.0,
		},
		{
			name: "One bomb, no defense",
			args: args{
				defenseCoverage: 0.0,
				bombs: []Bomb{
					{Quantity: 10, KillRate: 2.5},
				},
			},
			want: 2500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bomb{
				rules: &rules,
			}
			if got := b.getStructuresDestroyed(tt.args.defenseCoverage, tt.args.bombs); got != tt.want {
				t.Errorf("bomb.getStructuresDestroyed() = %v, want %v", got, tt.want)
			}
		})
	}
}
