import { WaypointTask, type Fleet } from './Fleet';
import type { Player } from './Player';

export type FilterOptions = {
	showIdleFleetsOnly: boolean;
	filterMyDesigns: boolean;
	filterAllyDesigns: boolean;
	filterEnemyDesigns: boolean;
	filterDesigns: string[];
	filterAllyShipClasses: ShipClass[];
	filterEnemyShipClasses: ShipClass[];
};

export type ShipClass = (typeof ShipClasses)[keyof typeof ShipClasses];

// ship classes used by filtering
export const ShipClasses = {
	Colony: 'Colony',
	Freighter: 'Freighter',
	Scout: 'Scout',
	Warship: 'Warship',
	Utility: 'Utility',
	Bomber: 'Bomber',
	Miner: 'Miner',
	FuelTransport: 'Fuel Transport'
} as const;

export function filterFleet(player: Player, fleet: Fleet, options: FilterOptions): boolean {
	return (
		filterIdleFleet(fleet, options.showIdleFleetsOnly) &&
		filterMyDesigns(player, fleet, options.filterMyDesigns, options.filterDesigns) &&
		filterAllyDesigns(player, fleet, options.filterAllyDesigns, options.filterAllyShipClasses) &&
		filterEnemyDesigns(player, fleet, options.filterEnemyDesigns, options.filterEnemyShipClasses)
	);
}

// This shows only your fleets that have no movement orders, and any active enemy ships (so you can match one with the other, if you wish).
export function filterIdleFleet(fleet: Fleet, enabled: boolean): boolean {
	if (!enabled) {
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

export function filterMyDesigns(
	player: Player,
	fleet: Fleet,
	enabled: boolean,
	showDesigns: string[]
): boolean {
	if (!enabled || fleet.playerNum !== player.num) {
		// no filter or not my fleet, show it (it may be filtered by a different function)
		return true;
	}

	return false;
}

export function filterEnemyDesigns(
	player: Player,
	fleet: Fleet,
	enabled: boolean,
	showShipClasses: ShipClass[]
): boolean {
	if (!enabled || !player.isEnemy(fleet.playerNum)) {
		// no filter or not an enemy fleet, show it (it may be filtered by a different function)
		return true;
	}

	return false;
}

export function filterAllyDesigns(
	player: Player,
	fleet: Fleet,
	enabled: boolean,
	showShipClasses: ShipClass[]
): boolean {
	if (!enabled || player.num == fleet.playerNum || !player.isFriend(fleet.playerNum)) {
		// no filter or not a friendly fleet, show it (it may be filtered by a different function)
		return true;
	}

	return false;
}
