import type { Universe } from '$lib/services/Universe';
import type { BattleRecordStats } from './Battle';
import type { Target } from './Fleet';
import { MapObjectType } from './MapObject';
import type { Player } from './Player';
import type { PlayerSettings } from './PlayerSettings';
import type { TechField } from './TechLevel';

export type Message = {
	type: MessageType;
	text: string;
	battleNum?: number;
	spec: PlayerMessageSpec;
} & Target;

export type PlayerMessageSpec = {
	amount: number;
	field: TechField;
	nextField: TechField;
	techGained: string;
	battle: BattleRecordStats;
};

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
	FleetTargetLost
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

function getBattleMessage(message: Message, universe: Universe, player: Player): string {
	const stats = message.spec.battle;
	const battle = universe.getBattle(message.battleNum);
	if (battle) {
		const location = universe.getBattleLocation(battle) ?? 'unknown';
		let text = `A battle took place at ${location}.`;

		const myShips = stats.numShipsByPlayer ? stats.numShipsByPlayer[player.num] ?? 0 : 0;

		const myShipsLost = stats.shipsDestroyedByPlayer
			? stats.shipsDestroyedByPlayer[player.num] ?? 0
			: 0;

		const myRemainingShips = myShips - myShipsLost;

		if (myShipsLost === myShips) {
			text += ' All of your forces were destroyed by enemy forces.';
		} else if (myShipsLost === 0) {
			text += ' None of your forces were destroyed.';
		} else if (myRemainingShips === 1) {
			text += ` Only one of your ships survived.`;
		} else if (myRemainingShips > 1) {
			text += ` ${myRemainingShips} of your ships surived.`;
		}

		return text;
	} else {
		return `A battle took place at an unknown location`;
	}
}

export function getMessageText(message: Message, universe: Universe, player: Player): string {
	switch (message.type) {
		case MessageType.Battle:
			return getBattleMessage(message, universe, player);
		default:
			return message.text;
	}
}
