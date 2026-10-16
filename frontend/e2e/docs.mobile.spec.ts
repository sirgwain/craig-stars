import { test, expect } from '@playwright/test';
import { exerciseGameCreationExample } from './helpers/docs';

test.use({ storageState: { cookies: [], origins: [] } });

test('public docs and interactive setup fit a phone @safari', async ({ page }, testInfo) => {
	await exerciseGameCreationExample(page, testInfo);
	await page.getByLabel('Player guide', { exact: true }).selectOption('/docs/mcp');
	await expect(page.getByRole('heading', { name: 'MCP setup' })).toBeVisible();
	await expect(page).toHaveURL('/docs/mcp');
});
