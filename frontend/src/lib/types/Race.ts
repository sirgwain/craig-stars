import type { HabType, LRT, PRT, Race } from './cs';
import {
	AR,
	ARM,
	BET,
	CA,
	CE,
	GR,
	Grav,
	HE,
	IFE,
	IS,
	ISB,
	IT,
	JoaT,
	LRTNone,
	LSP,
	MA,
	NAS,
	NRSE,
	OBRM,
	PP,
	QueueItemTypeAutoDefenses,
	QueueItemTypeAutoFactories,
	QueueItemTypeAutoMaxTerraform,
	QueueItemTypeAutoMineralAlchemy,
	QueueItemTypeAutoMineralPacket,
	QueueItemTypeAutoMines,
	QueueItemTypeAutoMinTerraform,
	QueueItemTypeBoraniumMineralPacket,
	QueueItemTypeDefenses,
	QueueItemTypeFactory,
	QueueItemTypeGenesisDevice,
	QueueItemTypeGermaniumMineralPacket,
	QueueItemTypeIroniumMineralPacket,
	QueueItemTypeMine,
	QueueItemTypeMineralAlchemy,
	QueueItemTypeMixedMineralPacket,
	QueueItemTypePlanetaryScanner,
	QueueItemTypeTerraformEnvironment,
	Rad,
	ResearchCostStandard,
	RS,
	SD,
	SpendLeftoverPointsOnSurfaceMinerals,
	SS,
	Temp,
	TT,
	UR,
	WM,
	type Hab
} from './cs';

export const lrts = [IFE, TT, ARM, ISB, GR, UR, NRSE, OBRM, NAS, LSP, BET, RS, MA, CE] as const;

export const humanoid = (): Race => ({
	id: 0,
	createdAt: '',
	updatedAt: '',
	name: 'Humanoid',
	pluralName: 'Humanoids',
	spendLeftoverPointsOn: SpendLeftoverPointsOnSurfaceMinerals,
	prt: JoaT,
	lrts: LRTNone,
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
		energy: ResearchCostStandard,
		weapons: ResearchCostStandard,
		propulsion: ResearchCostStandard,
		construction: ResearchCostStandard,
		electronics: ResearchCostStandard,
		biotechnology: ResearchCostStandard
	},
	spec: {
		newTechCostFactor: 1,
		miniaturizationMax: 0.75,
		miniaturizationPerLevel: 0.04,
		scanRangeFactor: 1,
		builtInScanner: { normalMulti: {}, penMulti: {} },
		habCenter: {
			grav: 50,
			temp: 50,
			rad: 50
		},
		costs: {
			[QueueItemTypeAutoDefenses]: {
				ironium: 5,
				boranium: 5,
				germanium: 5,
				resources: 15
			},
			[QueueItemTypeAutoFactories]: {
				germanium: 4,
				resources: 10
			},
			[QueueItemTypeAutoMaxTerraform]: {
				resources: 100
			},
			[QueueItemTypeAutoMinTerraform]: {
				resources: 100
			},
			[QueueItemTypeAutoMineralAlchemy]: {
				resources: 100
			},
			[QueueItemTypeAutoMineralPacket]: {
				ironium: 40,
				boranium: 40,
				germanium: 40,
				resources: 10
			},
			[QueueItemTypeAutoMines]: {
				resources: 5
			},
			[QueueItemTypeBoraniumMineralPacket]: {
				boranium: 100,
				resources: 10
			},
			[QueueItemTypeDefenses]: {
				ironium: 5,
				boranium: 5,
				germanium: 5,
				resources: 15
			},
			[QueueItemTypeFactory]: {
				germanium: 4,
				resources: 10
			},
			[QueueItemTypeGermaniumMineralPacket]: {
				germanium: 100,
				resources: 10
			},
			[QueueItemTypeIroniumMineralPacket]: {
				ironium: 100,
				resources: 10
			},
			[QueueItemTypeMine]: {
				resources: 5
			},
			[QueueItemTypeMineralAlchemy]: {
				resources: 100
			},
			[QueueItemTypeMixedMineralPacket]: {
				ironium: 40,
				boranium: 40,
				germanium: 40,
				resources: 10
			},
			[QueueItemTypePlanetaryScanner]: {
				ironium: 10,
				boranium: 10,
				germanium: 70,
				resources: 100
			},
			[QueueItemTypeGenesisDevice]: {
				ironium: 0,
				boranium: 0,
				germanium: 0,
				resources: 5000
			},
			[QueueItemTypeTerraformEnvironment]: {
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
						name: 'Armed Probe',
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
		innateScannerFactor: 1,
		canBuildDefenses: true,
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
		case HE:
			return 'Hyper Expansion';
		case SS:
			return 'Super Stealth';
		case WM:
			return 'Warmonger';
		case CA:
			return 'Claim Adjuster';
		case IS:
			return 'Inner Strength';
		case SD:
			return 'Space Demolition';
		case PP:
			return 'Packet Physics';
		case IT:
			return 'Interstellar Traveler';
		case AR:
			return 'Alternate Reality';
		case JoaT:
			return 'Jack of All Trades';
		default:
			return toString();
	}
};

export const getLabelForLRT = (lrt: LRT): string => {
	switch (lrt) {
		case IFE:
			return 'Improved Fuel Efficiency';
		case TT:
			return 'Total Terraforming';
		case ARM:
			return 'Advanced Remote Mining';
		case ISB:
			return 'Improved Starbases';
		case GR:
			return 'Generalized Research';
		case UR:
			return 'Ultimate Recycling';
		case NRSE:
			return 'No Ram Scoop Engines';
		case OBRM:
			return 'Only Basic Remote Mining';
		case NAS:
			return 'No Advanced Scanners';
		case LSP:
			return 'Low Starting Population';
		case BET:
			return 'Bleeding Edge Technology';
		case RS:
			return 'Regenerating Shields';
		case MA:
			return 'Mineral Alchemy';
		case CE:
			return 'Cheap Engines';
		default:
			return toString();
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
		case Grav:
			return race.immuneGrav ?? false;
		case Temp:
			return race.immuneTemp ?? false;
		case Rad:
			return race.immuneRad ?? false;
	}

	return false;
}
