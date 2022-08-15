import type { Cost } from './Cost';

export interface ShipDesign {
	id: number;
	createdAt: string;
	updatedAt: string;

	gameId: number;
	playerId: number;
	playerNum: number;
	name: string;
	version: number;
	hull: string;
	hullSetNumber: number;
	armor: number;
	slots: ShipDesignSlot[];
	spec: Spec;
}

export interface ShipDesignSlot {
	hullComponent: string;
	hullSlotIndex: number;
	quantity: number;
}

export interface Spec {
	weaponSlots: null;
	computed: boolean;
	engine: string;
	numEngines: number;
	cost: Cost;
	mass: number;
	armor: number;
	fuelCapacity: number;
	scanRange: number;
	scanRangePen: number;
	torpedoInaccuracyFactor: number;
	initiative: number;
	movement: number;
	bombs: null;
	smartBombs: null;
	retroBombs: null;
	scanner: boolean;
	mineLayingRateByMineType: null;
}
