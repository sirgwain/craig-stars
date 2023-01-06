import type { Cargo } from './Cargo';
import type { Fleet } from './Fleet';
import type { Hab } from './Hab';
import { MapObjectType, type MapObject } from './MapObject';
import type { Mineral } from './Mineral';
import type { Vector } from './Vector';

export const Unexplored = -1;

export interface Planet extends MapObject {
	hab?: Hab;
	baseHab?: Hab;
	terraformedAmount?: Hab;
	mineralConcentration?: Mineral;
	mineYears?: Mineral;
	cargo?: Cargo;
	population?: number;
	mines?: number;
	factories?: number;
	defenses?: number;
	contributesOnlyLeftoverToResearch?: boolean;
	homeworld?: boolean;
	scanner?: boolean;
	reportAge: number;
	productionQueue?: ProductionQueueItem[];
	spec: PlanetSpec;
}

export interface ProductionQueueItem {
	id?: number;
	createdAt?: string;
	updatedAt?: string;

	type: QueueItemType;
	quantity: number;
	designName?: string;
}

export class CommandedPlanet implements Planet {
	id = 0;
	gameId = 0;
	createdAt?: string | undefined;
	updatedAt?: string | undefined;
	readonly type = MapObjectType.Planet;

	hab: Hab = { grav: 0, temp: 0, rad: 0 };
	baseHab: Hab = { grav: 0, temp: 0, rad: 0 };
	terraformedAmount = { grav: 0, temp: 0, rad: 0 };
	mineralConcentration: Mineral = { ironium: 0, boranium: 0, germanium: 0 };
	mineYears: Mineral = { ironium: 0, boranium: 0, germanium: 0 };
	cargo: Cargo = { ironium: 0, boranium: 0, germanium: 0, colonists: 0 };
	mines = 0;
	factories = 0;
	defenses = 0;
	contributesOnlyLeftoverToResearch = false;
	homeworld = false;
	scanner = false;
	reportAge = 0;
	productionQueue: ProductionQueueItem[] = [];
	position: Vector = { x: 0, y: 0 };
	name = '';
	num = 0;
	playerNum = 0;
	spec: PlanetSpec = {
		habitability: 0,
		habitabilityTerraformed: 0,
		maxMines: 0,
		maxPossibleMines: 0,
		maxFactories: 0,
		maxPossibleFactories: 0,
		maxDefenses: 0,
		populationDensity: 0,
		maxPopulation: 0,
		growthAmount: 0,
		mineralOutput: { ironium: 0, boranium: 0, germanium: 0 },
		resourcesPerYear: 0,
		resourcesPerYearAvailable: 0,
		resourcesPerYearResearch: 0,
		defense: '',
		defenseCoverage: 0,
		defenseCoverageSmart: 0,
		scanner: '',
		scanRange: 0,
		scanRangePen: 0,
		canTerraform: false,
		terraformAmount: { grav: 0, temp: 0, rad: 0 },
		hasMassDriver: false,
		hasStarbase: false,
		dockCapacity: 0
	};
}

export const fromQueueItemType = (type: QueueItemType): ProductionQueueItem => ({
	type,
	quantity: 0
});

export enum QueueItemType {
	AutoMines = 'AutoMines',
	AutoFactories = 'AutoFactories',
	AutoDefenses = 'AutoDefenses',
	AutoMineralAlchemy = 'AutoMineralAlchemy',
	AutoMinTerraform = 'AutoMinTerraform',
	AutoMaxTerraform = 'AutoMaxTerraform',
	AutoMineralPacket = 'AutoMineralPacket',
	Factory = 'Factory',
	Mine = 'Mine',
	Defenses = 'Defenses',
	MineralAlchemy = 'MineralAlchemy',
	TerraformEnvironment = 'TerraformEnvironment',
	IroniumMineralPacket = 'IroniumMineralPacket',
	BoraniumMineralPacket = 'BoraniumMineralPacket',
	GermaniumMineralPacket = 'GermaniumMineralPacket',
	MixedMineralPacket = 'MixedMineralPacket',
	ShipToken = 'ShipToken',
	Starbase = 'Starbase'
}

export const stringToQueueItemType = (value: string): QueueItemType => {
	return QueueItemType[value as keyof typeof QueueItemType];
};

/**
 * Determine if a ProductionQueueItem is an auto item
 * @param type The type to check
 * @returns
 */
export const isAuto = (type: QueueItemType): boolean => {
	switch (type) {
		case QueueItemType.AutoMines:
		case QueueItemType.AutoFactories:
		case QueueItemType.AutoDefenses:
		case QueueItemType.AutoMineralAlchemy:
		case QueueItemType.AutoMinTerraform:
		case QueueItemType.AutoMaxTerraform:
		case QueueItemType.AutoMineralPacket:
			return true;
		default:
			return false;
	}
};

export interface PlanetSpec {
	habitability?: number;
	habitabilityTerraformed?: number;
	maxMines: number;
	maxPossibleMines: number;
	maxFactories: number;
	maxPossibleFactories: number;
	maxDefenses: number;
	populationDensity: number;
	maxPopulation: number;
	growthAmount: number;
	mineralOutput: Mineral;
	resourcesPerYear: number;
	resourcesPerYearAvailable: number;
	resourcesPerYearResearch: number;
	defense: string;
	defenseCoverage: number;
	defenseCoverageSmart: number;
	scanner: string;
	scanRange: number;
	scanRangePen: number;
	canTerraform: boolean;
	terraformAmount?: Hab;
	hasMassDriver: boolean;
	hasStarbase: boolean;
	dockCapacity: number;
}
