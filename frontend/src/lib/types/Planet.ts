import type { Cargo } from './Cargo';
import type { Fleet } from './Fleet';
import type { Hab } from './Hab';
import type { MapObject } from './MapObject';
import type { Mineral } from './Mineral';

export const Unexplored = -1;

export interface Planet extends MapObject {
	hab?: Hab;
	baseHab?: Hab;
	terraformedAmount?: Hab;
	mineralConcentration?: Mineral;
	mineYears?: Mineral;
	cargo?: Cargo;
	population?: number;
	playerID?: number;
	mines?: number;
	factories?: number;
	defenses?: number;
	contributesOnlyLeftoverToResearch?: boolean;
	homeworld?: boolean;
	scanner?: boolean;
	reportAge: number;
	productionQueue?: ProductionQueueItem[];
	starbase?: Fleet;
	spec?: PlanetSpec;
}

export interface ProductionQueueItem {
	id?: number;
	createdAt?: string;
	updatedAt?: string;

	type: QueueItemType;
	quantity: number;
	designName?: string;
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
	terraformAmount: boolean;
	hasMassDriver: boolean;
	hasStarbase: boolean;
	dockCapacity: number;
}
