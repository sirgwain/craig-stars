package cs

import (
	"testing"
)

func TestMapObject_String(t *testing.T) {

	tests := []struct {
		name string
		mo   MapObject
		want string
	}{
		{"MapObject String()", MapObject{GameID: 1, ID: 2, Num: 3, Name: "Bob's Revenge"},
			"GameID:     1, ID:     2, Num:   3 Bob's Revenge"},
		{"MapObject String()", MapObject{GameID: 12345, ID: 23456, Num: 120, Name: "Craig's Planet"},
			"GameID: 12345, ID: 23456, Num: 120 Craig's Planet"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mo.String(); got != tt.want {
				t.Errorf("MapObject.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
