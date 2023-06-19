import type { GameContext } from '$lib/types/GameContext';
import type { Planet } from '$lib/types/Planet';
import type { Player } from '$lib/types/Player';
import { Service } from './Service';

export class PlayerService extends Service {
	constructor(private player: Player) {
		super();
	}

	async submitTurn(): Promise<GameContext | undefined> {
		const response = await fetch(`/api/games/${this.player.gameId}/submit-turn`, {
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
