import type { Cost } from './Cost';
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
	Normal = 'Normal', // regular start
	MidGame = 'MidGame', // further tech levels, pop growth
	LateGame = 'LateGame',
	EndGame = 'EndGame'
}

export enum NewGamePlayerType {
	Host = 'Host',
	Invite = 'Invite',
	Open = 'Open',
	AI = 'AI'
}

export interface NewGamePlayer {
	type: NewGamePlayerType;
	userId?: number;
	raceId?: number;
}

export interface GameSettings {
	name: string;
	quickStartTurns?: number;
	size: Size;
	area?: Vector;
	density: Density;
	playerPositions: PlayerPositions;
	randomEvents: boolean;
	computerPlayersFormAlliances: boolean;
	publicPlayerScores: boolean;
	startMode: GameStartMode;
	year?: number;
	state?: string;
	victoryConditions?: VictoryConditions;
	players: NewGamePlayer[];
}

export interface Game {
	id: number;
	createdAt: string;
	updatedat: string;
	deletedAt: null;
	name: string;
	quickStartTurns: number;
	size: Size;
	area: Vector;
	density: Density;
	playerPositions: PlayerPositions;
	randomEvents: boolean;
	computerPlayersFormAlliances: boolean;
	publicPlayerScores: boolean;
	startMode: GameStartMode;
	year: number;
	state: string;
	victoryConditions: VictoryConditions;
	victorDeclared: boolean;
	rules: Rules;
}

export interface Rules {
	id: number;
	createdAt: string;
	updatedat: string;
	deletedAt: null;
	gameId: number;
	seed: number;
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
	randomEventChances: null;
	randomMineralDepositBonusRange: number[];
	wormholeCloak: number;
	wormholeMinDistance: number;
	wormholeStatsByStability: null;
	wormholePairsForSize: null;
	mineFieldStatsByType: null;
	repairRates: null;
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
	packetDecayRate: null;
	maxTechLevel: number;
	techBaseCost: null;
	prtSpecs: null;
	lrtSpecs: null;
}

export interface VictoryConditions {
	conditions: string[];
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
}
