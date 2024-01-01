import { goto } from '$app/navigation';
import { commandMapObject, selectMapObject, zoomToMapObject } from '$lib/services/Stores';
import type { Universe } from '$lib/services/Universe';
import { kebabCase } from 'lodash-es';
import type { BattleRecordStats } from './Battle';
import type { Cost } from './Cost';
import type { Fleet, Target } from './Fleet';
import type { Hab } from './Hab';
import { MapObjectType, None, ownedBy } from './MapObject';
import type { Mineral } from './Mineral';
import type { PlayerSettings } from './PlayerSettings';
import type { TechField } from './TechLevel';
import type { QueueItemType } from './QueueItemType';

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
	FleetColonistDieoff,
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

// goto a message target
export function gotoTarget(
	message: Message,
	gameId: number,
	playerNum: number,
	universe: Universe
) {
	const targetType = message.targetType ?? MessageTargetType.None;
	let moType = MapObjectType.None;

	if (message.battleNum) {
		goto(`/games/${gameId}/battles/${message.battleNum}`);
		return;
	}

	if (message.type === MessageType.GainTechLevel) {
		goto(`/games/${gameId}/research`);
	}

	if (message.type === MessageType.TechGained && message.spec.techGained) {
		goto(`/games/${gameId}/techs/${kebabCase(message.spec.techGained)}`);
	}

	if (message.targetNum) {
		switch (targetType) {
			case MessageTargetType.Planet:
				moType = MapObjectType.Planet;
				break;
			case MessageTargetType.Fleet:
				moType = MapObjectType.Fleet;
				break;
			case MessageTargetType.Wormhole:
				moType = MapObjectType.Wormhole;
				break;
			case MessageTargetType.MineField:
				moType = MapObjectType.MineField;
				break;
			case MessageTargetType.MysteryTrader:
				moType = MapObjectType.MysteryTrader;
				break;
			case MessageTargetType.MineralPacket:
				moType = MapObjectType.MineralPacket;
				break;
			case MessageTargetType.Battle:
				break;
		}

		if (moType != MapObjectType.None) {
			const target = universe.getMapObject(message);
			if (target) {
				// if this is a fleet that we own, select the planet before we command the fleet
				if (target.type == MapObjectType.Fleet && target.playerNum == playerNum) {
					const orbitingPlanetNum = (target as Fleet).orbitingPlanetNum;
					if (orbitingPlanetNum && orbitingPlanetNum != None) {
						const orbiting = universe.getPlanet(orbitingPlanetNum);
						if (orbiting) {
							selectMapObject(orbiting);
						}
					}
				} else {
					selectMapObject(target);
				}
				if (ownedBy(target, playerNum)) {
					commandMapObject(target);
				}

				// zoom on goto
				zoomToMapObject(target);
				goto(`/games/${gameId}`);
			}
		}
	}
}
