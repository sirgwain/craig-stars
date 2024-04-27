import type { BattleRecordStats } from './Battle';
import type { Cost } from './Cost';
import type { Target } from './Fleet';
import type { Hab } from './Hab';
import { MapObjectType } from './MapObject';
import type { Mineral } from './Mineral';
import type { PlayerSettings } from './PlayerSettings';
import type { QueueItemType } from './QueueItemType';
import type { TechField } from './TechLevel';

export type Message = {
	type: MessageType;
	text: string;
	battleNum?: number;
	spec: PlayerMessageSpec;
} & Target;

export type PlayerMessageSpec = {
	amount?: number;
	name?: string;
	sourcePlayerNum?: number;
	destPlayerNum?: number;
	prevAmount?: number;
	cost?: Cost;
	field?: TechField;
	nextField?: TechField;
	techGained?: string;
	queueItemType?: QueueItemType;
	battle: BattleRecordStats;
	comet?: PlayerMessageSpecComet;
};

export type PlayerMessageSpecComet = {
	size: CometSize;
	mineralsAdded: Mineral;
	mineralConcentrationIncreased: Mineral;
	habChanged: Hab;
	colonistsKilled: number;
};

export enum CometSize {
	Small = 'Small',
	Medium = 'Medium',
	Large = 'Large',
	Huge = 'Huge'
}

export enum MessageTargetType {
	None = '',
	Planet = 'Planet',
	Fleet = 'Fleet',
	Wormhole = 'Wormhole',
	MineField = 'MineField',
	MysteryTrader = 'MysteryTrader',
	MineralPacket = 'MineralPacket',
	Battle = 'Battle'
}

export function getMapObjectTypeForMessageType(
	targetType: MessageTargetType | MapObjectType | undefined
): MapObjectType {
	switch (targetType) {
		case MessageTargetType.Planet:
			return MapObjectType.Planet;
		case MessageTargetType.Fleet:
			return MapObjectType.Fleet;
		case MessageTargetType.Wormhole:
			return MapObjectType.Wormhole;
		case MessageTargetType.MineField:
			return MapObjectType.MineField;
		case MessageTargetType.MysteryTrader:
			return MapObjectType.MysteryTrader;
		case MessageTargetType.MineralPacket:
			return MapObjectType.MineralPacket;
	}

	return MapObjectType.None;
}

export enum MessageType {
	None,
	Info,
	Error,
	HomePlanet,
	PlayerDiscovery,
	PlanetDiscovery,
	PlanetProductionQueueEmpty,
	PlanetProductionQueueComplete,
	BuiltMineralAlchemy,
	BuiltMine,
	BuiltFactory,
	BuiltDefense,
	BuiltShip,
	BuiltStarbase,
	BuiltScanner,
	BuiltMineralPacket,
	BuiltTerraform,
	FleetOrdersComplete,
	FleetEngineFailure,
	FleetOutOfFuel,
	FleetGeneratedFuel,
	FleetScrapped,
	FleetMerged,
	FleetInvalidMergeNotFleet,
	FleetInvalidMergeUnowned,
	FleetPatrolTargeted,
	FleetInvalidRouteNotFriendlyPlanet,
	FleetInvalidRouteNotPlanet,
	FleetInvalidRouteNoRouteTarget,
	FleetInvalidTransport,
	FleetRoute,
	Invalid,
	PlanetColonized,
	GainTechLevel,
	MyPlanetBombed,
	MyPlanetRetroBombed,
	EnemyPlanetBombed,
	EnemyPlanetRetroBombed,
	MyPlanetInvaded,
	EnemyPlanetInvaded,
	Battle,
	CargoTransferred,
	MinesSwept,
	MinesLaid,
	MineFieldHit,
	FleetDumpedCargo,
	FleetStargateDamaged,
	MineralPacketCaught,
	MineralPacketDamage,
	MineralPacketLanded,
	MineralPacketDiscovered,
	MineralPacketTargettingPlayerDiscovered,
	Victor,
	FleetReproduce,
	RandomMineralDeposit,
	Permaform,
	Instaform,
	PacketTerraform,
	PacketPermaform,
	RemoteMined,
	TechGained,
	FleetTargetLost,
	FleetRadiatingEngineDieoff,
	PlanetDiedOff,
	PlanetEmptied,
	PlanetDiscoveryHabitable,
	PlanetDiscoveryTerraformable,
	PlanetDiscoveryUninhabitable,
	BuildInvalidItem,
	BuildMineralPacketNoMassDriver,
	BuildMineralPacketNoTarget,
	PlanetPopulationDecreased,
	PlanetPopulationDecreasedOvercrowding,
	PlayerDead,
	PlayerNoPlanets,
	CometStrike,
	CometStrikeMyPlanet,
	FleetShipExceededSafeSpeed,
	BonusResearchArtifact,
	FleetTransferGiven,
	FleetTransferGivenFailed,
	FleetTransferGivenFailedColonists,
	FleetTransferGivenRefused,
	FleetTransferReceived,
	FleetTransferReceivedFailed,
	FleetTransferReceivedRefused,
	TechLevelGainedInvasion,
	TechLevelGainedScrapFleet,
	TechLevelGainedBattle,
	FleetDieoff,
	MysteryTraderDiscovered,
	MysteryTraderChangedCourse,
	MysteryTraderAgain,
	MysteryTraderMetWithReward,
	MysteryTraderMetWithoutReward
}

// get the next visible message taking into account filters
export function getNextVisibleMessageNum(
	num: number,
	showFilteredMessages: boolean,
	messages: Message[],
	settings: PlayerSettings
): number {
	for (let i = num + 1; i < messages.length; i++) {
		if (showFilteredMessages || settings.isMessageVisible(messages[i].type)) {
			return i;
		}
	}
	return num;
}
