import { describe, expect, it } from 'vitest';
import { StargateWarpSpeed } from './Consts';
import {
	canTransferCargoOrFuel,
	canTransferCargoType,
	CommandedFleet,
	getTokenDamage,
	moveDamagedTokens
} from './Fleet';
import {
	FleetSchema,
	MapObjectType,
	ResourceType,
	SalvageSchema,
	ShipDesignSchema,
	ShipDesignSpecSchema,
	ShipTokenSchema,
	type Planet,
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

describe('canTransferCargoOrFuel test', () => {
	const scout = create(FleetSchema, { mapObject: { playerNum: 1, type: MapObjectType.FLEET } });
	const freighter = create(FleetSchema, {
		mapObject: { playerNum: 1, type: MapObjectType.FLEET },
		spec: { shipDesignSpec: { cargoCapacity: 70 } }
	});

	it('always opens for a fleet with a cargo hold', () => {
		expect(canTransferCargoOrFuel(freighter, undefined)).toBe(true);
	});

	it('opens for a fleet with no hold only when it can trade fuel with our own fleet', () => {
		const ours = create(FleetSchema, { mapObject: { playerNum: 1, type: MapObjectType.FLEET } });
		const theirs = create(FleetSchema, { mapObject: { playerNum: 2, type: MapObjectType.FLEET } });
		const salvage = create(SalvageSchema, { mapObject: { type: MapObjectType.SALVAGE } });
		expect(canTransferCargoOrFuel(scout, ours)).toBe(true);
		expect(canTransferCargoOrFuel(scout, theirs)).toBe(false);
		expect(canTransferCargoOrFuel(scout, salvage)).toBe(false);
		expect(canTransferCargoOrFuel(scout, undefined)).toBe(false);
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

describe('trip to a waypoint', () => {
	const planet = (num: number, dockCapacity: number) =>
		({
			mapObject: { num, playerNum: 1 },
			spec: { planetStarbaseSpec: { dockCapacity } }
		}) as Planet;
	const planets = [planet(1, 0), planet(2, -1)];
	const universe = { getPlanet: (num: number) => planets.find((p) => p.mapObject?.num === num) };
	const player = { isFriend: (num: number) => num === 1 };
	const at = (x: number, warpSpeed: number, estFuelUsage: number, planetNum = 0) => ({
		position: { x, y: 0 },
		warpSpeed,
		estFuelUsage,
		mapObjectTarget: planetNum
			? { targetType: MapObjectType.PLANET, targetNum: planetNum }
			: undefined
	});
	// out through deep space to a planet with no dock, on to a starbase, then out again
	const fleet = new CommandedFleet(
		create(FleetSchema, {
			mapObject: { type: MapObjectType.FLEET, playerNum: 1, num: 1 },
			fleetOrders: {
				waypoints: [
					at(0, 5, 0),
					at(100, 5, 40),
					at(125, 5, 10, 1),
					at(200, 6, 60, 2),
					at(210, 9, 5)
				]
			}
		})
	);
	type Args = Parameters<typeof fleet.getFuelUsageToWaypoint>;
	const fuelTo = (index: number) =>
		fleet.getFuelUsageToWaypoint(
			player as unknown as Args[0],
			universe as unknown as Args[1],
			index
		);

	it('adds up fuel leg by leg', () => {
		expect(fuelTo(0)).toBe(0);
		expect(fuelTo(1)).toBe(40);
		expect(fuelTo(2)).toBe(50);
	});

	it('still shows the fuel it takes to reach a starbase', () => {
		expect(fuelTo(3)).toBe(110);
	});

	it('starts over after refueling, keeping the most needed at once', () => {
		// 5mg since the starbase is less than the 110mg it took to get there
		expect(fuelTo(4)).toBe(110);
	});

	it('knows the fuel left on arrival, filling up at starbases on the way', () => {
		const fuelAt = (fuel: number, capacity: number, index: number) => {
			const f = new CommandedFleet(fleet);
			f.fuel = fuel;
			f.spec = { shipDesignSpec: { fuelCapacity: capacity } } as typeof f.spec;
			return f.getFuelOnArrival(
				player as unknown as Args[0],
				universe as unknown as Args[1],
				index
			);
		};
		expect(fuelAt(120, 200, 0)).toBe(120);
		expect(fuelAt(120, 200, 1)).toBe(80);
		expect(fuelAt(120, 200, 2)).toBe(70);
		// arrives at the starbase with 10mg, then leaves it with a full tank
		expect(fuelAt(120, 200, 3)).toBe(10);
		expect(fuelAt(120, 200, 4)).toBe(195);
		// not enough for the third leg, so it never reaches the starbase to fill up
		expect(fuelAt(100, 200, 2)).toBe(50);
		expect(fuelAt(100, 200, 3)).toBe(-10);
		expect(fuelAt(100, 200, 4)).toBe(-10);
	});

	it('adds up the years of every leg', () => {
		expect(fleet.getYearsToWaypoint(0)).toBe(0);
		expect(fleet.getYearsToWaypoint(1)).toBe(4); // 100 ly at warp 5
		expect(fleet.getYearsToWaypoint(2)).toBe(5); // 25 ly more
		expect(fleet.getYearsToWaypoint(3)).toBe(8); // 75 ly at warp 6 is 3 years
		expect(fleet.getYearsToWaypoint(4)).toBe(9); // a short hop still takes a year
	});

	it('uses a stargate in a year and never arrives at warp 0', () => {
		const gate = new CommandedFleet(
			create(FleetSchema, {
				fleetOrders: { waypoints: [at(0, 5, 0), at(500, StargateWarpSpeed, 0), at(600, 0, 0)] }
			})
		);
		expect(gate.getYearsToWaypoint(1)).toBe(1);
		expect(gate.getYearsToWaypoint(2)).toBeUndefined();
	});
});
