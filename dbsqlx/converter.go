//go:generate go run github.com/jmattheis/goverter/cmd/goverter --packageName dbsqlx --output ./dbsqlx/generated.go --packagePath github.com/sirgwain/craig-stars/dbsqlx github.com/sirgwain/craigstars/dbsqlx
package dbsqlx

import (
	"time"

	"github.com/sirgwain/craig-stars/game"
)

// generate converter with

// goverter:converter
// goverter:extend TimeToTime
// goverter:extend RulesToGameRules
// goverter:extend GameRulesToRules
// goverter:extend RaceSpecToGameRaceSpec
// goverter:extend GameRaceSpecToRaceSpec
// goverter:extend RaceSpecToRaceSpec
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

	// goverter:mapExtend VictoryConditions ExtendVictoryConditions
	// goverter:mapExtend Area ExtendArea
	ConvertGame(source Game) game.Game

	// goverter:mapExtend VictoryConditions ExtendVictoryConditions
	// goverter:mapExtend Area ExtendArea
	ConvertGames(source []Game) []game.Game

	// goverter:map VictoryConditions.NumCriteriaRequired VictoryConditionsNumCriteriaRequired
	// goverter:map VictoryConditions.YearsPassed VictoryConditionsYearsPassed
	// goverter:map VictoryConditions.OwnPlanets VictoryConditionsOwnPlanets
	// goverter:map VictoryConditions.AttainTechLevel VictoryConditionsAttainTechLevel
	// goverter:map VictoryConditions.AttainTechLevelNumFields VictoryConditionsAttainTechLevelNumFields
	// goverter:map VictoryConditions.ExceedsScore VictoryConditionsExceedsScore
	// goverter:map VictoryConditions.ExceedsSecondPlaceScore VictoryConditionsExceedsSecondPlaceScore
	// goverter:map VictoryConditions.ProductionCapacity VictoryConditionsProductionCapacity
	// goverter:map VictoryConditions.OwnCapitalShips VictoryConditionsOwnCapitalShips
	// goverter:map VictoryConditions.HighestScoreAfterYears VictoryConditionsHighestScoreAfterYears
	// goverter:map VictoryConditions.Conditions VictoryConditionsConditions
	// goverter:map Area.X AreaX
	// goverter:map Area.Y AreaY
	ConvertGameGame(source *game.Game) *Game
}

func TimeToTime(t time.Time) time.Time {
	return t
}

func RulesToGameRules(r Rules) game.Rules {
	return game.Rules(r)
}

func GameRulesToRules(r game.Rules) Rules {
	return Rules(r)
}

func RaceSpecToGameRaceSpec(r *RaceSpec) *game.RaceSpec {
	return (*game.RaceSpec)(r)
}

func GameRaceSpecToRaceSpec(r *game.RaceSpec) *RaceSpec {
	return (*RaceSpec)(r)
}

func RaceSpecToRaceSpec(s *game.RaceSpec) *game.RaceSpec {
	return s
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

func ExtendVictoryConditions(source Game) game.VictoryConditions {
	return game.VictoryConditions{
		Conditions:               source.VictoryConditionsConditions,
		NumCriteriaRequired:      source.VictoryConditionsNumCriteriaRequired,
		YearsPassed:              source.VictoryConditionsYearsPassed,
		OwnPlanets:               source.VictoryConditionsOwnPlanets,
		AttainTechLevel:          source.VictoryConditionsAttainTechLevel,
		AttainTechLevelNumFields: source.VictoryConditionsAttainTechLevelNumFields,
		ExceedsScore:             source.VictoryConditionsExceedsScore,
		ExceedsSecondPlaceScore:  source.VictoryConditionsExceedsSecondPlaceScore,
		ProductionCapacity:       source.VictoryConditionsProductionCapacity,
		OwnCapitalShips:          source.VictoryConditionsOwnCapitalShips,
		HighestScoreAfterYears:   source.VictoryConditionsHighestScoreAfterYears,
	}
}

func ExtendArea(source Game) game.Vector {
	return game.Vector{
		X: source.AreaX,
		Y: source.AreaY,
	}
}
