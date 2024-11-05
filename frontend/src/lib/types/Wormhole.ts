import type { MapObject } from './MapObject';
import type { WormholeStability } from './Rules';

export type Wormhole = {
	destinationNum?: number;
	stability: WormholeStability;
} & MapObject;
