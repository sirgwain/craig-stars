import { apiErrorsFailTest, expect, test } from './setup';

test('ship designer - create', async ({ newGamePage }) => {
	const { page, id } = newGamePage;
	apiErrorsFailTest(page, id);

	await page.locator('label').filter({ hasText: 'Commands' }).click();
	await page.getByRole('link', { name: 'Ship Designer' }).click();

	await page.getByRole('link', { name: 'Create' }).click();
	await page.getByRole('link', { name: 'Small Freighter' }).click();
	await page.waitForURL(`/games/${id}/designer/create/small-freighter`);

	// change hull-set
	const icon = await page.locator('.tech-avatar').first();
	await expect(icon).toHaveClass(/hull-small-freighter-0/);
	await page.locator('[data-type="next-hull-set-button"]').click();
	await expect(icon).toHaveClass(/hull-small-freighter-1/);

	await page.getByRole('button', { name: 'Long Hump' }).getByRole('link').click();
	await page.getByRole('button', { name: 'Engine needs' }).click();
	await page.getByRole('button', { name: 'Shield or Armor Up to' }).click();
	await page.getByRole('button', { name: 'Cow-hide Shield' }).getByRole('link').click();
	await page.getByRole('button', { name: 'Scanner Elec Mech Up to' }).click();
	await page.getByRole('button', { name: 'Rhino Scanner' }).getByRole('link').click();

	// should have this cost
	await expect(
		page.getByText('Cost of one Small Freighter Ironium 21kT Boranium 0kT Germanium 20kT Resources')
	).toBeVisible();

	// remove cow-hide shield, replace with Crobmium
	await page.locator('[data-type="hull-component-button"][data-id="Cow-hide Shield"]').click();
	await page
		.locator('[data-type="delete-hull-component-button"][data-id="Cow-hide Shield"]')
		.click();
	await page.getByRole('button', { name: 'Crobmnium' }).getByRole('link').click();

	// should have new cost
	await expect(
		page.getByText('Cost of one Small Freighter Ironium 25kT Boranium 0kT Germanium 18kT Resources')
	).toBeVisible();

	await page.getByRole('textbox', { name: 'Name' }).fill('My Freighter');

	// start waiting for post response
	const response = page.waitForResponse(
		(response) =>
			response.url().includes(`/api/games/${id}/designs`) &&
			response.request().method() === 'POST' &&
			response.status() === 200,
		{ timeout: 10000 }
	);

	// click save
	await page.getByRole('button', { name: 'Save' }).click();

	// grab the response when it fires
	const { num } = await (await response).json();

	// should route to designer page
	await page.waitForURL(`/games/${id}/designer/${num}`);
});
