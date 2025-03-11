import techjson from '$lib/ssr/techs.json';
import { describe, expect, it } from 'vitest';
import { CommandedPlanet } from './Planet';
import { CommandedPlayer } from './Player';
import { humanoid } from './Race';
import { defaultRules } from './Rules';
import {
	QueueItemTypeAutoDefenses,
	QueueItemTypeAutoFactories,
	QueueItemTypeAutoMines,
	QueueItemTypeDefenses,
	QueueItemTypeFactory,
	QueueItemTypeMine,
	QueueItemTypePlanetaryScanner,
	type TechStore
} from './cs';

describe('Planet test', () => {
	const techStore = techjson as unknown as TechStore;

	it('getMaxPopulation', () => {
		const planet = new CommandedPlanet();
		const player = new CommandedPlayer();

		planet.hab = { grav: 50, temp: 50, rad: 50 };

		expect(planet.getMaxPopulation(defaultRules, player, 100)).toBe(1_200_000);
	});

	it('getInnateMines', () => {
		const planet = new CommandedPlanet();
		const race = humanoid();

		expect(planet.getInnateMines(race, 16000)).toBe(0);

		if (race.spec) {
			race.spec.innateMining = true;
			race.spec.innateScannerFactor = 0.1;
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
	})

	it('getMaxBuildable', () => {
		const planet = new CommandedPlanet();
		const player = new CommandedPlayer();

		planet.hab = { grav: 50, temp: 50, rad: 50 };
		planet.mines = 10;
		planet.factories = 10;
		planet.defenses = 10;
		planet.population = 100_000;

		expect(planet.getMaxBuildable(techStore, player, 1_000_000, QueueItemTypeMine)).toBe(990);
		expect(planet.getMaxBuildable(techStore, player, 1_000_000, QueueItemTypeAutoMines)).toBe(90);
		expect(planet.getMaxBuildable(techStore, player, 1_000_000, QueueItemTypeFactory)).toBe(990);
		expect(planet.getMaxBuildable(techStore, player, 1_000_000, QueueItemTypeAutoFactories)).toBe(
			90
		);
		expect(planet.getMaxBuildable(techStore, player, 1_000_000, QueueItemTypeDefenses)).toBe(90);
		expect(planet.getMaxBuildable(techStore, player, 1_000_000, QueueItemTypeAutoDefenses)).toBe(
			90
		);

		// should build a scanner
		expect(planet.getMaxBuildable(techStore, player, 1, QueueItemTypePlanetaryScanner)).toBe(1);
		planet.scanner = true;
		expect(planet.getMaxBuildable(techStore, player, 1, QueueItemTypePlanetaryScanner)).toBe(0);
	});


});
