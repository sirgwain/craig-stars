import '../../../../../../css/app.css';
import { writable } from 'svelte/store';
import { expect, it, vi } from 'vitest';
import { page } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';
import PlayerScoresGraph from './PlayerScoresGraph.svelte';

const state = vi.hoisted(() => ({ context: undefined as unknown }));
vi.mock('#lib/services/GameContext.js', () => ({ getGameContext: () => state.context }));

const names = ['Galaxy Improvement Initiative #973c', 'Rabbitoids', 'Triplicates'];
const colors = ['#0000ff', '#00ff00', '#ff0000'];

// Names used to be drawn on the lines, where they ran into each other.
it('lists each player in a legend and only counts whole starbases', async () => {
	await page.viewport(700, 600);
	const playerIntels = names.map((racePluralName, i) => ({
		num: i + 1,
		racePluralName,
		color: colors[i]
	}));
	state.context = {
		game: writable({
			year: 2410,
			rules: { universeGenerationRules: { startingYear: 2400 }, showPublicScoresAfterYears: 0 }
		}),
		player: writable({ num: 1 }),
		universe: writable({
			playerIntels,
			// one or two starbases each, so a default scale would have ticks at halves
			scoreIntels: names.map((_, i) => ({
				scoreHistory: [...Array(10).keys()].map((turn) => ({
					starbases: turn > 4 ? 1 + (i % 2) : 1
				}))
			})),
			getPlayerColor: (num: number) => colors[num - 1]
		})
	};

	const host = document.createElement('div');
	host.style.width = '600px';
	host.style.height = '400px';
	document.body.appendChild(host);
	// the page gives the graph a fixed height, so size it before mounting
	const target = host.appendChild(document.createElement('div'));
	target.style.height = '400px';
	target.style.padding = '20px 30px';
	target.style.boxSizing = 'border-box';
	const { container } = render(
		PlayerScoresGraph,
		{ target, props: { type: 'starbases' } },
		{ baseElement: host }
	);
	await new Promise((resolve) => requestAnimationFrame(() => requestAnimationFrame(resolve)));

	const legend = [...container.querySelectorAll('[data-type="legend"] li')];
	expect(legend.map((li) => li.textContent?.trim())).toEqual(names);
	// the legend sits under the graph, inside the space the graph was given
	const bounds = container.getBoundingClientRect();
	for (const li of legend) {
		expect(li.getBoundingClientRect().bottom).toBeLessThanOrEqual(bounds.bottom + 0.5);
	}

	// and clear of the year labels along the bottom of the graph
	const years = [...container.querySelectorAll('.x-axis .tick')];
	const lowestYear = Math.max(...years.map((t) => t.getBoundingClientRect().bottom));
	expect(legend[0].getBoundingClientRect().top).toBeGreaterThanOrEqual(lowestYear - 0.5);

	const ticks = [...container.querySelectorAll('.y-axis .tick')].map((t) => t.textContent?.trim());
	expect(ticks.length).toBeGreaterThan(1);
	for (const tick of ticks) {
		expect(Number.isInteger(Number(tick))).toBe(true);
	}
});
