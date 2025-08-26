import {
	MinefieldType,
	Prt,
	TechCategory,
	TechHullComponentSchema,
	TechLevelSchema,
	TechOrigin,
	type TechHullComponent
} from '$lib/types/cs-proto';
import { create } from '@bufbuild/protobuf';
import { describe, expect, it } from 'vitest';
import { CommandedPlayer, canLearnTech } from './Player';
import { IFE } from './Race';

const fuelMizer: TechHullComponent = create(TechHullComponentSchema, {
	tech: {
		name: 'Fuel Mizer',
		cost: {
			ironium: 8,
			resources: 11
		},
		requirements: {
			techLevel: {
				propulsion: 2
			},
			lrtsRequired: 1
		},
		ranking: 30,
		category: TechCategory.ENGINE
	},
	hullSlotType: 2,
	mass: 6,
	idealSpeed: 6,
	freeSpeed: 4,
	maxSafeSpeed: 9,
	fuelUsage: [0, 0, 0, 0, 0, 35, 120, 175, 235, 360, 420]
});

const speedTrap20: TechHullComponent = create(TechHullComponentSchema, {
	tech: {
		name: 'Speed Trap 20',
		cost: {
			ironium: 30,
			germanium: 12,
			resources: 60
		},
		requirements: {
			techLevel: {
				propulsion: 2,
				biotechnology: 2
			},
			prtsRequired: [Prt.SD, Prt.IS]
		},
		ranking: 70,
		category: TechCategory.MINE_LAYER
	},
	hullSlotType: 8192,
	mass: 100,
	minefieldType: MinefieldType.SPEED_BUMP,
	mineLayingRate: 20
});

const smartBomb: TechHullComponent = create(TechHullComponentSchema, {
	tech: {
		name: 'Smart Bomb',
		cost: {
			ironium: 1,
			boranium: 22,
			resources: 27
		},
		requirements: {
			techLevel: {
				weapons: 5,
				biotechnology: 7
			},
			prtsDenied: [Prt.IS]
		},
		ranking: 90,
		category: TechCategory.BOMB
	},
	hullSlotType: 16,
	mass: 50,
	killRate: 1.3,
	smart: true
});

const multiFunctionPod: TechHullComponent = create(TechHullComponentSchema, {
	tech: {
		name: 'Multi-Function Pod',
		cost: {
			ironium: 5,
			germanium: 5,
			resources: 15
		},
		requirements: {
			techLevel: {
				energy: 11,
				propulsion: 11,
				electronics: 11
			},
			acquirable: true
		},
		ranking: 35,
		category: TechCategory.ELECTRICAL,
		origin: TechOrigin.MYSTERY_TRADER
	},
	hullSlotType: 64,
	mass: 2,
	cloakUnits: 60,
	torpedoJamming: 0.1,
	movementBonus: 1
});

describe('player test', () => {
	it('checks tech requirements', () => {
		const player = new CommandedPlayer();

		expect(canLearnTech(player, fuelMizer)).toBe(false);

		// get the level but not the LRT
		player.techLevels.propulsion = 2;
		expect(canLearnTech(player, fuelMizer)).toBe(false);

		// make this available
		player.race.lrts = IFE;
		expect(canLearnTech(player, fuelMizer)).toBe(true);

		// IS can learn speed trap
		player.race.prt = Prt.IS;
		expect(canLearnTech(player, speedTrap20)).toBe(true);

		// IS cannot learn smart bomb
		player.race.prt = Prt.IS;
		expect(canLearnTech(player, smartBomb)).toBe(false);

		// SD can learn speed trap
		player.race.prt = Prt.SD;
		expect(canLearnTech(player, speedTrap20)).toBe(true);
	});

	it('checks has tech', () => {
		const player = new CommandedPlayer();
		player.techLevels.propulsion = 2;
		expect(player.hasTech(fuelMizer)).toBe(false);

		// make it available
		player.race.lrts = IFE;
		expect(player.hasTech(fuelMizer)).toBe(true);

		// player doesn't have MT tech until acquired
		player.techLevels = create(TechLevelSchema, { energy: 11, propulsion: 11, electronics: 11 });
		expect(player.hasTech(multiFunctionPod)).toBe(false);

		player.acquiredTechs[multiFunctionPod?.tech?.name ?? ''] = true;
		expect(player.hasTech(multiFunctionPod)).toBe(true);
	});
});
