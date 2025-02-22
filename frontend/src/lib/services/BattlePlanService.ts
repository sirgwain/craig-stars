import type { Fleet } from '$lib/types/cs';
import type { Player } from '$lib/types/cs';
import type { BattlePlan } from '$lib/types/cs';
import { Service } from './Service';

export class BattlePlanService {
	static async update(gameId: number | string, battlePlan: BattlePlan): Promise<BattlePlan> {
		return Service.update(battlePlan, `/api/games/${gameId}/battle-plans/${battlePlan.num}`);
	}

	static async create(gameId: number | string, battlePlan: BattlePlan): Promise<BattlePlan> {
		return Service.create(battlePlan, `/api/games/${gameId}/battle-plans`);
	}

	static async delete(
		gameId: number | string,
		num: number | string
	): Promise<{ player: Player; fleets: Fleet[]; starbases: Fleet[] }> {
		const url = `/api/games/${gameId}/battle-plans/${num}`;
		const response = await fetch(url, {
			method: 'DELETE',
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		return (await response.json()) as {
			player: Player;
			fleets: Fleet[];
			starbases: Fleet[];
		};
	}
}
