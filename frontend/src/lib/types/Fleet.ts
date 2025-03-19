import type {
	AnyFleet,
	AnyPlanet,
	AnyShipDesign,
	DesignFinder,
	Universe
} from '$lib/services/Universe';
import { get as pluck } from 'lodash-es';
import { totalCargo } from './Cargo';
import type { CargoDest } from './CargoTransferRequest.svelte';
import { owned, ownedBy } from './MapObject';
import type { CommandedPlayer } from './Player';
import { distance, equal } from './Vector';
import {
	type Cargo,
	type Engine,
	type Fleet,
	type FleetIntel,
	type FleetSpec,
	type MapObject,
	MapObjectTypeFleet,
	MapObjectTypeMineralPacket,
	MapObjectTypeNone,
	MapObjectTypePlanet,
	MapObjectTypeSalvage,
	None,
	type PlanetIntel,
	type ShipToken,
	StargateWarpSpeed,
	TransportActionFillPercent,
	TransportActionLoadAll,
	TransportActionLoadAmount,
	TransportActionLoadDunnage,
	TransportActionLoadOptimal,
	TransportActionNone,
	TransportActionSetAmountTo,
	TransportActionSetWaypointTo,
	TransportActionUnloadAll,
	TransportActionUnloadAmount,
	TransportActionWaitForPercent,
	type Vector,
	type Waypoint,
	type WaypointTask,
	WaypointTaskColonize,
	WaypointTaskLayMineField,
	WaypointTaskMergeWithFleet,
	WaypointTaskNone,
	WaypointTaskPatrol,
	WaypointTaskRemoteMining,
	WaypointTaskRoute,
	WaypointTaskScrapFleet,
	WaypointTaskTransferFleet,
	WaypointTaskTransport,
	type WaypointTaskTransportAction,
	type WaypointTransportTasks
} from './cs';

export const WaypointTasks: WaypointTask[] = [
	WaypointTaskNone,
	WaypointTaskTransport,
	WaypointTaskColonize,
	WaypointTaskRemoteMining,
	WaypointTaskMergeWithFleet,
	WaypointTaskScrapFleet,
	WaypointTaskLayMineField,
	WaypointTaskPatrol,
	WaypointTaskRoute,
	WaypointTaskTransferFleet
] as const;

export const WaypointTransportTaskActions: WaypointTaskTransportAction[] = [
	TransportActionNone,
	TransportActionLoadOptimal,
	TransportActionLoadAll,
	TransportActionUnloadAll,
	TransportActionLoadAmount,
	TransportActionUnloadAmount,
	TransportActionFillPercent,
	TransportActionWaitForPercent,
	TransportActionLoadDunnage,
	TransportActionSetAmountTo,
	TransportActionSetWaypointTo
] as const;

export const emptyWaypointTransportTasks = (): WaypointTransportTasks => ({
	fuel: {},
	ironium: {},
	boranium: {},
	germanium: {},
	colonists: {}
});

/** A destination for a waypoint - either a MapObject or a position in space (but not both) */
export type WaypointDest = { mo: MapObject; position?: never } | { mo?: never; position: Vector };

export function emptyTransportTasks(): WaypointTransportTasks {
	return {
		fuel: {
			action: TransportActionNone
		},
		ironium: {
			action: TransportActionNone
		},
		boranium: {
			action: TransportActionNone
		},
		germanium: {
			action: TransportActionNone
		},
		colonists: {
			action: TransportActionNone
		}
	};
}

export class CommandedFleet implements Fleet {
	id = 0;
	gameId = 0;
	createdAt = '';
	updatedAt = '';
	age = 0;
	readonly type = MapObjectTypeFleet;

	name = '';
	playerNum = 0;
	num = 0;

	planetNum = 0;
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
	spec = {} as FleetSpec;
	tags = {};

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

	getWaypointMapObjects(universe: Universe): MapObject[] {
		return this.waypoints.map((wp) => {
			const mo = universe.getMapObject(wp);
			if (mo) {
				return mo;
			} else {
				return {
					position: wp.position,
					type: wp.targetType ?? MapObjectTypeNone,
					name: wp.targetName ?? '',
					num: wp.targetNum ?? 0,
					playerNum: wp.targetPlayerNum ?? 0
				} as MapObject;
			}
		});
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
	 * Add a {@linkcode Waypoint} to this {@linkcode CommandedFleet}.
	 * @param player The {@linkcode Player} controlling the fleet.
	 * @param universe Universe object
	 * @param dest The {@linkcode WaypointDest|destination} of the newly placed waypoint.
	 * @param currentSelectedWaypointIndex
	 * @param highestShipMass the mass of the highest ship design in the fleet
	 * @param fastestWaypoint Whether to use the fastest warp speed or
	 * @returns The waypoint index of the newly placed waypoint
	 */
	addWaypoint(
		player: CommandedPlayer,
		universe: Universe,
		dest: WaypointDest,
		currentSelectedWaypointIndex: number,
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

		// get the fuel allocated up to the last waypoint
		const fuelAlreadyAllocated = this.getFuelAllocated(player, universe, waypointIndex);
		const orbiting =
			selectedWaypoint.targetType === MapObjectTypePlanet
				? universe.getPlanet(selectedWaypoint.targetNum ?? 0)
				: undefined;

		const dist = Math.floor(distance(selectedWaypoint.position, position));

		// if our destination is a planet, determine some stuff about it
		const { warpSpeed, canColonize, canRemoteMine } = this.getWarpSpeed(
			player,
			universe,
			dist,
			orbiting,
			dest,
			fuelAlreadyAllocated,
			highestShipMass,
			fastestWaypoint
		);

		const task = selectedWaypoint.task ?? WaypointTaskNone;
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
				wp.task = WaypointTaskColonize;
				wp.transportTasks = emptyTransportTasks;
			} else if (canRemoteMine) {
				wp.task = WaypointTaskRemoteMining;
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
		player: CommandedPlayer,
		universe: Universe,
		dest: WaypointDest,
		currentSelectedWaypointIndex: number,
		highestShipMass: number,
		fastestWaypoint: boolean
	): boolean {
		const { selectedWaypoint, previousWaypoint, waypointIndex } = this.getSelectedWaypointInfo(
			currentSelectedWaypointIndex
		);
		if (!previousWaypoint) {
			return false;
		}

		const mo = dest.mo;
		const position = dest.position ?? dest.mo.position;

		if (equal(position, previousWaypoint.position)) {
			// don't update a waypoint to be the same as a previous waypoint, this should just delete it
			return false;
		}

		// get the fuel allocated up to but not including this waypoint since we're moving it around
		const fuelAlreadyAllocated = this.getFuelAllocated(player, universe, waypointIndex - 1);

		const orbiting =
			previousWaypoint.targetType === MapObjectTypePlanet
				? universe.getPlanet(previousWaypoint.targetNum ?? 0)
				: undefined;

		const dist = Math.floor(distance(previousWaypoint?.position, position));

		// if our destination is a planet, determine some stuff about it
		const { warpSpeed, canColonize, canRemoteMine } = this.getWarpSpeed(
			player,
			universe,
			dist,
			orbiting,
			dest,
			fuelAlreadyAllocated,
			highestShipMass,
			fastestWaypoint
		);

		let task = selectedWaypoint.task ?? WaypointTaskNone;
		const emptyTransportTasks = emptyWaypointTransportTasks();
		const transportTasks = selectedWaypoint.transportTasks ?? emptyTransportTasks;

		// don't update the waypoint to a colonize/remote mine task if we can't do it on this new target
		if (
			(task == WaypointTaskColonize && !canColonize) ||
			(task == WaypointTaskRemoteMining && !canRemoteMine)
		) {
			task = WaypointTaskNone;
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
				selectedWaypoint.task = WaypointTaskColonize;
				selectedWaypoint.transportTasks = emptyTransportTasks;
			} else if (canRemoteMine) {
				selectedWaypoint.task = WaypointTaskRemoteMining;
				selectedWaypoint.transportTasks = emptyTransportTasks;
			}
		} else {
			selectedWaypoint.position = dest.position;
			selectedWaypoint.targetName = '';
			selectedWaypoint.targetPlayerNum = None;
			selectedWaypoint.targetNum = None;
			selectedWaypoint.targetType = MapObjectTypeNone;
			selectedWaypoint.warpSpeed = warpSpeed;
			selectedWaypoint.task = task;
			selectedWaypoint.transportTasks = transportTasks;
		}

		return true;
	}

	/**
	 * Get the fuel allocated up to a waypoint index accounting for refueling
	 * @param player
	 * @param universe
	 * @param waypointIndex
	 * @returns
	 */
	getFuelAllocated(player: CommandedPlayer, universe: Universe, waypointIndex: number): number {
		let fuelAlreadyAllocated = 0;
		for (let i = 0; i <= waypointIndex; i++) {
			fuelAlreadyAllocated += this.waypoints[i].estFuelUsage ?? 0;
			const wp = this.waypoints[i];
			const target =
				wp.targetType === MapObjectTypePlanet ? universe.getPlanet(wp.targetNum ?? 0) : undefined;
			if (target && this.canFuel(player, target)) {
				// our previous waypoint was a fuel point, reset already allocated fuel to 0
				fuelAlreadyAllocated = 0;
			}
		}
		return fuelAlreadyAllocated;
	}

	/**
	 * Get the fuel allocated up to a waypoint index accounting for refueling
	 * @param player
	 * @param universe
	 * @param waypointIndex
	 * @returns
	 */
	getFuelLeftover(player: CommandedPlayer, universe: Universe, waypointIndex: number): number {
		let fuel = this.fuel;
		for (let i = 0; i <= waypointIndex; i++) {
			fuel -= this.waypoints[i].estFuelUsage ?? 0;
			const wp = this.waypoints[i];
			const target =
				wp.targetType === MapObjectTypePlanet ? universe.getPlanet(wp.targetNum ?? 0) : undefined;
			if (target && this.canFuel(player, target)) {
				// our previous waypoint was a fuel point, reset already allocated fuel to 0
				fuel = this.spec.fuelCapacity ?? 0;
			}
		}
		return fuel;
	}

	/**
	 * Get the fuel allocated up to a waypoint index accounting for refueling
	 * @param player
	 * @param universe
	 * @param waypointIndex
	 * @returns
	 */
	willRunOutOfFuel(player: CommandedPlayer, universe: Universe): boolean {
		let fuel = this.fuel;
		for (let i = 0; i < this.waypoints.length; i++) {
			if (i > 0) {
				const wp1 = this.waypoints[i];
				const fuelUsed = this.getFuelCost(
					universe,
					player.race.spec?.fuelEfficiencyOffset ?? 0,
					wp1.warpSpeed ?? 0,
					distance(this.waypoints[i - 1].position, wp1.position),
					this.spec.cargoCapacity ?? 0
				);
				fuel -= fuelUsed;
			}

			if (fuel < 0) {
				return true;
			}
			const wp = this.waypoints[i];
			const target =
				wp.targetType === MapObjectTypePlanet ? universe.getPlanet(wp.targetNum ?? 0) : undefined;
			if (target && this.canFuel(player, target)) {
				// our previous waypoint was a fuel point, reset already allocated fuel to 0
				fuel = this.spec.fuelCapacity ?? 0;
			}
		}
		return false;
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
		player: CommandedPlayer,
		designFinder: DesignFinder,
		dist: number,
		orbiting: AnyPlanet | undefined,
		dest: WaypointDest,
		fuelAlreadyAllocated: number,
		highestShipMass: number,
		fastestWaypoint: boolean
	): { warpSpeed: number; canColonize: boolean; canRemoteMine: boolean } {
		const mo = dest.mo;

		let canColonize = false;
		let canRemoteMine = false;
		let canJump = false;
		let canFuel = false;
		if (mo && mo.type == MapObjectTypePlanet) {
			const target = mo as PlanetIntel;
			canColonize = this.canColonize(target);
			canRemoteMine = this.canRemoteMine(player, target);
			canJump = this.canJump(player, orbiting, target, dist, highestShipMass);
			canFuel = this.canFuel(player, target);
		}

		// set our warp speed to the most fuel efficient based on our engine idealSpeed (7 for Long Hump 7, 8 for Alpha Drive 8, etc)
		// or the fastest warp we can get there without running out of fuel
		const warpSpeed = canJump
			? StargateWarpSpeed // stargate speed if we can gate
			: canFuel || canColonize || fastestWaypoint // max speed if configured for that, or colonizing
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

	/** Return the highest useful speed less than or equal to a given warp speed
	to reach a given destination.
	 * TODO: Move this to backend so the AI can use it
	 * @param designFinder DesignFinder to find ship designs
	 * @param fuelEfficiencyOffset Sum of all racial fuel cost bonuses/penalties
	 * @param fuelAlreadyAllocated Amount of fuel already allocated for prior waypoints (cannot be spent)
	 * @param dist Distance to destination
	 * @param startSpeed Initial speed to start checking against
	 * @param freeSpeed Maximum free speed of engine
	 * @param maxSafeSpeed Maximum safe speed of engine
	 * @returns The highest useful warp speed we can go at to reach the destination
	*/
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
		let speed = freeSpeed;
		for (let i = speed + 1; i <= maxSafeSpeed; i++) {
			speed = i;
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
	 * Check if a fleet can colonize a planet.
	 * @param target the target planet to check
	 * @returns true if this fleet can colonize this planet
	 */
	canColonize(target: PlanetIntel): boolean {
		return !!(
			this.spec.colonizer &&
			this.cargo.colonists &&
			!owned(target) &&
			(target.spec.terraformedHabitability ?? 0) > 0
		);
	}

	/**
	 * Check if a fleet can remote mine a planet.
	 * @param player The CommandedPlayer commanding the fleet
	 * @param target the target planet to check
	 * @returns true if this fleet can remote mine this planet
	 */
	canRemoteMine(player: CommandedPlayer, target: PlanetIntel): boolean {
		// We can mine unowned planets (as well as self-owned ones)
		return (
			(this.spec.miningRate ?? 0) > 0 &&
			(!owned(target) ||
				(ownedBy(target, player.num) && !!player.race.spec.canRemoteMineOwnPlanets))
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
	canJump(
		player: CommandedPlayer,
		orbiting: AnyPlanet | undefined,
		targetPlanet: PlanetIntel,
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
			return !!destStargateSafe;
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
	 * @param player The fleet player
	 * @param targetPlanet the planet the fleet is targeting
	 * @returns true if the fleet will refuel at this planet
	 */
	canFuel(player: CommandedPlayer, targetPlanet: AnyPlanet | undefined): boolean {
		return !!(
			targetPlanet &&
			owned(targetPlanet) &&
			player.isFriend(targetPlanet.playerNum) &&
			(targetPlanet.spec.dockCapacity ?? 0) != 0
		);
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

	/**
	 *
	 * @param universe
	 * @returns The target for what we should transfer cargo to, based on wp0
	 */
	getCargoTransferTarget(universe: Universe): CargoDest {
		const wp0 = this.waypoints[0];
		if (
			wp0.targetNum == undefined ||
			wp0.targetNum == 0 ||
			wp0.targetType == undefined ||
			wp0.targetType == MapObjectTypeNone
		) {
			// return some salvage at this position
			return universe.getSalvageAtPosition(this);
		}
		switch (wp0.targetType) {
			case MapObjectTypePlanet:
				return universe.getPlanet(wp0.targetNum);
			case MapObjectTypeFleet:
				return universe.getFleet(wp0.targetPlayerNum, wp0.targetNum);
			case MapObjectTypeSalvage:
				return universe.getSalvageAtPosition(this);
			case MapObjectTypeMineralPacket:
				return universe.getMineralPacket(wp0.targetPlayerNum ?? 0, wp0.targetNum);
		}
	}
}

export function getDamagePercentForToken(
	token: ShipToken,
	design: AnyShipDesign | undefined
): number {
	const armor = design?.spec.armor ?? 0;
	const totalArmor = armor * token.quantity;
	const quantityDamaged =
		(token.quantityDamaged ?? 0) > (token.quantity ?? 0)
			? (token.quantity ?? 0)
			: (token.quantityDamaged ?? 0);
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
				.filter((mo) => mo.type === MapObjectTypeFleet)
				.map((mo) => mo as unknown as FleetIntel)
				.filter((f: FleetIntel) => f.freighter)
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
export function idleFleetsFilter(fleet: AnyFleet, showIdleFleetsOnly: boolean): boolean {
	if (!showIdleFleetsOnly) {
		// no filter, show all fleets
		return true;
	}

	// show our fleets that are idle
	if (
		'waypoints' in fleet &&
		fleet.waypoints &&
		fleet.waypoints.length == 1 &&
		fleet.waypoints[0].task == WaypointTaskNone
	) {
		return true;
	}

	// enemy fleet that is moving, show it so players can match idle fleets to moving fleets
	if (!('waypoints' in fleet) && fleet.warpSpeed) {
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
		TransportActionLoadOptimal,
		TransportActionLoadAll,
		TransportActionLoadAmount,
		TransportActionLoadDunnage,
		TransportActionFillPercent,
		TransportActionWaitForPercent
	].indexOf(action) != -1;

export const isUnloadAction = (action: WaypointTaskTransportAction) =>
	[TransportActionUnloadAll, TransportActionUnloadAmount].indexOf(action) != -1;

export const getLocation = (fleet: AnyFleet, universe: Universe) =>
	fleet.orbitingPlanetNum
		? (universe.getPlanet(fleet.orbitingPlanetNum)?.name ?? 'unknown')
		: `Space: (${fleet.position.x}, ${fleet.position.y})`;

export const getDestination = (fleet: Fleet, universe: Universe) => {
	if (fleet.waypoints?.length && fleet.waypoints?.length > 1) {
		return universe.getTargetName(fleet.waypoints[1]);
	}
	return '--';
};

export const getEta = (fleet: Fleet) => {
	if (fleet.waypoints?.length && fleet.waypoints?.length > 1) {
		if (fleet.waypoints[1].warpSpeed === 0) {
			return -1;
		} else if (fleet.waypoints[1].warpSpeed === StargateWarpSpeed) {
			return 1;
		} else {
			return Math.ceil(
				Math.floor(distance(fleet.waypoints[0].position, fleet.waypoints[1].position)) /
					(fleet.waypoints[1].warpSpeed * fleet.waypoints[1].warpSpeed)
			);
		}
	}
	return 0;
};

export function getTokenCount(mo: MapObject) {
	if (mo.type == MapObjectTypeFleet) {
		const fleet = mo as AnyFleet;
		return fleet.tokens ? fleet.tokens.reduce((count, t) => count + t.quantity, 0) : 0;
	}
	return 0;
}

export function hasDestination(mo: MapObject): boolean {
	const fleet = mo.type == MapObjectTypeFleet ? (mo as Fleet) : undefined;
	return (fleet?.waypoints?.length ?? 0) > 1;
}

// get the mass of a fleet or fleetintel
export function getMass(fleet: AnyFleet) {
	if ('mass' in fleet) {
		return fleet.mass ?? 0;
	}
	return fleet.spec?.mass ?? 0;
}

// fleetsSortBy returns a sortBy function for fleets by key. This is used by the fleets report page
// and sorting when cycling through Fleets
export function fleetsSortBy(
	key: string,
	universe: Universe
): ((a: AnyFleet, b: AnyFleet) => number) | undefined {
	switch (key) {
		case 'name':
			return (a, b) => a.name.localeCompare(b.name);
		case 'location':
			return (a, b) => getLocation(a, universe).localeCompare(getLocation(b, universe));
		case 'destination':
			return (a, b) =>
				'waypoints' in a && 'waypoints' in b
					? getDestination(a, universe).localeCompare(getDestination(b, universe))
					: 0;
		case 'eta':
			return (a, b) => ('waypoints' in a && 'waypoints' in b ? getEta(a) - getEta(b) : 0);
		case 'cargo':
			return (a, b) => totalCargo(a.cargo) - totalCargo(b.cargo);
		case 'mass':
			return (a, b) => getMass(a) - getMass(b);
		case 'fuel':
			return (a, b) => ('fuel' in a && 'fuel' in b ? a.fuel - b.fuel : 0);
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
