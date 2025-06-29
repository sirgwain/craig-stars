import {
	type AnyFleet,
	type AnyMineField,
	type AnyMineralPacket,
	type AnyPlanet
} from '$lib/services/Universe';
import type {
	MapObject,
	MapObjectTarget,
	MysteryTraderIntel,
	SalvageIntel,
	Vector,
	WormholeIntel
} from './cs';
import {
	MapObjectTypeFleet,
	MapObjectTypeMineField,
	MapObjectTypeMineralPacket,
	MapObjectTypeMysteryTrader,
	MapObjectTypePlanet,
	MapObjectTypeSalvage,
	MapObjectTypeWormhole,
	None
} from './cs';
import { getTokenCount, hasDestination } from './Fleet';
import { emptyVector } from './Vector';

export type MovingMapObject = {
	heading: Vector;
	warpSpeed: number;
} & MapObject;

export const emptyMapObject = (): MapObject => {
	return { type: '', name: '', position: emptyVector, num: 0, playerNum: 0, tags: {} };
};

/**
 * Get default name for a mapObject or fleet
 * @param mo The MapObject or fleet to check
 * @returns String containing name of object/fleet
 */
export function getMapObjectName(mo: MapObject | AnyFleet | undefined): string {
	if (!mo) {
		return '';
	}

	// for fleets, we want the name to indicate if it has ships
	if ('tokens' in mo) {
		const numShips = getTokenCount(mo);
		const numTokens = mo.tokens?.length;
		let name = mo.name;
		if ((numTokens ?? 0) > 1) {
			const fleetNumIndex = name.lastIndexOf(' #');
			name = name.substring(0, fleetNumIndex).concat('+', name.slice(fleetNumIndex));
		}
		return `${name}${numShips > 1 ? ` (${numShips})` : ''}${hasDestination(mo) ? '*' : ''}`;
	}
	return mo.name;
}

// get the underlying map object as a destructurable item
export function getUnderlyingMapObject(mo: MapObject | undefined) {
	return {
		planet: mo?.type === MapObjectTypePlanet ? (mo as AnyPlanet) : undefined,
		fleet: mo?.type === MapObjectTypeFleet ? (mo as AnyFleet) : undefined,
		wormhole: mo?.type === MapObjectTypeWormhole ? (mo as WormholeIntel) : undefined,
		mineField: mo?.type === MapObjectTypeMineField ? (mo as AnyMineField) : undefined,
		mysteryTrader: mo?.type === MapObjectTypeMysteryTrader ? (mo as MysteryTraderIntel) : undefined,
		salvage: mo?.type === MapObjectTypeSalvage ? (mo as SalvageIntel) : undefined,
		mineralPacket: mo?.type === MapObjectTypeMineralPacket ? (mo as AnyMineralPacket) : undefined
	};
}

/**
 * Check if this MapObject is owned by a player
 * @param mo The MapObject to check
 * @param playerNum The player
 * @returns true if this mapobject is owned by the player
 */
export function ownedBy(mo: MapObject, playerNum: number): boolean {
	return mo.playerNum === playerNum;
}

/**
 * Check if this MapObject is owned by any player
 * @param mo The MapObject to check
 * @returns true if this mapobject is owned
 */
export function owned(mo: MapObject): boolean {
	return mo.playerNum != None;
}

export function commandable(playerNum: number, mo: MapObject | undefined): boolean {
	if (!mo) return false;
	return (
		(mo.type === MapObjectTypeFleet || mo.type === MapObjectTypePlanet) &&
		mo.playerNum === playerNum
	);
}

export const positionKey = (pos: MapObject | Vector | undefined): string => {
	if (!pos) {
		return '';
	}
	const mo = pos as MapObject;
	const v = pos as Vector;
	if (mo && 'position' in mo) {
		return `${mo.position.x},${mo.position.y}`;
	} else if (v && 'x' in v) {
		return `${v.x},${v.y}`;
	}
	return '';
};

export const key = (mo: MapObject | undefined): string => {
	return `${mo?.type ?? ''}-${mo?.num ?? ''}-${mo?.playerNum ?? ''}`;
};

// compare two map objects for equivalence using their natural keys (num, type, playerNum)
export function equal(mo1: MapObject | undefined, mo2: MapObject | undefined): boolean {
	return !!(
		mo1 &&
		mo2 &&
		mo1.num === mo2.num &&
		mo1.type === mo2.type &&
		mo1.playerNum === mo2.playerNum
	);
}

export function equalsTarget(
	mo1: MapObject | undefined,
	target: MapObjectTarget | undefined
): boolean {
	return !!(
		mo1 &&
		target &&
		mo1.num === (target.targetNum ?? 0) &&
		mo1.type === (target.targetType ?? '') &&
		mo1.playerNum === (target.targetPlayerNum ?? 0)
	);
}

export function toTarget(mo: MapObject): MapObjectTarget {
	return {
		targetType: mo.type,
		targetPosition: mo.position,
		targetNum: mo.num,
		targetPlayerNum: mo.playerNum,
		targetName: mo.name
	};
}
// compare two map objects for equivalence using their natural keys (num, type, playerNum)
export function targetsEqual(
	mo1: MapObjectTarget | undefined,
	mo2: MapObjectTarget | undefined
): boolean {
	return !!(
		mo1 &&
		mo2 &&
		(mo1?.targetType ?? '') === (mo2?.targetType ?? '') &&
		mo1?.targetNum === mo2?.targetNum &&
		mo1?.targetPlayerNum === mo2?.targetPlayerNum
	);
}
