import { fromJson } from '@bufbuild/protobuf';
import { test, expect } from './setup';
import { NewGamePlayerType } from '../src/lib/types/cs-proto';
import { exerciseGameCreationExample } from './helpers/docs';
import { CreateGameResponseSchema } from '../src/lib/protogen/craig_stars/v1/gameservice_pb';

// A new browser context has no session, even though other E2E projects use cached auth.
test.describe('public documentation', () => {
	test.use({ storageState: { cookies: [], origins: [] } });

	test('game creation example works anonymously without API calls', async ({ page }, testInfo) => {
		await exerciseGameCreationExample(page, testInfo);
	});

	test('docs navigation preserves MCP and initializes auth when leaving docs', async ({ page }) => {
		await page.goto('/docs');
		const navigation = page.getByRole('navigation', { name: 'Documentation', exact: true });
		await navigation.getByRole('link', { name: 'MCP setup' }).click();
		await expect(page).toHaveURL('/docs/mcp');
		await expect(page.getByRole('heading', { name: 'MCP setup' })).toBeVisible();
		await navigation.getByRole('link', { name: 'Create a game' }).click();
		await page.getByRole('link', { name: 'Galaxy settings', exact: true }).click();
		await expect(page).toHaveURL(/#galaxy-settings$/);
		await expect(page.getByRole('heading', { name: 'Galaxy settings' })).toBeInViewport();
		const auth = page.waitForResponse((response) => response.url().includes('UserService/GetMe'));
		await page.getByRole('link', { name: 'Single Player', exact: true }).click();
		await auth;
		await expect(page.getByRole('button', { name: "I'm an admin" })).toBeVisible();
	});
});

test('signed-in users can open game creation from docs', async ({ page }) => {
	await page.goto('/docs/games/create');
	await page.getByRole('link', { name: 'Single Player', exact: true }).click();
	await expect(page.getByRole('button', { name: 'Create Game', exact: true })).toBeVisible();
	await expect(page.getByLabel(/^Host/)).toContainText('Humanoids');
});

test('host game creates a setup lobby with the shared form', async ({
	authenticatedPage: page
}) => {
	await page.goto('/host-game');
	await expect(page.getByLabel('Player 2')).toHaveValue(`${NewGamePlayerType.OPEN}`);
	await page
		.getByRole('textbox', { name: 'Name', exact: true })
		.fill('Documentation hosting regression');
	const created = page.waitForResponse(
		(response) => response.url().includes('GameService/CreateGame') && response.ok()
	);
	await page.getByRole('button', { name: 'Create Game', exact: true }).click();
	const { game } = fromJson(CreateGameResponseSchema, await (await created).json());
	if (!game?.game?.id) throw new Error('No hosted game returned');
	try {
		await expect(page).toHaveURL(`/games/${game.game.id}`);
		await expect(page.getByRole('button', { name: 'Start Game', exact: true })).toBeDisabled();
		await expect(page.getByRole('button', { name: 'Add Open Slot' })).toBeVisible();
	} finally {
		await page.request.post('/api/grpc/craig_stars.v1.GameService/DeleteGame', {
			data: { gameId: `${game.game.id}` }
		});
	}
});
