import type { Player } from '$lib/types/cs';
import type { ProductionPlan } from '$lib/types/cs';
import { Service } from './Service';

export class ProductionPlanService {
	static async update(
		gameId: number | string,
		productionPlan: ProductionPlan
	): Promise<ProductionPlan> {
		return Service.update(
			productionPlan,
			`/api/games/${gameId}/production-plans/${productionPlan.num}`
		);
	}

	static async create(
		gameId: number | string,
		productionPlan: ProductionPlan
	): Promise<ProductionPlan> {
		return Service.create(productionPlan, `/api/games/${gameId}/production-plans`);
	}

	static async delete(gameId: number | string, num: number | string): Promise<Player> {
		const url = `/api/games/${gameId}/production-plans/${num}`;
		const response = await fetch(url, {
			method: 'DELETE',
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		return (await response.json()) as Player;
	}
}
