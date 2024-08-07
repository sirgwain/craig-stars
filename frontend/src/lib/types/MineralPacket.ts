import type { Cargo } from './Cargo';
import type { MovingMapObject } from './MapObject';

export type MineralPacket = {
	targetPlanetNum?: number;
	cargo: Cargo;
	safeWarpSpeed?: number;
	scanRange: number;
	scanRangePen: number;
} & MovingMapObject;
