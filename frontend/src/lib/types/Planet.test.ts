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
});
