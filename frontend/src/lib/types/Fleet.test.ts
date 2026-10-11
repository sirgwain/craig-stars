import { describe, expect, it } from 'vitest';
import { canTransferCargoType, getTokenDamage, moveDamagedTokens } from './Fleet';
import {
	FleetSchema,
	MapObjectType,
	ResourceType,
	SalvageSchema,
	ShipDesignSchema,
	ShipDesignSpecSchema,
	ShipTokenSchema,
	type ShipTokenJson
} from '#lib/types/cs-proto.js';
import { create } from '@bufbuild/protobuf';

describe('ShipToken moveDamagedTokens test', () => {
	it('transfer no damaged tokens', () => {
		const srcToken: ShipTokenJson = { designNum: 1, quantity: 1 };
		const destToken: ShipTokenJson = { designNum: 1, quantity: 1 };
		moveDamagedTokens(srcToken, destToken, 1);
		expect(srcToken).toEqual({ designNum: 1, quantity: 1, quantityDamaged: 0, damage: 0 });
		expect(destToken).toEqual({ designNum: 1, quantity: 1, quantityDamaged: 0, damage: 0 });
	});
	it('transfer 1 damaged token into undamaged stack', () => {
		const srcToken: ShipTokenJson = { designNum: 1, quantity: 1, quantityDamaged: 1, damage: 10 };
		const destToken: ShipTokenJson = { designNum: 1, quantity: 1 };
		moveDamagedTokens(srcToken, destToken, 1);
		expect(srcToken).toEqual({ designNum: 1, quantity: 1, quantityDamaged: 0, damage: 0 });
		expect(destToken).toEqual({ designNum: 1, quantity: 1, quantityDamaged: 1, damage: 10 });
	});
	it('transfer 1 damaged token into damaged stack', () => {
		const srcToken: ShipTokenJson = { designNum: 1, quantity: 1, quantityDamaged: 1, damage: 10 };
		const destToken: ShipTokenJson = { designNum: 1, quantity: 1, quantityDamaged: 1, damage: 5 };
		moveDamagedTokens(srcToken, destToken, 1);
		expect(srcToken).toEqual({ designNum: 1, quantity: 1, quantityDamaged: 0, damage: 0 });
		expect(destToken).toEqual({ designNum: 1, quantity: 1, quantityDamaged: 2, damage: 7.5 });
	});
});

describe('canTransferCargoType test', () => {
	const fleet = create(FleetSchema, { mapObject: { playerNum: 1, type: MapObjectType.FLEET } });

	it('does not jettison colonists or fuel into deep space', () => {
		expect(canTransferCargoType(fleet, undefined, ResourceType.COLONISTS)).toBe(false);
		expect(canTransferCargoType(fleet, undefined, ResourceType.FUEL)).toBe(false);
		expect(canTransferCargoType(fleet, undefined, ResourceType.IRONIUM)).toBe(true);
	});

	it('does not put colonists in salvage', () => {
		const salvage = create(SalvageSchema, { mapObject: { type: MapObjectType.SALVAGE } });
		expect(canTransferCargoType(fleet, salvage, ResourceType.COLONISTS)).toBe(false);
		expect(canTransferCargoType(fleet, salvage, ResourceType.GERMANIUM)).toBe(true);
	});
});

describe('getTokenDamage test', () => {
	const design = create(ShipDesignSchema, { spec: create(ShipDesignSpecSchema, { armor: 500 }) });

	it('undamaged token has no damage', () => {
		expect(getTokenDamage(create(ShipTokenSchema, { quantity: 1 }), design)).toBeUndefined();
	});
	it('missing token or design has no damage', () => {
		expect(getTokenDamage(undefined, design)).toBeUndefined();
		expect(
			getTokenDamage(
				create(ShipTokenSchema, { quantity: 1, quantityDamaged: 1, damage: 100 }),
				undefined
			)
		).toBeUndefined();
	});
	it("damage is a percent of one ship's armor", () => {
		expect(
			getTokenDamage(
				create(ShipTokenSchema, { quantity: 5, quantityDamaged: 2, damage: 125 }),
				design
			)
		).toEqual({ quantity: 2, percent: 25 });
	});
	it('any damage shows as at least 1%', () => {
		expect(
			getTokenDamage(
				create(ShipTokenSchema, { quantity: 1, quantityDamaged: 1, damage: 1 }),
				design
			)
		).toEqual({ quantity: 1, percent: 1 });
	});
	it('damaged quantity is capped at the token quantity', () => {
		expect(
			getTokenDamage(
				create(ShipTokenSchema, { quantity: 1, quantityDamaged: 3, damage: 250 }),
				design
			)
		).toEqual({ quantity: 1, percent: 50 });
	});
});
