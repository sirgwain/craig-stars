import { page } from 'vitest/browser';
import { describe, expect, it } from 'vitest';
import { render } from 'vitest-browser-svelte';
import TechEngineGraph from './TechEngineGraph.svelte';

const superScoop = {
	idealSpeed: 9,
	freeSpeed: 7,
	maxSafeSpeed: 9,
	fuelUsage: [0, 0, 0, 0, 0, 0, 0, 0, 65, 90, 108]
};

describe('TechEngineGraph', () => {
	it('shades unsafe warp speeds and highlights the ideal speed', async () => {
		const { container } = render(TechEngineGraph, { engine: superScoop });

		await expect.element(page.getByText('9', { exact: true })).toHaveClass('text-primary');
		await expect.element(page.getByText('8', { exact: true })).not.toHaveClass('text-primary');

		await expect.element(page.getByText('200%')).toHaveClass('text-error');
		await expect.element(page.getByText('100%')).not.toHaveClass('text-error');
		expect(container.querySelectorAll('.y-axis .border-solid').length).toBe(1);

		const region = container.querySelector('rect.region-x');
		expect(region?.getAttribute('x')).toBe('90');
		expect(region?.getAttribute('width')).toBe('10');
	});

	it('does not shade engines that are safe at every warp speed', async () => {
		const { container } = render(TechEngineGraph, {
			engine: { ...superScoop, idealSpeed: 10, maxSafeSpeed: 10 }
		});

		await expect.element(page.getByText('10', { exact: true })).toHaveClass('text-primary');
		expect(container.querySelector('rect.region-x')).toBeNull();
	});
});
