import type { Race } from '$lib/types/Race';
import { Service } from './Service';

export class RaceService extends Service {
	async loadRaces(): Promise<Race[]> {
		return this.get<Race[]>('/api/races');
	}

	async loadRace(raceId: number): Promise<Race> {
		return this.get<Race>(`/api/races/${raceId}`);
	}
}
