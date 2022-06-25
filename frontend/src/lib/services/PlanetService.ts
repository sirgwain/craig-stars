import type { Planet } from '$lib/types/Planet';
import { Service } from './Service';

export class PlanetService extends Service {
	async updatePlanet(planet: Planet): Promise<Planet> {
		return this.update<Planet>(planet, `/api/planets/${planet.gameId}`);
	}
}
