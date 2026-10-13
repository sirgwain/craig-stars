import { expect, test, createTestGame } from './setup';

for (const viewport of [
	{ width: 390, height: 844 },
	{ width: 768, height: 1024 },
	{ width: 1536, height: 960 }
]) {
	test(`library screens fit ${viewport.width}px`, async ({ authenticatedPage: page }) => {
		await page.setViewportSize(viewport);
		await page.addInitScript(() => {
			localStorage.theme = 'business';
		});
		const created = await createTestGame(
			page,
			'Kitchen Sink',
			`A long single player galaxy name for layout ${viewport.width}`
		);
		try {
			for (const route of ['/', '/races', '/techs']) {
				await page.goto(route);
				await expect(
					page.getByRole('heading', {
						name: route === '/' ? 'Games' : route === '/races' ? 'Races' : 'Techs',
						exact: true
					})
				).toBeVisible();
				await expect(page.locator('.loading')).toHaveCount(0);
				expect(
					await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)
				).toBe(true);
				if (route === '/') {
					const solo = page.getByRole('region', { name: 'Single player games' });
					await expect(solo.locator('article')).toHaveCount(5);
					await solo.getByRole('button', { name: /^Show all/ }).click();
					expect(await solo.locator('article').count()).toBeGreaterThan(5);
					await solo.getByRole('button', { name: 'Show fewer' }).click();
					await page.getByLabel('Search single player games').fill(`layout ${viewport.width}`);
					expect(await solo.locator('article').count()).toBeGreaterThan(0);
					await page.getByLabel('Search single player games').fill('');
				}
				if (route === '/techs') {
					if (viewport.width < 1024)
						await page.getByLabel('Technology category').selectOption({ label: 'Engine' });
					else await page.getByRole('button', { name: 'Engine', exact: true }).click();
				}
				await page.screenshot({
					path: `../tmp/library-${route === '/' ? 'games' : route.slice(1)}-${viewport.width}.png`,
					fullPage: false
				});
			}
			await expect(page.getByRole('link', { name: 'Tritanium', exact: true })).toHaveCount(0);
			await expect(page.locator('[data-type="tech-card"]')).not.toHaveCount(0);
			// engines show their fuel usage graph in place
			await page.getByRole('button', { name: 'usage…', exact: true }).first().click();
			await expect(page.locator('[data-type="tech-graph"]').first()).toBeVisible();
			// long specs are clamped behind a "more…" toggle
			await page.getByLabel('Search technologies').fill('Mine Dispenser');
			if (viewport.width < 1024)
				await page.getByLabel('Technology category').selectOption({ label: 'All technologies' });
			else await page.getByRole('button', { name: 'All technologies', exact: true }).click();
			const more = page.getByRole('button', { name: 'more…', exact: true }).first();
			await more.click();
			await expect(page.getByRole('button', { name: 'less', exact: true }).first()).toHaveAttribute(
				'aria-expanded',
				'true'
			);
			// planetary defenses show their coverage graph in place
			await page.getByLabel('Search technologies').fill('Missile Battery');
			await page.getByRole('button', { name: 'coverage…', exact: true }).first().click();
			await expect(page.locator('[data-type="tech-graph"]').first()).toBeVisible();
			expect(
				await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)
			).toBe(true);
		} finally {
			await page.request.post('/api/grpc/craig_stars.v1.GameService/DeleteGame', {
				data: { gameId: String(created.game?.id) }
			});
		}
	});
}
