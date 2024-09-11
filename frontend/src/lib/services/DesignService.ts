import type { Fleet } from '$lib/types/Fleet';
import type { ShipDesign, Spec } from '$lib/types/ShipDesign';
import { Service } from './Service';

export class DesignService {
	static async load(gameId: number | string): Promise<ShipDesign[]> {
		return Service.get(`/api/games/${gameId}/designs`);
	}

	static async get(gameId: number | string, num: number | string): Promise<ShipDesign> {
		return Service.get(`/api/games/${gameId}/designs/${num}`);
	}

	static async update(gameId: number | string, design: ShipDesign): Promise<ShipDesign> {
		return Service.update(design, `/api/games/${gameId}/designs/${design.num}`);
	}

	static async create(gameId: number | string, design: ShipDesign): Promise<ShipDesign> {
		return Service.create(design, `/api/games/${gameId}/designs`);
	}

	static async delete(
		gameId: number | string,
		num: number | string
	): Promise<{ fleets: Fleet[]; starbases: Fleet[] }> {
		const url = `/api/games/${gameId}/designs/${num}`;
		const response = await fetch(url, {
			method: 'DELETE',
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		return (await response.json()) as { fleets: Fleet[]; starbases: Fleet[] };
	}
}
