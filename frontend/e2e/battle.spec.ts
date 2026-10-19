import { expect, submitTurn, test } from './setup';

// Battle Mixed Fleet puts two fleets together in deep space; they fight when the turn is generated.
test('battle viewer steps through a battle and plays it', async ({ testGamePage }) => {
	const { page, gameId } = await testGamePage('Battle Mixed Fleet');
	await submitTurn(page);
	await page.goto(`/games/${gameId}/battles/1`);

	const phase = page.getByRole('spinbutton', { name: 'phase' });
	const animations = page.getByRole('checkbox', { name: 'Animations' });
	const play = page.getByRole('button', { name: 'Play', exact: true });
	const pause = page.getByRole('button', { name: 'Pause' });

	// it waits at the start until it is played
	await expect(play).toBeVisible();
	await expect(phase).toHaveValue('0');

	// the next attack marks the damage on its target
	await page.getByRole('button', { name: 'next attack' }).click();
	await expect(phase).not.toHaveValue('0');
	await expect(page.getByTestId('battle-damage').first()).toBeVisible();
	await expect(page.getByLabel(/(shield|armor) damage/).first()).toBeVisible();

	await page.getByRole('button', { name: 'begin' }).click();
	await page.getByRole('button', { name: '4x' }).click();
	await play.click();
	await expect(pause).toBeVisible();
	await expect.poll(async () => Number(await phase.inputValue())).toBeGreaterThan(3);
	await pause.click();
	await expect(play).toBeVisible();

	// the next attack stops playback and steps through one attack at a time
	await play.click();
	await expect(pause).toBeVisible();
	await page.getByRole('button', { name: 'next attack' }).click();
	await expect(play).toBeVisible();

	// without animations, the next attack lands right away
	await animations.uncheck();
	const before = Number(await phase.inputValue());
	await page.getByRole('button', { name: 'next attack' }).click();
	await expect.poll(async () => Number(await phase.inputValue())).toBeGreaterThan(before);
	await expect(play).toBeVisible();

	// the speed and the animations are remembered for the next battle, which waits to be played
	await page.reload();
	await expect(animations).not.toBeChecked();
	await expect(page.getByRole('button', { name: '4x' })).toHaveAttribute('aria-pressed', 'true');
	await expect(play).toBeVisible();
});
