import type { DesignFinder } from '$lib/services/Universe';
import { describe, it } from 'vitest';
import { WaypointTask, type Fleet } from './Fleet';
import { MapObjectType } from './MapObject';
import type { ShipDesign } from './ShipDesign';

// test designfinder that just returns a Long Range Scout
export class TestDesignFinder implements DesignFinder {
	designs = [longRangeScoutDesign, santaMariaDesign, cottonPickerDesign];

	getDesign(playerNum: number, num: number): ShipDesign | undefined {
		return this.designs.find((d) => d.playerNum === playerNum && d.num === num);
	}
	getMyDesign(num: number | undefined): ShipDesign | undefined {
		return this.designs.find((d) => d.num === num);
	}
}

export const longRangeScoutDesign: ShipDesign = {
	id: 0,
	gameId: 0,
	num: 1,
	playerNum: 1,
	originalPlayerNum: 0,
	name: 'Long Range Scout',
	version: 0,
	hull: 'Scout',
	hullSetNumber: 0,
	slots: [
		{
			hullComponent: 'Long Hump 6',
			hullSlotIndex: 1,
			quantity: 1
		},
		{
			hullComponent: 'Rhino Scanner',
			hullSlotIndex: 2,
			quantity: 1
		},
		{
			hullComponent: 'Fuel Tank',
			hullSlotIndex: 3,
			quantity: 1
		}
	],
	purpose: 'Scout',
	spec: {
		armor: 20,
		beamBonus: 1,
		cost: {
			ironium: 17,
			boranium: 2,
			germanium: 7,
			resources: 22
		},
		engine: {
			idealSpeed: 6,
			freeSpeed: 1,
			maxSafeSpeed: 9,
			fuelUsage: [0, 0, 20, 60, 100, 100, 105, 450, 750, 900, 1080]
		},
		estimatedRange: 2272,
		estimatedRangeFull: 2272,
		fuelCapacity: 300,
		hullType: 'Scout',
		initiative: 1,
		mass: 25,
		movement: 4,
		movementFull: 4,
		numBuilt: 1,
		numEngines: 1,
		numInstances: 1,
		reduceCloaking: 1,
		scanner: true,
		scanRange: 66,
		scanRangePen: 30,
		techLevel: {
			propulsion: 3,
			electronics: 1
		}
	}
};

export const santaMariaDesign: ShipDesign = {
	id: 0,
	gameId: 0,
	num: 2,
	playerNum: 1,
	originalPlayerNum: 0,
	name: 'Santa Maria',
	version: 0,
	hull: 'Colony Ship',
	hullSetNumber: 0,
	slots: [
		{
			hullComponent: 'Long Hump 6',
			hullSlotIndex: 1,
			quantity: 1
		},
		{
			hullComponent: 'Colonization Module',
			hullSlotIndex: 2,
			quantity: 1
		}
	],
	purpose: 'Colonizer',
	spec: {
		armor: 20,
		beamBonus: 1,
		cargoCapacity: 25,
		colonizer: true,
		cost: {
			ironium: 25,
			boranium: 9,
			germanium: 23,
			resources: 33
		},
		engine: {
			idealSpeed: 6,
			freeSpeed: 1,
			maxSafeSpeed: 9,
			fuelUsage: [0, 0, 20, 60, 100, 100, 105, 450, 750, 900, 1080]
		},
		estimatedRange: 623,
		estimatedRangeFull: 442,
		fuelCapacity: 200,
		hullType: 'Colonizer',
		mass: 61,
		movement: 4,
		movementFull: 3,
		numBuilt: 1,
		numEngines: 1,
		numInstances: 1,
		reduceCloaking: 1,
		scanner: true,
		scanRangePen: -1,
		techLevel: {
			propulsion: 3
		}
	}
};

export const cottonPickerDesign: ShipDesign = {
	id: 0,
	gameId: 0,
	num: 4,
	playerNum: 1,
	originalPlayerNum: 0,
	name: 'Cotton Picker',
	version: 0,
	hull: 'Mini-Miner',
	hullSetNumber: 0,
	slots: [
		{
			hullComponent: 'Long Hump 6',
			hullSlotIndex: 1,
			quantity: 1
		},
		{
			hullComponent: 'Rhino Scanner',
			hullSlotIndex: 2,
			quantity: 1
		},
		{
			hullComponent: 'Robo-Mini-Miner',
			hullSlotIndex: 3,
			quantity: 1
		},
		{
			hullComponent: 'Robo-Mini-Miner',
			hullSlotIndex: 4,
			quantity: 1
		}
	],
	purpose: 'Miner',
	spec: {
		armor: 130,
		beamBonus: 1,
		cost: {
			ironium: 90,
			germanium: 23,
			resources: 249
		},
		engine: {
			idealSpeed: 6,
			freeSpeed: 1,
			maxSafeSpeed: 9,
			fuelUsage: [0, 0, 20, 60, 100, 100, 105, 450, 750, 900, 1080]
		},
		estimatedRange: 69,
		estimatedRangeFull: 69,
		fuelCapacity: 210,
		hullType: 'Miner',
		mass: 574,
		miningRate: 8,
		movement: 2,
		movementFull: 2,
		numBuilt: 1,
		numEngines: 1,
		numInstances: 1,
		reduceCloaking: 1,
		scanner: true,
		scanRange: 50,
		techLevel: {
			propulsion: 3,
			construction: 2,
			electronics: 1
		}
	}
};

export const longRangeScout: Fleet = {
	id: 0,
	gameId: 0,
	type: MapObjectType.Fleet,
	position: { x: 0, y: 0 },
	num: 1,
	playerNum: 1,
	warpSpeed: 6,
	name: 'Long Range Scout #1',
	waypoints: [
		{
			position: {
				x: 0,
				y: 0
			},
			warpSpeed: 6,
			task: WaypointTask.None,
			transportTasks: {
				fuel: {},
				ironium: {},
				boranium: {},
				germanium: {},
				colonists: {}
			},
			targetType: MapObjectType.None
		}
	],
	planetNum: 0,
	baseName: 'Long Range Scout',
	cargo: {},
	fuel: 300,
	tokens: [{ designNum: 1, quantity: 1 }],
	heading: { x: 0, y: 0 },
	spec: {
		armor: 20,
		cost: {
			ironium: 17,
			boranium: 2,
			germanium: 7,
			resources: 22
		},
		engine: {
			idealSpeed: 6,
			freeSpeed: 1,
			maxSafeSpeed: 9,
			fuelUsage: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
		},
		estimatedRange: 2272,
		fuelCapacity: 300,
		mass: 25,
		reduceCloaking: 1,
		scanner: true,
		scanRange: 66,
		scanRangePen: 30,
		baseCloakedCargo: 25,
		basePacketSpeed: 0,
		massEmpty: 25,
		purposes: { Scout: true },
		totalShips: 1
	}
};

export const santaMaria: Fleet = {
	id: 0,
	gameId: 0,
	type: MapObjectType.Fleet,
	position: { x: 0, y: 0 },
	num: 2,
	playerNum: 1,
	name: 'Santa Maria #2',
	waypoints: [
		{
			position: { x: 0, y: 0 },
			warpSpeed: 6,
			task: WaypointTask.None,
			transportTasks: {
				fuel: {},
				ironium: {},
				boranium: {},
				germanium: {},
				colonists: {}
			},
			targetType: MapObjectType.None
		}
	],
	planetNum: 0,
	baseName: 'Santa Maria',
	cargo: {},
	fuel: 200,
	tokens: [{ designNum: 2, quantity: 1 }],
	warpSpeed: 0,
	heading: { x: 0, y: 0 },
	orbitingPlanetNum: 0,
	spec: {
		armor: 20,
		cargoCapacity: 25,
		colonizer: true,
		cost: {
			ironium: 25,
			boranium: 9,
			germanium: 23,
			resources: 33
		},
		engine: {
			idealSpeed: 6,
			freeSpeed: 1,
			maxSafeSpeed: 9,
			fuelUsage: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
		},
		estimatedRange: 623,
		fuelCapacity: 200,
		mass: 61,
		reduceCloaking: 1,
		scanner: true,
		scanRangePen: -1,
		baseCloakedCargo: 61,
		basePacketSpeed: 0,
		massEmpty: 61,
		purposes: { Colonizer: true },
		totalShips: 1
	}
};

export const cottonPicker: Fleet = {
	id: 0,
	gameId: 0,
	type: MapObjectType.Fleet,
	position: { x: 0, y: 0 },
	num: 4,
	playerNum: 1,
	name: 'Cotton Picker #4',
	waypoints: [
		{
			position: { x: 0, y: 0 },
			warpSpeed: 6,
			task: WaypointTask.None,
			transportTasks: {
				fuel: {},
				ironium: {},
				boranium: {},
				germanium: {},
				colonists: {}
			},
			targetType: MapObjectType.None
		}
	],
	planetNum: 0,
	baseName: 'Cotton Picker',
	cargo: {},
	fuel: 210,
	tokens: [{ designNum: 4, quantity: 1 }],
	warpSpeed: 0,
	heading: { x: 0, y: 0 },
	orbitingPlanetNum: 0,
	spec: {
		armor: 130,
		cost: {
			ironium: 90,
			germanium: 23,
			resources: 249
		},
		engine: {
			idealSpeed: 6,
			freeSpeed: 1,
			maxSafeSpeed: 9,
			fuelUsage: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
		},
		estimatedRange: 69,
		fuelCapacity: 210,
		mass: 574,
		miningRate: 8,
		reduceCloaking: 1,
		scanner: true,
		scanRange: 50,
		baseCloakedCargo: 574,
		basePacketSpeed: 0,
		massEmpty: 574,
		purposes: {
			Miner: true
		},
		totalShips: 1
	}
};

// make vitest complain about at test file with no tests
describe('Mock test', () => {
	it('mock', () => {});
});
