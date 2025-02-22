import type { GameWithPlayers } from '$lib/types/cs';
import { User, type SessionUser } from '$lib/types/User';
import { Service } from './Service';

export class AdminService {
	static async loadGames(): Promise<GameWithPlayers[]> {
		return Service.get<GameWithPlayers[]>('/api/admin/games');
	}

	static async loadUsers(): Promise<User[]> {
		const response = await Service.get<SessionUser[]>('/api/admin/users');
		return response.map((su) => Object.assign(new User(), su));
	}

	static async loadUserGames(userId: number | string): Promise<GameWithPlayers[]> {
		return Service.get<GameWithPlayers[]>(`/api/admin/users/${userId}/games`);
	}
}
