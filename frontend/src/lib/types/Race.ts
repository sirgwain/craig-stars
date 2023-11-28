import type { Cost } from './Cost';
import type { PartialEnumDictionary } from './EnumDictionary';
import { HabTypes, type Hab, type HabType } from './Hab';
import { type QueueItemType, QueueItemTypes } from './QueueItemType';
import type { TechLevel } from './TechLevel';

export interface Race {
	id?: number;
	createdAt?: string;
	updatedAt?: string;

	name: string;
	pluralName: string;
	spendLeftoverPointsOn?: string;
	prt: PRT;
	lrts: number;
	habLow: Hab;
	habHigh: Hab;
	growthRate: number;
	popEfficiency: number;
	factoryOutput: number;
	factoryCost: number;
	factoriesCostLess?: boolean;
	numFactories: number;
	mineOutput: number;
	mineCost: number;
	numMines: number;
	researchCost: ResearchCost;
	techsStartHigh?: boolean;
	immuneGrav?: boolean;
	immuneTemp?: boolean;
	immuneRad?: boolean;
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
	None = ''
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

export const lrts = [
	LRT.IFE,
	LRT.TT,
	LRT.ARM,
	LRT.ISB,
	LRT.GR,
	LRT.UR,
	LRT.NRSE,
	LRT.OBRM,
	LRT.NAS,
	LRT.LSP,
	LRT.BET,
	LRT.RS,
	LRT.MA,
	LRT.CE
];

export enum SpendLeftoverPointsOn {
	SurfaceMinerals = 'SurfaceMinerals',
	MineralConcentrations = 'MineralConcentrations',
	Mines = 'Mines',
	Factories = 'Factories',
	Defenses = 'Defenses'
}

export enum ResearchCostLevel {
	Extra = 'Extra',
	Standard = 'Standard',
	Less = 'Less'
}

export interface ResearchCost {
	energy: ResearchCostLevel;
	weapons: ResearchCostLevel;
	propulsion: ResearchCostLevel;
	construction: ResearchCostLevel;
	electronics: ResearchCostLevel;
	biotechnology: ResearchCostLevel;
}

export type RaceSpec = {
	newTechCostFactor: number;
	miniaturizationMax: number;
	miniaturizationPerLevel: number;
	builtInScannerMultiplier: number;
	armorStrengthFactor: number;
	scanRangeFactor: number;
	habCenter?: Hab;
	costs: PartialEnumDictionary<QueueItemType, Cost>;
	startingTechLevels?: TechLevel;
	startingPlanets?: StartingPlanet[];
	techCostOffset: TechCostOffset;
	mineralsPerSingleMineralPacket?: number;
	mineralsPerMixedMineralPacket?: number;
	packetResourceCost?: number;
	packetMineralCostFactor?: number;
	packetReceiverFactor?: number;
	packetDecayFactor?: number;
	packetOverSafeWarpPenalty?: number;
	packetBuiltInScanner?: boolean;
	detectPacketDestinationStarbases?: boolean;
	detectAllPackets?: boolean;
	packetTerraformChance?: number;
	packetPermaformChance?: number;
	packetPermaTerraformSizeUnit?: number;
	canGateCargo?: boolean;
	canDetectStargatePlanets?: boolean;
	shipsVanishInVoid?: boolean;
	techsCostExtraLevel?: number;
	freighterGrowthFactor?: number;
	growthFactor?: number;
	maxPopulationOffset?: number;
	builtInCloakUnits?: number;
	stealsResearch?: StealsResearch;
	freeCargoCloaking?: boolean;
	mineFieldsAreScanners?: boolean;
	mineFieldRateMoveFactor?: number;
	mineFieldSafeWarpBonus?: number;
	mineFieldMinDecayFactor?: number;
	mineFieldBaseDecayRate?: number;
	mineFieldPlanetDecayRate?: number;
	mineFieldMaxDecayRate?: number;
	canDetonateMineFields?: boolean;
	mineFieldDetonateDecayRate?: number;
	discoverDesignOnScan?: boolean;
	canRemoteMineOwnPlanets?: boolean;
	invasionAttackBonus?: number;
	invasionDefendBonus?: number;
	movementBonus?: number;
	instaforming?: boolean;
	permaformChance?: number;
	permaformPopulation?: number;
	repairFactor?: number;
	starbaseRepairFactor?: number;
	innateMining?: boolean;
	innateResources?: boolean;
	innateScanner?: boolean;
	innatePopulationFactor?: number;
	canBuildDefenses?: boolean;
	livesOnStarbases?: boolean;
	fuelEfficiencyOffset?: number;
	terraformCostOffset?: Cost;
	mineralAlchemyCostOffset?: number;
	scrapMineralOffset?: number;
	scrapMineralOffsetStarbase?: number;
	scrapResourcesOffset?: number;
	scrapResourcesOffsetStarbase?: number;
	startingPopulationFactor?: number;
	starbaseBuiltInCloakUnits?: number;
	starbaseCostFactor?: number;
	researchFactor?: number;
	researchSplashDamage?: number;
	shieldStrengthFactor?: number;
	shieldRegenerationRate?: number;
	engineFailureRate?: number;
	engineReliableSpeed?: number;
};

declare interface StealsResearch {
	energy?: number;
	weapons?: number;
	propulsion?: number;
	construction?: number;
	electronics?: number;
	biotechnology?: number;
}

declare interface TechCostOffset {
	engine?: number;
	beamWeapon?: number;
	torpedo?: number;
	bomb?: number;
	planetaryDefense?: number;
}

declare interface StartingPlanet {
	population?: number;
	habPenaltyFactor?: number;
	hasStargate?: boolean;
	hasMassDriver?: boolean;
	starbaseDesignName?: string;
	starbaseHull?: string;
	startingFleets?: StartingFleet[];
}

declare interface StartingFleet {
	name?: string;
	hullName?: string;
	hullSetNumber?: number;
	purpose?: string;
}

export const humanoid = (): Race => ({
	name: 'Humanoid',
	pluralName: 'Humanoids',
	spendLeftoverPointsOn: SpendLeftoverPointsOn.SurfaceMinerals,
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
		energy: ResearchCostLevel.Standard,
		weapons: ResearchCostLevel.Standard,
		propulsion: ResearchCostLevel.Standard,
		construction: ResearchCostLevel.Standard,
		electronics: ResearchCostLevel.Standard,
		biotechnology: ResearchCostLevel.Standard
	},
	spec: {
		newTechCostFactor: 1,
		miniaturizationMax: 0.75,
		miniaturizationPerLevel: 0.04,
		builtInScannerMultiplier: 20,
		scanRangeFactor: 1,
		habCenter: {
			grav: 50,
			temp: 50,
			rad: 50
		},
		costs: {
			[QueueItemTypes.AutoDefenses]: {
				ironium: 5,
				boranium: 5,
				germanium: 5,
				resources: 15
			},
			[QueueItemTypes.AutoFactories]: {
				germanium: 4,
				resources: 10
			},
			[QueueItemTypes.AutoMaxTerraform]: {
				resources: 100
			},
			[QueueItemTypes.AutoMinTerraform]: {
				resources: 100
			},
			[QueueItemTypes.AutoMineralAlchemy]: {
				resources: 100
			},
			[QueueItemTypes.AutoMineralPacket]: {
				ironium: 40,
				boranium: 40,
				germanium: 40,
				resources: 10
			},
			[QueueItemTypes.AutoMines]: {
				resources: 5
			},
			[QueueItemTypes.BoraniumMineralPacket]: {
				boranium: 100,
				resources: 10
			},
			[QueueItemTypes.Defenses]: {
				ironium: 5,
				boranium: 5,
				germanium: 5,
				resources: 15
			},
			[QueueItemTypes.Factory]: {
				germanium: 4,
				resources: 10
			},
			[QueueItemTypes.GermaniumMineralPacket]: {
				germanium: 100,
				resources: 10
			},
			[QueueItemTypes.IroniumMineralPacket]: {
				ironium: 100,
				resources: 10
			},
			[QueueItemTypes.Mine]: {
				resources: 5
			},
			[QueueItemTypes.MineralAlchemy]: {
				resources: 100
			},
			[QueueItemTypes.MixedMineralPacket]: {
				ironium: 40,
				boranium: 40,
				germanium: 40,
				resources: 10
			},
			[QueueItemTypes.PlanetaryScanner]: {
				ironium: 10,
				boranium: 10,
				germanium: 70,
				resources: 100
			},
			[QueueItemTypes.TerraformEnvironment]: {
				resources: 100
			}
		},
		startingTechLevels: {
			energy: 3,
			weapons: 3,
			propulsion: 3,
			construction: 3,
			electronics: 3,
			biotechnology: 3
		},
		startingPlanets: [
			{
				population: 25000,
				starbaseDesignName: 'Starbase',
				starbaseHull: 'Space Station',
				startingFleets: [
					{
						name: 'Long Range Scout',
						hullName: 'Scout',
						purpose: 'Scout'
					},
					{
						name: 'Santa Maria',
						hullName: 'Colony Ship',
						purpose: 'Colonizer'
					},
					{
						name: 'Teamster',
						hullName: 'Medium Freighter',
						purpose: 'Freighter'
					},
					{
						name: 'Cotton Picker',
						hullName: 'Mini-Miner',
						purpose: 'Miner'
					},
					{
						name: 'Armored Probe',
						hullName: 'Scout',
						hullSetNumber: 1,
						purpose: 'FighterScout'
					},
					{
						name: 'Stalwart Defender',
						hullName: 'Destroyer',
						purpose: 'Fighter'
					}
				]
			}
		],
		techCostOffset: {},
		mineralsPerSingleMineralPacket: 100,
		mineralsPerMixedMineralPacket: 40,
		packetResourceCost: 10,
		packetMineralCostFactor: 1,
		packetReceiverFactor: 1,
		packetDecayFactor: 1,
		packetPermaTerraformSizeUnit: 100,
		shipsVanishInVoid: true,
		techsCostExtraLevel: 4,
		growthFactor: 1,
		maxPopulationOffset: 0.2,
		stealsResearch: {},
		mineFieldMinDecayFactor: 1,
		mineFieldBaseDecayRate: 0.02,
		mineFieldPlanetDecayRate: 0.04,
		mineFieldMaxDecayRate: 0.5,
		mineFieldDetonateDecayRate: 0.25,
		invasionAttackBonus: 1.1,
		invasionDefendBonus: 1,
		repairFactor: 1,
		starbaseRepairFactor: 1,
		innatePopulationFactor: 1,
		canBuildDefenses: true,
		terraformCostOffset: {},
		startingPopulationFactor: 1,
		starbaseCostFactor: 1,
		researchFactor: 1,
		armorStrengthFactor: 1,
		shieldStrengthFactor: 1,
		engineReliableSpeed: 10
	}
});

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

// Get the habitability of this race for a given planet's hab value
export function getPlanetHabitability(race: Race, hab: Hab): number {
	let planetValuePoints = 0;
	let redValue = 0;
	let ideality = 10000;

	const habValues: [number, number, number] = [hab.grav ?? 0, hab.temp ?? 0, hab.rad ?? 0];
	const habCenters: [number, number, number] = [
		race.spec?.habCenter?.grav ?? 0,
		race.spec?.habCenter?.temp ?? 0,
		race.spec?.habCenter?.rad ?? 0
	];
	const habLows: [number, number, number] = [
		race.habLow.grav ?? 0,
		race.habLow.temp ?? 0,
		race.habLow.rad ?? 0
	];
	const habHighs: [number, number, number] = [
		race.habHigh.grav ?? 0,
		race.habHigh.temp ?? 0,
		race.habHigh.rad ?? 0
	];
	const immune: [boolean, boolean, boolean] = [
		race.immuneGrav ?? false,
		race.immuneTemp ?? false,
		race.immuneRad ?? false
	];

	let fromIdeal: number, tmp: number, habRadius: number, poorPlanetMod: number, habRed: number;

	for (let i = 0; i < habValues.length; i++) {
		const habValue: number = habValues[i];
		const habLower: number = habLows[i];
		const habUpper: number = habHighs[i];
		const habCenter: number = habCenters[i];

		if (immune[i]) {
			planetValuePoints += 10000;
		} else {
			if (habLower <= habValue && habUpper >= habValue) {
				// green planet
				fromIdeal = Math.abs(habValue - habCenter) * 100;
				if (habCenter > habValue) {
					habRadius = habCenter - habLower;
					fromIdeal /= habRadius;
					tmp = habCenter - habValue;
				} else {
					habRadius = habUpper - habCenter;
					fromIdeal /= habRadius;
					tmp = habValue - habCenter;
				}
				poorPlanetMod = tmp * 2 - habRadius;
				fromIdeal = 100 - fromIdeal;
				planetValuePoints += fromIdeal * fromIdeal;
				if (poorPlanetMod > 0) {
					ideality *= habRadius * 2 - poorPlanetMod;
					ideality /= habRadius * 2;
				}
			} else {
				// red planet
				if (habLower <= habValue) {
					habRed = habValue - habUpper;
				} else {
					habRed = habLower - habValue;
				}

				if (habRed > 15) {
					habRed = 15;
				}

				redValue += habRed;
			}
		}
	}

	if (redValue !== 0) {
		return -redValue;
	}

	planetValuePoints = Math.sqrt(planetValuePoints / 3.0) + 0.9;
	planetValuePoints = (planetValuePoints * ideality) / 10000;

	return Math.floor(planetValuePoints);
}

export function getHabWidth(race: Race) {
	return {
		grav: (race.habHigh.grav ?? 0) - (race.habLow.grav ?? 0),
		temp: (race.habHigh.temp ?? 0) - (race.habLow.temp ?? 0),
		rad: (race.habHigh.rad ?? 0) - (race.habLow.rad ?? 0)
	};
}

export function getHabChance(race: Race): number {
	const habWidth = getHabWidth(race);
	// do a straight calc of hab width, so if we have a hab with widths of 50, 50% of planets will be habitable
	// so we get (.5 * .5 * .5) = .125, or 1 in 8 planets
	const gravChance = race.immuneGrav ? 1.0 : habWidth.grav / 100.0;
	const tempChance = race.immuneTemp ? 1.0 : habWidth.temp / 100.0;
	const radChance = race.immuneRad ? 1.0 : habWidth.rad / 100.0;
	return gravChance * tempChance * radChance;
}

export function isImmune(race: Race, habType: HabType): boolean {
	switch (habType) {
		case HabTypes.Gravity:
			return race.immuneGrav ?? false;
		case HabTypes.Temperature:
			return race.immuneTemp ?? false;
		case HabTypes.Radiation:
			return race.immuneRad ?? false;
	}
}
