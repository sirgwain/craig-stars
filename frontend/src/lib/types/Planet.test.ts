import techjson from '$lib/ssr/techs.json';
import { describe, expect, it } from 'vitest';
import { CommandedPlanet } from './Planet';
import { Player } from './Player';
import { QueueItemTypes } from './QueueItemType';
import { humanoid } from './Race';
import { defaultRules } from './Rules';
import type { TechStore } from './Tech';

describe('Planet test', () => {
	const techStore = techjson as TechStore;

	it('getMaxPopulation', () => {
		const planet = new CommandedPlanet();
		const player = new Player();

		planet.hab = { grav: 50, temp: 50, rad: 50 };

		expect(planet.getMaxPopulation(defaultRules, player, 100)).toBe(1_200_000);
	});

	it('getGrowthAmount', () => {
		const planet = new CommandedPlanet();
		const player = new Player();
		const race = player.race; // defaults to humanoid
		race.growthRate = 10; // 10% growth

		planet.hab = { grav: 50, temp: 50, rad: 50 };
		planet.population = 100_000;

		// should be 1_200_000
		const maxPopulation = planet.getMaxPopulation(defaultRules, player, 100);

		// should grow 10%
		expect(
			planet.getGrowthAmount(
				race,
				maxPopulation,
				defaultRules.populationOvercrowdDieoffRate,
				defaultRules.populationOvercrowdDieoffRateMax
			)
		).toBe(10000);
	});

	it('getProductivePopulation', () => {
		const planet = new CommandedPlanet();

		planet.population = 100_000;

		expect(planet.getProductivePopulation(1_000_000)).toBe(100_000);
		expect(planet.getProductivePopulation(10_000)).toBe(30_000);
	});

	it('getInnateMines', () => {
		const planet = new CommandedPlanet();
		const race = humanoid();

		expect(planet.getInnateMines(race, 16000)).toBe(0);

		if (race.spec) {
			race.spec.innateMining = true;
			race.spec.innatePopulationFactor = 0.1;
		}
		expect(planet.getInnateMines(race, 16000)).toBe(12);
	});

	it('getMaxMines', () => {
		const planet = new CommandedPlanet();
		const race = humanoid();

		expect(planet.getMaxMines(race, 10_000)).toBe(10);
		expect(planet.getMaxMines(race, 100_000)).toBe(100);
	});

	it('getMaxFactories', () => {
		const planet = new CommandedPlanet();
		const race = humanoid();

		expect(planet.getMaxFactories(race, 10_000)).toBe(10);
		expect(planet.getMaxFactories(race, 100_000)).toBe(100);
	});

	it('getResourcesPerYear', () => {
		const planet = new CommandedPlanet();
		const player = new Player();
		planet.hab = { grav: 50, temp: 50, rad: 50 };
		planet.population = 10_000;
		expect(planet.getResourcesAvailable(player)).toBe(10);

		// 1 resource per factory
		planet.factories = 1;
		expect(planet.getResourcesAvailable(player)).toBe(11);
	});

	it('getMaxBuildable', () => {
		const planet = new CommandedPlanet();
		const player = new Player();

		planet.hab = { grav: 50, temp: 50, rad: 50 };
		planet.mines = 10;
		planet.factories = 10;
		planet.defenses = 10;
		planet.population = 100_000;

		expect(planet.getMaxBuildable(techStore, player, 1_000_000, QueueItemTypes.Mine)).toBe(990);
		expect(planet.getMaxBuildable(techStore, player, 1_000_000, QueueItemTypes.AutoMines)).toBe(90);
		expect(planet.getMaxBuildable(techStore, player, 1_000_000, QueueItemTypes.Factory)).toBe(990);
		expect(planet.getMaxBuildable(techStore, player, 1_000_000, QueueItemTypes.AutoFactories)).toBe(
			90
		);
		expect(planet.getMaxBuildable(techStore, player, 1_000_000, QueueItemTypes.Defenses)).toBe(90);
		expect(planet.getMaxBuildable(techStore, player, 1_000_000, QueueItemTypes.AutoDefenses)).toBe(
			90
		);

		// should build a scanner
		expect(planet.getMaxBuildable(techStore, player, 1, QueueItemTypes.PlanetaryScanner)).toBe(1);
		planet.scanner = true;
		expect(planet.getMaxBuildable(techStore, player, 1, QueueItemTypes.PlanetaryScanner)).toBe(0);
	});

	it('grows', () => {
		const planet = new CommandedPlanet();
		const player = new Player();

		planet.hab = { grav: 50, temp: 50, rad: 50 };
		planet.population = 100_000;

		planet.grow(defaultRules, player);
		expect(planet.population).toBe(115_000);
	});

	it('mines', () => {
		const planet = new CommandedPlanet();
		const player = new Player();

		planet.population = 100_000;
		planet.mines = 10;
		planet.mineralConcentration = { ironium: 100, boranium: 100, germanium: 100 };

		planet.mine(defaultRules, player.race);
		expect(planet.cargo).toEqual({ ironium: 10, boranium: 10, germanium: 10, colonists: 1_000 });
	});

	it('reduceMineralConcentration', () => {
		const planet = new CommandedPlanet();

		planet.mines = 150;
		planet.mineralConcentration = { ironium: 100, boranium: 100, germanium: 100 };
		planet.mineYears = { ironium: 151, boranium: 151, germanium: 151 };

		planet.reduceMineralConcentration(defaultRules);
		expect(planet.mineralConcentration).toEqual({ ironium: 99, boranium: 99, germanium: 99 });
	});
});
