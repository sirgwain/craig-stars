import type { DesignFinder, Universe } from '$lib/services/Universe';
import { get as pluck } from 'lodash-es';
import { totalCargo, type Cargo } from './Cargo';
import type { Cost } from './Cost';
import { MapObjectType, None, type MovingMapObject } from './MapObject';
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
	battlePlanNum?: number;
	tokens?: ShipToken[];
	waypoints?: Waypoint[];

	repeatOrders?: boolean;
	mass?: number;
	orbitingPlanetNum?: number;
	starbase?: boolean;
	spec?: Spec;
} & MovingMapObject;

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
	torpedoInaccuracyFactor?: number;
	initiative?: number;
	movement?: number;
	powerRating?: number;
	bomber?: number;
	bombs?: number;
	smartBombs?: number;
	retroBombs?: number;
	scanner?: boolean;
	shield?: number;
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
};

export class CommandedFleet implements Fleet {
	id = 0;
	gameId = 0;
	createdAt?: string | undefined;
	updatedAt?: string | undefined;
	readonly type = MapObjectType.Planet;

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
			return (a, b) => `${pluck(a, key)}`.localeCompare(`${pluck(b, key)}`);
	}
}
