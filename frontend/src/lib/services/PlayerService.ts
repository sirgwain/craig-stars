import type { GameContext } from '$lib/types/GameContext';
import type { Planet } from '$lib/types/Planet';
import type { Player, PlayerOrders } from '$lib/types/Player';
import { Service } from './Service';

type UpdateOrdersResult = {
	player: Player;
	planets: Planet[];
};
export class PlayerService extends Service {
	static async updateOrders(player: Player): Promise<UpdateOrdersResult | undefined> {
		const orders: PlayerOrders = {
			researching: player.researching,
			nextResearchField: player.nextResearchField,
			researchAmount: player.researchAmount
		};
		const response = await fetch(`/api/games/${player.gameId}`, {
			method: 'PUT',
			body: JSON.stringify(orders),
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			return (await response.json()) as UpdateOrdersResult;
		} else {
			console.error(response);
		}
	}

	static async submitTurn(player: Player): Promise<GameContext | undefined> {
		const response = await fetch(`/api/games/${player.gameId}/submit-turn`, {
			method: 'POST',
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			return (await response.json()) as GameContext;
		} else {
			console.error(response);
		}
	}
}
