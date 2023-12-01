import techjson from '$lib/ssr/techs.json';
import { multiply, type Cost } from '$lib/types/Cost';
import { CommandedPlanet } from '$lib/types/Planet';
import { Player } from '$lib/types/Player';
import { QueueItemTypes } from '$lib/types/QueueItemType';
import { defaultRules } from '$lib/types/Rules';
import type { ShipDesign } from '$lib/types/ShipDesign';
import type { TechStore } from '$lib/types/Tech';
import { cloneDeep } from 'lodash-es';
import { describe, expect, it } from 'vitest';
import { getNumBuilt, getProductionEstimates, produce } from './Producer';
import type { DesignFinder } from './Universe';

describe('Producer test', () => {
	const techStore = techjson as TechStore;

	const factoryCost: Cost = {
		germanium: 4,
		resources: 10
	};

	class TestDesignFinder implements DesignFinder {
		getDesign(playerNum: number, num: number): ShipDesign | undefined {
			throw new Error('Method not implemented.');
		}
		getMyDesign(num: number | undefined): ShipDesign | undefined {
			return {
				id: 1085,
				gameId: 42,
				num: 1,
				playerNum: 1,
				originalPlayerNum: 0,
				name: 'Long Range Scout',
				version: 0,
				hull: 'Scout',
				hullSetNumber: 0,
				slots: [
					{
						hullComponent: 'Long Hump 6',
						hullSlotIndex: 1,
						quantity: 1
					},
					{
						hullComponent: 'Rhino Scanner',
						hullSlotIndex: 2,
						quantity: 1
					},
					{
						hullComponent: 'Fuel Tank',
						hullSlotIndex: 3,
						quantity: 1
					}
				],
				purpose: 'Scout',
				spec: {
					hullType: 'Scout',
					engine: {
						idealSpeed: 6,
						freeSpeed: 1,
						maxSafeSpeed: 9,
						fuelUsage: [0, 0, 20, 60, 100, 100, 105, 450, 750, 900, 1080]
					},
					numEngines: 1,
					cost: {
						ironium: 17,
						boranium: 2,
						germanium: 7,
						resources: 22
					},
					techLevel: {},
					mass: 25,
					armor: 20,
					fuelCapacity: 300,
					scanRange: 66,
					scanRangePen: 30,
					torpedoInaccuracyFactor: 1,
					initiative: 1,
					movement: 4,
					scanner: true,
					numInstances: 1,
					numBuilt: 1,
					estimatedRange: 2272,
					estimatedRangeFull: 2272
				}
			};
		}
	}

	it('processQueueItem - 1 factory', () => {
		// process building a single factory
		const result = getNumBuilt(
			{
				type: QueueItemTypes.Factory,
				quantity: 1
			},
			factoryCost,
			multiply(factoryCost, 10),
			100
		);
		expect(result.numBuilt).toBe(1);
	});

	it('processQueueItem - 20 factories, available resources for 10', () => {
		// process 20 factories with only resources available for 10
		const result = getNumBuilt(
			{
				type: QueueItemTypes.Factory,
				quantity: 20
			},
			factoryCost,
			multiply(factoryCost, 10),
			100
		);
		expect(result.numBuilt).toBe(10);
	});

	it('processQueueItem - 20 factories, available resources for 10.5, partial allocation', () => {
		// process 20 factories with only resources available for 10
		const result = getNumBuilt(
			{
				type: QueueItemTypes.Factory,
				quantity: 20
			},
			factoryCost,
			multiply(factoryCost, 10.5),
			100
		);
		expect(result.numBuilt).toBe(10);
	});

	it('produce - 1 mine, 1 factory and be done', () => {
		const planet = new CommandedPlanet();
		const player = new Player();
		planet.hab = { grav: 50, temp: 50, rad: 50 };
		planet.cargo = {
			ironium: 100,
			boranium: 100,
			germanium: 100,
			colonists: 1000
		};
		planet.productionQueue = [
			{ type: QueueItemTypes.Mine, quantity: 1 },
			{ type: QueueItemTypes.Factory, quantity: 1 }
		];

		const designFinder = new TestDesignFinder();

		// we should build one mine and be done
		const result = produce(defaultRules, techStore, planet, player, designFinder, 1);
		expect(result).toBeTruthy();
		expect(planet.mines).toBe(1);
		expect(planet.factories).toBe(1);
	});

	it('estimate 1 factory, no resources, infinite build time', () => {
		const planet = new CommandedPlanet();
		const player = new Player();
		planet.hab = { grav: 50, temp: 50, rad: 50 };
		planet.cargo = {
			ironium: 0,
			boranium: 0,
			germanium: 0,
			colonists: 2500
		};
		planet.productionQueue = [{ type: QueueItemTypes.Factory, quantity: 1 }];
		const startingPlanet = cloneDeep(planet);

		const designFinder = new TestDesignFinder();

		// we should build one mine and be done
		const items = getProductionEstimates(defaultRules, techStore, player, planet, designFinder);

		// make sure the planet didn't update
		expect(planet).toEqual(startingPlanet);

		// should return a new production queue with estimates in it
		expect(items).toMatchObject([
			{
				type: QueueItemTypes.Factory,
				yearsToBuildOne: 100,
				yearsToBuildAll: 100,
				yearsToSkipAuto: 100
			}
		]);
	});

	it('estimate 1 mine, 1 factory, both in a year', () => {
		const planet = new CommandedPlanet();
		const player = new Player();
		planet.hab = { grav: 50, temp: 50, rad: 50 };
		planet.cargo = {
			ironium: 100,
			boranium: 100,
			germanium: 100,
			colonists: 1000
		};
		planet.productionQueue = [
			{ type: QueueItemTypes.Mine, quantity: 1 },
			{ type: QueueItemTypes.Factory, quantity: 1 }
		];
		const startingPlanet = cloneDeep(planet);

		const designFinder = new TestDesignFinder();

		// we should build one mine and be done
		const items = getProductionEstimates(defaultRules, techStore, player, planet, designFinder);

		// make sure the planet didn't update
		expect(planet).toEqual(startingPlanet);

		// should return a new production queue with estimates in it
		expect(items).toMatchObject([
			{
				type: QueueItemTypes.Mine,
				yearsToBuildOne: 1,
				yearsToBuildAll: 1,
				yearsToSkipAuto: 100
			},
			{
				type: QueueItemTypes.Factory,
				yearsToBuildOne: 1,
				yearsToBuildAll: 1,
				yearsToSkipAuto: 100
			}
		]);
	});

	it('estimate 1 mine, 1 factory, 5 resources per year', () => {
		const planet = new CommandedPlanet();
		const player = new Player();
		planet.hab = { grav: 50, temp: 50, rad: 50 };

		// we have enough minerals to build factories, but not enough resources
		planet.cargo = {
			ironium: 100,
			boranium: 100,
			germanium: 100
		};
		planet.population = 5_000; // 5_000 colonists generate 5 resources per year

		planet.productionQueue = [
			{ type: QueueItemTypes.Mine, quantity: 1 },
			{ type: QueueItemTypes.Factory, quantity: 1 }
		];

		const designFinder = new TestDesignFinder();

		// we should build one mine and then take 2 years to build the factory
		// the factory costs 10 resources so it takes 2 years to build one
		const items = getProductionEstimates(defaultRules, techStore, player, planet, designFinder);

		// should return a new production queue with estimates in it
		expect(items).toMatchObject([
			{
				type: QueueItemTypes.Mine,
				yearsToBuildOne: 1,
				yearsToBuildAll: 1,
				yearsToSkipAuto: 100
			},
			{
				type: QueueItemTypes.Factory,
				yearsToBuildOne: 3,
				yearsToBuildAll: 3,
				yearsToSkipAuto: 100
			}
		]);
	});

	it('estimate 1 auto mine, 1 auto factory, 5 resources per year', () => {
		const planet = new CommandedPlanet();
		const player = new Player();
		planet.hab = { grav: 50, temp: 50, rad: 50 };

		// we have enough minerals to build factories, but not enough resources
		planet.cargo = {
			ironium: 100,
			boranium: 100,
			germanium: 100
		};
		planet.population = 5_000; // 5_000 colonists generate 5 resources per year

		planet.productionQueue = [
			{ type: QueueItemTypes.AutoMines, quantity: 1 },
			{ type: QueueItemTypes.AutoFactories, quantity: 1 }
		];

		const designFinder = new TestDesignFinder();

		// we build a mine each year, but we have to wait for our planet to grow to build
		// the factory, which costs 10 resources and we keep auto building the mine each year
		const items = getProductionEstimates(defaultRules, techStore, player, planet, designFinder);

		// should return a new production queue with estimates in it
		expect(items).toMatchObject([
			{
				type: QueueItemTypes.AutoMines,
				yearsToBuildOne: 1,
				yearsToBuildAll: 1,
				yearsToSkipAuto: 100
			},
			{
				type: QueueItemTypes.AutoFactories,
				yearsToBuildOne: 10,
				yearsToBuildAll: 10,
				yearsToSkipAuto: 100
			}
		]);
	});

	it('estimate 5 MaxTerraform, 20 auto factories, 100 auto mines', () => {
		const planet = new CommandedPlanet();
		const player = new Player();
		player.techLevels = {
			energy: 3,
			weapons: 3,
			propulsion: 3,
			construction: 3,
			electronics: 3,
			biotechnology: 3
		};

		// will terraform 1
		planet.baseHab = { grav: 49, temp: 50, rad: 50 };
		planet.hab = { grav: 49, temp: 50, rad: 50 };

		// we have enough minerals to build factories, but not enough resources
		planet.mineralConcentration = {
			ironium: 100,
			boranium: 100,
			germanium: 100
		};
		planet.cargo = {
			ironium: 0,
			boranium: 0,
			germanium: 100
		};
		planet.population = 100_000; // 100_000 colonists generate 100 resources per year

		planet.productionQueue = [
			{ type: QueueItemTypes.AutoMaxTerraform, quantity: 5 },
			{ type: QueueItemTypes.AutoFactories, quantity: 20 },
			{ type: QueueItemTypes.AutoMines, quantity: 100 }
		];

		const designFinder = new TestDesignFinder();

		// we build a mine each year, but we have to wait for our planet to grow to build
		// the factory, which costs 10 resources and we keep auto building the mine each year
		const items = getProductionEstimates(defaultRules, techStore, player, planet, designFinder);

		// should return a new production queue with estimates in it
		expect(items).toMatchObject([
			{
				type: QueueItemTypes.AutoMaxTerraform,
				yearsToBuildOne: 1,
				yearsToBuildAll: 100,
				yearsToSkipAuto: 1
			},
			{
				type: QueueItemTypes.AutoFactories,
				yearsToBuildOne: 2,
				yearsToBuildAll: 8,
				yearsToSkipAuto: 4
			},
			{
				type: QueueItemTypes.AutoMines,
				yearsToBuildOne: 3,
				yearsToBuildAll: 100,
				yearsToSkipAuto: 14
			}
		]);
	});

	it('estimate 20 Auto Factories, 20 auto mines, 20 auto max terraform', () => {
		const planet = new CommandedPlanet();
		const player = new Player();
		player.techLevels = {
			energy: 3,
			weapons: 3,
			propulsion: 3,
			construction: 3,
			electronics: 3,
			biotechnology: 3
		};

		// perfect planet
		planet.baseHab = { grav: 50, temp: 50, rad: 50 };
		planet.hab = { grav: 50, temp: 50, rad: 50 };

		// we have enough minerals to build factories, but not enough resources
		planet.mineralConcentration = {
			ironium: 100,
			boranium: 100,
			germanium: 100
		};
		planet.cargo = {
			ironium: 1000,
			boranium: 1000,
			germanium: 1000
		};
		planet.population = 100_000; // 100_000 colonists generate 100 resources per year

		planet.factories = 100;
		planet.mines = 20;
		planet.productionQueue = [
			{ type: QueueItemTypes.AutoFactories, quantity: 100 },
			{ type: QueueItemTypes.AutoMines, quantity: 100 },
			{ type: QueueItemTypes.AutoMaxTerraform, quantity: 10 }
		];

		const designFinder = new TestDesignFinder();

		// we build a mine each year, but we have to wait for our planet to grow to build
		// the factory, which costs 10 resources and we keep auto building the mine each year
		const items = getProductionEstimates(defaultRules, techStore, player, planet, designFinder);

		// should return a new production queue with estimates in it
		expect(items).toMatchObject([
			{
				type: QueueItemTypes.AutoFactories,
				yearsToBuildOne: 2,
				yearsToBuildAll: 100,
				yearsToSkipAuto: 1
			},
			{
				type: QueueItemTypes.AutoMines,
				yearsToBuildOne: 1,
				yearsToBuildAll: 100,
				yearsToSkipAuto: 13
			},
			{
				type: QueueItemTypes.AutoMaxTerraform,
				yearsToBuildOne: 100,
				yearsToBuildAll: 100,
				yearsToSkipAuto: 13
			}
		]);
	});
});
