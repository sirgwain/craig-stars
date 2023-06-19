import type { Cargo } from './Cargo';
import type { MapObject } from './MapObject';
import type { Vector } from './Vector';

export type MineralPacket = {
	targetPlanetNum?: number;
	cargo: Cargo;
	warpFactor: number;
	safeWarpSpeed?: number;
	heading: Vector;
} & MapObject;
