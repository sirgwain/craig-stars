import { expect, test } from './setup';

for (const width of [1280, 390]) {
	for (const route of ['single-player-game', 'host-game']) {
		test(`${route} setup at ${width}px`, async ({ authenticatedPage: page }, testInfo) => {
			await page.setViewportSize({ width, height: 844 });
			await page.emulateMedia({ colorScheme: 'dark' });
			await page.goto(`/${route}`);
			await expect(page.getByRole('heading', { name: 'Galaxy', exact: true })).toBeVisible();
			const slots = page.locator('[data-type="player-slot"]');
			await expect(slots).toHaveCount(3);
			if (route === 'host-game') {
				await expect(slots.nth(1).getByLabel('Difficulty')).toHaveCount(0);
				await expect(slots.nth(1).getByText('Open slot', { exact: false })).toBeVisible();
			} else {
				await expect(slots.nth(1).getByLabel('Difficulty')).toBeVisible();
			}
			await page.getByRole('button', { name: 'Slower tech (2×)' }).click();
			await expect(page.getByRole('slider', { name: 'Technology cost' })).toHaveValue('2');
			await page.getByRole('slider', { name: 'Technology cost' }).focus();
			await page.keyboard.press('ArrowRight');
			await page.keyboard.press('ArrowRight');
			if (width > 1000)
				await expect(page.getByRole('complementary', { name: 'Game summary' })).toContainText(
					'2.5×'
				);
			await page.getByRole('button', { name: 'Add player' }).click();
			await expect(slots).toHaveCount(4);
			await page.getByRole('button', { name: 'Remove slot 2' }).click();
			await expect(slots).toHaveCount(3);
			await expect(slots.nth(1).getByLabel('Difficulty')).toBeVisible();
			await page.locator('summary').filter({ hasText: 'Victory Conditions' }).click();
			await page.getByRole('spinbutton', { name: 'Victory tech level', exact: true }).fill('24');
			await expect(
				page.getByRole('spinbutton', { name: 'Victory tech level', exact: true })
			).toHaveValue('24');
			expect(
				await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)
			).toBe(true);
			await page.evaluate(() => window.scrollTo(0, 0));
			await testInfo.attach(`${route}-${width}`, {
				body: await page.screenshot({
					fullPage: true,
					style: '.sticky { position: static !important; }'
				}),
				contentType: 'image/png'
			});
		});
	}
}

test('technology cost persists through hosted game setup and reload', async ({
	authenticatedPage: page
}) => {
	await page.goto('/host-game');
	await page.getByRole('textbox', { name: 'Name', exact: true }).fill(`Slow tech ${Date.now()}`);
	await page.getByRole('button', { name: 'Slower tech (2×)' }).click();
	const created = page.waitForResponse(
		(r) => r.url().includes('GameService/CreateGame') && r.request().method() === 'POST'
	);
	await page.getByRole('button', { name: 'Create Game', exact: true }).click();
	const response = await created;
	expect(response.ok()).toBe(true);
	const result = await response.json();
	expect(result.game.game.techCostMultiplier).toBe(2);
	await page.waitForURL(/\/games\/\d+/);
	await page.reload();
	await expect(page.getByRole('slider', { name: 'Technology cost' })).toHaveValue('2');
});
