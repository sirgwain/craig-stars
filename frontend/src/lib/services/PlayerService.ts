import type { Game } from '$lib/types/Game';
import type { Planet } from '$lib/types/Planet';
import type { PlayerIntels, PlayerUniverse, PlayerOrders, PlayerResponse } from '$lib/types/Player';
import { CSError } from './Errors';
import { Service } from './Service';

type UpdateOrdersResult = {
	player: PlayerResponse;
	planets: Planet[];
};

type SubmitTurnResponse = {
	game: Game;
	player?: PlayerResponse;
	universe?: PlayerUniverse & PlayerIntels;
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

	static async updatePlans(player: PlayerResponse): Promise<PlayerResponse | undefined> {
		const response = await fetch(`/api/games/${player.gameId}/player/plans`, {
			method: 'PUT',
			body: JSON.stringify(player),
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			return (await response.json()) as PlayerResponse;
		} else {
			throw new CSError(await response.json());
		}
	}

	static async submitTurn(gameId: number | string): Promise<SubmitTurnResponse | undefined> {
		const response = await fetch(`/api/games/${gameId}/submit-turn`, {
			method: 'POST',
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			return response.json();
		} else {
			console.error(response);
		}
	}
}
