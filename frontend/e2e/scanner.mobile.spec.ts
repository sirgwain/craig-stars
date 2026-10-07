import { expect, type Page } from '@playwright/test';
import { tapMapObject } from './helpers/scanner';
import { test } from './setup';

// On phones the command pane is a drawer below the scanner instead of a side panel
function drawer(page: Page) {
	return page.locator('[data-type="command-drawer"]');
}

function drawerToggle(page: Page) {
	return drawer(page).getByRole('button', { name: 'show command pane button' });
}

// tap the drawer header to open or close it, and wait for it to finish sliding
async function toggleDrawer(page: Page, open: boolean) {
	await drawerToggle(page).tap();
	await expect(drawer(page)).toHaveAttribute('data-open', `${open}`);
	await expect(drawerToggle(page)).toHaveCount(1);
}

function scannerCanvas(page: Page) {
	return page.locator('[data-type="scanner-canvas"]');
}

test('phone layout shows the scanner above the command drawer', async ({ testGamePage }) => {
	const { page } = await testGamePage('Kitchen Sink');

	const scanner = await scannerCanvas(page).boundingBox();
	const commandDrawer = await drawer(page).boundingBox();
	const viewport = page.viewportSize();
	if (!scanner || !commandDrawer || !viewport) {
		throw new Error('scanner or drawer not rendered');
	}

	// the scanner fills the width of the screen and takes up most of it
	expect(scanner.width).toBeCloseTo(viewport.width, -1);
	expect(scanner.height).toBeGreaterThan(viewport.height / 3);

	// the drawer sits below the scanner, at the bottom of the screen
	expect(commandDrawer.y).toBeGreaterThanOrEqual(scanner.y + scanner.height - 1);
	expect(commandDrawer.y + commandDrawer.height).toBeCloseTo(viewport.height, -1);
});

test('tapping map objects shows them in the drawer', async ({ testGamePage }) => {
	const { page, universe } = await testGamePage('Kitchen Sink');

	const otherPlayerPlanet = universe.planets.find((p) => p.mapObject?.playerNum === 2);
	const otherPlayerFleet = universe.fleets.find((f) => f.mapObject?.playerNum === 2);
	const mineralPacket = universe.mineralPackets[0];

	for (const mo of [otherPlayerPlanet, otherPlayerFleet, mineralPacket]) {
		await tapMapObject(page, mo);
		await expect(drawer(page)).toContainText(mo?.mapObject?.name ?? 'missing map object');
	}
});

test('command drawer opens and closes', async ({ testGamePage }) => {
	const { page } = await testGamePage('Kitchen Sink');

	// the homeworld is commanded when the game loads
	await toggleDrawer(page, true);
	await expect(drawer(page).locator('#planet-production-tile')).toBeVisible();

	await toggleDrawer(page, false);
	await expect(drawer(page).locator('#planet-production-tile')).toBeHidden();
});

/**
 * The drawer changes height as different things are selected, which resizes the scanner. The
 * canvas must be resized and redrawn before the browser paints, otherwise the old image is
 * stretched to the new size for a frame and the scanner visibly bounces.
 */
test('scanner canvas resizes without stretching', async ({ testGamePage }) => {
	const { page, universe } = await testGamePage('Kitchen Sink');

	// Watch the scanner's container. Resize observers run in the order they were created, so this
	// runs after the scanner's own observer, right before paint. At that point the canvas bitmap
	// and css size must already match the container.
	await page.evaluate(() => {
		const canvas = document.querySelector<HTMLCanvasElement>('[data-type="scanner-canvas"]');
		const container = canvas?.parentElement;
		if (!canvas || !container) {
			throw new Error('scanner canvas not found');
		}
		const results = { resizes: 0, stretched: [] as string[] };
		(window as unknown as { scannerResizes: typeof results }).scannerResizes = results;
		new ResizeObserver(() => {
			results.resizes++;
			// container sizes can be fractional, so compare exact sizes, allowing for rounding
			const dpr = window.devicePixelRatio;
			const containerHeight = container.getBoundingClientRect().height;
			const cssHeight = canvas.getBoundingClientRect().height;
			if (
				Math.abs(canvas.height - containerHeight * dpr) > 1 ||
				Math.abs(cssHeight - containerHeight) > 0.5
			) {
				results.stretched.push(
					`container ${containerHeight}px, canvas css ${cssHeight}px, bitmap ${canvas.height}px`
				);
			}
		}).observe(container);
	});

	// select things with different summaries, then open and close the drawer
	const otherPlayerPlanet = universe.planets.find((p) => p.mapObject?.playerNum === 2);
	const otherPlayerFleet = universe.fleets.find((f) => f.mapObject?.playerNum === 2);
	for (const mo of [
		otherPlayerPlanet,
		otherPlayerFleet,
		universe.mineralPackets[0],
		universe.planets[0]
	]) {
		await tapMapObject(page, mo);
		await expect(drawer(page)).toContainText(mo?.mapObject?.name ?? 'missing map object');
	}
	await toggleDrawer(page, true);
	await toggleDrawer(page, false);

	const results = await page.evaluate(
		() =>
			(window as unknown as { scannerResizes: { resizes: number; stretched: string[] } })
				.scannerResizes
	);
	// the first callback is the initial observation, make sure we actually resized the scanner
	expect(results.resizes).toBeGreaterThan(1);
	expect(results.stretched).toEqual([]);
});
