import '../../../../../../css/app.css';
import { Universe } from '#lib/services/Universe.js';
import type { PlayerIntel } from '#lib/types/cs-proto.js';
import { writable } from 'svelte/store';
import { expect, it, vi } from 'vitest';
import { page } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';
import PlayerScores from './PlayerScores.svelte';

const state = vi.hoisted(() => ({ context: undefined as unknown }));
vi.mock('#lib/services/GameContext.js', () => ({ getGameContext: () => state.context }));

const names = [
	'Galaxy Improvement Initiative #973c',
	'muffin.man.the.first',
	'Triplicates',
	'Gwainions',
	'Valkonians'
];

// Long race names used to run into each other, and the row labels into the first player's scores.
it.each([390, 1280])('keeps names and scores apart at %i wide', async (width) => {
	await page.viewport(width, 800);
	const universe = new Universe();
	universe.playerIntels = names.map(
		(racePluralName, i) => ({ num: i + 1, racePluralName }) as PlayerIntel
	);
	universe.scoreIntels = names.map(
		() => ({ scoreHistory: [{ planets: 21, unarmedShips: 298 }] }) as never
	);
	state.context = { universe: writable(universe) };

	const { container } = render(PlayerScores);
	const rect = (el: Element) => el.getBoundingClientRect();

	// each name stays inside its own column
	const headers = [...container.querySelectorAll('thead td')];
	expect(headers).toHaveLength(names.length);
	headers.forEach((cell, i) => {
		const name = rect(cell.querySelector('div')!);
		expect(name.left).toBeGreaterThanOrEqual(rect(cell).left - 0.5);
		expect(name.right).toBeLessThanOrEqual(rect(cell).right + 0.5);
		if (i > 0) {
			expect(rect(cell).left).toBeGreaterThanOrEqual(rect(headers[i - 1]).right - 0.5);
		}
	});

	// row labels stay on one line, clear of the first score
	for (const row of container.querySelectorAll('tbody tr')) {
		const label = row.querySelector('th')!;
		expect(label.scrollWidth).toBeLessThanOrEqual(label.clientWidth);
		expect(rect(label).right).toBeLessThanOrEqual(rect(row.querySelector('td')!).left + 0.5);
	}

	// five players fit on a phone without scrolling
	const scroller = container.querySelector<HTMLElement>('.overflow-x-auto')!;
	expect(scroller.scrollWidth).toBeLessThanOrEqual(scroller.clientWidth);
});
