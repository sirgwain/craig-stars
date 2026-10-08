import { fromJson } from '@bufbuild/protobuf';
import { type Page } from '@playwright/test';
import {
	CreateGameResponseSchema,
	type CreateGameResponseJson
} from '../src/lib/protogen/craig_stars/v1/gameservice_pb';
import { apiErrorsFailTest, expect, test } from './setup';

// create a game where the host controls players 1 and 2
async function createHotSeatGame(page: Page, name: string): Promise<bigint> {
	await page.getByRole('link', { name: 'Single Player' }).click();
	await page.getByRole('textbox', { name: 'Name', exact: true }).fill(name);
	await page.getByLabel('Size').selectOption('Tiny');
	await page.getByLabel('Density').selectOption('Sparse');
	await page.getByLabel('Player 2').selectOption({ label: 'Host' });
	await page.locator('[data-type="delete-button"][data-id="Player 3"]').click();

	const createGameResponsePromise = page.waitForResponse(
		(resp) => resp.url().includes('/api/grpc/craig_stars.v1.GameService/CreateGame') && resp.ok()
	);
	await page.getByRole('button', { name: 'Create Game' }).click();
	const { game } = fromJson(
		CreateGameResponseSchema,
		(await (await createGameResponsePromise).json()) as CreateGameResponseJson
	);
	if (!game?.game?.id) {
		throw new Error('failed to create game');
	}
	return game.game.id;
}

test('hot seat game', async ({ authenticatedPage: page }) => {
	apiErrorsFailTest(page);
	const name = `Hot Seat Game ${Date.now()}`;
	const gameId = await createHotSeatGame(page, name);
	const gameLink = page.locator('[data-type="game-link"]').first();
	const playButtons = page.locator('[data-type="play-button"]');
	const submitTurnButton = page.getByRole('button', { name: 'Submit Turn' });

	// opening a hot seat game starts with choosing a player
	await page.waitForURL(`/games/${gameId}`);
	await expect(gameLink).toHaveText(`${name} - 2400`);
	await expect(page.getByText('Choose a player to play')).toBeVisible();
	await expect(playButtons).toHaveCount(2);
	await expect(submitTurnButton).not.toBeVisible();

	// play the first player
	await playButtons.first().click();
	await page.waitForURL(`/games/${gameId}?asPlayer=1`);
	await submitTurnButton.click();

	// after submitting, wait for the other player, who can be played from here
	await expect(page.getByText('Waiting for players to play')).toBeVisible();
	await expect(playButtons).toHaveCount(1);

	// the submitted turn can be viewed read-only
	await page.getByRole('button', { name: 'View Turn (read-only)' }).click();
	await expect(page.getByText('Viewing (read-only)')).toBeVisible();
	await expect(submitTurnButton).not.toBeVisible();
	await page.getByRole('button', { name: 'Done' }).click();
	await expect(page.getByText('Waiting for players to play')).toBeVisible();

	// play the second player, submitting generates the next turn
	await playButtons.first().click();
	await page.waitForURL(`/games/${gameId}?asPlayer=2`);
	await submitTurnButton.click();

	// a new turn starts with choosing a player again
	await expect(gameLink).toHaveText(`${name} - 2401`);
	await expect(page.getByText('Choose a player to play')).toBeVisible();
	await expect(playButtons).toHaveCount(2);
	await expect(page).toHaveURL(`/games/${gameId}`);
});
