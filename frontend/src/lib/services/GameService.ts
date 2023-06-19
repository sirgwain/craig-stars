import type { GameContext } from '$lib/types/GameContext';
import type { Player } from '$lib/types/Player';
import { Service } from './Service';

export class GameService extends Service {
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
