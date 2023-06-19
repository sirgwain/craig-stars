import type { Cargo } from '$lib/types/Cargo';
import { CommandedFleet, type Fleet, type Waypoint } from '$lib/types/Fleet';
import type { MapObject } from '$lib/types/MapObject';
import type { Planet } from '$lib/types/Planet';
import type { ErrorResponse } from './Errors';
import { Service } from './Service';

// orders sent to the server
export class FleetOrders {
	constructor(private waypoints: Waypoint[], private repeatOrders: boolean = false) {}
}

export class FleetService {
	static async load(gameId: number): Promise<Fleet[]> {
		return Service.get(`/api/games/${gameId}/fleets`);
	}

	static async get(gameId: number | string, num: number | string): Promise<CommandedFleet> {
		const fleet = await Service.get<Fleet>(`/api/games/${gameId}/fleets/${num}`);
		const commandedFleet = new CommandedFleet();
		return Object.assign(commandedFleet, fleet);
	}

	static async update(gameId: number | string, fleet: CommandedFleet): Promise<CommandedFleet> {
		const updated = Service.update(fleet, `/api/games/${gameId}/fleets/${fleet.num}`);
		return Object.assign(fleet, updated);
	}

	static async transferCargo(
		fleet: CommandedFleet,
		dest: Fleet | Planet | undefined,
		transferAmount: Cargo
	): Promise<Fleet> {
		const url = `/api/games/${fleet.gameId}/fleets/${fleet.num}/transfer-cargo`;
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				accept: 'application/json'
			},
			body: JSON.stringify({
				mo: dest as MapObject,
				transferAmount: transferAmount
			})
		});

		if (response.ok) {
			return await response.json();
		} else {
			const errorResponse = (await response.json()) as ErrorResponse;
			throw new Error(errorResponse.error);
		}
	}

	static async splitAll(gameId: number | string, fleet: Fleet): Promise<Fleet[]> {
		const url = `/api/games/${gameId}/fleets/${fleet.num}/split-all`;
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			return (await response.json()) as Fleet[];
		} else {
			const errorResponse = (await response.json()) as ErrorResponse;
			throw new Error(errorResponse.error);
		}
	}

	static async merge(fleet: CommandedFleet, fleetNums: number[]): Promise<CommandedFleet> {
		const url = `/api/games/${fleet.gameId}/fleets/${fleet.num}/merge`;
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				accept: 'application/json'
			},
			body: JSON.stringify({ fleetNums })
		});

		if (response.ok) {
			return Object.assign(fleet, await response.json());
		} else {
			const errorResponse = (await response.json()) as ErrorResponse;
			throw new Error(errorResponse.error);
		}
	}

	static async updateFleetOrders(fleet: CommandedFleet): Promise<Fleet> {
		const fleetOrders = new FleetOrders(fleet.waypoints ?? [], fleet.repeatOrders);

		const response = await fetch(`/api/games/${fleet.gameId}/fleets/${fleet.num}`, {
			method: 'PUT',
			headers: {
				accept: 'application/json'
			},
			body: JSON.stringify(fleetOrders)
		});

		if (response.ok) {
			return await response.json();
		} else {
			console.error(response);
		}
		return Promise.resolve(fleet);
	}
}
