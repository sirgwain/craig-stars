import type {
	Fleet,
	Minefield,
	MineralPacket,
	MysteryTrader,
	Planet,
	Salvage,
	Wormhole
} from '$lib/types/cs-proto';
import {
	type MapObject,
	MapObjectSchema,
	type MapObjectTarget,
	MapObjectTargetSchema,
	MapObjectType,
	type Vector,
	type VectorJson
} from '$lib/types/cs-proto';
import { create } from '@bufbuild/protobuf';
import { None } from './Consts';
import { getTokenCount, hasDestination } from './Fleet';
import { distance } from './Vector';

export type MapObjectLike = {
	mapObject?: MapObject;
};

export type MapObjectTargetLike = {
	targetPosition?: Vector;
	targetType: MapObjectType;
	targetName?: string;
	targetNum: number;
	targetPlayerNum?: number;
};

export type MovingMapObject = {
	heading: Vector;
	warpSpeed: number;
} & MapObjectLike;

export type Position =
	| {
			x?: number;
			y?: number;
	  }
	| VectorJson;

export const emptyMapObject = (): MapObject => {
	return create(MapObjectSchema);
};

/**
 * Get default name for a mapObject or fleet
 * @param mo The MapObject or fleet to check
 * @returns String containing name of object/fleet
 */
export function getMapObjectName(mo: MapObjectLike | Fleet | undefined): string {
	if (!mo) {
		return '';
	}

	// for fleets, we want the name to indicate if it has ships
	if ('tokens' in (mo as Fleet)) {
		const fleet = mo as Fleet;
		const numShips = getTokenCount(fleet);
		const numTokens = fleet.tokens?.length;
		let name = fleet.mapObject?.name ?? '';
		if ((numTokens ?? 0) > 1 && name) {
			const fleetNumIndex = name.lastIndexOf(' #');
			if (fleetNumIndex > -1) {
				name = name.substring(0, fleetNumIndex).concat('+', name.slice(fleetNumIndex));
			}
		}
		return `${name}${numShips > 1 ? ` (${numShips})` : ''}${hasDestination(fleet) ? '*' : ''}`;
	}
	return mo.mapObject?.name ?? '';
}

// get the underlying map object as a destructurable item
export function getUnderlyingMapObject(mo: MapObjectLike | undefined) {
	const t = mo?.mapObject?.type ?? '';
	return {
		planet: t === MapObjectType.PLANET ? (mo as Planet) : undefined,
		fleet: t === MapObjectType.FLEET ? (mo as Fleet) : undefined,
		wormhole: t === MapObjectType.WORMHOLE ? (mo as Wormhole) : undefined,
		minefield: t === MapObjectType.MINEFIELD ? (mo as Minefield) : undefined,
		mysteryTrader: t === MapObjectType.MYSTERY_TRADER ? (mo as MysteryTrader) : undefined,
		salvage: t === MapObjectType.SALVAGE ? (mo as Salvage) : undefined,
		mineralPacket: t === MapObjectType.MINERAL_PACKET ? (mo as MineralPacket) : undefined
	};
}

/**
 * Check if this MapObject is owned by a player
 * @param mo The MapObject to check
 * @param playerNum The player
 * @returns true if this mapobject is owned by the player
 */
export function ownedBy(mo: MapObjectLike, playerNum: number): boolean {
	return (mo.mapObject?.playerNum ?? None) === playerNum;
}

/**
 * Check if this MapObject is owned by any player
 * @param mo The MapObject to check
 * @returns true if this mapobject is owned
 */
export function owned(mo: MapObjectLike): boolean {
	return (mo.mapObject?.playerNum ?? None) != None;
}

export function commandable(playerNum: number, mo: MapObjectLike | undefined): boolean {
	return !!(
		mo?.mapObject?.type &&
		(mo.mapObject.type === MapObjectType.PLANET || mo.mapObject.type === MapObjectType.FLEET) &&
		mo.mapObject.playerNum === playerNum
	);
}

export const positionKey = (pos: Position | MapObjectLike | undefined): string => {
	if (!pos) {
		return '';
	}
	// MapObjectLike
	if ((pos as MapObjectLike)?.mapObject?.position) {
		const p = (pos as MapObjectLike).mapObject!.position!;
		return `${p.x ?? 0},${p.y ?? 0}`;
	}
	// Position
	if ((pos as Position)?.x !== undefined) {
		const v = pos as Position;
		return `${v.x ?? 0},${v.y ?? 0}`;
	}

	return '';
};

export const key = (mo: MapObjectLike | undefined): string => {
	return `${mo?.mapObject?.type ?? ''}-${mo?.mapObject?.num ?? ''}-${mo?.mapObject?.playerNum ?? ''}`;
};

// compare two map objects for equivalence using their natural keys (num, type, playerNum)
export function equal(mo1: MapObjectLike | undefined, mo2: MapObjectLike | undefined): boolean {
	return !!(
		mo1?.mapObject &&
		mo2?.mapObject &&
		mo1.mapObject.num === mo2.mapObject.num &&
		mo1.mapObject.type === mo2.mapObject.type &&
		mo1.mapObject.playerNum === mo2.mapObject.playerNum
	);
}

export function equalsTarget(
	mo1: MapObjectLike | undefined,
	target: MapObjectTargetLike | undefined
): boolean {
	return !!(
		mo1?.mapObject &&
		target &&
		(mo1.mapObject.num ?? 0) === (target.targetNum ?? 0) &&
		(mo1.mapObject.type ?? '') === (target.targetType ?? '') &&
		(mo1.mapObject.playerNum ?? 0) === (target.targetPlayerNum ?? 0)
	);
}

export function toTarget(mo: MapObjectLike): MapObjectTarget {
	return create(MapObjectTargetSchema, {
		targetType: mo.mapObject?.type,
		targetPosition: mo.mapObject?.position,
		targetNum: mo.mapObject?.num,
		targetPlayerNum: mo.mapObject?.playerNum,
		targetName: mo.mapObject?.name
	});
}
// compare two map objects for equivalence using their natural keys (num, type, playerNum)
export function targetsEqual(
	mo1: MapObjectTargetLike | undefined,
	mo2: MapObjectTargetLike | undefined
): boolean {
	return !!(
		mo1 &&
		mo2 &&
		(mo1?.targetType ?? '') === (mo2?.targetType ?? '') &&
		mo1?.targetNum === mo2?.targetNum &&
		mo1?.targetPlayerNum === mo2?.targetPlayerNum
	);
}

export function nearest(mo: MapObjectLike, mapObjects: MapObjectLike[]): MapObjectLike | undefined {
	let nearest: MapObjectLike | undefined;
	let nearestDist = Number.MAX_VALUE;

	mapObjects.forEach((other) => {
		if (equal(mo, other)) {
			return;
		}
		const dist = distance(other.mapObject?.position, mo.mapObject?.position);
		if (dist < nearestDist) {
			nearest = other;
			nearestDist = dist;
		}
	});

	return nearest;
}
