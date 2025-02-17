import { CommandedPlayer } from '$lib/types/Player';
import { defaultRules } from '$lib/types/Rules';
import type { GameWithPlayers, ShipDesign, Vector } from '$lib/types/cs';
import {
	DensityNormal,
	GameStartModeNormal,
	GameStateWaitingForPlayers,
	PlayerPositionsModerate,
	SizeSmall,
	type PlayerStatus,
	type VictoryConditions
} from '$lib/types/cs';
import { TechService } from './TechService';
import { Universe } from './Universe';

export class FullGame implements GameWithPlayers {
	id = 0;
	createdAt = '';
	updatedAt = '';
	hostId = 0;
	seed = 0;
	name = '';
	hash = '';
	state = GameStateWaitingForPlayers;
	numPlayers = 0;
	openPlayerSlots = 0;
	quickStartTurns = 0;
	size = SizeSmall;
	area: Vector = { x: 0, y: 0 };
	density = DensityNormal;
	playerPositions = PlayerPositionsModerate;
	randomEvents = false;
	computerPlayersFormAlliances = false;
	publicPlayerScores = false;
	maxMinerals = false;
	acceleratedPlay = false;
	startMode = GameStartModeNormal;
	year = 2400;
	victoryConditions: VictoryConditions = {
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
	};
	public = false;
	victorDeclared = false;
	archived = false;
	rules = defaultRules;
	players: PlayerStatus[] = [];

	// some data that is loaded
	player: CommandedPlayer = new CommandedPlayer();
	universe: Universe = new Universe();
	techs = new TechService();

	isMultiplayer(): boolean {
		// we are multi player if any of the players are not ai controlled and not us
		return (
			this.openPlayerSlots > 0 ||
			this.players.findIndex((p) => p.num != this.player.num && !p.aiControlled) == -1
		);
	}

	isSinglePlayer(): boolean {
		return !this.isMultiplayer();
	}

	validateDesign(design: ShipDesign): { valid: boolean; reason?: string } {
		// TODO: add more validations

		// if we have a design with this name already, it is invalid
		const designsWithName = this.universe.getMyDesigns().filter((d) => d.name === design.name);
		if (designsWithName.length > 1 || (designsWithName.length === 1 && !design.id)) {
			return { valid: false, reason: `Another design named ${design.name} exists` };
		}
		return { valid: true };
	}
}
