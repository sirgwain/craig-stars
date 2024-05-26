import { goto } from '$app/navigation';
import type { CargoTransferRequest } from '$lib/types/Cargo';
import { CommandedFleet, type Fleet, type ShipToken, type Waypoint } from '$lib/types/Fleet';
import type { Game, GameSettings } from '$lib/types/Game';
import { MapObjectType, None, equal, ownedBy, type MapObject, key } from '$lib/types/MapObject';
import {
	MessageTargetType,
	MessageType,
	getMapObjectTypeForMessageType,
	type Message
} from '$lib/types/Message';
import { CommandedPlanet, type Planet } from '$lib/types/Planet';
import {
	Player,
	type BattlePlan,
	type PlayerResponse,
	type ProductionPlan,
	type TransportPlan
} from '$lib/types/Player';
import { PlayerSettings } from '$lib/types/PlayerSettings';
import type { Salvage } from '$lib/types/Salvage';
import type { ShipDesign } from '$lib/types/ShipDesign';
import { findIndex, kebabCase } from 'lodash-es';
import { getContext } from 'svelte';
import {
	derived,
	get,
	writable,
	type Readable,
	type Writable,
	type Unsubscriber
} from 'svelte/store';
import { BattlePlanService } from './BattlePlanService';
import { DesignService } from './DesignService';
import { FleetService } from './FleetService';
import { FullGame } from './FullGame';
import { GameService } from './GameService';
import { rollover } from './Math';
import { PlanetService } from './PlanetService';
import { PlayerService } from './PlayerService';
import { ProductionPlanService } from './ProductionPlanService';
import { TransportPlanService } from './TransportPlanService';
import { Universe } from './Universe';

export const playerFinderKey = Symbol();
export const designFinderKey = Symbol();
export const gameKey = Symbol();

export type GameContext = {
	game: Readable<FullGame>;
	player: Readable<Player>;
	universe: Readable<Universe>;
	settings: Writable<PlayerSettings>;
	messageNum: Writable<number>;
	commandedPlanet: Readable<CommandedPlanet | undefined>;
	commandedFleet: Readable<CommandedFleet | undefined>;
	commandedMapObject: Readable<MapObject | undefined>;
	commandedMapObjectKey: Readable<string>;
	selectedMapObject: Readable<MapObject | undefined>;
	zoomTarget: Readable<MapObject | undefined>;
	selectedWaypoint: Readable<Waypoint | undefined>;
	highlightedMapObject: Readable<MapObject | undefined>;
	highlightedMapObjectPeers: Readable<MapObject[]>;
	mostRecentMapObject: Readable<MapObject | undefined>;
	currentSelectedWaypointIndex: Readable<number>;

	// scanner updates
	selectMapObject: (mo: MapObject) => void;
	selectNextMapObject: () => void;
	selectWaypoint: (wp: Waypoint) => void;
	commandMapObject: (mo: MapObject) => void;
	commandHomeWorld: () => void;
	previousMapObject: () => void;
	nextMapObject: () => void;
	highlightMapObject: (mo: MapObject | undefined) => void;
	zoomToMapObject: (mo: MapObject) => void;
	nextCommandableMapObjectAtPosition: () => void;

	// message
	gotoTarget: (message: Message, gameId: number, playerNum: number, universe: Universe) => void;

	// game updates
	updateGame: (game: FullGame | GameSettings | Game) => void;
	loadStatus: () => Promise<void>;
	startPollingStatus: (interval?: number) => void;
	stopPollingStatus: () => void;

	// game CRUD
	submitTurn: () => Promise<void>;
	forceGenerateTurn: () => Promise<void>;

	updatePlayerOrders: () => Promise<void>;
	updatePlayerRelations: () => Promise<void>;
	createBattlePlan: (plan: BattlePlan) => Promise<BattlePlan>;
	updateBattlePlan: (plan: BattlePlan) => Promise<BattlePlan>;
	deleteBattlePlan: (num: number) => Promise<void>;
	createProductionPlan: (plan: ProductionPlan) => Promise<ProductionPlan>;
	updateProductionPlan: (plan: ProductionPlan) => Promise<ProductionPlan>;
	deleteProductionPlan: (num: number) => Promise<void>;
	createTransportPlan: (plan: TransportPlan) => Promise<TransportPlan>;
	updateTransportPlan: (plan: TransportPlan) => Promise<TransportPlan>;
	deleteTransportPlan: (num: number) => Promise<void>;

	createDesign: (design: ShipDesign) => Promise<ShipDesign>;
	updateDesign: (design: ShipDesign) => Promise<void>;
	deleteDesign: (num: number) => Promise<void>;
	updateFleetOrders: (fleet: CommandedFleet) => Promise<void>;
	renameFleet: (fleet: CommandedFleet, name: string) => Promise<void>;
	updatePlanetOrders: (planet: CommandedPlanet) => Promise<void>;
	transferCargo: (
		fleet: CommandedFleet,
		dest: Fleet | Planet | Salvage,
		transferAmount: CargoTransferRequest
	) => Promise<void>;
	split: (
		src: CommandedFleet,
		dest: Fleet | undefined,
		srcTokens: ShipToken[],
		destTokens: ShipToken[],
		transferAmount: CargoTransferRequest
	) => Promise<void>;
	splitAll: (fleet: CommandedFleet) => Promise<void>;
	merge: (fleet: CommandedFleet, fleetNums: number[]) => Promise<void>;
	resetContext: (fg: FullGame) => void;
};

// init the game context with empty data
export const getGameContext = () => getContext<GameContext>(gameKey);

// update the game context after a load
export function createGameContext(fg: FullGame): GameContext {
	const gameId = fg.id;
	const unsubscribers: Unsubscriber[] = [];

	const game = writable(fg);
	const player = writable(fg.player);
	const universe = writable(fg.universe);

	const defaultSettings = loadSettingsOrDefault(gameId, fg.player.num);
	const settings = writable(defaultSettings);

	const commandedPlanet = writable<CommandedPlanet | undefined>();
	const commandedFleet = writable<CommandedFleet | undefined>();
	const commandedMapObject = writable<MapObject | undefined>();
	const commandedMapObjectKey = writable<string>(key(undefined));
	const selectedMapObject = writable<MapObject | undefined>();
	const selectedWaypoint = writable<Waypoint | undefined>();
	const highlightedMapObject = writable<MapObject | undefined>();
	const highlightedMapObjectPeers = writable<MapObject[]>([]);
	const mostRecentMapObject = writable<MapObject | undefined>();

	const zoomTarget = writable<MapObject | undefined>();

	const messageNum = writable(
		getNextVisibleMessageNum(-1, false, fg.player.messages, defaultSettings)
	);

	// reset the GameContext for a new game
	// this is called after a new game is loaded from the server while waiting for a turn to generate
	function resetContext(fg: FullGame) {
		const s = get(settings);
		game.set(fg);
		player.set(fg.player);
		universe.set(fg.universe);
		commandedPlanet.set(undefined);
		commandedFleet.set(undefined);
		commandedMapObject.set(undefined);
		selectedMapObject.set(undefined);
		selectedWaypoint.set(undefined);
		highlightedMapObject.set(undefined);
		highlightedMapObjectPeers.set([]);
		mostRecentMapObject.set(undefined);
		messageNum.set(getNextVisibleMessageNum(-1, false, fg.player.messages, s));
	}

	// make sure updates to settings save to localStorage
	unsubscribers.push(
		settings.subscribe((value) => {
			value.beforeSave();
			localStorage.setItem(value.key, JSON.stringify(value));
		})
	);

	// for some use cases (like resetting the selectedWaypointIndex) we want to know when the commanded MapObject
	// changes from one MapObject to another. We don't want to react if the existing MapObject is just updated though
	unsubscribers.push(
		commandedMapObject.subscribe((mo) => {
			const existingKey = get(commandedMapObjectKey);
			const updatedKey = key(mo);
			if (existingKey != updatedKey) {
				commandedMapObjectKey.set(updatedKey);
			}
		})
	);

	function getNextVisibleMessageNum(
		num: number,
		showFilteredMessages: boolean,
		messages: Message[],
		settings: PlayerSettings
	): number {
		for (let i = num + 1; i < messages.length; i++) {
			if (showFilteredMessages || settings.isMessageVisible(messages[i].type)) {
				return i;
			}
		}
		return num;
	}

	// TODO: remove this dep

	const currentCommandedMapObjectIndex = derived(
		[universe, commandedFleet, commandedPlanet, settings],
		([$universe, $commandedFleet, $commandedPlanet, $settings]) => {
			if ($commandedPlanet) {
				return $universe
					.getMyPlanets($settings.sortPlanetsKey, $settings.sortPlanetsDescending)
					.findIndex((p) => p.num === $commandedPlanet.num);
			}
			if ($commandedFleet) {
				return $universe
					.getMyFleets($settings.sortFleetsKey, $settings.sortFleetsDescending)
					.findIndex((f) => f.num === $commandedFleet.num);
			}
			return 0;
		}
	);

	const currentSelectedMapObjectIndex = derived(
		[universe, selectedMapObject],
		([$universe, $selectedMapObject]) => {
			if ($selectedMapObject) {
				const mos = $universe.getMapObjectsByPosition($selectedMapObject.position);
				return findIndex(mos, (mo) => equal($selectedMapObject, mo));
			}
			return -1;
		}
	);

	// derived store of all commandableMapObjects at the position of the current commandedMapObjects
	const commandableMapObjectsAtCommandedMapObjectPosition = derived(
		[universe, commandedMapObject],
		([$universe, $commandedMapObject]) => {
			if ($commandedMapObject) {
				return $universe.getMyMapObjectsByPosition($commandedMapObject);
			}
		}
	);

	// derived store of the current index
	const currentCommandedMapObjectPositionIndex = derived(
		[commandedMapObject, commandableMapObjectsAtCommandedMapObjectPosition],
		([$commandedMapObject, $commandableMapObjectsAtCommandedMapObjectPosition]) => {
			if ($commandedMapObject && $commandableMapObjectsAtCommandedMapObjectPosition) {
				return findIndex($commandableMapObjectsAtCommandedMapObjectPosition, (mo) =>
					equal($commandedMapObject, mo)
				);
			}
			return -1;
		}
	);

	function nextCommandableMapObjectAtPosition() {
		const index = get(currentCommandedMapObjectPositionIndex);
		const commandable = get(commandableMapObjectsAtCommandedMapObjectPosition);

		if (commandable && commandable?.length > 0) {
			if (index + 1 > commandable.length) {
				commandMapObject(commandable[0]);
			} else {
				commandMapObject(commandable[index + 1]);
			}
		}
	}

	// goto a message target
	function gotoTarget(message: Message, gameId: number, playerNum: number, universe: Universe) {
		const targetType = message.targetType ?? MessageTargetType.None;
		const targetTargetType = message.spec.targetType ?? MessageTargetType.None;
		let moType = MapObjectType.None;
		let targetTargetMapObjectType = MapObjectType.None;

		if (message.battleNum) {
			goto(`/games/${gameId}/battles/${message.battleNum}`);
			return;
		}

		if (message.type === MessageType.PlayerGainTechLevel) {
			goto(`/games/${gameId}/research`);
		}

		if (message.type === MessageType.PlayerTechGained && message.spec.techGained) {
			goto(`/games/${gameId}/techs/${kebabCase(message.spec.techGained)}`);
		}

		if (message.targetNum) {
			moType = getMapObjectTypeForMessageType(targetType);
			targetTargetMapObjectType = getMapObjectTypeForMessageType(targetTargetType);

			if (moType != MapObjectType.None) {
				const target = universe.getMapObject(message);
				const targetTarget = universe.getMapObject(message.spec);
				if (target) {
					// if this is a fleet that we own, select the planet before we command the fleet
					if (target.type == MapObjectType.Fleet) {
						if (target.playerNum == playerNum) {
							commandMapObject(target);
							const orbitingPlanetNum = (target as Fleet).orbitingPlanetNum;
							if (orbitingPlanetNum && orbitingPlanetNum != None) {
								const orbiting = universe.getPlanet(orbitingPlanetNum);
								if (orbiting) {
									selectMapObject(orbiting);
								}
							}
						} else {
							selectMapObject(target);
						}
					} else if (target.type == MapObjectType.Planet) {
						if (target.playerNum == playerNum) {
							commandMapObject(target);
							if (targetTarget) {
								selectMapObject(targetTarget);
							} else {
								// select the planet as well if we don't have another target
								// it's weird in the UI to go to a planet that sends a message
								// and see another planet selected
								selectMapObject(target);
							}
						} else {
							selectMapObject(target);
							if (targetTarget && targetTarget.playerNum == playerNum) {
								commandMapObject(targetTarget);
							}
						}
					} else {
						selectMapObject(target);
					}

					// zoom on goto
					zoomToMapObject(target);
					goto(`/games/${gameId}`);
				}
			}
		}
	}

	const currentSelectedWaypointIndex = derived(
		[selectedWaypoint, commandedFleet],
		([$selectedWaypoint, $commandedFleet]) => {
			if ($selectedWaypoint && $commandedFleet) {
				return findIndex($commandedFleet.waypoints, (wp) => wp === $selectedWaypoint);
			}
			return -1;
		}
	);

	function selectNextMapObject() {
		const u = get(universe);
		const selected = get(selectedMapObject);
		const index = get(currentSelectedMapObjectIndex);

		if (index != -1 && selected) {
			const mos = u.getMapObjectsByPosition(selected.position);
			if (mos) {
				if (index >= mos.length - 1) {
					selectMapObject(mos[0]);
				} else {
					selectMapObject(mos[index + 1]);
				}
			}
		}
	}

	// command the previous mapObject for this type, i.e. the previous planet or fleet
	function previousMapObject() {
		const u = get(universe);
		const i = get(currentCommandedMapObjectIndex);
		const mo = get(commandedMapObject);
		const s = get(settings);

		if (mo) {
			if (mo.type == MapObjectType.Planet) {
				const planets = u.getMyPlanets(s.sortPlanetsKey, s.sortPlanetsDescending);
				const prevIndex = rollover(i - 1, 0, planets.length - 1);
				const planet = planets[prevIndex];
				commandMapObject(planet);
				zoomToMapObject(planet);
				selectMapObject(planet);
			} else if (mo.type == MapObjectType.Fleet) {
				const fleets = u.getMyFleets(s.sortFleetsKey, s.sortFleetsDescending);
				const prevIndex = rollover(i - 1, 0, fleets.length - 1);
				commandMapObject(fleets[prevIndex]);
				zoomToMapObject(fleets[prevIndex]);

				const fleet = fleets[prevIndex];
				if (fleet.orbitingPlanetNum && fleet.orbitingPlanetNum != None) {
					const planet = u.getMapObject({
						targetType: MapObjectType.Planet,
						targetNum: fleet.orbitingPlanetNum
					});
					if (planet) {
						selectMapObject(planet);
					}
				} else {
					selectMapObject(fleet);
				}
			}
		}
	}

	// command the next mapObject for this type, i.e. the next planet or fleet
	function nextMapObject() {
		const u = get(universe);
		const i = get(currentCommandedMapObjectIndex);
		const mo = get(commandedMapObject);
		const s = get(settings);

		if (mo) {
			if (mo.type == MapObjectType.Planet) {
				const planets = u.getMyPlanets(s.sortPlanetsKey, s.sortPlanetsDescending);
				const nextIndex = rollover(i + 1, 0, planets.length - 1);
				const planet = planets[nextIndex];
				commandMapObject(planet);
				zoomToMapObject(planet);
				selectMapObject(planet);
			} else if (mo.type == MapObjectType.Fleet) {
				const fleets = u.getMyFleets(s.sortFleetsKey, s.sortFleetsDescending);

				const nextIndex = rollover(i + 1, 0, fleets.length - 1);
				const fleet = fleets[nextIndex];
				commandMapObject(fleets[nextIndex]);
				zoomToMapObject(fleets[nextIndex]);
				if (fleet.orbitingPlanetNum && fleet.orbitingPlanetNum != None) {
					const planet = u.getMapObject({
						targetType: MapObjectType.Planet,
						targetNum: fleet.orbitingPlanetNum
					});
					if (planet) {
						selectMapObject(planet);
					}
				} else {
					selectMapObject(fleet);
				}
			}
		}
	}

	function selectMapObject(mo: MapObject) {
		selectedMapObject.update(() => mo);
		mostRecentMapObject.update(() => mo);
	}

	function selectWaypoint(wp: Waypoint) {
		selectedWaypoint.update(() => wp);
	}

	function commandMapObject(mo: MapObject) {
		commandedMapObject.update(() => mo);
		mostRecentMapObject.update(() => mo);
		if (mo.type == MapObjectType.Planet) {
			commandedPlanet.update(() => Object.assign(new CommandedPlanet(), mo));
			commandedFleet.update(() => undefined);
		} else if (mo.type == MapObjectType.Fleet) {
			commandedFleet.update(() => Object.assign(new CommandedFleet(), mo));
			commandedPlanet.update(() => undefined);
			selectedWaypoint.update(() => {
				const fleet = mo as Fleet;
				if (fleet?.waypoints && fleet.waypoints.length) {
					return fleet.waypoints[0];
				}
				return undefined;
			});
		}
	}

	// command the player's homeworld (or the first planet they own, if their homeworld has been taken)
	function commandHomeWorld() {
		const u = get(universe);
		const s = get(settings);
		const homeworld = u.getHomeworld();
		if (homeworld) {
			commandMapObject(homeworld);
			selectMapObject(homeworld);
			zoomToMapObject(homeworld);
		} else {
			// command our first planet
			const planets = u.getMyPlanets(s.sortPlanetsKey, s.sortPlanetsDescending);
			if (planets.length > 0) {
				commandMapObject(planets[0]);
				selectMapObject(planets[0]);
				zoomToMapObject(planets[0]);
			}
		}
	}

	function highlightMapObject(mo: MapObject | undefined) {
		highlightedMapObject.update(() => mo);
	}

	function zoomToMapObject(mo: MapObject) {
		zoomTarget.update(() => mo);
		mostRecentMapObject.update(() => mo);
	}

	// update the game state from a server
	function updateGame(g: FullGame | GameSettings | Game) {
		game.set(Object.assign(get(game), g));
	}

	function updatePlayer(p: Player | PlayerResponse | undefined) {
		player.set(Object.assign(get(player), p));
	}

	// after a fleet is updated from the server, update the fleet in the universe, reset any commanded/selected
	// state and trigger reactivity
	function updateFleet(fleet: CommandedFleet | Fleet, updatedFleet: CommandedFleet | Fleet) {
		fleet = Object.assign(fleet, updatedFleet);
		const index = get(currentSelectedWaypointIndex);
		const u = get(universe);
		const cf = get(commandedFleet);

		// update the fleet in the universe
		u.updateFleet(fleet);

		// if we were commanding this fleet, recommand it to trigger reactivity
		if (equal(cf, fleet)) {
			commandMapObject(fleet);
		}

		// if we were selecting this fleet, reselect it to trigger reactivity
		if (equal(get(selectedMapObject), fleet)) {
			selectMapObject(fleet);
		}

		if (index > -1 && fleet.waypoints && fleet.waypoints.length > index) {
			selectWaypoint(fleet.waypoints[index]);
		}

		// trigger reactivity
		universe.set(u);
	}

	// after a planet is updated from the server, update the planet in the universe, reset any commanded/selected
	// state and trigger reactivity
	function updatePlanet(planet: CommandedPlanet | Planet, updatedPlanet: CommandedPlanet | Planet) {
		planet = Object.assign(planet, updatedPlanet);
		const u = get(universe);
		const cp = get(commandedPlanet);
		u.updatePlanet(planet);

		// if we were commanding this planet, recommand it to trigger reactivity
		if (cp?.num === planet.num) {
			commandMapObject(planet);
		}

		// if we were selecting this planet, reselect it to trigger reactivity
		if (equal(get(selectedMapObject), planet)) {
			selectMapObject(planet);
		}

		// trigger reactivity
		universe.set(u);
	}

	// load the status of a game, but not all the universe data
	async function loadStatus(): Promise<void> {
		const result = await GameService.loadGame(gameId);
		updateGame(result);
	}

	// start polling the server for player status
	let playerStatusPollingInterval: number | undefined = undefined;
	function startPollingStatus(interval = 10000) {
		if (!playerStatusPollingInterval) {
			playerStatusPollingInterval = window.setInterval(async () => {
				await loadStatus();
			}, interval);
		}
	}

	// stop polling the server for player status
	function stopPollingStatus() {
		if (playerStatusPollingInterval) {
			window.clearInterval(playerStatusPollingInterval);
			playerStatusPollingInterval = undefined;
		}
	}

	// game CRUD

	async function submitTurn(): Promise<void> {
		const result = await PlayerService.submitTurn(gameId);
		if (result) {
			updateGame(result.game);
			updatePlayer(result.player);
			if (result.universe) {
				universe.set(get(universe).resetData(get(player).num, result.universe));
			}
		}
	}

	async function forceGenerateTurn(): Promise<void> {
		const result = await GameService.forceGenerateTurn(gameId);
		updateGame(result.game);
		updatePlayer(result.player);
		if (result.universe) {
			universe.set(get(universe).resetData(get(player).num, result.universe));
		}
	}

	async function updatePlayerOrders(): Promise<void> {
		const result = await PlayerService.updateOrders(get(player));
		if (result) {
			player.set(Object.assign(get(player), result.player));

			const u = get(universe);
			result.planets.forEach((planet) => {
				u.planets[planet.num - 1] = planet;
				if (equal(get(selectedMapObject), planet)) {
					selectMapObject(planet);
				}
			});

			u.resetMapObjectsByPosition();
			u.resetMyMapObjectsByPosition();

			universe.set(u);
		}
	}

	async function updatePlayerRelations(): Promise<void> {
		const result = await PlayerService.updateRelations(get(player));
		if (result) {
			updatePlayer(result);
		}
	}

	async function createBattlePlan(plan: BattlePlan): Promise<BattlePlan> {
		const created = await BattlePlanService.create(gameId, plan);
		const p = get(player);
		p.battlePlans = [...p.battlePlans, created];
		player.set(p);

		return created;
	}

	async function updateBattlePlan(plan: BattlePlan): Promise<BattlePlan> {
		return await BattlePlanService.update(gameId, plan);
	}

	async function deleteBattlePlan(num: number): Promise<void> {
		const resp = await BattlePlanService.delete(gameId, num);
		updatePlayer(resp.player);
		const u = get(universe);
		u.fleets = resp.fleets;
		u.starbases = resp.starbases;
		universe.set(u);
	}

	async function createProductionPlan(plan: ProductionPlan): Promise<ProductionPlan> {
		const created = await ProductionPlanService.create(gameId, plan);
		const p = get(player);
		p.productionPlans = [...p.productionPlans, created];
		player.set(p);

		return created;
	}

	async function updateProductionPlan(plan: ProductionPlan): Promise<ProductionPlan> {
		return await ProductionPlanService.update(gameId, plan);
	}

	async function deleteProductionPlan(num: number): Promise<void> {
		const player = await ProductionPlanService.delete(gameId, num);
		Object.assign(player, player);
	}

	async function createTransportPlan(plan: TransportPlan): Promise<TransportPlan> {
		const created = await TransportPlanService.create(gameId, plan);
		const p = get(player);
		p.transportPlans = [...p.transportPlans, created];
		player.set(p);

		return created;
	}

	async function updateTransportPlan(plan: TransportPlan): Promise<TransportPlan> {
		return await TransportPlanService.update(gameId, plan);
	}

	async function deleteTransportPlan(num: number): Promise<void> {
		const player = await TransportPlanService.delete(gameId, num);
		Object.assign(player, player);
	}

	async function createDesign(design: ShipDesign): Promise<ShipDesign> {
		// update this design
		design = await DesignService.create(gameId, design);
		universe.set(get(universe).addDesign(design));
		return design;
	}

	async function updateDesign(design: ShipDesign): Promise<void> {
		// update this design
		design = await DesignService.update(gameId, design);
		universe.set(get(universe).updateDesign(design));
	}

	async function deleteDesign(num: number): Promise<void> {
		const { fleets, starbases } = await DesignService.delete(gameId, num);
		const u = get(universe);
		u.fleets = fleets;
		u.starbases = starbases;
		u.resetMapObjectsByPosition();
		u.resetMyMapObjectsByPosition();

		u.designs = u.designs.filter((d) => d.num != num);
		universe.set(u);

		// reset our view to the homeworld, in case the commanded fleet had our deleted design
		commandHomeWorld();
	}

	async function updateFleetOrders(fleet: CommandedFleet): Promise<void> {
		const updatedFleet = await FleetService.updateFleetOrders(fleet);
		updateFleet(fleet, updatedFleet);
	}

	async function renameFleet(fleet: CommandedFleet, name: string): Promise<void> {
		const updatedFleet = await FleetService.rename(fleet, name);
		updateFleet(fleet, updatedFleet);
	}

	async function updatePlanetOrders(planet: CommandedPlanet): Promise<void> {
		const resp = await PlanetService.updatePlanetOrders(planet);

		// changing the planet orders changes the player's spec
		updatePlayer(resp.player);
		updatePlanet(planet, resp.planet);
	}

	async function transferCargo(
		fleet: CommandedFleet,
		dest: Fleet | Planet | Salvage,
		transferAmount: CargoTransferRequest
	): Promise<void> {
		const result = await FleetService.transferCargo(fleet, dest, transferAmount);

		if (result.dest?.type == MapObjectType.Planet) {
			const planet = result.dest as Planet;
			updatePlanet(dest as Planet, planet);
		} else if (result.dest?.type == MapObjectType.Fleet) {
			// update the destination fleet in the universe
			const destFleet = result.dest as Fleet;
			updateFleet(dest as Fleet, destFleet);
		}

		if (result.salvages) {
			get(universe).updateSalvages(result.salvages);
		}

		updateFleet(fleet, result.fleet);
	}

	async function split(
		src: CommandedFleet,
		dest: Fleet | undefined,
		srcTokens: ShipToken[],
		destTokens: ShipToken[],
		transferAmount: CargoTransferRequest
	): Promise<void> {
		const response = await FleetService.split(src, dest, srcTokens, destTokens, transferAmount);

		const u = get(universe);

		// if the original source is different than the returned source
		// and we have no response.dest, this means the original source was deleted
		// and has become the source, i.e. we moved all tokens from source to dest, creating
		// a new fleet. Weird edge case.
		if (src.num != response.source.num) {
			u.removeFleets([src.num]);
		}

		// update the commanded fleet
		const source = Object.assign(new CommandedFleet(), response.source);
		u.updateFleet(response.source);
		commandMapObject(source);
		if (equal(get(selectedMapObject), source)) {
			selectMapObject(source);
		}

		// update or add the new fleets to the universe
		if (response.dest) {
			if (dest?.num == 0) {
				u.addFleets([response.dest]);
			} else {
				u.updateFleet(response.dest);
			}
		} else {
			// if we had a dest and it was deleted, remove it
			dest && u.removeFleets([dest.num]);
		}

		const index = get(currentSelectedWaypointIndex);
		if (index > -1 && source.waypoints && source.waypoints.length > index) {
			selectWaypoint(source.waypoints[index]);
		}
	}

	async function splitAll(fleet: CommandedFleet): Promise<void> {
		const updatedFleets = await FleetService.splitAll(fleet.gameId, fleet);
		const sourceFleet = updatedFleets.find((f) => f.num == fleet.num);
		if (sourceFleet) {
			fleet = Object.assign(fleet, sourceFleet);
			commandMapObject(fleet);
		}

		const u = get(universe);

		// update and add the new fleets to the universe
		u.updateFleet(fleet);
		u.addFleets(updatedFleets.filter((f) => f.num != fleet.num));

		const index = get(currentSelectedWaypointIndex);
		if (index > -1 && fleet.waypoints && fleet.waypoints.length > index) {
			selectWaypoint(fleet.waypoints[index]);
		}
	}

	async function merge(fleet: CommandedFleet, fleetNums: number[]): Promise<void> {
		const updatedFleet = await FleetService.merge(fleet, fleetNums);

		get(universe).removeFleets(fleetNums);
		updateFleet(fleet, updatedFleet);
	}

	return {
		game,
		player,
		universe,
		settings,
		messageNum,
		commandedPlanet,
		commandedFleet,
		commandedMapObject,
		commandedMapObjectKey,
		selectedMapObject,
		zoomTarget,
		selectedWaypoint,
		highlightedMapObject,
		highlightedMapObjectPeers,
		mostRecentMapObject,
		currentSelectedWaypointIndex,

		selectMapObject,
		selectNextMapObject,
		selectWaypoint,
		commandMapObject,
		commandHomeWorld,
		previousMapObject,
		nextMapObject,
		highlightMapObject,
		zoomToMapObject,
		nextCommandableMapObjectAtPosition,
		gotoTarget,

		updateGame,
		loadStatus,
		startPollingStatus,
		stopPollingStatus,

		submitTurn,
		forceGenerateTurn,

		updatePlayerOrders,
		updatePlayerRelations,
		createBattlePlan,
		updateBattlePlan,
		deleteBattlePlan,
		createProductionPlan,
		updateProductionPlan,
		deleteProductionPlan,
		createTransportPlan,
		updateTransportPlan,
		deleteTransportPlan,

		createDesign,
		updateDesign,
		deleteDesign,
		updateFleetOrders,
		renameFleet,
		updatePlanetOrders,
		transferCargo,
		split,
		splitAll,
		merge,
		resetContext
	};
}

function loadSettingsOrDefault(gameId: number, playerNum: number): PlayerSettings {
	const key = PlayerSettings.key(gameId, playerNum);

	const json = localStorage.getItem(key);
	if (json) {
		const settingsJSON = JSON.parse(json) as PlayerSettings;
		if (settingsJSON) {
			// create a new object
			const settings = new PlayerSettings(gameId, playerNum);
			Object.assign(settings, settingsJSON);
			settings.afterLoad();
			return settings;
		}
	}

	return new PlayerSettings(gameId, playerNum);
}
