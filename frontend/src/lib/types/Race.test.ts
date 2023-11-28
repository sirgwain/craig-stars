import { describe, expect, it } from 'vitest';
import { getPlanetHabitability, humanoid } from './Race';

describe('Race test', () => {
	it('getPlanetHabitability', () => {
		const race = humanoid();

		expect(getPlanetHabitability(race, { grav: 50, temp: 50, rad: 50 })).toBe(100);
		expect(getPlanetHabitability(race, { grav: 0, temp: 0, rad: 0 })).toBe(-45);
	});
});
