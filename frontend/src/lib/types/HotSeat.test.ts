import { PlayerStatusSchema, type PlayerStatus } from '#lib/types/cs-proto.js';
import { create } from '@bufbuild/protobuf';
import { describe, expect, it } from 'vitest';
import { getFirstHotSeatPlayer, isHotSeat } from './HotSeat';

const player = (num: number, userId: number, submittedTurn = false, aiControlled = false) =>
	create(PlayerStatusSchema, { num, userId: BigInt(userId), submittedTurn, aiControlled });

describe('HotSeat', () => {
	const userId = BigInt(1);

	it('detects hot seat games', () => {
		expect(isHotSeat([player(1, 1), player(2, 0, false, true)], userId)).toBe(false);
		expect(isHotSeat([player(1, 1), player(2, 2)], userId)).toBe(false);
		expect(isHotSeat([player(1, 1), player(2, 1)], userId)).toBe(true);
	});

	it('finds the first player to play', () => {
		const players: PlayerStatus[] = [player(1, 1, true), player(2, 2), player(3, 1)];
		expect(getFirstHotSeatPlayer(players, userId)?.num).toBe(3);
		expect(getFirstHotSeatPlayer([player(1, 1, true), player(2, 1, true)], userId)?.num).toBe(1);
	});
});
