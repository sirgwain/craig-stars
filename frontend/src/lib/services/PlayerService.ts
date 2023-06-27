import type { Planet } from '$lib/types/Planet';
import type { PlayerOrders, PlayerResponse } from '$lib/types/Player';
import type { TurnGenerationResponse } from './GameService';
import { Service } from './Service';

type UpdateOrdersResult = {
	player: PlayerResponse;
	planets: Planet[];
};

export class PlayerService extends Service {
	static async updateOrders(player: PlayerResponse): Promise<UpdateOrdersResult | undefined> {
		const orders: PlayerOrders = {
			researching: player.researching,
			nextResearchField: player.nextResearchField,
			researchAmount: player.researchAmount
		};
		const response = await fetch(`/api/games/${player.gameId}/player`, {
			method: 'PUT',
			body: JSON.stringify(orders),
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		return (await response.json()) as UpdateOrdersResult;
	}

	static async updatePlans(player: PlayerResponse): Promise<PlayerResponse | undefined> {
		const response = await fetch(`/api/games/${player.gameId}/player/plans`, {
			method: 'PUT',
			body: JSON.stringify(player),
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		return (await response.json()) as PlayerResponse;
	}

	static async submitTurn(gameId: number | string): Promise<TurnGenerationResponse | undefined> {
		const response = await fetch(`/api/games/${gameId}/submit-turn`, {
			method: 'POST',
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		return response.json();
	}

	static async unsubmitTurn(gameId: number | string): Promise<undefined> {
		const response = await fetch(`/api/games/${gameId}/unsubmit-turn`, {
			method: 'POST',
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		return response.json();
	}
}
