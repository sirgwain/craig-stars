import { PlanetSchema } from '#lib/types/cs-proto.js';
import { create } from '@bufbuild/protobuf';
import { describe, expect, it } from 'vitest';
import { planetsSortBy } from './Planet';

describe('planetsSortBy', () => {
	const planets = [
		create(PlanetSchema, {
			mapObject: { num: 1 },
			cargo: { ironium: 10, boranium: 300, germanium: 5 },
			mineralConcentration: { ironium: 90, boranium: 10, germanium: 50 },
			spec: { miningOutput: { ironium: 1, boranium: 2, germanium: 30 } }
		}),
		create(PlanetSchema, {
			mapObject: { num: 2 },
			cargo: { ironium: 200, boranium: 20, germanium: 50 },
			mineralConcentration: { ironium: 20, boranium: 80, germanium: 40 },
			spec: { miningOutput: { ironium: 20, boranium: 1, germanium: 3 } }
		})
	];
	const sorted = (key: string) =>
		[...planets].sort(planetsSortBy(key)).map((p) => p.mapObject?.num);

	it('sorts mineral columns by total minerals', () => {
		expect(sorted('minerals')).toEqual([2, 1]);
	});

	it.each([
		{ key: 'minerals.ironium', want: [1, 2] },
		{ key: 'minerals.boranium', want: [2, 1] },
		{ key: 'minerals.germanium', want: [1, 2] },
		{ key: 'miningRate.ironium', want: [1, 2] },
		{ key: 'miningRate.germanium', want: [2, 1] },
		{ key: 'mineralConcentration.ironium', want: [2, 1] },
		{ key: 'mineralConcentration.boranium', want: [1, 2] }
	])('sorts by $key', ({ key, want }) => {
		expect(sorted(key)).toEqual(want);
	});

	describe('all planets report', () => {
		const names: Record<number, string> = { 1: 'Zebras', 2: 'Aardvarks' };
		const finder = { getPlayerPluralName: (num: number | undefined) => names[num ?? 0] ?? '' };
		const planet = (name: string, playerNum: number, reportAge: number) =>
			create(PlanetSchema, { mapObject: { name, playerNum, reportAge } });
		const report = [
			planet('Unowned', 0, 3),
			planet('Zed', 1, 0),
			planet('Beta', 2, 10),
			planet('Unexplored', 0, -1),
			planet('Alpha', 2, 1)
		];
		const sortedNames = (key: string) =>
			[...report]
				.sort(planetsSortBy(key, finder as Parameters<typeof planetsSortBy>[1]))
				.map((p) => p.mapObject?.name);

		it('sorts owners by name, each with their planets in order, and unowned planets last', () => {
			expect(sortedNames('owner')).toEqual(['Alpha', 'Beta', 'Zed', 'Unexplored', 'Unowned']);
		});

		it('sorts report age with the newest first and unexplored planets last', () => {
			expect(sortedNames('reportAge')).toEqual(['Zed', 'Alpha', 'Unowned', 'Beta', 'Unexplored']);
		});
	});
});
