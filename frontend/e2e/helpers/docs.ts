import { expect, type Page } from '@playwright/test';
import { NewGamePlayerType, Size } from '../../src/lib/types/cs-proto';

export async function exerciseGameCreationExample(page: Page) {
	const apiCalls: string[] = [];
	page.on('request', (request) => {
		if (request.url().includes('/api/')) apiCalls.push(request.url());
	});
	await page.goto('/docs/games/create');
	await expect(page.getByRole('heading', { name: 'Create a game', exact: true })).toBeVisible();

	const example = page.getByRole('region', { name: 'Try game creation' });
	await example.getByRole('textbox', { name: 'Name', exact: true }).fill('Practice galaxy');
	await example.getByLabel('Size').selectOption({ label: 'Tiny' });
	await example.getByLabel(/^Host/).selectOption({ label: 'Example Settlers' });
	await example.getByLabel('Public Player Scores').check();
	await example.locator('input[type="number"]').first().fill('80');
	await example.getByRole('button', { name: 'Add player' }).click();
	await example.getByLabel('Player 4').selectOption(`${NewGamePlayerType.GUEST}`);
	await example.getByRole('button', { name: 'Preview setup' }).click();
	const preview = example.locator('[data-type="setup-preview"]');
	await expect(preview).toContainText('Practice galaxy');
	await expect(preview).toContainText('4 players');
	await expect(preview).toContainText('Example Settlers');
	await expect(preview).toContainText('Guest');
	await expect(preview).toContainText('No game was created');

	await example.getByRole('button', { name: 'Reset', exact: true }).click();
	await expect(example.getByRole('textbox', { name: 'Name', exact: true })).toHaveValue(
		'A Barefoot Jaywalk'
	);
	await expect(example.getByLabel('Size')).toHaveValue(`${Size.SMALL}`);
	await expect(example.getByLabel('Public Player Scores')).not.toBeChecked();
	await expect(example.getByLabel(/^Host/)).toHaveValue('1');
	await expect(example.locator('input[type="number"]').first()).toHaveValue('60');
	await expect(example.getByLabel('Player 4')).toHaveCount(0);
	await expect(preview).toHaveCount(0);

	await example.getByLabel('Example setup').selectOption('host');
	await expect(example.getByLabel('Player 2')).toHaveValue(`${NewGamePlayerType.OPEN}`);
	await example.getByRole('button', { name: 'Preview setup' }).click();
	await expect(preview).toContainText('Player 2: Open');
	await example.locator('[data-type="delete-button"][data-id="Player 3"]').click();
	await example.getByRole('button', { name: 'Preview setup' }).click();
	await expect(preview).toContainText('2 players');

	const width = await page.evaluate(() => ({
		viewport: window.innerWidth,
		content: document.documentElement.scrollWidth
	}));
	expect(width.content).toBeLessThanOrEqual(width.viewport);
	expect(apiCalls).toEqual([]);
}
