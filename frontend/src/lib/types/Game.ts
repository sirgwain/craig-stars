import type { PlayerStatus } from './Player';
import type { Race } from './Race';
import type { Rules } from './Rules';
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
	Guest = 'Guest',
	Open = 'Open',
	AI = 'AI'
}

export enum AIDifficulty {
	// Easy = 'Easy',
	Normal = 'Normal',
	// Hard = 'Hard',
	Cheater = 'Cheater'
}

export type NewGamePlayer = {
	type: NewGamePlayerType;
	aiDifficulty?: AIDifficulty;
	userId?: number;
	race?: Race;
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
	maxMinerals?: boolean;
	acceleratedPlay?: boolean;
	startMode?: GameStartMode;
	year?: number;
	victoryConditions: VictoryConditions;
};

export enum GameState {
	Setup = 'Setup',
	WaitingForPlayers = 'WaitingForPlayers',
	GeneratingUniverse = 'GeneratingUniverse',
	GeneratingTurn = 'GeneratingTurn',
	GeneratingTurnError = 'GeneratingTurnError'
}

export type Game = {
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
	maxMinerals: boolean;
	acceleratedPlay: boolean;
	public?: boolean;
	startMode: GameStartMode;
	year: number;
	victoryConditions: VictoryConditions;
	victorDeclared: boolean;
	archived: boolean;
	rules?: Rules;
	players: PlayerStatus[];
};

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
