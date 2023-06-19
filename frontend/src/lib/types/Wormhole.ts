import type { WormholeStability } from './Game';
import type { MapObject } from './MapObject';

export type Wormhole = {
	destinationNum?: number;
	stability: WormholeStability;
} & MapObject;
