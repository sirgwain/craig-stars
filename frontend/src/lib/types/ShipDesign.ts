import type { Cost } from './Cost';

export interface ShipDesign {
	id: number;
	createdAt: string;
	updatedAt: string;

	gameId: number;
	num: number;
	playerNum: number;
	name: string;
	uuid: string;
	version: number;
	hull: string;
	hullSetNumber: number;
	armor: number;
	slots: ShipDesignSlot[];
	purpose: string;
	spec: Spec;
}

export type ShipDesignSlot = {
	hullComponent: string;
	hullSlotIndex: number;
	quantity: number;
}

export type Bomb = {
	quantity?: number;
	killRate?: number;
	minKillRate?: number;
	structureDestroyRate?: number;
	unterraformRate?: number;
};

export type Spec = {
	idealSpeed?: number;
	engine?: string;
	fuelUsage?: [
		number,
		number,
		number,
		number,
		number,
		number,
		number,
		number,
		number,
		number,
		number
	];
	numEngines?: number;
	cost?: Cost;
	mass?: number;
	armor?: number;
	fuelCapacity?: number;
	cargoCapacity?: number;
	cloakUnits?: number;
	scanRange?: number;
	scanRangePen?: number;
	repairBonus?: number;
	torpedoInaccuracyFactor?: number;
	initiative?: number;
	movement?: number;
	powerRating?: number;
	bomber?: boolean;
	bombs?: Bomb[];
	smartBombs?: Bomb[];
	retroBombs?: Bomb[];
	scanner?: boolean;
	immuneToOwnDetonation?: boolean;
	mineLayingRateByMineType?: { [mineFieldType: string]: number };
	shield?: number;
	colonizer?: boolean;
	canLayMines?: boolean;
	spaceDock?: number;
	miningRate?: number;
	terraformRate?: number;
	mineSweep?: number;
	cloakPercent?: number;
	reduceCloaking?: number;
	canStealFleetCargo?: boolean;
	canStealPlanetCargo?: boolean;
	orbitalConstructionModule?: boolean;
	hasWeapons?: boolean;
	weaponSlots?: ShipDesignSlot[];
	safeHullMass?: number;
	safeRange?: number;
	maxHullMass?: number;
	maxRange?: number;
};
