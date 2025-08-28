import { WaypointTask } from '../src/lib/protogen/craig_stars/v1/fleet_pb';
import { key, nearest } from '../src/lib/types/MapObject';
import { apiErrorsFailTest, expect, loadGamePage, submitTurn, test } from './setup';

test('new game scout test', async ({ newGamePage }) => {
	const { page, id, universe } = newGamePage;
	apiErrorsFailTest(page, id);

	const homeworld = universe.planets?.find((p) => p.mapObject?.playerNum === 1 && p.homeworld);
	if (!homeworld) {
		throw new Error('failed to find homeworld');
	}

	const nearestPlanet = nearest(homeworld, universe.planets);
	if (!nearestPlanet) {
		throw new Error('failed to find nearest planet');
	}

	// click the homeworld once to cycle to the first ship, a scout
	await page.locator(`[data-id="${key(homeworld)}"]`).click({ force: true });

	await page.locator('[data-type="command-tile"][data-id="Long Range Scout #1"]').first();

	// shift/meta click the planet to set a waypoint
	await page
		.locator(`[data-id="${key(nearestPlanet)}"]`)
		.click({ force: true, modifiers: ['Shift', 'Meta'] });

	// waypoints tile should update
	const fleetWaypointsTile = await page
		.locator('[data-type="command-tile"][data-id="Fleet Waypoints"]')
		.first();
	await expect(
		fleetWaypointsTile.getByText(`Coming from ${homeworld.mapObject?.name}`).first()
	).toBeVisible();

	const { universe: updatedUniverse } = await submitTurn(page);

	// our fleet should have flown away
	expect(
		updatedUniverse?.fleets.find((f) => f.mapObject?.name === 'Long Range Scout #1')
			?.orbitingPlanetNum
	).not.toBe(homeworld.mapObject?.num);
});

test('Scout Test', async ({ authenticatedPage }) => {
	const { page, gameId, universe } = await loadGamePage(authenticatedPage, 'Scout Test');
	apiErrorsFailTest(page, gameId);

	const homeworld = universe.planets?.find((p) => p.mapObject?.playerNum === 1 && p.homeworld);
	if (!homeworld) {
		throw new Error('failed to find homeworld');
	}

	const planet2 = universe.planets[1];

	// click the homeworld once to cycle to the first ship, a scout
	await page.locator(`[data-id="${key(homeworld)}"]`).click({ force: true });

	await page.locator('[data-type="command-tile"][data-id="Long Range Scout #1"]').first();

	// meta click the planet to set a waypoint at max speed
	await page.locator(`[data-id="${key(planet2)}"]`).click({ force: true, modifiers: ['Meta'] });

	// waypoints tile should update
	const fleetWaypointsTile = await page
		.locator('[data-type="command-tile"][data-id="Fleet Waypoints"]')
		.first();
	await expect(
		fleetWaypointsTile.getByText(`Coming from ${homeworld.mapObject?.name}`).first()
	).toBeVisible();

	const { universe: updatedUniverse } = await submitTurn(page);

	// our fleet should have flown to the new planet
	expect(
		updatedUniverse?.fleets.find((f) => f.mapObject?.name === 'Long Range Scout #1')
			?.orbitingPlanetNum
	).toBe(planet2.mapObject?.num);

	// click the nearest planet twice to cycle to the scout
	await page.locator(`[data-id="${key(planet2)}"]`).click({ force: true });
	await page.locator(`[data-id="${key(planet2)}"]`).click({ force: true });

	await expect(
		page.locator('[data-type="command-tile"][data-id="Orbiting Planet 2"]').first()
	).toBeVisible();
});

test('Colonizer Test', async ({ authenticatedPage }) => {
	const { page, gameId, universe, player } = await loadGamePage(
		authenticatedPage,
		'Colonizer Test'
	);
	apiErrorsFailTest(page, gameId);

	const homeworld = universe.planets?.find((p) => p.mapObject?.playerNum === 1 && p.homeworld);
	if (!homeworld) {
		throw new Error('failed to find homeworld');
	}

	const planet2 = universe.planets[1];

	// click the homeworld once to cycle to the colonizer
	await page.locator(`[data-id="${key(homeworld)}"]`).click({ force: true });

	await page.locator('[data-type="command-tile"][data-id="Santa Maria #1"]').first();

	// open up cargo dialog
	await page.getByText('of 25kT').first().click();
	// transfer max colonists
	await page
		.locator(`[data-id="colonists"][data-type="transfer-to-source-button"]`)
		.click({ clickCount: 1, modifiers: ['Meta'] });

	const transferCargoResponse = page.waitForResponse(
		(resp) =>
			resp.url().includes('/api/grpc/craig_stars.v1.FleetService/TransferCargo') &&
			resp.request().method() === 'POST' &&
			resp.status() === 200
	);

	await page.getByRole('button', { name: 'Ok' }).click();

	// wait for the TransferCargo to complete
	await transferCargoResponse;

	// wait for the UpdateFleetOrders to complete
	const updateFleetOrderResponse = page.waitForResponse(
		(resp) =>
			resp.url().includes('/api/grpc/craig_stars.v1.FleetService/UpdateFleetOrders') &&
			resp.request().method() === 'POST' &&
			resp.status() === 200
	);

	// meta click the planet to set a waypoint at max speed
	await page.locator(`[data-id="${key(planet2)}"]`).click({ force: true, modifiers: ['Meta'] });

	await updateFleetOrderResponse;
	
	// waypoints tile should update
	const fleetWaypointsTile = await page
		.locator('[data-type="command-tile"][data-id="Fleet Waypoints"]')
		.first();
	await expect(
		fleetWaypointsTile.getByText(`Coming from ${homeworld.mapObject?.name}`).first()
	).toBeVisible();

	// waypoint task should showcolonize
	const fleetWaypointTaskTile = await page
		.locator('[data-type="command-tile"][data-id="Waypoint Task"]')
		.first();
	const selectTask = await fleetWaypointTaskTile
		.locator('[data-type="select-waypoint-task"][data-id="waypoint-task"]')
		.first();
	await expect(selectTask).toHaveValue(String(WaypointTask.COLONIZE));

	const { universe: updatedUniverse } = await submitTurn(page);

	const updatedPlanet2 = updatedUniverse?.planets[1];
	// we should colonize the new planet
	expect(updatedPlanet2?.mapObject?.playerNum).toBe(player.num);

	// click the second planet twice to command it
	await page.locator(`[data-id="${key(updatedPlanet2)}"]`).click({ force: true });
	await page.locator(`[data-id="${key(updatedPlanet2)}"]`).click({ force: true });

	await expect(
		page.locator('[data-type="command-tile"][data-id="Planet 2"]').first()
	).toBeVisible();

	// Planet 2 should be selected
	const mapObjectSummary = await page.locator('[data-type="map-object-summary"]').first();

	await expect(
		mapObjectSummary
			.locator('div')
			.filter({ hasText: /^Planet 2$/ })
			.first()
	).toBeVisible();

	await expect(mapObjectSummary).toContainText('Population: 2,500');
	await expect(mapObjectSummary).toContainText('Value: 100%');
});
