import type { Race } from '$lib/types/Race';
import { Service } from './Service';

export class RaceService {

	static async load(): Promise<Race[]> {
		return Service.get<Race[]>('/api/races');
	}

	static async get(id: number | string): Promise<Race> {
		return Service.get<Race>(`/api/races/${id}`);
	}

	static async delete(race: Race): Promise<void> {
		return Service.delete(`/api/races/${race.id}`);
	}
}
