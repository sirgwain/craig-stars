import { expect, test } from './setup';

test('create a new game', async ({ newGamePage }) => {
	const { page, name } = newGamePage;
	const gameLink = page.getByRole('link', { name: name });
	await expect(gameLink).toBeVisible();
	await expect(gameLink).toHaveText(`${name} - 2400`);
});

test('submit turn', async ({ newGamePage }) => {
	const { page, name } = newGamePage;
	// start with a new game, ensure we have year 2400
	const gameLink = page.getByRole('link', { name: name });
	await expect(gameLink).toBeVisible();
	await expect(gameLink).toHaveText(`${name} - 2400`);

	// submit turn, wait for year 2401
	await page.getByRole('button', { name: 'Submit Turn' }).click();
	await page.waitForResponse(
		(response) =>
			response.url().includes(`/submit-turn`) &&
			response.request().method() === 'POST' &&
			response.status() === 200
	);

	await expect(gameLink).toBeVisible();
	await expect(gameLink).toHaveText(`${name} - 2401`);
});
