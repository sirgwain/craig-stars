import type { Cargo } from './Cargo';
import { MapObjectType, type MapObject } from './MapObject';

export type Salvage = {
	cargo: Cargo;
} & MapObject;

export function newSalvage(): Salvage {
	return {
		type: MapObjectType.Salvage,
		playerNum: 0,
		num: 0,
		name: '',
		position: { x: 0, y: 0 },
		cargo: {}
	};
}
