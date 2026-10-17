import path from 'node:path';
import { expect, test } from './setup';
import type { Page } from '@playwright/test';
import { fromJson, toJson } from '@bufbuild/protobuf';
import {
	PlayerSchema,
	PlayerUniverseSchema,
	GetPlayerResponseSchema,
	GetUniverseResponseSchema
} from '../src/lib/types/cs-proto';

async function getPlayer(page: Page, gameId: bigint) {
	const response = await page.request.post('/api/grpc/craig_stars.v1.PlayerService/GetPlayer', {
		data: { gameId: gameId.toString() }
	});
	expect(response.ok()).toBeTruthy();
	return fromJson(GetPlayerResponseSchema, await response.json()).player!;
}
async function getUniverse(page: Page, gameId: bigint) {
	const response = await page.request.post('/api/grpc/craig_stars.v1.PlayerService/GetUniverse', {
		data: { gameId: gameId.toString() }
	});
	expect(response.ok()).toBeTruthy();
	return fromJson(GetUniverseResponseSchema, await response.json()).universe!;
}

// Screenshots are report attachments, not checked-in snapshot baselines.
test.beforeEach(async ({ page }) => {
	await page.addInitScript(() => {
		localStorage.theme = 'business';
	});
});

test('simulator runs real WASM battles, preserves separate plans, and replays a batch', async ({
	authenticatedPage: page
}, testInfo) => {
	const errors: string[] = [];
	page.on('pageerror', (e) => errors.push(e.message));
	await page.goto('/simulator');
	const a = page.getByRole('region', { name: 'Side A', exact: true });
	await expect(a.getByTestId('simulation-stack')).toHaveCount(1);
	await a.getByRole('button', { name: 'Duplicate stack' }).click();
	const second = a.getByTestId('simulation-stack').nth(1);
	await second.locator('summary').filter({ hasText: 'Battle plan' }).click();
	await second
		.getByRole('combobox', { name: 'Tactic', exact: true })
		.selectOption({ label: 'Disengage' });
	await a.getByLabel('Stack 1 quantity', { exact: true }).fill('12');
	await page.getByLabel('Simulation seed').fill('42');
	await page.getByLabel('Number of battles').selectOption('10');
	{
		const screenshot = path.resolve(
			'../tmp/battle-simulator',
			testInfo.project.name,
			'desktop-setup.png'
		);
		await page.evaluate(() => window.scrollTo(0, 0));
		await page.screenshot({ path: screenshot, fullPage: true });
		await testInfo.attach('desktop-setup', { path: screenshot, contentType: 'image/png' });
	}
	// A partial number edit must not break autosaving or the app.
	await a.getByLabel('Stack 1 quantity', { exact: true }).fill('');
	await a.getByLabel('Stack 1 quantity', { exact: true }).fill('12');
	await page.reload();
	await expect(a.getByTestId('simulation-stack')).toHaveCount(2);
	await expect(a.getByLabel('Stack 1 quantity', { exact: true })).toHaveValue('12');
	await second.locator('summary').filter({ hasText: 'Battle plan' }).click();
	await expect(second.getByRole('combobox', { name: 'Tactic', exact: true })).toHaveValue('1');
	await page.getByLabel('Simulation seed').fill('42');
	await page.getByLabel('Number of battles').selectOption('10');
	// The simulator requires no game service calls after it has loaded.
	await page.route('**/api/grpc/**', (route) => route.abort());
	await page.getByRole('button', { name: 'Simulate', exact: true }).click();
	await expect(
		page.getByTestId('simulation-summary').getByRole('heading', { name: '10 battles' })
	).toBeVisible();
	await expect(page.getByTestId('battle-board')).toBeVisible();
	const firstSummary = await page.getByTestId('simulation-summary').innerText();
	await page.getByRole('button', { name: /next attack/i }).click();
	await expect(page.getByRole('spinbutton', { name: 'phase', exact: true })).not.toHaveValue('0');
	await page.getByLabel('Battle replay').selectOption('1');
	await expect(page.getByRole('spinbutton', { name: 'phase', exact: true })).toHaveValue('0');
	{
		const screenshot = path.resolve(
			'../tmp/battle-simulator',
			testInfo.project.name,
			'desktop-results.png'
		);
		await page.evaluate(() => window.scrollTo(0, 0));
		await page.screenshot({ path: screenshot, fullPage: true });
		await testInfo.attach('desktop-results', { path: screenshot, contentType: 'image/png' });
	}
	await page.getByRole('button', { name: 'Run again', exact: true }).click();
	await expect(page.getByRole('button', { name: 'Run again', exact: true })).toBeEnabled();
	expect(await page.getByTestId('simulation-summary').innerText()).toEqual(firstSummary);
	expect(errors).toEqual([]);
});

test('simulator edits local designs and validates them without changing the original stack', async ({
	authenticatedPage: page
}) => {
	const errors: string[] = [];
	page.on('pageerror', (e) => errors.push(e.message));
	await page.goto('/simulator');
	const a = page.getByRole('region', { name: 'Side A', exact: true });
	await a.getByRole('button', { name: 'Edit design' }).click();
	const editor = page.getByRole('dialog', { name: 'Simulation ship designer' });
	await editor.getByRole('textbox', { name: 'Name', exact: true }).fill('Simulation prototype');
	await editor.getByRole('button', { name: 'Save simulation design' }).click();
	await expect(editor).toBeHidden();
	await expect(a.getByText('Simulation prototype', { exact: true })).toBeVisible();
	await a.getByRole('button', { name: 'Edit design' }).click();
	await editor.getByRole('textbox', { name: 'Name', exact: true }).fill('Discard this');
	await editor.getByRole('button', { name: 'Cancel', exact: true }).click();
	await expect(a.getByText('Simulation prototype', { exact: true })).toBeVisible();
	await a.getByRole('button', { name: 'New design' }).click();
	await editor.getByLabel('Simulation hull').selectOption('Scout');
	await editor.getByRole('button', { name: 'Design ship' }).click();
	await editor.getByRole('button', { name: 'Save simulation design' }).click();
	await expect(editor.getByText(/required/).last()).toBeVisible();
	await expect(editor).toBeVisible();
	expect(errors).toEqual([]);
});

test('game simulator imports stacks with current tech and leaves game data untouched', async ({
	testGamePage
}, testInfo) => {
	const { page, gameId, player, universe } = await testGamePage('Battle Mixed Fleet');
	const beforePlayer = toJson(PlayerSchema, await getPlayer(page, gameId));
	const beforeUniverse = toJson(PlayerUniverseSchema, await getUniverse(page, gameId));
	const fleet = universe.fleets.find((f) => f.mapObject?.playerNum === player.num)!;
	await page.goto(`/games/${gameId}/simulator?fleet=${fleet.mapObject?.num}`);
	const a = page.getByRole('region', { name: 'Side A', exact: true });
	await expect(a.getByTestId('simulation-stack')).toHaveCount(fleet.tokens.length);
	await a.locator('summary').filter({ hasText: 'Race & technology' }).click();
	await expect(a.getByLabel('A weapons', { exact: true })).toHaveValue(
		String(player.techLevels?.weapons ?? 0)
	);
	await a.getByLabel('Stack 1 quantity', { exact: true }).fill('3');
	await page.getByRole('button', { name: 'Simulate', exact: true }).click();
	await expect(page.getByTestId('battle-board')).toBeVisible();
	{
		const screenshot = path.resolve(
			'../tmp/battle-simulator',
			testInfo.project.name,
			'game-simulator-results.png'
		);
		await page.evaluate(() => window.scrollTo(0, 0));
		await page.screenshot({ path: screenshot, fullPage: true });
		await testInfo.attach('game-simulator-results', { path: screenshot, contentType: 'image/png' });
	}
	expect(toJson(PlayerSchema, await getPlayer(page, gameId))).toEqual(beforePlayer);
	expect(toJson(PlayerUniverseSchema, await getUniverse(page, gameId))).toEqual(beforeUniverse);
});
