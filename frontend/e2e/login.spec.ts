import { expect, test } from '@playwright/test';

test('login has discord button', async ({ page }) => {
	await page.goto('/');
	await expect(page.locator('#discord-login')).toBeVisible();
});
