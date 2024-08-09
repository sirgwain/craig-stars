//go:generate go run github.com/jmattheis/goverter/cmd/goverter@v1.4.0 gen ./

package db

import (
	"database/sql"
	"time"

	"github.com/sirgwain/craig-stars/cs"
)

// will be instanciated in ./converter.init.go
var c Converter

// goverter:converter
// goverter:output:package github.com/sirgwain/craig-stars/db
// goverter:output:file ./generated.go
// goverter:ignoreUnexported
// goverter:extend TimeToTime
// goverter:extend NullTimeToTime
// goverter:extend TimeToNullTime
// goverter:extend NullBoolToBool
// goverter:extend BoolToNullBool
// goverter:extend RulesToGameRules
// goverter:extend GameRulesToRules
// goverter:extend TagsToGameTags
// goverter:extend GameTagsToTags
// goverter:extend RaceSpecToGameRaceSpec
// goverter:extend GameRaceSpecToRaceSpec
// goverter:extend ProductionPlansToGameProductionPlans
// goverter:extend GameProductionPlansToProductionPlans
// goverter:extend BattlePlansToGameBattlePlans
// goverter:extend GameBattlePlansToBattlePlans
// goverter:extend TransportPlansToGameTransportPlans
// goverter:extend GameTransportPlansToTransportPlans
// goverter:extend PlayerRelationshipsToGamePlayerRelationships
// goverter:extend GamePlayerRelationshipsToPlayerRelationships
// goverter:extend PlayerMessagesToGamePlayerMessages
// goverter:extend GamePlayerMessagesToPlayerMessages
// goverter:extend PlayerScoresToGamePlayerScores
// goverter:extend GamePlayerScoresToPlayerScores
// goverter:extend AcquiredTechsToGameAcquiredTechs
// goverter:extend GameAcquiredTechsToAcquiredTechs
// goverter:extend BattleRecordsToGameBattleRecords
// goverter:extend GameBattleRecordsToBattleRecords
// goverter:extend PlayerIntelsToGamePlayerIntels
// goverter:extend GamePlayerIntelsToPlayerIntels
// goverter:extend ScoreIntelsToGameScoreIntels
// goverter:extend GameScoreIntelsToScoreIntels
// goverter:extend PlanetIntelsToGamePlanetIntels
// goverter:extend GamePlanetIntelsToPlanetIntels
// goverter:extend FleetIntelsToGameFleetIntels
// goverter:extend GameFleetIntelsToFleetIntels
// goverter:extend ShipDesignIntelsToGameShipDesignIntels
// goverter:extend GameShipDesignIntelsToShipDesignIntels
// goverter:extend MineralPacketIntelsToGameMineralPacketIntels
// goverter:extend GameMineralPacketIntelsToMineralPacketIntels
// goverter:extend SalvageIntelsToGameSalvageIntels
// goverter:extend GameSalvageIntelsToSalvageIntels
// goverter:extend MineFieldIntelsToGameMineFieldIntels
// goverter:extend GameMineFieldIntelsToMineFieldIntels
// goverter:extend WormholeIntelsToGameWormholeIntels
// goverter:extend GameWormholeIntelsToWormholeIntels
// goverter:extend MysteryTraderIntelsToGameMysteryTraderIntels
// goverter:extend GameMysteryTraderIntelsToMysteryTraderIntels
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
// goverter:extend ShipTokensToGameShipTokens
// goverter:extend GameShipTokensToShipTokens
// goverter:extend WaypointsToGameWaypoints
// goverter:extend GameWaypointsToWaypoints
// goverter:extend MineFieldSpecToGameMineFieldSpec
// goverter:extend GameMineFieldSpecToMineFieldSpec
// goverter:extend ShipDesignSpecToGameShipDesignSpec
// goverter:extend GameShipDesignSpecToShipDesignSpec
// goverter:extend ShipDesignSlotsToGameShipDesignSlots
// goverter:extend GameShipDesignSlotsToShipDesignSlots
// goverter:extend WormholeSpecToGameWormholeSpec
// goverter:extend GameWormholeSpecToWormholeSpec
// goverter:extend MysteryTraderSpecToGameMysteryTraderSpec
// goverter:extend GameMysteryTraderSpecToMysteryTraderSpec
// goverter:extend MysteryTraderRewardTypeToGameMysteryTraderRewardType
// goverter:extend GameMysteryTraderRewardTypeToMysteryTraderRewardType
// goverter:enum no
type Converter interface {
	// goverter:map . DBObject
	ConvertUser(source User) cs.User

	ConvertUsers(source []User) []cs.User

	// goverter:autoMap DBObject
	ConvertGameUser(source *cs.User) *User

	// goverter:map . DBObject
	// goverter:map . ResearchCost | ExtendResearchCost
	// goverter:map . HabLow | ExtendHabLow
	// goverter:map . HabHigh | ExtendHabHigh
	ConvertRace(source Race) cs.Race

	ConvertRaces(source []Race) []cs.Race

	// goverter:autoMap DBObject
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
	ConvertGameRace(source *cs.Race) *Race

	// goverter:map . DBObject
	// goverter:map . VictoryConditions | ExtendVictoryConditions
	// goverter:map . Area | ExtendArea
	// goverter:map . Rules | ExtendDefaultRules
	ConvertGame(source Game) cs.Game

	ConvertGames(source []Game) []cs.Game

	// goverter:autoMap DBObject
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
	ConvertGameGame(source *cs.Game) *Game

	// goverter:autoMap GameDBObject
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
	// goverter:autoMap PlayerOrders
	// goverter:map PlayerIntels.BattleRecords BattleRecords
	// goverter:map PlayerIntels.PlayerIntels PlayerIntels
	// goverter:map PlayerIntels.ScoreIntels ScoreIntels
	// goverter:map PlayerIntels.PlanetIntels PlanetIntels
	// goverter:map PlayerIntels.FleetIntels FleetIntels
	// goverter:map PlayerIntels.StarbaseIntels StarbaseIntels
	// goverter:map PlayerIntels.ShipDesignIntels ShipDesignIntels
	// goverter:map PlayerIntels.MineralPacketIntels MineralPacketIntels
	// goverter:map PlayerIntels.MineFieldIntels MineFieldIntels
	// goverter:map PlayerIntels.WormholeIntels WormholeIntels
	// goverter:map PlayerIntels.MysteryTraderIntels MysteryTraderIntels
	// goverter:map PlayerIntels.SalvageIntels	 SalvageIntels
	// goverter:autoMap PlayerPlans
	ConvertGamePlayer(source *cs.Player) *Player

	// goverter:map . GameDBObject
	// goverter:map . TechLevels | ExtendTechLevels
	// goverter:map . TechLevelsSpent | ExtendTechLevelsSpent
	// goverter:map . PlayerOrders
	// goverter:map . PlayerIntels
	// goverter:map . PlayerPlans
	// goverter:ignore Designs
	ConvertPlayer(source Player) cs.Player

	ConvertPlayers(source []Player) []cs.Player

	// goverter:autoMap MapObject.GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:autoMap PlanetOrders
	// goverter:autoMap Hab
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
	// goverter:autoMap Cargo
	ConvertGamePlanet(source *cs.Planet) *Planet

	// goverter:map . Hab
	// goverter:map . BaseHab | ExtendBaseHab
	// goverter:map . TerraformedAmount | ExtendTerraformedAmount
	// goverter:map . MineralConcentration | ExtendMineralConcentration
	// goverter:map . MineYears | ExtendMineYears
	// goverter:map . Cargo
	// goverter:map . MapObject | ExtendPlanetMapObject
	// goverter:map . PlanetOrders
	// goverter:ignore Starbase
	// goverter:ignore Dirty
	ConvertPlanet(source *Planet) *cs.Planet

	// goverter:autoMap MapObject.GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:autoMap FleetOrders
	// goverter:map Heading.X HeadingX
	// goverter:map Heading.Y HeadingY
	// goverter:map PreviousPosition.X PreviousPositionX
	// goverter:map PreviousPosition.Y PreviousPositionY
	// goverter:autoMap Cargo
	ConvertGameFleet(source *cs.Fleet) *Fleet

	// goverter:map . Heading | ExtendFleetHeading
	// goverter:map . PreviousPosition | ExtendFleetPreviousPosition
	// goverter:map . Cargo
	// goverter:map . MapObject | ExtendFleetMapObject
	// goverter:map . FleetOrders | ExtendFleetFleetOrders
	ConvertFleet(source *Fleet) *cs.Fleet

	// goverter:autoMap GameDBObject
	// goverter:ignore CanDelete
	ConvertGameShipDesign(source *cs.ShipDesign) *ShipDesign
	// goverter:ignore Delete
	// goverter:map . GameDBObject
	ConvertShipDesign(source *ShipDesign) *cs.ShipDesign

	// goverter:autoMap MapObject.GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	ConvertGameWormhole(source *cs.Wormhole) *Wormhole

	// goverter:map . MapObject
	ConvertWormhole(source *Wormhole) *cs.Wormhole

	// goverter:map . GameDBObject
	// goverter:map . Position
	// goverter:map Type | MapObjectTypeWormhole
	// goverter:ignore Delete
	// goverter:ignore PlayerNum
	wormHoleMapObject(source Wormhole) cs.MapObject

	// goverter:autoMap MapObject.GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:map Heading.X HeadingX
	// goverter:map Heading.Y HeadingY
	// goverter:map Destination.X DestinationX
	// goverter:map Destination.Y DestinationY
	ConvertGameMysteryTrader(source *cs.MysteryTrader) *MysteryTrader

	// goverter:map . MapObject | ExtendMysteryTraderMapObject
	// goverter:map . Heading | ExtendMysteryTraderHeading
	// goverter:map . Destination | ExtendMysteryTraderDestination
	ConvertMysteryTrader(source *MysteryTrader) *cs.MysteryTrader

	// goverter:autoMap MapObject.GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:autoMap Cargo
	ConvertGameSalvage(source *cs.Salvage) *Salvage

	// goverter:map . MapObject | ExtendSalvageMapObject
	// goverter:map . Cargo
	ConvertSalvage(source *Salvage) *cs.Salvage

	// goverter:autoMap MapObject.GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:map MineFieldOrders.Detonate Detonate
	ConvertGameMineField(source *cs.MineField) *MineField

	// goverter:map . MapObject | ExtendMineFieldMapObject
	// goverter:map . MineFieldOrders
	ConvertMineField(source *MineField) *cs.MineField

	// goverter:autoMap MapObject.GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:autoMap Cargo
	// goverter:map Cargo.Ironium Ironium
	// goverter:map Cargo.Boranium Boranium
	// goverter:map Cargo.Germanium Germanium
	// goverter:map Heading.X HeadingX
	// goverter:map Heading.Y HeadingY
	ConvertGameMineralPacket(source *cs.MineralPacket) *MineralPacket

	// goverter:map . MapObject | ExtendMineralPacketMapObject
	// goverter:map . Cargo
	// goverter:map . Heading | ExtendMineralPacketHeading
	ConvertMineralPacket(source *MineralPacket) *cs.MineralPacket

	// goverter:ignore Colonists
	salvageCargo(source Salvage) cs.Cargo

	// goverter:ignore Colonists
	mineralPaketCargo(source MineralPacket) cs.Cargo
}

func MapObjectTypeWormhole() cs.MapObjectType {
	return cs.MapObjectTypeWormhole
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

func NullBoolToBool(source sql.NullBool) bool {
	if source.Valid {
		return source.Bool
	}
	return false
}

func BoolToNullBool(source bool) sql.NullBool {
	return sql.NullBool{
		Valid: true,
		Bool:  source,
	}
}

func RulesToGameRules(source *Rules) cs.Rules {
	return cs.Rules(*source)
}

func GameRulesToRules(source cs.Rules) *Rules {
	return (*Rules)(&source)
}

func RaceSpecToGameRaceSpec(source *RaceSpec) cs.RaceSpec {
	return (cs.RaceSpec)(*source)
}

func GameRaceSpecToRaceSpec(source cs.RaceSpec) *RaceSpec {
	return (*RaceSpec)(&source)
}

func TagsToGameTags(source *Tags) cs.Tags {
	if source == nil {
		return cs.Tags{}
	}
	return (cs.Tags)(*source)
}

func GameTagsToTags(source cs.Tags) *Tags {
	return (*Tags)(&source)
}

func BattlePlansToGameBattlePlans(source *BattlePlans) []cs.BattlePlan {
	// return an empty slice for nil
	if source == nil {
		return []cs.BattlePlan{}
	}

	return ([]cs.BattlePlan)(*source)
}

func GameBattlePlansToBattlePlans(source []cs.BattlePlan) *BattlePlans {
	return (*BattlePlans)(&source)
}

func ProductionPlansToGameProductionPlans(source *ProductionPlans) []cs.ProductionPlan {
	// return an empty slice for nil
	if source == nil {
		return []cs.ProductionPlan{}
	}
	return ([]cs.ProductionPlan)(*source)
}

func GameProductionPlansToProductionPlans(source []cs.ProductionPlan) *ProductionPlans {
	return (*ProductionPlans)(&source)
}

func TransportPlansToGameTransportPlans(source *TransportPlans) []cs.TransportPlan {
	// return an empty slice for nil
	if source == nil {
		return []cs.TransportPlan{}
	}
	return ([]cs.TransportPlan)(*source)
}

func GameTransportPlansToTransportPlans(source []cs.TransportPlan) *TransportPlans {
	return (*TransportPlans)(&source)
}

func PlayerRelationshipsToGamePlayerRelationships(source *PlayerRelationships) []cs.PlayerRelationship {
	// return an empty slice for nil
	if source == nil {
		return []cs.PlayerRelationship{}
	}
	return ([]cs.PlayerRelationship)(*source)
}

func GamePlayerRelationshipsToPlayerRelationships(source []cs.PlayerRelationship) *PlayerRelationships {
	return (*PlayerRelationships)(&source)
}

func PlayerMessagesToGamePlayerMessages(source *PlayerMessages) []cs.PlayerMessage {
	// return an empty slice for nil
	if source == nil {
		return []cs.PlayerMessage{}
	}
	return ([]cs.PlayerMessage)(*source)
}

func GamePlayerMessagesToPlayerMessages(source []cs.PlayerMessage) *PlayerMessages {
	return (*PlayerMessages)(&source)
}

func PlayerScoresToGamePlayerScores(source *PlayerScores) []cs.PlayerScore {
	// return an empty slice for nil
	if source == nil {
		return []cs.PlayerScore{}
	}
	return ([]cs.PlayerScore)(*source)
}

func GamePlayerScoresToPlayerScores(source []cs.PlayerScore) *PlayerScores {
	return (*PlayerScores)(&source)
}

func AcquiredTechsToGameAcquiredTechs(source *AcquiredTechs) []string {
	// return an empty slice for nil
	if source == nil {
		return []string{}
	}
	return ([]string)(*source)
}

func GameAcquiredTechsToAcquiredTechs(source []string) *AcquiredTechs {
	return (*AcquiredTechs)(&source)
}

func BattleRecordsToGameBattleRecords(source *BattleRecords) []cs.BattleRecord {
	// return an empty slice for nil
	if source == nil {
		return []cs.BattleRecord{}
	}
	return ([]cs.BattleRecord)(*source)
}

func GameBattleRecordsToBattleRecords(source []cs.BattleRecord) *BattleRecords {
	return (*BattleRecords)(&source)
}

func PlayerIntelsToGamePlayerIntels(source *PlayerIntels) []cs.PlayerIntel {
	// return an empty slice for nil
	if source == nil {
		return []cs.PlayerIntel{}
	}
	return ([]cs.PlayerIntel)(*source)
}

func GamePlayerIntelsToPlayerIntels(source []cs.PlayerIntel) *PlayerIntels {
	return (*PlayerIntels)(&source)
}

func ScoreIntelsToGameScoreIntels(source *ScoreIntels) []cs.ScoreIntel {
	// return an empty slice for nil
	if source == nil {
		return []cs.ScoreIntel{}
	}
	return ([]cs.ScoreIntel)(*source)
}

func GameScoreIntelsToScoreIntels(source []cs.ScoreIntel) *ScoreIntels {
	return (*ScoreIntels)(&source)
}

func PlanetIntelsToGamePlanetIntels(source *PlanetIntels) []cs.PlanetIntel {
	// return an empty slice for nil
	if source == nil {
		return []cs.PlanetIntel{}
	}
	return ([]cs.PlanetIntel)(*source)
}

func GamePlanetIntelsToPlanetIntels(source []cs.PlanetIntel) *PlanetIntels {
	return (*PlanetIntels)(&source)
}

func FleetIntelsToGameFleetIntels(source *FleetIntels) []cs.FleetIntel {
	// return an empty slice for nil
	if source == nil {
		return []cs.FleetIntel{}
	}
	return ([]cs.FleetIntel)(*source)
}

func GameFleetIntelsToFleetIntels(source []cs.FleetIntel) *FleetIntels {
	return (*FleetIntels)(&source)
}

func ShipDesignIntelsToGameShipDesignIntels(source *ShipDesignIntels) []cs.ShipDesignIntel {
	// return an empty slice for nil
	if source == nil {
		return []cs.ShipDesignIntel{}
	}
	return ([]cs.ShipDesignIntel)(*source)
}

func GameShipDesignIntelsToShipDesignIntels(source []cs.ShipDesignIntel) *ShipDesignIntels {
	return (*ShipDesignIntels)(&source)
}

func MineralPacketIntelsToGameMineralPacketIntels(source *MineralPacketIntels) []cs.MineralPacketIntel {
	// return an empty slice for nil
	if source == nil {
		return []cs.MineralPacketIntel{}
	}
	return ([]cs.MineralPacketIntel)(*source)
}

func GameMineralPacketIntelsToMineralPacketIntels(source []cs.MineralPacketIntel) *MineralPacketIntels {
	return (*MineralPacketIntels)(&source)
}

func SalvageIntelsToGameSalvageIntels(source *SalvageIntels) []cs.SalvageIntel {
	// return an empty slice for nil
	if source == nil {
		return []cs.SalvageIntel{}
	}
	return ([]cs.SalvageIntel)(*source)
}

func GameSalvageIntelsToSalvageIntels(source []cs.SalvageIntel) *SalvageIntels {
	return (*SalvageIntels)(&source)
}

func MineFieldIntelsToGameMineFieldIntels(source *MineFieldIntels) []cs.MineFieldIntel {
	// return an empty slice for nil
	if source == nil {
		return []cs.MineFieldIntel{}
	}
	return ([]cs.MineFieldIntel)(*source)
}

func GameMineFieldIntelsToMineFieldIntels(source []cs.MineFieldIntel) *MineFieldIntels {
	return (*MineFieldIntels)(&source)
}

func WormholeIntelsToGameWormholeIntels(source *WormholeIntels) []cs.WormholeIntel {
	// return an empty slice for nil
	if source == nil {
		return []cs.WormholeIntel{}
	}
	return ([]cs.WormholeIntel)(*source)
}

func GameWormholeIntelsToWormholeIntels(source []cs.WormholeIntel) *WormholeIntels {
	return (*WormholeIntels)(&source)
}

func MysteryTraderIntelsToGameMysteryTraderIntels(source *MysteryTraderIntels) []cs.MysteryTraderIntel {
	// return an empty slice for nil
	if source == nil {
		return []cs.MysteryTraderIntel{}
	}
	return ([]cs.MysteryTraderIntel)(*source)
}

func GameMysteryTraderIntelsToMysteryTraderIntels(source []cs.MysteryTraderIntel) *MysteryTraderIntels {
	return (*MysteryTraderIntels)(&source)
}

func PlayerRaceToGameRace(source *PlayerRace) cs.Race {
	// return an empty object for nil to support partial loads
	if source == nil {
		return cs.Race{}
	}
	return cs.Race(*source)
}

func GameRaceToPlayerRace(source cs.Race) *PlayerRace {
	return (*PlayerRace)(&source)
}

func PlayerSpecToGamePlayerSpec(source *PlayerSpec) cs.PlayerSpec {
	// return an empty object for nil to support partial loads
	if source == nil {
		return cs.PlayerSpec{}
	}

	return (cs.PlayerSpec)(*source)
}

func GamePlayerSpecToPlayerSpec(source cs.PlayerSpec) *PlayerSpec {
	return (*PlayerSpec)(&source)
}

func PlayerStatsToGamePlayerStats(source *PlayerStats) *cs.PlayerStats {
	return (*cs.PlayerStats)(source)
}

func GamePlayerStatsToPlayerStats(source *cs.PlayerStats) *PlayerStats {
	return (*PlayerStats)(source)
}

func PlanetSpecToGamePlanetSpec(source *PlanetSpec) cs.PlanetSpec {
	return (cs.PlanetSpec)(*source)
}

func GamePlanetSpecToPlanetSpec(source cs.PlanetSpec) *PlanetSpec {
	return (*PlanetSpec)(&source)
}

func ProductionQueueItemsToGameProductionQueueItems(source *ProductionQueueItems) []cs.ProductionQueueItem {
	return ([]cs.ProductionQueueItem)(*source)
}

func GameProductionQueueItemsToProductionQueueItems(source []cs.ProductionQueueItem) *ProductionQueueItems {
	return (*ProductionQueueItems)(&source)
}

func FleetSpecToGameFleetSpec(source *FleetSpec) cs.FleetSpec {
	return (cs.FleetSpec)(*source)
}

func GameFleetSpecToFleetSpec(source cs.FleetSpec) *FleetSpec {
	return (*FleetSpec)(&source)
}

func ShipTokensToGameShipTokens(source *ShipTokens) []cs.ShipToken {
	return ([]cs.ShipToken)(*source)
}

func GameShipTokensToShipTokens(source []cs.ShipToken) *ShipTokens {
	return (*ShipTokens)(&source)
}

func WaypointsToGameWaypoints(source *Waypoints) []cs.Waypoint {
	return ([]cs.Waypoint)(*source)
}

func GameWaypointsToWaypoints(source []cs.Waypoint) *Waypoints {
	return (*Waypoints)(&source)
}

func MineFieldSpecToGameMineFieldSpec(source *MineFieldSpec) cs.MineFieldSpec {
	return (cs.MineFieldSpec)(*source)
}

func GameMineFieldSpecToMineFieldSpec(source cs.MineFieldSpec) *MineFieldSpec {
	return (*MineFieldSpec)(&source)
}

func ShipDesignSpecToGameShipDesignSpec(source *ShipDesignSpec) cs.ShipDesignSpec {
	return (cs.ShipDesignSpec)(*source)
}

func GameShipDesignSpecToShipDesignSpec(source cs.ShipDesignSpec) *ShipDesignSpec {
	return (*ShipDesignSpec)(&source)
}

func ShipDesignSlotsToGameShipDesignSlots(source *ShipDesignSlots) []cs.ShipDesignSlot {
	return ([]cs.ShipDesignSlot)(*source)
}

func GameShipDesignSlotsToShipDesignSlots(source []cs.ShipDesignSlot) *ShipDesignSlots {
	return (*ShipDesignSlots)(&source)
}

func WormholeSpecToGameWormholeSpec(source *WormholeSpec) cs.WormholeSpec {
	return (cs.WormholeSpec)(*source)
}

func GameWormholeSpecToWormholeSpec(source cs.WormholeSpec) *WormholeSpec {
	return (*WormholeSpec)(&source)
}

func MysteryTraderSpecToGameMysteryTraderSpec(source *MysteryTraderSpec) cs.MysteryTraderSpec {
	return (cs.MysteryTraderSpec)(*source)
}

func GameMysteryTraderSpecToMysteryTraderSpec(source cs.MysteryTraderSpec) *MysteryTraderSpec {
	return (*MysteryTraderSpec)(&source)
}

func MysteryTraderRewardTypeToGameMysteryTraderRewardType(source *MysteryTraderRewardType) cs.MysteryTraderRewardType {
	return (cs.MysteryTraderRewardType)(*source)
}

func GameMysteryTraderRewardTypeToMysteryTraderRewardType(source cs.MysteryTraderRewardType) *MysteryTraderRewardType {
	return (*MysteryTraderRewardType)(&source)
}

func ExtendResearchCost(source Race) cs.ResearchCost {
	return cs.ResearchCost{
		Energy:        source.ResearchCostEnergy,
		Weapons:       source.ResearchCostWeapons,
		Propulsion:    source.ResearchCostPropulsion,
		Construction:  source.ResearchCostConstruction,
		Electronics:   source.ResearchCostElectronics,
		Biotechnology: source.ResearchCostBiotechnology,
	}
}

func ExtendHabLow(source Race) cs.Hab {
	return cs.Hab{
		Grav: source.HabLowGrav,
		Temp: source.HabLowTemp,
		Rad:  source.HabLowRad,
	}
}

func ExtendHabHigh(source Race) cs.Hab {
	return cs.Hab{
		Grav: source.HabHighGrav,
		Temp: source.HabHighTemp,
		Rad:  source.HabHighRad,
	}
}

func ExtendVictoryConditions(source Game) cs.VictoryConditions {
	return cs.VictoryConditions{
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

func ExtendArea(source Game) cs.Vector {
	return cs.Vector{
		X: source.AreaX,
		Y: source.AreaY,
	}
}

func ExtendDefaultRules(source Game) cs.Rules {
	return cs.NewRules()
}

func ExtendTechLevels(source Player) cs.TechLevel {
	return cs.TechLevel{
		Energy:        source.TechLevelsEnergy,
		Weapons:       source.TechLevelsWeapons,
		Propulsion:    source.TechLevelsPropulsion,
		Construction:  source.TechLevelsConstruction,
		Electronics:   source.TechLevelsElectronics,
		Biotechnology: source.TechLevelsBiotechnology,
	}
}

func ExtendTechLevelsSpent(source Player) cs.TechLevel {
	return cs.TechLevel{
		Energy:        source.TechLevelsSpentEnergy,
		Weapons:       source.TechLevelsSpentWeapons,
		Propulsion:    source.TechLevelsSpentPropulsion,
		Construction:  source.TechLevelsSpentConstruction,
		Electronics:   source.TechLevelsSpentElectronics,
		Biotechnology: source.TechLevelsSpentBiotechnology,
	}
}

func ExtendPlanetMapObject(source Planet) cs.MapObject {
	return cs.MapObject{
		GameDBObject: cs.GameDBObject{
			ID:        source.ID,
			GameID:    source.GameID,
			CreatedAt: source.CreatedAt,
			UpdatedAt: source.UpdatedAt,
		},
		Type: cs.MapObjectTypePlanet,
		Position: cs.Vector{
			X: source.X,
			Y: source.Y,
		},
		Name:      source.Name,
		Num:       source.Num,
		PlayerNum: source.PlayerNum,
		Tags:      TagsToGameTags(source.Tags),
	}
}

func ExtendBaseHab(source Planet) cs.Hab {
	return cs.Hab{
		Grav: source.BaseGrav,
		Temp: source.BaseTemp,
		Rad:  source.BaseRad,
	}
}

func ExtendTerraformedAmount(source Planet) cs.Hab {
	return cs.Hab{
		Grav: source.TerraformedAmountGrav,
		Temp: source.TerraformedAmountTemp,
		Rad:  source.TerraformedAmountRad,
	}
}

func ExtendMineralConcentration(source Planet) cs.Mineral {
	return cs.Mineral{
		Ironium:   source.MineralConcIronium,
		Boranium:  source.MineralConcBoranium,
		Germanium: source.MineralConcGermanium,
	}
}

func ExtendMineYears(source Planet) cs.Mineral {
	return cs.Mineral{
		Ironium:   source.MineYearsIronium,
		Boranium:  source.MineYearsBoranium,
		Germanium: source.MineYearsGermanium,
	}
}

func ExtendFleetMapObject(source Fleet) cs.MapObject {
	return cs.MapObject{
		GameDBObject: cs.GameDBObject{
			ID:        source.ID,
			GameID:    source.GameID,
			CreatedAt: source.CreatedAt,
			UpdatedAt: source.UpdatedAt,
		},
		Type: cs.MapObjectTypeFleet,
		Position: cs.Vector{
			X: source.X,
			Y: source.Y,
		},
		Name:      source.Name,
		Num:       source.Num,
		PlayerNum: source.PlayerNum,
		Tags:      TagsToGameTags(source.Tags),
	}
}

func ExtendFleetFleetOrders(source Fleet) cs.FleetOrders {
	return cs.FleetOrders{
		BattlePlanNum: source.BattlePlanNum,
		Waypoints:     *source.Waypoints,
		RepeatOrders:  source.RepeatOrders,
		Purpose:       source.Purpose,
	}
}

func ExtendFleetHeading(source Fleet) cs.Vector {
	return cs.Vector{
		X: source.HeadingX,
		Y: source.HeadingY,
	}
}

func ExtendFleetPreviousPosition(source Fleet) *cs.Vector {
	if source.PreviousPositionX == nil || source.PreviousPositionY == nil {
		return nil
	}
	return &cs.Vector{
		X: *source.PreviousPositionX,
		Y: *source.PreviousPositionY,
	}
}

func ExtendMysteryTraderMapObject(source MysteryTrader) cs.MapObject {
	return cs.MapObject{
		GameDBObject: cs.GameDBObject{
			ID:        source.ID,
			GameID:    source.GameID,
			CreatedAt: source.CreatedAt,
			UpdatedAt: source.UpdatedAt,
		},
		Type: cs.MapObjectTypeMysteryTrader,
		Position: cs.Vector{
			X: source.X,
			Y: source.Y,
		},
		Name: source.Name,
		Num:  source.Num,
		Tags: TagsToGameTags(source.Tags),
	}
}

func ExtendMysteryTraderHeading(source MysteryTrader) cs.Vector {
	return cs.Vector{
		X: source.HeadingX,
		Y: source.HeadingY,
	}
}

func ExtendMysteryTraderDestination(source MysteryTrader) cs.Vector {
	return cs.Vector{
		X: source.DestinationX,
		Y: source.DestinationY,
	}
}

func ExtendSalvageMapObject(source Salvage) cs.MapObject {
	return cs.MapObject{
		GameDBObject: cs.GameDBObject{
			ID:        source.ID,
			GameID:    source.GameID,
			CreatedAt: source.CreatedAt,
			UpdatedAt: source.UpdatedAt,
		},
		Type: cs.MapObjectTypeSalvage,
		Position: cs.Vector{
			X: source.X,
			Y: source.Y,
		},
		Name:      source.Name,
		Num:       source.Num,
		PlayerNum: source.PlayerNum,
		Tags:      TagsToGameTags(source.Tags),
	}
}

func ExtendMineFieldMapObject(source MineField) cs.MapObject {
	return cs.MapObject{
		GameDBObject: cs.GameDBObject{
			ID:        source.ID,
			GameID:    source.GameID,
			CreatedAt: source.CreatedAt,
			UpdatedAt: source.UpdatedAt,
		},
		Type: cs.MapObjectTypeMineField,
		Position: cs.Vector{
			X: source.X,
			Y: source.Y,
		},
		Name:      source.Name,
		Num:       source.Num,
		PlayerNum: source.PlayerNum,
		Tags:      TagsToGameTags(source.Tags),
	}
}

func ExtendMineralPacketHeading(source MineralPacket) cs.Vector {
	return cs.Vector{
		X: source.HeadingX,
		Y: source.HeadingY,
	}
}

func ExtendMineralPacketMapObject(source MineralPacket) cs.MapObject {
	return cs.MapObject{
		GameDBObject: cs.GameDBObject{
			ID:        source.ID,
			GameID:    source.GameID,
			CreatedAt: source.CreatedAt,
			UpdatedAt: source.UpdatedAt,
		},
		Type: cs.MapObjectTypeMineralPacket,
		Position: cs.Vector{
			X: source.X,
			Y: source.Y,
		},
		Name:      source.Name,
		Num:       source.Num,
		PlayerNum: source.PlayerNum,
		Tags:      TagsToGameTags(source.Tags),
	}
}
