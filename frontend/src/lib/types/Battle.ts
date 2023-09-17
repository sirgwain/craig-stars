import type { Universe } from '$lib/services/Universe';
import { flatten, groupBy, sumBy } from 'lodash-es';
import type { Cargo } from './Cargo';
import type { MapObject } from './MapObject';
import type { Vector } from './Vector';

export type BattleRecord = {
	num: number;
	planetNum?: number;
	position: Vector;
	tokens: Token[];
	actionsPerRound: Array<TokenAction[]>;
	destroyedTokens: BattleRecordDestroyedToken[];
	stats: BattleRecordStats;
};

export type Token = {
	num: number;
	playerNum: number;
	designNum: number;
	position: Vector;
	startingQuantity: number;
	damage?: number;
	quantityDamaged?: number;
	quantity?: number;
	initiative?: number;
	movement?: number;
	tactic: BattleTactic | string;
	primaryTarget: BattleTarget | string;
	secondaryTarget: BattleTarget | string;
	attackWho: BattleAttackWho | string;
};

export type BattleRecordDestroyedToken = {
	num: number;
	playerNum: number;
	designNum: number;
	quantity: number;
};

export type TokenAction = {
	type: number;
	tokenNum: number;
	round: number;
	from: Vector;
	to: Vector;
	slot?: number;
	targetNum?: number;
	target?: Token;
	tokensDestroyed?: number;
	damageDoneShields?: number;
	damageDoneArmor?: number;
	torpedoHits?: number;
	torpedoMisses?: number;
};

export type BattleRecordStats = {
	numPlayers?: number;
	numShipsByPlayer?: { [key: number]: number };
	shipsDestroyedByPlayer?: { [key: number]: number };
	damageTakenByPlayer?: { [key: number]: number };
	cargoLostByPlayer?: { [key: number]: Cargo };
};

export enum BattleTactic {
	Disengage = 'Disengage',
	DisengageIfChallenged = 'DisengageIfChallenged',
	MinimizeDamageToSelf = 'MinimizeDamageToSelf',
	MaximizeNetDamage = 'MaximizeNetDamage',
	MaximizeDamageRatio = 'MaximizeDamageRatio',
	MaximizeDamage = 'MaximizeDamage'
}

export enum BattleTarget {
	None = '',
	Any = 'Any',
	Starbase = 'Starbase',
	ArmedShips = 'ArmedShips',
	BombersFreighters = 'BombersFreighters',
	UnarmedShips = 'UnarmedShips',
	FuelTransports = 'FuelTransports',
	Freighters = 'Freighters'
}

export enum BattleAttackWho {
	Enemies = 'Enemies',
	EnemiesAndNeutrals = 'EnemiesAndNeutrals',
	Everyone = 'Everyone'
}
export enum TokenActionType {
	Fire,
	BeamFire,
	TorpedoFire,
	Move,
	RanAway
}

// a phase token is a token combined with a position
export type PhaseToken = {
	action?: TokenAction;
	ranAway?: boolean;
	destroyedPhase?: number;
	target?: boolean;
} & Token &
	Vector;

type TokensByLocation = Record<string, PhaseToken[]>;

export class Battle implements BattleRecord {
	constructor(public num: number, public position: Vector, record?: BattleRecord) {
		Object.assign(this, record);
		this.totalPhases = sumBy(this.actionsPerRound, (a) => a.length);
		this.totalRounds = this.actionsPerRound.length;
		this.tokens.forEach((t) => (t.quantity = t.startingQuantity));
		this.buildPhaseTokensForBattle();
		this.tokens.sort((a, b) => a.num - b.num);
		this.actions = flatten(this.actionsPerRound);
	}

	destroyedTokens: BattleRecordDestroyedToken[] = [];
	stats: BattleRecordStats = {
		numPlayers: 0,
		numShipsByPlayer: {},
		shipsDestroyedByPlayer: {},
		damageTakenByPlayer: {},
		cargoLostByPlayer: {}
	};

	planetNum?: number | undefined;
	tokens: Token[] = [];
	actionsPerRound: TokenAction[][] = [];
	actions: TokenAction[] = [];
	totalPhases: number;
	totalRounds: number;

	private tokensByPhase: PhaseToken[][] = [];
	private tokensByPhaseByLocation: TokensByLocation[] = [];

	getToken(num: number): Token | undefined {
		if (num > 0 && num <= this.tokens.length) {
			return this.tokens[num - 1];
		}
	}

	// get all remaining tokens at a location for a phase
	getTokensAtLocation(phase: number, x: number, y: number): PhaseToken[] | undefined {
		const phaseTokens = this.tokensByPhaseByLocation[phase];
		if (phaseTokens) {
			const remainingPhaseTokens = phaseTokens[getTokenLocationKey(x, y)]?.filter(
				(t) => !(t.ranAway || (t.destroyedPhase && phase > t.destroyedPhase))
			);
			// return undefined if we don't have any remaining tokens at this location for this phase
			if (remainingPhaseTokens?.length) {
				return remainingPhaseTokens;
			}
		}
	}

	// get the token performing an action for a phase
	getActionToken(phase: number): PhaseToken | undefined {
		return flatten(Object.values(this.tokensByPhaseByLocation[phase])).find((t) => t.action);
	}

	// get the action being performed for a phase
	getActionForPhase(phase: number) {
		return this.getActionToken(phase)?.action;
	}

	getTokenForPhase(num: number, phase: number) {
		return this.tokensByPhase[phase].find((t) => t.num == num);
	}

	getTargetForPhase(phase: number): PhaseToken | undefined {
		return this.tokensByPhase[phase].find((t) => t.target);
	}

	private buildPhaseTokensForBattle() {
		this.tokensByPhaseByLocation = [];

		// starting token configuration
		let tokens: PhaseToken[] = this.tokens.map((t) => ({
			...t,
			x: t.position.x,
			y: t.position.y
		}));

		// set the first phase to our base tokens
		this.tokensByPhaseByLocation.push(groupBy(tokens, (t) => getTokenLocationKey(t.x, t.y)));
		this.tokensByPhase.push(tokens);

		// a phase is incremented per action
		let phase = 1;
		for (let round = 1; round < this.actionsPerRound.length; round++) {
			const roundActions = this.actionsPerRound[round];
			for (let actionIndex = 0; actionIndex < roundActions.length; actionIndex++, phase++) {
				// find the action for this phase
				const action = roundActions[actionIndex];
				const phaseTokens = tokens.map((t) => {
					// clone each token for this phase
					const phaseToken = structuredClone(t);
					phaseToken.target = false; // clear out active target

					// if this token is being acted upon, update it
					if (phaseToken.num == action.tokenNum) {
						phaseToken.action = action;
						if (action.type == TokenActionType.Move) {
							phaseToken.x = action.to?.x ?? phaseToken.x;
							phaseToken.y = action.to?.y ?? phaseToken.y;
						} else if (action.type == TokenActionType.RanAway) {
							phaseToken.ranAway = true;
						}
					} else {
						// this token is idle, remove the action
						phaseToken.action = undefined;
					}
					if (action.target && action.targetNum === phaseToken.num) {
						Object.assign(phaseToken, action.target);
						phaseToken.target = true;
						if (!action.target.quantity) {
							// token was destroyed
							phaseToken.destroyedPhase = phase;
						}
					}
					return phaseToken;
				});
				// keep track of our progress
				tokens = phaseTokens;
				this.tokensByPhase.push(tokens);
				this.tokensByPhaseByLocation.push(
					groupBy(phaseTokens, (t) => getTokenLocationKey(t.x, t.y))
				);
			}
		}
	}
}

export const getTokenLocationKey = (x: number, y: number): string => `${x}-${y}`;

export function getNumShips(record: BattleRecord): number {
	return Object.values(record.stats.numShipsByPlayer ?? {}).reduce((count, num) => count + num, 0);
}

export function getOurShips(record: BattleRecord, allies: Set<number>): number {
	let count = 0;
	allies.forEach(
		(ally) => (count += record.stats?.numShipsByPlayer ? record.stats?.numShipsByPlayer[ally] : 0)
	);
	return count;
}

export function getTheirShips(record: BattleRecord, allies: Set<number>): number {
	let count = 0;
	Object.entries(record.stats?.numShipsByPlayer ?? {}).forEach((entry) => {
		const playerNum: number = parseInt(entry[0]);
		const numShips = entry[1];
		if (!allies.has(playerNum)) {
			count += numShips;
		}
	});
	return count;
}

export function getOurDead(record: BattleRecord, allies: Set<number>): number {
	let count = 0;
	allies.forEach(
		(ally) =>
			(count += record.stats?.shipsDestroyedByPlayer
				? record.stats?.shipsDestroyedByPlayer[ally]
				: 0)
	);
	return count;
}

export function getTheirDead(record: BattleRecord, allies: Set<number>): number {
	let count = 0;
	Object.entries(record.stats?.shipsDestroyedByPlayer ?? {}).forEach((entry) => {
		const playerNum: number = parseInt(entry[0]);
		const numShips = entry[1];
		if (!allies.has(playerNum)) {
			count += numShips;
		}
	});
	return count;
}

// get a target for the scanner so we can "goto" a battle and select this mapobject
export function getScannerTarget(
	battle: BattleRecord,
	universe: Universe
): MapObject | undefined {
	if (battle.planetNum) {
		return universe.getPlanet(battle.planetNum);
	} else {
		const myMapObjectsAtPosition = universe.getMyMapObjectsByPosition(battle.position);
		const mapObjectsAtPosition = universe.getMapObjectsByPosition(battle.position);

		if (myMapObjectsAtPosition.length > 0) {
			return myMapObjectsAtPosition[0];
		} else if (mapObjectsAtPosition) {
			return mapObjectsAtPosition[0];
		}
	}
}
