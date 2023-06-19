import type { Game } from '$lib/types/Game';
import type { Planet } from '$lib/types/Planet';
import { type PlayerResponse, type PlayerOrders, Player } from '$lib/types/Player';
import { Service } from './Service';

type UpdateOrdersResult = {
	player: PlayerResponse;
	planets: Planet[];
};

type SubmitTurnResponse = {
	game: Game;
	player: PlayerResponse;
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

		if (response.ok) {
			return (await response.json()) as UpdateOrdersResult;
		} else {
			console.error(response);
		}
	}

	static async submitTurn(
		player: PlayerResponse
	): Promise<{ game: Game; player: Player } | undefined> {
		const response = await fetch(`/api/games/${player.gameId}/submit-turn`, {
			method: 'POST',
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			const json = (await response.json()) as SubmitTurnResponse;
			const player = new Player(json.player.gameId, json.player.num, json.player);
			return { game: json.game, player };
		} else {
			console.error(response);
		}
	}
}
