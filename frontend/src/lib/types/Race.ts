import type { Hab } from './Hab';
import type { TechLevels } from './Player';

export interface Race {
	id: number;
	createdAt: string;
	updatedat: string;
	deletedAt: null;
	playerId: number;
	name: string;
	pluralName: string;
	prt: string;
	habLow: Hab;
	habHigh: Hab;
	growthRate: number;
	popEfficiency: number;
	factoryOutput: number;
	factoryCost: number;
	numFactories: number;
	mineOutput: number;
	mineCost: number;
	numMines: number;
	researchCost: ResearchCost;
	spec: RaceSpec;
}

export interface ResearchCost {
	energy: string;
	weapons: string;
	propulsion: string;
	construction: string;
	electronics: string;
	biotechnology: string;
}

export interface RaceSpec {
	startingTechLevels: TechLevels;
	startingFleets: StartingFleet[];
	startingPlanets: StartingPlanet[];
	techCostOffset: TechLevels;
	mineralsPerSingleMineralPacket: number;
	mineralsPerMixedMineralPacket: number;
	packetResourceCost: number;
	packetMineralCostFactor: number;
	packetReceiverFactor: number;
	packetDecayFactor: number;
	packetPermaTerraformSizeUnit: number;
	shipsVanishInVoid: boolean;
	builtInScannerMultiplier: number;
	techsCostExtraLevel: number;
	growthFactor: number;
	maxPopulationOffset: number;
	stealsResearch: TechLevels;
	mineFieldMinDecayFactor: number;
	mineFieldBaseDecayRate: number;
	mineFieldPlanetDecayRate: number;
	mineFieldMaxDecayRate: number;
	mineFieldDetonateDecayRate: number;
	invasionAttackBonus: number;
	invasionDefendBonus: number;
	repairFactor: number;
	starbaseRepairFactor: number;
	canBuildDefenses: boolean;
	terraformCostOffset: TechLevels;
	starbaseCostFactor: number;
}

export interface StartingFleet {
	name: string;
	hullName: string;
	purpose: string;
}

export interface StartingPlanet {
	population: number;
	habPenaltyFactor: number;
	hasStargate: boolean;
	hasMassDriver: boolean;
	startingFleets: null;
}
