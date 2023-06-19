import type { Game } from '$lib/types/Game';
import { type PlayerResponse, type PlayerMapObjects, Player } from '$lib/types/Player';
import { Service } from './Service';

type playerStatusResult = {
	players: PlayerResponse[];
};

export class GameService {
	static async loadPlayerGames(): Promise<Game[]> {
		return Service.get<Game[]>('/api/games');
	}

	static async loadHostedGames(): Promise<Game[]> {
		return Service.get<Game[]>('/api/games/hosted');
	}

	static async loadOpenGames(): Promise<Game[]> {
		return Service.get<Game[]>('/api/games/open');
	}

	static async deleteGame(gameId: number): Promise<any> {
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

	static async loadGame(gameId: number): Promise<Game> {
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

	static async loadLightPlayer(gameId: number): Promise<Player> {
		const response = await fetch(`/api/games/${gameId}/player`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			const json = (await response.json()) as PlayerResponse;
			return new Player(json.gameId, json.num, json);
		} else {
			throw new Error('Failed to load game');
		}
	}

	static async loadFullPlayer(gameId: number | string): Promise<Player> {
		const response = await fetch(`/api/games/${gameId}/full-player`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			const json = (await response.json()) as PlayerResponse;
			return new Player(json.gameId, json.num, json);
		} else {
			throw new Error('Failed to load game');
		}
	}

	static async loadPlayerMapObjects(gameId: number): Promise<PlayerMapObjects> {
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

	static async loadPlayerStatuses(gameId: number): Promise<Player[]> {
		const response = await fetch(`/api/games/${gameId}/player-statuses`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			const json = ((await response.json()) as playerStatusResult).players;
			return json.map((data) => new Player(data.gameId, data.num, data));
		} else {
			throw new Error('Failed to load game');
		}
	}
}
