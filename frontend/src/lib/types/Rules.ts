import type { Cost } from './Cost';
import type { Size } from './Game';
import type { MineFieldStats, MineFieldType } from './MineField';
import type { LRT, PRT, SpendLeftoverPointsOn } from './Race';
import type { TechLevel } from './TechLevel';

export type RandomEvent = 'Comet' | 'MineralDeposit' | 'PlanetaryChange' | 'AncientArtifact';

export type CometSize = 'Small' | 'Medium' | 'Large' | 'Huge';

export type RepairRate =
	| 'None'
	| 'Moving'
	| 'Stopped'
	| 'Orbiting'
	| 'OrbitingOwnPlanet'
	| 'Starbase';

export type WormholeStability =
	| 'RockSolid'
	| 'Stable'
	| 'MostlyStable'
	| 'Average'
	| 'SlightlyVolatile'
	| 'Volatile'
	| 'ExtremelyVolatile';

export type StartingFleetHull =
	| ''
	| 'Colony Ship'
	| 'Destroyer'
	| 'Medium Freighter'
	| 'Mini-Colony Ship'
	| 'Mini Mine Layer'
	| 'Mini-Miner'
	| 'Midget-Miner'
	| 'Privateer'
	| 'Scout';

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

export type Rules = {
	id?: number;
	createdAt?: string; // ISO 8601 string for time
	updatedAt?: string;
	gameId?: number;
	cometStatsBySize?: Record<CometSize, CometStats>;
	fleetSafeSpeedExplosionChance?: number;
	invasionDefenseCoverageFactor?: number;
	lrtSpecs?: Partial<Record<LRT, LRTSpec>>;
	maxPopulation?: number;
	maxTechLevel?: number;
	mineFieldCloak?: number;
	mineFieldStatsByType?: Record<MineFieldType, MineFieldStats>;
	mineralDecayFactor?: number;
	minMaxPopulationPercent?: number;
	movesToRunAway?: number;
	mysteryTraderRules?: MysteryTraderRules;
	packetDecayRate?: Record<number, number>;
	packetMaxOverwarpSpeed?: number;
	packetMinDecay?: number;
	planetMinDistance?: number;
	populationOvercrowdDieoffRate?: number;
	populationOvercrowdDieoffRateMax?: number;
	populationScannerError?: number;
	prtSpecs?: Partial<Record<PRT, PRTSpec>>;
	raceStartingPoints?: number;
	radiatingImmune?: number;
	randomArtifactResearchBonusRange?: [number, number];
	randomCometMinYear?: number;
	randomCometMinYearPlayerWorld?: number;
	randomEventChances?: Record<RandomEvent, number>;
	randomMineralDepositBonusRange?: [number, number];
	remoteMiningMineOutput?: number;
	repairRates?: Record<RepairRate, number>;
	salvageDecayMin?: number;
	salvageDecayRate?: number;
	salvageFromBattleFactor?: number;
	scrapMineralAmount?: number;
	scrapResourceAmount?: number;
	showPublicScoresAfterYears?: number;
	smartDefenseCoverageFactor?: number;
	stargateMaxHullMassFactor?: number;
	stargateMaxRangeFactor?: number;
	tachyonCloakReduction?: number;
	tachyonMaxCloakReduction?: number;
	techsId?: number;
	techTradeChance?: number;
	torpedoSplashDamage?: number;
	wormholeCloak?: number;
	wormholePairsForSize?: Record<Size, number>;
	wormholeStatsByStability?: Record<WormholeStability, WormholeStats>;
} & CostRules &
	BattleRules &
	UniverseGenerationRules;

export type UniverseGenerationRules = {
	highRadMineralConcentrationBonusThreshold?: number;
	limitMineralConcentration?: number;
	maxExtraWorldDistance?: number;
	maxHab?: number;
	maxMineralConcentration?: number;
	maxStartingMineralConcentration?: number;
	maxStartingMineralSurface?: number;
	minExtraPlanetMineralConcentration?: number;
	minExtraWorldDistance?: number;
	minHab?: number;
	minHomeworldMineralConcentration?: number;
	minMineralConcentration?: number;
	minStartingMineralConcentration?: number;
	minStartingMineralSurface?: number;
	raceLeftoverPointsPerItem?: Record<SpendLeftoverPointsOn, number>;
	startingYear?: number;
	wormholeMinPlanetDistance?: number;
};

export type CostRules = {
	defenseCost?: Cost;
	factoryCostGermanium?: number;
	mineralAlchemyCost?: number;
	planetaryScannerCost?: Cost;
	starbaseComponentCostReduction?: number;
	terraformCost?: Cost;
	techBaseCost?: number[];
};

export type BattleRules = {
	beamRangeDropoff?: number;
	numBattleRounds?: number;
};

export type CometStats = {
	allMinerals?: number;
	allRandomMinerals?: number;
	bonusMinerals?: number;
	bonusRandomMinerals?: number;
	bonusMinConcentration?: number;
	bonusRandomConcentration?: number;
	bonusAffectsMinerals?: number;
	minTerraform?: number;
	randomTerraform?: number;
	affectsHabs?: number;
	popKilledPercent?: number;
};

export type MysteryTraderRules = {
	chanceSpawn?: number[];
	chanceMaxTechGetsPart?: number;
	chanceCourseChange?: number;
	chanceSpeedUpOnly?: number;
	chanceAgain?: number;
	minYear?: number;
	evenYearOnly?: boolean;
	minWarp?: number;
	maxWarp?: number;
	maxMysteryTraders?: number;
	requestedBoon?: number;
	genesisDeviceCost?: Cost;
	techBoon?: MysteryTraderTechBoonRules[];
};

export type MysteryTraderTechBoonRules = {
	techLevels?: number;
	rewards?: MysteryTraderTechBoonMineralsReward[];
};

export type MysteryTraderTechBoonMineralsReward = {
	mineralsGiven?: number;
	reward?: number;
};

export type WormholeStats = {
	yearsToDegrade: number;
	chanceToJump: number;
	jiggleDistance: number;
};

export type PRTSpec = {
	prt?: PRT;
	pointCost?: number;
	startingTechLevels?: TechLevel;
	startingPlanets?: StartingPlanet[];
	techCostOffset?: TechCostOffset;
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
	builtInScannerMultiplier?: number;
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
	starbaseCostFactor?: number;
	innateMining?: boolean;
	innateResources?: boolean;
	innateScanner?: boolean;
	innatePopulationFactor?: number;
	canBuildDefenses?: boolean;
	livesOnStarbases?: boolean;
};

export type LRTSpec = {
	lrt?: LRT;
	startingFleets?: StartingFleet[];
	pointCost?: number;
	startingTechLevels?: TechLevel;
	techCostOffset?: TechCostOffset;
	newTechCostFactorOffset?: number;
	miniaturizationMax?: number;
	miniaturizationPerLevel?: number;
	noAdvancedScanners?: boolean;
	scanRangeFactorOffset?: number;
	fuelEfficiencyOffset?: number;
	maxPopulationOffset?: number;
	mineralAlchemyCostOffset?: number;
	scrapMineralOffset?: number;
	scrapMineralOffsetStarbase?: number;
	scrapResourcesOffset?: number;
	scrapResourcesOffsetStarbase?: number;
	startingPopulationFactorDelta?: number;
	starbaseBuiltInCloakUnits?: number;
	starbaseCostFactor?: number;
	researchFactorOffset?: number;
	researchSplashDamage?: number;
	shieldStrengthFactorOffset?: number;
	shieldRegenerationRateOffset?: number;
	armorStrengthFactorOffset?: number;
	engineFailureRateOffset?: number;
	engineReliableSpeed?: number;
};

export type TechCostOffset = {
	engine?: number;
	beamWeapon?: number;
	torpedo?: number;
	bomb?: number;
	planetaryDefense?: number;
	stargate?: number;
	terraforming?: number;
};

export type StartingPlanet = {
	population?: number;
	mines?: number;
	factories?: number;
	defenses?: number;
	habPenaltyFactor?: number;
	hasStargate?: boolean;
	hasMassDriver?: boolean;
	starbaseDesignName?: string;
	starbaseHull?: string;
	startingFleets?: StartingFleet[];
	homeworld?: boolean;
};

export type StartingFleet = {
	name?: string;
	hullName?: StartingFleetHull;
	hullSetNumber?: number;
	purpose?: ShipDesignPurpose;
};

export type StealsResearch = {
	energy?: number;
	weapons?: number;
	propulsion?: number;
	construction?: number;
	electronics?: number;
	biotechnology?: number;
};

export const defaultRules: Rules = {
	defenseCost: {
		ironium: 5,
		boranium: 5,
		germanium: 5,
		resources: 15
	},
	factoryCostGermanium: 4,
	mineralAlchemyCost: 100,
	planetaryScannerCost: {
		ironium: 10,
		boranium: 10,
		germanium: 70,
		resources: 100
	},
	starbaseComponentCostReduction: 2,
	terraformCost: {
		resources: 100
	},
	techBaseCost: [
		0, 50, 80, 130, 210, 340, 550, 890, 1440, 2330, 3770, 6100, 9870, 13850, 18040, 22440, 27050,
		31870, 36900, 42140, 47590, 53250, 59120, 65200, 71490, 77990, 84700
	],
	beamRangeDropoff: 0.1,
	numBattleRounds: 16,
	highRadMineralConcentrationBonusThreshold: 90,
	limitMineralConcentration: 30,
	maxExtraWorldDistance: 180,
	maxHab: 99,
	maxMineralConcentration: 200,
	maxStartingMineralConcentration: 121,
	maxStartingMineralSurface: 1000,
	minExtraPlanetMineralConcentration: 30,
	minExtraWorldDistance: 130,
	minHab: 1,
	minHomeworldMineralConcentration: 30,
	minMineralConcentration: 1,
	minStartingMineralConcentration: 1,
	minStartingMineralSurface: 300,
	raceLeftoverPointsPerItem: {
		Defenses: 10,
		Factories: 5,
		MineralConcentrations: 3,
		Mines: 2,
		SurfaceMinerals: 10
	},
	startingYear: 2400,
	wormholeMinPlanetDistance: 30,
	createdAt: '0001-01-01T00:00:00Z',
	updatedAt: '0001-01-01T00:00:00Z',
	cometStatsBySize: {
		Huge: {
			allMinerals: 50,
			allRandomMinerals: 250,
			bonusMinerals: 3000,
			bonusRandomMinerals: 17000,
			bonusMinConcentration: 65,
			bonusRandomConcentration: 65,
			bonusAffectsMinerals: 3,
			minTerraform: 6,
			randomTerraform: 6,
			affectsHabs: 3,
			popKilledPercent: 0.85
		},
		Large: {
			allMinerals: 50,
			allRandomMinerals: 250,
			bonusMinerals: 3000,
			bonusRandomMinerals: 17000,
			bonusMinConcentration: 50,
			bonusRandomConcentration: 50,
			bonusAffectsMinerals: 3,
			minTerraform: 3,
			randomTerraform: 3,
			affectsHabs: 3,
			popKilledPercent: 0.65
		},
		Medium: {
			allMinerals: 50,
			allRandomMinerals: 250,
			bonusMinerals: 3000,
			bonusRandomMinerals: 17000,
			bonusMinConcentration: 50,
			bonusRandomConcentration: 50,
			bonusAffectsMinerals: 2,
			minTerraform: 3,
			randomTerraform: 3,
			affectsHabs: 2,
			popKilledPercent: 0.45
		},
		Small: {
			allMinerals: 50,
			allRandomMinerals: 250,
			bonusMinerals: 3000,
			bonusRandomMinerals: 17000,
			bonusMinConcentration: 50,
			bonusRandomConcentration: 50,
			bonusAffectsMinerals: 1,
			minTerraform: 3,
			randomTerraform: 3,
			affectsHabs: 1,
			popKilledPercent: 0.25
		}
	},
	fleetSafeSpeedExplosionChance: 0.1,
	invasionDefenseCoverageFactor: 0.75,
	lrtSpecs: {
		'1': {
			startingTechLevels: {
				propulsion: 1
			},
			techCostOffset: {},
			fuelEfficiencyOffset: -0.15
		},
		'2': {
			startingTechLevels: {},
			techCostOffset: {
				terraforming: -0.3
			}
		},
		'4': {
			startingFleets: [
				{
					name: 'Potato Bug',
					hullName: 'Midget-Miner',
					purpose: 'Miner'
				},
				{
					name: 'Potato Bug',
					hullName: 'Midget-Miner',
					purpose: 'Miner'
				}
			],
			startingTechLevels: {},
			techCostOffset: {}
		},
		'8': {
			startingTechLevels: {},
			techCostOffset: {},
			starbaseBuiltInCloakUnits: 40,
			starbaseCostFactor: 0.8
		},
		'16': {
			startingTechLevels: {},
			techCostOffset: {},
			researchFactorOffset: -0.5,
			researchSplashDamage: 0.15
		},
		'32': {
			startingTechLevels: {},
			techCostOffset: {},
			scrapMineralOffset: 0.11666666666666667,
			scrapMineralOffsetStarbase: 0.1,
			scrapResourcesOffset: 0.35,
			scrapResourcesOffsetStarbase: 0.7
		},
		'64': {
			startingTechLevels: {},
			techCostOffset: {}
		},
		'128': {
			startingTechLevels: {},
			techCostOffset: {},
			maxPopulationOffset: 0.1
		},
		'256': {
			startingTechLevels: {},
			techCostOffset: {},
			noAdvancedScanners: true,
			scanRangeFactorOffset: 1
		},
		'512': {
			startingTechLevels: {},
			techCostOffset: {},
			startingPopulationFactorDelta: -0.3
		},
		'1024': {
			startingTechLevels: {},
			techCostOffset: {},
			newTechCostFactorOffset: 1,
			miniaturizationMax: 0.05,
			miniaturizationPerLevel: 0.01
		},
		'2048': {
			startingTechLevels: {},
			techCostOffset: {},
			shieldStrengthFactorOffset: 0.4,
			shieldRegenerationRateOffset: 0.1,
			armorStrengthFactorOffset: -0.5
		},
		'4096': {
			startingTechLevels: {},
			techCostOffset: {},
			mineralAlchemyCostOffset: -75
		},
		'8192': {
			startingTechLevels: {
				propulsion: 1
			},
			techCostOffset: {
				engine: -0.5
			},
			engineFailureRateOffset: 0.1,
			engineReliableSpeed: 6
		}
	},
	maxPopulation: 1000000,
	maxTechLevel: 26,
	mineFieldCloak: 75,
	mineFieldStatsByType: {
		Heavy: {
			minDamagePerFleetRS: 2500,
			damagePerEngineRS: 600,
			maxSpeed: 6,
			chanceOfHit: 0.01,
			minDamagePerFleet: 2000,
			damagePerEngine: 500,
			sweepFactor: 1,
			minDecay: 10,
			canDetonate: false
		},
		SpeedBump: {
			minDamagePerFleetRS: 0,
			damagePerEngineRS: 0,
			maxSpeed: 5,
			chanceOfHit: 0.035,
			minDamagePerFleet: 0,
			damagePerEngine: 0,
			sweepFactor: 0.333333343,
			minDecay: 0,
			canDetonate: false
		},
		Standard: {
			minDamagePerFleetRS: 600,
			damagePerEngineRS: 125,
			maxSpeed: 4,
			chanceOfHit: 0.003,
			minDamagePerFleet: 500,
			damagePerEngine: 100,
			sweepFactor: 1,
			minDecay: 10,
			canDetonate: true
		}
	},
	mineralDecayFactor: 1500000,
	minMaxPopulationPercent: 0.05,
	movesToRunAway: 7,
	mysteryTraderRules: {
		chanceSpawn: [7, 7, 7, 7, 7, 7, 7, 4, 4, 3, 2],
		chanceMaxTechGetsPart: 5,
		chanceCourseChange: 20,
		chanceSpeedUpOnly: 3,
		chanceAgain: 2,
		minYear: 40,
		evenYearOnly: true,
		minWarp: 7,
		maxWarp: 13,
		maxMysteryTraders: 5,
		requestedBoon: 5000,
		genesisDeviceCost: {
			resources: 5000
		},
		techBoon: [
			{
				techLevels: 59,
				rewards: [
					{
						mineralsGiven: 5000,
						reward: 6
					},
					{
						mineralsGiven: 6200,
						reward: 7
					},
					{
						mineralsGiven: 7400,
						reward: 8
					},
					{
						mineralsGiven: 8600,
						reward: 9
					},
					{
						mineralsGiven: 9800,
						reward: 10
					}
				]
			},
			{
				techLevels: 71,
				rewards: [
					{
						mineralsGiven: 5000,
						reward: 5
					},
					{
						mineralsGiven: 6200,
						reward: 6
					},
					{
						mineralsGiven: 7400,
						reward: 7
					},
					{
						mineralsGiven: 8600,
						reward: 8
					},
					{
						mineralsGiven: 9800,
						reward: 9
					}
				]
			},
			{
				techLevels: 83,
				rewards: [
					{
						mineralsGiven: 5000,
						reward: 4
					},
					{
						mineralsGiven: 6200,
						reward: 5
					},
					{
						mineralsGiven: 7400,
						reward: 6
					},
					{
						mineralsGiven: 8600,
						reward: 7
					},
					{
						mineralsGiven: 9800,
						reward: 8
					}
				]
			},
			{
				techLevels: 95,
				rewards: [
					{
						mineralsGiven: 5000,
						reward: 3
					},
					{
						mineralsGiven: 6200,
						reward: 4
					},
					{
						mineralsGiven: 7400,
						reward: 5
					},
					{
						mineralsGiven: 8600,
						reward: 6
					},
					{
						mineralsGiven: 9800,
						reward: 7
					}
				]
			},
			{
				techLevels: 107,
				rewards: [
					{
						mineralsGiven: 5000,
						reward: 2
					},
					{
						mineralsGiven: 6200,
						reward: 2
					},
					{
						mineralsGiven: 7400,
						reward: 2
					},
					{
						mineralsGiven: 8600,
						reward: 2
					},
					{
						mineralsGiven: 9800,
						reward: 2
					}
				]
			},
			{
				techLevels: 108,
				rewards: [
					{
						mineralsGiven: 5000,
						reward: 1
					},
					{
						mineralsGiven: 6200,
						reward: 1
					},
					{
						mineralsGiven: 7400,
						reward: 1
					},
					{
						mineralsGiven: 8600,
						reward: 1
					},
					{
						mineralsGiven: 9800,
						reward: 1
					}
				]
			}
		]
	},
	packetDecayRate: {
		'1': 0.1,
		'2': 0.25,
		'3': 0.5
	},
	packetMaxOverwarpSpeed: 3,
	packetMinDecay: 10,
	planetMinDistance: 15,
	populationOvercrowdDieoffRate: 0.04,
	populationOvercrowdDieoffRateMax: 0.12,
	populationScannerError: 0.2,
	prtSpecs: {
		AR: {
			pointCost: 66,
			startingTechLevels: {
				energy: 1
			},
			startingPlanets: [
				{
					population: 25000,
					mines: 10,
					factories: 10,
					defenses: 10,
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
						}
					],
					homeworld: true
				}
			],
			techCostOffset: {},
			mineralsPerSingleMineralPacket: 100,
			mineralsPerMixedMineralPacket: 40,
			packetResourceCost: 10,
			packetMineralCostFactor: 1.1,
			packetReceiverFactor: 1,
			packetDecayFactor: 1,
			packetPermaTerraformSizeUnit: 100,
			shipsVanishInVoid: true,
			techsCostExtraLevel: 3,
			freighterGrowthFactor: -0.03,
			growthFactor: 1,
			stealsResearch: {},
			mineFieldMinDecayFactor: 1,
			mineFieldBaseDecayRate: 0.02,
			mineFieldPlanetDecayRate: 0.04,
			mineFieldMaxDecayRate: 0.5,
			mineFieldDetonateDecayRate: 0.25,
			canRemoteMineOwnPlanets: true,
			invasionAttackBonus: 1.1,
			invasionDefendBonus: 1,
			repairFactor: 1,
			starbaseRepairFactor: 1,
			starbaseCostFactor: 0.8,
			innateMining: true,
			innateResources: true,
			innateScanner: true,
			innatePopulationFactor: 0.1,
			livesOnStarbases: true
		},
		CA: {
			pointCost: 66,
			startingTechLevels: {
				energy: 1,
				weapons: 1,
				propulsion: 1,
				construction: 2,
				biotechnology: 6
			},
			startingPlanets: [
				{
					population: 25000,
					mines: 10,
					factories: 10,
					defenses: 10,
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
							name: 'Change of Heart',
							hullName: 'Mini-Miner',
							hullSetNumber: 1,
							purpose: 'Terraformer'
						}
					],
					homeworld: true
				}
			],
			techCostOffset: {},
			mineralsPerSingleMineralPacket: 100,
			mineralsPerMixedMineralPacket: 40,
			packetResourceCost: 10,
			packetMineralCostFactor: 1.1,
			packetReceiverFactor: 1,
			packetDecayFactor: 1,
			packetPermaTerraformSizeUnit: 100,
			shipsVanishInVoid: true,
			techsCostExtraLevel: 3,
			growthFactor: 1,
			stealsResearch: {},
			mineFieldMinDecayFactor: 1,
			mineFieldBaseDecayRate: 0.02,
			mineFieldPlanetDecayRate: 0.04,
			mineFieldMaxDecayRate: 0.5,
			mineFieldDetonateDecayRate: 0.25,
			invasionAttackBonus: 1.1,
			invasionDefendBonus: 1,
			instaforming: true,
			permaformChance: 0.1,
			permaformPopulation: 100000,
			repairFactor: 1,
			starbaseRepairFactor: 1,
			starbaseCostFactor: 1,
			innatePopulationFactor: 1,
			canBuildDefenses: true
		},
		HE: {
			pointCost: 66,
			startingTechLevels: {},
			startingPlanets: [
				{
					population: 25000,
					mines: 10,
					factories: 10,
					defenses: 10,
					starbaseDesignName: 'Starbase',
					starbaseHull: 'Space Station',
					startingFleets: [
						{
							name: 'Long Range Scout',
							hullName: 'Scout',
							purpose: 'Scout'
						},
						{
							name: 'Spore Cloud',
							hullName: 'Mini-Colony Ship',
							purpose: 'Colonizer'
						},
						{
							name: 'Spore Cloud',
							hullName: 'Mini-Colony Ship',
							purpose: 'Colonizer'
						},
						{
							name: 'Spore Cloud',
							hullName: 'Mini-Colony Ship',
							purpose: 'Colonizer'
						}
					],
					homeworld: true
				}
			],
			techCostOffset: {},
			mineralsPerSingleMineralPacket: 100,
			mineralsPerMixedMineralPacket: 40,
			packetResourceCost: 10,
			packetMineralCostFactor: 1.1,
			packetReceiverFactor: 1,
			packetDecayFactor: 1,
			packetPermaTerraformSizeUnit: 100,
			shipsVanishInVoid: true,
			techsCostExtraLevel: 3,
			growthFactor: 2,
			maxPopulationOffset: -0.5,
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
			starbaseCostFactor: 1,
			innatePopulationFactor: 1,
			canBuildDefenses: true
		},
		IS: {
			pointCost: 66,
			startingTechLevels: {},
			startingPlanets: [
				{
					population: 25000,
					mines: 10,
					factories: 10,
					defenses: 10,
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
						}
					],
					homeworld: true
				}
			],
			techCostOffset: {
				beamWeapon: 0.25,
				torpedo: 0.25,
				bomb: 0.25,
				planetaryDefense: -0.4
			},
			mineralsPerSingleMineralPacket: 100,
			mineralsPerMixedMineralPacket: 40,
			packetResourceCost: 10,
			packetMineralCostFactor: 1.1,
			packetReceiverFactor: 1,
			packetDecayFactor: 1,
			packetPermaTerraformSizeUnit: 100,
			shipsVanishInVoid: true,
			techsCostExtraLevel: 3,
			freighterGrowthFactor: 0.5,
			growthFactor: 1,
			stealsResearch: {},
			mineFieldMinDecayFactor: 1,
			mineFieldBaseDecayRate: 0.02,
			mineFieldPlanetDecayRate: 0.04,
			mineFieldMaxDecayRate: 0.5,
			mineFieldDetonateDecayRate: 0.25,
			invasionAttackBonus: 1.1,
			invasionDefendBonus: 2,
			repairFactor: 2,
			starbaseRepairFactor: 1.5,
			starbaseCostFactor: 1,
			innatePopulationFactor: 1,
			canBuildDefenses: true
		},
		IT: {
			pointCost: 66,
			startingTechLevels: {
				propulsion: 5,
				construction: 5
			},
			startingPlanets: [
				{
					population: 20000,
					mines: 10,
					factories: 10,
					defenses: 10,
					hasStargate: true,
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
							name: 'Swashbuckler',
							hullName: 'Privateer',
							purpose: 'ArmedFreighter'
						},
						{
							name: 'Stalwart Defender',
							hullName: 'Destroyer',
							purpose: 'Fighter'
						}
					],
					homeworld: true
				},
				{
					population: 10000,
					mines: 10,
					factories: 4,
					habPenaltyFactor: 1,
					hasStargate: true,
					starbaseDesignName: 'Accelerator Platform',
					starbaseHull: 'Orbital Fort',
					startingFleets: [
						{
							name: 'Long Range Scout',
							hullName: 'Scout',
							purpose: 'Scout'
						}
					]
				}
			],
			techCostOffset: {
				stargate: -0.25
			},
			mineralsPerSingleMineralPacket: 100,
			mineralsPerMixedMineralPacket: 40,
			packetResourceCost: 10,
			packetMineralCostFactor: 1.2,
			packetReceiverFactor: 0.5,
			packetDecayFactor: 1,
			packetOverSafeWarpPenalty: 1,
			packetPermaTerraformSizeUnit: 100,
			canGateCargo: true,
			canDetectStargatePlanets: true,
			techsCostExtraLevel: 3,
			growthFactor: 1,
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
			starbaseCostFactor: 1,
			innatePopulationFactor: 1,
			canBuildDefenses: true
		},
		JoaT: {
			pointCost: 66,
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
					mines: 10,
					factories: 10,
					defenses: 10,
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
					],
					homeworld: true
				}
			],
			techCostOffset: {},
			mineralsPerSingleMineralPacket: 100,
			mineralsPerMixedMineralPacket: 40,
			packetResourceCost: 10,
			packetMineralCostFactor: 1.1,
			packetReceiverFactor: 1,
			packetDecayFactor: 1,
			packetPermaTerraformSizeUnit: 100,
			shipsVanishInVoid: true,
			builtInScannerMultiplier: 20,
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
			starbaseCostFactor: 1,
			innatePopulationFactor: 1,
			canBuildDefenses: true
		},
		PP: {
			pointCost: 66,
			startingTechLevels: {
				energy: 4
			},
			startingPlanets: [
				{
					population: 20000,
					mines: 10,
					factories: 10,
					defenses: 10,
					hasMassDriver: true,
					starbaseDesignName: 'Starbase',
					starbaseHull: 'Space Station',
					startingFleets: [
						{
							name: 'Long Range Scout',
							hullName: 'Scout',
							purpose: 'Scout'
						},
						{
							name: 'Long Range Scout',
							hullName: 'Scout',
							purpose: 'Scout'
						},
						{
							name: 'Santa Maria',
							hullName: 'Colony Ship',
							purpose: 'Colonizer'
						}
					],
					homeworld: true
				},
				{
					population: 10000,
					mines: 10,
					factories: 4,
					habPenaltyFactor: 1,
					hasMassDriver: true,
					starbaseDesignName: 'Accelerator Platform',
					starbaseHull: 'Orbital Fort',
					startingFleets: [
						{
							name: 'Long Range Scout',
							hullName: 'Scout',
							purpose: 'Scout'
						}
					]
				}
			],
			techCostOffset: {},
			mineralsPerSingleMineralPacket: 70,
			mineralsPerMixedMineralPacket: 25,
			packetResourceCost: 5,
			packetMineralCostFactor: 1,
			packetReceiverFactor: 1,
			packetDecayFactor: 0.5,
			packetBuiltInScanner: true,
			detectPacketDestinationStarbases: true,
			detectAllPackets: true,
			packetTerraformChance: 0.5,
			packetPermaformChance: 0.001,
			packetPermaTerraformSizeUnit: 100,
			shipsVanishInVoid: true,
			techsCostExtraLevel: 3,
			growthFactor: 1,
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
			starbaseCostFactor: 1,
			innatePopulationFactor: 1,
			canBuildDefenses: true
		},
		SD: {
			pointCost: 66,
			startingTechLevels: {
				propulsion: 2,
				biotechnology: 2
			},
			startingPlanets: [
				{
					population: 25000,
					mines: 10,
					factories: 10,
					defenses: 10,
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
							name: 'Little Hen',
							hullName: 'Mini Mine Layer',
							purpose: 'DamageMineLayer'
						},
						{
							name: 'Speed Turtle',
							hullName: 'Mini Mine Layer',
							purpose: 'SpeedMineLayer'
						}
					],
					homeworld: true
				}
			],
			techCostOffset: {},
			mineralsPerSingleMineralPacket: 100,
			mineralsPerMixedMineralPacket: 40,
			packetResourceCost: 10,
			packetMineralCostFactor: 1.1,
			packetReceiverFactor: 1,
			packetDecayFactor: 1,
			packetPermaTerraformSizeUnit: 100,
			shipsVanishInVoid: true,
			techsCostExtraLevel: 3,
			growthFactor: 1,
			stealsResearch: {},
			mineFieldsAreScanners: true,
			mineFieldRateMoveFactor: 0.5,
			mineFieldSafeWarpBonus: 2,
			mineFieldMinDecayFactor: 0.25,
			mineFieldBaseDecayRate: 0.02,
			mineFieldPlanetDecayRate: 0.04,
			mineFieldMaxDecayRate: 0.5,
			canDetonateMineFields: true,
			mineFieldDetonateDecayRate: 0.25,
			invasionAttackBonus: 1.1,
			invasionDefendBonus: 1,
			repairFactor: 1,
			starbaseRepairFactor: 1,
			starbaseCostFactor: 1,
			innatePopulationFactor: 1,
			canBuildDefenses: true
		},
		SS: {
			pointCost: 66,
			startingTechLevels: {
				electronics: 5
			},
			startingPlanets: [
				{
					population: 25000,
					mines: 10,
					factories: 10,
					defenses: 10,
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
						}
					],
					homeworld: true
				}
			],
			techCostOffset: {},
			mineralsPerSingleMineralPacket: 100,
			mineralsPerMixedMineralPacket: 40,
			packetResourceCost: 10,
			packetMineralCostFactor: 1.1,
			packetReceiverFactor: 1,
			packetDecayFactor: 1,
			packetPermaTerraformSizeUnit: 100,
			shipsVanishInVoid: true,
			techsCostExtraLevel: 3,
			growthFactor: 1,
			builtInCloakUnits: 300,
			stealsResearch: {
				energy: 0.5,
				weapons: 0.5,
				propulsion: 0.5,
				construction: 0.5,
				electronics: 0.5,
				biotechnology: 0.5
			},
			freeCargoCloaking: true,
			mineFieldSafeWarpBonus: 1,
			mineFieldMinDecayFactor: 1,
			mineFieldBaseDecayRate: 0.02,
			mineFieldPlanetDecayRate: 0.04,
			mineFieldMaxDecayRate: 0.5,
			mineFieldDetonateDecayRate: 0.25,
			invasionAttackBonus: 1.1,
			invasionDefendBonus: 1,
			repairFactor: 1,
			starbaseRepairFactor: 1,
			starbaseCostFactor: 1,
			innatePopulationFactor: 1,
			canBuildDefenses: true
		},
		WM: {
			pointCost: 66,
			startingTechLevels: {
				energy: 1,
				weapons: 6,
				propulsion: 1
			},
			startingPlanets: [
				{
					population: 25000,
					mines: 10,
					factories: 10,
					defenses: 10,
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
							name: 'Armed Probe',
							hullName: 'Scout',
							hullSetNumber: 1,
							purpose: 'FighterScout'
						}
					],
					homeworld: true
				}
			],
			techCostOffset: {
				beamWeapon: -0.25,
				torpedo: -0.25,
				bomb: -0.25
			},
			mineralsPerSingleMineralPacket: 100,
			mineralsPerMixedMineralPacket: 40,
			packetResourceCost: 10,
			packetMineralCostFactor: 1.1,
			packetReceiverFactor: 1,
			packetDecayFactor: 1,
			packetPermaTerraformSizeUnit: 100,
			shipsVanishInVoid: true,
			techsCostExtraLevel: 3,
			growthFactor: 1,
			stealsResearch: {},
			mineFieldMinDecayFactor: 1,
			mineFieldBaseDecayRate: 0.02,
			mineFieldPlanetDecayRate: 0.04,
			mineFieldMaxDecayRate: 0.5,
			mineFieldDetonateDecayRate: 0.25,
			discoverDesignOnScan: true,
			invasionAttackBonus: 1.65,
			invasionDefendBonus: 1,
			movementBonus: 2,
			repairFactor: 1,
			starbaseRepairFactor: 1,
			starbaseCostFactor: 1,
			innatePopulationFactor: 1,
			canBuildDefenses: true
		}
	},
	raceStartingPoints: 1650,
	radiatingImmune: 85,
	randomArtifactResearchBonusRange: [120, 400],
	randomCometMinYear: 10,
	randomCometMinYearPlayerWorld: 20,
	randomEventChances: {
		AncientArtifact: 0.33,
		Comet: 0.05,
		MineralDeposit: 0.05,
		PlanetaryChange: 0.05
	},
	randomMineralDepositBonusRange: [20, 50],
	remoteMiningMineOutput: 10,
	repairRates: {
		Moving: 0.01,
		None: 0,
		Orbiting: 0.03,
		OrbitingOwnPlanet: 0.05,
		Starbase: 0.1,
		Stopped: 0.02
	},
	salvageDecayMin: 10,
	salvageDecayRate: 0.1,
	salvageFromBattleFactor: 0.3,
	scrapMineralAmount: 0.333333343,
	showPublicScoresAfterYears: 20,
	smartDefenseCoverageFactor: 0.5,
	stargateMaxHullMassFactor: 5,
	stargateMaxRangeFactor: 5,
	tachyonCloakReduction: 5,
	tachyonMaxCloakReduction: 81,
	techTradeChance: 0.5,
	torpedoSplashDamage: 0.125,
	wormholeCloak: 75,
	wormholePairsForSize: {
		Huge: 6,
		HugeWide: 6,
		Large: 5,
		LargeWide: 5,
		Medium: 4,
		MediumWide: 4,
		Small: 3,
		SmallWide: 3,
		Tiny: 1,
		TinyWide: 1
	},
	wormholeStatsByStability: {
		Average: {
			yearsToDegrade: 5,
			chanceToJump: 0.04,
			jiggleDistance: 10
		},
		ExtremelyVolatile: {
			yearsToDegrade: -1,
			chanceToJump: 0.04,
			jiggleDistance: 10
		},
		MostlyStable: {
			yearsToDegrade: 5,
			chanceToJump: 0.02,
			jiggleDistance: 10
		},
		RockSolid: {
			yearsToDegrade: 10,
			chanceToJump: 0,
			jiggleDistance: 10
		},
		SlightlyVolatile: {
			yearsToDegrade: 5,
			chanceToJump: 0.03,
			jiggleDistance: 10
		},
		Stable: {
			yearsToDegrade: 5,
			chanceToJump: 0.005,
			jiggleDistance: 10
		},
		Volatile: {
			yearsToDegrade: 5,
			chanceToJump: 0.06,
			jiggleDistance: 10
		}
	}
};
