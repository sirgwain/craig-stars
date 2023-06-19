import type { Game } from '$lib/types/Game';
import type { GameContext } from '$lib/types/GameContext';
import type { Player, PlayerMapObjects } from '$lib/types/Player';
import { Service } from './Service';

type playerStatusResult = {
	players: Player[];
};

export class GameService extends Service {
	async loadPlayerGames(): Promise<Game[]> {
		return Service.get<Game[]>('/api/games');
	}

	async loadHostedGames(): Promise<Game[]> {
		return Service.get<Game[]>('/api/games/hosted');
	}

	async loadOpenGames(): Promise<Game[]> {
		return Service.get<Game[]>('/api/games/open');
	}

	async deleteGame(gameId: number): Promise<any> {
		const response = await fetch(`/api/games/${gameId}`, {
			method: 'DELETE',
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			const error = await response.json();
			throw new Error('Failed to delete game', error);
		}
	}

	async loadGame(gameId: number): Promise<Game> {
		const response = await fetch(`/api/games/${gameId}`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			return (await response.json()) as Game;
		} else {
			throw new Error('Failed to load game');
		}
	}

	async loadLightPlayer(gameId: number): Promise<Player> {
		const response = await fetch(`/api/games/${gameId}/player`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			return (await response.json()) as Player;
		} else {
			throw new Error('Failed to load game');
		}
	}

	async loadFullPlayer(gameId: number): Promise<Player> {
		const response = await fetch(`/api/games/${gameId}/full-player`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			return (await response.json()) as Player;
		} else {
			throw new Error('Failed to load game');
		}
	}

	async loadPlayerMapObjects(gameId: number): Promise<PlayerMapObjects> {
		const response = await fetch(`/api/games/${gameId}/mapobjects`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			return (await response.json()) as PlayerMapObjects;
		} else {
			throw new Error('Failed to load game');
		}
	}

	async loadOpenGame(gameId: number): Promise<Game> {
		const response = await fetch(`/api/games/open/${gameId}`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			return (await response.json()) as Game;
		} else {
			throw new Error('Failed to load game');
		}
	}

	async loadPlayerStatuses(gameId: number): Promise<Player[]> {
		const response = await fetch(`/api/games/${gameId}/player-statuses`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			return ((await response.json()) as playerStatusResult).players;
		} else {
			throw new Error('Failed to load game');
		}
	}
}
