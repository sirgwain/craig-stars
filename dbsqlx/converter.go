//go:generate go run github.com/jmattheis/goverter/cmd/goverter --packageName dbsqlx --output ./dbsqlx/generated.go --packagePath github.com/sirgwain/craig-stars/dbsqlx github.com/sirgwain/craigstars/dbsqlx
package dbsqlx

import (
	"time"

	"github.com/sirgwain/craig-stars/game"
)

// generate converter with

// goverter:converter
// goverter:extend TimeToTime
// goverter:name GameConverter
type Converter interface {
	ConvertUser(source User) game.User
	ConvertGameUser(source *game.User) *User
	ConvertUsers(source []User) []game.User

	// goverter:mapExtend ResearchCost ExtendResearchCost
	// goverter:mapExtend HabLow ExtendHabLow
	// goverter:mapExtend HabHigh ExtendHabHigh
	ConvertRace(source Race) game.Race
	// goverter:mapExtend ResearchCost ExtendResearchCost
	// goverter:mapExtend HabLow ExtendHabLow
	// goverter:mapExtend HabHigh ExtendHabHigh
	ConvertRaces(source []Race) []game.Race

	// goverter:map HabHigh.Grav HabHighGrav
	// goverter:map HabHigh.Temp HabHighTemp
	// goverter:map HabHigh.Rad HabHighRad
	// goverter:map HabLow.Grav HabLowGrav
	// goverter:map HabLow.Temp HabLowTemp
	// goverter:map HabLow.Rad HabLowRad
	// goverter:map ResearchCost.Energy ResearchCostEnergy
	// goverter:map ResearchCost.Weapons ResearchCostWeapons
	// goverter:map ResearchCost.Propulsion ResearchCostPropulsion
	// goverter:map ResearchCost.Construction ResearchCostConstruction
	// goverter:map ResearchCost.Electronics ResearchCostElectronics
	// goverter:map ResearchCost.Biotechnology ResearchCostBiotechnology
	ConvertGameRace(source *game.Race) *Race
}

func TimeToTime(t time.Time) time.Time {
	return t
}

func ExtendResearchCost(source Race) game.ResearchCost {
	return game.ResearchCost{
		Energy:        source.ResearchCostEnergy,
		Weapons:       source.ResearchCostWeapons,
		Propulsion:    source.ResearchCostPropulsion,
		Construction:  source.ResearchCostConstruction,
		Electronics:   source.ResearchCostElectronics,
		Biotechnology: source.ResearchCostBiotechnology,
	}
}

func ExtendHabLow(source Race) game.Hab {
	return game.Hab{
		Grav: source.HabLowGrav,
		Temp: source.HabLowTemp,
		Rad:  source.HabLowRad,
	}
}

func ExtendHabHigh(source Race) game.Hab {
	return game.Hab{
		Grav: source.HabHighGrav,
		Temp: source.HabHighTemp,
		Rad:  source.HabHighRad,
	}
}
