import { expect, submitTurn, test } from './setup';

// Battle Mixed Fleet puts two fleets together in deep space; they fight when the turn is generated.
test('battle viewer steps through a battle and plays it', async ({ testGamePage }) => {
	const { page, gameId } = await testGamePage('Battle Mixed Fleet');
	await submitTurn(page);
	await page.goto(`/games/${gameId}/battles/1`);

	const phase = page.getByRole('spinbutton', { name: 'phase' });
	const style = page.getByRole('combobox', { name: 'battle style' });

	// the classic viewer: step to the first attack, which marks the damage on its target
	await expect(style).toHaveValue('classic');
	await page.getByRole('button', { name: 'next attack' }).click();
	await expect(phase).not.toHaveValue('0');
	await expect(page.getByTestId('battle-damage')).toBeVisible();
	await expect(page.getByLabel(/(shield|armor) damage/).first()).toBeVisible();

	// the original playback starts on its own, and is remembered for the next battle
	await style.selectOption('Original');
	await expect(page.getByRole('button', { name: /pause/i })).toBeVisible();
	await page.getByRole('button', { name: '4x' }).click();
	await expect.poll(async () => Number(await phase.inputValue())).toBeGreaterThan(3);
	await page.getByRole('button', { name: /pause/i }).click();
	await expect(page.getByRole('button', { name: /^play/i }).first()).toBeVisible();

	await page.reload();
	await expect(style).toHaveValue('original');
	await expect(page.getByRole('button', { name: /pause/i })).toBeVisible();
});
