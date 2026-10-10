import { fromJson } from '@bufbuild/protobuf';
import { expect, type Locator, type Page } from '@playwright/test';
import {
	UpdateFleetOrdersRequestSchema,
	UpdateFleetOrdersResponseSchema,
	type UpdateFleetOrdersRequestJson,
	type UpdateFleetOrdersResponseJson
} from '../src/lib/protogen/craig_stars/v1/fleetservice_pb';
import { StargateWarpSpeed } from '../src/lib/types/Consts';
import { tapMapObject } from './helpers/scanner';
import { test } from './setup';

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
