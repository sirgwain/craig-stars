import type { Vector } from './Vector';

export const None = 0;
export const Infinite = -1;
export const StargateWarpSpeed = 11;

export interface MapObject {
	id?: number;
	createdAt?: string;
	updatedAt?: string;
	type: MapObjectType;
	gameId?: number;
	position: Vector;
	name: string;
	num: number;
	playerNum: number;
}

export type MovingMapObject = {
	heading: Vector;
	warpSpeed?: number;
} & MapObject;

export enum MapObjectType {
	None = '',
	Planet = 'Planet',
	Fleet = 'Fleet',
	Wormhole = 'Wormhole',
	MineField = 'MineField',
	MysteryTrader = 'MysteryTrader',
	Salvage = 'Salvage',
	MineralPacket = 'MineralPacket',
	PositionWaypoint = 'PositionWaypoint'
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

export const positionKey = (pos: MapObject | Vector): string => {
	const mo = pos as MapObject;
	const v = pos as Vector;
	if (mo) {
		return `${mo.position.x},${mo.position.y}`;
	} else {
		return `${v.x},${v.y}`;
	}
	return '';
};

export const key = (mo: MapObject | undefined): string => {
	return `${mo?.type ?? ''}-${mo?.num ?? ''}-${mo?.playerNum ?? ''}`;
};

// compare two map objects for equivalence using their natural keys (num, type, playerNum)
export const equal = (mo1: MapObject | undefined, mo2: MapObject | undefined): boolean =>
	!!(
		mo1 &&
		mo2 &&
		mo1?.num === mo2?.num &&
		mo1?.type === mo2?.type &&
		mo1?.playerNum === mo2?.playerNum
	);
