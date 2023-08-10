import { flatten, groupBy, sumBy } from 'lodash-es';
import { levelsAbove } from './TechLevel';
import type { Vector } from './Vector';

export type BattleRecord = {
	num: number;
	planetNum?: number;
	position: Vector;
	tokens: Token[];
	actionsPerRound: Array<TokenAction[]>;
};

export type Token = {
	num: number;
	playerNum: number;
	designNum: number;
	position: Vector;
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
	destroyed?: boolean;
} & Token &
	Vector;

type TokensByLocation = Record<string, PhaseToken[]>;

export class Battle implements BattleRecord {
	constructor(public num: number, public position: Vector, record?: BattleRecord) {
		Object.assign(this, record);
		this.totalPhases = sumBy(this.actionsPerRound, (a) => a.length);
		this.totalRounds = this.actionsPerRound.length;
		this.buildPhaseTokensForBattle();
		this.tokens.sort((a, b) => a.num - b.num);
		this.actions = flatten(this.actionsPerRound);
	}

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
				(t) => !(t.ranAway || t.destroyed)
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

		for (let round = 1; round < this.actionsPerRound.length; round++) {
			const roundActions = this.actionsPerRound[round];
			for (let actionIndex = 0; actionIndex < roundActions.length; actionIndex++) {
				// find the action for this phase
				const action = roundActions[actionIndex];
				const phaseTokens = tokens.map((t) => {
					// clone each token for this phase
					const clonedToken = structuredClone(t);

					// if this token is being acted upon, update it
					if (clonedToken.num == action.tokenNum) {
						clonedToken.action = action;
						if (action.type == TokenActionType.Move) {
							clonedToken.x = action.to?.x ?? clonedToken.x;
							clonedToken.y = action.to?.y ?? clonedToken.y;
						} else if (action.type == TokenActionType.RanAway) {
							clonedToken.ranAway = true;
						}
					} else {
						// this token is idle, remove the action
						clonedToken.action = undefined;
					}
					if (action.target && action.targetNum === clonedToken.num) {
						Object.assign(clonedToken, action.target);
					}
					return clonedToken;
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
