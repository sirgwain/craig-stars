import { describe, expect, it } from 'vitest';
import { moveDamagedTokens } from './Fleet';
import { type ShipToken } from './cs';

describe('ShipToken moveDamagedTokens test', () => {
	it('transfer no damaged tokens', () => {
		const srcToken: ShipToken = { designNum: 1, quantity: 1 };
		const destToken: ShipToken = { designNum: 1, quantity: 1 };
		moveDamagedTokens(srcToken, destToken, 1);
		expect(srcToken).toEqual({ designNum: 1, quantity: 1, quantityDamaged: 0, damage: 0 });
		expect(destToken).toEqual({ designNum: 1, quantity: 1, quantityDamaged: 0, damage: 0 });
	});
	it('transfer 1 damaged token into undamaged stack', () => {
		const srcToken: ShipToken = { designNum: 1, quantity: 1, quantityDamaged: 1, damage: 10 };
		const destToken: ShipToken = { designNum: 1, quantity: 1 };
		moveDamagedTokens(srcToken, destToken, 1);
		expect(srcToken).toEqual({ designNum: 1, quantity: 1, quantityDamaged: 0, damage: 0 });
		expect(destToken).toEqual({ designNum: 1, quantity: 1, quantityDamaged: 1, damage: 10 });
	});
	it('transfer 1 damaged token into damaged stack', () => {
		const srcToken: ShipToken = { designNum: 1, quantity: 1, quantityDamaged: 1, damage: 10 };
		const destToken: ShipToken = { designNum: 1, quantity: 1, quantityDamaged: 1, damage: 5 };
		moveDamagedTokens(srcToken, destToken, 1);
		expect(srcToken).toEqual({ designNum: 1, quantity: 1, quantityDamaged: 0, damage: 0 });
		expect(destToken).toEqual({ designNum: 1, quantity: 1, quantityDamaged: 2, damage: 7.5 });
	});
});
