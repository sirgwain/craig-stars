import type { Cost } from './Cost';
import type { EnumDictionary } from './EnumDictionary';
import type { Size } from './Game';
import type { MineFieldStats, MineFieldType } from './MineField';
import type { TechStore } from './Tech';

export type Rules = {
	id?: number;
	createdAt?: string;
	updatedAt?: string;

	gameId?: number;
	tachyonCloakReduction: number;
	maxPopulation: number;
	minMaxPopulationPercent: number;
	populationOvercrowdDieoffRate: number;
	populationOvercrowdDieoffRateMax: number;
	populationScannerError: number;
	smartDefenseCoverageFactor: number;
	invasionDefenseCoverageFactor: number;
	numBattleRounds: number;
	movesToRunAway: number;
	beamRangeDropoff: number;
	torpedoSplashDamage: number;
	salvageDecayRate: number;
	salvageDecayMin: number;
	mineFieldCloak: number;
	stargateMaxRangeFactor: number;
	stargateMaxHullMassFactor: number;
	fleetSafeSpeedExplosionChance: number;
	randomEventChances: EnumDictionary<RandomEventType, number>;
	randomMineralDepositBonusRange: number[];
	randomArtifactResearchBonusRange: number[];
	randomCometMinYear: number;
	randomCometMinYearPlayerWorld: number;
	wormholeCloak: number;
	wormholeMinDistance: number;
	wormholeStatsByStability: EnumDictionary<WormholeStability, WormholeStats>;
	wormholePairsForSize: EnumDictionary<Size, number>;
	mineFieldStatsByType: EnumDictionary<MineFieldType, MineFieldStats>;
	repairRates: EnumDictionary<RepairRate, number>;
	maxPlayers: number;
	startingYear: number;
	showPublicScoresAfterYears: number;
	planetMinDistance: number;
	maxExtraWorldDistance: number;
	minExtraWorldDistance: number;
	minHomeworldMineralConcentration: number;
	minExtraPlanetMineralConcentration: number;
	minMineralConcentration: number;
	minStartingMineralConcentration: number;
	maxStartingMineralConcentration: number;
	highRadGermaniumBonus: number;
	highRadGermaniumBonusThreshold: number;
	maxStartingMineralSurface: number;
	minStartingMineralSurface: number;
	minHab: number;
	maxHab: number;
	maxMineralConcentration: number;
	radiatingImmune: number;
	mineralDecayFactor: number;
	remoteMiningMineOutput: number;
	startingMines: number;
	startingFactories: number;
	startingDefenses: number;
	raceStartingPoints: number;
	scrapMineralAmount: number;
	scrapResourceAmount: number;
	factoryCostGermanium: number;
	defenseCost: Cost;
	mineralAlchemyCost: number;
	terraformCost: Cost;
	starbaseComponentCostFactor: number;
	salvageFromBattleFactor: number;
	techTradeChance: number;
	packetDecayRate: { [key: number]: number };
	maxTechLevel: number;
	techBaseCost: number[];
	techs?: TechStore;
	planetaryScannerCost: Cost;
	cometStatsBySize: any;
	prtSpecs: any;
	lrtSpecs: any;
};

export enum WormholeStability {
	RockSolid = 'RockSolid',
	Stable = 'Stable',
	MostlyStable = 'MostlyStable',
	Average = 'Average',
	SlightlyVolatile = 'SlightlyVolatile',
	Volatile = 'Volatile',
	ExtremelyVolatile = 'ExtremelyVolatile'
}

export interface WormholeStats {
	yearsToDegrade: number;
	chanceToJump: number;
	jiggleDistance: number;
}

export enum RandomEventType {
	Comet = 'Comet',
	MineralDeposit = 'MineralDeposit',
	PlanetaryChange = 'PlanetaryChange',
	AncientArtifact = 'AncientArtifact',
	MysteryTrader = 'MysteryTrader'
}

export enum RandomCometSize {
	Small = 'Small',
	Medium = 'Medium',
	Large = 'Large',
	Huge = 'Huge'
}

export enum RepairRate {
	None = 'None',
	Moving = 'Moving',
	Stopped = 'Stopped',
	Orbiting = 'Orbiting',
	OrbitingOwnPlanet = 'OrbitingOwnPlanet',
	Starbase = 'Starbase' // the rate starbases repair, not the rate a fleet repairs at a starbase, that's handled by the TechHull
}

export const defaultRules: Rules = {
	tachyonCloakReduction: 5,
	maxPopulation: 1000000,
	minMaxPopulationPercent: 0.05,
	populationOvercrowdDieoffRate: 0.04,
	populationOvercrowdDieoffRateMax: 0.12,
	populationScannerError: 0.2,
	smartDefenseCoverageFactor: 0.5,
	invasionDefenseCoverageFactor: 0.75,
	numBattleRounds: 16,
	movesToRunAway: 7,
	beamRangeDropoff: 0.1,
	torpedoSplashDamage: 0.125,
	salvageDecayRate: 0.1,
	salvageDecayMin: 10,
	mineFieldCloak: 75,
	stargateMaxRangeFactor: 5,
	stargateMaxHullMassFactor: 5,
	fleetSafeSpeedExplosionChance: 0.1,
	randomEventChances: {
		AncientArtifact: 0.33,
		Comet: 0.05,
		MineralDeposit: 0.05,
		MysteryTrader: 0.05,
		PlanetaryChange: 0.05
	},
	randomMineralDepositBonusRange: [20, 50],
	randomArtifactResearchBonusRange: [120, 400],
	randomCometMinYear: 10,
	randomCometMinYearPlayerWorld: 20,
	cometStatsBySize: {
		Huge: {
			minMinerals: 50,
			randomMinerals: 250,
			bonusMinerals: 3000,
			bonusRandomMinerals: 17000,
			minConcentrationBonus: 65,
			randomConcentrationBonus: 65,
			affectsMinerals: 3,
			minTerraform: 6,
			randomTerraform: 6,
			affectsHabs: 3,
			popKilledPercent: 0.85
		},
		Large: {
			minMinerals: 50,
			randomMinerals: 250,
			bonusMinerals: 3000,
			bonusRandomMinerals: 17000,
			minConcentrationBonus: 50,
			randomConcentrationBonus: 50,
			affectsMinerals: 3,
			minTerraform: 3,
			randomTerraform: 3,
			affectsHabs: 3,
			popKilledPercent: 0.65
		},
		Medium: {
			minMinerals: 50,
			randomMinerals: 250,
			bonusMinerals: 3000,
			bonusRandomMinerals: 17000,
			minConcentrationBonus: 50,
			randomConcentrationBonus: 50,
			affectsMinerals: 2,
			minTerraform: 3,
			randomTerraform: 3,
			affectsHabs: 2,
			popKilledPercent: 0.45
		},
		Small: {
			minMinerals: 50,
			randomMinerals: 250,
			bonusMinerals: 3000,
			bonusRandomMinerals: 17000,
			minConcentrationBonus: 50,
			randomConcentrationBonus: 50,
			affectsMinerals: 1,
			minTerraform: 3,
			randomTerraform: 3,
			affectsHabs: 1,
			popKilledPercent: 0.25
		}
	},
	wormholeCloak: 75,
	wormholeMinDistance: 30,
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
	},
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
	repairRates: {
		Moving: 0.01,
		None: 0,
		Orbiting: 0.03,
		OrbitingOwnPlanet: 0.05,
		Starbase: 0.1,
		Stopped: 0.02
	},
	maxPlayers: 16,
	startingYear: 2400,
	showPublicScoresAfterYears: 20,
	planetMinDistance: 15,
	maxExtraWorldDistance: 180,
	minExtraWorldDistance: 130,
	minHomeworldMineralConcentration: 30,
	minExtraPlanetMineralConcentration: 30,
	minHab: 1,
	maxHab: 99,
	minMineralConcentration: 1,
	maxMineralConcentration: 200,
	minStartingMineralConcentration: 1,
	maxStartingMineralConcentration: 100,
	highRadGermaniumBonus: 5,
	highRadGermaniumBonusThreshold: 85,
	radiatingImmune: 85,
	maxStartingMineralSurface: 1000,
	minStartingMineralSurface: 300,
	mineralDecayFactor: 1500000,
	remoteMiningMineOutput: 10,
	startingMines: 10,
	startingFactories: 10,
	startingDefenses: 10,
	raceStartingPoints: 1650,
	scrapMineralAmount: 0.333333343,
	scrapResourceAmount: 0,
	factoryCostGermanium: 4,
	defenseCost: {
		ironium: 5,
		boranium: 5,
		germanium: 5,
		resources: 15
	},
	mineralAlchemyCost: 100,
	planetaryScannerCost: {
		ironium: 10,
		boranium: 10,
		germanium: 70,
		resources: 100
	},
	terraformCost: {
		resources: 100
	},
	starbaseComponentCostFactor: 0.5,
	salvageFromBattleFactor: 0.3,
	techTradeChance: 0,
	packetDecayRate: {
		'1': 0.1,
		'2': 0.25,
		'3': 0.5
	},
	maxTechLevel: 26,
	techBaseCost: [
		0, 50, 80, 130, 210, 340, 550, 890, 1440, 2330, 3770, 6100, 9870, 13850, 18040, 22440, 27050,
		31870, 36900, 42140, 47590, 53250, 59120, 65200, 71490, 77990, 84700
	],
	prtSpecs: {
		AR: {
			pointCost: 66,
			startingTechLevels: {
				energy: 1
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
			techsCostExtraLevel: 3,
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
					starbaseDesignName: 'Starbase',
					starbaseHull: 'Space Station',
					startingFleets: [
						{
							name: 'Deep Space Probe',
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
					]
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
			packetMineralCostFactor: 1,
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
					population: 25000,
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
					]
				},
				{
					population: 10000,
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
			techCostOffset: {},
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
					population: 25000,
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
					]
				},
				{
					population: 10000,
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
							name: 'Armored Probe',
							hullName: 'Scout',
							hullSetNumber: 1,
							purpose: 'FighterScout'
						}
					]
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
			packetMineralCostFactor: 1,
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
	lrtSpecs: {
		'1': {
			startingTechLevels: {
				propulsion: 1
			},
			techCostOffset: {},
			fuelEfficiencyOffset: -0.15,
			terraformCostOffset: {}
		},
		'2': {
			startingTechLevels: {},
			techCostOffset: {},
			terraformCostOffset: {
				resources: -30
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
			techCostOffset: {},
			terraformCostOffset: {}
		},
		'8': {
			startingTechLevels: {},
			techCostOffset: {},
			terraformCostOffset: {},
			starbaseBuiltInCloakUnits: 40,
			starbaseCostFactorOffset: -0.2
		},
		'16': {
			startingTechLevels: {},
			techCostOffset: {},
			terraformCostOffset: {},
			researchFactorOffset: -0.5,
			researchSplashDamage: 0.15
		},
		'32': {
			startingTechLevels: {},
			techCostOffset: {},
			terraformCostOffset: {},
			scrapMineralOffset: 0.11666666666666667,
			scrapMineralOffsetStarbase: 0.5666666666666667,
			scrapResourcesOffset: 0.35,
			scrapResourcesOffsetStarbase: 0.7
		},
		'64': {
			startingTechLevels: {},
			techCostOffset: {},
			terraformCostOffset: {}
		},
		'128': {
			startingTechLevels: {},
			techCostOffset: {},
			maxPopulationOffset: 0.1,
			terraformCostOffset: {}
		},
		'256': {
			startingTechLevels: {},
			techCostOffset: {},
			noAdvancedScanners: true,
			scanRangeFactorOffset: 1,
			terraformCostOffset: {}
		},
		'512': {
			startingTechLevels: {},
			techCostOffset: {},
			terraformCostOffset: {},
			startingPopulationFactorDelta: -0.3
		},
		'1024': {
			startingTechLevels: {},
			techCostOffset: {},
			newTechCostFactorOffset: 1,
			miniaturizationMax: 0.05,
			miniaturizationPerLevel: 0.01,
			terraformCostOffset: {}
		},
		'2048': {
			startingTechLevels: {},
			techCostOffset: {},
			terraformCostOffset: {},
			shieldStrengthFactorOffset: 0.4,
			shieldRegenerationRateOffset: 0.1,
			armorStrengthFactorOffset: -0.5
		},
		'4096': {
			startingTechLevels: {},
			techCostOffset: {},
			terraformCostOffset: {},
			mineralAlchemyCostOffset: -75
		},
		'8192': {
			startingTechLevels: {
				propulsion: 1
			},
			techCostOffset: {
				engine: -0.5
			},
			terraformCostOffset: {},
			engineFailureRateOffset: 0.1,
			engineReliableSpeed: 6
		}
	}
};
