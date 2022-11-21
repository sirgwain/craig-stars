//go:generate go run github.com/jmattheis/goverter/cmd/goverter --packageName dbsqlx --output ./dbsqlx/generated.go --packagePath github.com/sirgwain/craig-stars/dbsqlx --ignoreUnexportedFields github.com/sirgwain/craigstars/dbsqlx
package dbsqlx

import (
	"time"

	"github.com/sirgwain/craig-stars/game"
)

// goverter:converter
// goverter:extend TimeToTime
// goverter:extend RulesToGameRules
// goverter:extend GameRulesToRules
// goverter:extend RaceSpecToGameRaceSpec
// goverter:extend GameRaceSpecToRaceSpec
// goverter:extend ProductionPlansToGameProductionPlans
// goverter:extend GameProductionPlansToProductionPlans
// goverter:extend TransportPlansToGameTransportPlans
// goverter:extend GameTransportPlansToTransportPlans
// goverter:extend PlayerSpecToGamePlayerSpec
// goverter:extend GamePlayerSpecToPlayerSpec
// goverter:extend PlayerStatsToGamePlayerStats
// goverter:extend GamePlayerStatsToPlayerStats
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

	// goverter:mapExtend TechLevels ExtendTechLevels
	// goverter:mapExtend TechLevelsSpent ExtendTechLevelsSpent
	// goverter:ignore Race
	// goverter:ignore Messages
	// goverter:ignore BattlePlans
	// goverter:ignore Designs
	// goverter:ignore PlanetIntels
	// goverter:ignore FleetIntels
	// goverter:ignore DesignIntels
	// goverter:ignore MineralPacketIntels
	// goverter:ignore MineFieldIntels
	ConvertPlayer(source Player) game.Player
	// goverter:mapExtend TechLevels ExtendTechLevels
	// goverter:mapExtend TechLevelsSpent ExtendTechLevelsSpent
	ConvertPlayers(source []Player) []game.Player

	// goverter:map TechLevels.Energy TechLevelsEnergy
	// goverter:map TechLevels.Weapons TechLevelsWeapons
	// goverter:map TechLevels.Propulsion TechLevelsPropulsion
	// goverter:map TechLevels.Construction TechLevelsConstruction
	// goverter:map TechLevels.Electronics TechLevelsElectronics
	// goverter:map TechLevels.Biotechnology TechLevelsBiotechnology
	// goverter:map TechLevelsSpent.Energy TechLevelsSpentEnergy
	// goverter:map TechLevelsSpent.Weapons TechLevelsSpentWeapons
	// goverter:map TechLevelsSpent.Propulsion TechLevelsSpentPropulsion
	// goverter:map TechLevelsSpent.Construction TechLevelsSpentConstruction
	// goverter:map TechLevelsSpent.Electronics TechLevelsSpentElectronics
	// goverter:map TechLevelsSpent.Biotechnology TechLevelsSpentBiotechnology
	ConvertGamePlayer(source *game.Player) *Player
}

func TimeToTime(source time.Time) time.Time {
	return source
}

func RulesToGameRules(source Rules) game.Rules {
	return game.Rules(source)
}

func GameRulesToRules(source game.Rules) Rules {
	return Rules(source)
}

func RaceSpecToGameRaceSpec(source *RaceSpec) *game.RaceSpec {
	return (*game.RaceSpec)(source)
}

func GameRaceSpecToRaceSpec(source *game.RaceSpec) *RaceSpec {
	return (*RaceSpec)(source)
}

func ProductionPlansToGameProductionPlans(source ProductionPlans) []game.ProductionPlan {
	return ([]game.ProductionPlan)(source)
}

func GameProductionPlansToProductionPlans(source []game.ProductionPlan) ProductionPlans {
	return (ProductionPlans)(source)
}

func TransportPlansToGameTransportPlans(source TransportPlans) []game.TransportPlan {
	return ([]game.TransportPlan)(source)
}

func GameTransportPlansToTransportPlans(source []game.TransportPlan) TransportPlans {
	return (TransportPlans)(source)
}

func PlayerSpecToGamePlayerSpec(source *PlayerSpec) *game.PlayerSpec {
	return (*game.PlayerSpec)(source)
}

func GamePlayerSpecToPlayerSpec(source *game.PlayerSpec) *PlayerSpec {
	return (*PlayerSpec)(source)
}

func PlayerStatsToGamePlayerStats(source *PlayerStats) *game.PlayerStats {
	return (*game.PlayerStats)(source)
}

func GamePlayerStatsToPlayerStats(source *game.PlayerStats) *PlayerStats {
	return (*PlayerStats)(source)
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

func ExtendTechLevels(source Player) game.TechLevel {
	return game.TechLevel{
		Energy:        source.TechLevelsEnergy,
		Weapons:       source.TechLevelsWeapons,
		Propulsion:    source.TechLevelsPropulsion,
		Construction:  source.TechLevelsConstruction,
		Electronics:   source.TechLevelsElectronics,
		Biotechnology: source.TechLevelsBiotechnology,
	}
}

func ExtendTechLevelsSpent(source Player) game.TechLevel {
	return game.TechLevel{
		Energy:        source.TechLevelsSpentEnergy,
		Weapons:       source.TechLevelsSpentWeapons,
		Propulsion:    source.TechLevelsSpentPropulsion,
		Construction:  source.TechLevelsSpentConstruction,
		Electronics:   source.TechLevelsSpentElectronics,
		Biotechnology: source.TechLevelsSpentBiotechnology,
	}
}
