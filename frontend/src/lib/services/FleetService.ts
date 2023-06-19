import type { Cargo } from '$lib/types/Cargo';
import type { ErrorResponse } from '$lib/types/ErrorResponse';
import type { Fleet, Waypoint } from '$lib/types/Fleet';
import type { MapObject } from '$lib/types/MapObject';
import type { Planet } from '$lib/types/Planet';
import { Service } from './Service';

// orders sent to the server
export class FleetOrders {
	constructor(private waypoints: Waypoint[], private repeatOrders: boolean = false) {}
}

export class FleetService extends Service {
	static async transferCargo(
		gameId: number | string,
		fleet: Fleet,
		dest: Planet | Fleet | undefined,
		transferAmount: Cargo
	): Promise<Fleet> {
		const url = `/api/games/${gameId}/fleets/${fleet.num}/transfer-cargo`;
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
			return (await response.json()) as Fleet;
		} else {
			const errorResponse = (await response.json()) as ErrorResponse;
			throw new Error(errorResponse.error);
		}
	}

	static async updateFleetOrders(fleet: Fleet): Promise<Fleet> {
		const fleetOrders = new FleetOrders(fleet.waypoints ?? [], fleet.repeatOrders);

		const response = await fetch(`/api/fleets/${fleet.id}`, {
			method: 'PUT',
			headers: {
				accept: 'application/json'
			},
			body: JSON.stringify(fleetOrders)
		});

		if (response.ok) {
			return (await response.json()) as Fleet;
		} else {
			console.error(response);
		}
		return Promise.resolve(fleet);
	}
}
