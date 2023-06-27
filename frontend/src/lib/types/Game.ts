import type { Cost } from './Cost';
import type { EnumDictionary } from './EnumDictionary';
import type { MineFieldStats, MineFieldType } from './MineField';
import type { PlayerStatus } from './Player';
import type { TechStore } from './Tech';
import type { Vector } from './Vector';

export enum Size {
	Tiny = 'Tiny',
	TinyWide = 'TinyWide',
	Small = 'Small',
	SmallWide = 'SmallWide',
	Medium = 'Medium',
	MediumWide = 'MediumWide',
	Large = 'Large',
	LargeWide = 'LargeWide',
	Huge = 'Huge',
	HugeWide = 'HugeWide'
}

export enum Density {
	Sparse = 'Sparse',
	Normal = 'Normal',
	Dense = 'Dense',
	Packed = 'Packed'
}

export enum PlayerPositions {
	Close = 'Close',
	Moderate = 'Moderate',
	Farther = 'Farther',
	Distant = 'Distant'
}

export enum GameStartMode {
	Normal = '', // regular start
	MidGame = 'MidGame', // further tech levels, pop growth
	LateGame = 'LateGame',
	EndGame = 'EndGame'
}

export enum NewGamePlayerType {
	Host = 'Host',
	Open = 'Open',
	AI = 'AI'
}

export type NewGamePlayer = {
	type: NewGamePlayerType;
	userId?: number;
	raceId?: number;
	color?: string;
};

export type NewGamePlayers = {
	players: NewGamePlayer[];
};

export type GameSettings = {
	name: string;
	public: boolean;
	quickStartTurns?: number;
	size: Size;
	area?: Vector;
	density: Density;
	playerPositions: PlayerPositions;
	randomEvents?: boolean;
	computerPlayersFormAlliances?: boolean;
	publicPlayerScores?: boolean;
	startMode?: GameStartMode;
	year?: number;
	victoryConditions?: VictoryConditions;
};

export enum GameState {
	Setup = 'Setup',
	WaitingForPlayers = 'WaitingForPlayers',
	GeneratingUniverse = 'GeneratingUniverse',
	GeneratingTurn = 'GeneratingTurn',
	GeneratingTurnError = 'GeneratingTurnError'
}

export interface Game {
	id: number;
	createdAt: string;
	updatedAt: string;
	hostId: number;

	name: string;
	hash?: string;
	state: GameState;
	numPlayers: number;
	openPlayerSlots: number;
	quickStartTurns: number;
	size: Size;
	area: Vector;
	density: Density;
	playerPositions: PlayerPositions;
	randomEvents: boolean;
	computerPlayersFormAlliances: boolean;
	publicPlayerScores: boolean;
	public?: boolean;
	startMode: GameStartMode;
	year: number;
	victoryConditions: VictoryConditions;
	victorDeclared: boolean;
	rules: Rules;
	players: PlayerStatus[];
}

export interface Rules {
	id?: number;
	createdAt?: string;
	updatedAt?: string;

	gameId?: number;
	tachyonCloakReduction: number;
	maxPopulation: number;
	fleetsScanWhileMoving: boolean;
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
	randomEventChances: EnumDictionary<RandomEventType, number>;
	randomMineralDepositBonusRange: number[];
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
	packetDecayRate: { [key: number]: number };
	maxTechLevel: number;
	techBaseCost: number[];
	techs?: TechStore;
	prtSpecs: any;
	lrtSpecs: any;
}

export interface WormholeStats {
	yearsToDegrade: number;
	chanceToJump: number;
	jiggleDistance: number;
}

export type VictoryConditions = {
	conditions: number;
	numCriteriaRequired: number;
	yearsPassed: number;
	ownPlanets: number;
	attainTechLevel: number;
	attainTechLevelNumFields: number;
	exceedsScore: number;
	exceedsSecondPlaceScore: number;
	productionCapacity: number;
	ownCapitalShips: number;
	highestScoreAfterYears: number;
};

export enum VictoryCondition {
	None = 0,
	OwnPlanets = 1 << 0,
	AttainTechLevels = 1 << 1,
	ExceedsScore = 1 << 2,
	ExceedsSecondPlaceScore = 1 << 3,
	ProductionCapacity = 1 << 4,
	OwnCapitalShips = 1 << 5,
	HighestScoreAfterYears = 1 << 6
}

export enum WormholeStability {
	RockSolid = 'RockSolid',
	Stable = 'Stable',
	MostlyStable = 'MostlyStable',
	Average = 'Average',
	SlightlyVolatile = 'SlightlyVolatile',
	Volatile = 'Volatile',
	ExtremelyVolatile = 'ExtremelyVolatile'
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
	fleetsScanWhileMoving: false,
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
	randomEventChances: {
		AncientArtifact: 0.01,
		Comet: 0.01,
		MineralDeposit: 0.01,
		MysteryTrader: 0.01,
		PlanetaryChange: 0.01
	},
	randomMineralDepositBonusRange: [20, 50],
	wormholeCloak: 75,
	wormholeMinDistance: 30,
	wormholeStatsByStability: {
		Average: {
			yearsToDegrade: 0,
			chanceToJump: 0,
			jiggleDistance: 10
		},
		ExtremelyVolatile: {
			yearsToDegrade: 0,
			chanceToJump: 0,
			jiggleDistance: 10
		},
		MostlyStable: {
			yearsToDegrade: 0,
			chanceToJump: 0,
			jiggleDistance: 10
		},
		RockSolid: {
			yearsToDegrade: 0,
			chanceToJump: 0,
			jiggleDistance: 10
		},
		SlightlyVolatile: {
			yearsToDegrade: 0,
			chanceToJump: 0,
			jiggleDistance: 10
		},
		Stable: {
			yearsToDegrade: 0,
			chanceToJump: 0,
			jiggleDistance: 10
		},
		Volatile: {
			yearsToDegrade: 0,
			chanceToJump: 0,
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
	showPublicScoresAfterYears: 1,
	planetMinDistance: 15,
	maxExtraWorldDistance: 180,
	minExtraWorldDistance: 130,
	minHomeworldMineralConcentration: 30,
	minExtraPlanetMineralConcentration: 30,
	minMineralConcentration: 1,
	minStartingMineralConcentration: 1,
	maxStartingMineralConcentration: 100,
	highRadGermaniumBonus: 5,
	highRadGermaniumBonusThreshold: 85,
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
	terraformCost: {
		resources: 100
	},
	starbaseComponentCostFactor: 0.5,
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
			prt: '',
			pointCost: 66,
			startingTechLevels: {},
			startingFleets: null,
			startingPlanets: [
				{
					population: 25000,
					habPenaltyFactor: 0,
					hasStargate: false,
					hasMassDriver: false,
					startingFleets: null
				}
			],
			techCostOffset: {},
			mineralsPerSingleMineralPacket: 100,
			mineralsPerMixedMineralPacket: 40,
			packetResourceCost: 10,
			packetMineralCostFactor: 1,
			packetReceiverFactor: 1,
			packetDecayFactor: 1,
			packetOverSafeWarpPenalty: 0,
			packetBuiltInScanner: false,
			detectPacketDestinationStarbases: false,
			detectAllPackets: false,
			packetTerraformChance: 0,
			packetPermaformChance: 0,
			packetPermaTerraformSizeUnit: 100,
			canGateCargo: false,
			canDetectStargatePlanets: false,
			shipsVanishInVoid: true,
			builtInScannerMultiplier: 0,
			techsCostExtraLevel: 3,
			freighterGrowthFactor: 0,
			growthFactor: 1,
			maxPopulationOffset: 0,
			builtInCloakUnits: 0,
			stealsResearch: {},
			freeCargoCloaking: false,
			mineFieldsAreScanners: false,
			mineFieldRateMoveFactor: 0,
			mineFieldSafeWarpBonus: 0,
			mineFieldMinDecayFactor: 1,
			mineFieldBaseDecayRate: 0.02,
			mineFieldPlanetDecayRate: 0.04,
			mineFieldMaxDecayRate: 0.5,
			canDetonateMineFields: false,
			mineFieldDetonateDecayRate: 0.25,
			discoverDesignOnScan: false,
			canRemoteMineOwnPlanets: false,
			invasionAttackBonus: 1,
			invasionDefendBonus: 1,
			movementBonus: 0,
			instaforming: false,
			permaformChance: 0,
			permaformPopulation: 0,
			repairFactor: 1,
			starbaseRepairFactor: 1,
			starbaseCostFactor: 1,
			innateMining: false,
			innateResources: false,
			innateScanner: false,
			innatePopulationFactor: 1,
			canBuildDefenses: true,
			livesOnStarbases: false
		},
		CA: {
			prt: '',
			pointCost: 66,
			startingTechLevels: {},
			startingFleets: null,
			startingPlanets: [
				{
					population: 25000,
					habPenaltyFactor: 0,
					hasStargate: false,
					hasMassDriver: false,
					startingFleets: null
				}
			],
			techCostOffset: {},
			mineralsPerSingleMineralPacket: 100,
			mineralsPerMixedMineralPacket: 40,
			packetResourceCost: 10,
			packetMineralCostFactor: 1,
			packetReceiverFactor: 1,
			packetDecayFactor: 1,
			packetOverSafeWarpPenalty: 0,
			packetBuiltInScanner: false,
			detectPacketDestinationStarbases: false,
			detectAllPackets: false,
			packetTerraformChance: 0,
			packetPermaformChance: 0,
			packetPermaTerraformSizeUnit: 100,
			canGateCargo: false,
			canDetectStargatePlanets: false,
			shipsVanishInVoid: true,
			builtInScannerMultiplier: 0,
			techsCostExtraLevel: 3,
			freighterGrowthFactor: 0,
			growthFactor: 1,
			maxPopulationOffset: 0,
			builtInCloakUnits: 0,
			stealsResearch: {},
			freeCargoCloaking: false,
			mineFieldsAreScanners: false,
			mineFieldRateMoveFactor: 0,
			mineFieldSafeWarpBonus: 0,
			mineFieldMinDecayFactor: 1,
			mineFieldBaseDecayRate: 0.02,
			mineFieldPlanetDecayRate: 0.04,
			mineFieldMaxDecayRate: 0.5,
			canDetonateMineFields: false,
			mineFieldDetonateDecayRate: 0.25,
			discoverDesignOnScan: false,
			canRemoteMineOwnPlanets: false,
			invasionAttackBonus: 1,
			invasionDefendBonus: 1,
			movementBonus: 0,
			instaforming: false,
			permaformChance: 0,
			permaformPopulation: 0,
			repairFactor: 1,
			starbaseRepairFactor: 1,
			starbaseCostFactor: 1,
			innateMining: false,
			innateResources: false,
			innateScanner: false,
			innatePopulationFactor: 1,
			canBuildDefenses: true,
			livesOnStarbases: false
		},
		HE: {
			prt: '',
			pointCost: 66,
			startingTechLevels: {},
			startingFleets: null,
			startingPlanets: [
				{
					population: 25000,
					habPenaltyFactor: 0,
					hasStargate: false,
					hasMassDriver: false,
					startingFleets: null
				}
			],
			techCostOffset: {},
			mineralsPerSingleMineralPacket: 100,
			mineralsPerMixedMineralPacket: 40,
			packetResourceCost: 10,
			packetMineralCostFactor: 1,
			packetReceiverFactor: 1,
			packetDecayFactor: 1,
			packetOverSafeWarpPenalty: 0,
			packetBuiltInScanner: false,
			detectPacketDestinationStarbases: false,
			detectAllPackets: false,
			packetTerraformChance: 0,
			packetPermaformChance: 0,
			packetPermaTerraformSizeUnit: 100,
			canGateCargo: false,
			canDetectStargatePlanets: false,
			shipsVanishInVoid: true,
			builtInScannerMultiplier: 0,
			techsCostExtraLevel: 3,
			freighterGrowthFactor: 0,
			growthFactor: 1,
			maxPopulationOffset: 0,
			builtInCloakUnits: 0,
			stealsResearch: {},
			freeCargoCloaking: false,
			mineFieldsAreScanners: false,
			mineFieldRateMoveFactor: 0,
			mineFieldSafeWarpBonus: 0,
			mineFieldMinDecayFactor: 1,
			mineFieldBaseDecayRate: 0.02,
			mineFieldPlanetDecayRate: 0.04,
			mineFieldMaxDecayRate: 0.5,
			canDetonateMineFields: false,
			mineFieldDetonateDecayRate: 0.25,
			discoverDesignOnScan: false,
			canRemoteMineOwnPlanets: false,
			invasionAttackBonus: 1,
			invasionDefendBonus: 1,
			movementBonus: 0,
			instaforming: false,
			permaformChance: 0,
			permaformPopulation: 0,
			repairFactor: 1,
			starbaseRepairFactor: 1,
			starbaseCostFactor: 1,
			innateMining: false,
			innateResources: false,
			innateScanner: false,
			innatePopulationFactor: 1,
			canBuildDefenses: true,
			livesOnStarbases: false
		},
		IS: {
			prt: '',
			pointCost: 66,
			startingTechLevels: {},
			startingFleets: null,
			startingPlanets: [
				{
					population: 25000,
					habPenaltyFactor: 0,
					hasStargate: false,
					hasMassDriver: false,
					startingFleets: null
				}
			],
			techCostOffset: {},
			mineralsPerSingleMineralPacket: 100,
			mineralsPerMixedMineralPacket: 40,
			packetResourceCost: 10,
			packetMineralCostFactor: 1,
			packetReceiverFactor: 1,
			packetDecayFactor: 1,
			packetOverSafeWarpPenalty: 0,
			packetBuiltInScanner: false,
			detectPacketDestinationStarbases: false,
			detectAllPackets: false,
			packetTerraformChance: 0,
			packetPermaformChance: 0,
			packetPermaTerraformSizeUnit: 100,
			canGateCargo: false,
			canDetectStargatePlanets: false,
			shipsVanishInVoid: true,
			builtInScannerMultiplier: 0,
			techsCostExtraLevel: 3,
			freighterGrowthFactor: 0,
			growthFactor: 1,
			maxPopulationOffset: 0,
			builtInCloakUnits: 0,
			stealsResearch: {},
			freeCargoCloaking: false,
			mineFieldsAreScanners: false,
			mineFieldRateMoveFactor: 0,
			mineFieldSafeWarpBonus: 0,
			mineFieldMinDecayFactor: 1,
			mineFieldBaseDecayRate: 0.02,
			mineFieldPlanetDecayRate: 0.04,
			mineFieldMaxDecayRate: 0.5,
			canDetonateMineFields: false,
			mineFieldDetonateDecayRate: 0.25,
			discoverDesignOnScan: false,
			canRemoteMineOwnPlanets: false,
			invasionAttackBonus: 1,
			invasionDefendBonus: 1,
			movementBonus: 0,
			instaforming: false,
			permaformChance: 0,
			permaformPopulation: 0,
			repairFactor: 1,
			starbaseRepairFactor: 1,
			starbaseCostFactor: 1,
			innateMining: false,
			innateResources: false,
			innateScanner: false,
			innatePopulationFactor: 1,
			canBuildDefenses: true,
			livesOnStarbases: false
		},
		IT: {
			prt: '',
			pointCost: 66,
			startingTechLevels: {},
			startingFleets: null,
			startingPlanets: [
				{
					population: 25000,
					habPenaltyFactor: 0,
					hasStargate: false,
					hasMassDriver: false,
					startingFleets: null
				}
			],
			techCostOffset: {},
			mineralsPerSingleMineralPacket: 100,
			mineralsPerMixedMineralPacket: 40,
			packetResourceCost: 10,
			packetMineralCostFactor: 1,
			packetReceiverFactor: 1,
			packetDecayFactor: 1,
			packetOverSafeWarpPenalty: 0,
			packetBuiltInScanner: false,
			detectPacketDestinationStarbases: false,
			detectAllPackets: false,
			packetTerraformChance: 0,
			packetPermaformChance: 0,
			packetPermaTerraformSizeUnit: 100,
			canGateCargo: false,
			canDetectStargatePlanets: false,
			shipsVanishInVoid: true,
			builtInScannerMultiplier: 0,
			techsCostExtraLevel: 3,
			freighterGrowthFactor: 0,
			growthFactor: 1,
			maxPopulationOffset: 0,
			builtInCloakUnits: 0,
			stealsResearch: {},
			freeCargoCloaking: false,
			mineFieldsAreScanners: false,
			mineFieldRateMoveFactor: 0,
			mineFieldSafeWarpBonus: 0,
			mineFieldMinDecayFactor: 1,
			mineFieldBaseDecayRate: 0.02,
			mineFieldPlanetDecayRate: 0.04,
			mineFieldMaxDecayRate: 0.5,
			canDetonateMineFields: false,
			mineFieldDetonateDecayRate: 0.25,
			discoverDesignOnScan: false,
			canRemoteMineOwnPlanets: false,
			invasionAttackBonus: 1,
			invasionDefendBonus: 1,
			movementBonus: 0,
			instaforming: false,
			permaformChance: 0,
			permaformPopulation: 0,
			repairFactor: 1,
			starbaseRepairFactor: 1,
			starbaseCostFactor: 1,
			innateMining: false,
			innateResources: false,
			innateScanner: false,
			innatePopulationFactor: 1,
			canBuildDefenses: true,
			livesOnStarbases: false
		},
		JoaT: {
			prt: '',
			pointCost: 66,
			startingTechLevels: {
				energy: 3,
				weapons: 3,
				propulsion: 3,
				construction: 3,
				electronics: 3,
				biotechnology: 3
			},
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
					purpose: 'FighterScout'
				},
				{
					name: 'Stalwart Defender',
					hullName: 'Destroyer',
					purpose: 'FighterScout'
				}
			],
			startingPlanets: [
				{
					population: 25000,
					habPenaltyFactor: 0,
					hasStargate: false,
					hasMassDriver: false,
					startingFleets: null
				}
			],
			techCostOffset: {},
			mineralsPerSingleMineralPacket: 100,
			mineralsPerMixedMineralPacket: 40,
			packetResourceCost: 10,
			packetMineralCostFactor: 1,
			packetReceiverFactor: 1,
			packetDecayFactor: 1,
			packetOverSafeWarpPenalty: 0,
			packetBuiltInScanner: false,
			detectPacketDestinationStarbases: false,
			detectAllPackets: false,
			packetTerraformChance: 0,
			packetPermaformChance: 0,
			packetPermaTerraformSizeUnit: 100,
			canGateCargo: false,
			canDetectStargatePlanets: false,
			shipsVanishInVoid: true,
			builtInScannerMultiplier: 20,
			techsCostExtraLevel: 4,
			freighterGrowthFactor: 0,
			growthFactor: 1,
			maxPopulationOffset: 0.2,
			builtInCloakUnits: 0,
			stealsResearch: {},
			freeCargoCloaking: false,
			mineFieldsAreScanners: false,
			mineFieldRateMoveFactor: 0,
			mineFieldSafeWarpBonus: 0,
			mineFieldMinDecayFactor: 1,
			mineFieldBaseDecayRate: 0.02,
			mineFieldPlanetDecayRate: 0.04,
			mineFieldMaxDecayRate: 0.5,
			canDetonateMineFields: false,
			mineFieldDetonateDecayRate: 0.25,
			discoverDesignOnScan: false,
			canRemoteMineOwnPlanets: false,
			invasionAttackBonus: 1,
			invasionDefendBonus: 1,
			movementBonus: 0,
			instaforming: false,
			permaformChance: 0,
			permaformPopulation: 0,
			repairFactor: 1,
			starbaseRepairFactor: 1,
			starbaseCostFactor: 1,
			innateMining: false,
			innateResources: false,
			innateScanner: false,
			innatePopulationFactor: 1,
			canBuildDefenses: true,
			livesOnStarbases: false
		},
		PP: {
			prt: '',
			pointCost: 66,
			startingTechLevels: {
				energy: 4
			},
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
			startingPlanets: [
				{
					population: 25000,
					habPenaltyFactor: 0,
					hasStargate: false,
					hasMassDriver: true,
					startingFleets: null
				},
				{
					population: 10000,
					habPenaltyFactor: 1,
					hasStargate: false,
					hasMassDriver: true,
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
			packetOverSafeWarpPenalty: 0,
			packetBuiltInScanner: true,
			detectPacketDestinationStarbases: true,
			detectAllPackets: true,
			packetTerraformChance: 0.5,
			packetPermaformChance: 0.001,
			packetPermaTerraformSizeUnit: 100,
			canGateCargo: false,
			canDetectStargatePlanets: false,
			shipsVanishInVoid: true,
			builtInScannerMultiplier: 0,
			techsCostExtraLevel: 3,
			freighterGrowthFactor: 0,
			growthFactor: 1,
			maxPopulationOffset: 0,
			builtInCloakUnits: 0,
			stealsResearch: {},
			freeCargoCloaking: false,
			mineFieldsAreScanners: false,
			mineFieldRateMoveFactor: 0,
			mineFieldSafeWarpBonus: 0,
			mineFieldMinDecayFactor: 1,
			mineFieldBaseDecayRate: 0.02,
			mineFieldPlanetDecayRate: 0.04,
			mineFieldMaxDecayRate: 0.5,
			canDetonateMineFields: false,
			mineFieldDetonateDecayRate: 0.25,
			discoverDesignOnScan: false,
			canRemoteMineOwnPlanets: false,
			invasionAttackBonus: 1,
			invasionDefendBonus: 1,
			movementBonus: 0,
			instaforming: false,
			permaformChance: 0,
			permaformPopulation: 0,
			repairFactor: 1,
			starbaseRepairFactor: 1,
			starbaseCostFactor: 1,
			innateMining: false,
			innateResources: false,
			innateScanner: false,
			innatePopulationFactor: 1,
			canBuildDefenses: true,
			livesOnStarbases: false
		},
		SD: {
			prt: '',
			pointCost: 66,
			startingTechLevels: {},
			startingFleets: null,
			startingPlanets: [
				{
					population: 25000,
					habPenaltyFactor: 0,
					hasStargate: false,
					hasMassDriver: false,
					startingFleets: null
				}
			],
			techCostOffset: {},
			mineralsPerSingleMineralPacket: 100,
			mineralsPerMixedMineralPacket: 40,
			packetResourceCost: 10,
			packetMineralCostFactor: 1,
			packetReceiverFactor: 1,
			packetDecayFactor: 1,
			packetOverSafeWarpPenalty: 0,
			packetBuiltInScanner: false,
			detectPacketDestinationStarbases: false,
			detectAllPackets: false,
			packetTerraformChance: 0,
			packetPermaformChance: 0,
			packetPermaTerraformSizeUnit: 100,
			canGateCargo: false,
			canDetectStargatePlanets: false,
			shipsVanishInVoid: true,
			builtInScannerMultiplier: 0,
			techsCostExtraLevel: 3,
			freighterGrowthFactor: 0,
			growthFactor: 1,
			maxPopulationOffset: 0,
			builtInCloakUnits: 0,
			stealsResearch: {},
			freeCargoCloaking: false,
			mineFieldsAreScanners: false,
			mineFieldRateMoveFactor: 0,
			mineFieldSafeWarpBonus: 0,
			mineFieldMinDecayFactor: 1,
			mineFieldBaseDecayRate: 0.02,
			mineFieldPlanetDecayRate: 0.04,
			mineFieldMaxDecayRate: 0.5,
			canDetonateMineFields: false,
			mineFieldDetonateDecayRate: 0.25,
			discoverDesignOnScan: false,
			canRemoteMineOwnPlanets: false,
			invasionAttackBonus: 1,
			invasionDefendBonus: 1,
			movementBonus: 0,
			instaforming: false,
			permaformChance: 0,
			permaformPopulation: 0,
			repairFactor: 1,
			starbaseRepairFactor: 1,
			starbaseCostFactor: 1,
			innateMining: false,
			innateResources: false,
			innateScanner: false,
			innatePopulationFactor: 1,
			canBuildDefenses: true,
			livesOnStarbases: false
		},
		SS: {
			prt: '',
			pointCost: 66,
			startingTechLevels: {},
			startingFleets: null,
			startingPlanets: [
				{
					population: 25000,
					habPenaltyFactor: 0,
					hasStargate: false,
					hasMassDriver: false,
					startingFleets: null
				}
			],
			techCostOffset: {},
			mineralsPerSingleMineralPacket: 100,
			mineralsPerMixedMineralPacket: 40,
			packetResourceCost: 10,
			packetMineralCostFactor: 1,
			packetReceiverFactor: 1,
			packetDecayFactor: 1,
			packetOverSafeWarpPenalty: 0,
			packetBuiltInScanner: false,
			detectPacketDestinationStarbases: false,
			detectAllPackets: false,
			packetTerraformChance: 0,
			packetPermaformChance: 0,
			packetPermaTerraformSizeUnit: 100,
			canGateCargo: false,
			canDetectStargatePlanets: false,
			shipsVanishInVoid: true,
			builtInScannerMultiplier: 0,
			techsCostExtraLevel: 3,
			freighterGrowthFactor: 0,
			growthFactor: 1,
			maxPopulationOffset: 0,
			builtInCloakUnits: 0,
			stealsResearch: {},
			freeCargoCloaking: false,
			mineFieldsAreScanners: false,
			mineFieldRateMoveFactor: 0,
			mineFieldSafeWarpBonus: 0,
			mineFieldMinDecayFactor: 1,
			mineFieldBaseDecayRate: 0.02,
			mineFieldPlanetDecayRate: 0.04,
			mineFieldMaxDecayRate: 0.5,
			canDetonateMineFields: false,
			mineFieldDetonateDecayRate: 0.25,
			discoverDesignOnScan: false,
			canRemoteMineOwnPlanets: false,
			invasionAttackBonus: 1,
			invasionDefendBonus: 1,
			movementBonus: 0,
			instaforming: false,
			permaformChance: 0,
			permaformPopulation: 0,
			repairFactor: 1,
			starbaseRepairFactor: 1,
			starbaseCostFactor: 1,
			innateMining: false,
			innateResources: false,
			innateScanner: false,
			innatePopulationFactor: 1,
			canBuildDefenses: true,
			livesOnStarbases: false
		},
		WM: {
			prt: '',
			pointCost: 66,
			startingTechLevels: {},
			startingFleets: null,
			startingPlanets: [
				{
					population: 25000,
					habPenaltyFactor: 0,
					hasStargate: false,
					hasMassDriver: false,
					startingFleets: null
				}
			],
			techCostOffset: {},
			mineralsPerSingleMineralPacket: 100,
			mineralsPerMixedMineralPacket: 40,
			packetResourceCost: 10,
			packetMineralCostFactor: 1,
			packetReceiverFactor: 1,
			packetDecayFactor: 1,
			packetOverSafeWarpPenalty: 0,
			packetBuiltInScanner: false,
			detectPacketDestinationStarbases: false,
			detectAllPackets: false,
			packetTerraformChance: 0,
			packetPermaformChance: 0,
			packetPermaTerraformSizeUnit: 100,
			canGateCargo: false,
			canDetectStargatePlanets: false,
			shipsVanishInVoid: true,
			builtInScannerMultiplier: 0,
			techsCostExtraLevel: 3,
			freighterGrowthFactor: 0,
			growthFactor: 1,
			maxPopulationOffset: 0,
			builtInCloakUnits: 0,
			stealsResearch: {},
			freeCargoCloaking: false,
			mineFieldsAreScanners: false,
			mineFieldRateMoveFactor: 0,
			mineFieldSafeWarpBonus: 0,
			mineFieldMinDecayFactor: 1,
			mineFieldBaseDecayRate: 0.02,
			mineFieldPlanetDecayRate: 0.04,
			mineFieldMaxDecayRate: 0.5,
			canDetonateMineFields: false,
			mineFieldDetonateDecayRate: 0.25,
			discoverDesignOnScan: false,
			canRemoteMineOwnPlanets: false,
			invasionAttackBonus: 1,
			invasionDefendBonus: 1,
			movementBonus: 0,
			instaforming: false,
			permaformChance: 0,
			permaformPopulation: 0,
			repairFactor: 1,
			starbaseRepairFactor: 1,
			starbaseCostFactor: 1,
			innateMining: false,
			innateResources: false,
			innateScanner: false,
			innatePopulationFactor: 1,
			canBuildDefenses: true,
			livesOnStarbases: false
		}
	},
	lrtSpecs: {
		'1': {
			lrt: 0,
			pointCost: 0,
			startingTechLevels: {},
			techCostOffset: {},
			newTechCostFactor: 0,
			miniaturizationMax: 0,
			miniaturizationPerLevel: 0,
			noAdvancedScanners: false,
			scanRangeFactor: 0,
			fuelEfficiencyOffset: 0,
			maxPopulationOffset: 0,
			terraformCostOffset: {},
			mineralAlchemyCostOffset: 0,
			scrapMineralOffset: 0,
			scrapMineralOffsetStarbase: 0,
			scrapResourcesOffset: 0,
			scrapResourcesOffsetStarbase: 0,
			startingPopulationFactor: 0,
			starbaseBuiltInCloakUnits: 0,
			starbaseCostFactor: 0,
			researchFactor: 0,
			researchSplashDamage: 0,
			shieldStrengthFactor: 0,
			shieldRegenerationRate: 0,
			engineFailureRate: 0,
			engineReliableSpeed: 0
		},
		'1024': {
			lrt: 0,
			pointCost: 0,
			startingTechLevels: {},
			techCostOffset: {},
			newTechCostFactor: 0,
			miniaturizationMax: 0,
			miniaturizationPerLevel: 0,
			noAdvancedScanners: false,
			scanRangeFactor: 0,
			fuelEfficiencyOffset: 0,
			maxPopulationOffset: 0,
			terraformCostOffset: {},
			mineralAlchemyCostOffset: 0,
			scrapMineralOffset: 0,
			scrapMineralOffsetStarbase: 0,
			scrapResourcesOffset: 0,
			scrapResourcesOffsetStarbase: 0,
			startingPopulationFactor: 0,
			starbaseBuiltInCloakUnits: 0,
			starbaseCostFactor: 0,
			researchFactor: 0,
			researchSplashDamage: 0,
			shieldStrengthFactor: 0,
			shieldRegenerationRate: 0,
			engineFailureRate: 0,
			engineReliableSpeed: 0
		},
		'128': {
			lrt: 0,
			pointCost: 0,
			startingTechLevels: {},
			techCostOffset: {},
			newTechCostFactor: 0,
			miniaturizationMax: 0,
			miniaturizationPerLevel: 0,
			noAdvancedScanners: false,
			scanRangeFactor: 0,
			fuelEfficiencyOffset: 0,
			maxPopulationOffset: 0,
			terraformCostOffset: {},
			mineralAlchemyCostOffset: 0,
			scrapMineralOffset: 0,
			scrapMineralOffsetStarbase: 0,
			scrapResourcesOffset: 0,
			scrapResourcesOffsetStarbase: 0,
			startingPopulationFactor: 0,
			starbaseBuiltInCloakUnits: 0,
			starbaseCostFactor: 0,
			researchFactor: 0,
			researchSplashDamage: 0,
			shieldStrengthFactor: 0,
			shieldRegenerationRate: 0,
			engineFailureRate: 0,
			engineReliableSpeed: 0
		},
		'16': {
			lrt: 0,
			pointCost: 0,
			startingTechLevels: {},
			techCostOffset: {},
			newTechCostFactor: 0,
			miniaturizationMax: 0,
			miniaturizationPerLevel: 0,
			noAdvancedScanners: false,
			scanRangeFactor: 0,
			fuelEfficiencyOffset: 0,
			maxPopulationOffset: 0,
			terraformCostOffset: {},
			mineralAlchemyCostOffset: 0,
			scrapMineralOffset: 0,
			scrapMineralOffsetStarbase: 0,
			scrapResourcesOffset: 0,
			scrapResourcesOffsetStarbase: 0,
			startingPopulationFactor: 0,
			starbaseBuiltInCloakUnits: 0,
			starbaseCostFactor: 0,
			researchFactor: 0,
			researchSplashDamage: 0,
			shieldStrengthFactor: 0,
			shieldRegenerationRate: 0,
			engineFailureRate: 0,
			engineReliableSpeed: 0
		},
		'2': {
			lrt: 0,
			pointCost: 0,
			startingTechLevels: {},
			techCostOffset: {},
			newTechCostFactor: 0,
			miniaturizationMax: 0,
			miniaturizationPerLevel: 0,
			noAdvancedScanners: false,
			scanRangeFactor: 0,
			fuelEfficiencyOffset: 0,
			maxPopulationOffset: 0,
			terraformCostOffset: {},
			mineralAlchemyCostOffset: 0,
			scrapMineralOffset: 0,
			scrapMineralOffsetStarbase: 0,
			scrapResourcesOffset: 0,
			scrapResourcesOffsetStarbase: 0,
			startingPopulationFactor: 0,
			starbaseBuiltInCloakUnits: 0,
			starbaseCostFactor: 0,
			researchFactor: 0,
			researchSplashDamage: 0,
			shieldStrengthFactor: 0,
			shieldRegenerationRate: 0,
			engineFailureRate: 0,
			engineReliableSpeed: 0
		},
		'2048': {
			lrt: 0,
			pointCost: 0,
			startingTechLevels: {},
			techCostOffset: {},
			newTechCostFactor: 0,
			miniaturizationMax: 0,
			miniaturizationPerLevel: 0,
			noAdvancedScanners: false,
			scanRangeFactor: 0,
			fuelEfficiencyOffset: 0,
			maxPopulationOffset: 0,
			terraformCostOffset: {},
			mineralAlchemyCostOffset: 0,
			scrapMineralOffset: 0,
			scrapMineralOffsetStarbase: 0,
			scrapResourcesOffset: 0,
			scrapResourcesOffsetStarbase: 0,
			startingPopulationFactor: 0,
			starbaseBuiltInCloakUnits: 0,
			starbaseCostFactor: 0,
			researchFactor: 0,
			researchSplashDamage: 0,
			shieldStrengthFactor: 0,
			shieldRegenerationRate: 0,
			engineFailureRate: 0,
			engineReliableSpeed: 0
		},
		'256': {
			lrt: 0,
			pointCost: 0,
			startingTechLevels: {},
			techCostOffset: {},
			newTechCostFactor: 0,
			miniaturizationMax: 0,
			miniaturizationPerLevel: 0,
			noAdvancedScanners: false,
			scanRangeFactor: 0,
			fuelEfficiencyOffset: 0,
			maxPopulationOffset: 0,
			terraformCostOffset: {},
			mineralAlchemyCostOffset: 0,
			scrapMineralOffset: 0,
			scrapMineralOffsetStarbase: 0,
			scrapResourcesOffset: 0,
			scrapResourcesOffsetStarbase: 0,
			startingPopulationFactor: 0,
			starbaseBuiltInCloakUnits: 0,
			starbaseCostFactor: 0,
			researchFactor: 0,
			researchSplashDamage: 0,
			shieldStrengthFactor: 0,
			shieldRegenerationRate: 0,
			engineFailureRate: 0,
			engineReliableSpeed: 0
		},
		'32': {
			lrt: 0,
			pointCost: 0,
			startingTechLevels: {},
			techCostOffset: {},
			newTechCostFactor: 0,
			miniaturizationMax: 0,
			miniaturizationPerLevel: 0,
			noAdvancedScanners: false,
			scanRangeFactor: 0,
			fuelEfficiencyOffset: 0,
			maxPopulationOffset: 0,
			terraformCostOffset: {},
			mineralAlchemyCostOffset: 0,
			scrapMineralOffset: 0,
			scrapMineralOffsetStarbase: 0,
			scrapResourcesOffset: 0,
			scrapResourcesOffsetStarbase: 0,
			startingPopulationFactor: 0,
			starbaseBuiltInCloakUnits: 0,
			starbaseCostFactor: 0,
			researchFactor: 0,
			researchSplashDamage: 0,
			shieldStrengthFactor: 0,
			shieldRegenerationRate: 0,
			engineFailureRate: 0,
			engineReliableSpeed: 0
		},
		'4': {
			lrt: 0,
			pointCost: 0,
			startingTechLevels: {},
			techCostOffset: {},
			newTechCostFactor: 0,
			miniaturizationMax: 0,
			miniaturizationPerLevel: 0,
			noAdvancedScanners: false,
			scanRangeFactor: 0,
			fuelEfficiencyOffset: 0,
			maxPopulationOffset: 0,
			terraformCostOffset: {},
			mineralAlchemyCostOffset: 0,
			scrapMineralOffset: 0,
			scrapMineralOffsetStarbase: 0,
			scrapResourcesOffset: 0,
			scrapResourcesOffsetStarbase: 0,
			startingPopulationFactor: 0,
			starbaseBuiltInCloakUnits: 0,
			starbaseCostFactor: 0,
			researchFactor: 0,
			researchSplashDamage: 0,
			shieldStrengthFactor: 0,
			shieldRegenerationRate: 0,
			engineFailureRate: 0,
			engineReliableSpeed: 0
		},
		'4096': {
			lrt: 0,
			pointCost: 0,
			startingTechLevels: {},
			techCostOffset: {},
			newTechCostFactor: 0,
			miniaturizationMax: 0,
			miniaturizationPerLevel: 0,
			noAdvancedScanners: false,
			scanRangeFactor: 0,
			fuelEfficiencyOffset: 0,
			maxPopulationOffset: 0,
			terraformCostOffset: {},
			mineralAlchemyCostOffset: 0,
			scrapMineralOffset: 0,
			scrapMineralOffsetStarbase: 0,
			scrapResourcesOffset: 0,
			scrapResourcesOffsetStarbase: 0,
			startingPopulationFactor: 0,
			starbaseBuiltInCloakUnits: 0,
			starbaseCostFactor: 0,
			researchFactor: 0,
			researchSplashDamage: 0,
			shieldStrengthFactor: 0,
			shieldRegenerationRate: 0,
			engineFailureRate: 0,
			engineReliableSpeed: 0
		},
		'512': {
			lrt: 0,
			pointCost: 0,
			startingTechLevels: {},
			techCostOffset: {},
			newTechCostFactor: 0,
			miniaturizationMax: 0,
			miniaturizationPerLevel: 0,
			noAdvancedScanners: false,
			scanRangeFactor: 0,
			fuelEfficiencyOffset: 0,
			maxPopulationOffset: 0,
			terraformCostOffset: {},
			mineralAlchemyCostOffset: 0,
			scrapMineralOffset: 0,
			scrapMineralOffsetStarbase: 0,
			scrapResourcesOffset: 0,
			scrapResourcesOffsetStarbase: 0,
			startingPopulationFactor: 0,
			starbaseBuiltInCloakUnits: 0,
			starbaseCostFactor: 0,
			researchFactor: 0,
			researchSplashDamage: 0,
			shieldStrengthFactor: 0,
			shieldRegenerationRate: 0,
			engineFailureRate: 0,
			engineReliableSpeed: 0
		},
		'64': {
			lrt: 0,
			pointCost: 0,
			startingTechLevels: {},
			techCostOffset: {},
			newTechCostFactor: 0,
			miniaturizationMax: 0,
			miniaturizationPerLevel: 0,
			noAdvancedScanners: false,
			scanRangeFactor: 0,
			fuelEfficiencyOffset: 0,
			maxPopulationOffset: 0,
			terraformCostOffset: {},
			mineralAlchemyCostOffset: 0,
			scrapMineralOffset: 0,
			scrapMineralOffsetStarbase: 0,
			scrapResourcesOffset: 0,
			scrapResourcesOffsetStarbase: 0,
			startingPopulationFactor: 0,
			starbaseBuiltInCloakUnits: 0,
			starbaseCostFactor: 0,
			researchFactor: 0,
			researchSplashDamage: 0,
			shieldStrengthFactor: 0,
			shieldRegenerationRate: 0,
			engineFailureRate: 0,
			engineReliableSpeed: 0
		},
		'8': {
			lrt: 0,
			pointCost: 0,
			startingTechLevels: {},
			techCostOffset: {},
			newTechCostFactor: 0,
			miniaturizationMax: 0,
			miniaturizationPerLevel: 0,
			noAdvancedScanners: false,
			scanRangeFactor: 0,
			fuelEfficiencyOffset: 0,
			maxPopulationOffset: 0,
			terraformCostOffset: {},
			mineralAlchemyCostOffset: 0,
			scrapMineralOffset: 0,
			scrapMineralOffsetStarbase: 0,
			scrapResourcesOffset: 0,
			scrapResourcesOffsetStarbase: 0,
			startingPopulationFactor: 0,
			starbaseBuiltInCloakUnits: 0,
			starbaseCostFactor: 0,
			researchFactor: 0,
			researchSplashDamage: 0,
			shieldStrengthFactor: 0,
			shieldRegenerationRate: 0,
			engineFailureRate: 0,
			engineReliableSpeed: 0
		},
		'8192': {
			lrt: 0,
			pointCost: 0,
			startingTechLevels: {},
			techCostOffset: {},
			newTechCostFactor: 0,
			miniaturizationMax: 0,
			miniaturizationPerLevel: 0,
			noAdvancedScanners: false,
			scanRangeFactor: 0,
			fuelEfficiencyOffset: 0,
			maxPopulationOffset: 0,
			terraformCostOffset: {},
			mineralAlchemyCostOffset: 0,
			scrapMineralOffset: 0,
			scrapMineralOffsetStarbase: 0,
			scrapResourcesOffset: 0,
			scrapResourcesOffsetStarbase: 0,
			startingPopulationFactor: 0,
			starbaseBuiltInCloakUnits: 0,
			starbaseCostFactor: 0,
			researchFactor: 0,
			researchSplashDamage: 0,
			shieldStrengthFactor: 0,
			shieldRegenerationRate: 0,
			engineFailureRate: 0,
			engineReliableSpeed: 0
		}
	}
};
