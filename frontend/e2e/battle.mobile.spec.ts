import { expect, submitTurn, test } from './setup';

// the 690px battle board is scaled down to fit a phone, with its controls under it
test('battle viewer fits a phone', async ({ testGamePage }) => {
	const { page, gameId } = await testGamePage('Battle Mixed Fleet');
	await submitTurn(page);
	await page.goto(`/games/${gameId}/battles/1`);

	const width = page.viewportSize()?.width ?? 0;
	const board = page.getByTestId('battle-board');
	await expect(board).toBeVisible();
	const box = await board.boundingBox();
	expect(box?.x).toBeGreaterThanOrEqual(0);
	expect((box?.x ?? 0) + (box?.width ?? 0)).toBeLessThanOrEqual(width);
	expect(await page.evaluate(() => document.scrollingElement?.scrollWidth)).toBeLessThanOrEqual(
		width
	);

	// the controls are under the board and big enough to tap
	for (const name of ['begin', 'previous', 'next', 'end']) {
		const button = await page.getByRole('button', { name, exact: true }).boundingBox();
		expect(button?.y).toBeGreaterThan((box?.y ?? 0) + (box?.height ?? 0));
		expect(button?.height).toBeGreaterThanOrEqual(40);
		expect((button?.x ?? 0) + (button?.width ?? 0)).toBeLessThanOrEqual(width);
	}

	// tap through a few phases
	const phase = page.getByRole('spinbutton', { name: 'phase' });
	await page.getByRole('button', { name: 'next attack' }).tap();
	await expect(phase).not.toHaveValue('0');
	await page.getByRole('button', { name: 'Play', exact: true }).tap();
	await expect(page.getByRole('button', { name: /pause/i })).toBeVisible();
	await page.getByRole('button', { name: /pause/i }).tap();
});
