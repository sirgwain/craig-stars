import { MapObjectTypeSalvage, type SalvageIntel } from './cs';

export function newSalvage(): SalvageIntel {
	return {
		type: MapObjectTypeSalvage,
		name: '',
		position: { x: 0, y: 0 },
		cargo: {},
		num: 0,
		playerNum: 0,
		tags: {},
		reportAge: 0
	};
}
