import type { Fleet } from '$lib/types/Fleet';
import type { ShipDesign, Spec } from '$lib/types/ShipDesign';
import { Service } from './Service';

export class DesignService {
	static async load(gameId: number): Promise<ShipDesign[]> {
		return Service.get(`/api/games/${gameId}/designs`);
	}

	static async get(gameId: number | string, num: number | string): Promise<ShipDesign> {
		return Service.get(`/api/games/${gameId}/designs/${num}`);
	}

	static async update(gameId: number | string, design: ShipDesign): Promise<ShipDesign> {
		return Service.update(design, `/api/games/${gameId}/designs/${design.num}`);
	}

	/**
	 * compute a spec for a design
	 * @param gameId The gameId for the design
	 * @param design The design to compute a spec for
	 * @returns the newly computed spec
	 */
	static async computeSpec(gameId: number | string, design: ShipDesign): Promise<Spec> {
		const url = `/api/games/${gameId}/designs/spec`;
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				accept: 'application/json'
			},
			body: JSON.stringify(design)
		});

		if (!response.ok) {
			const result = await response.json();
			if ('error' in result) {
				console.error(`Failed to compute spec ${url}`, result);
			}
			throw new Error(`Failed to compute spec ${url} ${JSON.stringify(result)}`);
		}

		// update the spec
		return (await response.json()) as Spec;
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
			const result = await response.json();
			if ('error' in result) {
				console.error(`Failed to delete ${url}`, result);
			}
			throw new Error(`Failed to delete ${url} ${JSON.stringify(result)}`);
		}
		return (await response.json()) as { fleets: Fleet[]; starbases: Fleet[] };
	}
}
