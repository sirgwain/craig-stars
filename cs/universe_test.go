package cs

import "testing"

func TestUniverse_getNextFleetNum(t *testing.T) {
	rules := NewRules()

	tests := []struct {
		name   string
		fleets []*Fleet
		want   int
	}{
		{"No Fleets", []*Fleet{}, 1},
		{"Simple fleet", []*Fleet{{MapObject: MapObject{PlayerNum: 1, Num: 1}}}, 2},
		{"Skipped num in fleets", []*Fleet{
			{MapObject: MapObject{PlayerNum: 1, Num: 1}},
			{MapObject: MapObject{PlayerNum: 1, Num: 3}},
		}, 2},
		{"Skipped num in fleets 2", []*Fleet{
			{MapObject: MapObject{PlayerNum: 1, Num: 1}},
			{MapObject: MapObject{PlayerNum: 1, Num: 2}},
			{MapObject: MapObject{PlayerNum: 1, Num: 3}},
			{MapObject: MapObject{PlayerNum: 1, Num: 5}},
		}, 4},
		{"Many fleets", []*Fleet{
			{MapObject: MapObject{PlayerNum: 1, Num: 1}},
			{MapObject: MapObject{PlayerNum: 1, Num: 2}},
			{MapObject: MapObject{PlayerNum: 1, Num: 3}},
			{MapObject: MapObject{PlayerNum: 1, Num: 4}},
		}, 5},
		{"Out of order", []*Fleet{
			{MapObject: MapObject{PlayerNum: 1, Num: 4}},
			{MapObject: MapObject{PlayerNum: 1, Num: 2}},
			{MapObject: MapObject{PlayerNum: 1, Num: 3}},
			{MapObject: MapObject{PlayerNum: 1, Num: 1}},
		}, 5},
		{"Out of order, missing number", []*Fleet{
			{MapObject: MapObject{PlayerNum: 1, Num: 4}},
			{MapObject: MapObject{PlayerNum: 1, Num: 2}},
			{MapObject: MapObject{PlayerNum: 1, Num: 1}},
		}, 3},
		{"Multiple fleet num 1 for starbases", []*Fleet{
			{MapObject: MapObject{PlayerNum: 1, Num: 1}},
			{MapObject: MapObject{PlayerNum: 1, Num: 1}},
			{MapObject: MapObject{PlayerNum: 1, Num: 1}},
			{MapObject: MapObject{PlayerNum: 1, Num: 2}},
		}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Universe{
				Fleets: tt.fleets,
				rules: &rules,
			}
			if got := u.getNextFleetNum(1); got != tt.want {
				t.Errorf("Player.getNextFleetNum() = %v, want %v", got, tt.want)
			}
		})
	}
}
