import type { Game } from '$lib/types/Game';
import { User, type SessionUser } from '$lib/types/User';
import { Service } from './Service';

export class AdminService {
	static async loadGames(): Promise<Game[]> {
		return Service.get<Game[]>('/api/admin/games');
	}

	static async loadUsers(): Promise<User[]> {
		const response = await Service.get<SessionUser[]>('/api/admin/users');
		return response.map((su) => Object.assign(new User(), su));
	}
}
