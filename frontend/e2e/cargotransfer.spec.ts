import { MapObjectType } from '../src/lib/protogen/craig_stars/v1/common_pb';
import {
	expect,
	submitTurn,
	test
} from './setup';

test('Cargo Transfer Planet Owned', async ({ testGamePage }) => {
	const { page } = await testGamePage('Cargo Transfer Planet Owned');

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

	await submitTurn(page);

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
	await expect(statusTile.getByText('Population 287,300').first()).toBeVisible();

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

test('Cargo Transfer Fleets', async ({ testGamePage }) => {
	const { page } = await testGamePage('Cargo Transfer Fleets');

	const otherFleetsHereTile = await page
		.locator('[data-type="command-tile"][data-id="Other Entities Here"]')
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
	await submitTurn(page);

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

test('Cargo Transfer Split', async ({ testGamePage }) => {
	const { page } = await testGamePage('Cargo Transfer Split');

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
	await submitTurn(page);

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

test('Cargo Transfer Jettison', async ({ testGamePage }) => {
	const { page } = await testGamePage('Cargo Transfer Jettison');

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
	await submitTurn(page);

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

test('Cargo Transfer Salvage', async ({ testGamePage }) => {
	const { page } = await testGamePage('Cargo Transfer Salvage');

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
	await submitTurn(page);

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

test('Cargo Transfer MineralPacket', async ({ testGamePage }) => {
	const name = 'Cargo Transfer MineralPacket';
	const { page } = await testGamePage(name);

	const otherFleetsHereTile = await page
		.locator('[data-type="command-tile"][data-id="Other Entities Here"]')
		.first();

	await otherFleetsHereTile.getByRole('button', { name: 'Transfer' }).first().click();

	await page
		.locator('[data-id="ironium"][data-type="transfer-to-source-button"]')
		.click({ clickCount: 1 });
	await page
		.locator('[data-id="boranium"][data-type="transfer-to-source-button"]')
		.click({ clickCount: 2 });
	await page
		.locator('[data-id="germanium"][data-type="transfer-to-source-button"]')
		.click({ clickCount: 3 });
	await page.getByRole('button', { name: 'Ok' }).click();

	// select the mineralPacket in the map object summary
	let mapObjectSummary = await page.locator('[data-type="map-object-summary"]').first();
	await mapObjectSummary.locator('[data-type="cycle-selected-map-object-button"]').first().click();

	// MineralPacket #1 should be selected
	await expect(
		mapObjectSummary
			.locator('div')
			.filter({ hasText: /^Humanoids Mineral Packet #1$/ })
			.first()
	).toBeVisible();

	// mineralPacket cargo should be updated
	await mapObjectSummary
		.getByText('Location: (0, 0) Traveling at')
		.first()
		.scrollIntoViewIfNeeded();
	await expect(mapObjectSummary.getByText('Ironium 49kT').first()).toBeVisible();
	await expect(mapObjectSummary.getByText('Boranium 48kT').first()).toBeVisible();
	await expect(mapObjectSummary.getByText('Germanium 47kT').first()).toBeVisible();

	// submit turn
	await page.getByRole('button', { name: 'Submit Turn' }).click();

	// wait for turn submit to finish
	await page.locator('#loading-modal').waitFor({ state: 'visible' }); // wait for loading modal to show up
	await expect(page.locator('#loading-modal')).not.toHaveClass(/modal-open/); // ensure submit is done

	// fleet cargo should be updated
	const fuelAndCargoTile = await page
		.locator('[data-type="command-tile"][data-id="Fuel & Cargo"]')
		.first();
	await expect(fuelAndCargoTile.getByText('Ironium 11kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Boranium 12kT').first()).toBeVisible();
	await expect(fuelAndCargoTile.getByText('Germanium 13kT').first()).toBeVisible();

	// click the mineral packet
	await page.locator(`[data-id="${MapObjectType.MINERAL_PACKET}-1-1"]`).click({ force: true });

	// select the mineralPacket in the map object summary
	mapObjectSummary = await page.locator('[data-type="map-object-summary"]').first();

	// select mineralPacket #1 and check it
	await expect(
		mapObjectSummary
			.locator('div')
			.filter({ hasText: /^Humanoids Mineral Packet #1$/ })
			.first()
	).toBeVisible();
	await mapObjectSummary
		.getByText('Location: (25, 0) Traveling at')
		.first()
		.scrollIntoViewIfNeeded();
	// mineralPacket cargo should be decayed but also updated
	await expect(mapObjectSummary.getByText('Ironium 25kT').first()).toBeVisible();
	await expect(mapObjectSummary.getByText('Boranium 24kT').first()).toBeVisible();
	await expect(mapObjectSummary.getByText('Germanium 24kT').first()).toBeVisible();
});
