import type { User } from '$lib/types/cs';
import { Service } from './Service';

export class UserService {
	static async get(id: number | string): Promise<User> {
		return Service.get<User>(`/api/users/${id}`);
	}
}
