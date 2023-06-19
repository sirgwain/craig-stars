import type { Cargo } from './Cargo';
import type { MapObject } from './MapObject';

export type Salvage = {
	cargo: Cargo;
} & MapObject;
