import { describe, expect, it } from 'vitest';
import { CommandedFleet, moveDamagedTokens, type Fleet, type ShipToken } from './Fleet';
import { Player } from './Player';
import type { MapObject } from './MapObject';

describe('Fleet test', () => {
	const fleet = new CommandedFleet({
		spec: {
			engine: {
				maxSafeSpeed: 9
			}
		}
	} as Fleet);
	it('returns minimal speeds for distances', () => {
		const player = new Player();
		const origin = { position: { x: 0, y: 0 } } as MapObject;
		// one year to go 49 ly
		expect(
			fleet.getMinimalWarp(player, origin, { position: { x: 49, y: 0 } } as MapObject, 0, 7)
		).toBe(7);

		// two years at warp 7, two at warp 6 or warp 5, pick warp 5
		expect(
			fleet.getMinimalWarp(player, origin, { position: { x: 50, y: 0 } } as MapObject, 0, 7)
		).toBe(5);

		// warp 6 takes 3 years to 72, so in 73 we can make it in 2 at warp 7
		expect(
			fleet.getMinimalWarp(player, origin, { position: { x: 73, y: 0 } } as MapObject, 0, 7)
		).toBe(7);

		// // this is obvious
		expect(
			fleet.getMinimalWarp(player, origin, { position: { x: 36, y: 0 } } as MapObject, 0, 7)
		).toBe(6);

		// // might as well go warp 5
		expect(
			fleet.getMinimalWarp(player, origin, { position: { x: 25, y: 0 } } as MapObject, 0, 7)
		).toBe(5);
	});
});

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
