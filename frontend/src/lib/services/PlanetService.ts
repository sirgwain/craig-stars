import { CommandedPlanet, type Planet, type PlanetOrders } from '$lib/types/Planet';
import type { PlayerResponse } from '$lib/types/Player';
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

	static async updatePlanetOrders(
		planet: CommandedPlanet
	): Promise<{ planet: Planet; player: PlayerResponse }> {
		const planetOrders: PlanetOrders = {
			contributesOnlyLeftoverToResearch: planet.contributesOnlyLeftoverToResearch,
			routeTargetType: planet.routeTargetType,
			routeTargetNum: planet.routeTargetNum,
			routeTargetPlayerNum: planet.routeTargetPlayerNum,
			packetSpeed: planet.packetSpeed,
			packetTargetNum: planet.packetTargetNum,
			productionQueue: planet.productionQueue
		};

		const response = await fetch(`/api/games/${planet.gameId}/planets/${planet.num}`, {
			method: 'PUT',
			headers: {
				accept: 'application/json'
			},
			body: JSON.stringify(planetOrders)
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		return await response.json();
	}
}
