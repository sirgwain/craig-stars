import type { DesignFinder, Universe } from '$lib/services/Universe';
import { get as pluck } from 'lodash-es';
import { totalCargo, type Cargo } from './Cargo';
import type { Cost } from './Cost';
import { MapObjectType, None, type MovingMapObject, ownedBy } from './MapObject';
import type { MessageTargetType } from './Message';
import type { ShipDesign } from './ShipDesign';
import type { Engine } from './Tech';
import type { Vector } from './Vector';

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
	immuneToOwnDetonation: boolean;
	mineLayingRateByMineType?: null;
	weaponSlots?: null;
	purposes?: any;
	totalShips: number;
	massEmpty: number;
	basePacketSpeed: number;
	safePacketSpeed: number;
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
	colonizer?: number;
	canLayMines?: number;
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

	getMinimalWarp(dist: number, idealSpeed: number): number {
		let speed = idealSpeed;
		const freeSpeed = this.spec?.engine?.freeSpeed ?? 1;

		const yearsAtIdealSpeed = dist / (idealSpeed * idealSpeed);
		for (let i = idealSpeed; i > freeSpeed; i--) {
			const yearsAtSpeed = dist / (i * i);
			// it takes the same time to go slower, so go slower
			if (Math.ceil(yearsAtIdealSpeed) == Math.ceil(yearsAtSpeed)) {
				speed = i;
			}
		}

		return speed;
	}

	// get the max warp we have fuel for to make it to the dest
	getMaxWarp(dist: number, designFinder: DesignFinder, fuelEfficiencyOffset: number): number {
		// start at one above free speed
		const freeSpeed = this.spec?.engine?.freeSpeed ?? 1;
		let speed: number;
		for (speed = freeSpeed + 1; speed < 9; speed++) {
			const fuelUsed = this.getFuelCost(
				designFinder,
				fuelEfficiencyOffset,
				speed,
				dist,
				this.spec.cargoCapacity ?? 0
			);
			if (fuelUsed > this.fuel) {
				speed--;
				break;
			}
		}

		const idealSpeed = this.spec?.engine?.idealSpeed ?? 1;

		// if we are using a ramscoop, make sure we at least go the ideal
		// speed of the engine. If we run out, oh well, it'll drop to
		// the free speed
		if (freeSpeed > 1 && speed < idealSpeed) {
			speed = idealSpeed;
		}

		// don't go faster than we need
		return this.getMinimalWarp(dist, speed);
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
