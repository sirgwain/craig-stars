package cs

import "testing"

func Test_gravString(t *testing.T) {
	type args struct {
		grav int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{".12 min", args{1}, "0.12g"},
		{"8.00 max", args{100}, "8.00g"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gravString(tt.args.grav); got != tt.want {
				t.Errorf("gravString() = %v, want %v", got, tt.want)
			}
		})
	}
}
