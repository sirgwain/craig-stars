import type { Page } from '@playwright/test';
import { scannerOverlay } from './helpers/scanner';
import { expect, test } from './setup';

function noteDialog(page: Page) {
	return page.locator('[data-type="note-dialog"]');
}

// start add label mode and click the scanner
async function clickForLabel(page: Page, position: { x: number; y: number }) {
	await page.locator('[data-type="add-label"]:visible').click();
	await scannerOverlay(page).click({ position, force: true });
}

test('writes a note about a planet', async ({ testGamePage }) => {
	const { page, gamePage, universe } = await testGamePage('Kitchen Sink');
	const planet = universe.planets.find((p) => p.mapObject?.playerNum === 2);
	const note = gamePage.summary().locator('[data-type="map-object-note"]');

	await test.step('add a note', async () => {
		await gamePage.selectMapObject(planet);
		await gamePage.summary().getByRole('button', { name: 'Add note' }).click();
		await expect(noteDialog(page).getByRole('heading')).toHaveText(
			`Note: ${planet?.mapObject?.name}`
		);
		await noteDialog(page).getByRole('textbox').fill('invade this next');
		await noteDialog(page).getByRole('checkbox', { name: 'Show on map' }).check();
		await noteDialog(page).getByRole('button', { name: 'Save' }).click();

		await expect(noteDialog(page).getByRole('textbox')).toBeHidden();
		await expect(note).toHaveText('invade this next');
	});

	await test.step('the note is still there after reloading', async () => {
		await page.reload();
		await gamePage.selectMapObject(planet);
		await expect(note).toHaveText('invade this next');
	});

	await test.step('clearing the note deletes it', async () => {
		await gamePage.summary().getByRole('button', { name: 'Edit note' }).click();
		await expect(noteDialog(page).getByRole('textbox')).toHaveValue('invade this next');
		await noteDialog(page).getByRole('textbox').fill('');
		await noteDialog(page).getByRole('button', { name: 'Save' }).click();

		await expect(note).toBeHidden();
		await expect(gamePage.summary().getByRole('button', { name: 'Add note' })).toBeVisible();
	});
});

test('labels the map', async ({ testGamePage }) => {
	const { page } = await testGamePage('Kitchen Sink');
	const position = { x: 120, y: 90 };

	await test.step('add a label', async () => {
		await clickForLabel(page, position);
		await expect(noteDialog(page).getByRole('heading')).toHaveText('New Map Label');
		await noteDialog(page).getByRole('textbox').fill('the front');
		await noteDialog(page).getByRole('button', { name: 'Save' }).click();
		await expect(noteDialog(page).getByRole('textbox')).toBeHidden();
	});

	await test.step('the map selects things again after the label is placed', async () => {
		await expect(page.locator('[data-type="add-label"]:visible svg')).not.toHaveClass(
			/stroke-accent/
		);
	});

	await test.step('clicking the label in add label mode edits it', async () => {
		await page.reload();
		await clickForLabel(page, { x: position.x + 5, y: position.y + 5 });
		await expect(noteDialog(page).getByRole('heading')).toHaveText('Map Label');
		await expect(noteDialog(page).getByRole('textbox')).toHaveValue('the front');
		await noteDialog(page).getByRole('button', { name: 'Delete' }).click();
		await expect(noteDialog(page).getByRole('textbox')).toBeHidden();
	});

	await test.step('the label is gone', async () => {
		await clickForLabel(page, position);
		await expect(noteDialog(page).getByRole('heading')).toHaveText('New Map Label');
	});
});
