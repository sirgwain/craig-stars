import type { CargoTransferRequest } from '$lib/types/Cargo';
import type { CommandedFleet, Fleet } from '$lib/types/Fleet';
import {
	Density,
	GameStartMode,
	GameState,
	PlayerPositions,
	Size,
	type Game,
	type VictoryConditions
} from '$lib/types/Game';
import { MapObjectType } from '$lib/types/MapObject';
import type { CommandedPlanet, Planet } from '$lib/types/Planet';
import {
	Player,
	type BattlePlan,
	type PlayerIntels,
	type PlayerStatus,
	type PlayerUniverse,
	type ProductionPlan,
	type TransportPlan
} from '$lib/types/Player';
import { defaultRules } from '$lib/types/Rules';
import type { Salvage } from '$lib/types/Salvage';
import type { ShipDesign } from '$lib/types/ShipDesign';
import type { SessionUser } from '$lib/types/User';
import type { Vector } from '$lib/types/Vector';
import { get } from 'svelte/store';
import { BattlePlanService } from './BattlePlanService';
import { updateGame, updateGameContext, updatePlayer, updateUniverse } from './Contexts';
import { DesignService } from './DesignService';
import { FleetService } from './FleetService';
import { GameService } from './GameService';
import { PlanetService } from './PlanetService';
import { PlayerService } from './PlayerService';
import { ProductionPlanService } from './ProductionPlanService';
import {
	commandMapObject,
	commandedFleet,
	commandedPlanet,
	currentSelectedWaypointIndex,
	selectMapObject,
	selectWaypoint,
	zoomToMapObject
} from './Stores';
import { TechService } from './TechService';
import { TransportPlanService } from './TransportPlanService';
import { Universe } from './Universe';

export class FullGame implements Game {
	id = 0;
	createdAt = '';
	updatedAt = '';
	hostId = 0;
	name = '';
	hash = '';
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
	rules = defaultRules;
	players: PlayerStatus[] = [];

	// some data that is loaded
	player: Player = new Player();
	universe: Universe = new Universe();
	techs = new TechService();
	playerStatusPollingInterval: number | undefined;

	// load this game from the server
	async load(id: number | string) {
		this.id = parseInt(id.toString());
		let pmos: PlayerUniverse & PlayerIntels = {
			designs: [],
			players: [],
			scores: [],
			battles: [],
			planets: [],
			fleets: [],
			starbases: [],
			mineFields: [],
			mineralPackets: [],
			salvages: [],
			wormholes: [],
			mysteryTraders: []
		};
		const game = await GameService.loadGame(id);
		Object.assign(this, game);
		if (this.state != GameState.Setup) {
			await Promise.all([
				GameService.loadFullPlayer(id).then((data) => {
					this.player = data;
				}),
				GameService.loadUniverse(id).then((u) => {
					pmos = u;
				}),
				// load techs the first time as well
				this.techs.fetch()
			]);
		}

		// setup the universe
		this.universe.setData(this.player.num, pmos);
		updateGameContext(this, this.player, this.universe);
		return this;
	}

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

	// command the player's homeworld (or the first planet they own, if their homeworld has been taken)
	commandHomeWorld() {
		const homeworld = this.universe.getHomeworld(this.player.num);
		if (homeworld) {
			commandMapObject(homeworld);
			selectMapObject(homeworld);
			zoomToMapObject(homeworld);
		} else {
			// command our first planet
			const planets = this.universe.getMyPlanets();
			if (planets.length > 0) {
				commandMapObject(planets[0]);
				selectMapObject(planets[0]);
				zoomToMapObject(planets[0]);
			}
		}
	}

	async submitTurn(): Promise<FullGame> {
		const resp = await PlayerService.submitTurn(this.id);
		if (resp) {
			Object.assign(this, resp.game);
			Object.assign(this.player, resp.player);
			if (resp.universe) {
				this.universe.setData(this.player.num, resp.universe);
			}
			updateGameContext(this, this.player, this.universe);
		}
		return this;
	}

	async forceGenerateTurn(): Promise<FullGame> {
		const resp = await GameService.forceGenerateTurn(this.id);
		Object.assign(this, resp.game);
		Object.assign(this.player, resp.player);
		if (resp.universe) {
			this.universe.setData(this.player.num, resp.universe);
		}
		updateGameContext(this, this.player, this.universe);
		return this;
	}

	async loadStatus(): Promise<Game> {
		const result = await GameService.loadGame(this.id);
		Object.assign(this, result);
		updateGame(this);
		return this;
	}

	async loadGuest(playerNum: number): Promise<SessionUser> {
		return await GameService.loadGuest(this.id, playerNum);
	}

	// start polling the server for player status
	async startPollingStatus(interval = 10000) {
		if (!this.playerStatusPollingInterval) {
			this.playerStatusPollingInterval = window.setInterval(async () => {
				this.loadStatus();
			}, interval);
		}
	}

	// stop polling the server for player status
	stopPollingStatus() {
		if (this.playerStatusPollingInterval) {
			window.clearInterval(this.playerStatusPollingInterval);
			this.playerStatusPollingInterval = undefined;
		}
	}

	async updatePlayerOrders() {
		const result = await PlayerService.updateOrders(this.player);
		if (result) {
			Object.assign(this.player, result.player);
			this.universe.updatePlanets(result.planets);
			updateUniverse(this.universe);
			updatePlayer(this.player);
		}
	}

	async updatePlayerRelations() {
		const result = await PlayerService.updateRelations(this.player);
		if (result) {
			Object.assign(this.player, result);
		}
	}

	async createBattlePlan(plan: BattlePlan) {
		const created = await BattlePlanService.create(this.id, plan);
		this.player.battlePlans = [...this.player.battlePlans, created];
		updatePlayer(this.player);
		return created;
	}

	async updateBattlePlan(plan: BattlePlan) {
		await BattlePlanService.update(this.id, plan);
		updatePlayer(this.player);
	}

	async deleteBattlePlan(num: number) {
		const { player, fleets, starbases } = await BattlePlanService.delete(this.id, num);
		Object.assign(this.player, player);
		this.universe.fleets = fleets;
		this.universe.starbases = starbases;
		updatePlayer(this.player);
	}

	async createProductionPlan(plan: ProductionPlan) {
		const created = await ProductionPlanService.create(this.id, plan);
		this.player.productionPlans = [...this.player.productionPlans, created];
		updatePlayer(this.player);
		return created;
	}

	async updateProductionPlan(plan: ProductionPlan) {
		await ProductionPlanService.update(this.id, plan);
		updatePlayer(this.player);
	}

	async deleteProductionPlan(num: number) {
		const player = await ProductionPlanService.delete(this.id, num);
		Object.assign(this.player, player);
		updatePlayer(this.player);
	}

	async createTransportPlan(plan: TransportPlan) {
		const created = await TransportPlanService.create(this.id, plan);
		this.player.transportPlans = [...this.player.transportPlans, created];
		updatePlayer(this.player);
		return created;
	}

	async updateTransportPlan(plan: TransportPlan) {
		await TransportPlanService.update(this.id, plan);
		updatePlayer(this.player);
	}

	async deleteTransportPlan(num: number) {
		const player = await TransportPlanService.delete(this.id, num);
		Object.assign(this.player, player);
		updatePlayer(this.player);
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

	async createDesign(design: ShipDesign): Promise<ShipDesign> {
		// update this design
		design = await DesignService.create(this.id, design);
		this.universe.addDesign(design);
		updateUniverse(this.universe);
		return design;
	}

	async updateDesign(design: ShipDesign) {
		// update this design
		design = await DesignService.update(this.id, design);
		this.universe.updateDesign(design);
		updateUniverse(this.universe);
	}

	async deleteDesign(num: number) {
		const { fleets, starbases } = await DesignService.delete(this.id, num);
		this.universe.fleets = fleets;
		this.universe.starbases = starbases;
		this.universe.resetMapObjectsByPosition();
		this.universe.resetMyMapObjectsByPosition();

		this.universe.designs = this.universe.designs.filter((d) => d.num != num);
		updateUniverse(this.universe);

		// reset our view to the homeworld, in case the commanded fleet had our deleted design
		this.commandHomeWorld();
	}

	async updateFleetOrders(fleet: CommandedFleet) {
		const selectedWaypointIndex = get(currentSelectedWaypointIndex);
		const updatedFleet = await FleetService.updateFleetOrders(fleet);
		fleet = Object.assign(fleet, updatedFleet);
		this.universe.updateFleet(fleet);
		commandedFleet.update(() => fleet);

		if (
			selectedWaypointIndex > -1 &&
			fleet.waypoints &&
			fleet.waypoints.length > selectedWaypointIndex
		) {
			selectWaypoint(fleet.waypoints[selectedWaypointIndex]);
		}

		updateUniverse(this.universe);
	}

	async renameFleet(fleet: CommandedFleet, name: string) {
		const selectedWaypointIndex = get(currentSelectedWaypointIndex);
		const updatedFleet = await FleetService.rename(fleet, name);
		fleet = Object.assign(fleet, updatedFleet);
		this.universe.updateFleet(fleet);
		commandedFleet.update(() => fleet);

		if (
			selectedWaypointIndex > -1 &&
			fleet.waypoints &&
			fleet.waypoints.length > selectedWaypointIndex
		) {
			selectWaypoint(fleet.waypoints[selectedWaypointIndex]);
		}

		updateUniverse(this.universe);
	}

	async updatePlanetOrders(planet: CommandedPlanet) {
		const resp = await PlanetService.updatePlanetOrders(planet);

		// changing the planet orders changes the player's spec
		Object.assign(this.player, resp.player);
		updatePlayer(this.player);

		planet = Object.assign(planet, resp.planet);
		this.universe.updatePlanet(planet);
		commandedPlanet.update(() => planet);
		updateUniverse(this.universe);
	}

	async transferCargo(
		fleet: CommandedFleet,
		dest: Fleet | Planet | Salvage,
		transferAmount: CargoTransferRequest
	) {
		const selectedWaypointIndex = get(currentSelectedWaypointIndex);
		const result = await FleetService.transferCargo(fleet, dest, transferAmount);

		if (result.dest?.type == MapObjectType.Planet) {
			const planet = result.dest as Planet;
			this.universe.updatePlanet(planet);

			// if we are currently commanding the destination, make sure it updates
			const p = get(commandedPlanet);
			if (p?.num == planet.num) {
				commandMapObject(planet);
			}
		} else if (result.dest?.type == MapObjectType.Fleet) {
			// update the destination fleet in the universe
			const destFleet = result.dest as Fleet;
			this.universe.updateFleet(destFleet);
		}

		if (result.salvages) {
			this.universe.updateSalvages(result.salvages);
		}

		fleet = Object.assign(fleet, result.fleet);
		this.universe.updateFleet(fleet);
		commandedFleet.update(() => fleet);

		if (
			selectedWaypointIndex > -1 &&
			fleet.waypoints &&
			fleet.waypoints.length > selectedWaypointIndex
		) {
			selectWaypoint(fleet.waypoints[selectedWaypointIndex]);
		}

		updateUniverse(this.universe);
	}

	async splitAll(fleet: CommandedFleet) {
		const selectedWaypointIndex = get(currentSelectedWaypointIndex);
		const updatedFleets = await FleetService.splitAll(fleet.gameId, fleet);
		const sourceFleet = updatedFleets.find((f) => f.num == fleet.num);
		if (sourceFleet) {
			fleet = Object.assign(fleet, sourceFleet);
			commandedFleet.update(() => fleet);
		}
		// update and add the new fleets to the universe
		this.universe.updateFleet(fleet);
		this.universe.addFleets(updatedFleets.filter((f) => f.num != fleet.num));

		if (
			selectedWaypointIndex > -1 &&
			fleet.waypoints &&
			fleet.waypoints.length > selectedWaypointIndex
		) {
			selectWaypoint(fleet.waypoints[selectedWaypointIndex]);
		}

		updateUniverse(this.universe);
	}

	async merge(fleet: CommandedFleet, fleetNums: number[]) {
		const selectedWaypointIndex = get(currentSelectedWaypointIndex);
		const updatedFleet = await FleetService.merge(fleet, fleetNums);

		this.universe.removeFleets(fleetNums);
		fleet = Object.assign(fleet, updatedFleet);
		this.universe.updateFleet(fleet);
		commandedFleet.update(() => fleet);

		if (
			selectedWaypointIndex > -1 &&
			fleet.waypoints &&
			fleet.waypoints.length > selectedWaypointIndex
		) {
			selectWaypoint(fleet.waypoints[selectedWaypointIndex]);
		}

		updateUniverse(this.universe);
	}
}
