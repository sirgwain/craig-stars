import type { Game } from '$lib/types/Game';
import type { GameContext } from '$lib/types/GameContext';
import { Service } from './Service';

export class GameService extends Service {
	async loadPlayerGames(): Promise<Game[]> {
		return this.get<Game[]>('/api/games');
	}

	async loadHostedGames(): Promise<Game[]> {
		return this.get<Game[]>('/api/games/hosted');
	}

	async loadGame(gameId: number): Promise<GameContext> {
		const response = await fetch(`/api/games/${gameId}`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			return (await response.json()) as GameContext;
		} else {
			console.log(response);
			throw new Error('Failed to load game');
		}
	}
}
