import { describe, expect, it } from 'vitest';
import { CommandedFleet } from './Fleet';

describe('fleet test', () => {
	const fleet = new CommandedFleet();
	it('returns minimal speeds for distances', () => {
		// one year to go 49 ly
		expect(fleet.getMinimalWarp(49, 7)).toBe(7);
		// two years at warp 7, two at warp 6 or warp 5, pick warp 5
		expect(fleet.getMinimalWarp(50, 7)).toBe(5);

		// warp 6 takes 3 years to 72, so in 73 we can make it in 2 at warp 7
		expect(fleet.getMinimalWarp(73, 7)).toBe(7);

		// this is obvious
		expect(fleet.getMinimalWarp(36, 7)).toBe(6);

		// might as well go warp 5
		expect(fleet.getMinimalWarp(25, 7)).toBe(5);
	});
});
