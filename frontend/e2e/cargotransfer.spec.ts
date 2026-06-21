import { MapObjectType } from '../src/lib/protogen/craig_stars/v1/common_pb';
import { expect, submitTurn, test } from './setup';

test('Cargo Transfer Planet Owned', async ({ testGamePage }) => {
	const { page, gamePage } = await testGamePage('Cargo Transfer Planet Owned');

	await gamePage.clickTileButton('Fleets In Orbit', 'Goto');
	await gamePage.openTransferFromTile('Orbiting Planet 1');
	await gamePage.transferCargo({ ironium: 1, boranium: 1, germanium: 1, colonists: 1 }, 'source');
	await gamePage.confirmDialog();

	await gamePage.expectTileCargo('Fuel & Cargo', {
		ironium: '1kT',
		boranium: '1kT',
		germanium: '1kT',
		colonists: '1kT'
	});

	await submitTurn(page);

	await gamePage.expectTileCargo('Minerals on Hand', {
		ironium: '999kT',
		boranium: '999kT',
		germanium: '999kT'
	});
	await expect(gamePage.tile('Status').getByText('Population 287,300').first()).toBeVisible();

	await gamePage.clickTileButton('Fleets In Orbit', 'Goto');
	await gamePage.expectTileCargo('Fuel & Cargo', {
		ironium: '1kT',
		boranium: '1kT',
		germanium: '1kT',
		colonists: '1kT'
	});
});

test('Cargo Transfer Fleets', async ({ testGamePage }) => {
	const { page, gamePage } = await testGamePage('Cargo Transfer Fleets');

	await gamePage.clickTileButton('Other Entities Here', 'Transfer');
	await gamePage.transferCargo(
		{ fuel: 1, ironium: 1, boranium: 1, germanium: 1, colonists: 1 },
		'dest'
	);
	await gamePage.confirmDialog();

	await gamePage.expectTileCargo('Fuel & Cargo', {
		fuel: '99 of 450mg',
		ironium: '9kT',
		boranium: '9kT',
		germanium: '9kT',
		colonists: '9kT'
	});

	await gamePage.clickNextSelectedObject();
	await gamePage.expectTileCargo('Fuel & Cargo', {
		fuel: '11 of 200mg',
		ironium: '6kT',
		boranium: '6kT',
		germanium: '6kT',
		colonists: '6kT'
	});

	await submitTurn(page);

	await gamePage.expectTileCargo('Fuel & Cargo', {
		fuel: '99 of 450mg',
		ironium: '9kT',
		boranium: '9kT',
		germanium: '9kT',
		colonists: '9kT'
	});

	await gamePage.clickNextSelectedObject();
	await gamePage.expectTileCargo('Fuel & Cargo', {
		fuel: '11 of 200mg',
		ironium: '6kT',
		boranium: '6kT',
		germanium: '6kT',
		colonists: '6kT'
	});
});

test('Cargo Transfer Split', async ({ testGamePage }) => {
	const { page, gamePage } = await testGamePage('Cargo Transfer Split');

	await gamePage.openTransferFromTile('In Deep Space');
	await gamePage.transferCargo({ ironium: 10 }, 'dest');
	await gamePage.confirmDialog();

	await page.getByRole('button', { name: 'Split All' }).first().click();

	await gamePage.expectTileCargo('Fuel & Cargo', {
		ironium: '36kT',
		boranium: '45kT',
		germanium: '45kT'
	});

	await gamePage.clickNextSelectedObject();
	await gamePage.expectTileCargo('Fuel & Cargo', {
		ironium: '4kT',
		boranium: '5kT',
		germanium: '5kT'
	});

	await gamePage.openTransferFromTile('In Deep Space');
	await gamePage.transferCargo({ ironium: 10 }, 'source');
	await gamePage.confirmDialog();

	await submitTurn(page);

	await gamePage.expectTileCargo('Fuel & Cargo', {
		ironium: '36kT',
		boranium: '45kT',
		germanium: '45kT'
	});

	await gamePage.clickNextSelectedObject();
	await gamePage.expectTileCargo('Fuel & Cargo', {
		ironium: '14kT',
		boranium: '5kT',
		germanium: '5kT'
	});
});

test('Cargo Transfer Jettison', async ({ testGamePage }) => {
	const { page, gamePage } = await testGamePage('Cargo Transfer Jettison');

	await gamePage.openTransferFromTile('In Deep Space');
	await gamePage.transferCargo({ ironium: 10, boranium: 20, germanium: 30 }, 'dest');
	await gamePage.confirmDialog();

	await gamePage.expectTileCargo('Fuel & Cargo', {
		ironium: '40kT',
		boranium: '30kT',
		germanium: '20kT'
	});

	await submitTurn(page);

	await gamePage.expectTileCargo('Fuel & Cargo', {
		ironium: '40kT',
		boranium: '30kT',
		germanium: '20kT'
	});

	await gamePage.cycleSummarySelection();
	await gamePage.expectSummarySelection(/^Salvage #1$/);
	await gamePage.scrollSummaryText('Salvage #1 Humanoids Location');
	await gamePage.expectSummaryCargo({
		ironium: '0kT',
		boranium: '10kT',
		germanium: '20kT'
	});
});

test('Cargo Transfer Salvage', async ({ testGamePage }) => {
	const { page, gamePage } = await testGamePage('Cargo Transfer Salvage');

	await gamePage.openTransferFromTile('In Deep Space');
	await gamePage.transferCargo({ ironium: 1, boranium: 2, germanium: 3 }, 'dest');
	await gamePage.confirmDialog();

	await gamePage.cycleSummarySelection();
	await gamePage.expectSummarySelection(/^Salvage #1$/);
	await gamePage.scrollSummaryText('Salvage #1 Humanoids Location');
	await gamePage.expectSummaryCargo({
		ironium: '51kT',
		boranium: '52kT',
		germanium: '53kT'
	});

	await submitTurn(page);

	await gamePage.expectTileCargo('Fuel & Cargo', {
		ironium: '9kT',
		boranium: '8kT',
		germanium: '7kT'
	});

	await gamePage.cycleSummarySelection();
	await gamePage.expectSummarySelection(/^Salvage #1$/);
	await gamePage.scrollSummaryText('Salvage #1 Humanoids Location');
	await gamePage.expectSummaryCargo({
		ironium: '41kT',
		boranium: '42kT',
		germanium: '43kT'
	});
});

test('Cargo Transfer MineralPacket', async ({ testGamePage }) => {
	const { page, gamePage } = await testGamePage('Cargo Transfer MineralPacket');

	await gamePage.clickTileButton('Other Entities Here', 'Transfer');
	await gamePage.transferCargo({ ironium: 1, boranium: 2, germanium: 3 }, 'source');
	await gamePage.confirmDialog();

	await gamePage.cycleSummarySelection();
	await gamePage.expectSummarySelection(/^Humanoids Mineral Packet #1$/);
	await gamePage.scrollSummaryText('Location: (0, 0) Traveling at');
	await gamePage.expectSummaryCargo({
		ironium: '49kT',
		boranium: '48kT',
		germanium: '47kT'
	});

	await submitTurn(page);

	await gamePage.expectTileCargo('Fuel & Cargo', {
		ironium: '11kT',
		boranium: '12kT',
		germanium: '13kT'
	});

	await page.locator(`[data-id="${MapObjectType.MINERAL_PACKET}-1-1"]`).click({ force: true });

	await gamePage.expectSummarySelection(/^Humanoids Mineral Packet #1$/);
	await gamePage.scrollSummaryText('Location: (25, 0) Traveling at');
	await gamePage.expectSummaryCargo({
		ironium: '25kT',
		boranium: '24kT',
		germanium: '24kT'
	});
});
