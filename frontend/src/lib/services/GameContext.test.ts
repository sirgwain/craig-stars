import { clone, create } from '@bufbuild/protobuf';
import { get } from 'svelte/store';
import { beforeEach, expect, it, vi } from 'vitest';
import {
	FleetSchema,
	MapObjectType,
	WaypointDestSchema,
	WaypointSchema,
	type Fleet,
	type FleetOrders
} from '#lib/types/cs-proto.js';
import type { CS } from '#lib/wasm.js';
import { CommandedPlayer } from '#lib/types/Player.js';
import { errors } from './Errors';
import { FullGame } from './FullGame';
import { createGameContext } from './GameContext';
import { Universe } from './Universe';

// Tests for how GameContext edits fleet waypoints and saves fleet orders. Edits draw locally right
// away and save to the server in the background. Slow servers (like players see on phones) must not
// make the waypoints lag, get added out of order, or get overwritten by stale responses.
//
// The server and wasm calls are mocked. Tests hold calls open with deferred() promises to choose
// the order in which they finish.

const clients = vi.hoisted(() => ({
	fleetClient: { updateFleetOrders: vi.fn(), renameFleet: vi.fn() },
	playerClient: { submitTurn: vi.fn() }
}));
vi.mock('$app/navigation', () => ({ goto: vi.fn() }));
vi.mock('./connect', () => ({
	...clients,
	gameClient: {},
	minefieldClient: {},
	planetClient: {},
	battlePlanClient: {},
	productionPlanClient: {},
	shipDesignClient: {},
	transportPlanClient: {}
}));

// A promise the test resolves or rejects by hand, to simulate a slow server or wasm call.
function deferred<T>() {
	let resolve!: (value: T) => void;
	let reject!: (error: Error) => void;
	const promise = new Promise<T>((res, rej) => {
		resolve = res;
		reject = rej;
	});
	return { promise, resolve, reject };
}

// A fleet at (10, 10) with only its starting waypoint.
function fleet(num = 1): Fleet {
	return create(FleetSchema, {
		mapObject: { type: MapObjectType.FLEET, playerNum: 1, num, position: { x: 10, y: 10 } },
		fleetOrders: { waypoints: [{ position: { x: 10, y: 10 } }] }
	});
}

// Waypoints are identified by x position, so a route reads like [10, 20, 30].
const dest = (x: number) => create(WaypointDestSchema, { position: { x, y: 10 } });
// The server's reply to a save that accepted these orders.
const response = (orders: FleetOrders, num = 1) => ({
	fleet: create(FleetSchema, { ...fleet(num), fleetOrders: orders })
});

// Create a context commanding fleet 1. The wasm addWaypoint mock inserts the new waypoint after the
// selected one, the same way the real one does.
async function setup() {
	const u = new Universe();
	u.fleets = [fleet(1), fleet(2)];
	const p = new CommandedPlayer();
	p.num = 1;
	const fg = new FullGame();
	fg.id = 1n;
	const add = vi.fn(async ({ fleet, dest, currentSelectedWaypointIndex }) => {
		const updated = clone(FleetSchema, fleet);
		const index = currentSelectedWaypointIndex + 1;
		updated.fleetOrders!.waypoints.splice(
			index,
			0,
			create(WaypointSchema, { position: dest.position })
		);
		return { fleet: updated, index };
	});
	const cs = {
		wasmService: {
			computeRaceSpec: vi.fn(async () => ({})),
			setPlayer: vi.fn(),
			setDesigns: vi.fn(),
			setIntels: vi.fn(),
			addWaypoint: add
		}
	} as unknown as CS;
	const context = await createGameContext(cs, fg, p, u);
	context.commandMapObject(u.fleets[0]);
	return { context, u, cs, add };
}

beforeEach(() => {
	vi.resetAllMocks();
	errors.set([]);
	vi.stubGlobal('localStorage', { getItem: () => null, setItem: () => {} });
	clients.playerClient.submitTurn.mockResolvedValue({});
});

// A tap draws and selects the waypoint before the save finishes. The save sends its own copy of the
// orders, so later edits (like a warp speed change) don't change what was sent.
it('draws and selects a waypoint without waiting for the server', async () => {
	const save = deferred<{ fleet: Fleet }>();
	clients.fleetClient.updateFleetOrders.mockReturnValue(save.promise);
	const { context, u } = await setup();
	expect(await context.addWaypoint(dest(20), false)).toBe(true);
	expect(get(context.currentSelectedWaypointIndex)).toBe(1);
	expect(u.fleets[0].fleetOrders!.waypoints).toHaveLength(2);
	await vi.waitFor(() => expect(clients.fleetClient.updateFleetOrders).toHaveBeenCalledOnce());
	const orders = clients.fleetClient.updateFleetOrders.mock.calls[0][0].fleetOrders;
	get(context.commandedFleet)!.fleetOrders.waypoints[1].warpSpeed = 9;
	expect(orders.waypoints[1].warpSpeed).toBe(0);
	save.resolve(response(orders));
	await vi.waitFor(() =>
		expect(get(context.commandedFleet)!.fleetOrders.waypoints[1].warpSpeed).toBe(0)
	);
});

// The out-of-order bug: two quick taps. The second tap's calculation has to wait for the first
// one, so it inserts after the first tap's waypoint and not after the old selection.
it('uses the previous tap selection even when local calculations overlap', async () => {
	clients.fleetClient.updateFleetOrders.mockImplementation(async ({ fleetOrders }) =>
		response(fleetOrders)
	);
	const { context, add } = await setup();
	const calculation = deferred<Awaited<ReturnType<typeof add>>>();
	const original = add.getMockImplementation()!;
	add.mockImplementationOnce(async (request) => {
		const result = await original(request);
		await calculation.promise;
		return result;
	});
	const first = context.addWaypoint(dest(20), false);
	const second = context.addWaypoint(dest(30), false);
	await vi.waitFor(() => expect(add).toHaveBeenCalledOnce());
	calculation.resolve({ fleet: fleet(), index: 1 });
	expect(await first).toBe(true);
	expect(await second).toBe(true);
	expect(add.mock.calls[1][0].currentSelectedWaypointIndex).toBe(1);
	expect(get(context.commandedFleet)!.fleetOrders.waypoints.map((wp) => wp.position!.x)).toEqual([
		10, 20, 30
	]);
});

// Saves for a fleet go out one at a time. While a save is in flight, newer taps replace any queued
// save, since each save sends the full orders. A response for older orders must not undo newer
// local edits or the player's waypoint selection.
it('serializes saves and ignores older responses while newer orders are visible', async () => {
	const first = deferred<{ fleet: Fleet }>();
	const second = deferred<{ fleet: Fleet }>();
	clients.fleetClient.updateFleetOrders
		.mockReturnValueOnce(first.promise)
		.mockReturnValueOnce(second.promise);
	const { context } = await setup();
	await context.addWaypoint(dest(20), false);
	await context.addWaypoint(dest(30), false);
	await context.addWaypoint(dest(40), false);
	expect(clients.fleetClient.updateFleetOrders).toHaveBeenCalledOnce();
	expect(get(context.currentSelectedWaypointIndex)).toBe(3);
	first.resolve(response(clients.fleetClient.updateFleetOrders.mock.calls[0][0].fleetOrders));
	await vi.waitFor(() => expect(clients.fleetClient.updateFleetOrders).toHaveBeenCalledTimes(2));
	// The save for the second tap is superseded by the third tap's full orders.
	expect(clients.fleetClient.updateFleetOrders.mock.calls[1][0].fleetOrders.waypoints).toHaveLength(
		4
	);
	expect(get(context.commandedFleet)!.fleetOrders.waypoints).toHaveLength(4);
	expect(get(context.currentSelectedWaypointIndex)).toBe(3);
	// Selecting an earlier waypoint while saving must survive the final response too.
	context.selectWaypoint(get(context.commandedFleet)!.fleetOrders.waypoints[0]);
	second.resolve(response(clients.fleetClient.updateFleetOrders.mock.calls[1][0].fleetOrders));
	await context.submitTurn();
	expect(get(context.currentSelectedWaypointIndex)).toBe(0);
	expect(clients.fleetClient.updateFleetOrders).toHaveBeenCalledTimes(2);
});

// Each fleet has its own save queue. A save finishing for a fleet the player is no longer commanding
// updates that fleet without changing the commanded fleet.
it('allows different fleets to save independently without changing the commanded fleet', async () => {
	const first = deferred<{ fleet: Fleet }>();
	const second = deferred<{ fleet: Fleet }>();
	clients.fleetClient.updateFleetOrders
		.mockReturnValueOnce(first.promise)
		.mockReturnValueOnce(second.promise);
	const { context, u } = await setup();
	await context.addWaypoint(dest(20), false);
	context.commandMapObject(u.fleets[1]);
	await context.addWaypoint(dest(30), false);
	await vi.waitFor(() => expect(clients.fleetClient.updateFleetOrders).toHaveBeenCalledTimes(2));
	first.resolve(response(clients.fleetClient.updateFleetOrders.mock.calls[0][0].fleetOrders));
	second.resolve(response(clients.fleetClient.updateFleetOrders.mock.calls[1][0].fleetOrders, 2));
	await context.submitTurn();
	expect(get(context.commandedFleet)!.mapObject.num).toBe(2);
	expect(get(context.currentSelectedWaypointIndex)).toBe(1);
	expect(u.fleets.every((f) => f.fleetOrders!.waypoints.length === 2)).toBe(true);
});

// A failed save only matters if it was the latest one: then the route goes back to the last orders
// the server accepted and an error is shown. An older save failing is ignored because the newer
// save carries the full orders. Saving works again after a rollback.
it('reports and rolls back only failures of the latest save', async () => {
	const first = deferred<{ fleet: Fleet }>();
	const third = deferred<{ fleet: Fleet }>();
	clients.fleetClient.updateFleetOrders
		.mockReturnValueOnce(first.promise)
		.mockImplementationOnce(async ({ fleetOrders }) => response(fleetOrders))
		.mockReturnValueOnce(third.promise)
		.mockImplementation(async ({ fleetOrders }) => response(fleetOrders));
	const { context } = await setup();
	const xs = () => get(context.commandedFleet)!.fleetOrders.waypoints.map((wp) => wp.position!.x);
	await context.addWaypoint(dest(20), false);
	await context.addWaypoint(dest(30), false);
	// A superseded save failing is silent; the newer save carries the full orders.
	first.reject(new Error('Temporary failure'));
	await context.submitTurn();
	expect(get(errors)).toHaveLength(0);
	expect(xs()).toEqual([10, 20, 30]);

	await context.addWaypoint(dest(40), false);
	third.reject(new Error('Offline'));
	await vi.waitFor(() => expect(get(errors)).toHaveLength(1));
	expect(get(errors)[0].error).toContain('Restored the last saved orders');
	expect(xs()).toEqual([10, 20, 30]);
	expect(get(context.currentSelectedWaypointIndex)).toBe(2);

	await context.addWaypoint(dest(50), false);
	await context.submitTurn();
	expect(xs()).toEqual([10, 20, 30, 50]);
});

// Submitting waits for pending waypoint edits and saves, so the server never generates the turn
// with old orders. If a pending save fails, the turn isn't submitted.
it('submits the turn only after pending edits save successfully', async () => {
	const save = deferred<{ fleet: Fleet }>();
	const failed = deferred<{ fleet: Fleet }>();
	clients.fleetClient.updateFleetOrders
		.mockReturnValueOnce(save.promise)
		.mockReturnValueOnce(failed.promise);
	const { context } = await setup();
	const added = context.addWaypoint(dest(20), false);
	const submitted = context.submitTurn();
	await added;
	expect(clients.playerClient.submitTurn).not.toHaveBeenCalled();
	await vi.waitFor(() => expect(clients.fleetClient.updateFleetOrders).toHaveBeenCalledOnce());
	save.resolve(response(clients.fleetClient.updateFleetOrders.mock.calls[0][0].fleetOrders));
	await submitted;
	expect(clients.playerClient.submitTurn).toHaveBeenCalledOnce();

	await context.addWaypoint(dest(30), false);
	const rejected = expect(context.submitTurn()).rejects.toThrow('Offline');
	// Let submission start waiting on the save before rejecting it.
	await new Promise((resolve) => setTimeout(resolve, 0));
	failed.reject(new Error('Offline'));
	await rejected;
	expect(clients.playerClient.submitTurn).toHaveBeenCalledOnce();
});

// After the context is reset for a new turn, queued saves from the old turn aren't sent, and a late
// response from the old turn doesn't change the new turn's fleets.
it('discards queued saves and late responses after resetting the context', async () => {
	const save = deferred<{ fleet: Fleet }>();
	clients.fleetClient.updateFleetOrders.mockReturnValue(save.promise);
	const { context } = await setup();
	await context.addWaypoint(dest(20), false);
	await context.addWaypoint(dest(30), false);
	const u = new Universe();
	u.fleets = [fleet()];
	const p = new CommandedPlayer();
	p.num = 1;
	await context.resetContext(new FullGame(), p, u);
	context.commandMapObject(u.fleets[0]);
	save.resolve(response(clients.fleetClient.updateFleetOrders.mock.calls[0][0].fleetOrders));
	await new Promise((resolve) => setTimeout(resolve, 0));
	expect(clients.fleetClient.updateFleetOrders).toHaveBeenCalledOnce();
	expect(get(context.commandedFleet)!.fleetOrders.waypoints).toHaveLength(1);
	expect(get(context.currentSelectedWaypointIndex)).toBe(0);
});

// A delete tapped while an add is still being calculated waits for the add, so it deletes the new
// waypoint and not whichever waypoint was selected before.
it('queues a deletion behind a pending waypoint calculation', async () => {
	clients.fleetClient.updateFleetOrders.mockImplementation(async ({ fleetOrders }) =>
		response(fleetOrders)
	);
	const { context, add } = await setup();
	const calculation = deferred<void>();
	const original = add.getMockImplementation()!;
	add.mockImplementationOnce(async (request) => {
		const result = await original(request);
		await calculation.promise;
		return result;
	});
	const added = context.addWaypoint(dest(20), false);
	const deleted = context.deleteWaypoint();
	calculation.resolve();
	await added;
	await deleted;
	await context.submitTurn();
	expect(get(context.commandedFleet)!.fleetOrders.waypoints.map((wp) => wp.position!.x)).toEqual([
		10
	]);
});

// If a failed save rolls the route back while another add is being calculated, that add was based
// on orders that were thrown away. It is dropped instead of putting the rolled-back waypoints back.
it('drops a waypoint calculated from orders that were rolled back', async () => {
	const save = deferred<{ fleet: Fleet }>();
	clients.fleetClient.updateFleetOrders.mockReturnValueOnce(save.promise);
	const { context, add } = await setup();
	await context.addWaypoint(dest(20), false);
	const calculation = deferred<void>();
	const original = add.getMockImplementation()!;
	add.mockImplementationOnce(async (request) => {
		const result = await original(request);
		await calculation.promise;
		return result;
	});
	const second = context.addWaypoint(dest(30), false);
	await vi.waitFor(() => expect(add).toHaveBeenCalledTimes(2));
	save.reject(new Error('Offline'));
	await vi.waitFor(() => expect(get(errors)).toHaveLength(1));
	calculation.resolve();
	expect(await second).toBe(false);
	expect(get(context.commandedFleet)!.fleetOrders.waypoints.map((wp) => wp.position!.x)).toEqual([
		10
	]);
	expect(clients.fleetClient.updateFleetOrders).toHaveBeenCalledOnce();
});
