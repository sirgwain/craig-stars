import { goto } from '$app/navigation';
import { getScannerTarget } from '$lib/types/Battle';
import type { CargoTransferRequest } from '$lib/types/CargoTransferRequest.svelte';
import { type CargoDest } from '$lib/types/CargoTransferRequest.svelte';
import { None } from '$lib/types/Consts';
import {
	CargoSchema,
	FleetOrdersSchema,
	FleetSchema,
	MapObjectType,
	MineralPacketSchema,
	PlanetSchema,
	PlayerMessageTargetType,
	PlayerMessageType,
	SalvageIntelSchema,
	type BattlePlan,
	type CargoTransfers,
	type Fleet,
	type Game,
	type GameSettings,
	type GameWithPlayers,
	type Minefield,
	type MineralPacket,
	type Planet,
	type Player,
	type PlayerMessage,
	type PlayerRelationship,
	type ProductionPlan,
	type SalvageIntel,
	type ShipDesign,
	type ShipTokenJson as ShipToken,
	type TransportPlan,
	type Waypoint,
	type WaypointDest
} from '$lib/types/cs-proto';
import { CommandedFleet } from '$lib/types/Fleet';
import { getGameWithPlayersFlat, type GameWithPlayersFlat } from '$lib/types/Game';
import { equal, key, ownedBy, type MapObjectLike } from '$lib/types/MapObject';
import { getMapObjectTarget, getMapObjectTypeForMessageType } from '$lib/types/Message';
import { CommandedPlanet } from '$lib/types/Planet';
import { CommandedPlayer } from '$lib/types/Player';
import { PlayerSettings } from '$lib/types/PlayerSettings';
import type { CS } from '$lib/wasm';
import { create } from '@bufbuild/protobuf';
import { findIndex, kebabCase } from 'lodash-es';
import { getContext } from 'svelte';
import {
	derived,
	get,
	writable,
	type Readable,
	type Unsubscriber,
	type Writable
} from 'svelte/store';
import {
	battlePlanClient,
	fleetClient,
	gameClient,
	minefieldClient,
	planetClient,
	playerClient,
	productionPlanClient,
	shipDesignClient,
	transportPlanClient
} from './connect';
import { FullGame } from './FullGame';
import { rollover } from './Math';
import { Universe, type AnyFleet } from './Universe';

export const playerFinderKey = Symbol();
export const designFinderKey = Symbol();
export const gameKey = Symbol();

export type GameContext = {
	fullyLoaded: Readable<boolean>;
	cs: CS;
	game: Readable<FullGame>;
	player: Readable<CommandedPlayer>;
	universe: Readable<Universe>;
	settings: Writable<PlayerSettings>;
	messageNum: Writable<number>;
	commandedPlanet: Readable<CommandedPlanet | undefined>;
	commandedFleet: Readable<CommandedFleet | undefined>;
	commandedMapObject: Readable<MapObjectLike | undefined>;
	commandedMapObjectKey: Readable<string>;
	selectedMapObject: Readable<MapObjectLike | undefined>;
	zoomTarget: Readable<MapObjectLike | undefined>;
	selectedWaypoint: Readable<Waypoint | undefined>;
	highlightedMapObject: Readable<MapObjectLike | undefined>;
	highlightedMapObjectPeers: Readable<MapObjectLike[]>;
	mostRecentMapObject: Readable<MapObjectLike | undefined>;
	currentSelectedWaypointIndex: Readable<number>;

	// scanner updates
	selectMapObject: (mo: MapObjectLike) => void;
	selectNextMapObject: () => void;
	selectWaypoint: (wp: Waypoint) => void;
	commandMapObject: (mo: MapObjectLike) => void;
	commandHomeWorld: () => void;
	previousMapObject: () => void;
	nextMapObject: () => void;
	highlightMapObject: (mo: MapObjectLike | undefined) => void;
	zoomToMapObject: (mo: MapObjectLike) => void;

	// message
	gotoTarget: (
		message: PlayerMessage,
		gameId: bigint,
		playerNum: number,
		universe: Universe
	) => void;
	gotoBattle: (battleNum: number) => void;

	// game updates
	updateGame: (g: FullGame | Game | GameSettings | GameWithPlayersFlat | GameWithPlayers) => void;
	loadStatus: () => Promise<void>;
	startPollingStatus: (interval?: number) => void;
	stopPollingStatus: () => void;

	// game CRUD
	submitTurn: () => Promise<void>;
	forceGenerateTurn: () => Promise<void>;

	updatePlayerOrders: () => Promise<void>;
	updatePlayerRelationships: (relations: PlayerRelationship[]) => Promise<void>;
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
	updateDesign: (design: ShipDesign) => Promise<ShipDesign>;
	deleteDesign: (num: number) => Promise<void>;

	// fleet waypoint updates
	addWaypoint: (dest: WaypointDest, fastestWaypoint: boolean) => Promise<boolean>;
	updateWaypoint: (dest: WaypointDest, fastestWaypoint: boolean, done: boolean) => Promise<void>;
	deleteWaypoint: () => Promise<void>;

	updateFleetOrders: (fleet: CommandedFleet) => Promise<void>;
	renameFleet: (fleet: CommandedFleet, name: string) => Promise<void>;

	updatePlanetOrders: (planet: CommandedPlanet) => Promise<void>;
	updateMinefieldOrders: (minefield: Minefield) => Promise<void>;
	transferCargo: (
		fleet: CommandedFleet,
		dest: CargoDest,
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
	setFullyLoaded: (fullyLoaded: boolean) => void;
	resetContext: (fg: FullGame, p: CommandedPlayer, u: Universe) => Promise<void>;
};

// init the game context with empty data
export const getGameContext = () => getContext<GameContext>(gameKey);

// update the game context after a load
export async function createGameContext(
	cs: CS,
	fg: FullGame,
	p: CommandedPlayer,
	u: Universe
): Promise<GameContext> {
	// setup initial wasm state
	const { spec: raceSpec } = await cs.wasmService.computeRaceSpec({ race: p.race });
	p.race.spec = raceSpec ?? p.race.spec;
	await cs.wasmService.setPlayer({ player: p });
	await cs.wasmService.setDesigns({ designs: u.designs });
	await cs.wasmService.setIntels({ intels: u.intels });

	const gameId = fg.id;
	const unsubscribers: Unsubscriber[] = [];

	const fullyLoaded = writable(false);

	const game = writable(fg);
	const player = writable(p);
	const universe = writable(u);

	const defaultSettings = loadSettingsOrDefault(gameId, p.num);
	const settings = writable(defaultSettings);

	const commandedPlanet = writable<CommandedPlanet | undefined>();
	const commandedFleet = writable<CommandedFleet | undefined>();
	const commandedMapObject = writable<MapObjectLike | undefined>();
	const commandedMapObjectKey = writable<string>(key(undefined));
	const selectedMapObject = writable<MapObjectLike | undefined>();
	const selectedWaypoint = writable<Waypoint | undefined>();
	const highlightedMapObject = writable<MapObjectLike | undefined>();
	const highlightedMapObjectPeers = writable<MapObjectLike[]>([]);
	const mostRecentMapObject = writable<MapObjectLike | undefined>();

	const zoomTarget = writable<MapObjectLike | undefined>();

	const messageNum = writable(getNextVisibleMessageNum(-1, false, p.messages, defaultSettings));

	// reset the GameContext for a new game
	// this is called after a new game is loaded from the server while waiting for a turn to generate
	async function resetContext(fg: FullGame, p: CommandedPlayer, u: Universe) {
		const { spec: raceSpec } = await cs.wasmService.computeRaceSpec({ race: p.race });
		p.race.spec = raceSpec ?? p.race.spec;
		await cs.wasmService.setPlayer({ player: p });

		game.set(fg);
		await updatePlayer(p);
		universe.set(u);
		commandedPlanet.set(undefined);
		commandedFleet.set(undefined);
		commandedMapObject.set(undefined);
		selectedMapObject.set(undefined);
		selectedWaypoint.set(undefined);
		highlightedMapObject.set(undefined);
		highlightedMapObjectPeers.set([]);
		mostRecentMapObject.set(undefined);

		const s = get(settings);
		messageNum.set(getNextVisibleMessageNum(-1, false, p.messages, s));
	}

	function setFullyLoaded(value: boolean) {
		fullyLoaded.update(() => value);
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
		messages: PlayerMessage[],
		settings: PlayerSettings
	): number {
		for (let i = num + 1; i < messages.length; i++) {
			if (showFilteredMessages || settings.isMessageVisible(messages[i].type)) {
				return i;
			}
		}
		return num;
	}

	const currentCommandedMapObjectIndex = derived(
		[universe, commandedFleet, commandedPlanet, settings],
		([$universe, $commandedFleet, $commandedPlanet, $settings]) => {
			if ($commandedPlanet) {
				return $universe
					.getMyPlanets($settings.sortPlanetsKey, $settings.sortPlanetsDescending)
					.findIndex((p) => (p.mapObject?.num ?? 0) === ($commandedPlanet.mapObject?.num ?? 0));
			}
			if ($commandedFleet) {
				return $universe
					.getMyFleets($settings.sortFleetsKey, $settings.sortFleetsDescending)
					.findIndex((f) => (f.mapObject?.num ?? 0) === ($commandedFleet.mapObject?.num ?? 0));
			}
			return 0;
		}
	);

	const currentSelectedMapObjectIndex = derived(
		[universe, selectedMapObject],
		([$universe, $selectedMapObject]) => {
			if ($selectedMapObject) {
				const mos = $universe.getMapObjectsByPosition($selectedMapObject.mapObject?.position);
				return findIndex(mos, (mo) => equal($selectedMapObject, mo));
			}
			return -1;
		}
	);

	// goto a message target
	function gotoTarget(
		message: PlayerMessage,
		gameId: bigint,
		playerNum: number,
		universe: Universe
	) {
		const targetType: PlayerMessageTargetType =
			message.target?.targetType ?? PlayerMessageTargetType.UNSPECIFIED;
		let moType = MapObjectType.UNSPECIFIED;

		if (message.battleNum) {
			goto(`/games/${gameId}/battles/${message.battleNum}`);
			return;
		}

		if (message.type === PlayerMessageType.PLAYER_GAIN_TECH_LEVEL) {
			goto(`/games/${gameId}/research`);
			return;
		}

		if (message.type === PlayerMessageType.BATTLE_REPORTS) {
			goto(`/games/${gameId}/battles`);
			return;
		}

		if (message.type === PlayerMessageType.PLAYER_TECH_GAINED && message.spec?.techGained) {
			goto(`/games/${gameId}/techs/${kebabCase(message.spec?.techGained)}`);
			return;
		}

		if (
			message.type === PlayerMessageType.MYSTERY_TRADER_MET_WITH_REWARD &&
			message.spec?.mysteryTrader?.tech
		) {
			goto(`/games/${gameId}/techs/${kebabCase(message.spec?.mysteryTrader?.tech)}`);
			return;
		}

		// the MT gave us a fleet, command it
		if (
			message.type === PlayerMessageType.MYSTERY_TRADER_MET_WITH_REWARD &&
			message.spec?.mysteryTrader?.fleetNum
		) {
			const fleet = universe.getFleet(playerNum, message.spec?.mysteryTrader?.fleetNum);
			if (fleet) {
				gotoTargetFleet(fleet, playerNum, universe);
				zoomToMapObject(fleet);
				goto(`/games/${gameId}`);
				return;
			}
		}

		if (message.spec?.target?.targetType === MapObjectType.MINEFIELD) {
			const fleet = universe.getFleet(message.target?.targetPlayerNum, message.target?.targetNum);
			const mf = universe.getMinefield(
				message.spec?.target?.targetPlayerNum,
				message.spec?.target?.targetNum
			);
			if (fleet) {
				if (ownedBy(fleet, playerNum)) {
					commandMapObject(fleet);
				} else {
					selectMapObject(fleet);
					zoomToMapObject(fleet);
				}
				goto(`/games/${gameId}`);
				return;
			}
			if (mf) {
				selectMapObject(mf);
				zoomToMapObject(mf);
				goto(`/games/${gameId}`);
				return;
			}
		}

		if (message.target?.targetNum) {
			moType = getMapObjectTypeForMessageType(targetType);

			if (moType != MapObjectType.UNSPECIFIED) {
				const target = universe.getMapObject(getMapObjectTarget(message));
				const targetTarget = universe.getMapObject(message.spec?.target);
				if (target) {
					// if this is a fleet that we own, select the planet before we command the fleet
					if (target.mapObject?.type === MapObjectType.FLEET) {
						gotoTargetFleet(target, playerNum, universe);
					} else if (target.mapObject?.type === MapObjectType.PLANET) {
						gotoTargetPlanet(target, targetTarget, playerNum);
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

	// command/select logic for a goto of a fleet
	// will select the planet the fleet is orbiting if it has one
	function gotoTargetFleet(target: MapObjectLike, playerNum: number, universe: Universe) {
		if (target.mapObject?.playerNum == playerNum) {
			commandMapObject(target);
			const orbitingPlanetNum = (target as AnyFleet).orbitingPlanetNum;
			if (orbitingPlanetNum && orbitingPlanetNum != None) {
				const orbiting = universe.getPlanet(orbitingPlanetNum);
				if (orbiting) {
					selectMapObject(orbiting);
				}
			}
		} else {
			selectMapObject(target);
		}
	}

	// command/select logic for a goto of a planet
	// will select an additional target if specified, otherwise selects the planet
	function gotoTargetPlanet(
		target: MapObjectLike,
		targetTarget: MapObjectLike | undefined,
		playerNum: number
	) {
		if (target.mapObject?.playerNum == playerNum) {
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
			if (targetTarget && targetTarget.mapObject?.playerNum == playerNum) {
				commandMapObject(targetTarget);
			}
		}
	}

	// select a battle location in the scanner and make the battle message the current message
	function gotoBattle(battleNum: number) {
		const u = get(universe);
		const p = get(player);
		const battle = u.getBattle(battleNum);
		if (!battle) {
			return;
		}

		const target = getScannerTarget(battle, u);
		if (target) {
			selectMapObject(target);
			zoomToMapObject(target);
		} else {
			zoomToMapObject({ position: battle.position } as MapObjectLike);
		}

		const battleMessageNum = p.messages.findIndex((m) => m.battleNum == battle.num);
		if (battleMessageNum) {
			messageNum.set(battleMessageNum);
		}
	}

	const currentSelectedWaypointIndex = derived(
		[selectedWaypoint, commandedFleet],
		([$selectedWaypoint, $commandedFleet]) => {
			if ($selectedWaypoint && $commandedFleet) {
				return findIndex($commandedFleet.fleetOrders?.waypoints, (wp) => wp === $selectedWaypoint);
			}
			return -1;
		}
	);

	function selectNextMapObject() {
		const u = get(universe);
		const selected = get(selectedMapObject);
		const index = get(currentSelectedMapObjectIndex);

		if (index != -1 && selected) {
			const mos = u.getMapObjectsByPosition(selected.mapObject?.position);
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
			if (mo.mapObject?.type === MapObjectType.PLANET) {
				const planets = u.getMyPlanets(s.sortPlanetsKey, s.sortPlanetsDescending);
				const prevIndex = rollover(i - 1, planets.length - 1);
				const planet = planets[prevIndex];
				commandMapObject(planet);
				zoomToMapObject(planet);
				selectMapObject(planet);
			} else if (mo.mapObject?.type === MapObjectType.FLEET) {
				const fleets = u.getMyFleets(s.sortFleetsKey, s.sortFleetsDescending);
				const prevIndex = rollover(i - 1, fleets.length - 1);
				commandMapObject(fleets[prevIndex]);
				zoomToMapObject(fleets[prevIndex]);

				const fleet = fleets[prevIndex];
				if (fleet.orbitingPlanetNum && fleet.orbitingPlanetNum !== None) {
					const planet = u.getMapObject({
						targetType: MapObjectType.PLANET,
						targetNum: fleet.orbitingPlanetNum,
						targetPosition: fleet.mapObject?.position
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
			if (mo.mapObject?.type === MapObjectType.PLANET) {
				const planets = u.getMyPlanets(s.sortPlanetsKey, s.sortPlanetsDescending);
				const nextIndex = rollover(i + 1, planets.length - 1);
				const planet = planets[nextIndex];
				commandMapObject(planet);
				zoomToMapObject(planet);
				selectMapObject(planet);
			} else if (mo.mapObject?.type === MapObjectType.FLEET) {
				const fleets = u.getMyFleets(s.sortFleetsKey, s.sortFleetsDescending);

				const nextIndex = rollover(i + 1, fleets.length - 1);
				const fleet = fleets[nextIndex];
				commandMapObject(fleets[nextIndex]);
				zoomToMapObject(fleets[nextIndex]);
				if (fleet.orbitingPlanetNum && fleet.orbitingPlanetNum != None) {
					const planet = u.getMapObject({
						targetType: MapObjectType.PLANET,
						targetNum: fleet.orbitingPlanetNum,
						targetPosition: fleet.mapObject?.position
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

	function selectMapObject(mo: MapObjectLike) {
		selectedMapObject.update(() => mo);
		mostRecentMapObject.update(() => mo);
	}

	function selectWaypoint(wp: Waypoint) {
		selectedWaypoint.update(() => wp);
	}

	function commandMapObject(mo: MapObjectLike) {
		commandedMapObject.update(() => mo);
		mostRecentMapObject.update(() => mo);
		if (mo.mapObject?.type === MapObjectType.PLANET) {
			const u = get(universe);
			// make sure this planet's production queue estimates are up to date
			const planet = new CommandedPlanet(u.getPlanet(mo.mapObject?.num) as Planet);
			// ensure CommandedPlanet has mapObject when commanded
			// TODO: await this, make async infect everything
			planet.updateProductionQueueEstimates(cs);
			commandedPlanet.update(() => planet);
			commandedFleet.update(() => undefined);
		} else if (mo.mapObject?.type === MapObjectType.FLEET) {
			// Build a CommandedFleet from the selected fleet, preserving its mapObject
			const cf = new CommandedFleet(mo as Fleet);
			commandedFleet.update(() => cf);
			commandedPlanet.update(() => undefined);
			selectedWaypoint.update(() => {
				const fleet = mo as Fleet;
				if (fleet?.fleetOrders?.waypoints && fleet.fleetOrders.waypoints.length) {
					return fleet.fleetOrders.waypoints[fleet.fleetOrders.waypoints.length - 1];
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
			const fleets = u.getMyFleets(s.sortFleetsKey, s.sortFleetsDescending);
			if (planets.length > 0) {
				commandMapObject(planets[0]);
				selectMapObject(planets[0]);
				zoomToMapObject(planets[0]);
			} else if (fleets.length > 0) {
				commandMapObject(fleets[0]);
				selectMapObject(fleets[0]);
				zoomToMapObject(fleets[0]);
			} else {
				const planet = u.getPlanet(1);
				if (planet) {
					selectMapObject(planet);
					zoomToMapObject(planet);
				}
			}
		}
	}

	function highlightMapObject(mo: MapObjectLike | undefined) {
		highlightedMapObject.update(() => mo);
	}

	function zoomToMapObject(mo: MapObjectLike) {
		zoomTarget.update(() => mo);
		mostRecentMapObject.update(() => mo);
	}

	// update the game state from a server
	function updateGame(g: FullGame | Game | GameSettings | GameWithPlayersFlat | GameWithPlayers) {
		let updated = g;
		if ('game' in updated) {
			// flatten the GameWithPlayers
			updated = getGameWithPlayersFlat(g as GameWithPlayers);
		}

		game.set(Object.assign(get(game), updated));
	}

	async function updatePlayer(p: CommandedPlayer | Player | undefined) {
		if (!p) {
			return;
		}
		// new player, recompute spec
		const updated = Object.assign(get(player), p);
		const { spec } = await cs.wasmService.computeRaceSpec({ race: updated.race });
		updated.race.spec = spec!;
		player.set(updated);
		await cs.wasmService.setPlayer({ player: p });
	}

	// after a fleet is updated from the server, update the fleet in the universe, reset any commanded/selected
	// state and trigger reactivity
	function updateFleet(
		fleet: CommandedFleet | Fleet,
		updatedFleet: CommandedFleet | Fleet | undefined
	) {
		if (!updatedFleet) {
			return;
		}
		fleet = new CommandedFleet(Object.assign(fleet, updatedFleet));
		const index = get(currentSelectedWaypointIndex);
		const u = get(universe);
		const cf = get(commandedFleet);

		// update the fleet in the universe
		u.updateFleet(fleet);

		// if we were commanding this fleet, recommand it to trigger reactivity
		if (equal(cf, fleet)) {
			commandMapObject(fleet);
			// update the selected waypoint if we just updated the commandedFleet
			if (
				index > -1 &&
				fleet.fleetOrders?.waypoints &&
				fleet.fleetOrders.waypoints.length > index
			) {
				selectWaypoint(fleet.fleetOrders!.waypoints[index]);
			}
		}

		// if we were selecting this fleet, reselect it to trigger reactivity
		if (equal(get(selectedMapObject), fleet)) {
			selectMapObject(fleet);
		}

		// trigger reactivity
		universe.set(u);
	}

	// after a planet is updated from the server, update the planet in the universe, reset any commanded/selected
	// state and trigger reactivity
	function updatePlanet(planet: CommandedPlanet | Planet, updatedPlanet: Planet) {
		planet = Object.assign(planet, updatedPlanet);
		const u = get(universe);
		u.updatePlanet(planet);

		// if we were commanding this planet, recommand it to trigger reactivity
		if (equal(get(commandedPlanet), planet)) {
			commandMapObject(planet);
		}

		// if we were selecting this planet, reselect it to trigger reactivity
		if (equal(get(selectedMapObject), planet)) {
			selectMapObject(planet);
		}

		// trigger reactivity
		universe.set(u);
	}

	// after a minefield is updated from the server, update the minefield in the universe, reset any commanded/selected
	// state and trigger reactivity
	function updateMinefield(minefield: Minefield, updatedMinefield: Minefield) {
		minefield = Object.assign(minefield, updatedMinefield);
		const u = get(universe);
		u.updateMinefield(minefield);

		// if we were selecting this minefield, reselect it to trigger reactivity
		if (equal(get(selectedMapObject), minefield)) {
			selectMapObject(minefield);
		}

		// trigger reactivity
		universe.set(u);
	}

	// load the status of a game, but not all the universe data
	async function loadStatus(): Promise<void> {
		const resp = await gameClient.getGame({ gameId });
		if (resp.game) {
			updateGame(getGameWithPlayersFlat(resp.game));
		}
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

	async function submitTurn(): Promise<void> {
		const { player, game } = await playerClient.submitTurn({ gameId });
		await updatePlayer(player);
		if (game) {
			updateGame(getGameWithPlayersFlat(game));
		}
	}

	async function forceGenerateTurn(): Promise<void> {
		await gameClient.forceGenerateTurn({ gameId });
		await loadStatus();
	}

	async function updatePlayerOrders(): Promise<void> {
		const p = get(player);
		const { player: updated, planets } = await playerClient.updatePlayerOrders({
			gameId,
			orders: p.playerOrders
		});
		await updatePlayer(updated);

		const u = get(universe);
		u.planets = planets;
		u.resetAllPlanets();
		planets.forEach((planet) => {
			if (equal(get(selectedMapObject), planet)) {
				selectMapObject(planet);
			}
		});

		u.resetMapObjectsByPosition();
		u.resetMyMapObjectsByPosition();

		universe.set(u);
	}

	async function updatePlayerRelations(relations: PlayerRelationship[]): Promise<void> {
		const p = get(player);
		const result = await playerClient.updatePlayerRelations({ gameId, relations: relations });
		p.relations = result.relations;
	}

	async function createBattlePlan(plan: BattlePlan): Promise<BattlePlan> {
		const { plan: created } = await battlePlanClient.createBattlePlan({ gameId, plan });
		const p = get(player);
		p.playerPlans.battlePlans = [...p.playerPlans.battlePlans, created!];
		await updatePlayer(p);

		return created!;
	}

	async function updateBattlePlan(plan: BattlePlan): Promise<BattlePlan> {
		const { plan: updated } = await battlePlanClient.updateBattlePlan({ gameId, plan });

		const p = get(player);
		for (let i = 0; i < p.playerPlans.battlePlans.length; i++) {
			if (p.playerPlans.battlePlans[i].num === updated?.num) {
				p.playerPlans.battlePlans[i] = updated;
			}
		}
		await updatePlayer(p);

		return updated!;
	}

	async function deleteBattlePlan(num: number): Promise<void> {
		const { player, fleets, starbases } = await battlePlanClient.deleteBattlePlan({ gameId, num });
		await updatePlayer(player);
		const u = get(universe);
		u.fleets = fleets;
		u.starbases = starbases;
		universe.set(u);
	}

	async function createProductionPlan(plan: ProductionPlan): Promise<ProductionPlan> {
		const { plan: created } = await productionPlanClient.createProductionPlan({ gameId, plan });
		const p = get(player);
		p.playerPlans.productionPlans = [...p.playerPlans.productionPlans, created!];
		await updatePlayer(p);

		return created!;
	}

	async function updateProductionPlan(plan: ProductionPlan): Promise<ProductionPlan> {
		const { plan: updated } = await productionPlanClient.updateProductionPlan({ gameId, plan });

		const p = get(player);
		for (let i = 0; i < p.playerPlans.productionPlans.length; i++) {
			if (p.playerPlans.productionPlans[i].num === updated?.num) {
				p.playerPlans.productionPlans[i] = updated;
			}
		}
		await updatePlayer(p);

		return updated!;
	}

	async function deleteProductionPlan(num: number): Promise<void> {
		await productionPlanClient.deleteProductionPlan({ gameId, num });
		const p = get(player);
		p.playerPlans.productionPlans = p.playerPlans.productionPlans.filter(
			(plan) => plan.num !== num
		);
		await updatePlayer(p);
	}

	async function createTransportPlan(plan: TransportPlan): Promise<TransportPlan> {
		const { plan: created } = await transportPlanClient.createTransportPlan({ gameId, plan });
		const p = get(player);
		p.playerPlans.transportPlans = [...p.playerPlans.transportPlans, created!];
		await updatePlayer(p);

		return created!;
	}

	async function updateTransportPlan(plan: TransportPlan): Promise<TransportPlan> {
		const { plan: updated } = await transportPlanClient.updateTransportPlan({ gameId, plan });
		const p = get(player);
		for (let i = 0; i < p.playerPlans.transportPlans.length; i++) {
			if (p.playerPlans.transportPlans[i].num === updated?.num) {
				p.playerPlans.transportPlans[i] = updated;
			}
		}
		await updatePlayer(p);

		return updated!;
	}

	async function deleteTransportPlan(num: number): Promise<void> {
		await transportPlanClient.deleteTransportPlan({ gameId, num });
		const p = get(player);
		p.playerPlans.transportPlans = p.playerPlans.transportPlans.filter((plan) => plan.num !== num);
		await updatePlayer(p);
	}

	async function createDesign(design: ShipDesign): Promise<ShipDesign> {
		// update this design
		const { design: created } = await shipDesignClient.createShipDesign({ gameId, design });
		const u = get(universe);
		u.addDesign(created!);
		await cs.wasmService.setDesigns({ designs: u.getMyDesigns() });
		universe.set(u);

		return created!;
	}

	async function updateDesign(design: ShipDesign): Promise<ShipDesign> {
		// update this design
		const { design: updated } = await shipDesignClient.updateShipDesign({ gameId, design });
		const u = get(universe);
		u.updateDesign(updated!);
		await cs.wasmService.setDesigns({ designs: u.getMyDesigns() });
		universe.set(u);

		return updated!;
	}

	async function deleteDesign(num: number): Promise<void> {
		const { planets, fleets, starbases } = await shipDesignClient.deleteShipDesign({ gameId, num });
		const u = get(universe);
		// replace our fleets and starbases (but keep intel issue #146)
		u.fleets = fleets;
		u.starbases = starbases;
		u.planets = planets;
		u.resetAllPlanets();
		u.resetMapObjectsByPosition();
		u.resetMyMapObjectsByPosition();
		universe.set(u);

		u.designs = u.designs.filter((d) => d.num != num);
		await cs.wasmService.setDesigns({ designs: u.getMyDesigns() });

		// reset our view to the homeworld, in case the commanded fleet had our deleted design
		commandHomeWorld();
	}

	async function addWaypoint(dest: WaypointDest, fastestWaypoint: boolean): Promise<boolean> {
		const fleet = get(commandedFleet);
		const sw = get(selectedWaypoint);
		const currentIndex = get(currentSelectedWaypointIndex);
		const u = get(universe);
		const fastest = get(settings).fastestWaypoint || fastestWaypoint;
		if (!fleet) {
			return false;
		}

		const result = await cs.wasmService.addWaypoint({
			fleet,
			dest,
			currentSelectedWaypointIndex: currentIndex,
			fastestWaypoint: fastest
		});

		if (!result.index || !result.fleet?.fleetOrders?.waypoints) {
			return false;
		}

		fleet.fleetOrders = result.fleet.fleetOrders;
		await updateFleetOrders(fleet);

		// select the new waypoint
		selectWaypoint(fleet.fleetOrders?.waypoints[result.index]);
		if (sw && sw.mapObjectTarget?.targetType && sw.mapObjectTarget?.targetNum) {
			const mo = u.getMapObject(sw.mapObjectTarget);

			if (mo) {
				selectMapObject(mo);
			}
		}

		return true;
	}

	async function updateWaypoint(dest: WaypointDest, fastestWaypoint: boolean, done: boolean) {
		const fleet = get(commandedFleet);
		const sw = get(selectedWaypoint);
		const currentIndex = get(currentSelectedWaypointIndex);
		const u = get(universe);
		const fastest = get(settings).fastestWaypoint || fastestWaypoint;

		if (!fleet || !sw) {
			return;
		}

		const result = await cs.wasmService.updateWaypoint({
			fleet,
			dest,
			currentSelectedWaypointIndex: currentIndex,
			fastestWaypoint: fastest
		});

		if (!result.fleet) {
			console.error('error updating fleet waypoint, no fleet returned from wasm');
		}

		if (result.updated) {
			// update the selectedWaypoint while dragging
			selectedWaypoint.update(() =>
				Object.assign(sw, result.fleet?.fleetOrders?.waypoints[currentIndex])
			);
			// check if we are done updating this waypoint and should save it to the server
			if (done) {
				// don't dragging, update the fleet
				(fleet as Fleet).fleetOrders =
					result.fleet?.fleetOrders ??
					// create a valid FleetOrders value with the correct shape
					create(FleetOrdersSchema);
				await updateFleetOrders(fleet);

				// select the new waypoint
				selectWaypoint(fleet.fleetOrders?.waypoints[currentIndex]);
				if (sw && sw.mapObjectTarget?.targetType && sw.mapObjectTarget?.targetNum) {
					const mo = u.getMapObject(sw.mapObjectTarget);

					if (mo) {
						selectMapObject(mo);
					}
				}
			} else {
				// trigger reaction
				selectedWaypoint.update(() => sw);
			}
		} else {
			// TODO: this logic is hard to follow with deletes and all that
			if (done) {
				// we dragged a waypoint to the previous position, delete it
				deleteWaypoint();
			}
		}
	}

	async function deleteWaypoint() {
		const fleet = get(commandedFleet);
		const sw = get(selectedWaypoint);
		const selectedWaypointIndex = get(currentSelectedWaypointIndex);
		const u = get(universe);

		if (!fleet || !selectedWaypoint || selectedWaypointIndex == 0) {
			return;
		}

		fleet.fleetOrders.waypoints = fleet.fleetOrders?.waypoints.filter((wp) => wp != sw);

		// select the previous waypoint
		const wp = fleet.fleetOrders?.waypoints[selectedWaypointIndex - 1];
		selectWaypoint(wp);

		const mo = u.getMapObject(wp.mapObjectTarget);
		if (mo) {
			selectMapObject(mo);
		}

		updateFleetOrders(fleet);
	}

	async function updateFleetOrders(fleet: CommandedFleet): Promise<void> {
		const { fleet: updatedFleet } = await fleetClient.updateFleetOrders({
			gameId,
			fleetNum: fleet.mapObject?.num ?? 0,
			fleetOrders: fleet.fleetOrders
		});
		if (updatedFleet) {
			updateFleet(fleet, updatedFleet);
		}
	}

	async function renameFleet(fleet: CommandedFleet, name: string): Promise<void> {
		const { fleet: updatedFleet } = await fleetClient.renameFleet({
			gameId,
			fleetNum: fleet.mapObject?.num ?? 0,
			name: name
		});
		if (updatedFleet) {
			updateFleet(fleet, updatedFleet);
		}
	}

	async function updatePlanetOrders(planet: CommandedPlanet): Promise<void> {
		const resp = await planetClient.updatePlanetOrders({
			gameId,
			planetNum: planet.mapObject.num,
			planetOrders: planet.planetOrders
		});

		if (resp.player && resp.planet) {
			// changing the planet orders changes the player's spec
			await updatePlayer(resp.player);
			updatePlanet(planet, resp.planet);
		}
	}

	async function updateMinefieldOrders(minefield: Minefield): Promise<void> {
		const { minefield: updatedMinefield } = await minefieldClient.updateMinefieldOrders({
			gameId,
			minefieldNum: minefield.mapObject?.num ?? 0,
			minefieldOrders: minefield.minefieldOrders
		});
		if (updatedMinefield) {
			updateMinefield(minefield, updatedMinefield);
		}
	}

	async function transferCargo(
		fleet: CommandedFleet,
		dest: CargoDest,
		transferAmount: CargoTransferRequest
	): Promise<void> {
		const result = await fleetClient.transferCargo({
			gameId,
			fleetNum: fleet.mapObject?.num ?? 0,
			mo: dest?.mapObject,
			transferAmount: create(CargoSchema, transferAmount),
			fuelTransferAmount: transferAmount.fuel
		});
		const u = get(universe);

		await updatePlayer(result.player);
		updateFleet(fleet, result.fleet);

		if (result.dest?.value?.mapObject?.type === MapObjectType.PLANET) {
			const planet = create(PlanetSchema, result.dest.value as Planet);
			updatePlanet(dest as Planet, planet);
		} else if (result.dest?.value?.mapObject?.type === MapObjectType.FLEET) {
			// update the destination fleet in the universe
			const destFleet = create(FleetSchema, result.dest.value as Fleet);
			updateFleet(dest as Fleet, destFleet);
		} else if (result.dest?.value?.mapObject?.type === MapObjectType.MINERAL_PACKET) {
			// TODO: do we need to create different mineral packets here?
			const destMineralPacket = create(MineralPacketSchema, result.dest.value as MineralPacket);
			u.updateMineralPacket(destMineralPacket);
		} else if (result.dest?.value?.mapObject?.type === MapObjectType.SALVAGE) {
			const destSalvage = create(SalvageIntelSchema, result.dest.value as SalvageIntel);
			u.updateSalvage(destSalvage);
		}

		const smo = get(selectedMapObject);
		if (smo && smo.mapObject?.type === MapObjectType.SALVAGE) {
			const salvage = u.getSalvage(smo.mapObject?.num);
			if (salvage) {
				selectMapObject(salvage);
			}
		}
		if (smo && smo.mapObject?.type === MapObjectType.MINERAL_PACKET) {
			const mineralPacket = u.getMineralPacket(smo.mapObject?.playerNum, smo.mapObject?.num);
			if (mineralPacket) {
				selectMapObject(mineralPacket);
			}
		}
	}

	async function split(
		src: CommandedFleet,
		dest: Fleet | undefined,
		sourceTokens: ShipToken[],
		destTokens: ShipToken[],
		transferAmount: CargoTransferRequest
	): Promise<void> {
		const {
			source: updatedSource,
			dest: updatedDest,
			cargoTransfers
		} = await fleetClient.splitFleet({
			gameId,
			sourceFleetNum: src.mapObject.num,
			destFleetNum: dest?.mapObject?.num,
			destBaseName: dest?.baseName,
			sourceTokens,
			destTokens,
			transferAmount,
			fuelTransferAmount: transferAmount.fuel
		});

		const u = get(universe);

		// if the original source is different than the returned source
		// and we have no response.dest, this means the original source was deleted
		// and has become the source, i.e. we moved all tokens from source to dest, creating
		// a new fleet. Weird edge case.
		if (updatedSource && src.mapObject.num !== updatedSource.mapObject?.num) {
			u.removeFleets([src.mapObject.num]);
			u.addFleets([updatedSource]);
		}

		// update the commanded fleet
		const commandedFleet = Object.assign(new CommandedFleet(), updatedSource);
		u.updateFleet(updatedSource);
		commandMapObject(commandedFleet);
		if (equal(get(selectedMapObject), commandedFleet)) {
			selectMapObject(commandedFleet);
		}

		// update or add the new fleets to the universe
		if (updatedDest) {
			if (!dest?.mapObject?.num) {
				u.addFleets([updatedDest]);
			} else {
				u.updateFleet(updatedDest);
			}
		} else {
			// if we had a dest and it was deleted, remove it
			if (dest?.mapObject?.num) {
				u.removeFleets([dest.mapObject.num]);
			}
		}

		const index = get(currentSelectedWaypointIndex);
		if (
			index > -1 &&
			commandedFleet.fleetOrders?.waypoints &&
			commandedFleet.fleetOrders?.waypoints.length > index
		) {
			selectWaypoint(commandedFleet.fleetOrders?.waypoints[index]);
		}
		await updateCargoTransfers(cargoTransfers);
	}

	async function splitAll(fleet: CommandedFleet): Promise<void> {
		const { fleets, cargoTransfers } = await fleetClient.splitAllFleets({
			gameId,
			fleetNum: fleet.mapObject?.num ?? 0
		});

		const sourceFleet = fleets.find((f) => f.mapObject?.num == fleet.mapObject?.num);
		if (sourceFleet) {
			fleet = Object.assign(fleet, sourceFleet);
			commandMapObject(fleet);
		}

		const u = get(universe);

		// update and add the new fleets to the universe
		u.updateFleet(fleet);
		u.addFleets(fleets.filter((f) => f.mapObject?.num != fleet.mapObject?.num));

		const index = get(currentSelectedWaypointIndex);
		if (index > -1 && fleet.fleetOrders?.waypoints && fleet.fleetOrders?.waypoints.length > index) {
			selectWaypoint(fleet.fleetOrders?.waypoints[index]);
		}

		if (cargoTransfers) {
			await updateCargoTransfers(cargoTransfers);
		}
	}

	async function merge(fleet: CommandedFleet, fleetNums: number[]): Promise<void> {
		const { fleet: updatedFleet, cargoTransfers } = await fleetClient.mergeFleets({
			gameId,
			fleetNum: fleet.mapObject?.num ?? 0,
			fleetNums
		});

		get(universe).removeFleets(fleetNums);
		updateFleet(fleet, updatedFleet);
		await updateCargoTransfers(cargoTransfers);
	}

	async function updateCargoTransfers(cargoTransfers: { [key: string]: CargoTransfers }) {
		const p = get(player);
		p.playerOrders.cargoTransfers = cargoTransfers;
		await updatePlayer(p);
	}

	return {
		cs,
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
		gotoTarget,
		gotoBattle,

		updateGame,
		loadStatus,
		startPollingStatus,
		stopPollingStatus,

		submitTurn,
		forceGenerateTurn,

		updatePlayerOrders,
		updatePlayerRelationships: updatePlayerRelations,
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

		addWaypoint,
		updateWaypoint,
		deleteWaypoint,
		updateFleetOrders,
		renameFleet,

		updatePlanetOrders,
		updateMinefieldOrders,
		transferCargo,
		split,
		splitAll,
		merge,
		resetContext,
		fullyLoaded,
		setFullyLoaded
	};
}

function loadSettingsOrDefault(gameId: bigint, playerNum: number): PlayerSettings {
	const key = PlayerSettings.key(gameId, playerNum);

	const json = localStorage.getItem(key);
	if (json) {
		const settingsJSON = JSON.parse(json) as PlayerSettings;
		if (settingsJSON) {
			// create a new object
			const settings = new PlayerSettings(`${gameId}`, playerNum);
			Object.assign(settings, settingsJSON);
			settings.afterLoad();
			return settings;
		}
	}

	return new PlayerSettings(`${gameId}`, playerNum);
}
