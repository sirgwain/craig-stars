//go:generate go run github.com/jmattheis/goverter/cmd/goverter --packageName db --output ./db/generated.go --packagePath github.com/sirgwain/craig-stars/db --ignoreUnexportedFields ./db
package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/sirgwain/craig-stars/cs"
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
// goverter:extend PlayerRelationshipsToGamePlayerRelationships
// goverter:extend GamePlayerRelationshipsToPlayerRelationships
// goverter:extend PlayerMessagesToGamePlayerMessages
// goverter:extend GamePlayerMessagesToPlayerMessages
// goverter:extend PlayerIntelsToGamePlayerIntels
// goverter:extend GamePlayerIntelsToPlayerIntels
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
// goverter:extend WormholeIntelsToGameWormholeIntels
// goverter:extend GameWormholeIntelsToWormholeIntels
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
// goverter:extend WormholeSpecToGameWormholeSpec
// goverter:extend GameWormholeSpecToWormholeSpec
type Converter interface {
	ConvertUser(source User) cs.User
	ConvertGameUser(source *cs.User) *User
	ConvertUsers(source []User) []cs.User

	// goverter:mapExtend ResearchCost ExtendResearchCost
	// goverter:mapExtend HabLow ExtendHabLow
	// goverter:mapExtend HabHigh ExtendHabHigh
	ConvertRace(source Race) cs.Race
	// goverter:mapExtend ResearchCost ExtendResearchCost
	// goverter:mapExtend HabLow ExtendHabLow
	// goverter:mapExtend HabHigh ExtendHabHigh
	ConvertRaces(source []Race) []cs.Race

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

	// goverter:mapExtend VictoryConditions ExtendVictoryConditions
	// goverter:mapExtend Area ExtendArea
	// goverter:mapExtend Rules ExtendDefaultRules
	ConvertGame(source Game) cs.Game

	// goverter:mapExtend VictoryConditions ExtendVictoryConditions
	// goverter:mapExtend Area ExtendArea
	ConvertGames(source []Game) []cs.Game

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
	// goverter:map PlayerOrders.Researching Researching
	// goverter:map PlayerOrders.NextResearchField NextResearchField
	// goverter:map PlayerOrders.ResearchAmount ResearchAmount
	// goverter:map PlayerIntels.PlayerIntels PlayerIntels
	// goverter:map PlayerIntels.PlanetIntels PlanetIntels
	// goverter:map PlayerIntels.FleetIntels FleetIntels
	// goverter:map PlayerIntels.ShipDesignIntels ShipDesignIntels
	// goverter:map PlayerIntels.MineralPacketIntels MineralPacketIntels
	// goverter:map PlayerIntels.MineFieldIntels MineFieldIntels
	// goverter:map PlayerIntels.WormholeIntels WormholeIntels
	// goverter:map PlayerPlans.BattlePlans BattlePlans
	// goverter:map PlayerPlans.ProductionPlans ProductionPlans
	// goverter:map PlayerPlans.TransportPlans TransportPlans
	ConvertGamePlayer(source *cs.Player) *Player

	// goverter:mapExtend TechLevels ExtendTechLevels
	// goverter:mapExtend TechLevelsSpent ExtendTechLevelsSpent
	// goverter:mapExtend PlayerOrders ExtendPlayerPlayerOrders
	// goverter:mapExtend PlayerIntels ExtendPlayerPlayerIntels
	// goverter:mapExtend PlayerPlans ExtendPlayerPlayerPlans
	// goverter:ignore Designs
	ConvertPlayer(source Player) cs.Player
	// goverter:mapExtend TechLevels ExtendTechLevels
	// goverter:mapExtend TechLevelsSpent ExtendTechLevelsSpent
	// goverter:mapExtend PlayerOrders ExtendPlayerPlayerOrders
	// goverter:mapExtend PlayerIntels ExtendPlayerPlayerIntels
	// goverter:mapExtend PlayerPlans ExtendPlayerPlayerPlans
	ConvertPlayers(source []Player) []cs.Player

	// goverter:map MapObject.ID ID
	// goverter:map MapObject.GameID GameID
	// goverter:map MapObject.CreatedAt CreatedAt
	// goverter:map MapObject.UpdatedAt UpdatedAt
	// goverter:map MapObject.Type Type
	// goverter:map MapObject.Dirty Dirty
	// goverter:map MapObject.Delete Delete
	// goverter:map MapObject.Position.X X
	// goverter:map MapObject.Position.Y Y
	// goverter:map MapObject.Name Name
	// goverter:map MapObject.Num Num
	// goverter:map MapObject.PlayerNum	 PlayerNum
	// goverter:map PlanetOrders.ContributesOnlyLeftoverToResearch ContributesOnlyLeftoverToResearch
	// goverter:map PlanetOrders.ProductionQueue ProductionQueue
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
	ConvertGamePlanet(source *cs.Planet) *Planet

	// goverter:mapExtend Hab ExtendHab
	// goverter:mapExtend BaseHab ExtendBaseHab
	// goverter:mapExtend TerraformedAmount ExtendTerraformedAmount
	// goverter:mapExtend MineralConcentration ExtendMineralConcentration
	// goverter:mapExtend MineYears ExtendMineYears
	// goverter:mapExtend Cargo ExtendPlanetCargo
	// goverter:mapExtend MapObject ExtendPlanetMapObject
	// goverter:mapExtend PlanetOrders ExtendPlanetPlanetOrders
	ConvertPlanet(source *Planet) *cs.Planet

	// goverter:map MapObject.ID ID
	// goverter:map MapObject.GameID GameID
	// goverter:map MapObject.CreatedAt CreatedAt
	// goverter:map MapObject.UpdatedAt UpdatedAt
	// goverter:map MapObject.Type Type
	// goverter:map MapObject.Dirty Dirty
	// goverter:map MapObject.Delete Delete
	// goverter:map MapObject.Position.X X
	// goverter:map MapObject.Position.Y Y
	// goverter:map MapObject.Name Name
	// goverter:map MapObject.Num Num
	// goverter:map MapObject.PlayerNum	 PlayerNum
	// goverter:map FleetOrders.BattlePlanName BattlePlanName
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
	ConvertGameFleet(source *cs.Fleet) *Fleet

	// goverter:mapExtend Heading ExtendFleetHeading
	// goverter:mapExtend PreviousPosition ExtendFleetPreviousPosition
	// goverter:mapExtend Cargo ExtendFleetCargo
	// goverter:mapExtend MapObject ExtendFleetMapObject
	// goverter:mapExtend FleetOrders ExtendFleetFleetOrders
	// goverter:ignore Tokens
	ConvertFleet(source *Fleet) *cs.Fleet

	ConvertGameShipDesign(source *cs.ShipDesign) *ShipDesign
	// goverter:ignore Dirty
	ConvertShipDesign(source *ShipDesign) *cs.ShipDesign

	ConvertGameShipToken(source cs.ShipToken) ShipToken
	// goverter:ignore Design
	ConvertShipToken(source ShipToken) cs.ShipToken

	// goverter:map MapObject.ID ID
	// goverter:map MapObject.GameID GameID
	// goverter:map MapObject.CreatedAt CreatedAt
	// goverter:map MapObject.UpdatedAt UpdatedAt
	// goverter:map MapObject.Type Type
	// goverter:map MapObject.Dirty Dirty
	// goverter:map MapObject.Delete Delete
	// goverter:map MapObject.Position.X X
	// goverter:map MapObject.Position.Y Y
	// goverter:map MapObject.Name Name
	// goverter:map MapObject.Num Num
	// goverter:ignore PlayerNum
	ConvertGameWormhole(source *cs.Wormhole) *Wormhole

	// goverter:mapExtend MapObject ExtendWormholeMapObject
	ConvertWormhole(source *Wormhole) *cs.Wormhole

	// goverter:map MapObject.ID ID
	// goverter:map MapObject.GameID GameID
	// goverter:map MapObject.CreatedAt CreatedAt
	// goverter:map MapObject.UpdatedAt UpdatedAt
	// goverter:map MapObject.Type Type
	// goverter:map MapObject.Dirty Dirty
	// goverter:map MapObject.Delete Delete
	// goverter:map MapObject.Position.X X
	// goverter:map MapObject.Position.Y Y
	// goverter:map MapObject.Name Name
	// goverter:map MapObject.Num Num
	// goverter:map MapObject.PlayerNum	 PlayerNum
	// goverter:ignore Tags
	// goverter:map Cargo.Ironium Ironium
	// goverter:map Cargo.Boranium Boranium
	// goverter:map Cargo.Germanium Germanium
	ConvertGameSalvage(source *cs.Salvage) *Salvage

	// goverter:mapExtend MapObject ExtendSalvageMapObject
	// goverter:mapExtend Cargo ExtendSalvageCargo
	ConvertSalvage(source *Salvage) *cs.Salvage

	// goverter:map MapObject.ID ID
	// goverter:map MapObject.GameID GameID
	// goverter:map MapObject.CreatedAt CreatedAt
	// goverter:map MapObject.UpdatedAt UpdatedAt
	// goverter:map MapObject.Dirty Dirty
	// goverter:map MapObject.Delete Delete
	// goverter:map MapObject.Position.X X
	// goverter:map MapObject.Position.Y Y
	// goverter:map MapObject.Name Name
	// goverter:map MapObject.Num Num
	// goverter:map MapObject.PlayerNum PlayerNum
	// goverter:map MineFieldOrders.Detonate Detonate
	// goverter:ignore Tags
	ConvertGameMineField(source *cs.MineField) *MineField

	// goverter:mapExtend MapObject ExtendMineFieldMapObject
	// goverter:mapExtend MineFieldOrders ExtendMineFieldMineFieldOrders
	ConvertMineField(source *MineField) *cs.MineField

	// goverter:map MapObject.ID ID
	// goverter:map MapObject.GameID GameID
	// goverter:map MapObject.CreatedAt CreatedAt
	// goverter:map MapObject.UpdatedAt UpdatedAt
	// goverter:map MapObject.Type Type
	// goverter:map MapObject.Dirty Dirty
	// goverter:map MapObject.Delete Delete
	// goverter:map MapObject.Position.X X
	// goverter:map MapObject.Position.Y Y
	// goverter:map MapObject.Name Name
	// goverter:map MapObject.Num Num
	// goverter:map MapObject.PlayerNum	 PlayerNum
	// goverter:ignore Tags
	// goverter:map Cargo.Ironium Ironium
	// goverter:map Cargo.Boranium Boranium
	// goverter:map Cargo.Germanium Germanium
	// goverter:map Heading.X HeadingX
	// goverter:map Heading.Y HeadingY
	ConvertGameMineralPacket(source *cs.MineralPacket) *MineralPacket

	// goverter:mapExtend MapObject ExtendMineralPacketMapObject
	// goverter:mapExtend Cargo ExtendMineralPacketCargo
	// goverter:mapExtend Heading ExtendMineralPacketHeading
	ConvertMineralPacket(source *MineralPacket) *cs.MineralPacket
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

func ExtendPlayerPlayerOrders(source Player) cs.PlayerOrders {
	return cs.PlayerOrders{
		Researching:       source.Researching,
		NextResearchField: source.NextResearchField,
		ResearchAmount:    source.ResearchAmount,
	}
}

func ExtendPlayerPlayerPlans(source Player) cs.PlayerPlans {
	plans := cs.PlayerPlans{}

	if source.ProductionPlans != nil {
		plans.ProductionPlans = *source.ProductionPlans
	}

	if source.BattlePlans != nil {
		plans.BattlePlans = *source.BattlePlans
	}

	if source.TransportPlans != nil {
		plans.TransportPlans = *source.TransportPlans
	}

	return plans
}

func ExtendPlayerPlayerIntels(source Player) cs.PlayerIntels {
	intels := cs.PlayerIntels{}

	if source.PlayerIntels != nil {
		intels.PlayerIntels = *source.PlayerIntels
	}

	if source.PlanetIntels != nil {
		intels.PlanetIntels = *source.PlanetIntels
	}

	if source.FleetIntels != nil {
		intels.FleetIntels = *source.FleetIntels
	}

	if source.ShipDesignIntels != nil {
		intels.ShipDesignIntels = *source.ShipDesignIntels
	}

	if source.MineralPacketIntels != nil {
		intels.MineralPacketIntels = *source.MineralPacketIntels
	}

	if source.MineFieldIntels != nil {
		intels.MineFieldIntels = *source.MineFieldIntels
	}

	if source.WormholeIntels != nil {
		intels.WormholeIntels = *source.WormholeIntels
	}

	return intels
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

func WaypointsToGameWaypoints(source *Waypoints) []cs.Waypoint {
	return ([]cs.Waypoint)(*source)
}

func GameWaypointsToWaypoints(source []cs.Waypoint) *Waypoints {
	return (*Waypoints)(&source)
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
		Type:      cs.MapObjectTypePlanet,
		ID:        source.ID,
		GameID:    source.GameID,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
		Position: cs.Vector{
			X: source.X,
			Y: source.Y,
		},
		Name:      source.Name,
		Num:       source.Num,
		PlayerNum: source.PlayerNum,
		// Tags:      source.Tags,
	}
}

func ExtendPlanetPlanetOrders(source Planet) cs.PlanetOrders {
	return cs.PlanetOrders{
		ContributesOnlyLeftoverToResearch: source.ContributesOnlyLeftoverToResearch,
		ProductionQueue:                   *source.ProductionQueue,
	}
}

func ExtendHab(source Planet) cs.Hab {
	return cs.Hab{
		Grav: source.Grav,
		Temp: source.Temp,
		Rad:  source.Rad,
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

func ExtendPlanetCargo(source Planet) cs.Cargo {
	return cs.Cargo{
		Ironium:   source.Ironium,
		Boranium:  source.Boranium,
		Germanium: source.Germanium,
		Colonists: source.Colonists,
	}
}

func ExtendFleetMapObject(source Fleet) cs.MapObject {
	return cs.MapObject{
		Type:      cs.MapObjectTypeFleet,
		ID:        source.ID,
		GameID:    source.GameID,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
		Position: cs.Vector{
			X: source.X,
			Y: source.Y,
		},
		Name:      source.Name,
		Num:       source.Num,
		PlayerNum: source.PlayerNum,
		// Tags:      source.Tags,
	}
}

func ExtendFleetFleetOrders(source Fleet) cs.FleetOrders {
	return cs.FleetOrders{
		BattlePlanName: source.BattlePlanName,
		Waypoints:      *source.Waypoints,
		RepeatOrders:   source.RepeatOrders,
	}
}

func ExtendFleetCargo(source Fleet) cs.Cargo {
	return cs.Cargo{
		Ironium:   source.Ironium,
		Boranium:  source.Boranium,
		Germanium: source.Germanium,
		Colonists: source.Colonists,
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

func ExtendWormholeMapObject(source Wormhole) cs.MapObject {
	return cs.MapObject{
		Type:      cs.MapObjectTypeWormhole,
		ID:        source.ID,
		GameID:    source.GameID,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
		Position: cs.Vector{
			X: source.X,
			Y: source.Y,
		},
		Name: source.Name,
		Num:  source.Num,
	}
}

func ExtendSalvageMapObject(source Salvage) cs.MapObject {
	return cs.MapObject{
		Type:      cs.MapObjectTypeSalvage,
		ID:        source.ID,
		GameID:    source.GameID,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
		Position: cs.Vector{
			X: source.X,
			Y: source.Y,
		},
		Name:      source.Name,
		Num:       source.Num,
		PlayerNum: source.PlayerNum,
		// Tags:      source.Tags,
	}
}

func ExtendSalvageCargo(source Salvage) cs.Cargo {
	return cs.Cargo{
		Ironium:   source.Ironium,
		Boranium:  source.Boranium,
		Germanium: source.Germanium,
	}
}

func ExtendMineFieldMapObject(source MineField) cs.MapObject {
	return cs.MapObject{
		Type:      cs.MapObjectTypeMineField,
		ID:        source.ID,
		GameID:    source.GameID,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
		Position: cs.Vector{
			X: source.X,
			Y: source.Y,
		},
		Name:      source.Name,
		Num:       source.Num,
		PlayerNum: source.PlayerNum,
		// Tags:      source.Tags,
	}
}

func ExtendMineFieldMineFieldOrders(source MineField) cs.MineFieldOrders {
	return cs.MineFieldOrders{
		Detonate: source.Detonate,
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
		Type:      cs.MapObjectTypeMineralPacket,
		ID:        source.ID,
		GameID:    source.GameID,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
		Position: cs.Vector{
			X: source.X,
			Y: source.Y,
		},
		Name:      source.Name,
		Num:       source.Num,
		PlayerNum: source.PlayerNum,
		// Tags:      source.Tags,
	}
}

func ExtendMineralPacketCargo(source MineralPacket) cs.Cargo {
	return cs.Cargo{
		Ironium:   source.Ironium,
		Boranium:  source.Boranium,
		Germanium: source.Germanium,
	}
}
