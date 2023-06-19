import type { Race } from '$lib/types/Race';
import { Service } from './Service';

export class RaceService extends Service {
	static async load(): Promise<Race[]> {
		return Service.get<Race[]>('/api/races');
	}

	static async delete(id: number): Promise<void> {
		return Service.delete(id, '/api/races');
	}
}
