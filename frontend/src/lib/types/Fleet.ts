import type { AnyFleet, AnyPlanet, AnyShipDesign, Universe } from '$lib/services/Universe';
import { get as pluck } from 'lodash-es';
import { totalCargo } from './Cargo';
import type { CargoDest } from './CargoTransferRequest.svelte';
import { owned } from './MapObject';
import type { CommandedPlayer } from './Player';
import { distance } from './Vector';
import {
	type Cargo,
	type CargoType,
	Colonists,
	type Fleet,
	type FleetSpec,
	Fuel,
	type MapObject,
	MapObjectTypeFleet,
	MapObjectTypeMineralPacket,
	MapObjectTypeNone,
	MapObjectTypePlanet,
	MapObjectTypeSalvage,
	None,
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
				fuel -= wp1.estFuelUsage ?? 0;
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
			(targetPlanet.spec?.dockCapacity ?? 0) != 0
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
export function canTransferCargo(fleet: Fleet): boolean {
	return (fleet.spec?.cargoCapacity ?? 0) > 0;
}

// true if this fleet can transfer this cargo type
// used to stop stealing colonists or fuel
export function canTransferCargoType(fleet: Fleet, dest: CargoDest, cargoType: CargoType): boolean {
	if (dest?.type === MapObjectTypeFleet) {
		switch (cargoType) {
			case Colonists:
			case Fuel:
				return fleet.playerNum === dest?.playerNum;
		}
	}
	if (dest?.type === MapObjectTypeSalvage || dest?.type == MapObjectTypeMineralPacket) {
		switch (cargoType) {
			case Colonists:
			case Fuel:
				return false;
		}
	}

	return true;
}

// true if this fleet can load cargo
export function canLoadFuelOrCargo(fleet: Fleet, dest: CargoDest): boolean {
	// can always load from our own stuff, or empty stuff
	if (dest?.playerNum === fleet.playerNum || dest?.playerNum === None) {
		return true;
	}

	if (dest?.type === MapObjectTypePlanet) {
		// we can only load from this planet if we can steal planet cargo
		return !!fleet.spec.canStealPlanetCargo;
	}

	if (dest?.type === MapObjectTypeFleet) {
		// we can only load from this fleet if we can steal fleet cargo
		return !!fleet.spec.canStealFleetCargo;
	}
	return false;
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
		case 'task':
			return (a, b) =>
				'waypoints' in a && 'waypoints' in b
					? (a.waypoints[a.waypoints.length - 1].task ?? '').localeCompare(
							b.waypoints[b.waypoints.length - 1].task ?? ''
						)
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
