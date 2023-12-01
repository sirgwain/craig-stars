import { Player } from '$lib/types/Player';
import { describe, expect, it } from 'vitest';
import { getTerraformAmount } from './Terraformer';
import type { TechStore } from '$lib/types/Tech';
import techjson from '$lib/ssr/techs.json';

describe('Terraformer test', () => {
	const techStore = techjson as TechStore;

	it('getTerraformAmount - no ability', () => {
		const player = new Player();

		expect(
			getTerraformAmount(
				techStore,
				{ grav: 50, temp: 50, rad: 50 },
				{ grav: 50, temp: 50, rad: 50 },
				player
			)
		).toEqual({ grav: 0, temp: 0, rad: 0 });

		expect(
			getTerraformAmount(
				techStore,
				{ grav: 47, temp: 50, rad: 50 },
				{ grav: 47, temp: 50, rad: 50 },
				player
			)
		).toEqual({ grav: 0, temp: 0, rad: 0 });
	});

	it('getTerraformAmount - 3 ability', () => {
		const player = new Player();
		player.techLevels = {
			energy: 3,
			weapons: 3,
			propulsion: 3,
			construction: 3,
			electronics: 3,
			biotechnology: 3
		};

		expect(
			getTerraformAmount(
				techStore,
				{ grav: 50, temp: 50, rad: 50 },
				{ grav: 50, temp: 50, rad: 50 },
				player
			)
		).toEqual({ grav: 0, temp: 0, rad: 0 });

		expect(
			getTerraformAmount(
				techStore,
				{ grav: 47, temp: 48, rad: 49 },
				{ grav: 47, temp: 48, rad: 49 },
				player
			)
		).toEqual({ grav: 3, temp: 2, rad: 1 });

		expect(
			getTerraformAmount(
				techStore,
				{ grav: 53, temp: 52, rad: 51 },
				{ grav: 53, temp: 52, rad: 51 },
				player
			)
		).toEqual({ grav: -3, temp: -2, rad: -1 });

		expect(
			getTerraformAmount(
				techStore,
				{ grav: 53, temp: 48, rad: 51 },
				{ grav: 53, temp: 48, rad: 51 },
				player
			)
		).toEqual({ grav: -3, temp: 2, rad: -1 });
	});

	it('getTerraformAmount - 3 ability, already partially terraformed', () => {
		const player = new Player();
		player.techLevels = {
			energy: 3,
			weapons: 3,
			propulsion: 3,
			construction: 3,
			electronics: 3,
			biotechnology: 3
		};

		expect(
			getTerraformAmount(
				techStore,
				{ grav: 48, temp: 50, rad: 50 }, // hab
				{ grav: 47, temp: 50, rad: 50 }, // baseHab
				player
			)
		).toEqual({ grav: 2, temp: 0, rad: 0 });

		expect(
			getTerraformAmount(
				techStore,
				{ grav: 48, temp: 53, rad: 50 }, // hab
				{ grav: 47, temp: 54, rad: 50 }, // baseHab
				player
			)
		).toEqual({ grav: 2, temp: -2, rad: 0 });
	});
});
