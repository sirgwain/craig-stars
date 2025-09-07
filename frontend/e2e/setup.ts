import { create, fromJson } from '@bufbuild/protobuf';
import { test as base, expect, Page } from '@playwright/test';
// Import for direct HTTP calls instead of Connect client
import { toJson } from '@bufbuild/protobuf';
import { type GameWithPlayers } from '../src/lib/protogen/craig_stars/v1/game_pb';
import {
	CreateGameResponseJson,
	CreateGameResponseSchema
} from '../src/lib/protogen/craig_stars/v1/gameservice_pb';
import { type Player, type PlayerUniverse } from '../src/lib/protogen/craig_stars/v1/player_pb';
import {
	GetPlayerResponseSchema,
	GetUniverseResponseJson,
	GetUniverseResponseSchema,
	SubmitTurnResponseSchema,
	type GetPlayerResponseJson,
	type SubmitTurnResponseJson
} from '../src/lib/protogen/craig_stars/v1/playerservice_pb';
import {
	CreateTestGameRequestSchema,
	CreateTestGameResponseSchema
} from '../src/lib/protogen/craig_stars/v1/testservice_pb';

let gameNum = 1;

// Helper function to create test games via authenticated page request
export async function createTestGame(
	page: Page,
	testGameName: string,
	gameName: string
): Promise<GameWithPlayers> {
	const request = create(CreateTestGameRequestSchema, {
		testGameName: testGameName,
		gameName: gameName
	});

	const response = await page.request.post('/api/grpc/craig_stars.v1.TestService/CreateTestGame', {
		headers: {
			'Content-Type': 'application/json'
		},
		data: toJson(CreateTestGameRequestSchema, request)
	});

	if (!response.ok()) {
		throw new Error(`Failed to create test game: ${response.status()} ${response.statusText()}`);
	}

	const responseJson = await response.json();
	const createResponse = fromJson(CreateTestGameResponseSchema, responseJson);

	if (!createResponse.game) {
		throw new Error('No game returned from CreateTestGame');
	}

	return createResponse.game;
}

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
	testGamePage: (testGameName: string) => Promise<{
		page: Page;
		gameId: bigint;
		universe: PlayerUniverse;
		player: Player;
	}>;
}>({
	authenticatedPage: async ({ page }, use) => {
		// Navigate to homepage - authentication is already cached via storageState
		await page.goto('/');

		// Wait for the homepage to be visible to confirm we're logged in
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

		const gameLink = await authenticatedPage.locator('[data-type="game-link"]').first();
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
	},

	testGamePage: async ({ authenticatedPage }, use) => {
		const createdGames: bigint[] = [];

		const loadTestGame = async (testGameName: string) => {
			const result = await loadTestGamePage(authenticatedPage, testGameName);
			createdGames.push(result.gameId);
			return result;
		};

		await use(loadTestGame);

		// cleanup all created test games
		for (const gameId of createdGames) {
			try {
				await authenticatedPage.goto('/');
				const deleteButton = await authenticatedPage.locator(
					`[data-type="delete-button"][data-id="${gameId}"]`
				);
				if (await deleteButton.isVisible({ timeout: 1000 }).catch(() => false)) {
					// confirm delete
					authenticatedPage.once('dialog', async (dialog) => {
						await dialog.accept();
					});

					await deleteButton.click();
					await expect(deleteButton).not.toBeVisible();
				}
			} catch (error) {
				console.warn(`Failed to cleanup test game ${gameId}:`, error);
			}
		}
	}
});

export async function loadTestGamePage(page: Page, testGameName: string) {
	const name = `Test Game ${gameNum++}`;

	const game = await createTestGame(page, testGameName, name);

	if (!game.game) {
		throw new Error('failed to create test game');
	}

	// fail if any api calls to this game fail
	apiErrorsFailTest(page);

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
	const gameId = game.game.id;
	await page.goto(`/games/${gameId}`);

	// wait for both to finish, order doesn't matter
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

export async function apiErrorsFailTest(page: Page) {
	page.on('response', async (response) => {
		if (response.url().includes(`/api/grpc/craig_stars.v1`) && !response.ok()) {
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
