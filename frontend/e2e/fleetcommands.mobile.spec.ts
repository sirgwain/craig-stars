import { create, fromJson, toJson } from '@bufbuild/protobuf';
import { expect, type Locator, type Page } from '@playwright/test';
import {
	UpdateFleetOrdersRequestSchema,
	UpdateFleetOrdersResponseSchema,
	type UpdateFleetOrdersRequestJson,
	type UpdateFleetOrdersResponseJson
} from '../src/lib/protogen/craig_stars/v1/fleetservice_pb';
import { StargateWarpSpeed } from '../src/lib/types/Consts';
import { WaypointTask, WaypointTaskTransportAction } from '../src/lib/types/cs-proto';
import { tapMapObject } from './helpers/scanner';
import { test } from './setup';

test(
	'toggling the drawer preserves a fleet on a repeating round trip',
	{ tag: '@safari' },
	async ({ testGamePage }) => {
		const { page, universe, gameId } = await testGamePage('Kitchen Sink');
		const homeworld = universe.planets[0];
		const target = universe.planets[1];
		const drawer = page.locator('[data-type="command-drawer"]');
		const toggle = drawer.getByRole('button', { name: 'show command pane button' });
		const endpoint = '/craig_stars.v1.FleetService/UpdateFleetOrders';

		await tapMapObject(page, homeworld);
		await page.locator('#add-waypoint').tap();
		const saved = page.waitForResponse((response) => response.url().endsWith(endpoint));
		await tapMapObject(page, target);
		const { fleet } = fromJson(
			UpdateFleetOrdersResponseSchema,
			(await (await saved).json()) as UpdateFleetOrdersResponseJson
		);
		// Seed the return leg and repeat flag to focus on drawer and waypoint selection.
		const orders = fleet!.fleetOrders!;
		orders.waypoints.push({ ...orders.waypoints[0], warpSpeed: 6 });
		orders.repeatOrders = true;
		const response = await page.request.post(`/api/grpc${endpoint}`, {
			data: toJson(
				UpdateFleetOrdersRequestSchema,
				create(UpdateFleetOrdersRequestSchema, {
					gameId,
					fleetNum: fleet!.mapObject!.num,
					fleetOrders: orders
				})
			)
		});
		expect(response.ok()).toBe(true);
		await page.reload();
		await tapMapObject(page, homeworld);
		await toggle.tap();

		const waypoints = drawer.locator('[data-type="command-tile"][data-id="Fleet Waypoints"]');
		const waypointTarget = drawer
			.locator('[data-type="command-tile"][data-id="Waypoint Task"]')
			.getByRole('combobox')
			.first();
		await expect(waypoints.locator('li')).toHaveText([
			homeworld.mapObject!.name,
			target.mapObject!.name,
			homeworld.mapObject!.name
		]);
		await expect(waypoints.getByRole('checkbox', { name: 'Repeat Orders' })).toBeChecked();

		// Selecting the return waypoint shows the homeworld in the collapsed summary. It must
		// not replace the commanded fleet, or reset a different selected waypoint, on a toggle.
		await waypoints.getByRole('button', { name: target.mapObject!.name, exact: true }).tap();
		await expect(waypointTarget.locator('option:checked')).toHaveText(target.mapObject!.name);
		await waypoints
			.getByRole('button', { name: homeworld.mapObject!.name, exact: true })
			.last()
			.tap();
		for (const selected of [2, 1]) {
			if (selected === 1) {
				await waypoints.getByRole('button', { name: target.mapObject!.name, exact: true }).tap();
			}
			await expect(waypointTarget.locator('option:checked')).toHaveText(
				selected === 1 ? target.mapObject!.name : homeworld.mapObject!.name
			);
			await toggle.tap();
			await expect(drawer).toHaveAttribute('data-open', 'false');
			await expect(toggle).toHaveCount(1);
			await toggle.tap();
			await expect(drawer).toHaveAttribute('data-open', 'true');
			await expect(toggle).toHaveCount(1);
			await expect(toggle).toHaveText(fleet!.mapObject!.name);
			await expect(waypoints.locator('li').nth(selected)).toHaveClass(/bg-primary-focus/);
			await expect(waypoints.getByRole('checkbox', { name: 'Repeat Orders' })).toBeChecked();
			await expect(waypointTarget.locator('option:checked')).toHaveText(
				selected === 1 ? target.mapObject!.name : homeworld.mapObject!.name
			);
		}
	}
);

test('editing transport orders preserves the plan and keeps the command drawer open', async ({
	testGamePage
}) => {
	const { page, universe, player } = await testGamePage('Kitchen Sink');
	const plan = player.playerPlans!.transportPlans.find((p) => p.name === 'Load Colonists')!;
	const drawer = page.locator('[data-type="command-drawer"]');
	await tapMapObject(page, universe.planets[0]);
	await drawer.getByRole('button', { name: 'show command pane button' }).tap();
	const task = drawer.locator('[data-type="command-tile"][data-id="Waypoint Task"]');
	const saveEndpoint = '/craig_stars.v1.FleetService/UpdateFleetOrders';
	let saved = page.waitForResponse((response) => response.url().endsWith(saveEndpoint));
	await task
		.locator('[data-type="select-waypoint-task"]')
		.selectOption(String(WaypointTask.TRANSPORT));
	await saved;
	saved = page.waitForResponse((response) => response.url().endsWith(saveEndpoint));
	await task.getByRole('combobox').last().selectOption(String(plan.num));
	await saved;

	const dialog = page.locator('[data-type="dialog"][data-id="Transport Orders"]');
	await task.getByRole('button', { name: 'Edit transport orders' }).tap();
	await expect(dialog.locator('[name="actionColonists"]')).toHaveValue(
		String(WaypointTaskTransportAction.LOAD_ALL)
	);
	await dialog.getByText('Transport Orders', { exact: true }).tap();
	await expect(drawer).toHaveAttribute('data-open', 'true');

	// Accepting the existing orders must preserve the applied plan.
	saved = page.waitForResponse((response) => response.url().endsWith(saveEndpoint));
	await dialog.getByRole('button', { name: 'OK', exact: true }).tap();
	const { fleet } = fromJson(
		UpdateFleetOrdersResponseSchema,
		(await (await saved).json()) as UpdateFleetOrdersResponseJson
	);
	expect(fleet!.fleetOrders!.waypoints[0].transportTasks).toEqual(plan.tasks);
	await expect(dialog).toBeHidden();
	await expect(drawer).toHaveAttribute('data-open', 'true');

	await task.getByRole('button', { name: 'Edit transport orders' }).tap();
	await dialog
		.locator('[name="actionColonists"]')
		.selectOption(String(WaypointTaskTransportAction.UNLOAD_AMOUNT));
	await dialog.locator('[name="amountColonists"]').fill('25');
	await expect(drawer).toHaveAttribute('data-open', 'true');
	await dialog.getByRole('button', { name: 'Cancel', exact: true }).tap();
	await expect(drawer).toHaveAttribute('data-open', 'true');
	await task.getByRole('button', { name: 'Edit transport orders' }).tap();
	await expect(dialog.locator('[name="actionColonists"]')).toHaveValue(
		String(WaypointTaskTransportAction.LOAD_ALL)
	);
	await dialog.getByRole('button', { name: 'Cancel', exact: true }).tap();
	await drawer.getByRole('button', { name: 'Close command pane backdrop', exact: true }).tap({
		position: { x: 10, y: page.viewportSize()!.height * 0.15 }
	});
	await expect(drawer).toHaveAttribute('data-open', 'false');
});

test('rapid waypoint taps update before a delayed save and persist in tap order', async ({
	testGamePage
}) => {
	const { page, universe } = await testGamePage('Kitchen Sink');
	const homeworld = universe.planets[0];
	const target = universe.planets[1];
	const secondTarget = universe.planets[2];
	const endpoint = '**/api/grpc/craig_stars.v1.FleetService/UpdateFleetOrders';
	let releaseSave!: () => void;
	const heldSave = new Promise<void>((resolve) => {
		releaseSave = resolve;
	});
	const requests: number[][] = [];
	await page.route(endpoint, async (route) => {
		const { fleetOrders } = fromJson(
			UpdateFleetOrdersRequestSchema,
			route.request().postDataJSON() as UpdateFleetOrdersRequestJson
		);
		requests.push(fleetOrders!.waypoints.map((wp) => wp.mapObjectTarget?.targetNum ?? 0));
		if (requests.length === 1) await heldSave;
		await route.continue();
	});

	try {
		await tapMapObject(page, homeworld);
		await page.locator('#add-waypoint').tap();
		await tapMapObject(page, target);
		await expect.poll(() => requests.length).toBe(1);
		await tapMapObject(page, secondTarget);
		const drawer = page.locator('[data-type="command-drawer"]');
		await drawer.getByRole('button', { name: 'show command pane button' }).tap();
		const waypoints = drawer.locator('[data-type="command-tile"][data-id="Fleet Waypoints"]');
		await expect(waypoints.locator('li')).toHaveText([
			homeworld!.mapObject!.name,
			target.mapObject!.name,
			secondTarget.mapObject!.name
		]);
		// Both taps are already visible, but only the first save has reached the network.
		expect(requests).toHaveLength(1);
		const saved = page.waitForResponse(
			(res) =>
				res.url().endsWith('/craig_stars.v1.FleetService/UpdateFleetOrders') &&
				requests.length === 2
		);
		releaseSave();
		const { fleet } = fromJson(
			UpdateFleetOrdersResponseSchema,
			(await (await saved).json()) as UpdateFleetOrdersResponseJson
		);
		expect(fleet!.fleetOrders!.waypoints.map((wp) => wp.mapObjectTarget?.targetNum)).toEqual([
			homeworld!.mapObject!.num,
			target.mapObject!.num,
			secondTarget.mapObject!.num
		]);
		await expect(waypoints.locator('li')).toHaveCount(3);
	} finally {
		releaseSave();
		await page.unrouteAll({ behavior: 'wait' });
	}
});

/**
 * Drag a finger across an element, from one fraction of its width to another. Playwright only
 * taps, so send the touches through the chrome devtools protocol. Chrome turns these into pointer
 * events too, like a real phone.
 */
async function touchDrag(page: Page, locator: Locator, from: number, to: number) {
	const box = await locator.boundingBox();
	if (!box) {
		throw new Error('element not rendered');
	}
	const y = box.y + box.height / 2;
	const point = (fraction: number) => [{ x: box.x + box.width * fraction, y }];

	const cdp = await page.context().newCDPSession(page);
	await cdp.send('Input.dispatchTouchEvent', { type: 'touchStart', touchPoints: point(from) });
	const steps = 5;
	for (let i = 1; i <= steps; i++) {
		await cdp.send('Input.dispatchTouchEvent', {
			type: 'touchMove',
			touchPoints: point(from + ((to - from) * i) / steps)
		});
	}
	await cdp.send('Input.dispatchTouchEvent', { type: 'touchEnd', touchPoints: [] });
	await cdp.detach();
}

test('dragging the warp gauge to use a stargate saves the waypoint', async ({ testGamePage }) => {
	const { page, universe } = await testGamePage('Stargate Test');

	const planet1 = universe.planets.find((p) => p.mapObject?.name === 'Planet 1');

	// record each warp speed the server saves for the scout's waypoint to Planet 2
	const savedWarpSpeeds: number[] = [];
	page.on('response', async (response) => {
		if (response.url().includes('/api/grpc/craig_stars.v1.FleetService/UpdateFleetOrders')) {
			const { fleet } = fromJson(
				UpdateFleetOrdersResponseSchema,
				(await response.json()) as UpdateFleetOrdersResponseJson
			);
			savedWarpSpeeds.push(fleet?.fleetOrders?.waypoints[1]?.warpSpeed ?? 0);
		}
	});

	// tap the homeworld to cycle to the scout, then open the drawer to its commands
	const drawer = page.locator('[data-type="command-drawer"]');
	await tapMapObject(page, planet1);
	await drawer.getByRole('button', { name: 'show command pane button' }).tap();

	const waypoints = drawer.locator('[data-type="command-tile"][data-id="Fleet Waypoints"]');
	await waypoints.getByRole('button', { name: 'Planet 2' }).tap();
	const gauge = waypoints.getByRole('slider', { name: 'Warp speed' });
	await expect(gauge).toHaveAttribute('aria-valuenow', '5');

	// drag from warp 5 all the way to the right
	await touchDrag(page, gauge, 5 / StargateWarpSpeed, 1);
	await expect(gauge).toHaveText('Use Stargate');

	// the last save is the speed we dragged to
	await expect.poll(() => savedWarpSpeeds.at(-1)).toBe(StargateWarpSpeed);
});
