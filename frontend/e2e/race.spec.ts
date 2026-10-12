import { expect, test } from './setup';

test('view races', async ({ authenticatedPage: page }) => {
	await page.getByRole('link', { name: 'Races' }).click();
	await page.waitForURL(`/races`);

	// wait for races to show up
	await expect(page.getByRole('link', { name: 'Humanoids' })).toBeVisible();
});

test('edit race', async ({ newRacePage }) => {
	const { page } = newRacePage;
	// humanoids have 25 points left
	await expect(page.getByText('Points 25')).toBeVisible();

	// click the LRT Improved Fuel Efficiency
	await page.getByRole('button', { name: 'Improved Fuel Efficiency' }).click();

	// negative points
	await expect(page.getByText('Points -53')).toBeVisible();

	// grayed out save button
	await expect(page.getByRole('button', { name: 'Save' })).toBeDisabled();

	// make expensive research and adjust hab
	await page.locator('input[name="EnergyResearchCost"]').first().click();
	await page.locator('[data-type="grav-left-button"]').click();

	await page.getByRole('textbox', { name: 'Description' }).fill('slow start, strong late');

	// save button is good, click it
	await expect(page.getByRole('button', { name: 'Save' })).toBeEnabled();
	const saved = page.waitForResponse(
		(response) => response.url().includes('RaceService/UpdateRace') && response.status() === 200
	);
	await page.getByRole('button', { name: 'Save' }).click();
	await saved;

	// the description is saved with the race and shown in the list of races
	await page.goto('/races');
	await expect(page.locator('[data-type="race-description"]').first()).toHaveText(
		'slow start, strong late'
	);
});
