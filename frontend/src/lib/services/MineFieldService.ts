import { type Minefield, type MinefieldOrders } from '$lib/types/cs';
import { Service } from './Service';

// orders sent to the server
export class MinefieldOrdersRequest implements MinefieldOrders {
	constructor(public detonate: boolean) {}
}

export class MinefieldService {
	static async updateMinefieldOrders(minefield: Minefield): Promise<Minefield> {
		const minefieldOrders = new MinefieldOrdersRequest(minefield.detonate ?? false);

		const response = await fetch(`/api/games/${minefield.gameId}/mine-fields/${minefield.num}`, {
			method: 'PUT',
			headers: {
				accept: 'application/json'
			},
			body: JSON.stringify(minefieldOrders)
		});

		if (!response.ok) {
			await Service.throwError(response);
		}

		return await response.json();
	}
}
