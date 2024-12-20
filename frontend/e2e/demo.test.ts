import { expect, test } from '@playwright/test';

test('home page has expected div', async ({ page }) => {
	await page.goto('/');
	await expect(page.locator('body > div')).toBeVisible();
});
