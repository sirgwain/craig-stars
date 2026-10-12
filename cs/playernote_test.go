package cs

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlayerNote_Validate(t *testing.T) {
	area := Vector{X: 400, Y: 400}
	tests := []struct {
		name    string
		note    PlayerNote
		want    PlayerNote
		wantErr bool
	}{
		{
			name: "planet note is trimmed and loses its target player",
			note: PlayerNote{TargetType: MapObjectTypePlanet, TargetNum: 3, TargetPlayerNum: 2, Text: "  germ rich \n"},
			want: PlayerNote{TargetType: MapObjectTypePlanet, TargetNum: 3, Text: "germ rich"},
		},
		{
			name: "label is always shown on the map",
			note: PlayerNote{Position: Vector{X: 10, Y: 20}, Text: "the front"},
			want: PlayerNote{Position: Vector{X: 10, Y: 20}, Text: "the front", ShowOnMap: true},
		},
		{name: "empty text", note: PlayerNote{TargetType: MapObjectTypePlanet, TargetNum: 3, Text: "  "}, wantErr: true},
		{name: "too long", note: PlayerNote{Position: Vector{X: 1, Y: 1}, Text: strings.Repeat("a", PlayerNoteMaxLength+1)}, wantErr: true},
		{name: "fleet nums are reused", note: PlayerNote{TargetType: MapObjectTypeFleet, TargetNum: 1, TargetPlayerNum: 1, Text: "decoy"}, wantErr: true},
		{name: "planet without a num", note: PlayerNote{TargetType: MapObjectTypePlanet, Text: "hi"}, wantErr: true},
		{name: "label outside the universe", note: PlayerNote{Position: Vector{X: 401, Y: 20}, Text: "hi"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.note.Validate(area)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, tt.note)
		})
	}
}
