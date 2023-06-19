import type { Cargo } from './Cargo';
import type { Hab } from './Hab';
import type { MapObject } from './MapObject';
import type { Mineral } from './Mineral';

export interface Planet extends MapObject {
	hab: Hab;
	baseHab: Hab;
	terraformedAmount: Hab;
	mineralConcentration: Mineral;
	mineYears: Mineral;
	cargo: Cargo;
	playerID: number;
	mines?: number;
	factories?: number;
	defenses?: number;
	contributesOnlyLeftoverToResearch: boolean;
	homeworld?: boolean;
	scanner?: boolean;
	reportAge?: number;
	spec?: PlanetSpec;
}

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
}
