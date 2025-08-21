import type { Universe } from '$lib/services/Universe';
import {
	CargoSchema,
	type Fleet,
	type FleetOrders,
	FleetOrdersSchema,
	FleetSchema,
	FleetSpecSchema,
	MapObjectSchema,
	MapObjectType,
	type Planet,
	ResourceType,
	type ShipDesign,
	type ShipToken,
	type ShipTokenJson,
	type Vector,
	VectorSchema,
	WaypointTask,
	WaypointTaskTransportAction,
	type WaypointTransportTasks,
	WaypointTransportTasksSchema
} from '$lib/types/cs-proto';
import { create, merge, type UnknownField } from '@bufbuild/protobuf';
import { get as pluck } from 'lodash-es';
import { type CargoType, totalCargo } from './Cargo';
import type { CargoDest } from './CargoTransferRequest.svelte';
import { None, StargateWarpSpeed } from './Consts';
import { type MapObjectLike, owned } from './MapObject';
import type { CommandedPlayer } from './Player';
import { distance, emptyVector } from './Vector';

export const WaypointTasks: WaypointTask[] = [
	WaypointTask.UNSPECIFIED,
	WaypointTask.TRANSPORT,
	WaypointTask.COLONIZE,
	WaypointTask.REMOTE_MINING,
	WaypointTask.MERGE_WITH_FLEET,
	WaypointTask.SCRAP_FLEET,
	WaypointTask.LAY_MINEFIELD,
	WaypointTask.PATROL,
	WaypointTask.ROUTE,
	WaypointTask.TRANSFER_FLEET
] as const;

export function emptyTransportTasks(): WaypointTransportTasks {
	return create(WaypointTransportTasksSchema, {
		fuel: {
			action: WaypointTaskTransportAction.UNSPECIFIED
		},
		ironium: {
			action: WaypointTaskTransportAction.UNSPECIFIED
		},
		boranium: {
			action: WaypointTaskTransportAction.UNSPECIFIED
		},
		germanium: {
			action: WaypointTaskTransportAction.UNSPECIFIED
		},
		colonists: {
			action: WaypointTaskTransportAction.UNSPECIFIED
		}
	});
}

export class CommandedFleet implements Fleet {
	$typeName: 'craig_stars.v1.Fleet';
	$unknown?: UnknownField[] | undefined;
	readonly type = MapObjectType.FLEET;

	mapObject = create(MapObjectSchema);
	fleetOrders: FleetOrders = create(FleetOrdersSchema);
	age = 0;
	baseName = '';
	cargo = create(CargoSchema);
	damage = 0;
	fuel = 0;
	heading = create(VectorSchema);
	mass = 0;
	orbitingPlanetNum = None;
	planetNum = 0;
	previousPosition?: Vector | undefined;
	starbase = false;
	tokens: ShipToken[] = [];
	warpSpeed = 0;
	spec = create(FleetSpecSchema);

	constructor(data?: Fleet) {
		this.$typeName = 'craig_stars.v1.Fleet';
		if (data) {
			merge(FleetSchema, this, data);
		}
	}

	getWaypointMapObjects(universe: Universe): MapObjectLike[] {
		return this.fleetOrders?.waypoints
			.map((wp) => universe.getMapObject(wp.mapObjectTarget))
			.filter((m): m is MapObjectLike => !!m);
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
		for (let i = 0; i < this.fleetOrders?.waypoints.length; i++) {
			if (i > 0) {
				const wp1 = this.fleetOrders?.waypoints[i];
				fuel -= wp1.estFuelUsage ?? 0;
			}

			if (fuel < 0) {
				return true;
			}
			const wp = this.fleetOrders?.waypoints[i];
			const target =
				wp.mapObjectTarget?.targetType === MapObjectType.PLANET
					? universe.getPlanet(wp.mapObjectTarget?.targetNum)
					: undefined;
			if (target && this.canFuel(player, target)) {
				// our previous waypoint was a fuel point, reset already allocated fuel to 0
				fuel = this.spec.shipDesignSpec?.fuelCapacity ?? 0;
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
	canFuel(player: CommandedPlayer, targetPlanet: Planet | undefined): boolean {
		return !!(
			targetPlanet &&
			owned(targetPlanet) &&
			player.isFriend(targetPlanet.mapObject?.playerNum ?? None) &&
			(targetPlanet.spec?.planetStarbaseSpec?.dockCapacity ?? 0) != 0
		);
	}

	/**
	 *
	 * @returns The total number of mines laid per year for all types of minefields this fleet can lay
	 */
	getTotalMinesLaidPerYear() {
		if (this.spec.shipDesignSpec?.mineLayingRateByMineType) {
			return Object.values(this.spec.shipDesignSpec?.mineLayingRateByMineType).reduce(
				(count, n) => count + n,
				0
			);
		}
		return 0;
	}

	/**
	 *
	 * @param universe
	 * @returns The target for what we should transfer cargo to, based on wp0
	 */
	getCargoTransferTarget(universe: Universe): CargoDest {
		const wp0 = this.fleetOrders?.waypoints[0];
		if (
			wp0.mapObjectTarget?.targetNum == undefined ||
			wp0.mapObjectTarget?.targetNum == 0 ||
			wp0.mapObjectTarget?.targetType == undefined ||
			wp0.mapObjectTarget?.targetType === MapObjectType.UNSPECIFIED
		) {
			// return some salvage at this position
			return universe.getSalvageAtPosition(this.mapObject?.position);
		}
		switch (wp0.mapObjectTarget?.targetType) {
			case MapObjectType.PLANET:
				return universe.getPlanet(wp0.mapObjectTarget?.targetNum);
			case MapObjectType.FLEET:
				return universe.getFleet(
					wp0.mapObjectTarget?.targetPlayerNum,
					wp0.mapObjectTarget?.targetNum
				);
			case MapObjectType.SALVAGE:
				return universe.getSalvageAtPosition(this.mapObject?.position);
			case MapObjectType.MINERAL_PACKET:
				return universe.getMineralPacket(
					wp0.mapObjectTarget?.targetPlayerNum ?? 0,
					wp0.mapObjectTarget?.targetNum
				);
		}
	}
}

export function getDamagePercentForToken(token: ShipToken, design: ShipDesign | undefined): number {
	const armor = design?.spec?.armor ?? 0;
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
	return (fleet.spec?.shipDesignSpec?.cargoCapacity ?? 0) > 0;
}

// true if this fleet can transfer this cargo type
// used to stop stealing colonists or fuel
export function canTransferCargoType(fleet: Fleet, dest: CargoDest, cargoType: CargoType): boolean {
	if (dest?.mapObject?.type === MapObjectType.FLEET) {
		switch (cargoType) {
			case ResourceType.COLONISTS:
			case ResourceType.FUEL:
				return fleet.mapObject?.playerNum === dest?.mapObject?.playerNum;
		}
	}
	if (
		dest?.mapObject?.type === MapObjectType.SALVAGE ||
		dest?.mapObject?.type === MapObjectType.MINERAL_PACKET
	) {
		switch (cargoType) {
			case ResourceType.COLONISTS:
			case ResourceType.FUEL:
				return false;
		}
	}

	return true;
}

// true if this fleet can load cargo
export function canLoadFuelOrCargo(fleet: Fleet, dest: CargoDest): boolean {
	// can always load from our own stuff, or empty stuff
	if (
		dest?.mapObject?.playerNum === fleet.mapObject?.playerNum ||
		dest?.mapObject?.playerNum === None
	) {
		return true;
	}

	if (dest?.mapObject?.type === MapObjectType.PLANET) {
		// we can only load from this planet if we can steal planet cargo
		return !!fleet.spec?.shipDesignSpec?.canStealPlanetCargo;
	}

	if (dest?.mapObject?.type === MapObjectType.FLEET) {
		// we can only load from this fleet if we can steal fleet cargo
		return !!fleet.spec?.shipDesignSpec?.canStealFleetCargo;
	}
	return false;
}

export const isLoadAction = (action: WaypointTaskTransportAction) =>
	[
		WaypointTaskTransportAction.LOAD_OPTIMAL,
		WaypointTaskTransportAction.LOAD_ALL,
		WaypointTaskTransportAction.LOAD_AMOUNT,
		WaypointTaskTransportAction.LOAD_DUNNAGE,
		WaypointTaskTransportAction.FILL_PERCENT,
		WaypointTaskTransportAction.WAIT_FOR_PERCENT
	].indexOf(action) != -1;

export const isUnloadAction = (action: WaypointTaskTransportAction) =>
	[WaypointTaskTransportAction.UNLOAD_ALL, WaypointTaskTransportAction.UNLOAD_AMOUNT].indexOf(
		action
	) != -1;

export const getLocation = (fleet: Fleet, universe: Universe) =>
	fleet.orbitingPlanetNum
		? (universe.getPlanet(fleet.orbitingPlanetNum)?.mapObject?.name ?? 'unknown')
		: `Space: (${fleet.mapObject?.position?.x ?? 0}, ${fleet.mapObject?.position?.y ?? 0})`;

export const getDestination = (fleet: Fleet, universe: Universe) => {
	const wps = fleet.fleetOrders?.waypoints ?? [];
	if (wps.length > 1) {
		return universe.getTargetName(wps[1]);
	}
	return '--';
};

export const getEta = (fleet: Fleet) => {
	const wps = fleet.fleetOrders?.waypoints ?? [];
	if (wps.length > 1) {
		if (wps[1].warpSpeed === 0) {
			return -1;
		} else if (wps[1].warpSpeed === StargateWarpSpeed) {
			return 1;
		} else {
			return Math.ceil(
				Math.floor(distance(wps[0].position ?? emptyVector(), wps[1].position)) /
					(wps[1].warpSpeed * wps[1].warpSpeed)
			);
		}
	}
	return 0;
};

export function getTokenCount(mo: MapObjectLike) {
	if (mo.mapObject?.type === MapObjectType.FLEET) {
		const fleet = mo as Fleet;
		return fleet.tokens ? fleet.tokens.reduce((count, t) => count + t.quantity, 0) : 0;
	}
	return 0;
}

export function hasDestination(mo: MapObjectLike): boolean {
	const fleet = mo.mapObject?.type === MapObjectType.FLEET ? (mo as Fleet) : undefined;
	return (fleet?.fleetOrders?.waypoints?.length ?? 0) > 1;
}

// get the mass of a fleet or Fleet
export function getMass(fleet: Fleet): number {
	return fleet.spec?.shipDesignSpec?.mass ?? 0;
}

// fleetsSortBy returns a sortBy function for fleets by key. This is used by the fleets report page
// and sorting when cycling through Fleets
export function fleetsSortBy(
	key: string,
	universe: Universe
): ((a: Fleet, b: Fleet) => number) | undefined {
	switch (key) {
		case 'name':
			return (a, b) => (a.mapObject?.name ?? '').localeCompare(b.mapObject?.name ?? '');
		case 'location':
			return (a, b) => getLocation(a, universe).localeCompare(getLocation(b, universe));
		case 'destination':
			return (a, b) =>
				'fleetOrders' in a && 'fleetOrders' in b
					? getDestination(a, universe).localeCompare(getDestination(b, universe))
					: 0;
		case 'task':
			return (a, b) =>
				'fleetOrders' in a && 'fleetOrders' in b
					? (a.fleetOrders?.waypoints[a.fleetOrders?.waypoints.length - 1].task ?? 0) -
						(b.fleetOrders?.waypoints[b.fleetOrders?.waypoints.length - 1].task ?? 0)
					: 0;
		case 'eta':
			return (a, b) => ('fleetOrders' in a && 'fleetOrders' in b ? getEta(a) - getEta(b) : 0);
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
export function moveDamagedTokens(
	srcToken: ShipTokenJson,
	destToken: ShipTokenJson,
	quantity: number
) {
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
