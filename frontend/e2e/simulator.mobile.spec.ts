import { expect, test } from './setup';

test('simulator fits a phone with separate side editors and full-width replay @safari', async ({
	authenticatedPage: page
}) => {
	await page.goto('/simulator');
	const a = page.getByRole('region', { name: 'Side A', exact: true });
	const b = page.getByRole('region', { name: 'Side B', exact: true });
	await expect(a).toBeVisible();
	await expect(b).toBeHidden();
	await page.getByRole('button', { name: 'B · Opponent', exact: true }).tap();
	await expect(b).toBeVisible();
	await expect(a).toBeHidden();
	await b.getByLabel('Stack 1 quantity', { exact: true }).fill('8');
	expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
	await page.getByRole('button', { name: 'A · Your forces', exact: true }).tap();
	await a.getByRole('button', { name: 'Edit design' }).tap();
	const editor = page.getByRole('dialog', { name: 'Simulation ship designer' });
	await editor.getByRole('textbox', { name: 'Name', exact: true }).fill('Phone prototype');
	expect(await editor.evaluate((el) => el.scrollWidth <= el.clientWidth)).toBe(true);
	await editor.getByRole('button', { name: 'Save simulation design' }).tap();
	await expect(editor).toBeHidden();
	await page.getByLabel('Simulation seed').fill('42');
	await page.getByRole('button', { name: 'Simulate', exact: true }).tap();
	const board = page.getByTestId('battle-board');
	await expect(board).toBeVisible();
	const box = await board.boundingBox();
	expect(box?.x).toBeGreaterThanOrEqual(0);
	expect((box?.x ?? 0) + (box?.width ?? 0)).toBeLessThanOrEqual(page.viewportSize()!.width);
	expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
	await page.getByRole('button', { name: /next attack/i }).tap();
	await expect(page.getByRole('spinbutton', { name: 'phase', exact: true })).not.toHaveValue('0');
});
