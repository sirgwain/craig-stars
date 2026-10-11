import '../../../../css/app.css';
import type { Race } from '#lib/types/cs-proto.js';
import type { WasmClient } from '#lib/wasm.js';
import { expect, it } from 'vitest';
import { page } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';
import RacePoints from './RacePoints.svelte';

// The points box floats over the race form. It used to take clicks across its whole row, and
// stayed full size however far the page was scrolled.
it('shrinks into the corner when the page scrolls, and only the box takes clicks', async () => {
	await page.viewport(390, 700);
	const wasmClient = {
		calculateRacePoints: async () => ({ points: 25 })
	} as unknown as WasmClient;
	const filler = document.createElement('div');
	filler.style.height = '3000px';
	const { container } = render(RacePoints, { wasmClient, race: {} as Race });
	document.body.appendChild(filler);

	const row = container.querySelector<HTMLElement>('[data-type="race-points"]')!;
	const box = row.querySelector<HTMLElement>('.stats')!;
	await expect.element(page.getByText('25')).toBeInTheDocument();
	expect(getComputedStyle(row).pointerEvents).toBe('none');
	expect(getComputedStyle(box).pointerEvents).toBe('auto');

	const full = box.getBoundingClientRect();
	window.scrollTo(0, 400);
	await expect.poll(() => box.dataset.compact).toBe('true');
	// wait out the shrink animation
	await expect.poll(() => box.getBoundingClientRect().width).toBeLessThan(full.width * 0.6);
	const small = box.getBoundingClientRect();
	// it stays against the right edge
	expect(small.right).toBeCloseTo(full.right, 0);

	window.scrollTo(0, 0);
	await expect.poll(() => box.dataset.compact).toBe('false');
	filler.remove();
});
