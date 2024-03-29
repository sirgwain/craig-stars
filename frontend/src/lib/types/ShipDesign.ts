import type { Cost } from './Cost';
import type { Engine } from './Tech';
import type { TechLevel } from './TechLevel';

export type ShipDesign = {
	id?: number;
	gameId: number;
	createdAt?: Date;
	updatedAt?: Date;
	num?: number;
	playerNum: number;
	originalPlayerNum: number;
	name: string;
	version: number;
	hull: string;
	hullSetNumber: number;
	cannotDelete?: boolean;
	slots: ShipDesignSlot[];
	purpose?: string;
	reportAge?: number;
	spec: Spec;
};

export type ShipDesignSlot = {
	hullComponent: string;
	hullSlotIndex: number;
	quantity: number;
};

export type Bomb = {
	quantity?: number;
	killRate?: number;
	minKillRate?: number;
	structureDestroyRate?: number;
	unterraformRate?: number;
};

export type Spec = {
	engine: Engine;
	hullType?: string;
	techLevel: TechLevel;
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
	torpedoJamming?: number;
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
	shields?: number;
	colonizer?: boolean;
	starbase?: boolean;
	canLayMines?: boolean;
	spaceDock?: number;
	miningRate?: number;
	terraformRate?: number;
	mineSweep?: number;
	cloakPercent?: number;
	cloakPercentFullCargo?: number;
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
	maxPopulation?: number;
	numInstances?: number;
	numBuilt?: number;
	estimatedRange?: number;
	estimatedRangeFull?: number;
};
