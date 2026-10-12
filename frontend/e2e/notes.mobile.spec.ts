import { expect, type Page } from '@playwright/test';
import { scannerOverlay, tapMapObject } from './helpers/scanner';
import { test } from './setup';

function noteDialog(page: Page) {
	return page.locator('[data-type="note-dialog"]');
}

// On a phone there's no room for add label in the toolbar, it's in the menu of collapsed items
test('labels the map from the toolbar menu on a phone', async ({ testGamePage }) => {
	const { page } = await testGamePage('Kitchen Sink');
	const addLabel = page.getByRole('link', { name: 'Add map label' });

	await page.locator('a[href="#planet-view-states"]').tap();
	await addLabel.tap();
	// the menu covers the map until it has finished closing
	await expect(addLabel).toBeHidden();
	await scannerOverlay(page).tap({ position: { x: 200, y: 300 }, force: true });

	await expect(noteDialog(page).getByRole('heading')).toHaveText('New Map Label');
	// the dialog is as wide as the screen once it has finished opening
	const width = page.viewportSize()?.width ?? 0;
	await expect
		.poll(async () => (await noteDialog(page).locator('.modal-box').boundingBox())?.width)
		.toBeCloseTo(width, -1);

	await noteDialog(page).locator('textarea').fill('the front');
	await noteDialog(page).getByRole('button', { name: 'Save' }).tap();
	await expect(noteDialog(page).locator('textarea')).toBeHidden();
});

test('writes a note about a planet from the drawer on a phone', async ({ testGamePage }) => {
	const { page, universe } = await testGamePage('Kitchen Sink');
	const drawer = page.locator('[data-type="command-drawer"]');
	const planet = universe.planets.find((p) => p.mapObject?.playerNum === 2);

	await tapMapObject(page, planet);
	await drawer.getByRole('button', { name: 'Add notes and tags' }).tap();
	await expect(noteDialog(page).getByRole('heading')).toHaveText(
		`Notes: ${planet?.mapObject?.name}`
	);
	// the dialog stays in the top of the screen, the on screen keyboard covers the rest
	const height = page.viewportSize()?.height ?? 0;
	const box = () => noteDialog(page).locator('.modal-box').boundingBox();
	// it slides in from above the screen
	await expect.poll(async () => (await box())?.y).toBeCloseTo(0, 0);
	const opened = await box();
	expect((opened?.y ?? 0) + (opened?.height ?? height)).toBeLessThan(height * 0.55);
	await noteDialog(page).locator('textarea').fill('invade this next');
	await noteDialog(page).getByRole('button', { name: 'Save' }).tap();

	await expect(drawer.getByRole('button', { name: 'Edit notes and tags' })).toHaveText(
		'invade this next'
	);
});
