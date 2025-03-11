import { expect, test } from './setup';

const cargoTransferGameId = 2;

test('load cargo to colony ship', async ({ newGamePage }) => {
	const { page, id, name } = newGamePage;

	page.on('response', async (response) => {
		if (response.url().includes(`/api/games/${id}`) && !response.ok()) {
			// fail any api requests
			throw new Error(`API request failed: ${response.url()} - Status: ${response.status()}`);
		}
	});

	const gameLink = page.getByRole('link', { name: name });
	await expect(gameLink).toBeVisible();

	await page.locator('[data-type="fleets-in-orbit-select"]').first().selectOption('1');
	await page.getByRole('button', { name: 'Goto' }).click();
	await page.getByRole('button', { name: 'Transfer' }).first().click();
	// add {1, 2, 3, 10} cargo
	await page.locator('[data-id="ironium"][data-type="transfer-to-source-button"]').click();
	await page
		.locator('[data-id="boranium"][data-type="transfer-to-source-button"]')
		.click({ clickCount: 2 });
	await page
		.locator('[data-id="germanium"][data-type="transfer-to-source-button"]')
		.click({ clickCount: 3 });
	await page
		.locator('[data-id="colonists"][data-type="transfer-to-source-button"]')
		.click({ modifiers: ['Shift'] });
	await page.getByRole('button', { name: 'Ok' }).click();

	await expect(page.getByText('Ironium 1kT').first()).toBeVisible();
	await expect(page.getByText('Boranium 2kT').first()).toBeVisible();
	await expect(page.getByText('Germanium 3kT').first()).toBeVisible();
	await expect(page.getByText('Colonists 10kT').first()).toBeVisible();
});

test('load cargo between ships', async ({ newGamePage }) => {
	const { page, id, name } = newGamePage;

	page.on('response', async (response) => {
		if (response.url().includes(`/api/games/${id}`) && !response.ok()) {
			// fail any api requests
			throw new Error(`API request failed: ${response.url()} - Status: ${response.status()}`);
		}
	});

	const gameLink = page.getByRole('link', { name: name });
	await expect(gameLink).toBeVisible();

	await page.locator('[data-type="fleets-in-orbit-select"]').first().selectOption('1');
	await page.getByRole('button', { name: 'Goto' }).click();
	await page.getByRole('button', { name: 'Transfer' }).first().click();
	// add 10 ironium to the santa maria
	await page
		.locator('[data-id="ironium"][data-type="transfer-to-source-button"]')
		.click({ modifiers: ['Shift'] });
	await page.getByRole('button', { name: 'Ok' }).click();
	await expect(page.getByText('Ironium 10kT').first()).toBeVisible();

	// select teamster and transfer cargo to it
	await page.locator('[data-type="other-fleets-here-select"]').first().selectOption('1');
	await page.getByRole('button', { name: 'Transfer' }).nth(1).click();
	await page.getByRole('button', { name: 'x10', exact: true }).click();
	await page.locator('[data-id="ironium"][data-type="transfer-to-dest-button"]').first().click();
	await page.getByRole('button', { name: 'Ok' }).click();
	await expect(page.getByText('Ironium 0kT').first()).toBeVisible();

	// goto teamster and verify it transfered
	await page.getByRole('button', { name: 'Goto' }).nth(1).click();
	await expect(page.getByText('Ironium 10kT').first()).toBeVisible();
});

test('salvage cargo transfer', async ({ authenticatedPage: page }) => {
	page.on('response', async (response) => {
		if (response.url().includes(`/api/games/${cargoTransferGameId}`) && !response.ok()) {
			// fail any api requests
			throw new Error(`API request failed: ${response.url()} - Status: ${response.status()}`);
		}
	});

	// go to the cargo transfer test game
	await page.goto(`/games/${cargoTransferGameId}`, { waitUntil: 'networkidle' });
	await expect(page.getByRole('link', { name: 'Cargo Transfer Test -' })).toBeVisible();

	await page.getByRole('button', { name: 'search' }).first().click();
	await page.getByRole('searchbox', { name: 'search' }).fill('salvager');
	await page.getByRole('button', { name: 'Humanoids Teamster Salvager #' }).click();

	await page.getByText('of 210kT').first().click();
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
	await page.locator('[data-type="cycle-selected-map-object-button"]').first().click();

	// Salvage #1 should be selected
	await expect(
		page
			.locator('div')
			.filter({ hasText: /^Salvage #1$/ })
			.first()
	).toBeVisible();
	// salvage cargo should be updated
	await page.getByText('Salvage #1 Humanoids Location').first().scrollIntoViewIfNeeded();
	await expect(page.getByText('Ironium 11kT').first()).toBeVisible();
	await expect(page.getByText('Boranium 12kT').first()).toBeVisible();
	await expect(page.getByText('Germanium 13kT').first()).toBeVisible();
});
