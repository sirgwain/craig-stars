import { describe, expect, it } from 'vitest';
import type { DesignFinder } from '#lib/services/Universe.js';
import { PlayerRelation, ShipClass, type Fleet, type Minefield, type ShipDesign } from './cs-proto';
import {
	filterFleet,
	filterMinefield,
	getMinefieldOwner,
	getShipClass,
	getVisibleShipCount,
	ShipClasses
} from './Filter';
import { PlayerSettings } from './PlayerSettings';

import { CommandedPlayer } from './Player';

const player = new CommandedPlayer();
player.num = 1;
player.relations = [
	{ relation: PlayerRelation.FRIEND },
	{ relation: PlayerRelation.FRIEND },
	{ relation: PlayerRelation.ENEMY },
	{ relation: PlayerRelation.NEUTRAL }
] as never;

// every other player's design 1 is a scout and design 2 is a battleship
const designs: DesignFinder = {
	getDesign: (playerNum, num) =>
		({ spec: { shipClass: num === 1 ? ShipClass.SCOUT : ShipClass.WARSHIP } }) as ShipDesign,
	getMyDesign: () => undefined
};

function fleet(playerNum: number, ...tokens: { designNum: number; quantity: number }[]): Fleet {
	return { mapObject: { playerNum }, tokens } as never;
}

describe('ship classes', () => {
	it('lists every class but unspecified', () => {
		expect(ShipClasses).toHaveLength(8);
		expect(ShipClasses).not.toContain(ShipClass.UNSPECIFIED);
	});

	it('reads the class from the design', () => {
		expect(getShipClass({ spec: { shipClass: ShipClass.MINER } } as ShipDesign)).toBe(
			ShipClass.MINER
		);
		// designs without a class are utility ships
		expect(getShipClass({ spec: {} } as ShipDesign)).toBe(ShipClass.UTILITY);
		expect(getShipClass(undefined)).toBe(ShipClass.UTILITY);
	});
});

describe('design filters', () => {
	const mine = fleet(1, { designNum: 1, quantity: 2 }, { designNum: 2, quantity: 3 });
	const ally = fleet(2, { designNum: 1, quantity: 2 }, { designNum: 2, quantity: 3 });
	const enemy = fleet(3, { designNum: 1, quantity: 2 }, { designNum: 2, quantity: 3 });
	const neutral = fleet(4, { designNum: 2, quantity: 4 });

	it('shows everything with no filters', () => {
		for (const f of [mine, ally, enemy]) {
			expect(filterFleet(player, f, {}, designs)).toBe(true);
			expect(getVisibleShipCount(player, f, {}, designs)).toBe(5);
		}
	});

	it('ignores the chosen designs while the filter is off', () => {
		const options = { filterDesigns: [1], filterEnemyShipClasses: [ShipClass.SCOUT] };
		expect(getVisibleShipCount(player, mine, options, designs)).toBe(5);
		expect(getVisibleShipCount(player, enemy, options, designs)).toBe(5);
	});

	it('shows only my chosen designs', () => {
		const options = { filterMyDesigns: true, filterDesigns: [2] };
		expect(getVisibleShipCount(player, mine, options, designs)).toBe(3);
		expect(filterFleet(player, mine, options, designs)).toBe(true);
		expect(filterFleet(player, fleet(1, { designNum: 1, quantity: 9 }), options, designs)).toBe(
			false
		);
		// other players are not affected
		expect(getVisibleShipCount(player, enemy, options, designs)).toBe(5);
	});

	it('hides all my fleets when no designs are chosen', () => {
		expect(filterFleet(player, mine, { filterMyDesigns: true }, designs)).toBe(false);
	});

	it('shows only the chosen classes of enemy and neutral ships', () => {
		const options = { filterEnemyDesigns: true, filterEnemyShipClasses: [ShipClass.SCOUT] };
		expect(getVisibleShipCount(player, enemy, options, designs)).toBe(2);
		expect(filterFleet(player, neutral, options, designs)).toBe(false);
		expect(getVisibleShipCount(player, ally, options, designs)).toBe(5);
		expect(getVisibleShipCount(player, mine, options, designs)).toBe(5);
	});

	it('shows only the chosen classes of ally ships', () => {
		const options = { filterAllyDesigns: true, filterAllyShipClasses: [ShipClass.WARSHIP] };
		expect(getVisibleShipCount(player, ally, options, designs)).toBe(3);
		expect(getVisibleShipCount(player, enemy, options, designs)).toBe(5);
	});
});

describe('minefield filters', () => {
	const minefield = (playerNum: number) => ({ mapObject: { playerNum } }) as Minefield;
	const shown = (options: Parameters<typeof filterMinefield>[2]) =>
		[1, 2, 3, 4].filter((playerNum) => filterMinefield(player, minefield(playerNum), options));

	it('knows whose minefield it is', () => {
		expect(
			[1, 2, 3, 4].map((playerNum) => getMinefieldOwner(player, minefield(playerNum)))
		).toEqual(['mine', 'friends', 'enemies', 'neutrals']);
	});

	it('shows the chosen minefields', () => {
		expect(shown({})).toEqual([1, 2, 3, 4]);
		expect(shown({ showMinefields: false })).toEqual([]);
		expect(shown({ showMinefieldsOf: ['mine', 'enemies'] })).toEqual([1, 3]);
		expect(shown({ showMinefields: false, showMinefieldsOf: ['mine'] })).toEqual([]);
	});

	it('turns minefields off when there are none to show, and back on', () => {
		const settings = new PlayerSettings();
		settings.setShownMinefieldOwners(['neutrals']);
		expect(shown(settings)).toEqual([4]);
		settings.toggleMinefields();
		expect(shown(settings)).toEqual([]);
		settings.toggleMinefields();
		expect(shown(settings)).toEqual([4]);

		settings.setShownMinefieldOwners([]);
		expect(settings.showMinefields).toBe(false);
		settings.toggleMinefields();
		expect(shown(settings)).toEqual([1, 2, 3, 4]);
	});
});
