import { expect, test } from './setup';

test('all games', async ({ authenticatedPage: page }) => {
	await page.locator('#menu svg').click();
	await page.getByRole('link', { name: 'All Games' }).click();

	await page.waitForURL(`/admin/games`);

	// wait for races to show up
	await expect(page.getByRole('link', { name: 'Single Unit Game' })).toBeVisible();
});

test('users', async ({ authenticatedPage: page }) => {
	await page.locator('#menu svg').click();
	await page.getByRole('link', { name: 'Users' }).click();
	await page.waitForURL(`/admin/users`);

	await expect(page.getByRole('cell', { name: 'admin', exact: true })).toBeVisible();
});
