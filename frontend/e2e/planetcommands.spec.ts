import { apiErrorsFailTest, expect, submitTurn, test } from './setup';

test('production tile clear', async ({ newGamePage }) => {
	const { page } = newGamePage;
	apiErrorsFailTest(page);

	// confirm delete
	page.once('dialog', async (dialog) => {
		await dialog.accept();
	});

	const productionTile = await page
		.locator('[data-type="command-tile"][data-id="Production"]')
		.first();

	await productionTile.getByRole('button', { name: 'Clear' }).click();

	await expect(productionTile.getByRole('list').filter({ hasText: /^$/ })).toBeVisible();
});

test('production queue change', async ({ newGamePage }) => {
	const { page, universe } = newGamePage;
	apiErrorsFailTest(page);

	const homeworld = universe.planets.find((p) => p.homeworld);
	if (!homeworld) {
		throw new Error('unable to find homeworld');
	}

	// open the production dialog
	const productionTile = await page
		.locator('[data-type="command-tile"][data-id="Production"]')
		.first();
	await productionTile.getByRole('button', { name: 'Change' }).click();

	const productionQueueDialog = await page
		.locator('[data-type="dialog"][data-id="Production Queue"]')
		.first();

	await expect(productionQueueDialog).toBeVisible();
	await expect(productionQueueDialog).toContainClass('modal-open');

	// clear out the queue
	await productionQueueDialog.getByRole('button', { name: 'Clear' }).click();

	await expect(productionQueueDialog).toContainText('-- Top of the Queue --');

	// add a factory
	await productionQueueDialog.getByRole('button', { name: 'Factory', exact: true }).click();
	await productionQueueDialog.getByRole('button', { name: 'Add' }).click();
	await expect(productionQueueDialog).toContainText('-- Top of the Queue -- Factory 1');

	// add 10 factories
	await productionQueueDialog.getByRole('button', { name: 'x10', exact: true }).click();
	await productionQueueDialog.getByRole('button', { name: 'Add' }).click();
	await expect(productionQueueDialog).toContainText('-- Top of the Queue -- Factory 11');

	// add 100 mines
	await productionQueueDialog.getByRole('button', { name: 'Mine', exact: true }).click();
	await productionQueueDialog.getByRole('button', { name: 'x100' }).click();
	await productionQueueDialog.getByRole('button', { name: 'Add' }).click();
	await expect(productionQueueDialog).toContainText('-- Top of the Queue -- Factory 11Mine 100');

	// move mine up
	await productionQueueDialog.getByRole('button', { name: 'Item Up' }).click();
	await expect(productionQueueDialog).toContainText('-- Top of the Queue -- Mine 100Factory 11');

	// move mine down
	await productionQueueDialog.getByRole('button', { name: 'Item Down' }).click();
	await expect(productionQueueDialog).toContainText('-- Top of the Queue -- Factory 11Mine 100');

	// remove mine
	await productionQueueDialog.getByRole('button', { name: 'Remove' }).click();
	await expect(productionQueueDialog).toContainText('-- Top of the Queue -- Factory 11');

	// add a ship to the top of the queue
	await productionQueueDialog.getByRole('button', { name: 'x1', exact: true }).click();
	await productionQueueDialog.getByRole('button', { name: '-- Top of the Queue --' }).click();
	await productionQueueDialog.getByRole('button', { name: 'Long Range Scout' }).click();
	await productionQueueDialog.getByRole('button', { name: 'Add' }).click();
	await expect(productionQueueDialog).toContainText(
		'-- Top of the Queue -- Long Range Scout 1Factory 11'
	);

	// save it
	await productionQueueDialog.getByRole('button', { name: 'Ok' }).click();

	// production queue tile should update
	await expect(productionTile).toContainText('Long Range Scout 1Factory 11');

	await submitTurn(page);

	// should have built a long range scout and a factory
	await expect(productionTile).toContainText('Factory 10');
});

test('production queue starbases add ships and packets', async ({ testGamePage }) => {
	// the homeworld has an orbital fort, with no dock or mass driver
	const { page } = await testGamePage('Production Starbases');
	apiErrorsFailTest(page);

	const productionTile = await page
		.locator('[data-type="command-tile"][data-id="Production"]')
		.first();
	await productionTile.getByRole('button', { name: 'Change' }).click();

	const productionQueueDialog = await page
		.locator('[data-type="dialog"][data-id="Production Queue"]')
		.first();
	await expect(productionQueueDialog).toBeVisible();
	await productionQueueDialog.getByRole('button', { name: 'Clear' }).click();

	const scout = productionQueueDialog.getByRole('button', { name: 'Long Range Scout' });
	const packet = productionQueueDialog.getByRole('button', {
		name: 'Mixed Mineral Packet',
		exact: true
	});

	// no ships or packets without a dock or mass driver
	await expect(scout).toHaveCount(0);
	await expect(packet).toHaveCount(0);

	// a space station has a dock, so we can build ships after it
	await productionQueueDialog.getByRole('button', { name: 'Station', exact: true }).click();
	await productionQueueDialog.getByRole('button', { name: 'Add' }).click();
	await expect(productionQueueDialog).toContainText('-- Top of the Queue -- Station 1');
	await expect(scout).toBeVisible();
	await expect(packet).toHaveCount(0);

	// a station with a mass driver after it, so we can build packets after it
	// it's priced as an upgrade from the station, not the orbital fort
	await productionQueueDialog.getByRole('button', { name: 'Flinger', exact: true }).click();
	await productionQueueDialog.getByRole('button', { name: 'Add' }).click();
	await expect(productionQueueDialog).toContainText('-- Top of the Queue -- Station 1Flinger 1');
	await expect(productionQueueDialog).toContainText(
		'Cost of Flinger x 1 Ironium 24kT Boranium 20kT Germanium 20kT Resources 70'
	);
	await expect(packet).toBeVisible();

	// add a packet after the flinger and save it
	await packet.click();
	await productionQueueDialog.getByRole('button', { name: 'Add' }).click();
	await expect(productionQueueDialog).toContainText(
		'-- Top of the Queue -- Station 1Flinger 1Mixed Mineral Packet 1'
	);
	await productionQueueDialog.getByRole('button', { name: 'Ok' }).click();
	await expect(productionTile).toContainText('Station 1Flinger 1Mixed');
});
