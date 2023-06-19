import type { Cargo } from './Cargo';
import type { MapObject } from './MapObject';
import type { Vector } from './Vector';

export type MineralPacket = {
	targetPlanetNum?: number;
	cargo: Cargo;
	warpSpeed: number;
	safeWarpSpeed?: number;
	heading: Vector;
	scanRange: number;
	scanRangePen: number;
} & MapObject;
