import type { CargoDest, CargoTransferRequest } from '$lib/types/CargoTransferRequest.svelte';
import { CommandedFleet } from '$lib/types/Fleet';
import type {
	Cargo,
	CargoTransfers,
	MapObject,
	MineralPacketIntel,
	Player,
	SalvageIntel
} from '$lib/types/cs';
import { type Fleet, type FleetOrders, type ShipToken, type Waypoint } from '$lib/types/cs';
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
	player: Player | undefined;
	salvages?: SalvageIntel[];
	mineralPackets?: MineralPacketIntel[];
};

type SplitFleetResponse = {
	source: Fleet;
	dest?: Fleet;
	cargoTransfers: CargoTransfers;
};

type SplitAllResponse = {
	fleets: Fleet[];
	cargoTransfers: CargoTransfers;
};

type MergeResponse = {
	fleet: Fleet;
	cargoTransfers: CargoTransfers;
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
		dest: CargoDest,
		transferAmount: Cargo & { fuel: number }
	): Promise<TransferCargoResponse> {
		const url = `/api/games/${fleet.gameId}/fleets/${fleet.num}/transfer-cargo`;
		const body = JSON.stringify({
			mo: dest as MapObject,
			transferAmount: transferAmount
		});
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				accept: 'application/json'
			},
			body: body
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		return (await response.json()) as TransferCargoResponse;
	}

	static async split(
		source: CommandedFleet,
		dest: Fleet | undefined,
		sourceTokens: ShipToken[],
		destTokens: ShipToken[],
		transferAmount: CargoTransferRequest
	): Promise<SplitFleetResponse> {
		const url = `/api/games/${source.gameId}/fleets/${source.num}/split`;
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				accept: 'application/json'
			},
			body: JSON.stringify({
				sourceFleetNum: source.num,
				destFleetNum: dest?.num,
				destBaseName: dest?.baseName,
				sourceTokens,
				destTokens,
				transferAmount: transferAmount.jsonData()
			})
		});

		if (!response.ok) {
			await Service.throwError(response);
		}

		return await response.json();
	}

	static async splitAll(gameId: number | string, fleet: Fleet): Promise<SplitAllResponse> {
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
		return (await response.json()) as SplitAllResponse;
	}

	static async merge(fleet: CommandedFleet, fleetNums: number[]): Promise<MergeResponse> {
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

		return (await response.json()) as MergeResponse;
	}

	static async updateFleetOrders(fleet: CommandedFleet): Promise<Fleet> {
		const fleetOrders = new FleetOrdersRequest(
			fleet.waypoints ?? [],
			fleet.repeatOrders,
			fleet.battlePlanNum
		);

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
