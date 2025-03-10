import { expect, test } from './setup';

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
