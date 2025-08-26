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

	// save button is good, click it
	await expect(page.getByRole('button', { name: 'Save' })).toBeEnabled();
	await page.getByRole('button', { name: 'Save' }).click();
});
