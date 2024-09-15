import { describe, expect, it } from 'vitest';
import { CommandedFleet, moveDamagedTokens, type ShipToken } from './Fleet';
import { None } from './MapObject';
import { cottonPicker, longRangeScout, santaMaria, TestDesignFinder } from './Mock.test';
import { CommandedPlanet } from './Planet';
import { Player, PlayerRelation } from './Player';
import type { RaceSpec } from './Race';

describe('Fleet test', () => {
	it('getFuelUsage', () => {
		const fleet = new CommandedFleet(longRangeScout);
		const designFinder = new TestDesignFinder();
		// 5mg fuel for one year at warp6
		expect(fleet.getFuelCost(designFinder, 0, 6, 36, 0)).toBe(5);

		// 338mg fuel for 300 ly at warp 9
		expect(fleet.getFuelCost(designFinder, 0, 7, 300, 0)).toBe(169);

		// 338mg fuel for 300 ly at warp 9
		expect(fleet.getFuelCost(designFinder, 0, 8, 300, 0)).toBe(282);

		// 338mg fuel for 300 ly at warp 9
		expect(fleet.getFuelCost(designFinder, 0, 9, 300, 0)).toBe(338);
	});

	it('tests canColonize', () => {
		const scout = new CommandedFleet(longRangeScout);
		const colonizer = new CommandedFleet(santaMaria);

		const ownedPlanet = new CommandedPlanet();
		ownedPlanet.playerNum = 2;

		const goodPlanet = new CommandedPlanet();
		goodPlanet.playerNum = None;
		goodPlanet.spec.terraformedHabitability = 100;

		// scouts can't colonize
		expect(scout.canColonize(goodPlanet)).toBe(false);

		// colonizers need colonists
		colonizer.cargo.colonists = 0;
		expect(colonizer.canColonize(goodPlanet)).toBe(false);

		// colonizer with colonists
		colonizer.cargo.colonists = 25;
		expect(colonizer.canColonize(goodPlanet)).toBe(true);

		// colonizer can't colonize owned planet
		colonizer.cargo.colonists = 25;
		expect(colonizer.canColonize(ownedPlanet)).toBe(false);

		// colonizer can't colonize bad planet
		const badPlanet = new CommandedPlanet();
		badPlanet.spec.terraformedHabitability = 0;
		colonizer.cargo.colonists = 25;
		expect(colonizer.canColonize(badPlanet)).toBe(false);
	});

	it('tests canRemoteMine', () => {
		const scout = new CommandedFleet(longRangeScout);
		const remoteMiner = new CommandedFleet(cottonPicker);

		const goodPlanet = new CommandedPlanet();
		goodPlanet.playerNum = None;

		const ownedPlanet = new CommandedPlanet();
		ownedPlanet.playerNum = 1;

		// create an AR player that can remote mine their own planets
		const arPlayer = new Player();
		arPlayer.race.spec = Object.assign(arPlayer.race.spec ?? {}, {
			canRemoteMineOwnPlanets: true
		} as RaceSpec);
		arPlayer.num = 1;

		// mine away!
		expect(remoteMiner.canRemoteMine(new Player(), goodPlanet)).toBe(true);
		// scouts can't remote mine
		expect(scout.canRemoteMine(new Player(), goodPlanet)).toBe(false);

		// can't remote mine an owned planet
		expect(remoteMiner.canRemoteMine(new Player(), ownedPlanet)).toBe(false);

		// can remote mine an owned planet if we are AR
		expect(remoteMiner.canRemoteMine(arPlayer, ownedPlanet)).toBe(true);

		// ar can't remote mine an owned planet if it's owned by someone else
		const ownedBySomeoneElsePlanet = new CommandedPlanet();
		ownedBySomeoneElsePlanet.playerNum = 2;
		expect(remoteMiner.canRemoteMine(arPlayer, ownedBySomeoneElsePlanet)).toBe(false);
	});

	it('tests canGate', () => {
		const scout = new CommandedFleet(longRangeScout);

		// make a new player that is friendly to player 2, not to player 3
		const player = new Player();
		player.num = 1;
		player.relations = [
			{
				relation: PlayerRelation.Friend
			},
			{
				relation: PlayerRelation.Friend
			},
			{
				relation: PlayerRelation.Enemy
			}
		];

		// make a couple gate planets
		const orbiting = new CommandedPlanet();
		orbiting.playerNum = 1;
		orbiting.spec.hasStargate = true;
		orbiting.spec.safeHullMass = 100;
		orbiting.spec.safeRange = 100;

		const target = new CommandedPlanet();
		target.playerNum = 1;
		target.spec.hasStargate = true;
		target.spec.safeHullMass = 100;
		target.spec.safeRange = 100;

		// can gate
		expect(scout.canGate(player, orbiting, target, 100, 100)).toBe(true);

		// can't gate, not orbiting
		expect(scout.canGate(player, undefined, target, 100, 100)).toBe(false);

		// can gate, have jump device
		scout.spec.canJump = true;
		expect(scout.canGate(player, undefined, target, 100, 100)).toBe(true);
		scout.spec.canJump = false;

		// can't gate, dest is unfriendly
		target.playerNum = 3;
		expect(scout.canGate(player, orbiting, target, 100, 100)).toBe(false);
		target.playerNum = 1;

		// can gate, dest is friendly
		target.playerNum = 2;
		expect(scout.canGate(player, orbiting, target, 100, 100)).toBe(true);
		target.playerNum = 1;

		// can't gate, too far
		expect(scout.canGate(player, orbiting, target, 200, 100)).toBe(false);

		// can't gate, too heavy
		expect(scout.canGate(player, orbiting, target, 100, 200)).toBe(false);

		// can't gate, have cargo
		scout.cargo.colonists = 10;
		expect(scout.canGate(player, orbiting, target, 100, 100)).toBe(false);
		scout.cargo.colonists = 0;

		// can gate, have cargo, can gate cargo
		const itPlayer = new Player();
		itPlayer.race.spec = Object.assign({}, { canGateCargo: true }) as RaceSpec;
		scout.cargo.colonists = 10;
		expect(scout.canGate(itPlayer, orbiting, target, 100, 100)).toBe(true);
		scout.cargo.colonists = 0;
	});

	it('returns minimal speeds for distances', () => {
		const fleet = new CommandedFleet(longRangeScout);
		// one year to go 49 ly
		expect(fleet.getMinimalWarp(49, 7, 1, 9)).toBe(7);

		// two years at warp 7, two at warp 6 or warp 5, pick warp 5
		expect(fleet.getMinimalWarp(50, 7, 1, 9)).toBe(5);

		// warp 6 takes 3 years to 72, so in 73 we can make it in 2 at warp 7
		expect(fleet.getMinimalWarp(73, 7, 1, 9)).toBe(7);

		// this is obvious
		expect(fleet.getMinimalWarp(36, 7, 1, 9)).toBe(6);

		// might as well go warp 5
		expect(fleet.getMinimalWarp(25, 7, 1, 9)).toBe(5);

		// getMinimalWarp is called by getMaxWarp as well to find the mininmum from some
		// max starting point. Make sure we don't exceed safe warp
		// if we have enough fuel to go warp 10, but our engine is only safe at warp 9, slow it down
		expect(fleet.getMinimalWarp(100, 10, 1, 9)).toBe(9);
	});

	it('returns fastest speed without running out of fuel', () => {
		const fleet = new CommandedFleet(longRangeScout);
		const designFinder = new TestDesignFinder();
		// one year to go 81 ly, plenty `of fuel
		fleet.fuel = 300;
		expect(fleet.getMaxWarp(designFinder, 0, 81, 1, 9)).toBe(9);

		// one year to go 64 ly, plenty of fuel
		fleet.fuel = 300;
		expect(fleet.getMaxWarp(designFinder, 0, 64, 1, 9)).toBe(8);

		// don't go faster than we need, go warp 5 for 25ly
		fleet.fuel = 300;
		expect(fleet.getMaxWarp(designFinder, 0, 25, 1, 9)).toBe(5);

		// make sure we don't run out of fuel over long distances
		// 300 ly at warp 9 would take 338mg of fuel, so go warp 8 (using 282mg fuel)
		fleet.fuel = 300;
		expect(fleet.getMaxWarp(designFinder, 0, 300, 1, 9)).toBe(8);
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
