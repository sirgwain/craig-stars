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
		const name = 'Test Game';
		await authenticatedPage.getByRole('textbox', { name: 'Name', exact: true }).fill(name);
		await authenticatedPage.getByLabel('Size').selectOption('Tiny');
		await authenticatedPage.getByLabel('Density').selectOption('Sparse');
		await authenticatedPage.getByRole('checkbox', { name: 'Public Player Scores' }).click();

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

		const gameLink = authenticatedPage.getByRole('link', { name: name });
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
		await authenticatedPage.getByRole('link', { name: 'Create' }).click();

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
	const gameLink = page.getByRole('link', { name: name });
	await expect(gameLink).toBeVisible();

	// fail if any api calls to this game fail
	const gameId = await gameLink.getAttribute('data-id');
	apiErrorsFailTest(page, gameId);

	// open the game
	await gameLink.click();

	const playerResponse = await page.waitForResponse(
		(resp) =>
			resp.url().includes('/api/grpc/craig_stars.v1.PlayerService/GetPlayer') &&
			resp.request().method() === 'POST' &&
			resp.status() === 200
	);

	const universeResponse = await page.waitForResponse(
		(resp) =>
			resp.url().includes('/api/grpc/craig_stars.v1.PlayerService/GetUniverse') &&
			resp.request().method() === 'POST' &&
			resp.status() === 200
	);

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
	await page.getByRole('button', { name: 'Submit Turn' }).click();

	// wait for turn submit to finish
	const submitTurnResponse = await page.waitForResponse(
		(resp) =>
			resp.url().includes('/api/grpc/craig_stars.v1.PlayerService/SubmitTurn') &&
			resp.request().method() === 'POST' &&
			resp.status() === 200
	);

	const universeResponse = await page.waitForResponse(
		(resp) =>
			resp.url().includes('/api/grpc/craig_stars.v1.PlayerService/GetUniverse') &&
			resp.request().method() === 'POST' &&
			resp.status() === 200
	);

	const playerResponse = await page.waitForResponse(
		(resp) =>
			resp.url().includes('/api/grpc/craig_stars.v1.PlayerService/GetPlayer') &&
			resp.request().method() === 'POST' &&
			resp.status() === 200
	);

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

	// wait for turn submit to finish
	await page.locator('#loading-modal').waitFor({ state: 'visible' }); // wait for loading modal to show up
	await expect(page.locator('#loading-modal')).not.toHaveClass(/modal-open/); // ensure submit is done

	return { game, player, universe };
}

// no js errors allowed
test.beforeEach(async ({ page }) => {
	page.on('pageerror', (err) => {
		throw new Error(`🚨 JavaScript error in the browser: ${err.message}`);
	});
});

export { expect };
