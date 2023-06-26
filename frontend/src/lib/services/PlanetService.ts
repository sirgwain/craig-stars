import { CommandedPlanet, type Planet, type PlanetOrders } from '$lib/types/Planet';
import { CSError, type ErrorResponse } from './Errors';
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

	static async getPlanetProductionEstimates(planet: CommandedPlanet): Promise<CommandedPlanet> {
		const response = await fetch(`/api/calculators/planet-production-estimate`, {
			method: 'POST',
			headers: {
				accept: 'application/json'
			},
			body: JSON.stringify(planet)
		});

		if (!response.ok) {
			await Service.raiseError(response);
		}
		const updated = await response.json();
		return Object.assign(new CommandedPlanet(), updated);
	}

	static async updatePlanetOrders(planet: CommandedPlanet): Promise<Planet> {
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
			await Service.raiseError(response);
		}
		return await response.json();
	}
}
