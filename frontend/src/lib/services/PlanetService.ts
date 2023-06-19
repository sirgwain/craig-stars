import { CommandedPlanet, type Planet } from '$lib/types/Planet';
import { Service } from './Service';

export class PlanetService {
	static async load(gameId: number): Promise<Planet[]> {
		return Service.get(`/api/games/${gameId}/planets`);
	}

	static async get(gameId: number | string, num: number | string): Promise<CommandedPlanet> {
		const planet = await Service.get<Planet>(`/api/games/${gameId}/planets/${num}`);
		const commandedPlanet = new CommandedPlanet();
		return Object.assign(commandedPlanet, planet);
	}

	static async update(gameId: number | string, planet: CommandedPlanet): Promise<CommandedPlanet> {
		const updated = Service.update(planet, `/api/games/${gameId}/planets/${planet.num}`);
		return Object.assign(planet, updated);
	}
}
