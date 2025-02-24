import { expect, test } from './setup';

test('create a new game', async ({ newGamePage }) => {
	const { page, name } = newGamePage;
	const gameLink = page.getByRole('link', { name: name });
	await expect(gameLink).toBeVisible();
	await expect(gameLink).toHaveText(`${name} - 2400`);
});

test('submit turn', async ({ newGamePage }) => {
	const { page, id, name } = newGamePage;

	page.on('response', async (response) => {
		if (response.url().includes(`/api/games/${id}`) && !response.ok()) {
			// fail any api requests
			throw new Error(`API request failed: ${response.url()} - Status: ${response.status()}`);
		}
	});

	// start with a new game, ensure we have year 2400
	const gameLink = page.getByRole('link', { name: name });
	await expect(gameLink).toBeVisible();
	await expect(gameLink).toHaveText(`${name} - 2400`);

	await page.getByRole('button', { name: 'Submit Turn' }).click();

	await expect(gameLink).toBeVisible();
	await expect(gameLink).toHaveText(`${name} - 2401`);
});

test('research page', async ({ newGamePage }) => {
	const { page, id } = newGamePage;

	page.on('response', async (response) => {
		if (response.url().includes(`/api/games/${id}`) && !response.ok()) {
			// fail any api requests
			throw new Error(`API request failed: ${response.url()} - Status: ${response.status()}`);
		}
	});

	await page.locator('label').filter({ hasText: 'Commands' }).click();
	await page.getByRole('link', { name: 'Research' }).click();
	await page.getByRole('radio', { name: 'Weapons' }).click();
	// select weapons
	await expect(page.getByText('Weapons 4')).toBeVisible();

	// select prop
	await page.getByRole('radio', { name: 'Propulsion' }).click();
	await expect(page.getByText('Propulsion 4')).toBeVisible();

	// select con
	await page.getByRole('radio', { name: 'Construction' }).click();
	await expect(page.getByText('Construction 4')).toBeVisible();

	// select elec
	await page.getByRole('radio', { name: 'Electronics' }).click();
	await expect(page.getByText('Electronics 4')).toBeVisible();

	// select energy
	await page.getByRole('radio', { name: 'Energy' }).click();
	await expect(page.getByText('Energy 4')).toBeVisible();

	// increase research
	await page.locator('[data-type="spin-number-increase-button"]').click();
	await expect(page.getByText('Research Budget 16 %')).toBeVisible();

	// decrease research
	await page.locator('[data-type="spin-number-decrease-button"]').click();
	await expect(page.getByText('Research Budget 15 %')).toBeVisible();
});

test('relations page', async ({ newGamePage }) => {
	const { page, id } = newGamePage;

	page.on('response', async (response) => {
		if (response.url().includes(`/api/games/${id}`) && !response.ok()) {
			// fail any api requests
			throw new Error(`API request failed: ${response.url()} - Status: ${response.status()}`);
		}
	});

	await page.locator('label').filter({ hasText: 'Commands' }).click();
	await page.getByRole('link', { name: 'Relations' }).click();

	await page.getByRole('radio', { name: 'Friend' }).first().check();
	await page.getByRole('radio', { name: 'Neutral' }).first().check();
	await page.getByRole('radio', { name: 'Enemy' }).first().check();
	await page.locator('input[name="player-relation-2-share-map"]').check();
});

test('battle plans page', async ({ newGamePage }) => {
	const { page, id } = newGamePage;
	const name = 'Test Battle Plan';

	page.on('response', async (response) => {
		if (response.url().includes(`/api/games/${id}`) && !response.ok()) {
			// fail any api requests
			throw new Error(`API request failed: ${response.url()} - Status: ${response.status()}`);
		}
	});

	await page.locator('label').filter({ hasText: 'Commands' }).click();
	await page.getByRole('link', { name: 'Battle Plans' }).click();

	await page.getByRole('link', { name: 'Create' }).click();
	await page.getByRole('textbox', { name: 'Name' }).fill(name);
	await page.getByLabel('Primary Target').selectOption('Any');
	await page.getByLabel('Secondary Target').selectOption('Starbase');
	await page.getByLabel('Tactic DisengageDisengage If').selectOption('Disengage');
	await page.getByLabel('Attack Who EnemiesEnemies And').selectOption('Everyone');
	await page.getByRole('button', { name: 'Save' }).click();
	await page.getByRole('link', { name: 'Battle Plans' }).click();

	// delete the plan we just created
	const deleteButton = await page.locator(`[data-type="delete-button"][data-id="${name}"]`);
	await expect(deleteButton).toBeVisible();

	// confirm delete
	page.once('dialog', async (dialog) => {
		await dialog.accept();
	});
	await deleteButton.click();
	await expect(page.getByRole('link', { name: name })).not.toBeVisible();
});

test('production plans page', async ({ newGamePage }) => {
	const { page, id } = newGamePage;
	const name = 'Test Production Plan';

	page.on('response', async (response) => {
		if (response.url().includes(`/api/games/${id}`) && !response.ok()) {
			// fail any api requests
			throw new Error(`API request failed: ${response.url()} - Status: ${response.status()}`);
		}
	});

	await page.locator('label').filter({ hasText: 'Commands' }).click();
	await page.getByRole('link', { name: 'Production Plans' }).click();

	await page.getByRole('link', { name: 'Create' }).click();
	await page.getByRole('textbox', { name: 'Name' }).fill(name);

	await page.getByRole('button', { name: 'Factory (Auto)' }).click();
	await page.getByRole('button', { name: 'Add' }).click();
	await page.getByRole('button', { name: 'Mine (Auto)' }).click();
	await page.getByRole('button', { name: 'Add' }).click();
	await page.getByRole('button', { name: 'Defenses (Auto)' }).click();
	await page.getByRole('button', { name: 'Add' }).click();
	await page.getByRole('button', { name: 'Remove' }).click();
	await page.getByRole('button', { name: 'Item Up' }).click();
	await page.getByRole('button', { name: 'Item Down' }).click();
	await page.getByRole('button', { name: 'Clear' }).click();
	await page.getByRole('button', { name: 'Alchemy (Auto)' }).click();
	await page.getByRole('button', { name: 'Add' }).click();
	await page.getByText('Contributes Only Leftover to').click();

	await page.getByRole('button', { name: 'Save' }).click();
	await page.getByRole('link', { name: 'Production Plans' }).click();

	// delete the plan we just created
	const deleteButton = await page.locator(`[data-type="delete-button"][data-id="${name}"]`);
	await expect(deleteButton).toBeVisible();

	// confirm delete
	page.once('dialog', async (dialog) => {
		await dialog.accept();
	});
	await deleteButton.click();
	await expect(page.getByRole('link', { name: name })).not.toBeVisible();
});
