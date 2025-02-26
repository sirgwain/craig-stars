import {
	GameStateSetup,
	type Game,
	type GameSettings,
	type GameWithPlayers,
	type Player,
	type PlayerIntels,
	type PlayerStatus
} from '$lib/types/cs';
import { CommandedPlayer } from '$lib/types/Player';
import type { SessionUser } from '$lib/types/User';
import { FullGame } from './FullGame';
import { Service } from './Service';
import type { PlayerUniverse } from './Universe';

export type TurnGenerationResponse = {
	game: Game;
	player?: Player;
	universe?: PlayerUniverse & PlayerIntels;
};

type UniverseResponse = PlayerUniverse & PlayerIntels;

export class GameService {
	static async updateSettings(id: number, settings: GameSettings): Promise<GameSettings> {
		return Service.update(settings, `/api/games/${id}`);
	}

	static async addOpenPlayerSlot(id: number): Promise<Game> {
		return Service.post(undefined, `/api/games/${id}/add-open-player-slot`);
	}

	static async addGuestPlayer(id: number): Promise<Game> {
		return Service.post(undefined, `/api/games/${id}/add-guest-player`);
	}

	static async addAIPlayer(id: number): Promise<Game> {
		return Service.post(undefined, `/api/games/${id}/add-ai-player`);
	}

	static async kickPlayer(id: number, playerNum: number): Promise<Game> {
		return Service.post({ playerNum }, `/api/games/${id}/kick-player`);
	}

	static async deletePlayer(id: number, playerNum: number): Promise<Game> {
		return Service.post({ playerNum }, `/api/games/${id}/delete-player`);
	}

	static async updatePlayer(id: number, player: PlayerStatus): Promise<Game> {
		return Service.post(player, `/api/games/${id}/update-player`);
	}

	static async loadPlayerGames(): Promise<GameWithPlayers[]> {
		return Service.get<GameWithPlayers[]>('/api/games');
	}

	static async loadHostedGames(): Promise<GameWithPlayers[]> {
		return Service.get<GameWithPlayers[]>('/api/games/hosted');
	}

	static async loadOpenGames(): Promise<GameWithPlayers[]> {
		return Service.get<GameWithPlayers[]>('/api/games/open');
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

	static async loadGame(gameId: number | string): Promise<GameWithPlayers> {
		const response = await fetch(`/api/games/${gameId}`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		return (await response.json()) as GameWithPlayers;
	}

	static async loadGuest(gameId: number | string, playerNum: number): Promise<SessionUser> {
		const response = await fetch(`/api/games/${gameId}/guest/${playerNum}`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		return await response.json();
	}

	static async loadGameByHash(hash: string): Promise<GameWithPlayers[]> {
		const response = await fetch(`/api/games/invite/${hash}`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		return (await response.json()) as GameWithPlayers[];
	}

	static async loadLightPlayer(gameId: number): Promise<CommandedPlayer> {
		const response = await fetch(`/api/games/${gameId}/player`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		const json = (await response.json()) as Player;
		return new CommandedPlayer(json);
	}

	static async loadFullPlayer(gameId: number | string): Promise<CommandedPlayer> {
		const response = await fetch(`/api/games/${gameId}/full-player`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		const json = (await response.json()) as Player;
		return new CommandedPlayer(json);
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

	// load in a full game with universe, techs, and player data
	static async loadFullGame(gameId: number): Promise<FullGame> {
		const id = parseInt(gameId.toString());
		const game = await GameService.loadGame(id);
		const fg = Object.assign(new FullGame(), game);
		if (fg.state != GameStateSetup) {
			await Promise.all([
				GameService.loadFullPlayer(id).then((data) => {
					fg.player = data;
				}),
				GameService.loadUniverse(id).then((u) => {
					fg.universe.setData(u);
				}),
				// load techs the first time as well
				fg.techs.fetch()
			]);
		}

		// configure the universe for the player after the player is loaded
		fg.universe.setPlayer(fg.player.num);
		return fg;
	}
}
