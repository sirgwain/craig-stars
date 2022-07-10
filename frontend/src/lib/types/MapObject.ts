import type { Vector } from './Vector';

export interface MapObject {
	id?: number;
	createdAt?: string;
	updatedat?: string;

	gameId?: number;
	position: Vector;
	name: string;
	num: number;
	playerNum: number | null;
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
	return mo.playerNum != null;
}
