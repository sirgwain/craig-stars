import { type MineField, type MineFieldOrders } from '$lib/types/MineField';
import { Service } from './Service';

// orders sent to the server
export class MineFieldOrdersRequest implements MineFieldOrders {
	constructor(public detonate: boolean) {}
}

export class MineFieldService {
	static async updateMineFieldOrders(mineField: MineField): Promise<MineField> {
		const mineFieldOrders = new MineFieldOrdersRequest(mineField.detonate ?? false);

		const response = await fetch(`/api/games/${mineField.gameId}/mine-fields/${mineField.num}`, {
			method: 'PUT',
			headers: {
				accept: 'application/json'
			},
			body: JSON.stringify(mineFieldOrders)
		});

		if (!response.ok) {
			await Service.throwError(response);
		}

		return await response.json();
	}
}
