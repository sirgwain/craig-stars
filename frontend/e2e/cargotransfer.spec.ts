import { apiErrorsFailTest, expect, loadGamePage, test } from './setup';

test('Cargo Transfer Planet Owned', async ({ authenticatedPage }) => {
	const { page, gameId } = await loadGamePage(authenticatedPage, 'Cargo Transfer Planet Owned');
	apiErrorsFailTest(page, gameId);

	let fleetsInOrbitTile = await page
		.locator('[data-type="command-tile"][data-id="Fleets In Orbit"]')
		.first();

	await fleetsInOrbitTile.getByRole('button', { name: 'Goto' }).first().click();

	// open transfer dialog
	const fleetOrbitingTile = await page
		.locator('[data-type="command-tile"][data-id="Orbiting Planet 1"]')
		.first();
	await fleetOrbitingTile
		.getByRole('button', { name: /Jettison|Transfer/ })
		.first()
		.click();

	// transfer one of each to the fleet
	for (const cargoType of ['ironium', 'boranium', 'germanium', 'colonists']) {
		await page
			.locator(`[data-id="${cargoType}"][data-type="transfer-to-source-button"]`)
			.click({ clickCount: 1 });
	}
	await page.getByRole('button', { name: 'Ok' }).click();

	// verify cargo
	let fuelAndCargoTile = await page
		.locator('[data-type="command-tile"][data-id="Fuel & Cargo"]')
		.first();
	await expect(fuelAndCargoTile.getByText('Ironium 1kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Boranium 1kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Germanium 1kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Colonists 1kT').first()).toBeVisible();

	// submit turn and make sure it "sticks"
	await page.getByRole('button', { name: 'Submit Turn' }).click();

	// verify cargo
	const mineralsOnHandTile = await page
		.locator('[data-type="command-tile"][data-id="Minerals on Hand"]')
		.first();
	await expect(mineralsOnHandTile.getByText('Ironium 999kT').first()).toBeVisible();
	await expect(mineralsOnHandTile.getByText('Boranium 999kT').first()).toBeVisible();
	await expect(mineralsOnHandTile.getByText('Germanium 999kT').first()).toBeVisible();

	// pop grows
	// TODO: fragile test?
	const statusTile = await page.locator('[data-type="command-tile"][data-id="Status"]').first();
	await expect(statusTile.getByText('Population 287,400').first()).toBeVisible();

	// goto fleet
	fleetsInOrbitTile = await page
		.locator('[data-type="command-tile"][data-id="Fleets In Orbit"]')
		.first();
	await fleetsInOrbitTile.getByRole('button', { name: 'Goto' }).first().click();

	// verify cargo after turn submit
	fuelAndCargoTile = await page
		.locator('[data-type="command-tile"][data-id="Fuel & Cargo"]')
		.first();
	await expect(fuelAndCargoTile.getByText('Ironium 1kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Boranium 1kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Germanium 1kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Colonists 1kT').first()).toBeVisible();
});

test('Cargo Transfer Fleets', async ({ authenticatedPage }) => {
	const { page, gameId } = await loadGamePage(authenticatedPage, 'Cargo Transfer Fleets');
	apiErrorsFailTest(page, gameId);

	const otherFleetsHereTile = await page
		.locator('[data-type="command-tile"][data-id="Other Fleets Here"]')
		.first();

	await otherFleetsHereTile.getByRole('button', { name: 'Transfer' }).first().click();

	// transfer one of each
	// teamster starts with {10, 10, 10, 10} 100 fuel
	// colony ship starts with {5, 5, 5, 5} 10 fuel
	for (const cargoType of ['fuel', 'ironium', 'boranium', 'germanium', 'colonists']) {
		await page
			.locator(`[data-id="${cargoType}"][data-type="transfer-to-dest-button"]`)
			.click({ clickCount: 1 });
	}
	await page.getByRole('button', { name: 'Ok' }).click();

	// verify cargo
	let fuelAndCargoTile = await page
		.locator('[data-type="command-tile"][data-id="Fuel & Cargo"]')
		.first();
	await expect(fuelAndCargoTile.getByText('Ironium 9kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Boranium 9kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Germanium 9kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Colonists 9kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Fuel 99 of 450mg').first().first()).toBeVisible();

	// select colony ship
	await page.getByRole('button', { name: 'Next', exact: true }).first().click();
	await expect(fuelAndCargoTile.getByText('Ironium 6kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Boranium 6kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Germanium 6kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Colonists 6kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Fuel 11 of 200mg').first().first()).toBeVisible();

	// submit turn and make sure it "sticks"
	await page.getByRole('button', { name: 'Submit Turn' }).click();

	// verify cargo
	fuelAndCargoTile = await page
		.locator('[data-type="command-tile"][data-id="Fuel & Cargo"]')
		.first();
	await expect(fuelAndCargoTile.getByText('Ironium 9kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Boranium 9kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Germanium 9kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Colonists 9kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Fuel 99 of 450mg').first().first()).toBeVisible();

	// select colony ship
	await page.getByRole('button', { name: 'Next', exact: true }).first().click();
	await expect(fuelAndCargoTile.getByText('Ironium 6kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Boranium 6kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Germanium 6kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Colonists 6kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Fuel 11 of 200mg').first().first()).toBeVisible();
});

test('Cargo Transfer Split', async ({ authenticatedPage }) => {
	const { page, gameId } = await loadGamePage(authenticatedPage, 'Cargo Transfer Split');
	apiErrorsFailTest(page, gameId);

	const fleetOrbitingTile = await page
		.locator('[data-type="command-tile"][data-id="In Deep Space"]')
		.first();

	// jettison 10 ironium
	await fleetOrbitingTile
		.getByRole('button', { name: /Jettison|Transfer/ })
		.first()
		.click();
	await page
		.locator('[data-id="ironium"][data-type="transfer-to-dest-button"]')
		.click({ clickCount: 10 });
	await page.getByRole('button', { name: 'Ok' }).click();

	// split fleet
	await page.getByRole('button', { name: 'Split All' }).first().click();

	// verify cargo
	let fuelAndCargoTile = await page
		.locator('[data-type="command-tile"][data-id="Fuel & Cargo"]')
		.first();
	await expect(fuelAndCargoTile.getByText('Ironium 36kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Boranium 45kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Germanium 45kT').first()).toBeVisible();

	// select colony ship
	await page.getByRole('button', { name: 'Next', exact: true }).first().click();
	await expect(fuelAndCargoTile.getByText('Ironium 4kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Boranium 5kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Germanium 5kT').first()).toBeVisible();

	// grab 10 ironium back
	await fleetOrbitingTile
		.getByRole('button', { name: /Jettison|Transfer/ })
		.first()
		.click();
	await page
		.locator('[data-id="ironium"][data-type="transfer-to-source-button"]')
		.click({ clickCount: 10 });
	await page.getByRole('button', { name: 'Ok' }).click();

	// submit turn
	await page.getByRole('button', { name: 'Submit Turn' }).click();

	fuelAndCargoTile = await page
		.locator('[data-type="command-tile"][data-id="Fuel & Cargo"]')
		.first();
	await expect(fuelAndCargoTile.getByText('Ironium 36kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Boranium 45kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Germanium 45kT').first()).toBeVisible();

	// select colony ship
	await page.getByRole('button', { name: 'Next', exact: true }).first().click();
	await expect(fuelAndCargoTile.getByText('Ironium 14kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Boranium 5kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Germanium 5kT').first()).toBeVisible();
});

test('Cargo Transfer Jettison', async ({ authenticatedPage }) => {
	const { page, gameId } = await loadGamePage(authenticatedPage, 'Cargo Transfer Jettison');
	apiErrorsFailTest(page, gameId);

	const fleetOrbitingTile = await page
		.locator('[data-type="command-tile"][data-id="In Deep Space"]')
		.first();

	// click the cargo button
	await fleetOrbitingTile
		.getByRole('button', { name: /Jettison|Transfer/ })
		.first()
		.click();
	await page
		.locator('[data-id="ironium"][data-type="transfer-to-dest-button"]')
		.click({ clickCount: 10 });
	await page
		.locator('[data-id="boranium"][data-type="transfer-to-dest-button"]')
		.click({ clickCount: 20 });
	await page
		.locator('[data-id="germanium"][data-type="transfer-to-dest-button"]')
		.click({ clickCount: 30 });
	await page.getByRole('button', { name: 'Ok' }).click();

	// fleet cargo should be updated
	let fuelAndCargoTile = await page
		.locator('[data-type="command-tile"][data-id="Fuel & Cargo"]')
		.first();
	await expect(fuelAndCargoTile.getByText('Ironium 40kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Boranium 30kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Germanium 20kT').first()).toBeVisible();

	// submit turn
	await page.getByRole('button', { name: 'Submit Turn' }).click();

	// fleet cargo should be updated
	fuelAndCargoTile = await page
		.locator('[data-type="command-tile"][data-id="Fuel & Cargo"]')
		.first();
	await expect(fuelAndCargoTile.getByText('Ironium 40kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Boranium 30kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Germanium 20kT').first()).toBeVisible();

	// select the salvage in the map object summary
	const mapObjectSummary = await page.locator('[data-type="map-object-summary"]').first();

	await mapObjectSummary.locator('[data-type="cycle-selected-map-object-button"]').click();

	// select salvage #1 and check it
	await expect(
		mapObjectSummary
			.locator('div')
			.filter({ hasText: /^Salvage #1$/ })
			.first()
	).toBeVisible();
	await mapObjectSummary
		.getByText('Salvage #1 Humanoids Location')
		.first()
		.scrollIntoViewIfNeeded();
	// salvage cargo should exist but also be decayed
	await expect(mapObjectSummary.getByText('Ironium 0kT').first()).toBeVisible();
	await expect(mapObjectSummary.getByText('Boranium 10kT').first()).toBeVisible();
	await expect(mapObjectSummary.getByText('Germanium 20kT').first()).toBeVisible();
});

test('Cargo Transfer Salvage', async ({ authenticatedPage }) => {
	const { page, gameId } = await loadGamePage(authenticatedPage, 'Cargo Transfer Salvage');
	apiErrorsFailTest(page, gameId);

	const fleetOrbitingTile = await page
		.locator('[data-type="command-tile"][data-id="In Deep Space"]')
		.first();

	// click the cargo button
	await fleetOrbitingTile
		.getByRole('button', { name: /Jettison|Transfer/ })
		.first()
		.click();
	await page
		.locator('[data-id="ironium"][data-type="transfer-to-dest-button"]')
		.click({ clickCount: 1 });
	await page
		.locator('[data-id="boranium"][data-type="transfer-to-dest-button"]')
		.click({ clickCount: 2 });
	await page
		.locator('[data-id="germanium"][data-type="transfer-to-dest-button"]')
		.click({ clickCount: 3 });
	await page.getByRole('button', { name: 'Ok' }).click();

	// select the salvage in the map object summary
	let mapObjectSummary = await page.locator('[data-type="map-object-summary"]').first();
	await mapObjectSummary.locator('[data-type="cycle-selected-map-object-button"]').first().click();

	// Salvage #1 should be selected
	await expect(
		mapObjectSummary
			.locator('div')
			.filter({ hasText: /^Salvage #1$/ })
			.first()
	).toBeVisible();
	// salvage cargo should be updated
	await mapObjectSummary
		.getByText('Salvage #1 Humanoids Location')
		.first()
		.scrollIntoViewIfNeeded();
	await expect(mapObjectSummary.getByText('Ironium 51kT').first()).toBeVisible();
	await expect(mapObjectSummary.getByText('Boranium 52kT').first()).toBeVisible();
	await expect(mapObjectSummary.getByText('Germanium 53kT').first()).toBeVisible();

	// submit turn
	await page.getByRole('button', { name: 'Submit Turn' }).click();

	// fleet cargo should be updated
	const fuelAndCargoTile = await page
		.locator('[data-type="command-tile"][data-id="Fuel & Cargo"]')
		.first();
	await expect(fuelAndCargoTile.getByText('Ironium 9kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Boranium 8kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Germanium 7kT').first()).toBeVisible();

	// select the salvage in the map object summary
	mapObjectSummary = await page.locator('[data-type="map-object-summary"]').first();

	await mapObjectSummary.locator('[data-type="cycle-selected-map-object-button"]').click();

	// select salvage #1 and check it
	await expect(
		mapObjectSummary
			.locator('div')
			.filter({ hasText: /^Salvage #1$/ })
			.first()
	).toBeVisible();
	await mapObjectSummary
		.getByText('Salvage #1 Humanoids Location')
		.first()
		.scrollIntoViewIfNeeded();
	// salvage cargo should be decayed but also updated
	await expect(mapObjectSummary.getByText('Ironium 41kT').first()).toBeVisible();
	await expect(mapObjectSummary.getByText('Boranium 42kT').first()).toBeVisible();
	await expect(mapObjectSummary.getByText('Germanium 43kT').first()).toBeVisible();
});
