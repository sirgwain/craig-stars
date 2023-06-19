import { expect, test } from '@playwright/test';

test('login page has expected h2', async ({ page }) => {
	await page.goto('/');
	expect(await page.textContent('h2')).toBe('Login');
});
