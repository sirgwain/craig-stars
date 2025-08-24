import { apiErrorsFailTest, expect, submitTurn, test } from './setup';
import { WaypointTaskTransportAction } from '../src/lib/protogen/craig_stars/v1/fleet_pb';

test('create a new game', async ({ newGamePage }) => {
	const { page, name } = newGamePage;
	const gameLink = await page.locator('[data-type="game-link"]').first();
	await expect(gameLink).toBeVisible();
	await expect(gameLink).toHaveText(`${name} - 2400`);
});

test('submit turn', async ({ newGamePage }) => {
	const { page, id, name } = newGamePage;
	apiErrorsFailTest(page, id);

	// start with a new game, ensure we have year 2400
	const gameLink = await page.locator('[data-type="game-link"]').first();
	await expect(gameLink).toBeVisible();
	await expect(gameLink).toHaveText(`${name} - 2400`);

	await page.getByRole('button', { name: 'Submit Turn' }).click();

	await expect(gameLink).toBeVisible();
	await expect(gameLink).toHaveText(`${name} - 2401`);
});

test('research page', async ({ newGamePage }) => {
	const { page, id } = newGamePage;
	apiErrorsFailTest(page, id);

	await page.locator('label').filter({ hasText: 'Commands' }).click();
	await page.getByRole('link', { name: 'Research' }).click();
	await page.getByRole('radio', { name: 'Weapons' }).click();
	// select weapons
	await expect(page.getByText('Weapons 4')).toBeVisible();

	// select prop
	await page.getByRole('radio', { name: 'Propulsion' }).click();
	await expect(page.getByText('Propulsion 4')).toBeVisible();

	// select con
	await page.getByRole('radio', { name: 'Construction' }).click();
	await expect(page.getByText('Construction 4')).toBeVisible();

	// select elec
	await page.getByRole('radio', { name: 'Electronics' }).click();
	await expect(page.getByText('Electronics 4')).toBeVisible();

	// select energy
	await page.getByRole('radio', { name: 'Energy' }).click();
	await expect(page.getByText('Energy 4')).toBeVisible();

	// increase research
	await page.locator('[data-type="spin-number-increase-button"]').click();
	await expect(page.getByText('Research Budget 16 %')).toBeVisible();

	// decrease research
	await page.locator('[data-type="spin-number-decrease-button"]').click();
	await expect(page.getByText('Research Budget 15 %')).toBeVisible();
});

test('relations page', async ({ newGamePage }) => {
	const { page, id } = newGamePage;
	apiErrorsFailTest(page, id);

	await page.locator('label').filter({ hasText: 'Commands' }).click();
	await page.getByRole('link', { name: 'Relations' }).click();
	await page.waitForURL(`/games/${id}/relations`);

	await page.getByRole('radio', { name: 'Friend' }).first().check();
	await page.getByRole('radio', { name: 'Neutral' }).first().check();
	await page.getByRole('radio', { name: 'Enemy' }).first().check();
	await page.locator('input[name="player-relation-2-share-map"]').check();
});

test('battle plans page', async ({ newGamePage }) => {
	const { page, id } = newGamePage;
	const name = 'Test Battle Plan';
	apiErrorsFailTest(page, id);

	await page.locator('label').filter({ hasText: 'Commands' }).click();
	await page.getByRole('link', { name: 'Battle Plans' }).click();
	await page.waitForURL(`/games/${id}/battle-plans`);

	await page.getByRole('link', { name: 'Create' }).click();
	await page.getByRole('textbox', { name: 'Name' }).fill(name);
	await page.getByLabel('Primary Target').selectOption('Any');
	await page.getByLabel('Secondary Target').selectOption('Starbase');
	await page.getByLabel('Tactic DisengageDisengage If').selectOption('Disengage');
	await page.getByLabel('Attack Who EnemiesEnemies And').selectOption('Everyone');
	await page.getByRole('button', { name: 'Save' }).click();
	await page.getByRole('link', { name: 'Battle Plans' }).click();

	// delete the plan we just created
	const deleteButton = await page.locator(`[data-type="delete-button"][data-id="${name}"]`);
	await expect(deleteButton).toBeVisible();

	// confirm delete
	page.once('dialog', async (dialog) => {
		await dialog.accept();
	});
	await deleteButton.click();
	await expect(page.getByRole('link', { name: name })).not.toBeVisible();
});

test('production plans page', async ({ newGamePage }) => {
	const { page, id } = newGamePage;
	const name = 'Test Production Plan';
	apiErrorsFailTest(page, id);

	await page.locator('label').filter({ hasText: 'Commands' }).click();
	await page.getByRole('link', { name: 'Production Plans' }).click();
	await page.waitForURL(`/games/${id}/production-plans`);

	await page.getByRole('link', { name: 'Create' }).click();
	await page.getByRole('textbox', { name: 'Name' }).fill(name);

	await page.getByRole('button', { name: 'Factory (Auto)' }).click();
	await page.getByRole('button', { name: 'Add' }).click();
	await page.getByRole('button', { name: 'Mine (Auto)' }).click();
	await page.getByRole('button', { name: 'Add' }).click();
	await page.getByRole('button', { name: 'Defenses (Auto)' }).click();
	await page.getByRole('button', { name: 'Add' }).click();
	await page.getByRole('button', { name: 'Remove' }).click();
	await page.getByRole('button', { name: 'Item Up' }).click();
	await page.getByRole('button', { name: 'Item Down' }).click();
	await page.getByRole('button', { name: 'Clear' }).click();
	await page.getByRole('button', { name: 'Alchemy (Auto)' }).click();
	await page.getByRole('button', { name: 'Add' }).click();
	await page.getByText('Contributes Only Leftover to').click();

	await page.getByRole('button', { name: 'Save' }).click();
	await page.getByRole('link', { name: 'Production Plans' }).click();

	// delete the plan we just created
	const deleteButton = await page.locator(`[data-type="delete-button"][data-id="${name}"]`);
	await expect(deleteButton).toBeVisible();

	// confirm delete
	page.once('dialog', async (dialog) => {
		await dialog.accept();
	});
	await deleteButton.click();
	await expect(page.getByRole('link', { name: name })).not.toBeVisible();
});

test('transport plans page', async ({ newGamePage }) => {
	const { page, id } = newGamePage;
	const name = 'Test Transport Plan';
	apiErrorsFailTest(page, id);

	await page.locator('label').filter({ hasText: 'Commands' }).click();
	await page.getByRole('link', { name: 'Transport Plans' }).click();
	await page.waitForURL(`/games/${id}/transport-plans`);

	await page.getByRole('link', { name: 'Create' }).click();
	await page.getByRole('textbox', { name: 'Name' }).fill(name);

	await page
		.getByLabel('Action Fuel NoneLoad')
		.selectOption(String(WaypointTaskTransportAction.LOAD_OPTIMAL));

	await page
		.getByLabel('Action Ironium NoneLoad')
		.selectOption(String(WaypointTaskTransportAction.LOAD_AMOUNT));
	await page.getByLabel('Amount Ironium').fill('1');
	await page
		.getByLabel('Action Boranium NoneLoad')
		.selectOption(String(WaypointTaskTransportAction.UNLOAD_AMOUNT));
	await page.getByLabel('Amount Boranium').fill('1');
	await page
		.getByLabel('Action Germanium NoneLoad')
		.selectOption(String(WaypointTaskTransportAction.FILL_PERCENT));
	await page.getByLabel('Amount Germanium').fill('3');
	await page
		.getByLabel('Action Colonists NoneLoad')
		.selectOption(String(WaypointTaskTransportAction.SET_AMOUNT_TO));
	await page.getByLabel('Amount Colonists').fill('10');

	await page.getByRole('button', { name: 'Save' }).click();
	await page.getByRole('link', { name: 'Transport Plans' }).click();

	// delete the plan we just created
	const deleteButton = await page.locator(`[data-type="delete-button"][data-id="${name}"]`);
	await expect(deleteButton).toBeVisible();

	// confirm delete
	page.once('dialog', async (dialog) => {
		await dialog.accept();
	});
	await deleteButton.click();
	await expect(page.getByRole('link', { name: name })).not.toBeVisible();
});

test('planets report page', async ({ newGamePage }) => {
	const { page, id } = newGamePage;

	await page.locator('label').filter({ hasText: 'Reports' }).click();
	await page.getByRole('link', { name: 'Planets' }).click();
	await page.waitForURL(`/games/${id}/planets`);

	// sort
	await page.getByRole('button', { name: 'Name' }).click();
	await page.getByRole('button', { name: 'Starbase' }).click();
	await page.getByRole('button', { name: 'Population' }).click();
	await page.getByRole('button', { name: 'Cap' }).click();
	await page.getByRole('button', { name: 'Growth' }).click();
	await page.getByRole('button', { name: 'Value' }).click();
	await page.getByRole('button', { name: 'Production' }).click();
	await page.getByRole('button', { name: 'Mine', exact: true }).click();
	await page.getByRole('button', { name: 'Factories' }).click();
	await page.getByRole('columnheader', { name: 'Defense' }).click();
	await page.getByRole('columnheader', { name: 'Minerals' }).click();
	await page.getByRole('button', { name: 'Mining Rate' }).click();

	// show all
	await page.getByRole('checkbox', { name: 'Show All' }).check();

	// sort
	await page.getByRole('button', { name: 'Name' }).click();
	await page.getByRole('button', { name: 'Owner' }).click();
	await page.getByRole('button', { name: 'Report Age' }).click();
	await page.getByRole('button', { name: 'Starbase' }).click();
	await page.getByRole('button', { name: 'Population' }).click();
	await page.getByRole('button', { name: 'Value' }).click();
	await page.getByRole('button', { name: 'Defense' }).click();
	await page.getByRole('button', { name: 'Minerals' }).click();
	await page.getByRole('button', { name: 'Mineral Concentration' }).click();

	// show just ours
	await page.getByRole('checkbox', { name: 'Show All' }).check();
});

test('fleets report page', async ({ newGamePage }) => {
	const { page, id } = newGamePage;

	await page.locator('label').filter({ hasText: 'Reports' }).click();
	await page.getByRole('link', { name: 'Fleets', exact: true }).click();
	await page.waitForURL(`/games/${id}/fleets`);

	// sort
	await page.getByRole('button', { name: 'Name' }).click();
	await page.getByRole('button', { name: 'ID' }).click();
	await page.getByRole('button', { name: 'Location' }).click();
	await page.getByRole('button', { name: 'Destination' }).click();
	await page.getByRole('button', { name: 'ETA' }).click();
	await page.getByRole('button', { name: 'Fuel' }).click();
	await page.getByRole('button', { name: 'Cargo' }).click();
	await page.getByText('Composition').click();
	await page.getByRole('button', { name: 'Cloak' }).click();
	await page.getByRole('button', { name: 'Battle Plan' }).click();
	await page.getByRole('button', { name: 'Mass' }).click();
});

test('designs report page', async ({ newGamePage }) => {
	const { page, id } = newGamePage;

	await page.locator('label').filter({ hasText: 'Reports' }).click();
	await page.getByRole('link', { name: 'Designs' }).click();
	await page.waitForURL(`/games/${id}/designs`);

	// sort
	await page.getByRole('button', { name: 'Player' }).click();
	await page.getByRole('button', { name: 'ID', exact: true }).first().click();
	await page.getByRole('button', { name: 'Name' }).click();
	await page.getByRole('button', { name: 'Hull' }).click();
	await page.getByRole('button', { name: 'Rating' }).click();
	await page.getByRole('button', { name: 'Armor' }).click();
	await page.getByRole('button', { name: 'Shields' }).click();
	await page.getByRole('button', { name: 'Initiative' }).click();
	await page.getByRole('button', { name: 'Movement' }).click();
	await page.getByRole('button', { name: 'Mass' }).click();

	// popup
	await page.getByRole('button', { name: 'Scout' }).first().click();

	// filter by player

	await page
		.getByRole('row', { name: 'Humanoids 1 Opens ship design' })
		.getByRole('link')
		.first()
		.click();
	await page.waitForURL(`/games/${id}/designs/1`);
});

test('messages report page', async ({ newGamePage }) => {
	const { page, id } = newGamePage;

	await page.locator('label').filter({ hasText: 'Reports' }).click();
	await page.getByRole('link', { name: 'Messages' }).click();
	await page.waitForURL(`/games/${id}/messages`);

	// sort

	await page.getByRole('button', { name: 'Target' }).click();
	await page.getByRole('columnheader', { name: 'Target' }).click();
	await page.getByRole('button', { name: 'Text' }).click();
	await page.getByRole('button', { name: 'Text' }).click();
	await page.getByRole('checkbox', { name: 'Show All Messages' }).check();
	await page.getByRole('checkbox', { name: 'Show All Messages' }).uncheck();
	await page.getByRole('textbox', { name: 'search' }).fill('home planet');

	// go to home world
	await page.locator('[data-type="goto-target-button"]').first().click();
	await page.waitForURL(`/games/${id}`);
});

test('battles report page', async ({ newGamePage }) => {
	const { page, id } = newGamePage;

	await page.locator('label').filter({ hasText: 'Reports' }).click();
	await page.getByRole('link', { name: 'Battles' }).click();
	await page.waitForURL(`/games/${id}/battles`);

	// sort
	await page.getByRole('button', { name: 'Location' }).click();
	await page.getByRole('button', { name: 'Location' }).click();
	await page.getByRole('button', { name: 'Players' }).click();
	await page.getByRole('button', { name: 'Players' }).click();
	await page.getByRole('button', { name: 'Ships' }).click();
	await page.getByRole('button', { name: 'Ships' }).click();
	await page.getByRole('button', { name: 'Ours', exact: true }).click();
	await page.getByRole('button', { name: 'Ours', exact: true }).click();
	await page.getByRole('button', { name: 'Theirs', exact: true }).click();
	await page.getByRole('button', { name: 'Theirs', exact: true }).click();
	await page.getByRole('button', { name: 'Our Dead' }).click();
	await page.getByRole('button', { name: 'Our Dead' }).click();
	await page.getByRole('button', { name: 'Their Dead' }).click();
	await page.getByRole('button', { name: 'Their Dead' }).click();
	await page.getByRole('button', { name: 'Ours Left' }).click();
	await page.getByRole('button', { name: 'Ours Left' }).click();
	await page.getByRole('button', { name: 'Theirs Left' }).click();
	await page.getByRole('button', { name: 'Theirs Left' }).click();
});

test('players report page', async ({ newGamePage }) => {
	const { page, id } = newGamePage;

	// generate some data
	await submitTurn(page);

	await page.locator('label').filter({ hasText: 'Reports' }).click();
	await page.getByRole('link', { name: 'Players' }).click();
	await page.waitForURL(`/games/${id}/players`);

	await page.getByText('1 admin').click();
	await expect(page.getByRole('main')).toContainText('1 admin');
	await expect(page.getByRole('main')).toContainText('Humanoids');
	await expect(page.getByRole('main')).toContainText('Playing');
});

test('game techs page', async ({ newGamePage }) => {
	const { page, id } = newGamePage;

	// browse to player techs page, show all, search for scout, click scout
	await page.locator('#menu svg').click();
	await page.getByRole('link', { name: 'Techs' }).click();
	await page.waitForURL(`/games/${id}/techs`);

	await page.getByRole('checkbox', { name: 'Show All' }).check();
	await page.getByRole('checkbox', { name: 'Show All' }).uncheck();
	await page.getByRole('textbox', { name: 'search' }).fill('scout');
	await page.getByRole('link', { name: 'Scout' }).click();

	await page.waitForURL(`/games/${id}/techs/scout`);
});

test('game race page', async ({ newGamePage }) => {
	const { page, id } = newGamePage;

	// browse to player techs page, show all, search for scout, click scout
	await page.locator('#menu svg').click();
	await page.getByRole('link', { name: 'Race' }).click();
	await page.waitForURL(`/games/${id}/race`);
});
