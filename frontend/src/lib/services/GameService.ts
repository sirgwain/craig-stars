import type { Game } from '$lib/types/Game';
import {
	Player,
	type PlayerIntels, type PlayerResponse,
	type PlayerUniverse
} from '$lib/types/Player';
import { CSError } from './Errors';
import { Service } from './Service';

type playerStatusResult = {
	game: Game;
	players: PlayerResponse[];
};

type UniverseResponse = PlayerUniverse & PlayerIntels;

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

	static async deleteGame(gameId: number): Promise<void> {
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

	static async loadGame(gameId: number | string): Promise<Game> {
		const response = await fetch(`/api/games/${gameId}`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			return (await response.json()) as Game;
		} else {
			throw new CSError(response.statusText, response.status);
		}
	}

	static async loadGameByHash(hash: string): Promise<Game[]> {
		const response = await fetch(`/api/games/invite/${hash}`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			return (await response.json()) as Game[];
		} else {
			throw new CSError(response.statusText, response.status);
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
			return new Player(json);
		} else {
			throw new CSError(response.statusText, response.status);
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
			return new Player(json);
		} else {
			throw new CSError(response.statusText, response.status);
		}
	}

	static async loadUniverse(gameId: number | string): Promise<UniverseResponse> {
		const response = await fetch(`/api/games/${gameId}/universe`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			return (await response.json()) as UniverseResponse;
		} else {
			throw new CSError(response.statusText, response.status);
		}
	}

	static async loadPlayersStatus(gameId: number): Promise<playerStatusResult> {
		const response = await fetch(`/api/games/${gameId}/players-status`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			return (await response.json()) as playerStatusResult;
		} else {
			throw new CSError(response.statusText, response.status);
		}
	}
}
