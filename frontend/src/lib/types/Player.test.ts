import { describe, expect, it } from 'vitest';
import { canLearnTech, Player } from './Player';
import { LRT } from './Race';
import { TechCategory, type TechEngine } from './Tech';

const fuelMizer: TechEngine = {
	name: 'Fuel Mizer',
	cost: {
		ironium: 8,
		resources: 11
	},
	requirements: {
		propulsion: 2,
		lrtsRequired: 1
	},
	ranking: 30,
	category: TechCategory.Engine,
	hullSlotType: 2,
	mass: 6,
	idealSpeed: 6,
	freeSpeed: 4,
	fuelUsage: [0, 0, 0, 0, 0, 35, 120, 175, 235, 360, 420]
};

describe('player test', () => {
	it('checks tech requirements', () => {
		const player = new Player();

		expect(canLearnTech(player, fuelMizer)).toBe(false);

		// get the level but not the LRT
		player.techLevels.propulsion = 2;
		expect(canLearnTech(player, fuelMizer)).toBe(false);

		// make this available
		player.race.lrts = LRT.IFE;
		expect(canLearnTech(player, fuelMizer)).toBe(true);
	});
});
