import type { PlayerResponse, TransportPlan } from '$lib/types/Player';
import { CSError, type ErrorResponse } from './Errors';
import { Service } from './Service';

export class TransportPlanService {
	static async update(
		gameId: number | string,
		transportPlan: TransportPlan
	): Promise<TransportPlan> {
		return Service.update(
			transportPlan,
			`/api/games/${gameId}/transport-plans/${transportPlan.num}`
		);
	}

	static async create(
		gameId: number | string,
		transportPlan: TransportPlan
	): Promise<TransportPlan> {
		return Service.create(transportPlan, `/api/games/${gameId}/transport-plans`);
	}

	static async delete(gameId: number | string, num: number | string): Promise<PlayerResponse> {
		const url = `/api/games/${gameId}/transport-plans/${num}`;
		const response = await fetch(url, {
			method: 'DELETE',
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		return (await response.json()) as PlayerResponse;
	}
}
