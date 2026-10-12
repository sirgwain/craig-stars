import { apiErrorsFailTest, expect, test } from './setup';
import { CreateShipDesignResponseJson } from '../src/lib/protogen/craig_stars/v1/shipdesignservice_pb';

test('ship designer - create', async ({ newGamePage }) => {
	const { page, id } = newGamePage;
	apiErrorsFailTest(page);

	await page.locator('label').filter({ hasText: 'Commands' }).click();
	await page.getByRole('link', { name: 'Ship Designer' }).click();
	await page.waitForURL(`/games/${id}/designer`);

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
	await page.getByRole('textbox', { name: 'Notes' }).fill('early colony supply runs');

	// start waiting for post response
	const response = page.waitForResponse(
		(response) =>
			response.url().includes(`/api/grpc/craig_stars.v1.ShipDesignService`) &&
			response.request().method() === 'POST' &&
			response.status() === 200,
		{ timeout: 10000 }
	);

	// click save
	await page.getByRole('button', { name: 'Save' }).click();

	// grab the response when it fires
	const { design } = (await (await response).json()) as CreateShipDesignResponseJson;
	if (!design?.num) {
		throw new Error('failed to create ship design');
	}

	// should route to designer page
	await page.waitForURL(`/games/${id}/designer/${design.num}`);
	await expect(page.getByRole('textbox', { name: 'Design notes' })).toHaveValue(
		'early colony supply runs'
	);
});

test('ship designer - notes on a design that is in use', async ({ testGamePage }) => {
	const { page, gameId, universe } = await testGamePage('Kitchen Sink');
	apiErrorsFailTest(page);
	const design = universe.designs.find((d) => d.playerNum === 1 && d.spec?.numInstances);
	if (!design) {
		throw new Error('failed to find a design with ships');
	}

	await page.goto(`/games/${gameId}/designer/${design.num}`);
	const notes = page.getByRole('textbox', { name: 'Design notes' });
	await notes.fill('fast scout for the early game');
	await page.getByRole('button', { name: 'Save notes' }).click();
	await expect(page.getByRole('button', { name: 'Save notes' })).toBeHidden();

	await page.reload();
	await expect(notes).toHaveValue('fast scout for the early game');

	// the notes are shown with the design in the list of designs
	await page.goto(`/games/${gameId}/designer`);
	await expect(page.locator('[data-type="design-description"]')).toHaveText(
		'fast scout for the early game'
	);
});
