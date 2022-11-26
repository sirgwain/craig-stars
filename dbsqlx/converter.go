//go:generate go run github.com/jmattheis/goverter/cmd/goverter --packageName dbsqlx --output ./dbsqlx/generated.go --packagePath github.com/sirgwain/craig-stars/dbsqlx --ignoreUnexportedFields github.com/sirgwain/craigstars/dbsqlx
package dbsqlx

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/sirgwain/craig-stars/game"
)

// goverter:converter
// goverter:extend TimeToTime
// goverter:extend NullTimeToTime
// goverter:extend TimeToNullTime
// goverter:extend UUIDToUUID
// goverter:extend RulesToGameRules
// goverter:extend GameRulesToRules
// goverter:extend RaceSpecToGameRaceSpec
// goverter:extend GameRaceSpecToRaceSpec
// goverter:extend ProductionPlansToGameProductionPlans
// goverter:extend GameProductionPlansToProductionPlans
// goverter:extend BattlePlansToGameBattlePlans
// goverter:extend GameBattlePlansToBattlePlans
// goverter:extend TransportPlansToGameTransportPlans
// goverter:extend GameTransportPlansToTransportPlans
// goverter:extend PlayerMessagesToGamePlayerMessages
// goverter:extend GamePlayerMessagesToPlayerMessages
// goverter:extend PlanetIntelsToGamePlanetIntels
// goverter:extend GamePlanetIntelsToPlanetIntels
// goverter:extend FleetIntelsToGameFleetIntels
// goverter:extend GameFleetIntelsToFleetIntels
// goverter:extend ShipDesignIntelsToGameShipDesignIntels
// goverter:extend GameShipDesignIntelsToShipDesignIntels
// goverter:extend MineralPacketIntelsToGameMineralPacketIntels
// goverter:extend GameMineralPacketIntelsToMineralPacketIntels
// goverter:extend MineFieldIntelsToGameMineFieldIntels
// goverter:extend GameMineFieldIntelsToMineFieldIntels
// goverter:extend GameRaceToPlayerRace
// goverter:extend PlayerRaceToGameRace
// goverter:extend PlayerSpecToGamePlayerSpec
// goverter:extend GamePlayerSpecToPlayerSpec
// goverter:extend PlayerStatsToGamePlayerStats
// goverter:extend GamePlayerStatsToPlayerStats
// goverter:name GameConverter
// goverter:extend PlanetSpecToGamePlanetSpec
// goverter:extend GamePlanetSpecToPlanetSpec
// goverter:extend ProductionQueueItemsToGameProductionQueueItems
// goverter:extend GameProductionQueueItemsToProductionQueueItems
// goverter:extend FleetSpecToGameFleetSpec
// goverter:extend GameFleetSpecToFleetSpec
// goverter:extend WaypointsToGameWaypoints
// goverter:extend GameWaypointsToWaypoints
// goverter:extend ShipDesignSpecToGameShipDesignSpec
// goverter:extend GameShipDesignSpecToShipDesignSpec
// goverter:extend ShipDesignSlotsToGameShipDesignSlots
// goverter:extend GameShipDesignSlotsToShipDesignSlots
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
	// goverter:mapExtend Rules ExtendDefaultRules
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
	// goverter:ignore Designs
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

	// goverter:map MapObject.ID ID
	// goverter:map MapObject.GameID GameID
	// goverter:map MapObject.CreatedAt CreatedAt
	// goverter:map MapObject.UpdatedAt UpdatedAt
	// goverter:map MapObject.Type Type
	// goverter:map MapObject.PlayerID PlayerID
	// goverter:map MapObject.Dirty Dirty
	// goverter:map MapObject.Delete Delete
	// goverter:map MapObject.Position.X X
	// goverter:map MapObject.Position.Y Y
	// goverter:map MapObject.Name Name
	// goverter:map MapObject.Num Num
	// goverter:map MapObject.PlayerNum	 PlayerNum
	// goverter:ignore Tags
	// goverter:map Hab.Grav Grav
	// goverter:map Hab.Temp Temp
	// goverter:map Hab.Rad Rad
	// goverter:map BaseHab.Grav BaseGrav
	// goverter:map BaseHab.Temp BaseTemp
	// goverter:map BaseHab.Rad BaseRad
	// goverter:map TerraformedAmount.Grav TerraformedAmountGrav
	// goverter:map TerraformedAmount.Temp TerraformedAmountTemp
	// goverter:map TerraformedAmount.Rad TerraformedAmountRad
	// goverter:map MineralConcentration.Ironium MineralConcIronium
	// goverter:map MineralConcentration.Boranium MineralConcBoranium
	// goverter:map MineralConcentration.Germanium MineralConcGermanium
	// goverter:map MineYears.Ironium MineYearsIronium
	// goverter:map MineYears.Boranium MineYearsBoranium
	// goverter:map MineYears.Germanium MineYearsGermanium
	// goverter:map Cargo.Ironium Ironium
	// goverter:map Cargo.Boranium Boranium
	// goverter:map Cargo.Germanium Germanium
	// goverter:map Cargo.Colonists Colonists
	ConvertGamePlanet(source *game.Planet) *Planet

	// goverter:mapExtend Hab ExtendHab
	// goverter:mapExtend BaseHab ExtendBaseHab
	// goverter:mapExtend TerraformedAmount ExtendTerraformedAmount
	// goverter:mapExtend MineralConcentration ExtendMineralConcentration
	// goverter:mapExtend MineYears ExtendMineYears
	// goverter:mapExtend Cargo ExtendPlanetCargo
	// goverter:mapExtend MapObject ExtendPlanetMapObject
	ConvertPlanet(source *Planet) *game.Planet

	// goverter:map MapObject.ID ID
	// goverter:map MapObject.GameID GameID
	// goverter:map MapObject.CreatedAt CreatedAt
	// goverter:map MapObject.UpdatedAt UpdatedAt
	// goverter:map MapObject.Type Type
	// goverter:map MapObject.PlayerID PlayerID
	// goverter:map MapObject.Dirty Dirty
	// goverter:map MapObject.Delete Delete
	// goverter:map MapObject.Position.X X
	// goverter:map MapObject.Position.Y Y
	// goverter:map MapObject.Name Name
	// goverter:map MapObject.Num Num
	// goverter:map MapObject.PlayerNum	 PlayerNum
	// goverter:map FleetOrders.Waypoints Waypoints
	// goverter:map FleetOrders.RepeatOrders RepeatOrders
	// goverter:ignore Tags
	// goverter:map Heading.X HeadingX
	// goverter:map Heading.Y HeadingY
	// goverter:map PreviousPosition.X PreviousPositionX
	// goverter:map PreviousPosition.Y PreviousPositionY
	// goverter:map Cargo.Ironium Ironium
	// goverter:map Cargo.Boranium Boranium
	// goverter:map Cargo.Germanium Germanium
	// goverter:map Cargo.Colonists Colonists
	ConvertGameFleet(source *game.Fleet) *Fleet

	// goverter:mapExtend Heading ExtendFleetHeading
	// goverter:mapExtend PreviousPosition ExtendFleetPreviousPosition
	// goverter:mapExtend Cargo ExtendFleetCargo
	// goverter:mapExtend MapObject ExtendFleetMapObject
	// goverter:mapExtend FleetOrders ExtendFleetFleetOrders
	// goverter:ignore Tokens
	ConvertFleet(source *Fleet) *game.Fleet

	ConvertGameShipDesign(source *game.ShipDesign) *ShipDesign
	// goverter:ignore Dirty
	ConvertShipDesign(source *ShipDesign) *game.ShipDesign

	ConvertGameShipToken(source game.ShipToken) ShipToken
	// goverter:ignore Design
	ConvertShipToken(source ShipToken) game.ShipToken
}

func TimeToTime(source time.Time) time.Time {
	return source
}

func NullTimeToTime(source sql.NullTime) time.Time {
	if source.Valid {
		return source.Time
	}
	return time.Time{}
}

func TimeToNullTime(source time.Time) sql.NullTime {
	return sql.NullTime{
		Valid: true,
		Time:  source,
	}
}

func UUIDToUUID(source uuid.UUID) uuid.UUID {
	return source
}

func RulesToGameRules(source *Rules) game.Rules {
	return game.Rules(*source)
}

func GameRulesToRules(source game.Rules) *Rules {
	return (*Rules)(&source)
}

func RaceSpecToGameRaceSpec(source *RaceSpec) *game.RaceSpec {
	return (*game.RaceSpec)(source)
}

func GameRaceSpecToRaceSpec(source *game.RaceSpec) *RaceSpec {
	return (*RaceSpec)(source)
}

func BattlePlansToGameBattlePlans(source *BattlePlans) []game.BattlePlan {
	return ([]game.BattlePlan)(*source)
}

func GameBattlePlansToBattlePlans(source []game.BattlePlan) *BattlePlans {
	return (*BattlePlans)(&source)
}

func ProductionPlansToGameProductionPlans(source *ProductionPlans) []game.ProductionPlan {
	return ([]game.ProductionPlan)(*source)
}

func GameProductionPlansToProductionPlans(source []game.ProductionPlan) *ProductionPlans {
	return (*ProductionPlans)(&source)
}

func TransportPlansToGameTransportPlans(source *TransportPlans) []game.TransportPlan {
	return ([]game.TransportPlan)(*source)
}

func GameTransportPlansToTransportPlans(source []game.TransportPlan) *TransportPlans {
	return (*TransportPlans)(&source)
}

func PlayerMessagesToGamePlayerMessages(source *PlayerMessages) []game.PlayerMessage {
	return ([]game.PlayerMessage)(*source)
}

func GamePlayerMessagesToPlayerMessages(source []game.PlayerMessage) *PlayerMessages {
	return (*PlayerMessages)(&source)
}

func PlanetIntelsToGamePlanetIntels(source *PlanetIntels) []game.PlanetIntel {
	return ([]game.PlanetIntel)(*source)
}

func GamePlanetIntelsToPlanetIntels(source []game.PlanetIntel) *PlanetIntels {
	return (*PlanetIntels)(&source)
}

func FleetIntelsToGameFleetIntels(source *FleetIntels) []game.FleetIntel {
	return ([]game.FleetIntel)(*source)
}

func GameFleetIntelsToFleetIntels(source []game.FleetIntel) *FleetIntels {
	return (*FleetIntels)(&source)
}

func ShipDesignIntelsToGameShipDesignIntels(source *ShipDesignIntels) []game.ShipDesignIntel {
	return ([]game.ShipDesignIntel)(*source)
}

func GameShipDesignIntelsToShipDesignIntels(source []game.ShipDesignIntel) *ShipDesignIntels {
	return (*ShipDesignIntels)(&source)
}

func MineralPacketIntelsToGameMineralPacketIntels(source *MineralPacketIntels) []game.MineralPacketIntel {
	return ([]game.MineralPacketIntel)(*source)
}

func GameMineralPacketIntelsToMineralPacketIntels(source []game.MineralPacketIntel) *MineralPacketIntels {
	return (*MineralPacketIntels)(&source)
}

func MineFieldIntelsToGameMineFieldIntels(source *MineFieldIntels) []game.MineFieldIntel {
	return ([]game.MineFieldIntel)(*source)
}

func GameMineFieldIntelsToMineFieldIntels(source []game.MineFieldIntel) *MineFieldIntels {
	return (*MineFieldIntels)(&source)
}


func PlayerRaceToGameRace(source *PlayerRace) game.Race {
	return game.Race(*source)
}

func GameRaceToPlayerRace(source game.Race) *PlayerRace {
	return (*PlayerRace)(&source)
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

func PlanetSpecToGamePlanetSpec(source *PlanetSpec) *game.PlanetSpec {
	return (*game.PlanetSpec)(source)
}

func GamePlanetSpecToPlanetSpec(source *game.PlanetSpec) *PlanetSpec {
	return (*PlanetSpec)(source)
}

func ProductionQueueItemsToGameProductionQueueItems(source *ProductionQueueItems) []game.ProductionQueueItem {
	return ([]game.ProductionQueueItem)(*source)
}

func GameProductionQueueItemsToProductionQueueItems(source []game.ProductionQueueItem) *ProductionQueueItems {
	return (*ProductionQueueItems)(&source)
}

func FleetSpecToGameFleetSpec(source *FleetSpec) *game.FleetSpec {
	return (*game.FleetSpec)(source)
}

func GameFleetSpecToFleetSpec(source *game.FleetSpec) *FleetSpec {
	return (*FleetSpec)(source)
}

func WaypointsToGameWaypoints(source *Waypoints) []game.Waypoint {
	return ([]game.Waypoint)(*source)
}

func GameWaypointsToWaypoints(source []game.Waypoint) *Waypoints {
	return (*Waypoints)(&source)
}

func ShipDesignSpecToGameShipDesignSpec(source *ShipDesignSpec) *game.ShipDesignSpec {
	return (*game.ShipDesignSpec)(source)
}

func GameShipDesignSpecToShipDesignSpec(source *game.ShipDesignSpec) *ShipDesignSpec {
	return (*ShipDesignSpec)(source)
}

func ShipDesignSlotsToGameShipDesignSlots(source *ShipDesignSlots) []game.ShipDesignSlot {
	return ([]game.ShipDesignSlot)(*source)
}

func GameShipDesignSlotsToShipDesignSlots(source []game.ShipDesignSlot) *ShipDesignSlots {
	return (*ShipDesignSlots)(&source)
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
		Conditions:               *source.VictoryConditionsConditions,
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

func ExtendDefaultRules(source Game) game.Rules {
	return game.NewRules()
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

func ExtendPlanetMapObject(source Planet) game.MapObject {
	return game.MapObject{
		Type:      game.MapObjectTypePlanet,
		ID:        source.ID,
		GameID:    source.GameID,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
		PlayerID:  source.PlayerID,
		Position: game.Vector{
			X: source.X,
			Y: source.Y,
		},
		Name:      source.Name,
		Num:       source.Num,
		PlayerNum: source.PlayerNum,
		// Tags:      source.Tags,
	}
}

func ExtendHab(source Planet) game.Hab {
	return game.Hab{
		Grav: source.Grav,
		Temp: source.Temp,
		Rad:  source.Rad,
	}
}

func ExtendBaseHab(source Planet) game.Hab {
	return game.Hab{
		Grav: source.BaseGrav,
		Temp: source.BaseTemp,
		Rad:  source.BaseRad,
	}
}

func ExtendTerraformedAmount(source Planet) game.Hab {
	return game.Hab{
		Grav: source.TerraformedAmountGrav,
		Temp: source.TerraformedAmountTemp,
		Rad:  source.TerraformedAmountRad,
	}
}

func ExtendMineralConcentration(source Planet) game.Mineral {
	return game.Mineral{
		Ironium:   source.MineralConcIronium,
		Boranium:  source.MineralConcBoranium,
		Germanium: source.MineralConcGermanium,
	}
}

func ExtendMineYears(source Planet) game.Mineral {
	return game.Mineral{
		Ironium:   source.MineYearsIronium,
		Boranium:  source.MineYearsBoranium,
		Germanium: source.MineYearsGermanium,
	}
}

func ExtendPlanetCargo(source Planet) game.Cargo {
	return game.Cargo{
		Ironium:   source.Ironium,
		Boranium:  source.Boranium,
		Germanium: source.Germanium,
		Colonists: source.Colonists,
	}
}

func ExtendFleetMapObject(source Fleet) game.MapObject {
	return game.MapObject{
		Type:      game.MapObjectTypeFleet,
		ID:        source.ID,
		GameID:    source.GameID,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
		PlayerID:  source.PlayerID,
		Position: game.Vector{
			X: source.X,
			Y: source.Y,
		},
		Name:      source.Name,
		Num:       source.Num,
		PlayerNum: source.PlayerNum,
		// Tags:      source.Tags,
	}
}

func ExtendFleetFleetOrders(source Fleet) game.FleetOrders {
	return game.FleetOrders{
		Waypoints:    *source.Waypoints,
		RepeatOrders: source.RepeatOrders,
	}
}

func ExtendFleetCargo(source Fleet) game.Cargo {
	return game.Cargo{
		Ironium:   source.Ironium,
		Boranium:  source.Boranium,
		Germanium: source.Germanium,
		Colonists: source.Colonists,
	}
}

func ExtendFleetHeading(source Fleet) game.Vector {
	return game.Vector{
		X: source.HeadingX,
		Y: source.HeadingY,
	}
}

func ExtendFleetPreviousPosition(source Fleet) *game.Vector {
	if source.PreviousPositionX == nil || source.PreviousPositionY == nil {
		return nil
	}
	return &game.Vector{
		X: *source.PreviousPositionX,
		Y: *source.PreviousPositionY,
	}
}

