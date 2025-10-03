//go:generate go tool github.com/jmattheis/goverter/cmd/goverter gen ./
package db

import (
	"database/sql"
	"time"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db/generated"
)

// will be instanciated in ./converter.init.go
var c Converter

// goverter:converter
// goverter:name GameConverter
// goverter:output:package github.com/sirgwain/craig-stars/db
// goverter:output:file ./converter.generated.go
// goverter:extend AcquiredTechsToGameAcquiredTechs
// goverter:extend BattlePlansToGameBattlePlans
// goverter:extend BattleRecordsToGameBattleRecords
// goverter:extend CargoTransfersToGameCargoTransfers
// goverter:extend FleetIntelsToGameFleetIntels
// goverter:extend FleetSpecToGameFleetSpec
// goverter:extend Float64ToNullFloat64
// goverter:extend GameAcquiredTechsToAcquiredTechs
// goverter:extend GameBattlePlansToBattlePlans
// goverter:extend GameBattleRecordsToBattleRecords
// goverter:extend GameCargoTransfersToCargoTransfers
// goverter:extend GameFleetIntelsToFleetIntels
// goverter:extend GameFleetSpecToFleetSpec
// goverter:extend GameMinefieldIntelsToMinefieldIntels
// goverter:extend GameMineralPacketIntelsToMineralPacketIntels
// goverter:extend GameMysteryTraderIntelsToMysteryTraderIntels
// goverter:extend GameMysteryTraderPlayersRewardedToMysteryTraderPlayersRewarded
// goverter:extend GameMysteryTraderSpecToMysteryTraderSpec
// goverter:extend GamePlanetIntelsToPlanetIntels
// goverter:extend GamePlanetSpecToPlanetSpec
// goverter:extend GamePlayerIntelsToPlayerIntels
// goverter:extend GamePlayerMessagesToPlayerMessages
// goverter:extend GamePlayerRelationshipsToPlayerRelationships
// goverter:extend GamePlayerScoresToPlayerScores
// goverter:extend GamePlayerSpecToPlayerSpec
// goverter:extend GamePlayerStatsToPlayerStats
// goverter:extend GameProductionPlansToProductionPlans
// goverter:extend GameProductionQueueItemsToProductionQueueItems
// goverter:extend GameRaceToPlayerRace
// goverter:extend GameSalvageIntelsToSalvageIntels
// goverter:extend GameScoreIntelsToScoreIntels
// goverter:extend GameShipDesignIntelsToShipDesignIntels
// goverter:extend GameShipDesignSlotsToShipDesignSlots
// goverter:extend GameShipDesignSpecToShipDesignSpec
// goverter:extend GameShipTokensToShipTokens
// goverter:extend GameTagsToTags
// goverter:extend GameTransportPlansToTransportPlans
// goverter:extend GameWaypointsToWaypoints
// goverter:extend GameWormholeIntelsToWormholeIntels
// goverter:extend GameWormholeSpecToWormholeSpec
// goverter:extend Int64ToInt
// goverter:extend IntToInt64
// goverter:extend Float64ToInt
// goverter:extend IntToFloat64
// goverter:extend NullFloat64ToInt
// goverter:extend IntToNullFloat64
// goverter:extend IntToNullInt64
// goverter:extend MinefieldIntelsToGameMinefieldIntels
// goverter:extend MineralPacketIntelsToGameMineralPacketIntels
// goverter:extend MysteryTraderIntelsToGameMysteryTraderIntels
// goverter:extend MysteryTraderPlayersRewardedToGameMysteryTraderPlayersRewarded
// goverter:extend MysteryTraderSpecToGameMysteryTraderSpec
// goverter:extend NullBoolToBool
// goverter:extend NullInt64ToInt
// goverter:extend NullInt64ToInt64
// goverter:extend NullStringToString
// goverter:extend NullTimeToTime
// goverter:extend PlanetIntelsToGamePlanetIntels
// goverter:extend PlanetSpecToGamePlanetSpec
// goverter:extend PlayerIntelsToGamePlayerIntels
// goverter:extend PlayerMessagesToGamePlayerMessages
// goverter:extend PlayerRaceToGameRace
// goverter:extend PlayerRelationshipsToGamePlayerRelationships
// goverter:extend PlayerScoresToGamePlayerScores
// goverter:extend PlayerSpecToGamePlayerSpec
// goverter:extend PlayerStatsToGamePlayerStats
// goverter:extend ProductionPlansToGameProductionPlans
// goverter:extend ProductionQueueItemsToGameProductionQueueItems
// goverter:extend SalvageIntelsToGameSalvageIntels
// goverter:extend ScoreIntelsToGameScoreIntels
// goverter:extend ShipDesignIntelsToGameShipDesignIntels
// goverter:extend ShipDesignSlotsToGameShipDesignSlots
// goverter:extend ShipDesignSpecToGameShipDesignSpec
// goverter:extend ShipTokensToGameShipTokens
// goverter:extend TagsToGameTags
// goverter:extend TransportPlansToGameTransportPlans
// goverter:extend WaypointsToGameWaypoints
// goverter:extend WormholeIntelsToGameWormholeIntels
// goverter:extend WormholeSpecToGameWormholeSpec
// goverter:enum no
// goverter:ignoreUnexported
// goverter:matchIgnoreCase
// goverter:skipCopySameType
// goverter:useZeroValueOnPointerInconsistency
type Converter interface {
	// goverter:autoMap DBObject
	// goverter:autoMap UserSettings
	ConvertGameUserToCreateParams(source *cs.User) generated.CreateUserParams

	// goverter:autoMap DBObject
	// goverter:autoMap UserSettings
	ConvertGameUserToUpdateParams(source *cs.User) generated.UpdateUserParams

	// goverter:map . DBObject
	// goverter:map . UserSettings
	ConvertUser(source generated.User) cs.User

	ConvertUsers(source []generated.User) []cs.User

	// goverter:map . DBObject
	// goverter:map . ResearchCost | ExtendResearchCost
	// goverter:map . HabLow | ExtendHabLow
	// goverter:map . HabHigh | ExtendHabHigh
	// goverter:ignore Spec
	ConvertRace(source generated.Race) cs.Race

	ConvertRaces(source []generated.Race) []cs.Race

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
	ConvertGameRace(source *cs.Race) generated.Race

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
	ConvertGameRaceToCreateParams(source *cs.Race) generated.CreateRaceParams

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
	ConvertGameRaceToUpdateParams(source *cs.Race) generated.UpdateRaceParams

	// goverter:map . DBObject
	// goverter:map . VictoryConditions | ExtendVictoryConditions
	// goverter:map . Area | ExtendArea
	// goverter:ignore Rules
	ConvertGame(source generated.Game) cs.Game

	ConvertGames(source []generated.Game) []cs.Game

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
	ConvertGameGame(source *cs.Game) generated.Game

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
	ConvertGameGameToCreateParams(source *cs.Game) generated.CreateGameParams

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
	ConvertGameGameToUpdateParams(source *cs.Game) generated.UpdateGameParams

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
	// goverter:autoMap Intels
	// goverter:autoMap PlayerPlans
	ConvertGamePlayer(source *cs.Player) generated.Player

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
	// goverter:autoMap Intels
	// goverter:autoMap PlayerPlans
	ConvertGamePlayerToCreateParams(source *cs.Player) generated.CreatePlayerParams

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
	// goverter:autoMap Intels
	// goverter:autoMap PlayerPlans
	ConvertGamePlayerToUpdateParams(source *cs.Player) generated.UpdatePlayerParams

	// goverter:map . GameDBObject | ExtendPlayerGameDBObject
	// goverter:map . TechLevels | ExtendTechLevels
	// goverter:map . TechLevelsSpent | ExtendTechLevelsSpent
	// goverter:map . PlayerOrders
	// goverter:map . Intels
	// goverter:map . PlayerPlans
	// goverter:ignore Designs
	// goverter:ignore Spec
	// goverter:ignore TechsJustGained
	ConvertPlayer(source generated.Player) cs.Player

	// goverter:map . GameDBObject | ExtendLightPlayerGameDBObject
	// goverter:map . TechLevels | ExtendTechLevelsLight
	// goverter:map . TechLevelsSpent | ExtendTechLevelsSpentLight
	// goverter:map . PlayerOrders
	// goverter:map . PlayerPlans
	// goverter:ignore Intels
	// goverter:ignore Designs
	// goverter:ignore Spec
	// goverter:ignore TechsJustGained
	ConvertLightPlayer(source generated.GetLightPlayerForGameRow) cs.Player

	// goverter:map . GameDBObject | ExtendPlayerStatusGameDBObject
	// goverter:ignoreMissing
	ConvertPlayerStatus(source generated.GetPlayersStatusForGameRow) cs.Player
	ConvertPlayerStatuses(source []generated.GetPlayersStatusForGameRow) []*cs.Player

	ConvertGetGamesWithPlayersRowToPlayerStatus(source generated.GetGamesWithPlayersRow) cs.GamePlayer
	ConvertGetGamesWithPlayersForUserRowToPlayerStatus(source generated.GetGamesWithPlayersForUserRow) cs.GamePlayer
	ConvertGetGameWithPlayersRowToPlayerStatus(source generated.GetGameWithPlayersRow) cs.GamePlayer
	ConvertGetPlayerForGameRowToShipDesign(source generated.GetPlayerForGameRow) generated.ShipDesign
	ConvertGetPlayersWithDesignsForGameRowToShipDesign(source generated.GetPlayersWithDesignsForGameRow) generated.ShipDesign

	ConvertPlayers(source []generated.Player) []*cs.Player

	// goverter:autoMap GameDBObject
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
	ConvertGamePlanet(source *cs.Planet) generated.Planet

	// goverter:autoMap GameDBObject
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
	ConvertGamePlanetToCreateParams(source *cs.Planet) generated.CreatePlanetParams

	// goverter:autoMap GameDBObject
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
	ConvertGamePlanetToUpdateParams(source *cs.Planet) generated.UpdatePlanetParams

	// goverter:map . Hab
	// goverter:map . BaseHab | ExtendBaseHab
	// goverter:map . TerraformedAmount | ExtendTerraformedAmount
	// goverter:map . MineralConcentration | ExtendMineralConcentration
	// goverter:map . MineYears | ExtendMineYears
	// goverter:map . Cargo
	// goverter:map . GameDBObject
	// goverter:map . MapObject | ExtendPlanetMapObject
	// goverter:map . PlanetOrders
	// goverter:ignore Starbase
	// goverter:ignore Dirty
	ConvertPlanet(source generated.Planet) *cs.Planet
	ConvertPlanets(source []generated.Planet) []*cs.Planet

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:autoMap FleetOrders
	// goverter:map Heading.X HeadingX
	// goverter:map Heading.Y HeadingY
	// goverter:map PreviousPosition.X PreviousPositionX
	// goverter:map PreviousPosition.Y PreviousPositionY
	// goverter:autoMap Cargo
	ConvertGameFleet(source *cs.Fleet) generated.Fleet

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:autoMap FleetOrders
	// goverter:map Heading.X HeadingX
	// goverter:map Heading.Y HeadingY
	// goverter:map PreviousPosition.X PreviousPositionX
	// goverter:map PreviousPosition.Y PreviousPositionY
	// goverter:autoMap Cargo
	ConvertGameFleetToCreateParams(source *cs.Fleet) generated.CreateFleetParams

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:autoMap FleetOrders
	// goverter:map Heading.X HeadingX
	// goverter:map Heading.Y HeadingY
	// goverter:map PreviousPosition.X PreviousPositionX
	// goverter:map PreviousPosition.Y PreviousPositionY
	// goverter:autoMap Cargo
	ConvertGameFleetToUpdateParams(source *cs.Fleet) generated.UpdateFleetParams

	// goverter:map . Heading | ExtendFleetHeading
	// goverter:map . PreviousPosition | ExtendFleetPreviousPosition
	// goverter:map . Cargo
	// goverter:map . GameDBObject
	// goverter:map . MapObject | ExtendFleetMapObject
	// goverter:map . FleetOrders | ExtendFleetFleetOrders
	ConvertFleet(source generated.Fleet) *cs.Fleet

	ConvertFleets(source []generated.Fleet) []*cs.Fleet

	// goverter:autoMap GameDBObject
	ConvertGameShipDesign(source *cs.ShipDesign) generated.ShipDesign

	// goverter:autoMap GameDBObject
	ConvertGameShipDesignToCreateParams(source *cs.ShipDesign) generated.CreateShipDesignParams

	// goverter:autoMap GameDBObject
	ConvertGameShipDesignToUpdateParams(source *cs.ShipDesign) generated.UpdateShipDesignParams

	// goverter:ignore Delete
	// goverter:map . GameDBObject
	ConvertShipDesign(source generated.ShipDesign) *cs.ShipDesign

	ConvertShipDesigns(source []generated.ShipDesign) []*cs.ShipDesign

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	ConvertGameWormhole(source *cs.Wormhole) generated.Wormhole

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	ConvertGameWormholeToCreateParams(source *cs.Wormhole) generated.CreateWormholeParams

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	ConvertGameWormholeToUpdateParams(source *cs.Wormhole) generated.UpdateWormholeParams

	// goverter:map . GameDBObject
	// goverter:map . MapObject
	ConvertWormhole(source generated.Wormhole) *cs.Wormhole
	ConvertWormholes(source []generated.Wormhole) []*cs.Wormhole

	// goverter:map . Position
	// goverter:map Type | MapObjectTypeWormhole
	// goverter:ignore Delete
	// goverter:ignore PlayerNum
	wormHoleMapObject(source generated.Wormhole) cs.MapObject

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:map Heading.X HeadingX
	// goverter:map Heading.Y HeadingY
	// goverter:map Destination.X DestinationX
	// goverter:map Destination.Y DestinationY
	ConvertGameMysteryTrader(source *cs.MysteryTrader) generated.MysteryTrader

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:map Heading.X HeadingX
	// goverter:map Heading.Y HeadingY
	// goverter:map Destination.X DestinationX
	// goverter:map Destination.Y DestinationY
	ConvertGameMysteryTraderToCreateParams(source *cs.MysteryTrader) generated.CreateMysteryTraderParams

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:map Heading.X HeadingX
	// goverter:map Heading.Y HeadingY
	// goverter:map Destination.X DestinationX
	// goverter:map Destination.Y DestinationY
	ConvertGameMysteryTraderToUpdateParams(source *cs.MysteryTrader) generated.UpdateMysteryTraderParams

	// goverter:map . GameDBObject
	// goverter:map . MapObject | ExtendMysteryTraderMapObject
	// goverter:map . Heading | ExtendMysteryTraderHeading
	// goverter:map . Destination | ExtendMysteryTraderDestination
	ConvertMysteryTrader(source generated.MysteryTrader) *cs.MysteryTrader
	ConvertMysteryTraders(source []generated.MysteryTrader) []*cs.MysteryTrader

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:autoMap Cargo
	ConvertGameSalvage(source *cs.Salvage) generated.Salvage

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:autoMap Cargo
	ConvertGameSalvageToCreateParams(source *cs.Salvage) generated.CreateSalvageParams

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:autoMap Cargo
	ConvertGameSalvageToUpdateParams(source *cs.Salvage) generated.UpdateSalvageParams

	// goverter:map . GameDBObject
	// goverter:map . MapObject | ExtendSalvageMapObject
	// goverter:map . Cargo
	ConvertSalvage(source generated.Salvage) *cs.Salvage
	ConvertSalvages(source []generated.Salvage) []*cs.Salvage

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:map MinefieldOrders.Detonate Detonate
	ConvertGameMinefield(source *cs.Minefield) generated.Minefield

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:map MinefieldOrders.Detonate Detonate
	ConvertGameMinefieldToCreateParams(source *cs.Minefield) generated.CreateMinefieldParams

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:map MinefieldOrders.Detonate Detonate
	ConvertGameMinefieldToUpdateParams(source *cs.Minefield) generated.UpdateMinefieldParams

	// goverter:map . GameDBObject
	// goverter:map . MapObject | ExtendMinefieldMapObject
	// goverter:map . MinefieldOrders
	ConvertMinefield(source generated.Minefield) *cs.Minefield
	ConvertMinefields(source []generated.Minefield) []*cs.Minefield

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:autoMap Cargo
	// goverter:map Cargo.Ironium Ironium
	// goverter:map Cargo.Boranium Boranium
	// goverter:map Cargo.Germanium Germanium
	// goverter:map Heading.X HeadingX
	// goverter:map Heading.Y HeadingY
	ConvertGameMineralPacket(source *cs.MineralPacket) generated.MineralPacket

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:autoMap Cargo
	// goverter:map Cargo.Ironium Ironium
	// goverter:map Cargo.Boranium Boranium
	// goverter:map Cargo.Germanium Germanium
	// goverter:map Heading.X HeadingX
	// goverter:map Heading.Y HeadingY
	ConvertGameMineralPacketToCreateParams(source *cs.MineralPacket) generated.CreateMineralPacketParams

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:autoMap Cargo
	// goverter:map Cargo.Ironium Ironium
	// goverter:map Cargo.Boranium Boranium
	// goverter:map Cargo.Germanium Germanium
	// goverter:map Heading.X HeadingX
	// goverter:map Heading.Y HeadingY
	ConvertGameMineralPacketToUpdateParams(source *cs.MineralPacket) generated.UpdateMineralPacketParams

	// goverter:map . GameDBObject
	// goverter:map . MapObject | ExtendMineralPacketMapObject
	// goverter:map . Cargo
	// goverter:map . Heading | ExtendMineralPacketHeading
	ConvertMineralPacket(source generated.MineralPacket) *cs.MineralPacket
	ConvertMineralPackets(source []generated.MineralPacket) []*cs.MineralPacket

	// goverter:ignore Colonists
	salvageCargo(source generated.Salvage) cs.Cargo

	// goverter:ignore Colonists
	mineralPacketCargo(source generated.MineralPacket) cs.Cargo
}

func MapObjectTypeWormhole() cs.MapObjectType {
	return cs.MapObjectTypeWormhole
}

func NullTimeToTime(source sql.NullTime) time.Time {
	if source.Valid {
		return source.Time
	}
	return time.Time{}
}

func NullBoolToBool(source sql.NullBool) bool {
	if source.Valid {
		return source.Bool
	}
	return false
}

func NullStringToString(source sql.NullString) string {
	if source.Valid {
		return source.String
	}
	return ""
}

func Float64ToNullFloat64(source float64) sql.NullFloat64 {
	return sql.NullFloat64{
		Valid:   true,
		Float64: source,
	}
}

func NullInt64ToInt64(source sql.NullInt64) int64 {
	if source.Valid {
		return source.Int64
	}
	return 0
}

func NullInt64ToInt(source sql.NullInt64) int {
	if source.Valid {
		return int(source.Int64)
	}
	return 0
}

func IntToNullInt64(source int) sql.NullInt64 {
	return sql.NullInt64{
		Valid: true,
		Int64: int64(source),
	}
}

func Float64ToInt(source float64) int {
	return int(source)
}

func IntToFloat64(source int) float64 {
	return float64(source)
}

func Int64ToInt(source int64) int {
	return int(source)
}

func IntToInt64(source int) int64 {
	return int64(source)
}

func NullFloat64ToInt(source sql.NullFloat64) int {
	if source.Valid {
		return int(source.Float64)
	}
	return 0
}

func IntToNullFloat64(source int) sql.NullFloat64 {
	return sql.NullFloat64{
		Valid:   true,
		Float64: float64(source),
	}
}

func TagsToGameTags(source *generated.Tags) cs.Tags {
	if source == nil {
		return cs.Tags{}
	}
	return (cs.Tags)(*source)
}

func GameTagsToTags(source cs.Tags) *generated.Tags {
	return (*generated.Tags)(&source)
}

func CargoTransfersToGameCargoTransfers(source *generated.CargoTransfers) cs.CargoTransfers {
	// return an empty slice for nil
	if source == nil {
		return cs.CargoTransfers{}
	}

	return (cs.CargoTransfers)(*source)
}

func GameCargoTransfersToCargoTransfers(source cs.CargoTransfers) *generated.CargoTransfers {
	return (*generated.CargoTransfers)(&source)
}

func BattlePlansToGameBattlePlans(source *generated.BattlePlans) []cs.BattlePlan {
	// return an empty slice for nil
	if source == nil {
		return []cs.BattlePlan{}
	}

	return ([]cs.BattlePlan)(*source)
}

func GameBattlePlansToBattlePlans(source []cs.BattlePlan) *generated.BattlePlans {
	return (*generated.BattlePlans)(&source)
}

func ProductionPlansToGameProductionPlans(source *generated.ProductionPlans) []cs.ProductionPlan {
	// return an empty slice for nil
	if source == nil {
		return []cs.ProductionPlan{}
	}
	return ([]cs.ProductionPlan)(*source)
}

func GameProductionPlansToProductionPlans(source []cs.ProductionPlan) *generated.ProductionPlans {
	return (*generated.ProductionPlans)(&source)
}

func TransportPlansToGameTransportPlans(source *generated.TransportPlans) []cs.TransportPlan {
	// return an empty slice for nil
	if source == nil {
		return []cs.TransportPlan{}
	}
	return ([]cs.TransportPlan)(*source)
}

func GameTransportPlansToTransportPlans(source []cs.TransportPlan) *generated.TransportPlans {
	return (*generated.TransportPlans)(&source)
}

func PlayerRelationshipsToGamePlayerRelationships(source *generated.PlayerRelationships) []cs.PlayerRelationship {
	// return an empty slice for nil
	if source == nil {
		return []cs.PlayerRelationship{}
	}
	return ([]cs.PlayerRelationship)(*source)
}

func GamePlayerRelationshipsToPlayerRelationships(source []cs.PlayerRelationship) *generated.PlayerRelationships {
	return (*generated.PlayerRelationships)(&source)
}

func PlayerMessagesToGamePlayerMessages(source *generated.PlayerMessages) []cs.PlayerMessage {
	// return an empty slice for nil
	if source == nil {
		return []cs.PlayerMessage{}
	}
	return ([]cs.PlayerMessage)(*source)
}

func GamePlayerMessagesToPlayerMessages(source []cs.PlayerMessage) *generated.PlayerMessages {
	return (*generated.PlayerMessages)(&source)
}

func PlayerScoresToGamePlayerScores(source *generated.PlayerScores) []cs.PlayerScore {
	// return an empty slice for nil
	if source == nil {
		return []cs.PlayerScore{}
	}
	return ([]cs.PlayerScore)(*source)
}

func GamePlayerScoresToPlayerScores(source []cs.PlayerScore) *generated.PlayerScores {
	return (*generated.PlayerScores)(&source)
}

func AcquiredTechsToGameAcquiredTechs(source *generated.AcquiredTechs) map[string]bool {
	// return an empty slice for nil
	if source == nil {
		return map[string]bool{}
	}
	return (map[string]bool)(*source)
}

func GameAcquiredTechsToAcquiredTechs(source map[string]bool) *generated.AcquiredTechs {
	return (*generated.AcquiredTechs)(&source)
}

func BattleRecordsToGameBattleRecords(source *generated.BattleRecords) []cs.BattleRecord {
	// return an empty slice for nil
	if source == nil {
		return []cs.BattleRecord{}
	}
	return ([]cs.BattleRecord)(*source)
}

func GameBattleRecordsToBattleRecords(source []cs.BattleRecord) *generated.BattleRecords {
	return (*generated.BattleRecords)(&source)
}

func PlayerIntelsToGamePlayerIntels(source *generated.PlayerIntels) []cs.PlayerIntel {
	// return an empty slice for nil
	if source == nil {
		return []cs.PlayerIntel{}
	}
	return ([]cs.PlayerIntel)(*source)
}

func GamePlayerIntelsToPlayerIntels(source []cs.PlayerIntel) *generated.PlayerIntels {
	return (*generated.PlayerIntels)(&source)
}

func ScoreIntelsToGameScoreIntels(source *generated.ScoreIntels) []cs.ScoreIntel {
	// return an empty slice for nil
	if source == nil {
		return []cs.ScoreIntel{}
	}
	return ([]cs.ScoreIntel)(*source)
}

func GameScoreIntelsToScoreIntels(source []cs.ScoreIntel) *generated.ScoreIntels {
	return (*generated.ScoreIntels)(&source)
}

func PlanetIntelsToGamePlanetIntels(source *generated.PlanetIntels) []*cs.Planet {
	// return an empty slice for nil
	if source == nil {
		return []*cs.Planet{}
	}
	return ([]*cs.Planet)(*source)
}

func GamePlanetIntelsToPlanetIntels(source []*cs.Planet) *generated.PlanetIntels {
	return (*generated.PlanetIntels)(&source)
}

func FleetIntelsToGameFleetIntels(source *generated.FleetIntels) []*cs.Fleet {
	// return an empty slice for nil
	if source == nil {
		return []*cs.Fleet{}
	}
	return ([]*cs.Fleet)(*source)
}

func GameFleetIntelsToFleetIntels(source []*cs.Fleet) *generated.FleetIntels {
	return (*generated.FleetIntels)(&source)
}

func ShipDesignIntelsToGameShipDesignIntels(source *generated.ShipDesignIntels) []*cs.ShipDesign {
	// return an empty slice for nil
	if source == nil {
		return []*cs.ShipDesign{}
	}
	return ([]*cs.ShipDesign)(*source)
}

func GameShipDesignIntelsToShipDesignIntels(source []*cs.ShipDesign) *generated.ShipDesignIntels {
	return (*generated.ShipDesignIntels)(&source)
}

func MineralPacketIntelsToGameMineralPacketIntels(source *generated.MineralPacketIntels) []*cs.MineralPacket {
	// return an empty slice for nil
	if source == nil {
		return []*cs.MineralPacket{}
	}
	return ([]*cs.MineralPacket)(*source)
}

func GameMineralPacketIntelsToMineralPacketIntels(source []*cs.MineralPacket) *generated.MineralPacketIntels {
	return (*generated.MineralPacketIntels)(&source)
}

func SalvageIntelsToGameSalvageIntels(source *generated.SalvageIntels) []*cs.Salvage {
	// return an empty slice for nil
	if source == nil {
		return []*cs.Salvage{}
	}
	return ([]*cs.Salvage)(*source)
}

func GameSalvageIntelsToSalvageIntels(source []*cs.Salvage) *generated.SalvageIntels {
	return (*generated.SalvageIntels)(&source)
}

func MinefieldIntelsToGameMinefieldIntels(source *generated.MinefieldIntels) []*cs.Minefield {
	// return an empty slice for nil
	if source == nil {
		return []*cs.Minefield{}
	}
	return ([]*cs.Minefield)(*source)
}

func GameMinefieldIntelsToMinefieldIntels(source []*cs.Minefield) *generated.MinefieldIntels {
	return (*generated.MinefieldIntels)(&source)
}

func WormholeIntelsToGameWormholeIntels(source *generated.WormholeIntels) []*cs.Wormhole {
	// return an empty slice for nil
	if source == nil {
		return []*cs.Wormhole{}
	}
	return ([]*cs.Wormhole)(*source)
}

func GameWormholeIntelsToWormholeIntels(source []*cs.Wormhole) *generated.WormholeIntels {
	return (*generated.WormholeIntels)(&source)
}

func MysteryTraderIntelsToGameMysteryTraderIntels(source *generated.MysteryTraderIntels) []*cs.MysteryTrader {
	// return an empty slice for nil
	if source == nil {
		return []*cs.MysteryTrader{}
	}
	return ([]*cs.MysteryTrader)(*source)
}

func GameMysteryTraderIntelsToMysteryTraderIntels(source []*cs.MysteryTrader) *generated.MysteryTraderIntels {
	return (*generated.MysteryTraderIntels)(&source)
}

func PlayerRaceToGameRace(source *generated.PlayerRace) cs.Race {
	// return an empty object for nil to support partial loads
	if source == nil {
		return cs.Race{}
	}
	return cs.Race(*source)
}

func GameRaceToPlayerRace(source cs.Race) *generated.PlayerRace {
	return (*generated.PlayerRace)(&source)
}

func PlayerSpecToGamePlayerSpec(source *generated.PlayerSpec) cs.PlayerSpec {
	// return an empty object for nil to support partial loads
	if source == nil {
		return cs.PlayerSpec{}
	}

	return (cs.PlayerSpec)(*source)
}

func GamePlayerSpecToPlayerSpec(source cs.PlayerSpec) *generated.PlayerSpec {
	return (*generated.PlayerSpec)(&source)
}

func PlayerStatsToGamePlayerStats(source *generated.PlayerStats) *cs.PlayerStats {
	return (*cs.PlayerStats)(source)
}

func GamePlayerStatsToPlayerStats(source *cs.PlayerStats) *generated.PlayerStats {
	return (*generated.PlayerStats)(source)
}

func PlanetSpecToGamePlanetSpec(source *generated.PlanetSpec) cs.PlanetSpec {
	return (cs.PlanetSpec)(*source)
}

func GamePlanetSpecToPlanetSpec(source cs.PlanetSpec) *generated.PlanetSpec {
	return (*generated.PlanetSpec)(&source)
}

func ProductionQueueItemsToGameProductionQueueItems(source *generated.ProductionQueueItems) []cs.ProductionQueueItem {
	return ([]cs.ProductionQueueItem)(*source)
}

func GameProductionQueueItemsToProductionQueueItems(source []cs.ProductionQueueItem) *generated.ProductionQueueItems {
	return (*generated.ProductionQueueItems)(&source)
}

func FleetSpecToGameFleetSpec(source *generated.FleetSpec) cs.FleetSpec {
	return (cs.FleetSpec)(*source)
}

func GameFleetSpecToFleetSpec(source cs.FleetSpec) *generated.FleetSpec {
	return (*generated.FleetSpec)(&source)
}

func ShipTokensToGameShipTokens(source *generated.ShipTokens) []cs.ShipToken {
	if source == nil {
		return nil
	}
	return ([]cs.ShipToken)(*source)
}

func GameShipTokensToShipTokens(source []cs.ShipToken) *generated.ShipTokens {
	return (*generated.ShipTokens)(&source)
}

func WaypointsToGameWaypoints(source *generated.Waypoints) []cs.Waypoint {
	return ([]cs.Waypoint)(*source)
}

func GameWaypointsToWaypoints(source []cs.Waypoint) *generated.Waypoints {
	return (*generated.Waypoints)(&source)
}

func ShipDesignSpecToGameShipDesignSpec(source *generated.ShipDesignSpec) cs.ShipDesignSpec {
	return (cs.ShipDesignSpec)(*source)
}

func GameShipDesignSpecToShipDesignSpec(source cs.ShipDesignSpec) *generated.ShipDesignSpec {
	return (*generated.ShipDesignSpec)(&source)
}

func ShipDesignSlotsToGameShipDesignSlots(source *generated.ShipDesignSlots) []cs.ShipDesignSlot {
	return ([]cs.ShipDesignSlot)(*source)
}

func GameShipDesignSlotsToShipDesignSlots(source []cs.ShipDesignSlot) *generated.ShipDesignSlots {
	return (*generated.ShipDesignSlots)(&source)
}

func WormholeSpecToGameWormholeSpec(source *generated.WormholeSpec) cs.WormholeSpec {
	return (cs.WormholeSpec)(*source)
}

func GameWormholeSpecToWormholeSpec(source cs.WormholeSpec) *generated.WormholeSpec {
	return (*generated.WormholeSpec)(&source)
}

func MysteryTraderSpecToGameMysteryTraderSpec(source *generated.MysteryTraderSpec) cs.MysteryTraderSpec {
	return (cs.MysteryTraderSpec)(*source)
}

func GameMysteryTraderSpecToMysteryTraderSpec(source cs.MysteryTraderSpec) *generated.MysteryTraderSpec {
	return (*generated.MysteryTraderSpec)(&source)
}

func MysteryTraderPlayersRewardedToGameMysteryTraderPlayersRewarded(source *generated.MysteryTraderPlayersRewarded) map[int]bool {
	// return an empty slice for nil
	if source == nil {
		return map[int]bool{}
	}
	return (map[int]bool)(*source)
}

func GameMysteryTraderPlayersRewardedToMysteryTraderPlayersRewarded(source map[int]bool) *generated.MysteryTraderPlayersRewarded {
	return (*generated.MysteryTraderPlayersRewarded)(&source)
}

func ExtendResearchCost(source generated.Race) cs.ResearchCost {
	return cs.ResearchCost{
		Energy:        cs.ResearchCostLevel(source.ResearchCostEnergy),
		Weapons:       cs.ResearchCostLevel(source.ResearchCostWeapons),
		Propulsion:    cs.ResearchCostLevel(source.ResearchCostPropulsion),
		Construction:  cs.ResearchCostLevel(source.ResearchCostConstruction),
		Electronics:   cs.ResearchCostLevel(source.ResearchCostElectronics),
		Biotechnology: cs.ResearchCostLevel(source.ResearchCostBiotechnology),
	}
}

func ExtendHabLow(source generated.Race) cs.Hab {
	return cs.Hab{
		Grav: int(source.HabLowGrav),
		Temp: int(source.HabLowTemp),
		Rad:  int(source.HabLowRad),
	}
}

func ExtendHabHigh(source generated.Race) cs.Hab {
	return cs.Hab{
		Grav: int(source.HabHighGrav),
		Temp: int(source.HabHighTemp),
		Rad:  int(source.HabHighRad),
	}
}
func ExtendVictoryConditions(source generated.Game) cs.VictoryConditions {
	return cs.VictoryConditions{
		Conditions:               source.VictoryConditionsConditions,
		NumCriteriaRequired:      int(source.VictoryConditionsNumCriteriaRequired),
		YearsPassed:              int(source.VictoryConditionsYearsPassed),
		OwnPlanets:               int(source.VictoryConditionsOwnPlanets),
		AttainTechLevel:          int(source.VictoryConditionsAttainTechLevel),
		AttainTechLevelNumFields: int(source.VictoryConditionsAttainTechLevelNumFields),
		ExceedsScore:             int(source.VictoryConditionsExceedsScore),
		ExceedsSecondPlaceScore:  int(source.VictoryConditionsExceedsSecondPlaceScore),
		ProductionCapacity:       int(source.VictoryConditionsProductionCapacity),
		OwnCapitalShips:          int(source.VictoryConditionsOwnCapitalShips),
		HighestScoreAfterYears:   int(source.VictoryConditionsHighestScoreAfterYears),
	}
}

func ExtendArea(source generated.Game) cs.Vector {
	return cs.Vector{
		X: int(source.AreaX),
		Y: int(source.AreaY),
	}
}

func ExtendTechLevels(source generated.Player) cs.TechLevel {
	return cs.TechLevel{
		Energy:        int(source.TechLevelsEnergy),
		Weapons:       int(source.TechLevelsWeapons),
		Propulsion:    int(source.TechLevelsPropulsion),
		Construction:  int(source.TechLevelsConstruction),
		Electronics:   int(source.TechLevelsElectronics),
		Biotechnology: int(source.TechLevelsBiotechnology),
	}
}

func ExtendTechLevelsSpent(source generated.Player) cs.TechLevel {
	return cs.TechLevel{
		Energy:        int(source.TechLevelsSpentEnergy),
		Weapons:       int(source.TechLevelsSpentWeapons),
		Propulsion:    int(source.TechLevelsSpentPropulsion),
		Construction:  int(source.TechLevelsSpentConstruction),
		Electronics:   int(source.TechLevelsSpentElectronics),
		Biotechnology: int(source.TechLevelsSpentBiotechnology),
	}
}

func ExtendPlayerGameDBObject(source generated.Player) cs.GameDBObject {
	return cs.GameDBObject{
		ID:        source.ID,
		GameID:    source.GameID,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
	}
}

func ExtendLightPlayerGameDBObject(source generated.GetLightPlayerForGameRow) cs.GameDBObject {
	return cs.GameDBObject{
		ID:        source.ID,
		GameID:    source.GameID,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
	}
}

func ExtendPlayerStatusGameDBObject(source generated.GetPlayersStatusForGameRow) cs.GameDBObject {
	return cs.GameDBObject{
		ID:        source.ID,
		GameID:    source.GameID,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
	}
}

func ExtendTechLevelsLight(source generated.GetLightPlayerForGameRow) cs.TechLevel {
	return cs.TechLevel{
		Energy:        int(source.TechLevelsEnergy),
		Weapons:       int(source.TechLevelsWeapons),
		Propulsion:    int(source.TechLevelsPropulsion),
		Construction:  int(source.TechLevelsConstruction),
		Electronics:   int(source.TechLevelsElectronics),
		Biotechnology: int(source.TechLevelsBiotechnology),
	}
}

func ExtendTechLevelsSpentLight(source generated.GetLightPlayerForGameRow) cs.TechLevel {
	return cs.TechLevel{
		Energy:        int(source.TechLevelsSpentEnergy),
		Weapons:       int(source.TechLevelsSpentWeapons),
		Propulsion:    int(source.TechLevelsSpentPropulsion),
		Construction:  int(source.TechLevelsSpentConstruction),
		Electronics:   int(source.TechLevelsSpentElectronics),
		Biotechnology: int(source.TechLevelsSpentBiotechnology),
	}
}

func ExtendPlanetGameDBObject(source generated.Planet) cs.GameDBObject {
	return cs.GameDBObject{
		ID:        source.ID,
		GameID:    source.GameID,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
	}
}

func ExtendPlanetMapObject(source generated.Planet) cs.MapObject {
	return cs.MapObject{
		Type: cs.MapObjectTypePlanet,
		Position: cs.Vector{
			X: int(source.X),
			Y: int(source.Y),
		},
		Name:      source.Name,
		Num:       int(source.Num),
		PlayerNum: int(source.PlayerNum),
		Tags:      TagsToGameTags(source.Tags),
	}
}

func ExtendBaseHab(source generated.Planet) cs.Hab {
	return cs.Hab{
		Grav: int(source.BaseGrav),
		Temp: int(source.BaseTemp),
		Rad:  int(source.BaseRad),
	}
}

func ExtendTerraformedAmount(source generated.Planet) cs.Hab {
	return cs.Hab{
		Grav: int(source.TerraformedAmountGrav),
		Temp: int(source.TerraformedAmountTemp),
		Rad:  int(source.TerraformedAmountRad),
	}
}

func ExtendMineralConcentration(source generated.Planet) cs.Mineral {
	return cs.Mineral{
		Ironium:   int(source.MineralConcIronium),
		Boranium:  int(source.MineralConcBoranium),
		Germanium: int(source.MineralConcGermanium),
	}
}

func ExtendMineYears(source generated.Planet) cs.Mineral {
	return cs.Mineral{
		Ironium:   int(source.MineYearsIronium),
		Boranium:  int(source.MineYearsBoranium),
		Germanium: int(source.MineYearsGermanium),
	}
}

func ExtendFleetMapObject(source generated.Fleet) cs.MapObject {
	return cs.MapObject{
		Type: cs.MapObjectTypeFleet,
		Position: cs.Vector{
			X: int(source.X),
			Y: int(source.Y),
		},
		Name:      source.Name,
		Num:       int(source.Num),
		PlayerNum: int(source.PlayerNum),
		Tags:      TagsToGameTags(source.Tags),
	}
}

func ExtendFleetFleetOrders(source generated.Fleet) cs.FleetOrders {
	return cs.FleetOrders{
		BattlePlanNum: int(source.BattlePlanNum),
		Waypoints:     *source.Waypoints,
		RepeatOrders:  source.RepeatOrders,
		Purpose:       *source.Purpose,
	}
}

func ExtendFleetHeading(source generated.Fleet) cs.VectorFloat64 {
	return cs.VectorFloat64{
		X: source.HeadingX,
		Y: source.HeadingY,
	}
}

func ExtendFleetPreviousPosition(source generated.Fleet) *cs.Vector {
	if !source.PreviousPositionX.Valid || !source.PreviousPositionY.Valid {
		return nil
	}
	return &cs.Vector{
		X: int(source.PreviousPositionX.Float64),
		Y: int(source.PreviousPositionY.Float64),
	}
}

func ExtendMysteryTraderMapObject(source generated.MysteryTrader) cs.MapObject {
	return cs.MapObject{
		Type: cs.MapObjectTypeMysteryTrader,
		Position: cs.Vector{
			X: int(source.X),
			Y: int(source.Y),
		},
		Name: source.Name,
		Num:  int(source.Num),
		Tags: TagsToGameTags(source.Tags),
	}
}

func ExtendMysteryTraderHeading(source generated.MysteryTrader) cs.VectorFloat64 {
	return cs.VectorFloat64{
		X: source.HeadingX,
		Y: source.HeadingY,
	}
}

func ExtendMysteryTraderDestination(source generated.MysteryTrader) cs.Vector {
	return cs.Vector{
		X: int(source.DestinationX),
		Y: int(source.DestinationY),
	}
}

func ExtendSalvageMapObject(source generated.Salvage) cs.MapObject {
	return cs.MapObject{
		Type: cs.MapObjectTypeSalvage,
		Position: cs.Vector{
			X: int(source.X),
			Y: int(source.Y),
		},
		Name:      source.Name,
		Num:       int(source.Num),
		PlayerNum: int(source.PlayerNum),
		Tags:      TagsToGameTags(source.Tags),
	}
}

func ExtendMinefieldMapObject(source generated.Minefield) cs.MapObject {
	return cs.MapObject{
		Type: cs.MapObjectTypeMinefield,
		Position: cs.Vector{
			X: int(source.X),
			Y: int(source.Y),
		},
		Name:      source.Name,
		Num:       int(source.Num),
		PlayerNum: int(source.PlayerNum),
		Tags:      TagsToGameTags(source.Tags),
	}
}

func ExtendMineralPacketHeading(source generated.MineralPacket) cs.VectorFloat64 {
	return cs.VectorFloat64{
		X: source.HeadingX,
		Y: source.HeadingY,
	}
}

func ExtendMineralPacketMapObject(source generated.MineralPacket) cs.MapObject {
	return cs.MapObject{
		Type: cs.MapObjectTypeMineralPacket,
		Position: cs.Vector{
			X: int(source.X),
			Y: int(source.Y),
		},
		Name:      source.Name,
		Num:       int(source.Num),
		PlayerNum: int(source.PlayerNum),
		Tags:      TagsToGameTags(source.Tags),
	}
}
