import type { Cost } from './Cost';
import type { MineFieldType } from './MineField';
import type { Engine, TechHullType } from './Tech';
import type { TechLevel } from './TechLevel';

export type ShipDesignPurpose =
	| ''
	| 'Scout'
	| 'Colonizer'
	| 'Bomber'
	| 'StructureBomber'
	| 'SmartBomber'
	| 'Fighter'
	| 'FighterScout'
	| 'CapitalShip'
	| 'Freighter'
	| 'ColonistFreighter'
	| 'FuelFreighter'
	| 'MultiPurposeFreighter'
	| 'ArmedFreighter'
	| 'Miner'
	| 'Terraformer'
	| 'DamageMineLayer'
	| 'SpeedMineLayer'
	| 'Starbase'
	| 'FuelDepot'
	| 'StarbaseQuarter'
	| 'StarbaseHalf'
	| 'PacketThrower'
	| 'Stargater'
	| 'Fort'
	| 'StarterColony';

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
	mysteryTrader?: boolean;
	slots: ShipDesignSlot[];
	purpose?: string;
	reportAge?: number;
	spec: ShipDesignSpec;
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

export type ShipDesignSpec = {
	additionalMassDrivers?: number;
	armor?: number;
	basePacketSpeed?: number;
	beamBonus?: number;
	beamDefense?: number;
	bomber?: boolean;
	bombs?: Bomb[];
	canJump?: boolean;
	canLayMines?: boolean;
	canStealFleetCargo?: boolean;
	canStealPlanetCargo?: boolean;
	cargoCapacity?: number;
	cloakPercent?: number;
	cloakPercentFullCargo?: number;
	cloakUnits?: number;
	colonizer?: boolean;
	cost: Cost;
	engine: Engine;
	estimatedRange?: number;
	estimatedRangeFull?: number;
	fuelCapacity?: number;
	fuelGeneration?: number;
	hasWeapons?: boolean;
	hullType?: TechHullType;
	immuneToOwnDetonation?: boolean;
	initiative?: number;
	innateScanRangePenFactor?: number;
	mass?: number;
	massDriver?: string;
	maxHullMass?: number;
	maxPopulation?: number;
	maxRange?: number;
	mineLayingRateByMineType?: Record<MineFieldType, number>;
	mineSweep?: number;
	miningRate?: number;
	movement?: number;
	movementBonus?: number;
	movementFull?: number;
	numBuilt?: number;
	numEngines?: number;
	numInstances?: number;
	orbitalConstructionModule?: boolean;
	powerRating?: number;
	radiating?: boolean;
	reduceCloaking?: number;
	reduceMovement?: number;
	repairBonus?: number;
	retroBombs?: Bomb[];
	safeHullMass?: number;
	safePacketSpeed?: number;
	safeRange?: number;
	scanner?: boolean;
	scanRange?: number;
	scanRangePen?: number;
	shields?: number;
	smartBombs?: Bomb[];
	spaceDock?: number;
	starbase?: boolean;
	stargate?: string;
	techLevel?: TechLevel;
	terraformRate?: number;
	torpedoBonus?: number;
	torpedoJamming?: number;
	weaponSlots?: ShipDesignSlot[];
};
