import type { Cargo } from '$lib/types/Cargo';
import { CommandedFleet, type Fleet, type FleetOrders, type Waypoint } from '$lib/types/Fleet';
import type { MapObject } from '$lib/types/MapObject';
import type { Planet } from '$lib/types/Planet';
import type { PlayerRelation, PlayerResponse } from '$lib/types/Player';
import type { Salvage } from '$lib/types/Salvage';
import { Service } from './Service';

// orders sent to the server
export class FleetOrdersRequest implements FleetOrders {
	constructor(
		public waypoints: Waypoint[],
		public repeatOrders: boolean = false,
		public battlePlanNum: number = 1
	) {}
}

type TransferCargoResponse = {
	fleet: Fleet;
	dest: MapObject | undefined;
	salvages?: Salvage[];
};

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

	static async rename(fleet: CommandedFleet, name: string): Promise<CommandedFleet> {
		// rename the fleet and update it
		const updated = await Service.post<{ name: string }, Fleet>(
			{ name },
			`/api/games/${fleet.gameId}/fleets/${fleet.num}/rename`
		);
		return Object.assign(fleet, updated);
	}

	static async transferCargo(
		fleet: CommandedFleet,
		dest: Fleet | Planet | Salvage,
		transferAmount: Cargo
	): Promise<TransferCargoResponse> {
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

		if (!response.ok) {
			await Service.throwError(response);
		}
		return (await response.json()) as TransferCargoResponse;
	}

	static async splitAll(gameId: number | string, fleet: Fleet): Promise<Fleet[]> {
		const url = `/api/games/${gameId}/fleets/${fleet.num}/split-all`;
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		return (await response.json()) as Fleet[];
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

		if (!response.ok) {
			await Service.throwError(response);
		}

		return Object.assign(fleet, await response.json());
	}

	static async updateFleetOrders(fleet: CommandedFleet): Promise<Fleet> {
		const fleetOrders = new FleetOrdersRequest(fleet.waypoints ?? [], fleet.repeatOrders, fleet.battlePlanNum);

		const response = await fetch(`/api/games/${fleet.gameId}/fleets/${fleet.num}`, {
			method: 'PUT',
			headers: {
				accept: 'application/json'
			},
			body: JSON.stringify(fleetOrders)
		});

		if (!response.ok) {
			await Service.throwError(response);
		}

		return await response.json();
	}
}
