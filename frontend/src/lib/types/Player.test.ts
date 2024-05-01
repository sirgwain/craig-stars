import techjson from '$lib/ssr/techs.json';
import { describe, expect, it } from 'vitest';
import { Player, canLearnTech } from './Player';
import { LRT } from './Race';
import type { ShipDesign } from './ShipDesign';
import { TechCategory, type TechEngine, type TechStore } from './Tech';

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

export const baseStationDesign: ShipDesign = {
	num: 1,
	playerNum: 1,
	gameId: 1,
	originalPlayerNum: 0,
	version: 0,
	name: 'Base Station',
	hull: 'Space Station',
	hullSetNumber: 0,
	slots: [],
	spec: {
		hullType: 'Starbase',
		engine: {},
		cost: {
			ironium: 92,
			boranium: 61,
			germanium: 190,
			resources: 456
		},
		techLevel: {},
		armor: 500,
		scanRange: -1,
		scanRangePen: -1,
		repairBonus: 0.15,
		torpedoInaccuracyFactor: 1,
		initiative: 14,
		starbase: true,
		spaceDock: -1,
		cloakPercentFullCargo: -4611686018427388000,
		maxPopulation: 1000000
	}
};

export const orbitalFort2Design: ShipDesign = {
	id: 2006,
	gameId: 60,
	num: 20,
	playerNum: 1,
	originalPlayerNum: 0,
	name: 'Orbital Fort II',
	version: 0,
	hull: 'Orbital Fort',
	hullSetNumber: 0,
	slots: [
		{
			hullComponent: 'Beta Torpedo',
			hullSlotIndex: 2,
			quantity: 12
		},
		{
			hullComponent: 'Beta Torpedo',
			hullSlotIndex: 4,
			quantity: 12
		},
		{
			hullComponent: 'Wolverine Diffuse Shield',
			hullSlotIndex: 3,
			quantity: 12
		},
		{
			hullComponent: 'Organic Armor',
			hullSlotIndex: 5,
			quantity: 12
		}
	],
	spec: {
		hullType: 'OrbitalFort',
		engine: {},
		cost: {
			ironium: 463,
			boranium: 144,
			germanium: 230,
			resources: 493
		},
		techLevel: {
			energy: 6,
			weapons: 5,
			propulsion: 1,
			biotechnology: 7
		},
		mass: 792,
		armor: 2200,
		scanRange: -1,
		scanRangePen: -1,
		repairBonus: 0.03,
		torpedoInaccuracyFactor: 1,
		initiative: 10,
		powerRating: 288,
		shields: 720,
		starbase: true,
		hasWeapons: true,
		weaponSlots: [
			{
				hullComponent: 'Beta Torpedo',
				hullSlotIndex: 2,
				quantity: 12
			},
			{
				hullComponent: 'Beta Torpedo',
				hullSlotIndex: 4,
				quantity: 12
			}
		],
		maxPopulation: 250000
	}
};

const techStore = techjson as TechStore;

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

	it('getTerraformAbility', () => {
		const player = new Player();

		expect(player.getTerraformAbility(techStore)).toEqual({ grav: 0, temp: 0, rad: 0 });

		// get some tech
		player.techLevels = {
			energy: 3,
			weapons: 3,
			propulsion: 3,
			construction: 3,
			electronics: 3,
			biotechnology: 3
		};
		expect(player.getTerraformAbility(techStore)).toEqual({ grav: 3, temp: 3, rad: 3 });
	});

	it('getStarbaseUpgradeCost', () => {
		const player = new Player();

		expect(player.getStarbaseUpgradeCost(techStore, orbitalFort2Design, baseStationDesign)).toEqual(
			{ ironium: 0, boranium: 0, germanium: 0, resources: 0 }
		);
	});
});
