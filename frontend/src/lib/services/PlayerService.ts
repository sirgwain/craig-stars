import type { Game } from '$lib/types/cs';
import type { Planet } from '$lib/types/cs';
import type { Player } from '$lib/types/cs';
import type { PlayerOrders } from '$lib/types/cs';
import type { TechLevel } from '$lib/types/cs';
import type { TurnGenerationResponse } from './GameService';
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

	static async updatePlans(player: Player): Promise<Player | undefined> {
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
		return (await response.json()) as Player;
	}

	static async updateRelations(player: Player): Promise<Player | undefined> {
		const response = await fetch(`/api/games/${player.gameId}/player/relations`, {
			method: 'PUT',
			body: JSON.stringify(player),
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		return (await response.json()) as Player;
	}

	static async archiveGame(gameId: number): Promise<Game> {
		const response = await fetch(`/api/games/${gameId}/archive-game`, {
			method: 'POST',
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			await Service.throwError(response);
		}

		return await response.json();
	}

	static async unArchiveGame(gameId: number): Promise<Game> {
		const response = await fetch(`/api/games/${gameId}/unarchive-game`, {
			method: 'POST',
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			await Service.throwError(response);
		}

		return await response.json();
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

	static async getResearchCost(
		gameId: number,
		techLevel: TechLevel
	): Promise<{ resources: number }> {
		const response = await fetch(`/api/games/${gameId}/research-cost`, {
			method: 'POST',
			headers: {
				accept: 'application/json'
			},
			body: JSON.stringify(techLevel)
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		return await response.json();
	}
}
