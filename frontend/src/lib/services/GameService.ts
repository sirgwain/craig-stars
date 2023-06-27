import type { Game, GameSettings, NewGamePlayers } from '$lib/types/Game';
import {
	Player,
	type PlayerIntels,
	type PlayerResponse,
	type PlayerUniverse
} from '$lib/types/Player';
import { Service } from './Service';

export type TurnGenerationResponse = {
	game: Game;
	player?: PlayerResponse;
	universe?: PlayerUniverse & PlayerIntels;
};

type UniverseResponse = PlayerUniverse & PlayerIntels;

export class GameService {
	static async updateSettings(
		id: number,
		settings: GameSettings,
		players?: NewGamePlayers
	): Promise<GameSettings> {
		return Service.update(settings, `/api/games/${id}`);
	}

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
			await Service.throwError(response);
		}
	}

	static async loadGame(gameId: number | string): Promise<Game> {
		const response = await fetch(`/api/games/${gameId}`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		return (await response.json()) as Game;
	}

	static async loadGameByHash(hash: string): Promise<Game[]> {
		const response = await fetch(`/api/games/invite/${hash}`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		return (await response.json()) as Game[];
	}

	static async loadLightPlayer(gameId: number): Promise<Player> {
		const response = await fetch(`/api/games/${gameId}/player`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		const json = (await response.json()) as PlayerResponse;
		return new Player(json);
	}

	static async loadFullPlayer(gameId: number | string): Promise<Player> {
		const response = await fetch(`/api/games/${gameId}/full-player`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		const json = (await response.json()) as PlayerResponse;
		return new Player(json);
	}

	static async loadUniverse(gameId: number | string): Promise<UniverseResponse> {
		const response = await fetch(`/api/games/${gameId}/universe`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		return (await response.json()) as UniverseResponse;
	}

	static async forceGenerateTurn(gameId: number): Promise<TurnGenerationResponse> {
		const response = await fetch(`/api/games/${gameId}/generate-turn`, {
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
}
