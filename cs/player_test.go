package cs

import (
	"reflect"
	"testing"
)

// create a new test player with humanoid race and computed specs
func testPlayer() *Player {
	race := Humanoids()
	race.Spec = computeRaceSpec(&race, &rules)
	return NewPlayer(1, &race).withSpec(&rules)
}

func TestPlayer_HasTech(t *testing.T) {

	type args struct {
		player *Player
		tech   *Tech
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Everyone can learn the Viewer50", args{NewPlayer(1, NewRace()), &Viewer50.Tech}, true},
		{"Player is missing tech levels for the Viewer90", args{NewPlayer(1, NewRace()), &Viewer90.Tech}, false},
		{"Player is missing tech levels for the TransGalacticDrive", args{NewPlayer(1, NewRace()), &TransGalacticDrive.Tech}, false},
		{"Player has tech levels for the TransGalacticDrive", args{NewPlayer(1, NewRace()).WithTechLevels(TechLevel{Propulsion: 9}), &TransGalacticDrive.Tech}, true},
		{"Only players with IFE can learn the FuelMizer", args{NewPlayer(1, NewRace()), &FuelMizer.Tech}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.player.HasTech(tt.args.tech); got != tt.want {
				t.Errorf("Player.HasTech() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayer_CanLearnTech(t *testing.T) {

	type args struct {
		player *Player
		tech   *Tech
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Everyone can learn the Viewer50", args{NewPlayer(1, NewRace()), &Viewer50.Tech}, true},
		{"Everyone can learn the Viewer90", args{NewPlayer(1, NewRace()), &Viewer90.Tech}, true},
		{"Only players with IFE can learn the FuelMizer", args{NewPlayer(1, NewRace()), &FuelMizer.Tech}, false},
		{"Only players with IFE can learn the FuelMizer", args{NewPlayer(1, NewRace().WithLRT(IFE)), &FuelMizer.Tech}, true},
		{"Players with NRSE cannot learn the ramscoops", args{NewPlayer(1, NewRace().WithLRT(NRSE)), &FuelMizer.Tech}, false},
		{"IS can learn speed trap 20", args{NewPlayer(1, NewRace().WithPRT(IS)), &SpeedTrap20.Tech}, true},
		{"SD can learn speed trap 20", args{NewPlayer(1, NewRace().WithPRT(SD)), &SpeedTrap20.Tech}, true},
		{"WM cannot learn speed trap 20", args{NewPlayer(1, NewRace().WithPRT(WM)), &SpeedTrap20.Tech}, false},
		{"IS cannot learn smart bombs", args{NewPlayer(1, NewRace().WithPRT(IS)), &SmartBomb.Tech}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.player.CanLearnTech(tt.args.tech); got != tt.want {
				t.Errorf("Player.CanLearnTech() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUniverse_getNextFleetNum(t *testing.T) {
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
			player := testPlayer()
			if got := player.getNextFleetNum(tt.fleets); got != tt.want {
				t.Errorf("Player.getNextFleetNum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayer_defaultRelationships(t *testing.T) {
	type fields struct {
		aiControlled     bool
		aiFormsAlliances bool
	}
	type args struct {
		players []*Player
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []PlayerRelationship
	}{
		{
			name: "single player", args: args{players: []*Player{}},
			want: []PlayerRelationship{{Relation: PlayerRelationFriend}},
		},
		{
			name: "against ai", args: args{players: []*Player{NewPlayer(0, NewRace()).WithNum(2).WithAIControlled(true)}},
			want: []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationEnemy}},
		},
		{
			name: "against player and ai",
			args: args{players: []*Player{
				NewPlayer(0, NewRace()).WithNum(2).WithAIControlled(false),
				NewPlayer(0, NewRace()).WithNum(3).WithAIControlled(true),
			}},
			want: []PlayerRelationship{
				{Relation: PlayerRelationFriend},
				{Relation: PlayerRelationNeutral},
				{Relation: PlayerRelationEnemy},
			},
		},
		{
			name: "ai against ai", fields: fields{aiControlled: true}, args: args{players: []*Player{NewPlayer(0, NewRace()).WithNum(2).WithAIControlled(true)}},
			want: []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationEnemy}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPlayer(0, NewRace()).WithNum(1).WithAIControlled(tt.fields.aiControlled)

			if got := p.defaultRelationships(append([]*Player{p}, tt.args.players...), tt.fields.aiFormsAlliances); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Player.defaultRelationships() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayer_GetResearchCost(t *testing.T) {
	type fields struct {
		techLevels      TechLevel
		techLevelsSpent TechLevel
	}
	type args struct {
		techLevel TechLevel
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name:   "No tech, no cost",
			fields: fields{techLevels: TechLevel{}, techLevelsSpent: TechLevel{}},
			args:   args{techLevel: TechLevel{}},
			want:   0,
		},
		{
			name:   "first level",
			fields: fields{techLevels: TechLevel{}, techLevelsSpent: TechLevel{}},
			args:   args{techLevel: TechLevel{Energy: 1}},
			want:   50,
		},
		{
			name:   "first level, spent some",
			fields: fields{techLevels: TechLevel{}, techLevelsSpent: TechLevel{Energy: 30}},
			args:   args{techLevel: TechLevel{Energy: 1}},
			want:   20,
		},
		{
			name:   "multiple fields, multiple levels",
			fields: fields{techLevels: TechLevel{}, techLevelsSpent: TechLevel{}},
			args:   args{techLevel: TechLevel{Energy: 2, Weapons: 3}},
			want:   490,
		},
		{
			name:   "multiple fields, multiple levels, spent some",
			fields: fields{techLevels: TechLevel{}, techLevelsSpent: TechLevel{Energy: 10, Weapons: 20}},
			args:   args{techLevel: TechLevel{Energy: 2, Weapons: 3}},
			want:   460,
		},
		{
			name:   "multiple fields, multiple levels, only need 2 energy",
			fields: fields{techLevels: TechLevel{Energy: 1, Weapons: 3}, techLevelsSpent: TechLevel{}},
			args:   args{techLevel: TechLevel{Energy: 3, Weapons: 3}},
			want:   300,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPlayer(0, NewRace().WithSpec(&rules))
			p.TechLevels = tt.fields.techLevels
			p.TechLevelsSpent = tt.fields.techLevelsSpent

			if got := p.GetResearchCost(&rules, tt.args.techLevel); got != tt.want {
				t.Errorf("Player.GetResearchCost() = %v, want %v", got, tt.want)
			}
		})
	}
}


