import type { Cost } from '$lib/types/Cost';
import { CommandedPlanet, type Planet, type PlanetOrders } from '$lib/types/Planet';
import type { Player } from '$lib/types/Player';
import type { ShipDesign } from '$lib/types/ShipDesign';
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

	static async getPlanetProductionEstimates(planet: CommandedPlanet, player: Player): Promise<CommandedPlanet> {

		const response = await fetch(`/api/calculators/planet-production-estimate`, {
			method: 'POST',
			headers: {
				accept: 'application/json'
			},
			body: JSON.stringify({planet, player})
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		const updated = await response.json();
		return Object.assign(new CommandedPlanet(), updated);
	}

	static async getStarbaseUpgradeCost(design: ShipDesign, newDesign: ShipDesign): Promise<Cost> {
		const response = await fetch(`/api/calculators/starbase-upgrade-cost`, {
			method: 'POST',
			headers: {
				accept: 'application/json'
			},
			body: JSON.stringify({ design, newDesign })
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		return await response.json();
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
			await Service.throwError(response);
		}
		return await response.json();
	}
}
