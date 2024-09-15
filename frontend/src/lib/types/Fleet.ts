import type { DesignFinder, Universe } from '$lib/services/Universe';
import { get as pluck } from 'lodash-es';
import { totalCargo, type Cargo } from './Cargo';
import type { Cost } from './Cost';
import {
	MapObjectType,
	None,
	owned,
	ownedBy,
	StargateWarpSpeed,
	type MapObject,
	type MovingMapObject
} from './MapObject';
import type { MessageTargetType } from './Message';
import type { MineFieldType } from './MineField';
import type { Planet } from './Planet';
import type { Player } from './Player';
import type { ShipDesign } from './ShipDesign';
import type { Engine } from './Tech';
import { distance, equal, type Vector } from './Vector';

export type Fleet = {
	playerNum: number; // override mapObject fleets always have a player.
	planetNum?: number;
	baseName: string;
	fuel?: number;
	cargo?: Cargo;
	damage?: number;
	tokens?: ShipToken[];
	mass?: number;
	scanRange?: number; // discoverable for allies when scanning
	scanRangePen?: number;
	freighter?: boolean;
	orbitingPlanetNum?: number;
	starbase?: boolean;
	spec?: Spec;
} & MovingMapObject &
	FleetOrders;

export type FleetOrders = {
	waypoints?: Waypoint[];
	repeatOrders?: boolean;
	battlePlanNum?: number;
};

export type ShipToken = {
	id?: number;
	createdAt?: string;
	updatedAt?: string;

	gameId?: number;
	designNum: number;
	quantity: number;
	damage?: number;
	quantityDamaged?: number;
};

export type Target = {
	targetType?: MapObjectType | MessageTargetType;
	targetPosition?: Vector;
	targetPlayerNum?: number;
	targetNum?: number;
	targetName?: string;
};

export type Waypoint = {
	position: Vector;
	warpSpeed: number;
	estFuelUsage?: number;
	task: WaypointTask;
	waitAtWaypoint?: boolean;
	layMineFieldDuration?: number;
	patrolRange?: number;
	patrolWarpSpeed?: number;
	transferToPlayer?: number;
	partiallyComplete?: boolean;
	transportTasks: WaypointTransportTasks;
} & Target;

export enum WaypointTask {
	None = '',
	Transport = 'Transport',
	Colonize = 'Colonize',
	RemoteMining = 'RemoteMining',
	MergeWithFleet = 'MergeWithFleet',
	ScrapFleet = 'ScrapFleet',
	LayMineField = 'LayMineField',
	Patrol = 'Patrol',
	Route = 'Route',
	TransferFleet = 'TransferFleet'
}

export type WaypointTransportTasks = {
	fuel: WaypointTransportTask;
	ironium: WaypointTransportTask;
	boranium: WaypointTransportTask;
	germanium: WaypointTransportTask;
	colonists: WaypointTransportTask;
};

export const emptyWaypointTransportTasks = (): WaypointTransportTasks => ({
	fuel: {},
	ironium: {},
	boranium: {},
	germanium: {},
	colonists: {}
});

export type WaypointTransportTask = {
	action?: WaypointTaskTransportAction;
	amount?: number;
};

export enum WaypointTaskTransportAction {
	None = '',
	LoadOptimal = 'LoadOptimal',
	LoadAll = 'LoadAll',
	UnloadAll = 'UnloadAll',
	LoadAmount = 'LoadAmount',
	UnloadAmount = 'UnloadAmount',
	FillPercent = 'FillPercent',
	WaitForPercent = 'WaitForPercent',
	LoadDunnage = 'LoadDunnage',
	SetAmountTo = 'SetAmountTo',
	SetWaypointTo = 'SetWaypointTo'
}

export type Spec = {
	engine: Engine;
	cost: Cost;
	mass: number;
	armor: number;
	fuelCapacity: number;
	immuneToOwnDetonation?: boolean;
	mineLayingRateByMineType?: Record<MineFieldType, number>;
	weaponSlots?: null;
	purposes?: any;
	totalShips: number;
	massEmpty: number;
	basePacketSpeed: number;
	safePacketSpeed?: number;
	baseCloakedCargo: number;
	stargate?: string;
	massDriver?: string;

	numEngines?: number;
	estimatedRange?: number;
	cargoCapacity?: number;
	cloakUnits?: number;
	scanRange?: number;
	scanRangePen?: number;
	repairBonus?: number;
	torpedoBonus?: number;
	torpedoJamming?: number;
	initiative?: number;
	movement?: number;
	powerRating?: number;
	bomber?: number;
	bombs?: number;
	smartBombs?: number;
	retroBombs?: number;
	scanner?: boolean;
	shields?: number;
	colonizer?: boolean;
	canLayMines?: boolean;
	spaceDock?: number;
	miningRate?: number;
	terraformRate?: number;
	mineSweep?: number;
	cloakPercent?: number;
	reduceCloaking?: number;
	canStealFleetCargo?: number;
	canStealPlanetCargo?: number;
	orbitalConstructionModule?: number;
	hasWeapons?: boolean;
	hasStargate?: boolean;
	hasMassDriver?: boolean;
	canJump?: boolean; // TODO: actually implement this
	maxPopulation?: number;
};

export function emptyTransportTasks(): WaypointTransportTasks {
	return {
		fuel: {
			action: WaypointTaskTransportAction.None
		},
		ironium: {
			action: WaypointTaskTransportAction.None
		},
		boranium: {
			action: WaypointTaskTransportAction.None
		},
		germanium: {
			action: WaypointTaskTransportAction.None
		},
		colonists: {
			action: WaypointTaskTransportAction.None
		}
	};
}

export class CommandedFleet implements Fleet {
	id = 0;
	gameId = 0;
	createdAt?: string | undefined;
	updatedAt?: string | undefined;
	readonly type = MapObjectType.Fleet;

	name = '';
	playerNum = 0;
	num = 0;

	planetNum = undefined;
	baseName = '';
	fuel = 0;
	cargo: Cargo = {};
	damage = 0;
	battlePlanNum = 0;
	tokens: ShipToken[] = [];
	waypoints: Waypoint[] = [];
	repeatOrders = false;
	heading = { x: 0, y: 0 };
	warpSpeed = 0;
	mass = 0;
	orbitingPlanetNum = None;
	starbase = false;
	position = { x: 0, y: 0 };
	spec = {} as Spec;

	constructor(data?: Fleet) {
		Object.assign(this, data);
	}

	getFuelCost(
		designFinder: DesignFinder,
		fuelEfficiencyOffset: number,
		warpSpeed: number,
		distance: number,
		cargoCapacity: number
	): number {
		const efficiencyFactor: number = 1 + fuelEfficiencyOffset;
		let fuelCost = 0;

		for (const token of this.tokens) {
			const design = designFinder.getDesign(this.playerNum, token.designNum);
			if (design?.spec) {
				let mass: number = (design.spec.mass ?? 0) * token.quantity;
				const fleetCargo: number = totalCargo(this.cargo);
				const stackCapacity: number = (design.spec.cargoCapacity ?? 0) * token.quantity;

				if (cargoCapacity > 0) {
					mass += Math.floor((fleetCargo * stackCapacity) / cargoCapacity);
				}

				const engine: Engine = design.spec.engine;
				fuelCost += getFuelCostForEngine(engine, warpSpeed, mass, distance, efficiencyFactor);
			}
		}

		return fuelCost;
	}

	getSelectedWaypointInfo(currentSelectedWaypointIndex: number = -1): {
		selectedWaypoint: Waypoint;
		nextWaypoint: Waypoint | undefined;
		previousWaypoint: Waypoint | undefined;
		waypointIndex: number;
	} {
		// default to the first waypoint if the current selected waypoint index is invalid
		const index: number =
			currentSelectedWaypointIndex == -1 || currentSelectedWaypointIndex >= this.waypoints.length
				? 0
				: currentSelectedWaypointIndex;

		return {
			selectedWaypoint: this.waypoints[index],
			nextWaypoint: index < this.waypoints.length - 1 ? this.waypoints[index + 1] : undefined,
			previousWaypoint: index > 0 ? this.waypoints[index - 1] : undefined,
			waypointIndex: index
		};
	}

	/**
	 * Add a waypoint to a destination, returning the index of the newly added waypoint
	 * @param dest
	 * @param orbiting
	 */
	addWaypoint(
		player: Player,
		designFinder: DesignFinder,
		dest: { mo: MapObject; position?: never } | { mo?: never; position: Vector },
		currentSelectedWaypointIndex: number,
		orbiting: Planet | undefined,
		highestShipMass: number,
		fastestWaypoint: boolean
	): number | undefined {
		const { selectedWaypoint, nextWaypoint, waypointIndex } = this.getSelectedWaypointInfo(
			currentSelectedWaypointIndex
		);
		const mo = dest.mo;
		const position = dest.position ?? dest.mo.position;

		if (
			equal(position, selectedWaypoint.position) ||
			(nextWaypoint && equal(position, nextWaypoint.position))
		) {
			// don't add a duplicate waypoint waypoints
			return;
		}

		let fuelAlreadyAllocated = 0;
		for (let i = 0; i < waypointIndex; i++) {
			fuelAlreadyAllocated += this.waypoints[i].estFuelUsage ?? 0;
		}

		// if our destination is a planet, determine some stuff about it
		const { warpSpeed, canColonize, canRemoteMine } = this.getWarpSpeed(
			player,
			designFinder,
			dest,
			orbiting,
			fuelAlreadyAllocated,
			highestShipMass,
			fastestWaypoint
		);

		const task = selectedWaypoint.task ?? WaypointTask.None;
		const emptyTransportTasks = emptyWaypointTransportTasks();
		const transportTasks = selectedWaypoint.transportTasks ?? emptyTransportTasks;

		if (mo) {
			// create a waypoint with a MapObject as a target
			const wp: Waypoint = {
				position: mo.position,
				targetName: mo.name,
				targetPlayerNum: mo.playerNum,
				targetNum: mo.num,
				targetType: mo.type,
				warpSpeed: warpSpeed,
				task: task,
				transportTasks: transportTasks
			};
			this.waypoints.splice(waypointIndex + 1, 0, wp);

			// if this is a colonizer and the target is a habitable planet
			if (canColonize) {
				wp.task = WaypointTask.Colonize;
				wp.transportTasks = emptyTransportTasks;
			} else if (canRemoteMine) {
				wp.task = WaypointTask.RemoteMining;
				wp.transportTasks = emptyTransportTasks;
			}
		} else {
			this.waypoints.splice(waypointIndex + 1, 0, {
				position: dest.position,
				warpSpeed: warpSpeed,
				task: task,
				transportTasks: transportTasks
			});
		}

		return waypointIndex + 1;
	}

	/**
	 * Add a waypoint to a destination, returning the index of the newly added waypoint
	 * @param dest
	 * @param orbiting
	 */
	updateWaypoint(
		player: Player,
		designFinder: DesignFinder,
		dest: { mo: MapObject; position?: never } | { mo?: never; position: Vector },
		currentSelectedWaypointIndex: number,
		orbiting: Planet | undefined,
		highestShipMass: number,
		fastestWaypoint: boolean
	): boolean {
		const { selectedWaypoint, previousWaypoint, waypointIndex } = this.getSelectedWaypointInfo(
			currentSelectedWaypointIndex
		);
		const mo = dest.mo;
		const position = dest.position ?? dest.mo.position;

		if (previousWaypoint && equal(position, previousWaypoint.position)) {
			// don't update a waypoint to be the same as a previous waypoint, this should just delete it
			return false;
		}

		let fuelAlreadyAllocated = 0;
		for (let i = 0; i < waypointIndex; i++) {
			fuelAlreadyAllocated += this.waypoints[i].estFuelUsage ?? 0;
		}

		// if our destination is a planet, determine some stuff about it
		const { warpSpeed, canColonize, canRemoteMine } = this.getWarpSpeed(
			player,
			designFinder,
			dest,
			orbiting,
			fuelAlreadyAllocated,
			highestShipMass,
			fastestWaypoint
		);

		let task = selectedWaypoint.task ?? WaypointTask.None;
		const emptyTransportTasks = emptyWaypointTransportTasks();
		const transportTasks = selectedWaypoint.transportTasks ?? emptyTransportTasks;

		// don't update the waypoint to a colonize/remote mine task if we can't do it on this new target
		if (
			(task == WaypointTask.Colonize && !canColonize) ||
			(task == WaypointTask.RemoteMining && !canRemoteMine)
		) {
			task = WaypointTask.None;
		}

		if (mo) {
			selectedWaypoint.position = mo.position;
			selectedWaypoint.targetName = mo.name;
			selectedWaypoint.targetPlayerNum = mo.playerNum;
			selectedWaypoint.targetNum = mo.num;
			selectedWaypoint.targetType = mo.type;
			selectedWaypoint.warpSpeed = warpSpeed;
			selectedWaypoint.task = task;
			selectedWaypoint.transportTasks = transportTasks;

			// if this is a colonizer and the target is a habitable planet
			if (canColonize) {
				selectedWaypoint.task = WaypointTask.Colonize;
				selectedWaypoint.transportTasks = emptyTransportTasks;
			} else if (canRemoteMine) {
				selectedWaypoint.task = WaypointTask.RemoteMining;
				selectedWaypoint.transportTasks = emptyTransportTasks;
			}
		} else {
			selectedWaypoint.position = dest.position;
			selectedWaypoint.targetName = '';
			selectedWaypoint.targetPlayerNum = None;
			selectedWaypoint.targetNum = None;
			selectedWaypoint.targetType = MapObjectType.None;
			selectedWaypoint.warpSpeed = warpSpeed;
			selectedWaypoint.task = task;
			selectedWaypoint.transportTasks = transportTasks;
		}

		return true;
	}

	/**
	 * Get the warpSpeed of a waypoint to a destination. Also return whether we can colonize or remote mine the dest
	 * @param player
	 * @param designFinder
	 * @param dest
	 * @param orbiting
	 * @param highestShipMass
	 * @param fastestWaypoint
	 * @returns
	 */
	getWarpSpeed(
		player: Player,
		designFinder: DesignFinder,
		dest: { mo: MapObject; position?: never } | { mo?: never; position: Vector },
		orbiting: Planet | undefined,
		fuelAlreadyAllocated: number,
		highestShipMass: number,
		fastestWaypoint: boolean
	): { warpSpeed: number; canColonize: boolean; canRemoteMine: boolean } {
		const mo = dest.mo;
		const position = dest.position ?? dest.mo.position;
		const dist = Math.floor(distance(this.position, position));

		let canColonize = false;
		let canRemoteMine = false;
		let canGate = false;
		if (mo && mo.type == MapObjectType.Planet) {
			const target = mo as Planet;
			canColonize = this.canColonize(target);
			canRemoteMine = this.canRemoteMine(player, target);
			canGate = this.canGate(player, orbiting, target, dist, highestShipMass);
		}

		// set our warp speed to the most fuel efficient based on our engine idealSpeed (7 for Long Hump 7, 8 for Alpha Drive 8, etc)
		// or the fastest warp we can get there without running out of fuel
		let warpSpeed = canGate
			? StargateWarpSpeed // stargate speed if we can gate
			: canColonize || fastestWaypoint // max speed if configured for that, or colonizing
				? this.getMaxWarp(
						designFinder,
						player.race.spec?.fuelEfficiencyOffset ?? 0,
						fuelAlreadyAllocated,
						dist,
						this.spec.engine.freeSpeed ?? 1,
						this.spec?.engine?.maxSafeSpeed ?? 9
					)
				: this.getMinimalWarp(
						designFinder,
						player.race.spec?.fuelEfficiencyOffset ?? 0,
						fuelAlreadyAllocated,
						dist,
						this.spec.engine.idealSpeed ?? 0,
						this.spec.engine.freeSpeed ?? 1,
						this.spec?.engine?.maxSafeSpeed ?? 9
					);
		return { warpSpeed, canColonize, canRemoteMine };
	}

	// get the highest useful speed less than or equal to a given warp speed
	// needed to reach the destination
	getMinimalWarp(
		designFinder: DesignFinder,
		fuelEfficiencyOffset: number,
		fuelAlreadyAllocated: number,
		dist: number,
		startSpeed: number,
		freeSpeed: number,
		maxSafeSpeed: number
	): number {
		const yearsAtIdealSpeed = Math.ceil(dist / (startSpeed * startSpeed));

		// start checking 1 warp speed below our maximum assigned speed
		let speed = startSpeed;
		for (let i = startSpeed; i > freeSpeed; i--) {
			const yearsAtSpeed = Math.ceil(dist / (i * i));
			// if it takes the same time to go slower, go slower
			if (Math.ceil(yearsAtIdealSpeed) == Math.ceil(yearsAtSpeed)) {
				speed = i;
			}
		}

		// start at speed and go backwards if we would run out of fuel at this speed
		while (speed >= freeSpeed) {
			const fuelUsed = this.getFuelCost(
				designFinder,
				fuelEfficiencyOffset,
				speed,
				dist,
				this.spec.cargoCapacity ?? 0
			);
			if (fuelUsed + fuelAlreadyAllocated > this.fuel) {
				// ran out of fuel, go slower and try again
				speed--;
				continue;
			}
			break;
		}

		return Math.min(maxSafeSpeed, speed);
	}

	// get the max warp we have fuel for to make it to the destination
	getMaxWarp(
		designFinder: DesignFinder,
		fuelEfficiencyOffset: number,
		fuelAlreadyAllocated: number,
		dist: number,
		freeSpeed: number,
		maxSafeSpeed: number
	): number {
		// start at one above free speed and add to it until we run out of fuel
		let speed: number;
		for (speed = freeSpeed + 1; speed <= maxSafeSpeed; speed++) {
			const fuelUsed = this.getFuelCost(
				designFinder,
				fuelEfficiencyOffset,
				speed,
				dist,
				this.spec.cargoCapacity ?? 0
			);
			if (fuelUsed + fuelAlreadyAllocated > this.fuel || speed > maxSafeSpeed) {
				// ran out of fuel, go back one speed and we're done
				speed--;
				break;
			}
		}

		const idealSpeed = this.spec?.engine?.idealSpeed ?? 5;
		const idealFuelUsed = this.getFuelCost(
			designFinder,
			fuelEfficiencyOffset,
			idealSpeed,
			dist,
			this.spec.cargoCapacity ?? 0
		);

		// if we are using a ramscoop, make sure we at least go the ideal
		// speed of the engine if we can. If we run out,
		// it'll drop to the free speed
		if (freeSpeed > 1 && speed < idealSpeed && idealFuelUsed > this.fuel) {
			speed = idealSpeed;
		}

		// don't go faster than we need
		return this.getMinimalWarp(
			designFinder,
			fuelEfficiencyOffset,
			fuelAlreadyAllocated,
			dist,
			speed,
			freeSpeed,
			maxSafeSpeed
		);
	}

	/**
	 *
	 * @param target the target planet to check
	 * @returns true if this fleet can colonize this planet
	 */
	canColonize(target: Planet): boolean {
		return !!(
			this.spec.colonizer &&
			this.cargo.colonists &&
			!owned(target) &&
			(target.spec.terraformedHabitability ?? 0) > 0
		);
	}

	/**
	 *
	 * @param target the target planet to check
	 * @returns true if this fleet can colonize this planet
	 */
	canRemoteMine(player: Player, target: Planet): boolean {
		return !!(
			this.spec.miningRate &&
			this.spec.miningRate > 0 &&
			(!owned(target) || (ownedBy(target, player.num) && player.race.spec?.canRemoteMineOwnPlanets))
		);
	}

	/**
	 *
	 * @param player The fleet player
	 * @param orbiting the planet the fleet is orbiting (or undefined if not orbiting a planet)
	 * @param targetPlanet the planet the fleet is targeting
	 * @param dist the distance away of the target planet
	 * @param highestShipMass the highest mass of any ship in the fleet
	 * @returns true if the fleet can gate to this planet
	 */
	canGate(
		player: Player,
		orbiting: Planet | undefined,
		targetPlanet: Planet,
		dist: number,
		highestShipMass: number
	): boolean {
		const destSafeHullMass = targetPlanet.spec.safeHullMass ?? 0;
		const destSafeRange = targetPlanet.spec.safeRange ?? 0;
		const destStargateSafe =
			targetPlanet.spec.hasStargate &&
			owned(targetPlanet) &&
			player.isFriend(targetPlanet.playerNum ?? 0) &&
			destSafeRange >= dist &&
			highestShipMass <= destSafeHullMass;

		if (this.spec?.canJump) {
			// we have a jump gate installed in our ship, we only care about the destination gate
			return destStargateSafe;
		} else {
			if (!orbiting || !orbiting.spec.hasStargate) {
				return false;
			}
			const canGateCargo = player.race.spec?.canGateCargo;
			const sourceSafeHullMass = orbiting.spec.safeHullMass ?? 0;
			const sourceSafeRange = orbiting.spec.safeRange ?? 0;
			const sourceStargateSafe =
				(canGateCargo || totalCargo(this.cargo) == 0) &&
				owned(orbiting) &&
				player.isFriend(targetPlanet.playerNum ?? 0) &&
				sourceSafeRange >= dist &&
				highestShipMass <= sourceSafeHullMass;
			return !!(destStargateSafe && sourceStargateSafe);
		}
	}

	/**
	 *
	 * @returns The total number of mines laid per year for all types of minefields this fleet can lay
	 */
	getTotalMinesLaidPerYear() {
		if (this.spec.mineLayingRateByMineType) {
			return Object.values(this.spec.mineLayingRateByMineType).reduce((count, n) => count + n, 0);
		}
		return 0;
	}
}

export function getDamagePercentForToken(token: ShipToken, design: ShipDesign | undefined): number {
	const armor = design?.spec.armor ?? 0;
	const totalArmor = armor * token.quantity;
	const quantityDamaged =
		(token.quantityDamaged ?? 0) > (token.quantity ?? 0)
			? token.quantity ?? 0
			: token.quantityDamaged ?? 0;
	const totalDamage = quantityDamaged * (token.damage ?? 0);
	if (totalArmor > 0 && totalDamage > 0) {
		return (totalDamage / totalArmor) * 100;
	}
	return 0;
}

// true if this fleet can transfer cargo
export function canTransferCargo(fleet: Fleet, universe: Universe): boolean {
	if (!fleet.spec?.cargoCapacity) {
		return false;
	}
	if (fleet.orbitingPlanetNum) {
		const planet = universe.getPlanet(fleet.orbitingPlanetNum);
		if (planet && !ownedBy(planet, fleet.playerNum)) {
			// if any of these fleets can transport, it's a contested planet
			const orbitingForeignFreighters = universe
				.getMapObjectsByPosition(planet)
				.filter((mo) => mo.type === MapObjectType.Fleet)
				.map((mo) => mo as Fleet)
				.filter((f: Fleet) => f.freighter)
				.filter((f) => f.playerNum !== fleet.playerNum);

			// don't allow manual transfers over contested planets
			if (orbitingForeignFreighters.length > 0) {
				return false;
			}
		}
	}
	return true;
}

// This shows only your fleets that have no movement orders, and any active enemy ships (so you can match one with the other, if you wish).
export function idleFleetsFilter(fleet: Fleet, showIdleFleetsOnly: boolean): boolean {
	if (!showIdleFleetsOnly) {
		// no filter, show all fleets
		return true;
	}

	// show our fleets that are idle
	if (
		fleet.waypoints &&
		fleet.waypoints.length == 1 &&
		fleet.waypoints[0].task == WaypointTask.None
	) {
		return true;
	}

	// enemy fleet that is moving, show it so players can match idle fleets to moving fleets
	if (!fleet.waypoints && fleet.warpSpeed) {
		return true;
	}

	// don't show this fleet if we got here, it's our fleet and moving, or an enemy fleet and idle
	return false;
}

function getFuelCostForEngine(
	engine: Engine,
	warpSpeed: number,
	mass: number,
	dist: number,
	ifeFactor: number
): number {
	if (warpSpeed === 0 || engine.fuelUsage == undefined || warpSpeed >= engine.fuelUsage.length) {
		return 0;
	}

	const distanceCeiling: number = Math.ceil(dist);
	const engineEfficiency: number = Math.ceil(ifeFactor * engine.fuelUsage[warpSpeed]);
	const teorFuel: number = Math.floor((mass * engineEfficiency * distanceCeiling) / 2000) / 10;
	const intFuel: number = Math.ceil(teorFuel);

	return intFuel;
}

export const isLoadAction = (action: WaypointTaskTransportAction) =>
	[
		WaypointTaskTransportAction.LoadOptimal,
		WaypointTaskTransportAction.LoadAll,
		WaypointTaskTransportAction.LoadAmount,
		WaypointTaskTransportAction.LoadDunnage,
		WaypointTaskTransportAction.FillPercent,
		WaypointTaskTransportAction.WaitForPercent
	].indexOf(action) != -1;

export const isUnloadAction = (action: WaypointTaskTransportAction) =>
	[WaypointTaskTransportAction.UnloadAll, WaypointTaskTransportAction.UnloadAmount].indexOf(
		action
	) != -1;

export const getLocation = (fleet: Fleet, universe: Universe) =>
	fleet.orbitingPlanetNum
		? universe.getPlanet(fleet.orbitingPlanetNum)?.name ?? 'unknown'
		: `Space: (${fleet.position.x}, ${fleet.position.y})`;

export const getDestination = (fleet: Fleet, universe: Universe) => {
	if (fleet.waypoints?.length && fleet.waypoints?.length > 1) {
		return universe.getTargetName(fleet.waypoints[1]);
	}
	return '--';
};

export function getTokenCount(mo: MapObject) {
	if (mo.type == MapObjectType.Fleet) {
		const fleet = mo as Fleet;
		return fleet.tokens ? fleet.tokens.reduce((count, t) => count + t.quantity, 0) : 0;
	}
	return 0;
}

export function hasDestination(mo: MapObject): boolean {
	const fleet = mo.type == MapObjectType.Fleet ? (mo as Fleet) : undefined;
	return (fleet?.waypoints?.length ?? 0) > 1;
}

// fleetsSortBy returns a sortBy function for fleets by key. This is used by the fleets report page
// and sorting when cycling through Fleets
export function fleetsSortBy(
	key: string,
	universe: Universe
): ((a: Fleet, b: Fleet) => number) | undefined {
	switch (key) {
		case 'name':
			return (a, b) => a.name.localeCompare(b.name);
		case 'location':
			return (a, b) => getLocation(a, universe).localeCompare(getLocation(b, universe));
		case 'destination':
			return (a, b) => getDestination(a, universe).localeCompare(getDestination(b, universe));
		case 'cargo':
			return (a, b) => totalCargo(a.cargo) - totalCargo(b.cargo);
		case 'mass':
			return (a, b) => (a.spec?.mass ?? 0) - (b.spec?.mass ?? 0);
		case 'fuel':
			return (a, b) => (a.fuel ?? 0) - (b.fuel ?? 0);
		default:
			return (a, b) => {
				const aVal = pluck(a, key);
				const bVal = pluck(b, key);
				if (typeof aVal == 'number' && typeof bVal == 'number') {
					return aVal - bVal;
				}
				return `${aVal}`.localeCompare(`${bVal}`);
			};
	}
}

/**
 * Move a postitive quantity of damaged tokens from a source to a dest token
 * @param srcToken the source to move damaged tokens from
 * @param destToken the dest to move damaged tokens to
 * @param quantity a positive quanityt to move
 */
export function moveDamagedTokens(srcToken: ShipToken, destToken: ShipToken, quantity: number) {
	const quantityDamagedToMove = Math.min(quantity, srcToken.quantityDamaged ?? 0);

	// figure out how much total damage we are moving over and how much current damage there is
	// the idea is if we have a stack on each side like this:
	//
	// src = 1 damaged token @10 damage
	// dest = 1 damaged token @5 damage
	//
	// after moving 1 damaged token from src to dest, we have 2 damaged tokens with 15 total damage between (i.e 7.5 damage / token)
	const damageToMove = quantityDamagedToMove * (srcToken.damage ?? 0);
	const currentDestDamage = (destToken.quantityDamaged ?? 0) * (destToken.damage ?? 0);

	// Move up to quantity damaged tokens
	destToken.quantityDamaged = (destToken.quantityDamaged ?? 0) + quantityDamagedToMove;
	if (destToken.quantityDamaged) {
		destToken.damage = (currentDestDamage + damageToMove) / destToken.quantityDamaged;
	} else {
		destToken.damage = 0;
	}

	// move the damaged tokens away from the source and zero out the damage if we have no damaged tokens remaining
	srcToken.quantityDamaged = (srcToken.quantityDamaged ?? 0) - quantityDamagedToMove;
	if (srcToken.quantityDamaged == 0) {
		srcToken.damage = 0;
	}
}
