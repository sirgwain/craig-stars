import type { GameWithPlayers } from '$lib/types/cs';
import { UserSession } from '$lib/types/User';
import { Service } from './Service';

export class AdminService {
	static async loadGames(): Promise<GameWithPlayers[]> {
		return Service.get<GameWithPlayers[]>('/api/admin/games');
	}

	static async loadUsers(): Promise<UserSession[]> {
		const response = await Service.get<UserSession[]>('/api/admin/users');
		return response.map((su) => Object.assign(new UserSession(), su));
	}

	static async loadUserGames(userId: number | string): Promise<GameWithPlayers[]> {
		return Service.get<GameWithPlayers[]>(`/api/admin/users/${userId}/games`);
	}
}
