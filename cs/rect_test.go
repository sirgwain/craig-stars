package cs

import (
	"math"
	"reflect"
	"testing"
)

func TestRect_Center(t *testing.T) {
	type fields struct {
		X      float64
		Y      float64
		Width  float64
		Height float64
	}
	tests := []struct {
		name   string
		fields fields
		want   Vector
	}{
		{"0, 0", fields{-1, -1, 2, 2}, Vector{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rect := Rect{
				X:      tt.fields.X,
				Y:      tt.fields.Y,
				Width:  tt.fields.Width,
				Height: tt.fields.Height,
			}
			if got := rect.Center(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Rect.Center() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRect_PointInRectangle(t *testing.T) {
	type fields struct {
		X      float64
		Y      float64
		Width  float64
		Height float64
	}
	type args struct {
		point Vector
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"in", fields{-1, -1, 2, 2}, args{Vector{1, 1}}, true},
		{"out", fields{-1, -1, 2, 2}, args{Vector{2, 1}}, false},
		{"out", fields{-1, -1, 2, 2}, args{Vector{-2, 1}}, false},
		{"out", fields{-1, -1, 2, 2}, args{Vector{1, 2}}, false},
		{"out", fields{-1, -1, 2, 2}, args{Vector{1, -2}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rect := Rect{
				X:      tt.fields.X,
				Y:      tt.fields.Y,
				Width:  tt.fields.Width,
				Height: tt.fields.Height,
			}
			if got := rect.PointInRectangle(tt.args.point); got != tt.want {
				t.Errorf("Rect.PointInRectangle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRect_PointInRotatedRectangle(t *testing.T) {
	type fields struct {
		X      float64
		Y      float64
		Width  float64
		Height float64
	}
	type args struct {
		point     Vector
		rectAngle float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"0, 0 in", fields{-1, -1, 2, 2}, args{Vector{0, 0}, math.Pi / 4}, true},
		{"1.4ish, 0 in", fields{-1, -1, 2, 2}, args{Vector{math.Sqrt(2)-.01, 0}, -math.Pi / 4}, true},
		{"out", fields{-1, -1, 2, 2}, args{Vector{1, 1}, math.Pi / 4}, false},
		{"out", fields{-1, -1, 2, 2}, args{Vector{-1, 1}, math.Pi / 4}, false},
		{"out", fields{-1, -1, 2, 2}, args{Vector{1, 1}, math.Pi / 4}, false},
		{"out", fields{-1, -1, 2, 2}, args{Vector{1, -1}, math.Pi / 4}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rect := Rect{
				X:      tt.fields.X,
				Y:      tt.fields.Y,
				Width:  tt.fields.Width,
				Height: tt.fields.Height,
			}
			if got := rect.PointInRotatedRectangle(tt.args.point, tt.args.rectAngle); got != tt.want {
				t.Errorf("Rect.PointInRotatedRectangle() = %v, want %v", got, tt.want)
			}
		})
	}
}
