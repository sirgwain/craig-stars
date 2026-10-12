import type { Page } from '@playwright/test';
import { clickMapObject, scannerOverlay } from './helpers/scanner';
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
		await gamePage.summary().getByRole('button', { name: 'Add notes and tags' }).click();
		await expect(noteDialog(page).getByRole('heading')).toHaveText(
			`Notes: ${planet?.mapObject?.name}`
		);
		await noteDialog(page).locator('textarea').fill('invade this next');
		await noteDialog(page).getByRole('checkbox', { name: 'Show note on map' }).check();
		await noteDialog(page).getByRole('button', { name: 'Save' }).click();

		await expect(noteDialog(page).locator('textarea')).toBeHidden();
		await expect(note).toHaveText('invade this next');
	});

	await test.step('the note is still there after reloading', async () => {
		await page.reload();
		await gamePage.selectMapObject(planet);
		await expect(note).toHaveText('invade this next');
	});

	await test.step('clearing the note deletes it', async () => {
		await gamePage.summary().getByRole('button', { name: 'Edit notes and tags' }).click();
		await expect(noteDialog(page).locator('textarea')).toHaveValue('invade this next');
		await noteDialog(page).locator('textarea').fill('');
		await noteDialog(page).getByRole('button', { name: 'Save' }).click();

		await expect(note).toBeHidden();
		await expect(
			gamePage.summary().getByRole('button', { name: 'Add notes and tags' })
		).toBeVisible();
	});
});

test('labels the map', async ({ testGamePage }) => {
	const { page } = await testGamePage('Kitchen Sink');
	const position = { x: 120, y: 90 };

	await test.step('add a label', async () => {
		await clickForLabel(page, position);
		await expect(noteDialog(page).getByRole('heading')).toHaveText('New Map Label');
		await noteDialog(page).locator('textarea').fill('the front');
		await noteDialog(page).getByRole('button', { name: 'Save' }).click();
		await expect(noteDialog(page).locator('textarea')).toBeHidden();
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
		await expect(noteDialog(page).locator('textarea')).toHaveValue('the front');
		await noteDialog(page).getByRole('button', { name: 'Delete' }).click();
		await expect(noteDialog(page).locator('textarea')).toBeHidden();
	});

	await test.step('the label is gone', async () => {
		await clickForLabel(page, position);
		await expect(noteDialog(page).getByRole('heading')).toHaveText('New Map Label');
	});
});

test('add label on a planet writes a note about the planet', async ({ testGamePage }) => {
	const { page, gamePage, universe } = await testGamePage('Kitchen Sink');
	const planet = universe.planets.find((p) => p.mapObject?.playerNum === 2);

	await page.locator('[data-type="add-label"]:visible').click();
	await clickMapObject(page, planet);
	await expect(noteDialog(page).getByRole('heading')).toHaveText(
		`Notes: ${planet?.mapObject?.name}`
	);
	// they asked for a note on the map
	await expect(noteDialog(page).getByRole('checkbox', { name: 'Show note on map' })).toBeChecked();
	await noteDialog(page).locator('textarea').fill('their homeworld');
	await noteDialog(page).getByRole('button', { name: 'Save' }).click();
	await expect(noteDialog(page).locator('textarea')).toBeHidden();

	await test.step('the note belongs to the planet', async () => {
		await gamePage.selectMapObject(planet);
		await expect(gamePage.summary().locator('[data-type="map-object-note"]')).toHaveText(
			'their homeworld'
		);
	});

	await test.step('clicking the planet again edits the same note', async () => {
		await page.locator('[data-type="add-label"]:visible').click();
		await clickMapObject(page, planet);
		await expect(noteDialog(page).locator('textarea')).toHaveValue('their homeworld');
	});
});

test('tags planets and fleets', async ({ testGamePage }) => {
	const { page, gamePage, universe } = await testGamePage('Kitchen Sink');
	const planet = universe.planets.find((p) => p.mapObject?.playerNum === 2);
	// a fleet of ours that's alone in deep space
	const fleet = universe.fleets.find(
		(f) =>
			f.mapObject?.playerNum === 1 &&
			f.mapObject.position?.x === 0 &&
			f.mapObject.position?.y === 50
	);
	const tags = gamePage.summary().locator('[data-type="map-object-tags"] [data-type="tag"]');
	const toggle = (name: string) => noteDialog(page).getByRole('button', { name, exact: true });

	await test.step('create tags on a planet', async () => {
		await gamePage.selectMapObject(planet);
		await gamePage.summary().getByRole('button', { name: 'Add notes and tags' }).click();
		const search = noteDialog(page).getByRole('textbox', { name: 'Find or create a tag' });
		await search.fill('target');
		await search.press('Enter');
		await search.fill('gate hub');
		await noteDialog(page).getByRole('button', { name: 'Create “gate hub”' }).click();
		await expect(toggle('target')).toHaveAttribute('aria-pressed', 'true');
		await expect(toggle('gate hub')).toHaveAttribute('aria-pressed', 'true');
		await noteDialog(page).getByRole('button', { name: 'Save' }).click();

		await expect(tags).toHaveText(['target', 'gate hub']);
	});

	await test.step('toggle an existing tag on a fleet', async () => {
		await page.reload();
		await gamePage.selectMapObject(fleet);
		await gamePage.summary().getByRole('button', { name: 'Add notes and tags' }).click();
		await expect(noteDialog(page).getByRole('heading')).toHaveText(
			`Notes: ${fleet?.mapObject?.name}`
		);
		// fleets move, so their notes aren't shown on the map
		await expect(noteDialog(page).getByRole('checkbox')).toBeHidden();
		await expect(toggle('target')).toHaveAttribute('aria-pressed', 'false');
		await toggle('target').click();
		await noteDialog(page).locator('textarea').fill('decoy');
		await noteDialog(page).getByRole('button', { name: 'Save' }).click();

		await expect(tags).toHaveText(['target']);
		await expect(gamePage.summary().locator('[data-type="map-object-note"]')).toHaveText('decoy');
	});

	// a fleet in transit can't be selected on the map, so its notes are also in the command pane
	await test.step('edit the fleet note from the composition tile', async () => {
		const tile = page.locator('[data-type="command-tile"][data-id="Fleet Composition"]:visible');
		// clicking the selected fleet commands it
		await gamePage.selectMapObject(fleet);
		await expect(tile.locator('[data-type="fleet-note"]')).toHaveText('decoy');
		await expect(tile.locator('[data-type="tag"]')).toHaveText(['target']);
		await tile.getByRole('button', { name: 'Edit notes and tags' }).click();
		await noteDialog(page).locator('textarea').fill('bait');
		await noteDialog(page).getByRole('button', { name: 'Save' }).click();
		await expect(tile.locator('[data-type="fleet-note"]')).toHaveText('bait');
	});

	await test.step('the pencil beside the fleet name opens its notes too', async () => {
		await page.locator('[data-type="fleet-title-note-button"]:visible').click();
		await expect(noteDialog(page).locator('textarea')).toHaveValue('bait');
		await noteDialog(page).getByRole('button', { name: 'Cancel' }).click();
	});

	await test.step('the notes filter lists the tags', async () => {
		const filter = page.locator('[data-filter="notes"]:visible');
		// the filters are in a menu when the toolbar is too narrow for them
		const openFilters = () => page.locator('a[href="#planet-view-states"]').click();
		await openFilters();
		await filter.click();
		const menu = page.locator('.dropdown:has([data-filter="notes"]:visible) ul');
		await expect(menu.getByRole('menuitemcheckbox', { name: 'gate hub' })).toBeChecked();
		await menu.getByRole('menuitemcheckbox', { name: 'No Notes' }).click();
		await expect(filter).not.toHaveClass(/btn-primary/);
		// the filter is remembered
		await page.reload();
		await openFilters();
		await expect(filter).not.toHaveClass(/btn-primary/);
	});

	await test.step('removing the last tag removes it from the planet', async () => {
		await gamePage.selectMapObject(planet);
		await expect(tags).toHaveText(['target', 'gate hub']);
		await gamePage.summary().getByRole('button', { name: 'Edit notes and tags' }).click();
		await toggle('target').click();
		await toggle('gate hub').click();
		await noteDialog(page).getByRole('button', { name: 'Save' }).click();

		await expect(tags).toHaveCount(0);
		await expect(
			gamePage.summary().getByRole('button', { name: 'Add notes and tags' })
		).toBeVisible();
	});
});

test('notes report', async ({ testGamePage }) => {
	const { page, gamePage, gameId, universe } = await testGamePage('Kitchen Sink');
	const planet = universe.planets.find((p) => p.mapObject?.playerNum === 2);
	const report = page.locator('[data-type="notes-report"]');
	const notes = report.locator('[data-type="note-text"]');

	await test.step('write a note with a tag, and a label', async () => {
		await gamePage.selectMapObject(planet);
		await gamePage.summary().getByRole('button', { name: 'Add notes and tags' }).click();
		await noteDialog(page).locator('textarea').fill('invade this next');
		const search = noteDialog(page).getByRole('textbox', { name: 'Find or create a tag' });
		await search.fill('target');
		await search.press('Enter');
		await noteDialog(page).getByRole('button', { name: 'Save' }).click();
		await expect(noteDialog(page).locator('textarea')).toBeHidden();

		await clickForLabel(page, { x: 120, y: 90 });
		await noteDialog(page).locator('textarea').fill('the front');
		await noteDialog(page).getByRole('button', { name: 'Save' }).click();
		await expect(noteDialog(page).locator('textarea')).toBeHidden();
	});

	await test.step('the report lists both', async () => {
		await page.locator('label').filter({ hasText: 'Reports' }).click();
		await page.getByRole('link', { name: 'Notes', exact: true }).first().click();
		await expect(page).toHaveURL(new RegExp(`/games/${gameId}/notes$`));
		await expect(notes).toHaveText(['the front', 'invade this next']);
	});

	await test.step('filter by tag, type and text', async () => {
		await report.getByRole('button', { name: 'Filter by target' }).click();
		await expect(notes).toHaveText(['invade this next']);
		await report.getByRole('button', { name: 'Filter by target' }).click();
		await expect(notes).toHaveCount(2);

		await report.getByLabel('Type of note').selectOption('Label');
		await expect(notes).toHaveText(['the front']);
		await report.getByLabel('Type of note').selectOption('');

		await report.getByPlaceholder('search').fill('invade');
		await expect(notes).toHaveText(['invade this next']);
		await report.getByPlaceholder('search').fill('');
	});

	await test.step('edit a note from the report', async () => {
		await report.getByRole('button', { name: 'Edit' }).last().click();
		await expect(noteDialog(page).locator('textarea')).toHaveValue('invade this next');
		await noteDialog(page).locator('textarea').fill('invaded');
		await noteDialog(page).getByRole('button', { name: 'Save' }).click();
		await expect(notes).toHaveText(['the front', 'invaded']);
	});

	await test.step('delete a note from the report', async () => {
		page.once('dialog', (dialog) => dialog.accept());
		await report.getByRole('button', { name: 'Delete' }).first().click();
		await expect(notes).toHaveText(['invaded']);
	});

	await test.step('go to the planet', async () => {
		await report.getByRole('button', { name: planet?.mapObject?.name }).click();
		await expect(page).toHaveURL(new RegExp(`/games/${gameId}$`));
		await expect(gamePage.summary().locator('[data-type="map-object-note"]')).toHaveText('invaded');
	});
});
