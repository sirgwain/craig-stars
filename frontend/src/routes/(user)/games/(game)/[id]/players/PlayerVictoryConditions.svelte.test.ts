import '../../../../../../css/app.css';
import { VictoryConditionAttainTechLevels, VictoryConditionOwnPlanets } from '#lib/types/Consts.js';
import { writable } from 'svelte/store';
import { expect, it, vi } from 'vitest';
import { page } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';
import PlayerVictoryConditions from './PlayerVictoryConditions.svelte';

const state = vi.hoisted(() => ({ context: undefined as unknown }));
vi.mock('#lib/services/GameContext.js', () => ({ getGameContext: () => state.context }));

const names = [
	'Claim Adjusters',
	'Plaguebringers',
	'Udines',
	'Necro Barbarians',
	'Techno Vikings',
	'Berserkers'
];

// With six players on a phone the condition column used to be squeezed to one letter per line,
// and the slanted names ran into each other.
it('keeps conditions readable and names apart on a phone', async () => {
	await page.viewport(390, 800);
	const conditions = VictoryConditionOwnPlanets | VictoryConditionAttainTechLevels;
	state.context = {
		game: writable({
			victoryConditions: {
				conditions,
				ownPlanets: 60,
				attainTechLevel: 22,
				attainTechLevelNumFields: 4,
				numCriteriaRequired: 1,
				yearsPassed: 50
			}
		}),
		universe: writable({
			playerIntels: names.map((racePluralName, i) => ({ num: i + 1, racePluralName })),
			planetIntels: Array(76),
			planets: Array(76),
			getPlayerPluralName: (num: number) => names[num - 1],
			getPlayerScore: (num: number) => ({ achievedVictoryConditions: num === 2 ? conditions : 0 })
		})
	};

	const { container } = render(PlayerVictoryConditions);
	const rect = (el: Element) => el.getBoundingClientRect();

	// each condition fits on a couple of lines
	for (const label of container.querySelectorAll('tbody th')) {
		expect(rect(label).width).toBeGreaterThanOrEqual(140);
		expect(rect(label).height).toBeLessThan(80);
	}

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
});
