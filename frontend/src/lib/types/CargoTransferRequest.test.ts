import { describe, expect, it } from 'vitest';
import {
	clampCargoTransfer,
	clampFuelTransfer,
	suggestCargoTransfer,
	suggestFuelTransfer
} from './CargoTransferRequest';

describe('suggestFuelTransfer', () => {
	it('clamps fuel transfer', () => {
		// source has 200/200, dest has 0/200
		// suggest -100 to even them out
		expect(suggestFuelTransfer(200, 200, 0, 200)).toEqual(-100);

		// source has 50/200, dest has 0/200
		// suggest -25 to even them out
		expect(suggestFuelTransfer(50, 200, 0, 200)).toEqual(-25);

		// source has 50/200, dest has 100/200
		// suggest 25 to even them out
		expect(suggestFuelTransfer(50, 200, 100, 200)).toEqual(25);

		// source has 300/200, dest has 0/200
		// suggest -150 to even them out
		expect(suggestFuelTransfer(300, 200, 0, 200)).toEqual(-150);
	});
});

describe('clampFuelTransfer', () => {
	it('clamps fuel transfer', () => {
		// transfer 200 from source to dest
		// source has 200/200, dest has 0/200
		expect(clampFuelTransfer(-200, 200, 200, 0, 200)).toEqual(-200);

		// transfer 200 from source to dest
		// source has 200/200, dest has 10/200
		expect(clampFuelTransfer(-200, 200, 200, 10, 200)).toEqual(-190);

		// transfer 50 from source to dest
		// source has 150/200, dest has 150/200
		expect(clampFuelTransfer(-50, 150, 200, 150, 200)).toEqual(-50);

		// transfer 75 from source to dest
		// source has 150/200, dest has 150/200, can only give -50
		expect(clampFuelTransfer(-75, 150, 200, 150, 200)).toEqual(-50);

		// transfer 75 from dest to source
		// source has 150/200, dest has 150/200, can only take 50
		expect(clampFuelTransfer(75, 150, 200, 150, 200)).toEqual(50);

		// transfer 50 from dest to source
		// source has 150/200, dest has 150/200, can only take 50
		expect(clampFuelTransfer(75, 150, 200, 150, 200)).toEqual(50);
	});
});

describe('suggestCargoTransfer', () => {
	it('suggests per-lane deltas to balance by capacity', () => {
		// equal capacities: src has 200 iron, dest has 0 iron => suggest move 100 to dest => delta -100
		expect(
			suggestCargoTransfer(
				{ ironium: 200, boranium: 0, germanium: 0, colonists: 0 },
				200,
				{ ironium: 0, boranium: 0, germanium: 0, colonists: 0 },
				200
			)
		).toEqual({ ironium: -100, boranium: 0, germanium: 0, colonists: 0 });

		// equal capacities: src has 50 iron, dest has 0 iron => suggest -25
		expect(
			suggestCargoTransfer(
				{ ironium: 50, boranium: 0, germanium: 0, colonists: 0 },
				200,
				{ ironium: 0, boranium: 0, germanium: 0, colonists: 0 },
				200
			)
		).toEqual({ ironium: -25, boranium: 0, germanium: 0, colonists: 0 });

		// equal capacities: src has 50 iron, dest has 100 iron => total 150 => desired dest 75 => delta = 100-75 = +25
		expect(
			suggestCargoTransfer(
				{ ironium: 50, boranium: 0, germanium: 0, colonists: 0 },
				200,
				{ ironium: 100, boranium: 0, germanium: 0, colonists: 0 },
				200
			)
		).toEqual({ ironium: 25, boranium: 0, germanium: 0, colonists: 0 });
	});

	it('respects shared capacity via clampCargoTransfer', () => {
		// dest only has 10 free total capacity; suggestion would try to move 100 iron,
		// but clamp should reduce to -10.
		expect(
			suggestCargoTransfer(
				{ ironium: 200, boranium: 0, germanium: 0, colonists: 0 },
				200,
				{ ironium: 190, boranium: 0, germanium: 0, colonists: 0 },
				200
			)
		).toEqual({ ironium: -5, boranium: 0, germanium: 0, colonists: 0 });
	});
});

describe('clampCargoTransfer', () => {
	it('clamps per-component and total capacity', () => {
		const src = { ironium: 200, boranium: 0, germanium: 0, colonists: 0 };
		const dest = { ironium: 10, boranium: 0, germanium: 0, colonists: 0 };

		// try to send 200 ironium (src->dest), but dest only has 190 free total space (cap 200, total 10)
		expect(
			clampCargoTransfer(
				{ ironium: -200, boranium: 0, germanium: 0, colonists: 0 },
				src,
				200,
				dest,
				200
			)
		).toEqual({ ironium: -190, boranium: 0, germanium: 0, colonists: 0 });

		// mixed components: try to send 100 iron + 150 colonists (sum -250), but dest only has 200 free total
		expect(
			clampCargoTransfer(
				{ ironium: -100, boranium: 0, germanium: 0, colonists: -150 },
				{ ironium: 100, boranium: 0, germanium: 0, colonists: 150 },
				300,
				{ ironium: 0, boranium: 0, germanium: 0, colonists: 0 },
				200
			)
		).toEqual({ ironium: -100, boranium: 0, germanium: 0, colonists: -100 }); // pulled back toward 0
	});
});
