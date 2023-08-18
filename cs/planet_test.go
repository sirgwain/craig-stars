package cs

import (
	"reflect"
	"testing"
)

func newTestPlayerPlanet() (player *Player, planet *Planet) {
	player = NewPlayer(1, NewRace())
	player.Race.Spec = computeRaceSpec(&player.Race, &rules)
	planet = &Planet{}
	planet.PlayerNum = player.Num

	player.Spec = computePlayerSpec(player, &rules, []*Planet{planet})

	return player, planet
}

func TestPlanet_String(t *testing.T) {

	tests := []struct {
		name string
		p    *Planet
		want string
	}{
		{"MapObject String()", &Planet{MapObject: MapObject{GameDBObject: GameDBObject{GameID: 1, ID: 2}, Num: 3, Name: "Bob's Revenge"}},
			"Planet GameID:     1, ID:     2, Num:   3 Bob's Revenge"},
		{"MapObject String()", &Planet{MapObject: MapObject{GameDBObject: GameDBObject{GameID: 12345, ID: 23456}, Num: 120, Name: "Craig's Planet"}},
			"Planet GameID: 12345, ID: 23456, Num: 120 Craig's Planet"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.String(); got != tt.want {
				t.Errorf("MapObject.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanet_innateMines(t *testing.T) {
	player := NewPlayer(1, &Race{Spec: RaceSpec{InnateMining: false}})
	planet := Planet{}
	planet.setPopulation(16000)

	if got := planet.innateMines(player); got != 0 {
		t.Errorf("Planet.GetInnateMines() = %v, want %v", got, 0)
	}

	// should get 40 mines for 16k pop when the player has innate mining
	player.Race.Spec.InnateMining = true
	player.Race.Spec.InnatePopulationFactor = .1
	if got := planet.innateMines(player); got != 12 {
		t.Errorf("Planet.GetInnateMines() = %v, want %v", got, 40)
	}

}

func TestPlanet_innateScanner(t *testing.T) {
	player := NewPlayer(1, &Race{Spec: RaceSpec{InnateMining: false}})
	planet := Planet{}
	planet.setPopulation(67300)

	if got := planet.innateScanner(player); got != 0 {
		t.Errorf("Planet.GetInnateMines() = %v, want %v", got, 0)
	}

	// should get 40 mines for 16k pop when the player has innate mining
	player.Race.Spec.InnateScanner = true
	player.Race.Spec.InnatePopulationFactor = .1
	if got := planet.innateScanner(player); got != 82 {
		t.Errorf("Planet.GetInnateMines() = %v, want %v", got, 82)
	}

}

func TestPlanet_getGrowthAmount(t *testing.T) {
	type fields struct {
		Hab        Hab
		Population int
	}
	type args struct {
		player        *Player
		maxPopulation int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{name: "empty planet", fields: fields{Hab{50, 50, 50}, 0}, args: args{NewPlayer(1, NewRace()), 1_000_000}, want: 0},
		{name: "less than 25% cap, grows at full 10% growth rate", fields: fields{Hab{50, 50, 50}, 100_000}, args: args{NewPlayer(1, NewRace()), 1_200_000}, want: 10_000},
		{name: "at 50% cap, it slows down in growth", fields: fields{Hab{50, 50, 50}, 600_000}, args: args{NewPlayer(1, NewRace()), 1_200_000}, want: 26_700},
		{name: "we are basicallly at capacity, we only grow a tiny amount", fields: fields{Hab{50, 50, 50}, 1_180_000}, args: args{NewPlayer(1, NewRace()), 1_200_000}, want: 100},
		{name: "no more growth past a certain capacity", fields: fields{Hab{50, 50, 50}, 1_190_000}, args: args{NewPlayer(1, NewRace()), 1_200_000}, want: 0},
		{name: "hostile planets kill off colonists", fields: fields{Hab{10, 15, 15}, 2500}, args: args{NewPlayer(1, NewRace()), 0}, want: -100},
		{name: "super hostile planet with 100k people, should be -45% habitable, so should kill off -4.5% of the pop", fields: fields{Hab{}, 100_000}, args: args{NewPlayer(1, NewRace()), 0}, want: -4500},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Planet{
				Cargo: Cargo{Colonists: tt.fields.Population / 100},
				Hab:   tt.fields.Hab,
			}
			// 10% growth for easier math
			tt.args.player.Race.GrowthRate = 10
			tt.args.player.Race.Spec = computeRaceSpec(&tt.args.player.Race, &rules)
			if got := p.getGrowthAmount(tt.args.player, tt.args.maxPopulation); got != tt.want {
				t.Errorf("Planet.getGrowthAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueueItemType_IsAuto(t *testing.T) {
	tests := []struct {
		name string
		tr   QueueItemType
		want bool
	}{
		{"Not Auto", QueueItemTypeFactory, false},
		{"Not Auto", QueueItemTypeAutoMaxTerraform, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.IsAuto(); got != tt.want {
				t.Errorf("QueueItemType.IsAuto() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanet_reduceMineralConcentration(t *testing.T) {
	type args struct {
		rules Rules
	}
	tests := []struct {
		name   string
		planet *Planet
		args   args
		want   Mineral
	}{
		{"Redcue empty planet min conc", NewPlanet(), args{rules}, Mineral{1, 1, 1}},
		{
			"150 mines should reduce 100% conc by 1 if we have 151 mineyears",
			NewPlanet().
				WithMineralConcentration(Mineral{100, 100, 100}).
				WithMines(150).
				WithMineYears(Mineral{151, 151, 151}),
			args{rules},
			Mineral{99, 99, 99},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.planet.reduceMineralConcentration(&tt.args.rules)

			if got := tt.planet.MineralConcentration; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Planet.reduceMineralConcentration() = %v, want %v", got, tt.want)
			}

		})
	}
}

func Test_getMaxPopulation(t *testing.T) {
	type args struct {
		hab                 int
		maxPopulationOffset float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"plain homeworld", args{100, 0}, 1_000_000},
		{"joat homeworld", args{100, .2}, 1_200_000},
		{"low hab world", args{1, 0}, 50_000},
		{"bad hab world", args{-45, 0}, 50_000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMaxPopulation(&rules, tt.args.hab, tt.args.maxPopulationOffset); got != tt.want {
				t.Errorf("getMaxPopulation() = %v, want %v", got, tt.want)
			}
		})
	}
}
