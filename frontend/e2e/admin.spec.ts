import { create, toJson } from '@bufbuild/protobuf';
import { DeleteGameRequestSchema } from '../src/lib/protogen/craig_stars/v1/gameservice_pb';
import { createTestGame, expect, test } from './setup';

test('all games', async ({ authenticatedPage: page }) => {
	const game1 = await createTestGame(page, 'Single Unit Game', 'All Games Test Game 1');
	const game2 = await createTestGame(page, 'Single Unit Game', 'All Games Test Game 2');

	if (!game1.game?.id || !game2.game?.id) {
		throw new Error('failed to create games for all games test');
	}

	await page.locator('#menu svg').click();
	await page.getByRole('link', { name: 'All Games' }).click();

	await page.waitForURL(`/admin/games`);

	// wait for races to show up
	await expect(page.getByRole('link', { name: game1.game?.name })).toBeVisible();
	await expect(page.getByRole('link', { name: game2.game?.name })).toBeVisible();

	const response1 = await page.request.post('/api/grpc/craig_stars.v1.GameService/DeleteGame', {
		headers: {
			'Content-Type': 'application/json'
		},
		data: toJson(
			DeleteGameRequestSchema,
			create(DeleteGameRequestSchema, { gameId: game1.game.id })
		)
	});

	if (!response1.ok) {
		throw new Error('failed to delete game1 for all games test');
	}

	const response2 = await page.request.post('/api/grpc/craig_stars.v1.GameService/DeleteGame', {
		headers: {
			'Content-Type': 'application/json'
		},
		data: toJson(
			DeleteGameRequestSchema,
			create(DeleteGameRequestSchema, { gameId: game2.game.id })
		)
	});

	if (!response2.ok) {
		throw new Error('failed to delete game2 for all games test');
	}
});

test('users', async ({ authenticatedPage: page }) => {
	await page.locator('#menu svg').click();
	await page.getByRole('link', { name: 'Users' }).click();
	await page.waitForURL(`/admin/users`);

	await expect(page.getByRole('cell', { name: 'admin', exact: true })).toBeVisible();
});
