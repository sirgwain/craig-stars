import { expect, test } from '@playwright/test';
import techs from '../src/lib/ssr/techs.json' with { type: 'json' };

for (const width of [390, 768, 1344, 1536]) {
	test(`tech layout uses available width at ${width}px`, async ({ page }) => {
		await page.setViewportSize({ width, height: 960 });
		await page.addInitScript(() => {
			localStorage.theme = 'business';
		});
		await page.route('**/*UserService/GetMe*', (route) =>
			route.fulfill({ json: { user: { id: '1', username: 'Preview', role: 'USER_ROLE_USER' } } })
		);
		await page.route('**/*TechService/GetTechs*', (route) => route.fulfill({ json: techs }));
		await page.goto('/techs');
		const cards = page.locator('[data-type="tech-card"]');
		await expect(cards.first()).toBeVisible();
		expect(
			await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)
		).toBe(true);
		if (width >= 1024) {
			const heading = await page.getByRole('heading', { name: 'Techs', exact: true }).boundingBox();
			expect(heading!.width).toBeGreaterThan(width * 0.9);
			const first = await cards.nth(0).boundingBox();
			const second = await cards.nth(1).boundingBox();
			const third = await cards.nth(2).boundingBox();
			expect(first!.width).toBeGreaterThan(450);
			expect(second!.y).toBe(first!.y);
			expect(third!.y).toBeGreaterThan(first!.y);
			const costs = await cards
				.first()
				.locator('dt')
				.evaluateAll((labels) =>
					labels.filter((l) => l.textContent !== 'Mass').map((l) => l.getBoundingClientRect().top)
				);
			expect(costs).toHaveLength(4);
			expect(new Set(costs).size).toBe(1);
		}
		if (width === 1344) await page.screenshot({ path: '../tmp/tech-wide-fixed.png' });
	});
}
