import type { CommandedFleet } from '$lib/types/Fleet';
import {
	defaultRules,
	Density,
	GameStartMode,
	GameState,
	PlayerPositions,
	Size,
	type Game,
	type VictoryConditions
} from '$lib/types/Game';
import { Player, type PlayerMapObjects } from '$lib/types/Player';
import type { Vector } from '$lib/types/Vector';
import { commandedFleet } from './Context';
import { DesignService } from './DesignService';
import { FleetService } from './FleetService';
import { GameService } from './GameService';
import { PlayerService } from './PlayerService';
import { TechService } from './TechService';
import { Universe } from './Universe';

export class FullGame implements Game {
	id = 0;
	createdAt = '';
	updatedAt = '';
	hostId = 0;
	name = '';
	state = GameState.WaitingForPlayers;
	numPlayers = 0;
	openPlayerSlots = 0;
	quickStartTurns = 0;
	size = Size.Small;
	area: Vector = { x: 0, y: 0 };
	density = Density.Normal;
	playerPositions = PlayerPositions.Moderate;
	randomEvents = false;
	computerPlayersFormAlliances = false;
	publicPlayerScores = false;
	startMode = GameStartMode.Normal;
	year = 2400;
	victoryConditions: VictoryConditions = {
		conditions: [],
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
	victorDeclared = false;
	rules = defaultRules;

	// some data that is loaded
	player: Player = new Player();
	universe: Universe = new Universe();
	techs = new TechService();

	constructor(json?: Game) {
		Object.assign(this, json);
	}

	// load this game from the server
	async load(id: number | string) {
		this.id = parseInt(id.toString());
		let pmos: PlayerMapObjects = {
			planets: [],
			fleets: [],
			starbases: []
		};
		await Promise.all([
			GameService.loadGame(id).then((game) => Object.assign(this, game)),
			GameService.loadFullPlayer(id).then((data) => {
				this.player = data;
			}),
			GameService.loadUniverse(id).then((u) => {
				pmos = u;
			}),
			// load techs the first time as well
			this.techs.fetch()
		]);

		// setup the universe
		this.universe.playerNum = this.player.num;
		this.universe.setMapObjects(pmos);
		this.universe.setIntels(this.player);
		return this;
	}

	async submitTurn(): Promise<FullGame> {
		const resp = await PlayerService.submitTurn(this.id);
		if (resp) {
			Object.assign(this, resp.game);
			Object.assign(this.player, resp.player);
			this.universe.setMapObjects(resp.mapObjects);
			this.universe.setIntels(this.player);
		}
		return this;
	}

	async updatePlayerOrders() {
		const result = await PlayerService.updateOrders(this.player);
		if (result) {
			Object.assign(this.player, result.player);
			this.universe.planets = result.planets;
		}
		return this.player;
	}

	async updatePlayerPlans() {
		const result = await PlayerService.updatePlans(this.player);
		if (result) {
			Object.assign(this.player, result);
		}
		return this.player;
	}

	async deleteDesign(num: number) {
		const { fleets, starbases } = await DesignService.delete(this.id, num);
		this.universe.fleets = fleets;
		this.universe.starbases = starbases;
		this.universe.resetMyMapObjectsByPosition();

		this.player.designs = this.player.designs.filter((d) => d.num != num);
	}

	async updateFleetOrders(fleet: CommandedFleet) {
		const updatedFleet = await FleetService.updateFleetOrders(fleet);
		fleet = Object.assign(fleet, updatedFleet);
		this.universe.updateFleet(fleet);
		commandedFleet.update(() => fleet);
	}

	async splitAll(fleet: CommandedFleet) {
		const updatedFleets = await FleetService.splitAll(fleet.gameId, fleet);
		const sourceFleet = updatedFleets.find((f) => f.num == fleet.num);
		if (sourceFleet) {
			fleet = Object.assign(fleet, sourceFleet);
			commandedFleet.update(() => fleet);
		}
		// update and add the new fleets to the universe
		this.universe.updateFleet(fleet);
		this.universe.addFleets(updatedFleets.filter((f) => f.num != fleet.num));
	}

	async merge(fleet: CommandedFleet, fleetNums: number[]) {
		const updatedFleet = await FleetService.merge(fleet, fleetNums);

		this.universe.removeFleets(fleetNums);
		fleet = Object.assign(fleet, updatedFleet);
		this.universe.updateFleet(fleet);
		commandedFleet.update(() => fleet);
	}

	getPlanet(num: number) {
		return this.universe.getPlanet(num);
	}

	getPlayerName(playerNum: number | undefined) {
		if (playerNum && playerNum > 0 && playerNum <= this.player.playerIntels.length) {
			const intel = this.player.playerIntels[playerNum - 1];
			return intel.racePluralName ?? intel.name;
		}
		return 'unknown';
	}

	getPlayerColor(playerNum: number | undefined) {
		if (playerNum && playerNum > 0 && playerNum <= this.player.playerIntels.length) {
			const intel = this.player.playerIntels[playerNum - 1];
			return intel.color ?? '#FF0000';
		}
		return '#FF0000';
	}
}
