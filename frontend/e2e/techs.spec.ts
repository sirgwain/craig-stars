import { expect, test } from './setup';

test('view techs', async ({ authenticatedPage: page }) => {
	await page.getByRole('link', { name: 'Techs' }).click();

	// wait for cards to show up
	await expect(page.getByRole('link', { name: 'Tritanium' })).toBeVisible();

	// should have a page of tech cards
	const techCards = await page.locator('[data-type="tech-card"]');
	const count = await techCards.count();
	await expect(count).toBeGreaterThan(1);
});

test('filter techs', async ({ authenticatedPage: page }) => {
	await page.getByRole('link', { name: 'Techs' }).click();

	// wait for cards to show up
	await expect(page.getByRole('link', { name: 'Tritanium' })).toBeVisible();

	await page.getByRole('textbox', { name: 'search' }).click();
	await page.getByRole('textbox', { name: 'search' }).fill('tritanium');
	await expect(page.getByRole('link', { name: 'Scout' })).toBeHidden();

	// should have only one card
	const techCards = await page.locator('[data-type="tech-card"]');
	await expect(techCards).toHaveCount(1);
});

test('browse to tritanium', async ({ authenticatedPage: page }) => {
	await page.getByRole('link', { name: 'Techs' }).click();

	// wait for cards to show up
	await expect(page.getByRole('link', { name: 'Tritanium' })).toBeVisible();

	// click tritanium
	await page.getByRole('link', { name: 'Tritanium' }).click();
	await page.waitForURL('/techs/tritanium');

	await expect(page.getByText('Armor Strength: 50')).toBeVisible();
});

test('browse to scout', async ({ authenticatedPage: page }) => {
	await page.getByRole('link', { name: 'Techs' }).click();

	// wait for cards to show up
	await expect(page.getByRole('link', { name: 'Scout' })).toBeVisible();

	// click scout
	await page.getByRole('link', { name: 'Scout' }).click();
	await page.waitForURL('/techs/scout');

	// should have a hull view
	await expect(page.getByRole('heading', { name: 'Hull' })).toBeVisible();
});
