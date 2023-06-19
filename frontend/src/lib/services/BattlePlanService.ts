import type { Fleet } from '$lib/types/Fleet';
import { Service } from './Service';
import type { BattlePlan, PlayerResponse } from '$lib/types/Player';
import { CSError, type ErrorResponse } from './Errors';

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
	): Promise<{ player: PlayerResponse; fleets: Fleet[]; starbases: Fleet[] }> {
		const url = `/api/games/${gameId}/battle-plans/${num}`;
		const response = await fetch(url, {
			method: 'DELETE',
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			const result = (await response.json()) as ErrorResponse;
			throw new CSError(result);
		}
		return (await response.json()) as {
			player: PlayerResponse;
			fleets: Fleet[];
			starbases: Fleet[];
		};
	}
}
