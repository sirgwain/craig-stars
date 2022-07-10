import type { Hab } from './Hab';
import type { Cost } from './Cost';
import type { QueueItemType } from './Planet';
import type { TechLevel } from './Player';
import type { EnumDictionary } from './EnumDictionary';

export interface Race {
	id?: number;
	createdAt?: string;
	updatedat?: string;
	deletedAt?: null;
	playerId?: number;
	name: string;
	pluralName: string;
	prt: PRT;
	lrts: number;
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
	spec?: RaceSpec;
}

export enum PRT {
	/// Hyper Expansion
	HE = 'HE',

	/// Super Stealth
	SS = 'SS',

	/// Warmonger
	WM = 'WM',

	/// Claim Adjuster
	CA = 'CA',

	/// Inner Strength
	IS = 'IS',

	/// Space Demolition
	SD = 'SD',

	/// Packet Physics
	PP = 'PP',

	/// Interstellar Traveler
	IT = 'IT',

	/// Alternate Reality
	AR = 'AR',

	/// Jack of All Trades
	JoaT = 'JoaT',

	/// This is only for tech requirements
	None = 'None'
}

export enum LRT {
	// Only used for TechRequirements
	None = 0,

	// Improved Fuel Efficiency
	IFE = 1 << 0,

	// Total Terraforming
	TT = 1 << 1,

	// Advanced Remote Mining
	ARM = 1 << 2,

	// Improved Starbases
	ISB = 1 << 3,

	// Generalized Research
	GR = 1 << 4,

	// Ultimate Recycling
	UR = 1 << 5,

	// No Ramscoop Engines
	NRSE = 1 << 6,

	// Only Basic Remote Mining
	OBRM = 1 << 7,

	// No Advanced Scanners
	NAS = 1 << 8,

	// Low Starting Population
	LSP = 1 << 9,

	// Bleeding Edge Technology
	BET = 1 << 10,

	// Regenerating Shields
	RS = 1 << 11,

	// Mineral Alchemy
	MA = 1 << 12,

	// Cheap Engines
	CE = 1 << 13
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
	costs: EnumDictionary<QueueItemType, Cost>;
	startingTechLevels: TechLevel;
	startingFleets: StartingFleet[];
	startingPlanets: StartingPlanet[];
	techCostOffset: TechLevel;
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
	stealsResearch: TechLevel;
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
	terraformCostOffset: TechLevel;
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

export const humanoid: Race = {
	name: 'Humanoid',
	pluralName: 'Humanoids',
	prt: PRT.JoaT,
	lrts: LRT.None,
	habLow: { grav: 15, temp: 15, rad: 15 },
	habHigh: { grav: 85, temp: 85, rad: 85 },
	growthRate: 15,
	popEfficiency: 10,
	factoryOutput: 10,
	factoryCost: 10,
	numFactories: 10,
	mineOutput: 10,
	mineCost: 5,
	numMines: 10,
	researchCost: {
		energy: 'Standard',
		weapons: 'Standard',
		propulsion: 'Standard',
		construction: 'Standard',
		electronics: 'Standard',
		biotechnology: 'Standard'
	}
};

export const getLabelForPRT = (prt: PRT): string => {
	switch (prt) {
		case PRT.HE:
			return 'Hyper Expansion';
		case PRT.SS:
			return 'Super Stealth';
		case PRT.WM:
			return 'Warmonger';
		case PRT.CA:
			return 'Claim Adjuster';
		case PRT.IS:
			return 'Inner Strength';
		case PRT.SD:
			return 'Space Demolition';
		case PRT.PP:
			return 'Packet Physics';
		case PRT.IT:
			return 'Interstellar Traveler';
		case PRT.AR:
			return 'Alternate Reality';
		case PRT.JoaT:
			return 'Jack of All Trades';
		default:
			return prt.toString();
	}
};

export const getLabelForLRT = (lrt: LRT): string => {
	switch (lrt) {
		case LRT.IFE:
			return 'Improved Fuel Efficiency';
		case LRT.TT:
			return 'Total Terraforming';
		case LRT.ARM:
			return 'Advanced Remote Mining';
		case LRT.ISB:
			return 'Improved Starbases';
		case LRT.GR:
			return 'Generalized Research';
		case LRT.UR:
			return 'Ultimate Recycling';
		case LRT.NRSE:
			return 'No Ram Scoop Engines';
		case LRT.OBRM:
			return 'Only Basic Remote Mining';
		case LRT.NAS:
			return 'No Advanced Scanners';
		case LRT.LSP:
			return 'Low Starting Population';
		case LRT.BET:
			return 'Bleeding Edge Technology';
		case LRT.RS:
			return 'Regenerating Shields';
		case LRT.MA:
			return 'Mineral Alchemy';
		case LRT.CE:
			return 'Cheap Engines';
		default:
			return lrt.toString();
	}
};
