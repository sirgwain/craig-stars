import { describe, expect, it } from 'vitest';
import { Battle, type BattleRecord } from './Battle';

const record: BattleRecord = {
	num: 1,
	position: { x: 1, y: 2 },
	tokens: [
		{
			num: 1,
			playerNum: 1,
			token: { designNum: 1, quantity: 1 },
			startingPosition: { x: 1, y: 4 }
		},
		{
			num: 2,
			playerNum: 2,
			token: { designNum: 3, quantity: 1 },
			startingPosition: { x: 8, y: 5 }
		}
	],
	actionsPerRound: [
		[],
		[
			{ type: 3, tokenNum: 2, from: { x: 8, y: 5 }, to: { x: 9, y: 5 } },
			{ type: 3, tokenNum: 1, from: { x: 1, y: 4 }, to: { x: 0, y: 5 } },
			{ type: 3, tokenNum: 2, from: { x: 9, y: 5 }, to: { x: 9, y: 4 } },
			{ type: 3, tokenNum: 1, from: { x: 0, y: 5 }, to: { x: 0, y: 4 } }
		],
		[
			{ type: 3, tokenNum: 2, from: { x: 9, y: 4 }, to: { x: 9, y: 3 } },
			{ type: 3, tokenNum: 1, from: { x: 0, y: 4 }, to: { x: 0, y: 5 } },
			{ type: 3, tokenNum: 2, from: { x: 9, y: 3 }, to: { x: 9, y: 2 } },
			{ type: 3, tokenNum: 1, from: { x: 0, y: 5 }, to: { x: -1, y: 6 } },
			{ type: 3, tokenNum: 2, from: { x: 9, y: 2 }, to: { x: 8, y: 2 } },
			{ type: 3, tokenNum: 1, from: { x: -1, y: 6 }, to: { x: 0, y: 5 } }
		]
	]
};

describe('battle test', () => {
	it('build a battle with actions', () => {
		const battle = new Battle(record.num, record.position, record);

		expect(battle.totalPhases, '10 phases').toBe(10);

		expect(battle.getTokensAtLocation(0, 1, 4)?.length, '1 token at 1, 4').toBe(1);
		expect(battle.getTokensAtLocation(0, 1, 4), 'token 1 at 1, 4').toEqual([
			{
				num: 1,
				playerNum: 1,
				token: { designNum: 1, quantity: 1 },
				startingPosition: { x: 1, y: 4 },
				x: 1,
				y: 4
			}
		]);

		expect(battle.getTokensAtLocation(0, 8, 5)?.length, '1 token at 8, 5').toBe(1);
		expect(battle.getTokensAtLocation(0, 8, 5), 'token 2 at 8, 5').toEqual([
			{
				num: 2,
				playerNum: 2,
				token: { designNum: 3, quantity: 1 },
				startingPosition: { x: 8, y: 5 },
				x: 8,
				y: 5
			}
		]);

		// phase 1, token 2 moves
		expect(battle.getActionForPhase(1), 'token 2 moved to 9, 5').toEqual({
			type: 3,
			tokenNum: 2,
			from: { x: 8, y: 5 },
		to: { x: 9, y: 5 }
		});

		// phase 2, token 1 moves
		expect(battle.getActionForPhase(2), 'token 1 moved to 0, 5').toEqual({
			type: 3,
			tokenNum: 1,
			from: { x: 1, y: 4 },
			to: { x: 0, y: 5 }
		});
	});
});
