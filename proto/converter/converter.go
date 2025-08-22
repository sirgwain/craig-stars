//go:generate go run ../../generators/enum-mapper-gen ../../enums.yaml
//go:generate go tool github.com/jmattheis/goverter/cmd/goverter gen ./
package converter

import (
	"time"

	"github.com/sirgwain/craig-stars/cs"
	craig_starsv1 "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

// will be instanciated in ./converter.init.go
var C Converter

// goverter:converter
// goverter:name ProtoConverter
// goverter:output:package github.com/sirgwain/craig-stars/proto/converter
// goverter:output:file ./converter.generated.go
// goverter:ignoreUnexported
// goverter:matchIgnoreCase
// goverter:skipCopySameType
// goverter:useZeroValueOnPointerInconsistency
// goverter:extend Int32ToInt
// goverter:extend IntToInt32
// goverter:extend TimestampToTime
// goverter:extend TimeToTimestamp
// goverter:extend TimepToTimestamp
// goverter:extend TimestampToTimep
// goverter:extend Int32ToBitmask
// goverter:extend BitmaskToInt32
// goverter:extend Int32ToFloat64
// goverter:extend Float64ToInt32
// goverter:extend Uint32ToUint
// goverter:extend UintToUint32
// goverter:extend MinefieldTypeToInt32
// goverter:extend Int32ToMinefieldType
// goverter:extend CargoTransfersToCSCargoTransfers
// goverter:extend CSCargoTransfersToCargoTransfers
// goverter:extend ActionsPerRoundToCSActionsPerRound
// goverter:extend CSActionsPerRoundToActionsPerRound
//
// goverter:extend AIDifficultyToCSAIDifficulty
// goverter:extend CSAIDifficultyToAIDifficulty
// goverter:extend BattleAttackWhoToCSBattleAttackWho
// goverter:extend CSBattleAttackWhoToBattleAttackWho
// goverter:extend BattleRecordTokenActionTypeToCSBattleRecordTokenActionType
// goverter:extend CSBattleRecordTokenActionTypeToBattleRecordTokenActionType
// goverter:extend BattleTacticToCSBattleTactic
// goverter:extend CSBattleTacticToBattleTactic
// goverter:extend BattleTargetToCSBattleTarget
// goverter:extend CSBattleTargetToBattleTarget
// goverter:extend CargoTransferStatusToCSCargoTransferStatus
// goverter:extend CSCargoTransferStatusToCargoTransferStatus
// goverter:extend CometSizeToCSCometSize
// goverter:extend CSCometSizeToCometSize
// goverter:extend DensityToCSDensity
// goverter:extend CSDensityToDensity
// goverter:extend FleetPurposeToCSFleetPurpose
// goverter:extend CSFleetPurposeToFleetPurpose
// goverter:extend GameStartModeToCSGameStartMode
// goverter:extend CSGameStartModeToGameStartMode
// goverter:extend GameStateToCSGameState
// goverter:extend CSGameStateToGameState
// goverter:extend MapObjectTypeToCSMapObjectType
// goverter:extend CSMapObjectTypeToMapObjectType
// goverter:extend MinefieldTypeToCSMinefieldType
// goverter:extend CSMinefieldTypeToMinefieldType
// goverter:extend MysteryTraderRewardTypeToCSMysteryTraderRewardType
// goverter:extend CSMysteryTraderRewardTypeToMysteryTraderRewardType
// goverter:extend NewGamePlayerTypeToCSNewGamePlayerType
// goverter:extend CSNewGamePlayerTypeToNewGamePlayerType
// goverter:extend NextResearchFieldToCSNextResearchField
// goverter:extend CSNextResearchFieldToNextResearchField
// goverter:extend PlayerMessageTargetTypeToCSPlayerMessageTargetType
// goverter:extend CSPlayerMessageTargetTypeToPlayerMessageTargetType
// goverter:extend PlayerMessageTypeToCSPlayerMessageType
// goverter:extend CSPlayerMessageTypeToPlayerMessageType
// goverter:extend PlayerPositionsToCSPlayerPositions
// goverter:extend CSPlayerPositionsToPlayerPositions
// goverter:extend PlayerRelationToCSPlayerRelation
// goverter:extend CSPlayerRelationToPlayerRelation
// goverter:extend PRTToCSPRT
// goverter:extend CSPRTToPRT
// goverter:extend QueueItemTypeToCSQueueItemType
// goverter:extend CSQueueItemTypeToQueueItemType
// goverter:extend ResearchCostLevelToCSResearchCostLevel
// goverter:extend CSResearchCostLevelToResearchCostLevel
// goverter:extend ResourceTypeToCSResourceType
// goverter:extend CSResourceTypeToResourceType
// goverter:extend ShipDesignPurposeToCSShipDesignPurpose
// goverter:extend CSShipDesignPurposeToShipDesignPurpose
// goverter:extend SizeToCSSize
// goverter:extend CSSizeToSize
// goverter:extend SpendLeftoverPointsOnToCSSpendLeftoverPointsOn
// goverter:extend CSSpendLeftoverPointsOnToSpendLeftoverPointsOn
// goverter:extend TechCategoryToCSTechCategory
// goverter:extend CSTechCategoryToTechCategory
// goverter:extend TechFieldToCSTechField
// goverter:extend CSTechFieldToTechField
// goverter:extend TechHullTypeToCSTechHullType
// goverter:extend CSTechHullTypeToTechHullType
// goverter:extend TechOriginToCSTechOrigin
// goverter:extend CSTechOriginToTechOrigin
// goverter:extend TerraformHabTypeToCSTerraformHabType
// goverter:extend CSTerraformHabTypeToTerraformHabType
// goverter:extend UserRoleToCSUserRole
// goverter:extend CSUserRoleToUserRole
// goverter:extend WaypointTaskToCSWaypointTask
// goverter:extend CSWaypointTaskToWaypointTask
// goverter:extend WaypointTaskTransportActionToCSWaypointTaskTransportAction
// goverter:extend CSWaypointTaskTransportActionToWaypointTaskTransportAction
// goverter:extend WormholeStabilityToCSWormholeStability
// goverter:extend CSWormholeStabilityToWormholeStability
//
// goverter:extend MinefieldTypeMapToIntMap
// goverter:extend IntMapToMinefieldTypeMap
// goverter:extend MinefieldTypeToInt32
// goverter:extend Int32ToMinefieldType
// goverter:extend SpendLeftoverPointsOnMapToIntMap
// goverter:extend IntMapToSpendLeftoverPointsOnMap
// goverter:extend CometSizeMapToCometStatsMap
// goverter:extend CometStatsMapToCometSizeMap
// goverter:extend StringMapToPRTSpecMap
// goverter:extend PRTSpecMapToStringMap
// goverter:extend IntMapToLRTSpecMap
// goverter:extend LRTSpecMapToIntMap
// goverter:extend MinefieldTypeMapToMinefieldStatsMap
// goverter:extend MinefieldStatsMapToMinefieldTypeMap
// goverter:extend PacketDecayRateMapTocketDecayRateMap
// goverter:extend ProtoToPacketDecayRateMap
// goverter:extend RandomEventChancesMapTondomEventChancesMap
// goverter:extend ProtoToRandomEventChancesMap
// goverter:extend RepairRatesMapTopairRatesMap
// goverter:extend ProtoToRepairRatesMap
// goverter:extend SizeMapToIntMap
// goverter:extend IntMapToSizeMap
// goverter:extend WormholeStabilityMapToIntMap
// goverter:extend IntMapToWormholeStabilityMap
// goverter:extend TechCostOffsetToIntMap
// goverter:extend IntMapToTechCostOffset
// goverter:extend TerraformHabTypeMapToIntMap
// goverter:extend IntMapToTerraformHabTypeMap
// goverter:extend ShipDesignPurposeMapToBoolMap
// goverter:extend BoolMapToShipDesignPurposeMap
// goverter:extend QueueItemTypeMapToIntMap
// goverter:extend IntMapToQueueItemTypeMap
type Converter interface {
	ConvertCargo(source *craig_starsv1.Cargo) cs.Cargo
	ConvertCSCargo(source cs.Cargo) *craig_starsv1.Cargo
	ConvertCost(source *craig_starsv1.Cost) cs.Cost
	ConvertCSCost(source cs.Cost) *craig_starsv1.Cost
	ConvertHab(source *craig_starsv1.Hab) cs.Hab
	ConvertTech(source *craig_starsv1.Tech) cs.Tech
	ConvertTechLevel(source *craig_starsv1.TechLevel) cs.TechLevel
	ConvertRules(source *craig_starsv1.Rules) *cs.Rules
	ConvertCSRules(source *cs.Rules) *craig_starsv1.Rules
	ConvertVector(source *craig_starsv1.Vector) cs.Vector
	ConvertWaypointDest(source *craig_starsv1.WaypointDest) cs.WaypointDest

	// goverter:ignore Delete
	ConvertMapObject(source *craig_starsv1.MapObject) cs.MapObject

	// goverter:map . DBObject
	// goverter:map . UserSettings
	// goverter:ignore Password
	// goverter:ignore Email
	ConvertUser(source *craig_starsv1.User) cs.User
	ConvertUserSettings(source *craig_starsv1.UserSettings) cs.UserSettings

	// goverter:autoMap DBObject
	// goverter:autoMap UserSettings
	ConvertCSUser(source cs.User) *craig_starsv1.User
	ConvertCSUsers(source []cs.User) []*craig_starsv1.User

	// goverter:map . DBObject
	// goverter:ignore Spec
	ConvertRace(source craig_starsv1.Race) cs.Race
	ConvertRaceP(source *craig_starsv1.Race) *cs.Race

	// goverter:autoMap DBObject
	ConvertCSRace(source cs.Race) *craig_starsv1.Race
	ConvertCSRaces(source []cs.Race) []*craig_starsv1.Race
	ConvertCSRaceSpec(source cs.RaceSpec) *craig_starsv1.RaceSpec

	// goverter:map . DBObject
	// goverter:ignore Rules
	ConvertGame(source *craig_starsv1.Game) *cs.Game
	ConvertVictoryConditions(source *craig_starsv1.VictoryConditions) cs.VictoryConditions

	// goverter:ignore Rules
	// goverter:ignore TechStore
	ConvertGameSettings(source *craig_starsv1.GameSettings) cs.GameSettings

	// goverter:autoMap DBObject
	ConvertCSGame(source cs.Game) *craig_starsv1.Game

	ConvertCSGameWithPlayers(source *cs.GameWithPlayers) *craig_starsv1.GameWithPlayers
	ConvertCSGamesWithPlayers(source []cs.GameWithPlayers) []*craig_starsv1.GameWithPlayers

	ConvertPlayerStatus(source *craig_starsv1.PlayerStatus) cs.GamePlayer

	// goverter:map Race | ExtendProtoRace
	// goverter:ignore Intels
	// goverter:ignore Designs
	// goverter:ignore Spec
	// goverter:ignore TechsJustGained
	ConvertPlayer(source *craig_starsv1.Player) *cs.Player

	ConvertIntels(source *craig_starsv1.Intels) cs.Intels
	ConvertCSIntels(source cs.Intels) *craig_starsv1.Intels
	ConvertCSPlayerIntels(source []cs.PlayerIntel) []*craig_starsv1.PlayerIntel
	ConvertCSScoreIntels(source []cs.ScoreIntel) []*craig_starsv1.ScoreIntel
	ConvertCSBattleRecords(source []cs.BattleRecord) []*craig_starsv1.BattleRecord

	// goverter:map Race | ExtendPlayerRace
	ConvertCSPlayer(source *cs.Player) *craig_starsv1.Player
	ConvertCSCargoTransfers(source cs.CargoTransfers) map[string]*craig_starsv1.CargoTransfers

	ConvertPlayerResearchSpec(source *craig_starsv1.PlayerResearchSpec) cs.PlayerResearchSpec
	ConvertCSPlayerResearchSpec(source cs.PlayerResearchSpec) *craig_starsv1.PlayerResearchSpec

	ConvertByHandCargoTransfer(source *craig_starsv1.ByHandCargoTransfer) cs.ByHandCargoTransfer
	ConvertCSByHandCargoTransfer(source cs.ByHandCargoTransfer) *craig_starsv1.ByHandCargoTransfer

	ConvertBattlePlan(source *craig_starsv1.BattlePlan) cs.BattlePlan
	ConvertCSBattlePlan(source cs.BattlePlan) *craig_starsv1.BattlePlan

	ConvertTransportPlan(source *craig_starsv1.TransportPlan) cs.TransportPlan
	ConvertCSTransportPlan(source cs.TransportPlan) *craig_starsv1.TransportPlan

	ConvertProductionPlan(source *craig_starsv1.ProductionPlan) cs.ProductionPlan
	ConvertCSProductionPlan(source cs.ProductionPlan) *craig_starsv1.ProductionPlan

	ConvertPlayerRelations(source []*craig_starsv1.PlayerRelationship) []cs.PlayerRelationship
	ConvertPlayerOrders(source *craig_starsv1.PlayerOrders) cs.PlayerOrders

	// goverter:ignore Delete
	ConvertShipDesign(source craig_starsv1.ShipDesign) cs.ShipDesign
	ConvertShipDesignP(source *craig_starsv1.ShipDesign) *cs.ShipDesign
	ConvertShipDesigns(source []*craig_starsv1.ShipDesign) []*cs.ShipDesign

	ConvertShipDesignSpec(source *craig_starsv1.ShipDesignSpec) *cs.ShipDesignSpec
	ConvertCSShipDesignSpec(source cs.ShipDesignSpec) *craig_starsv1.ShipDesignSpec

	ConvertCSShipDesign(source *cs.ShipDesign) *craig_starsv1.ShipDesign
	ConvertCSShipDesigns(source []*cs.ShipDesign) []*craig_starsv1.ShipDesign

	ConvertFleet(source *craig_starsv1.Fleet) *cs.Fleet
	ConvertFleetOrders(source *craig_starsv1.FleetOrders) *cs.FleetOrders

	ConvertCSFleet(source *cs.Fleet) *craig_starsv1.Fleet
	ConvertCSFleets(source []*cs.Fleet) []*craig_starsv1.Fleet

	ConvertFleetSpec(source *craig_starsv1.FleetSpec) cs.FleetSpec
	ConvertCSFleetSpec(source cs.FleetSpec) *craig_starsv1.FleetSpec

	ConvertShipToken(source *craig_starsv1.ShipToken) cs.ShipToken
	ConvertShipTokens(source []*craig_starsv1.ShipToken) []cs.ShipToken

	// goverter:ignore RandomArtifact
	// goverter:ignore Starbase
	// goverter:ignore Dirty
	ConvertPlanet(source craig_starsv1.Planet) cs.Planet
	ConvertPlanetP(source *craig_starsv1.Planet) *cs.Planet
	ConvertPlanets(source []*craig_starsv1.Planet) []*cs.Planet
	ConvertPlanetOrders(source *craig_starsv1.PlanetOrders) *cs.PlanetOrders

	ConvertCSPlanet(source *cs.Planet) *craig_starsv1.Planet
	ConvertCSPlanets(source []*cs.Planet) []*craig_starsv1.Planet

	ConvertPlanetSpec(source *craig_starsv1.PlanetSpec) cs.PlanetSpec
	ConvertCSPlanetSpec(source cs.PlanetSpec) *craig_starsv1.PlanetSpec

	ConvertProductionQueueItem(source *craig_starsv1.ProductionQueueItem) cs.ProductionQueueItem
	ConvertCSProductionQueueItem(source cs.ProductionQueueItem) *craig_starsv1.ProductionQueueItem

	ConvertCSBattleRecord(source cs.BattleRecord) *craig_starsv1.BattleRecord
	ConvertBattleRecordTokenAction(source *craig_starsv1.BattleRecordTokenAction) cs.BattleRecordTokenAction
	ConvertCSBattleRecordTokenAction(source cs.BattleRecordTokenAction) *craig_starsv1.BattleRecordTokenAction

	// goverter:map . MysteryTraderReward
	ConvertPlayerMessageSpecMysteryTrader(source *craig_starsv1.PlayerMessageSpecMysteryTrader) *cs.PlayerMessageSpecMysteryTrader

	// goverter:autoMap MysteryTraderReward
	ConvertCSPlayerMessageSpecMysteryTrader(source *cs.PlayerMessageSpecMysteryTrader) *craig_starsv1.PlayerMessageSpecMysteryTrader

	// goverter:map . Engine
	ConvertTechHullComponent(source craig_starsv1.TechHullComponent) cs.TechHullComponent
	// goverter:autoMap Engine
	ConvertCSTechHullComponent(source cs.TechHullComponent) craig_starsv1.TechHullComponent

	ConvertTechHullComponentP(source *craig_starsv1.TechHullComponent) *cs.TechHullComponent
	ConvertCSTechHullComponentP(source *cs.TechHullComponent) *craig_starsv1.TechHullComponent
	ConvertCSTechHullComponents(source []cs.TechHullComponent) []*craig_starsv1.TechHullComponent

	ConvertTechHull(source *craig_starsv1.TechHull) *cs.TechHull
	ConvertCSTechHull(source *cs.TechHull) *craig_starsv1.TechHull
	ConvertCSTechHulls(source []cs.TechHull) []*craig_starsv1.TechHull

	ConvertTechPlanetary(source *craig_starsv1.TechPlanetary) *cs.TechPlanetary
	ConvertCSTechPlanetary(source *cs.TechPlanetary) *craig_starsv1.TechPlanetary
	ConvertCSTechPlanetaries(source []cs.TechPlanetary) []*craig_starsv1.TechPlanetary

	ConvertTechPlanetaryScanner(source *craig_starsv1.TechPlanetaryScanner) *cs.TechPlanetaryScanner
	ConvertCSTechPlanetaryScanner(source *cs.TechPlanetaryScanner) *craig_starsv1.TechPlanetaryScanner
	ConvertCSTechPlanetaryScanners(source []cs.TechPlanetaryScanner) []*craig_starsv1.TechPlanetaryScanner

	ConvertTechDefense(source *craig_starsv1.TechDefense) *cs.TechDefense
	ConvertCSTechDefense(source *cs.TechDefense) *craig_starsv1.TechDefense
	ConvertCSTechDefenses(source []cs.TechDefense) []*craig_starsv1.TechDefense

	ConvertTechTerraform(source *craig_starsv1.TechTerraform) *cs.TechTerraform
	ConvertCSTechTerraform(source *cs.TechTerraform) *craig_starsv1.TechTerraform
	ConvertCSTechTerraforms(source []cs.TechTerraform) []*craig_starsv1.TechTerraform

	ConvertMinefield(source *craig_starsv1.Minefield) *cs.Minefield
	ConvertMinefieldOrders(source *craig_starsv1.MinefieldOrders) *cs.MinefieldOrders

	ConvertCSMinefield(source *cs.Minefield) *craig_starsv1.Minefield
	ConvertCSMinefields(source []*cs.Minefield) []*craig_starsv1.Minefield
	ConvertCSMinefieldSpec(source cs.MinefieldSpec) *craig_starsv1.MinefieldSpec

	ConvertMineralPacket(source *craig_starsv1.MineralPacket) *cs.MineralPacket
	ConvertCSMineralPacket(source *cs.MineralPacket) *craig_starsv1.MineralPacket
	ConvertCSMineralPackets(source []*cs.MineralPacket) []*craig_starsv1.MineralPacket

	ConvertCSMysteryTraders(source []*cs.MysteryTrader) []*craig_starsv1.MysteryTrader

	ConvertWormhole(source *craig_starsv1.Wormhole) *cs.Wormhole
	ConvertCSWormhole(source *cs.Wormhole) *craig_starsv1.Wormhole
	ConvertCSWormholes(source []*cs.Wormhole) []*craig_starsv1.Wormhole

	ConvertSalvage(source *craig_starsv1.Salvage) *cs.Salvage
	ConvertCSSalvage(source *cs.Salvage) *craig_starsv1.Salvage
	ConvertCSSalvages(source []*cs.Salvage) []*craig_starsv1.Salvage

	ConvertMinefieldStats(source *craig_starsv1.MinefieldStats) cs.MinefieldStats
	ConvertCSMinefieldStats(source cs.MinefieldStats) *craig_starsv1.MinefieldStats

	ConvertWormholeStats(source *craig_starsv1.WormholeStats) cs.WormholeStats
	ConvertCSWormholeStats(source cs.WormholeStats) *craig_starsv1.WormholeStats

	ConvertCometStats(source *craig_starsv1.CometStats) cs.CometStats
	ConvertCSCometStats(source cs.CometStats) *craig_starsv1.CometStats

	ConvertPRTSpec(source *craig_starsv1.PRTSpec) cs.PRTSpec
	ConvertCSPRTSpec(source cs.PRTSpec) *craig_starsv1.PRTSpec

	ConvertLRTSpec(source *craig_starsv1.LRTSpec) cs.LRTSpec
	ConvertCSLRTSpec(source cs.LRTSpec) *craig_starsv1.LRTSpec
}

func TimestampToTime(source *timestamppb.Timestamp) time.Time {
	if source == nil {
		return time.Time{}
	}
	return source.AsTime()
}

func TimepToTimestamp(source *time.Time) *timestamppb.Timestamp {
	if source == nil {
		return nil
	}
	return timestamppb.New(*source)
}

func TimestampToTimep(source *timestamppb.Timestamp) *time.Time {
	if source == nil {
		return nil
	}
	t := source.AsTime()
	return &t
}

func TimeToTimestamp(source time.Time) *timestamppb.Timestamp {
	return timestamppb.New(source)
}
func Int32ToInt(source int32) int {
	return int(source)
}

func IntToInt32(source int) int32 {
	return int32(source)
}

func Uint32ToUint(source uint32) uint {
	return uint(source)
}

func UintToUint32(source uint) uint32 {
	return uint32(source)
}

func Int32ToBitmask(source int32) cs.Bitmask {
	return cs.Bitmask(source)
}

func BitmaskToInt32(source cs.Bitmask) int32 {
	return int32(source)
}

func Int32ToFloat64(source int32) float64 {
	return float64(source)
}

func Float64ToInt32(source float64) int32 {
	return int32(source)
}

func CargoTransfersToCSCargoTransfers(c Converter, source *craig_starsv1.CargoTransfers) []cs.ByHandCargoTransfer {
	var dest []cs.ByHandCargoTransfer
	if source == nil {
		return dest
	}

	for _, v := range source.Transfers {
		dest = append(dest, c.ConvertByHandCargoTransfer(v))
	}
	return dest
}

func CSCargoTransfersToCargoTransfers(c Converter, source []cs.ByHandCargoTransfer) *craig_starsv1.CargoTransfers {
	if len(source) == 0 {
		return nil
	}
	var dest craig_starsv1.CargoTransfers
	for _, v := range source {
		dest.Transfers = append(dest.Transfers, c.ConvertCSByHandCargoTransfer(v))
	}
	return &dest
}

func ActionsPerRoundToCSActionsPerRound(c Converter, source *craig_starsv1.ActionsPerRound) []cs.BattleRecordTokenAction {
	var dest []cs.BattleRecordTokenAction
	if source == nil {
		return dest
	}

	for _, v := range source.Actions {
		dest = append(dest, c.ConvertBattleRecordTokenAction(v))
	}
	return dest
}

func CSActionsPerRoundToActionsPerRound(c Converter, source []cs.BattleRecordTokenAction) *craig_starsv1.ActionsPerRound {
	if len(source) == 0 {
		return nil
	}
	var dest craig_starsv1.ActionsPerRound
	for _, v := range source {
		dest.Actions = append(dest.Actions, c.ConvertCSBattleRecordTokenAction(v))
	}
	return &dest
}

func ExtendPlayerRace(c Converter, source cs.Race) *craig_starsv1.Race {
	return c.ConvertCSRace(source)
}

func ExtendProtoRace(c Converter, source *craig_starsv1.Race) cs.Race {
	return *c.ConvertRaceP(source)
}

func ResourceTypeToCSResourceType(m craig_starsv1.ResourceType) cs.ResourceType {
	switch m {
	case craig_starsv1.ResourceType_RESOURCE_TYPE_IRONIUM:
		return cs.Ironium
	case craig_starsv1.ResourceType_RESOURCE_TYPE_BORANIUM:
		return cs.Boranium
	case craig_starsv1.ResourceType_RESOURCE_TYPE_GERMANIUM:
		return cs.Germanium
	case craig_starsv1.ResourceType_RESOURCE_TYPE_COLONISTS:
		return cs.Colonists
	case craig_starsv1.ResourceType_RESOURCE_TYPE_FUEL:
		return cs.Fuel
	case craig_starsv1.ResourceType_RESOURCE_TYPE_RESOURCES:
		return cs.Resources
	default:
		return cs.ResourceType(0) // TODO: this is ironium....
	}
}

func CSResourceTypeToResourceType(m cs.ResourceType) craig_starsv1.ResourceType {
	switch m {
	case cs.Ironium:
		return craig_starsv1.ResourceType_RESOURCE_TYPE_IRONIUM
	case cs.Boranium:
		return craig_starsv1.ResourceType_RESOURCE_TYPE_BORANIUM
	case cs.Germanium:
		return craig_starsv1.ResourceType_RESOURCE_TYPE_GERMANIUM
	case cs.Colonists:
		return craig_starsv1.ResourceType_RESOURCE_TYPE_COLONISTS
	case cs.Fuel:
		return craig_starsv1.ResourceType_RESOURCE_TYPE_FUEL
	case cs.Resources:
		return craig_starsv1.ResourceType_RESOURCE_TYPE_RESOURCES
	default:
		return craig_starsv1.ResourceType_RESOURCE_TYPE_UNSPECIFIED
	}
}
