import { MapObjectTypeSalvage, type Salvage } from './cs';

export function newSalvage(): Salvage {
	return {
		id: 0,
		createdAt: '',
		updatedAt: '',
		gameId: 0,
		type: MapObjectTypeSalvage,
		name: '',
		position: { x: 0, y: 0 },
		cargo: {},
		num: 0,
		playerNum: 0,
		tags: {}
	};
}
