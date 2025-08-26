import {
	Density,
	GameStartMode,
	GameState,
	PlayerPositions,
	Size,
	VectorSchema,
	type Vector
} from '$lib/types/cs-proto';
import {
	GameWithPlayersSchema,
	VictoryConditionsSchema,
	type Game,
	type GameWithPlayers,
	type VictoryConditions
} from '$lib/types/cs-proto';
import type { PlayerStatus } from '$lib/types/cs-proto';
import { defaultRules } from '$lib/types/Rules';
import { create } from '@bufbuild/protobuf';
import { TimestampSchema } from '@bufbuild/protobuf/wkt';

export class FullGame implements Game {
	$typeName: 'craig_stars.v1.Game';
	$unknown = undefined;
	id = BigInt(0);
	createdAt = create(TimestampSchema, {});
	updatedAt = create(TimestampSchema, {});
	hostId = BigInt(0);
	seed = BigInt(0);
	name = '';
	hash = '';
	state = GameState.WAITING_FOR_PLAYERS;
	numPlayers = 0;
	openPlayerSlots = 0;
	quickStartTurns = 0;
	size = Size.SMALL;
	area: Vector = create(VectorSchema, { x: 0, y: 0 });
	density = Density.NORMAL;
	playerPositions = PlayerPositions.MODERATE;
	randomEvents = false;
	computerPlayersFormAlliances = false;
	publicPlayerScores = false;
	maxMinerals = false;
	startMode = GameStartMode.UNSPECIFIED; // Normal
	year = 2400;
	victoryConditions: VictoryConditions = create(VictoryConditionsSchema, {
		conditions: 0,
		numCriteriaRequired: 0,
		yearsPassed: 0,
		ownPlanets: 0,
		attainTechLevel: 0,
		attainTechLevelNumFields: 0,
		exceedsScore: 0,
		exceedsSecondPlaceScore: 0,
		productionCapacity: 0,
		ownCapitalShips: 0,
		highestScoreAfterYears: 0
	});
	public = false;
	victorDeclared = false;
	archived = false;
	rules = defaultRules;
	players: PlayerStatus[] = [];

	constructor() {
		this.$typeName = 'craig_stars.v1.Game';
	}

	isMultiplayer(): boolean {
		// we are multi player if any of the players are not ai controlled and not us
		return (
			this.openPlayerSlots > 0 ||
			this.players.reduce((count, p) => count + (p.aiControlled ? 0 : 1), 0) > 1
		);
	}

	isSinglePlayer(): boolean {
		return !this.isMultiplayer();
	}

	toGameWithPlayers(): GameWithPlayers {
		return create(GameWithPlayersSchema, {
			game: this,
			players: this.players
		});
	}
}
