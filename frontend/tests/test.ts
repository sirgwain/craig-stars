import { expect, test } from '@playwright/test';

test('login page has expected h3', async ({ page }) => {
	await page.goto('/login');
	expect(await page.textContent('h3')).toBe('Login');
});
