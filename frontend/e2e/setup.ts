import { test as base, expect, Page } from '@playwright/test';

export const test = base.extend<{
	authenticatedPage: Page;
	newGamePage: { page: Page; id: string; name: string };
	newRacePage: { page: Page; id: string; name: string };
}>({
	authenticatedPage: async ({ page }, use) => {
		// Navigate to login page
		await page.goto('/');
		await page.getByRole('button', { name: "I'm an admin" }).click();
		await page.getByRole('textbox', { name: 'Username' }).click();
		await page.getByRole('textbox', { name: 'Username' }).fill('admin');
		await page.getByRole('textbox', { name: 'Username' }).press('Tab');
		await page.getByRole('textbox', { name: 'Password' }).fill('admin');
		await page.getByRole('button', { name: 'Submit' }).click();

		// Wait for the homepage to be visible
		await expect(page.getByRole('link', { name: 'Single Player' })).toBeVisible();

		// Use this logged-in page in tests
		await use(page);
	},

	newGamePage: async ({ authenticatedPage }, use) => {
		await authenticatedPage.getByRole('link', { name: 'Single Player' }).click();

		// fill in some fields
		const name = 'Test Game';
		await authenticatedPage.getByRole('textbox', { name: 'Name', exact: true }).fill(name);
		await authenticatedPage.getByLabel('Size').selectOption('Tiny');
		await authenticatedPage.getByLabel('Density').selectOption('Sparse');
		await authenticatedPage.getByRole('checkbox', { name: 'Public Player Scores' }).click();

		authenticatedPage.getByRole('button', { name: 'Create Game' }).click();
		const response = await authenticatedPage.waitForResponse(
			(response) =>
				response.url().includes('/api/games') &&
				response.request().method() === 'POST' &&
				response.status() === 200
		);

		const { id } = await response.json();

		const gameLink = authenticatedPage.getByRole('link', { name: name });
		await expect(gameLink).toBeVisible();
		await expect(gameLink).toHaveText(`${name} - 2400`);

		// do whatever our subtest wants
		await use({ page: authenticatedPage, id, name: name });

		// delete the game
		await authenticatedPage.goto('/');
		const deleteButton = await authenticatedPage.locator(
			`[data-type="delete-button"][data-id="${id}"]`
		);
		await expect(deleteButton).toBeVisible();

		// confirm delete
		authenticatedPage.once('dialog', async (dialog) => {
			await dialog.accept();
		});

		await deleteButton.click();
		await expect(deleteButton).not.toBeVisible();
	},

	newRacePage: async ({ authenticatedPage }, use) => {
		await authenticatedPage.getByRole('link', { name: 'Races' }).click();
		await authenticatedPage.getByRole('link', { name: 'Create' }).click();

		const name = 'Test Race';
		await authenticatedPage.getByRole('textbox', { name: 'Name', exact: true }).fill(name);
		await authenticatedPage.getByRole('textbox', { name: 'Plural Name' }).fill(name + 's');

		authenticatedPage.getByRole('button', { name: 'Save' }).click();
		const response = await authenticatedPage.waitForResponse(
			(response) =>
				response.url().includes('/api/races') &&
				response.request().method() === 'POST' &&
				response.status() === 200
		);

		const { id } = await response.json();

		// do whatever our subtest wants
		await use({ page: authenticatedPage, id, name: name });

		// delete the game
		await authenticatedPage.goto('/races');
		const deleteButton = await authenticatedPage.locator(
			`[data-type="delete-button"][data-id="${id}"]`
		);
		await expect(deleteButton).toBeVisible();

		// confirm delete
		authenticatedPage.once('dialog', async (dialog) => {
			await dialog.accept();
		});

		// delete race we created
		await deleteButton.click();
		await expect(deleteButton).not.toBeVisible();
	}
});

export async function loadGamePage(page: Page, name: string) {
	const gameLink = page.getByRole('link', { name: name });
	await expect(gameLink).toBeVisible();

	// fail if any api calls to this game fail
	const gameId = await gameLink.getAttribute('data-id');
	apiErrorsFailTest(page, gameId);

	// open the game
	await gameLink.click();
	await expect(page.locator(`[data-type="game-view"][data-id="${gameId}"]`)).toBeVisible();

	return { page, gameId };
}

export async function apiErrorsFailTest(page: Page, gameId: string | null) {
	if (gameId === null) {
		throw new Error(`invalid gameId for page ${page.url()}`);
	}
	page.on('response', async (response) => {
		if (response.url().includes(`/api/games/${gameId}`) && !response.ok()) {
			// fail any api requests
			throw new Error(`API request failed: ${response.url()} - Status: ${response.status()}`);
		}
	});
}

// no js errors allowed
test.beforeEach(async ({ page }) => {
	page.on('pageerror', (err) => {
		throw new Error(`🚨 JavaScript error in the browser: ${err.message}`);
	});
});

export { expect };
