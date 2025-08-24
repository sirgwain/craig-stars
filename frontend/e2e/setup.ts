import { fromJson } from '@bufbuild/protobuf';
import { test as base, expect, Page } from '@playwright/test';
import { type GameWithPlayers } from '../src/lib/protogen/craig_stars/v1/game_pb';
import {
	CreateGameResponseJson,
	CreateGameResponseSchema
} from '../src/lib/protogen/craig_stars/v1/gameservice_pb';
import { type PlayerUniverse } from '../src/lib/protogen/craig_stars/v1/player_pb';
import {
	GetPlayerResponseSchema,
	GetUniverseResponseJson,
	GetUniverseResponseSchema,
	SubmitTurnResponseSchema,
	type GetPlayerResponseJson,
	type SubmitTurnResponseJson
} from '../src/lib/protogen/craig_stars/v1/playerservice_pb';

let gameNum = 1;

export const test = base.extend<{
	authenticatedPage: Page;
	newGamePage: {
		page: Page;
		id: string;
		name: string;
		game: GameWithPlayers;
		universe: PlayerUniverse;
	};
	newRacePage: { page: Page; id: string; name: string };
}>({
	authenticatedPage: async ({ page }, use) => {
		// Navigate to login page
		await page.goto('/');
		await page.getByRole('button', { name: "I'm an admin" }).click();
		await page.getByRole('textbox', { name: 'Username' }).click();
		await page.getByRole('textbox', { name: 'Username' }).fill('admin');
		await page.getByRole('textbox', { name: 'Username' }).press('Tab');
		await page.getByRole('textbox', { name: 'Password' }).fill('admin');
		await page.getByRole('button', { name: 'Submit' }).click();

		// Wait for the homepage to be visible
		await expect(page.getByRole('link', { name: 'Single Player' })).toBeVisible();

		// Use this logged-in page in tests
		await use(page);
	},

	newGamePage: async ({ authenticatedPage }, use) => {
		await authenticatedPage.getByRole('link', { name: 'Single Player' }).click();

		// fill in some fields
		const name = `Test Game ${gameNum++}`;
		await authenticatedPage.getByRole('textbox', { name: 'Name', exact: true }).fill(name);
		await authenticatedPage.getByLabel('Size').selectOption('Tiny');
		await authenticatedPage.getByLabel('Density').selectOption('Sparse');
		await authenticatedPage.getByRole('checkbox', { name: 'Public Player Scores' }).click();
		await authenticatedPage.locator('[data-type="delete-button"][data-id="Player 3"]').click();

		authenticatedPage.getByRole('button', { name: 'Create Game' }).click();
		const response = await authenticatedPage.waitForResponse(
			(response) =>
				response.url().includes('/api/grpc/craig_stars.v1.GameService') &&
				response.request().method() === 'POST' &&
				response.status() === 200
		);

		const { game } = fromJson(
			CreateGameResponseSchema,
			(await response.json()) as CreateGameResponseJson
		);
		if (!game?.game?.id) {
			throw new Error('failed to create game');
		}

		const universeResponse = await authenticatedPage.waitForResponse(
			(resp) =>
				resp.url().includes('/api/grpc/craig_stars.v1.PlayerService/GetUniverse') &&
				resp.request().method() === 'POST' &&
				resp.status() === 200
		);

		const { universe } = fromJson(
			GetUniverseResponseSchema,
			(await universeResponse.json()) as GetUniverseResponseJson
		);
		if (!universe) {
			throw new Error('failed to get universe for game');
		}
		await authenticatedPage.waitForURL(`/games/${game.game.id}`);

		const gameLink = authenticatedPage.getByRole('link', {
			name: `${name} - ${game.game.year}`,
			exact: true
		});
		await expect(gameLink).toBeVisible();
		await expect(gameLink).toHaveText(`${name} - 2400`);

		// do whatever our subtest wants
		await use({
			page: authenticatedPage,
			id: `${game.game.id}`,
			name: name,
			game: game,
			universe: universe
		});

		// delete the game
		await authenticatedPage.goto('/');
		const deleteButton = await authenticatedPage.locator(
			`[data-type="delete-button"][data-id="${game.game.id}"]`
		);
		await expect(deleteButton).toBeVisible();

		// confirm delete
		authenticatedPage.once('dialog', async (dialog) => {
			await dialog.accept();
		});

		await deleteButton.click();
		await expect(deleteButton).not.toBeVisible();
	},

	newRacePage: async ({ authenticatedPage }, use) => {
		await authenticatedPage.getByRole('link', { name: 'Races' }).click();
		await authenticatedPage.waitForURL(`/races`);

		await authenticatedPage.getByRole('link', { name: 'Create' }).click();
		await authenticatedPage.waitForURL(`/races/new`);

		const name = 'Test Race';
		await authenticatedPage.getByRole('textbox', { name: 'Name', exact: true }).fill(name);
		await authenticatedPage.getByRole('textbox', { name: 'Plural Name' }).fill(name + 's');

		authenticatedPage.getByRole('button', { name: 'Save' }).click();
		const response = await authenticatedPage.waitForResponse(
			(response) =>
				response.url().includes('/api/grpc/craig_stars.v1.RaceService') &&
				response.request().method() === 'POST' &&
				response.status() === 200
		);

		const { race } = await response.json();

		// do whatever our subtest wants
		await use({ page: authenticatedPage, id: race.id, name: name });

		// delete the game
		await authenticatedPage.goto('/races');
		const deleteButton = await authenticatedPage.locator(
			`[data-type="delete-button"][data-id="${race.id}"]`
		);
		await expect(deleteButton).toBeVisible();

		// confirm delete
		authenticatedPage.once('dialog', async (dialog) => {
			await dialog.accept();
		});

		// delete race we created
		await deleteButton.click();
		await expect(deleteButton).not.toBeVisible();
	}
});

export async function loadGamePage(page: Page, name: string) {
	const gameLink = page.getByRole('link', { name, exact: true });
	await expect(gameLink).toBeVisible();

	// fail if any api calls to this game fail
	const gameId = await gameLink.getAttribute('data-id');
	apiErrorsFailTest(page, gameId);

	// kick off the response waiters before clicking
	const playerResponsePromise = page.waitForResponse(
		(resp) =>
			resp.url().includes('/api/grpc/craig_stars.v1.PlayerService/GetPlayer') &&
			resp.request().method() === 'POST' &&
			resp.status() === 200
	);
	const universeResponsePromise = page.waitForResponse(
		(resp) =>
			resp.url().includes('/api/grpc/craig_stars.v1.PlayerService/GetUniverse') &&
			resp.request().method() === 'POST' &&
			resp.status() === 200
	);

	// open the game (this will trigger both requests)
	await gameLink.click();

	// wait for both to finish, order doesn’t matter
	const [playerResponse, universeResponse] = await Promise.all([
		playerResponsePromise,
		universeResponsePromise
	]);

	const { universe } = fromJson(
		GetUniverseResponseSchema,
		(await universeResponse.json()) as GetUniverseResponseJson
	);

	const { player } = fromJson(
		GetPlayerResponseSchema,
		(await playerResponse.json()) as GetPlayerResponseJson
	);

	if (!player || !universe) {
		throw new Error('failed to load universe and player');
	}

	await page.waitForURL(`/games/${gameId}`);
	await expect(page.locator(`[data-type="game-view"][data-id="${gameId}"]`)).toBeVisible();

	return { page, gameId, universe, player };
}

export async function apiErrorsFailTest(page: Page, gameId: string | null) {
	if (gameId === null) {
		throw new Error(`invalid gameId for page ${page.url()}`);
	}
	page.on('response', async (response) => {
		if (response.url().includes(`/api/grpc/craig_stars.v1.GameService`) && !response.ok()) {
			// fail any api requests
			throw new Error(`API request failed: ${response.url()} - Status: ${response.status()}`);
		}
	});
}

export async function submitTurn(page: Page) {
	// start waiting for all three before clicking
	const submitTurnResponsePromise = page.waitForResponse(
		(resp) =>
			resp.url().includes('/api/grpc/craig_stars.v1.PlayerService/SubmitTurn') &&
			resp.request().method() === 'POST' &&
			resp.status() === 200
	);

	const universeResponsePromise = page.waitForResponse(
		(resp) =>
			resp.url().includes('/api/grpc/craig_stars.v1.PlayerService/GetUniverse') &&
			resp.request().method() === 'POST' &&
			resp.status() === 200
	);

	const playerResponsePromise = page.waitForResponse(
		(resp) =>
			resp.url().includes('/api/grpc/craig_stars.v1.PlayerService/GetPlayer') &&
			resp.request().method() === 'POST' &&
			resp.status() === 200
	);

	// trigger the requests
	await page.getByRole('button', { name: 'Submit Turn' }).click();

	// await submit turn (must finish first)
	const submitTurnResponse = await submitTurnResponsePromise;

	// then wait for the others in parallel
	const [universeResponse, playerResponse] = await Promise.all([
		universeResponsePromise,
		playerResponsePromise
	]);

	const { game } = fromJson(
		SubmitTurnResponseSchema,
		(await submitTurnResponse.json()) as SubmitTurnResponseJson
	);

	const { universe } = fromJson(
		GetUniverseResponseSchema,
		(await universeResponse.json()) as GetUniverseResponseJson
	);

	const { player } = fromJson(
		GetPlayerResponseSchema,
		(await playerResponse.json()) as GetPlayerResponseJson
	);

	// wait for turn submit UI flow to finish
	await page.locator('#loading-modal').waitFor({ state: 'visible' });
	await expect(page.locator('#loading-modal')).not.toHaveClass(/modal-open/);

	return { game, player, universe };
}

// no js errors allowed
test.beforeEach(async ({ page }) => {
	page.on('pageerror', (err) => {
		throw new Error(`🚨 JavaScript error in the browser: ${err.message}`);
	});
});

export { expect };
