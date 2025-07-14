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
// goverter:output:package github.com/sirgwain/craig-stars/db
// goverter:output:file ./converter.generated.go
// goverter:ignoreUnexported
// goverter:extend AcquiredTechsToGameAcquiredTechs
// goverter:extend BattlePlansToGameBattlePlans
// goverter:extend BattleRecordsToGameBattleRecords
// goverter:extend BoolToNullBool
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
// goverter:extend GameMineFieldIntelsToMineFieldIntels
// goverter:extend GameMineFieldSpecToMineFieldSpec
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
// goverter:extend GameRaceSpecToRaceSpec
// goverter:extend GameRaceToPlayerRace
// goverter:extend GameRulesToRules
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
// goverter:extend Int64ToNullInt64
// goverter:extend IntToNullInt64
// goverter:extend IntToInt64
// goverter:extend MineFieldIntelsToGameMineFieldIntels
// goverter:extend MineFieldSpecToGameMineFieldSpec
// goverter:extend MineralPacketIntelsToGameMineralPacketIntels
// goverter:extend MysteryTraderIntelsToGameMysteryTraderIntels
// goverter:extend MysteryTraderPlayersRewardedToGameMysteryTraderPlayersRewarded
// goverter:extend MysteryTraderSpecToGameMysteryTraderSpec
// goverter:extend NullBoolToBool
// goverter:extend NullFloat64ToFloat64
// goverter:extend NullInt64ToInt64
// goverter:extend NullInt64ToInt
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
// goverter:extend RaceSpecToGameRaceSpec
// goverter:extend RulesToGameRules
// goverter:extend SalvageIntelsToGameSalvageIntels
// goverter:extend ScoreIntelsToGameScoreIntels
// goverter:extend ShipDesignIntelsToGameShipDesignIntels
// goverter:extend ShipDesignSlotsToGameShipDesignSlots
// goverter:extend ShipDesignSpecToGameShipDesignSpec
// goverter:extend ShipTokensToGameShipTokens
// goverter:extend StringToNullString
// goverter:extend TagsToGameTags
// goverter:extend TimeToNullTime
// goverter:extend TransportPlansToGameTransportPlans
// goverter:extend WaypointsToGameWaypoints
// goverter:extend WormholeIntelsToGameWormholeIntels
// goverter:extend WormholeSpecToGameWormholeSpec
// goverter:enum no
// goverter:matchIgnoreCase
// goverter:useZeroValueOnPointerInconsistency
// goverter:useUnderlyingTypeMethods
// goverter:skipCopySameType
// goverter:name GameConverter
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
	ConvertRace(source generated.Race) cs.Race

	ConvertRaces(source []generated.Race) []cs.Race

	// goverter:autoMap DBObject
	// goverter:map HabHigh.Grav Habhighgrav
	// goverter:map HabHigh.Temp Habhightemp
	// goverter:map HabHigh.Rad Habhighrad
	// goverter:map HabLow.Grav Hablowgrav
	// goverter:map HabLow.Temp Hablowtemp
	// goverter:map HabLow.Rad Hablowrad
	// goverter:map ResearchCost.Energy Researchcostenergy
	// goverter:map ResearchCost.Weapons Researchcostweapons
	// goverter:map ResearchCost.Propulsion Researchcostpropulsion
	// goverter:map ResearchCost.Construction Researchcostconstruction
	// goverter:map ResearchCost.Electronics Researchcostelectronics
	// goverter:map ResearchCost.Biotechnology Researchcostbiotechnology
	ConvertGameRace(source *cs.Race) generated.Race

	// goverter:autoMap DBObject
	// goverter:map HabHigh.Grav Habhighgrav
	// goverter:map HabHigh.Temp Habhightemp
	// goverter:map HabHigh.Rad Habhighrad
	// goverter:map HabLow.Grav Hablowgrav
	// goverter:map HabLow.Temp Hablowtemp
	// goverter:map HabLow.Rad Hablowrad
	// goverter:map ResearchCost.Energy Researchcostenergy
	// goverter:map ResearchCost.Weapons Researchcostweapons
	// goverter:map ResearchCost.Propulsion Researchcostpropulsion
	// goverter:map ResearchCost.Construction Researchcostconstruction
	// goverter:map ResearchCost.Electronics Researchcostelectronics
	// goverter:map ResearchCost.Biotechnology Researchcostbiotechnology
	ConvertGameRaceToCreateParams(source *cs.Race) generated.CreateRaceParams

	// goverter:autoMap DBObject
	// goverter:map HabHigh.Grav Habhighgrav
	// goverter:map HabHigh.Temp Habhightemp
	// goverter:map HabHigh.Rad Habhighrad
	// goverter:map HabLow.Grav Hablowgrav
	// goverter:map HabLow.Temp Hablowtemp
	// goverter:map HabLow.Rad Hablowrad
	// goverter:map ResearchCost.Energy Researchcostenergy
	// goverter:map ResearchCost.Weapons Researchcostweapons
	// goverter:map ResearchCost.Propulsion Researchcostpropulsion
	// goverter:map ResearchCost.Construction Researchcostconstruction
	// goverter:map ResearchCost.Electronics Researchcostelectronics
	// goverter:map ResearchCost.Biotechnology Researchcostbiotechnology
	ConvertGameRaceToUpdateParams(source *cs.Race) generated.UpdateRaceParams

	// goverter:map . DBObject
	// goverter:map . VictoryConditions | ExtendVictoryConditions
	// goverter:map . Area | ExtendArea
	// goverter:ignore Rules
	ConvertGame(source generated.Game) cs.Game

	ConvertGames(source []generated.Game) []cs.Game

	// goverter:autoMap DBObject
	// goverter:map VictoryConditions.NumCriteriaRequired Victoryconditionsnumcriteriarequired
	// goverter:map VictoryConditions.YearsPassed Victoryconditionsyearspassed
	// goverter:map VictoryConditions.OwnPlanets Victoryconditionsownplanets
	// goverter:map VictoryConditions.AttainTechLevel Victoryconditionsattaintechlevel
	// goverter:map VictoryConditions.AttainTechLevelNumFields Victoryconditionsattaintechlevelnumfields
	// goverter:map VictoryConditions.ExceedsScore Victoryconditionsexceedsscore
	// goverter:map VictoryConditions.ExceedsSecondPlaceScore Victoryconditionsexceedssecondplacescore
	// goverter:map VictoryConditions.ProductionCapacity Victoryconditionsproductioncapacity
	// goverter:map VictoryConditions.OwnCapitalShips Victoryconditionsowncapitalships
	// goverter:map VictoryConditions.HighestScoreAfterYears Victoryconditionshighestscoreafteryears
	// goverter:map VictoryConditions.Conditions Victoryconditionsconditions
	// goverter:map Area.X Areax
	// goverter:map Area.Y Areay
	ConvertGameGame(source *cs.Game) generated.Game

	// goverter:autoMap DBObject
	// goverter:map VictoryConditions.NumCriteriaRequired Victoryconditionsnumcriteriarequired
	// goverter:map VictoryConditions.YearsPassed Victoryconditionsyearspassed
	// goverter:map VictoryConditions.OwnPlanets Victoryconditionsownplanets
	// goverter:map VictoryConditions.AttainTechLevel Victoryconditionsattaintechlevel
	// goverter:map VictoryConditions.AttainTechLevelNumFields Victoryconditionsattaintechlevelnumfields
	// goverter:map VictoryConditions.ExceedsScore Victoryconditionsexceedsscore
	// goverter:map VictoryConditions.ExceedsSecondPlaceScore Victoryconditionsexceedssecondplacescore
	// goverter:map VictoryConditions.ProductionCapacity Victoryconditionsproductioncapacity
	// goverter:map VictoryConditions.OwnCapitalShips Victoryconditionsowncapitalships
	// goverter:map VictoryConditions.HighestScoreAfterYears Victoryconditionshighestscoreafteryears
	// goverter:map VictoryConditions.Conditions Victoryconditionsconditions
	// goverter:map Area.X Areax
	// goverter:map Area.Y Areay
	ConvertGameGameToCreateParams(source *cs.Game) generated.CreateGameParams

	// goverter:autoMap DBObject
	// goverter:map VictoryConditions.NumCriteriaRequired Victoryconditionsnumcriteriarequired
	// goverter:map VictoryConditions.YearsPassed Victoryconditionsyearspassed
	// goverter:map VictoryConditions.OwnPlanets Victoryconditionsownplanets
	// goverter:map VictoryConditions.AttainTechLevel Victoryconditionsattaintechlevel
	// goverter:map VictoryConditions.AttainTechLevelNumFields Victoryconditionsattaintechlevelnumfields
	// goverter:map VictoryConditions.ExceedsScore Victoryconditionsexceedsscore
	// goverter:map VictoryConditions.ExceedsSecondPlaceScore Victoryconditionsexceedssecondplacescore
	// goverter:map VictoryConditions.ProductionCapacity Victoryconditionsproductioncapacity
	// goverter:map VictoryConditions.OwnCapitalShips Victoryconditionsowncapitalships
	// goverter:map VictoryConditions.HighestScoreAfterYears Victoryconditionshighestscoreafteryears
	// goverter:map VictoryConditions.Conditions Victoryconditionsconditions
	// goverter:map Area.X Areax
	// goverter:map Area.Y Areay
	ConvertGameGameToUpdateParams(source *cs.Game) generated.UpdateGameParams

	// goverter:autoMap GameDBObject
	// goverter:map TechLevels.Energy Techlevelsenergy
	// goverter:map TechLevels.Weapons Techlevelsweapons
	// goverter:map TechLevels.Propulsion Techlevelspropulsion
	// goverter:map TechLevels.Construction Techlevelsconstruction
	// goverter:map TechLevels.Electronics Techlevelselectronics
	// goverter:map TechLevels.Biotechnology Techlevelsbiotechnology
	// goverter:map TechLevelsSpent.Energy Techlevelsspentenergy
	// goverter:map TechLevelsSpent.Weapons Techlevelsspentweapons
	// goverter:map TechLevelsSpent.Propulsion Techlevelsspentpropulsion
	// goverter:map TechLevelsSpent.Construction Techlevelsspentconstruction
	// goverter:map TechLevelsSpent.Electronics Techlevelsspentelectronics
	// goverter:map TechLevelsSpent.Biotechnology Techlevelsspentbiotechnology
	// goverter:autoMap PlayerOrders
	// goverter:map PlayerIntels.BattleRecords Battlerecords
	// goverter:map PlayerIntels.PlayerIntels Playerintels
	// goverter:map PlayerIntels.ScoreIntels Scoreintels
	// goverter:map PlayerIntels.PlanetIntels Planetintels
	// goverter:map PlayerIntels.FleetIntels Fleetintels
	// goverter:map PlayerIntels.ShipDesignIntels Shipdesignintels
	// goverter:map PlayerIntels.MineralPacketIntels Mineralpacketintels
	// goverter:map PlayerIntels.MineFieldIntels Minefieldintels
	// goverter:map PlayerIntels.WormholeIntels Wormholeintels
	// goverter:map PlayerIntels.MysteryTraderIntels Mysterytraderintels
	// goverter:map PlayerIntels.SalvageIntels	 Salvageintels
	// goverter:autoMap PlayerPlans
	ConvertGamePlayer(source *cs.Player) generated.Player

	// goverter:autoMap GameDBObject
	// goverter:map TechLevels.Energy Techlevelsenergy
	// goverter:map TechLevels.Weapons Techlevelsweapons
	// goverter:map TechLevels.Propulsion Techlevelspropulsion
	// goverter:map TechLevels.Construction Techlevelsconstruction
	// goverter:map TechLevels.Electronics Techlevelselectronics
	// goverter:map TechLevels.Biotechnology Techlevelsbiotechnology
	// goverter:map TechLevelsSpent.Energy Techlevelsspentenergy
	// goverter:map TechLevelsSpent.Weapons Techlevelsspentweapons
	// goverter:map TechLevelsSpent.Propulsion Techlevelsspentpropulsion
	// goverter:map TechLevelsSpent.Construction Techlevelsspentconstruction
	// goverter:map TechLevelsSpent.Electronics Techlevelsspentelectronics
	// goverter:map TechLevelsSpent.Biotechnology Techlevelsspentbiotechnology
	// goverter:autoMap PlayerOrders
	// goverter:map PlayerIntels.BattleRecords Battlerecords
	// goverter:map PlayerIntels.PlayerIntels Playerintels
	// goverter:map PlayerIntels.ScoreIntels Scoreintels
	// goverter:map PlayerIntels.PlanetIntels Planetintels
	// goverter:map PlayerIntels.FleetIntels Fleetintels
	// goverter:map PlayerIntels.ShipDesignIntels Shipdesignintels
	// goverter:map PlayerIntels.MineralPacketIntels Mineralpacketintels
	// goverter:map PlayerIntels.MineFieldIntels Minefieldintels
	// goverter:map PlayerIntels.WormholeIntels Wormholeintels
	// goverter:map PlayerIntels.MysteryTraderIntels Mysterytraderintels
	// goverter:map PlayerIntels.SalvageIntels	 Salvageintels
	// goverter:autoMap PlayerPlans
	ConvertGamePlayerToCreateParams(source *cs.Player) generated.CreatePlayerParams

	// goverter:autoMap GameDBObject
	// goverter:map TechLevels.Energy Techlevelsenergy
	// goverter:map TechLevels.Weapons Techlevelsweapons
	// goverter:map TechLevels.Propulsion Techlevelspropulsion
	// goverter:map TechLevels.Construction Techlevelsconstruction
	// goverter:map TechLevels.Electronics Techlevelselectronics
	// goverter:map TechLevels.Biotechnology Techlevelsbiotechnology
	// goverter:map TechLevelsSpent.Energy Techlevelsspentenergy
	// goverter:map TechLevelsSpent.Weapons Techlevelsspentweapons
	// goverter:map TechLevelsSpent.Propulsion Techlevelsspentpropulsion
	// goverter:map TechLevelsSpent.Construction Techlevelsspentconstruction
	// goverter:map TechLevelsSpent.Electronics Techlevelsspentelectronics
	// goverter:map TechLevelsSpent.Biotechnology Techlevelsspentbiotechnology
	// goverter:autoMap PlayerOrders
	// goverter:map PlayerIntels.BattleRecords Battlerecords
	// goverter:map PlayerIntels.PlayerIntels Playerintels
	// goverter:map PlayerIntels.ScoreIntels Scoreintels
	// goverter:map PlayerIntels.PlanetIntels Planetintels
	// goverter:map PlayerIntels.FleetIntels Fleetintels
	// goverter:map PlayerIntels.ShipDesignIntels Shipdesignintels
	// goverter:map PlayerIntels.MineralPacketIntels Mineralpacketintels
	// goverter:map PlayerIntels.MineFieldIntels Minefieldintels
	// goverter:map PlayerIntels.WormholeIntels Wormholeintels
	// goverter:map PlayerIntels.MysteryTraderIntels Mysterytraderintels
	// goverter:map PlayerIntels.SalvageIntels	 Salvageintels
	// goverter:autoMap PlayerPlans
	ConvertGamePlayerToUpdateParams(source *cs.Player) generated.UpdatePlayerParams

	// goverter:map . GameDBObject
	// goverter:map . TechLevels | ExtendTechLevels
	// goverter:map . TechLevelsSpent | ExtendTechLevelsSpent
	// goverter:map . PlayerOrders
	// goverter:map . PlayerIntels
	// goverter:map . PlayerPlans
	// goverter:ignore Designs
	ConvertPlayer(source generated.Player) cs.Player

	// goverter:map . GameDBObject
	// goverter:map . TechLevels | ExtendTechLevelsLight
	// goverter:map . TechLevelsSpent | ExtendTechLevelsSpentLight
	// goverter:map . PlayerOrders
	// goverter:map . PlayerPlans
	// goverter:ignore Messages
	// goverter:ignore PlayerIntels
	// goverter:ignore Designs
	ConvertLightPlayer(source generated.GetLightPlayerForGameRow) cs.Player

	// goverter:map . GameDBObject
	// goverter:ignoreMissing
	ConvertPlayerStatus(source generated.GetPlayersStatusForGameRow) cs.Player
	ConvertPlayerStatuses(source []generated.GetPlayersStatusForGameRow) []*cs.Player

	ConvertGetGamesWithPlayersRowToPlayerStatus(source generated.GetGamesWithPlayersRow) cs.PlayerStatus
	ConvertGetGamesWithPlayersForUserRowToPlayerStatus(source generated.GetGamesWithPlayersForUserRow) cs.PlayerStatus
	ConvertGetGameWithPlayersRowToPlayerStatus(source generated.GetGameWithPlayersRow) cs.PlayerStatus

	ConvertGetPlayerForGameRowToShipDesign(source generated.GetPlayerForGameRow) generated.Shipdesign
	ConvertGetPlayerForGameAndUserRowToShipDesign(source generated.GetPlayerForGameAndUserRow) generated.Shipdesign
	ConvertGetPlayersWithDesignsForGameRowToShipDesign(source generated.GetPlayersWithDesignsForGameRow) generated.Shipdesign

	ConvertPlayers(source []generated.Player) []*cs.Player

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:autoMap PlanetOrders
	// goverter:autoMap Hab
	// goverter:map BaseHab.Grav Basegrav
	// goverter:map BaseHab.Temp Basetemp
	// goverter:map BaseHab.Rad Baserad
	// goverter:map TerraformedAmount.Grav Terraformedamountgrav
	// goverter:map TerraformedAmount.Temp Terraformedamounttemp
	// goverter:map TerraformedAmount.Rad Terraformedamountrad
	// goverter:map MineralConcentration.Ironium Mineralconcironium
	// goverter:map MineralConcentration.Boranium Mineralconcboranium
	// goverter:map MineralConcentration.Germanium Mineralconcgermanium
	// goverter:map MineYears.Ironium Mineyearsironium
	// goverter:map MineYears.Boranium Mineyearsboranium
	// goverter:map MineYears.Germanium Mineyearsgermanium
	// goverter:autoMap Cargo
	ConvertGamePlanet(source *cs.Planet) generated.Planet

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:autoMap PlanetOrders
	// goverter:autoMap Hab
	// goverter:map BaseHab.Grav Basegrav
	// goverter:map BaseHab.Temp Basetemp
	// goverter:map BaseHab.Rad Baserad
	// goverter:map TerraformedAmount.Grav Terraformedamountgrav
	// goverter:map TerraformedAmount.Temp Terraformedamounttemp
	// goverter:map TerraformedAmount.Rad Terraformedamountrad
	// goverter:map MineralConcentration.Ironium Mineralconcironium
	// goverter:map MineralConcentration.Boranium Mineralconcboranium
	// goverter:map MineralConcentration.Germanium Mineralconcgermanium
	// goverter:map MineYears.Ironium Mineyearsironium
	// goverter:map MineYears.Boranium Mineyearsboranium
	// goverter:map MineYears.Germanium Mineyearsgermanium
	// goverter:autoMap Cargo
	ConvertGamePlanetToCreateParams(source *cs.Planet) generated.CreatePlanetParams

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:autoMap PlanetOrders
	// goverter:autoMap Hab
	// goverter:map BaseHab.Grav Basegrav
	// goverter:map BaseHab.Temp Basetemp
	// goverter:map BaseHab.Rad Baserad
	// goverter:map TerraformedAmount.Grav Terraformedamountgrav
	// goverter:map TerraformedAmount.Temp Terraformedamounttemp
	// goverter:map TerraformedAmount.Rad Terraformedamountrad
	// goverter:map MineralConcentration.Ironium Mineralconcironium
	// goverter:map MineralConcentration.Boranium Mineralconcboranium
	// goverter:map MineralConcentration.Germanium Mineralconcgermanium
	// goverter:map MineYears.Ironium Mineyearsironium
	// goverter:map MineYears.Boranium Mineyearsboranium
	// goverter:map MineYears.Germanium Mineyearsgermanium
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
	// goverter:map Heading.X Headingx
	// goverter:map Heading.Y Headingy
	// goverter:map PreviousPosition.X Previouspositionx
	// goverter:map PreviousPosition.Y Previouspositiony
	// goverter:autoMap Cargo
	ConvertGameFleet(source *cs.Fleet) generated.Fleet

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:autoMap FleetOrders
	// goverter:map Heading.X Headingx
	// goverter:map Heading.Y Headingy
	// goverter:map PreviousPosition.X Previouspositionx
	// goverter:map PreviousPosition.Y Previouspositiony
	// goverter:autoMap Cargo
	ConvertGameFleetToCreateParams(source *cs.Fleet) generated.CreateFleetParams

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:autoMap FleetOrders
	// goverter:map Heading.X Headingx
	// goverter:map Heading.Y Headingy
	// goverter:map PreviousPosition.X Previouspositionx
	// goverter:map PreviousPosition.Y Previouspositiony
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
	// goverter:ignore Candelete
	ConvertGameShipDesign(source *cs.ShipDesign) generated.Shipdesign

	// goverter:autoMap GameDBObject
	// goverter:ignore Candelete
	ConvertGameShipDesignToCreateParams(source *cs.ShipDesign) generated.CreateShipDesignParams

	// goverter:autoMap GameDBObject
	// goverter:ignore Candelete
	ConvertGameShipDesignToUpdateParams(source *cs.ShipDesign) generated.UpdateShipDesignParams

	// goverter:ignore Delete
	// goverter:map . GameDBObject
	ConvertShipDesign(source generated.Shipdesign) *cs.ShipDesign

	ConvertShipDesigns(source []generated.Shipdesign) []*cs.ShipDesign

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
	// goverter:map Heading.X Headingx
	// goverter:map Heading.Y Headingy
	// goverter:map Destination.X Destinationx
	// goverter:map Destination.Y Destinationy
	ConvertGameMysteryTrader(source *cs.MysteryTrader) generated.Mysterytrader

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:map Heading.X Headingx
	// goverter:map Heading.Y Headingy
	// goverter:map Destination.X Destinationx
	// goverter:map Destination.Y Destinationy
	ConvertGameMysteryTraderToCreateParams(source *cs.MysteryTrader) generated.CreateMysteryTraderParams

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:map Heading.X Headingx
	// goverter:map Heading.Y Headingy
	// goverter:map Destination.X Destinationx
	// goverter:map Destination.Y Destinationy
	ConvertGameMysteryTraderToUpdateParams(source *cs.MysteryTrader) generated.UpdateMysteryTraderParams

	// goverter:map . GameDBObject
	// goverter:map . MapObject | ExtendMysteryTraderMapObject
	// goverter:map . Heading | ExtendMysteryTraderHeading
	// goverter:map . Destination | ExtendMysteryTraderDestination
	ConvertMysteryTrader(source generated.Mysterytrader) *cs.MysteryTrader
	ConvertMysteryTraders(source []generated.Mysterytrader) []*cs.MysteryTrader

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
	// goverter:map MineFieldOrders.Detonate Detonate
	ConvertGameMineField(source *cs.MineField) generated.Minefield

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:map MineFieldOrders.Detonate Detonate
	ConvertGameMineFieldToCreateParams(source *cs.MineField) generated.CreateMineFieldParams

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:map MineFieldOrders.Detonate Detonate
	ConvertGameMineFieldToUpdateParams(source *cs.MineField) generated.UpdateMineFieldParams

	// goverter:map . GameDBObject
	// goverter:map . MapObject | ExtendMineFieldMapObject
	// goverter:map . MineFieldOrders
	ConvertMineField(source generated.Minefield) *cs.MineField
	ConvertMineFields(source []generated.Minefield) []*cs.MineField

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:autoMap Cargo
	// goverter:map Cargo.Ironium Ironium
	// goverter:map Cargo.Boranium Boranium
	// goverter:map Cargo.Germanium Germanium
	// goverter:map Heading.X Headingx
	// goverter:map Heading.Y Headingy
	ConvertGameMineralPacket(source *cs.MineralPacket) generated.Mineralpacket

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:autoMap Cargo
	// goverter:map Cargo.Ironium Ironium
	// goverter:map Cargo.Boranium Boranium
	// goverter:map Cargo.Germanium Germanium
	// goverter:map Heading.X Headingx
	// goverter:map Heading.Y Headingy
	ConvertGameMineralPacketToCreateParams(source *cs.MineralPacket) generated.CreateMineralPacketParams

	// goverter:autoMap GameDBObject
	// goverter:autoMap MapObject.Position
	// goverter:autoMap MapObject
	// goverter:autoMap Cargo
	// goverter:map Cargo.Ironium Ironium
	// goverter:map Cargo.Boranium Boranium
	// goverter:map Cargo.Germanium Germanium
	// goverter:map Heading.X Headingx
	// goverter:map Heading.Y Headingy
	ConvertGameMineralPacketToUpdateParams(source *cs.MineralPacket) generated.UpdateMineralPacketParams

	// goverter:map . GameDBObject
	// goverter:map . MapObject | ExtendMineralPacketMapObject
	// goverter:map . Cargo
	// goverter:map . Heading | ExtendMineralPacketHeading
	ConvertMineralPacket(source generated.Mineralpacket) *cs.MineralPacket
	ConvertMineralPackets(source []generated.Mineralpacket) []*cs.MineralPacket

	// goverter:ignore Colonists
	salvageCargo(source generated.Salvage) cs.Cargo

	// goverter:ignore Colonists
	mineralPacketCargo(source generated.Mineralpacket) cs.Cargo
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

func NullStringToString(source sql.NullString) string {
	if source.Valid {
		return source.String
	}
	return ""
}

func StringToNullString(source string) sql.NullString {
	return sql.NullString{
		Valid:  true,
		String: source,
	}
}

func NullFloat64ToFloat64(source sql.NullFloat64) float64 {
	if source.Valid {
		return source.Float64
	}
	return 0
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

func Int64ToNullInt64(source int64) sql.NullInt64 {
	return sql.NullInt64{
		Valid: true,
		Int64: source,
	}
}

func IntToNullInt64(source int) sql.NullInt64 {
	return sql.NullInt64{
		Valid: true,
		Int64: int64(source),
	}
}

func Int64ToInt(source int64) int {
	return int(source)
}

func IntToInt64(source int) int64 {
	return int64(source)
}

func RulesToGameRules(source *generated.Rules) cs.Rules {
	return cs.Rules(*source)
}

func GameRulesToRules(source cs.Rules) *generated.Rules {
	return (*generated.Rules)(&source)
}

func RaceSpecToGameRaceSpec(source *generated.RaceSpec) cs.RaceSpec {
	return (cs.RaceSpec)(*source)
}

func GameRaceSpecToRaceSpec(source cs.RaceSpec) *generated.RaceSpec {
	return (*generated.RaceSpec)(&source)
}

func RaceGenSpecToGameRaceSpec(source *generated.RaceSpec) cs.RaceSpec {
	return (cs.RaceSpec)(*source)
}

func GameRaceSpecToRaceGenSpec(source cs.RaceSpec) *generated.RaceSpec {
	return (*generated.RaceSpec)(&source)
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

func PlanetIntelsToGamePlanetIntels(source *generated.PlanetIntels) []cs.PlanetIntel {
	// return an empty slice for nil
	if source == nil {
		return []cs.PlanetIntel{}
	}
	return ([]cs.PlanetIntel)(*source)
}

func GamePlanetIntelsToPlanetIntels(source []cs.PlanetIntel) *generated.PlanetIntels {
	return (*generated.PlanetIntels)(&source)
}

func FleetIntelsToGameFleetIntels(source *generated.FleetIntels) []cs.FleetIntel {
	// return an empty slice for nil
	if source == nil {
		return []cs.FleetIntel{}
	}
	return ([]cs.FleetIntel)(*source)
}

func GameFleetIntelsToFleetIntels(source []cs.FleetIntel) *generated.FleetIntels {
	return (*generated.FleetIntels)(&source)
}

func ShipDesignIntelsToGameShipDesignIntels(source *generated.ShipDesignIntels) []cs.ShipDesignIntel {
	// return an empty slice for nil
	if source == nil {
		return []cs.ShipDesignIntel{}
	}
	return ([]cs.ShipDesignIntel)(*source)
}

func GameShipDesignIntelsToShipDesignIntels(source []cs.ShipDesignIntel) *generated.ShipDesignIntels {
	return (*generated.ShipDesignIntels)(&source)
}

func MineralPacketIntelsToGameMineralPacketIntels(source *generated.MineralPacketIntels) []cs.MineralPacketIntel {
	// return an empty slice for nil
	if source == nil {
		return []cs.MineralPacketIntel{}
	}
	return ([]cs.MineralPacketIntel)(*source)
}

func GameMineralPacketIntelsToMineralPacketIntels(source []cs.MineralPacketIntel) *generated.MineralPacketIntels {
	return (*generated.MineralPacketIntels)(&source)
}

func SalvageIntelsToGameSalvageIntels(source *generated.SalvageIntels) []cs.SalvageIntel {
	// return an empty slice for nil
	if source == nil {
		return []cs.SalvageIntel{}
	}
	return ([]cs.SalvageIntel)(*source)
}

func GameSalvageIntelsToSalvageIntels(source []cs.SalvageIntel) *generated.SalvageIntels {
	return (*generated.SalvageIntels)(&source)
}

func MineFieldIntelsToGameMineFieldIntels(source *generated.MineFieldIntels) []cs.MineFieldIntel {
	// return an empty slice for nil
	if source == nil {
		return []cs.MineFieldIntel{}
	}
	return ([]cs.MineFieldIntel)(*source)
}

func GameMineFieldIntelsToMineFieldIntels(source []cs.MineFieldIntel) *generated.MineFieldIntels {
	return (*generated.MineFieldIntels)(&source)
}

func WormholeIntelsToGameWormholeIntels(source *generated.WormholeIntels) []cs.WormholeIntel {
	// return an empty slice for nil
	if source == nil {
		return []cs.WormholeIntel{}
	}
	return ([]cs.WormholeIntel)(*source)
}

func GameWormholeIntelsToWormholeIntels(source []cs.WormholeIntel) *generated.WormholeIntels {
	return (*generated.WormholeIntels)(&source)
}

func MysteryTraderIntelsToGameMysteryTraderIntels(source *generated.MysteryTraderIntels) []cs.MysteryTraderIntel {
	// return an empty slice for nil
	if source == nil {
		return []cs.MysteryTraderIntel{}
	}
	return ([]cs.MysteryTraderIntel)(*source)
}

func GameMysteryTraderIntelsToMysteryTraderIntels(source []cs.MysteryTraderIntel) *generated.MysteryTraderIntels {
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

func MineFieldSpecToGameMineFieldSpec(source *generated.MineFieldSpec) cs.MineFieldSpec {
	return (cs.MineFieldSpec)(*source)
}

func GameMineFieldSpecToMineFieldSpec(source cs.MineFieldSpec) *generated.MineFieldSpec {
	return (*generated.MineFieldSpec)(&source)
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
		Energy:        source.Researchcostenergy,
		Weapons:       source.Researchcostweapons,
		Propulsion:    source.Researchcostpropulsion,
		Construction:  source.Researchcostconstruction,
		Electronics:   source.Researchcostelectronics,
		Biotechnology: source.Researchcostbiotechnology,
	}
}

func ExtendHabLow(source generated.Race) cs.Hab {
	return cs.Hab{
		Grav: int(source.Hablowgrav.Int64),
		Temp: int(source.Hablowtemp.Int64),
		Rad:  int(source.Hablowrad.Int64),
	}
}

func ExtendHabHigh(source generated.Race) cs.Hab {
	return cs.Hab{
		Grav: int(source.Habhighgrav.Int64),
		Temp: int(source.Habhightemp.Int64),
		Rad:  int(source.Habhighrad.Int64),
	}
}
func ExtendVictoryConditions(source generated.Game) cs.VictoryConditions {
	return cs.VictoryConditions{
		Conditions:               source.Victoryconditionsconditions,
		NumCriteriaRequired:      int(source.Victoryconditionsnumcriteriarequired.Int64),
		YearsPassed:              int(source.Victoryconditionsyearspassed.Int64),
		OwnPlanets:               int(source.Victoryconditionsownplanets.Int64),
		AttainTechLevel:          int(source.Victoryconditionsattaintechlevel.Int64),
		AttainTechLevelNumFields: int(source.Victoryconditionsattaintechlevelnumfields.Int64),
		ExceedsScore:             int(source.Victoryconditionsexceedsscore.Int64),
		ExceedsSecondPlaceScore:  int(source.Victoryconditionsexceedssecondplacescore.Int64),
		ProductionCapacity:       int(source.Victoryconditionsproductioncapacity.Int64),
		OwnCapitalShips:          int(source.Victoryconditionsowncapitalships.Int64),
		HighestScoreAfterYears:   int(source.Victoryconditionshighestscoreafteryears.Int64),
	}
}

func ExtendArea(source generated.Game) cs.Vector {
	return cs.Vector{
		X: source.Areax.Float64,
		Y: source.Areay.Float64,
	}
}

func ExtendTechLevels(source generated.Player) cs.TechLevel {
	return cs.TechLevel{
		Energy:        int(source.Techlevelsenergy.Int64),
		Weapons:       int(source.Techlevelsweapons.Int64),
		Propulsion:    int(source.Techlevelspropulsion.Int64),
		Construction:  int(source.Techlevelsconstruction.Int64),
		Electronics:   int(source.Techlevelselectronics.Int64),
		Biotechnology: int(source.Techlevelsbiotechnology.Int64),
	}
}

func ExtendTechLevelsSpent(source generated.Player) cs.TechLevel {
	return cs.TechLevel{
		Energy:        int(source.Techlevelsspentenergy.Int64),
		Weapons:       int(source.Techlevelsspentweapons.Int64),
		Propulsion:    int(source.Techlevelsspentpropulsion.Int64),
		Construction:  int(source.Techlevelsspentconstruction.Int64),
		Electronics:   int(source.Techlevelsspentelectronics.Int64),
		Biotechnology: int(source.Techlevelsspentbiotechnology.Int64),
	}
}

func ExtendTechLevelsLight(source generated.GetLightPlayerForGameRow) cs.TechLevel {
	return cs.TechLevel{
		Energy:        int(source.Techlevelsenergy.Int64),
		Weapons:       int(source.Techlevelsweapons.Int64),
		Propulsion:    int(source.Techlevelspropulsion.Int64),
		Construction:  int(source.Techlevelsconstruction.Int64),
		Electronics:   int(source.Techlevelselectronics.Int64),
		Biotechnology: int(source.Techlevelsbiotechnology.Int64),
	}
}

func ExtendTechLevelsSpentLight(source generated.GetLightPlayerForGameRow) cs.TechLevel {
	return cs.TechLevel{
		Energy:        int(source.Techlevelsspentenergy.Int64),
		Weapons:       int(source.Techlevelsspentweapons.Int64),
		Propulsion:    int(source.Techlevelsspentpropulsion.Int64),
		Construction:  int(source.Techlevelsspentconstruction.Int64),
		Electronics:   int(source.Techlevelsspentelectronics.Int64),
		Biotechnology: int(source.Techlevelsspentbiotechnology.Int64),
	}
}

func ExtendPlanetGameDBObject(source generated.Planet) cs.GameDBObject {
	return cs.GameDBObject{
		ID:        source.ID,
		GameID:    source.Gameid,
		CreatedAt: source.Createdat,
		UpdatedAt: source.Updatedat,
	}
}

func ExtendPlanetMapObject(source generated.Planet) cs.MapObject {
	return cs.MapObject{
		Type: cs.MapObjectTypePlanet,
		Position: cs.Vector{
			X: source.X.Float64,
			Y: source.Y.Float64,
		},
		Name:      source.Name,
		Num:       int(source.Num.Int64),
		PlayerNum: int(source.Playernum.Int64),
		Tags:      TagsToGameTags(source.Tags),
	}
}

func ExtendBaseHab(source generated.Planet) cs.Hab {
	return cs.Hab{
		Grav: int(source.Basegrav.Int64),
		Temp: int(source.Basetemp.Int64),
		Rad:  int(source.Baserad.Int64),
	}
}

func ExtendTerraformedAmount(source generated.Planet) cs.Hab {
	return cs.Hab{
		Grav: int(source.Terraformedamountgrav.Int64),
		Temp: int(source.Terraformedamounttemp.Int64),
		Rad:  int(source.Terraformedamountrad.Int64),
	}
}

func ExtendMineralConcentration(source generated.Planet) cs.Mineral {
	return cs.Mineral{
		Ironium:   int(source.Mineralconcironium.Int64),
		Boranium:  int(source.Mineralconcboranium.Int64),
		Germanium: int(source.Mineralconcgermanium.Int64),
	}
}

func ExtendMineYears(source generated.Planet) cs.Mineral {
	return cs.Mineral{
		Ironium:   int(source.Mineyearsironium.Int64),
		Boranium:  int(source.Mineyearsboranium.Int64),
		Germanium: int(source.Mineyearsgermanium.Int64),
	}
}

func ExtendFleetMapObject(source generated.Fleet) cs.MapObject {
	return cs.MapObject{
		Type: cs.MapObjectTypeFleet,
		Position: cs.Vector{
			X: source.X.Float64,
			Y: source.Y.Float64,
		},
		Name:      source.Name,
		Num:       int(source.Num.Int64),
		PlayerNum: int(source.Playernum.Int64),
		Tags:      TagsToGameTags(source.Tags),
	}
}

func ExtendFleetFleetOrders(source generated.Fleet) cs.FleetOrders {
	return cs.FleetOrders{
		BattlePlanNum: int(source.Battleplannum),
		Waypoints:     *source.Waypoints,
		RepeatOrders:  source.Repeatorders.Bool,
		Purpose:       *source.Purpose,
	}
}

func ExtendFleetHeading(source generated.Fleet) cs.Vector {
	return cs.Vector{
		X: source.Headingx.Float64,
		Y: source.Headingy.Float64,
	}
}

func ExtendFleetPreviousPosition(source generated.Fleet) *cs.Vector {
	if !source.Previouspositionx.Valid || !source.Previouspositiony.Valid {
		return nil
	}
	return &cs.Vector{
		X: source.Previouspositionx.Float64,
		Y: source.Previouspositiony.Float64,
	}
}

func ExtendMysteryTraderMapObject(source generated.Mysterytrader) cs.MapObject {
	return cs.MapObject{
		Type: cs.MapObjectTypeMysteryTrader,
		Position: cs.Vector{
			X: source.X.Float64,
			Y: source.Y.Float64,
		},
		Name: source.Name,
		Num:  int(source.Num.Int64),
		Tags: TagsToGameTags(source.Tags),
	}
}

func ExtendMysteryTraderHeading(source generated.Mysterytrader) cs.Vector {
	return cs.Vector{
		X: source.Headingx.Float64,
		Y: source.Headingy.Float64,
	}
}

func ExtendMysteryTraderDestination(source generated.Mysterytrader) cs.Vector {
	return cs.Vector{
		X: source.Destinationx.Float64,
		Y: source.Destinationy.Float64,
	}
}

func ExtendSalvageMapObject(source generated.Salvage) cs.MapObject {
	return cs.MapObject{
		Type: cs.MapObjectTypeSalvage,
		Position: cs.Vector{
			X: source.X.Float64,
			Y: source.Y.Float64,
		},
		Name:      source.Name,
		Num:       int(source.Num.Int64),
		PlayerNum: int(source.Playernum.Int64),
		Tags:      TagsToGameTags(source.Tags),
	}
}

func ExtendMineFieldMapObject(source generated.Minefield) cs.MapObject {
	return cs.MapObject{
		Type: cs.MapObjectTypeMineField,
		Position: cs.Vector{
			X: source.X.Float64,
			Y: source.Y.Float64,
		},
		Name:      source.Name,
		Num:       int(source.Num.Int64),
		PlayerNum: int(source.Playernum.Int64),
		Tags:      TagsToGameTags(source.Tags),
	}
}

func ExtendMineralPacketHeading(source generated.Mineralpacket) cs.Vector {
	return cs.Vector{
		X: source.Headingx.Float64,
		Y: source.Headingy.Float64,
	}
}

func ExtendMineralPacketMapObject(source generated.Mineralpacket) cs.MapObject {
	return cs.MapObject{
		Type: cs.MapObjectTypeMineralPacket,
		Position: cs.Vector{
			X: source.X.Float64,
			Y: source.Y.Float64,
		},
		Name:      source.Name,
		Num:       int(source.Num.Int64),
		PlayerNum: int(source.Playernum.Int64),
		Tags:      TagsToGameTags(source.Tags),
	}
}
