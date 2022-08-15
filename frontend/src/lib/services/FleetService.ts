import type { Cargo } from '$lib/types/Cargo';
import type { ErrorResponse } from '$lib/types/ErrorResponse';
import type { Fleet } from '$lib/types/Fleet';
import type { MapObject } from '$lib/types/MapObject';
import type { Planet } from '$lib/types/Planet';
import { Service } from './Service';

export class FleetService extends Service {
	async transferCargo(
		fleet: Fleet,
		dest: Planet | Fleet | undefined,
		transferAmount: Cargo
	): Promise<Fleet> {
		const url = `/api/fleets/${fleet.id}/transfer-cargo`;
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

	async updateFleet(fleet: Fleet): Promise<Fleet> {
		return this.update<Fleet>(fleet, `/api/fleets/${fleet.gameId}`);
	}
}
