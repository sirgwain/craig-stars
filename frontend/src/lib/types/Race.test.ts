import { describe, expect, it } from 'vitest';
import type { Race } from '#lib/types/cs-proto.js';
import { clone, create } from '@bufbuild/protobuf';
import { HabSchema, RulesSchema } from '#lib/types/cs-proto.js';
import { getHabitablePlanetOdds } from './Race';
import { defaultRules } from './Rules';

const hab = (grav: number, temp: number, rad: number) => ({ grav, temp, rad });
const race = (low: number[], high: number[], immune: boolean[] = []) =>
	({
		habLow: hab(low[0], low[1], low[2]),
		habHigh: hab(high[0], high[1], high[2]),
		immuneGrav: immune[0] ?? false,
		immuneTemp: immune[1] ?? false,
		immuneRad: immune[2] ?? false
	}) as Race;

// the default rules with the same dropoff range for every hab type
const withDropoff = (range: number) => {
	const rules = clone(RulesSchema, defaultRules);
	rules.universeGenerationRules!.habDropoffRange = create(HabSchema, {
		grav: range,
		temp: range,
		rad: range
	});
	return rules;
};

// These match the "1 in N planets" line in the original race designer.
describe('getHabitablePlanetOdds', () => {
	it.each([
		{ name: 'Humanoids', race: race([10, 35, 13], [56, 81, 53]), odds: 9 },
		// the narrowest ranges the race designer allows, in the middle of each range
		{ name: 'narrowest ranges', race: race([40, 40, 40], [60, 60, 60]), odds: 93 },
		// planets with extreme gravity are rare, so the same width at the edge finds fewer
		{ name: 'narrow gravity at the low edge', race: race([0, 40, 40], [20, 60, 60]), odds: 126 },
		{ name: 'immune to everything', race: race([0, 0, 0], [0, 0, 0], [true, true, true]), odds: 0 },
		{ name: 'full width ranges', race: race([0, 0, 0], [100, 100, 100]), odds: 0 },
		{ name: 'almost every planet', race: race([0, 0, 0], [99, 0, 99], [false, true]), odds: 1 }
	])('$name is 1 in $odds', ({ race, odds }) => {
		expect(getHabitablePlanetOdds(race)).toBe(odds);
	});

	// The game's rules decide how planet values are spread out.
	it('uses the dropoff ranges from the rules', () => {
		const narrow = race([40, 40, 40], [60, 60, 60]);
		// with no dropoff every value is equally likely, so each range is a straight 20%
		const flat = withDropoff(0);
		expect(getHabitablePlanetOdds(narrow, flat)).toBe(125);
		// a wider dropoff makes the middle values more common, so a centered range finds more planets
		const steep = withDropoff(30);
		expect(getHabitablePlanetOdds(narrow, steep)).toBeLessThan(93);
		// at the edge the same dropoff finds fewer
		expect(getHabitablePlanetOdds(race([0, 0, 0], [20, 20, 20]), steep)).toBeGreaterThan(
			getHabitablePlanetOdds(race([0, 0, 0], [20, 20, 20]))
		);
	});
});
