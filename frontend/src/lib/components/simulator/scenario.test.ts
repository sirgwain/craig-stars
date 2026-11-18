import { describe, it, expect } from 'vitest';
import { clone, create } from '@bufbuild/protobuf';
import {
	BattleSimulationStackSchema,
	SimulateBattleResponseSchema,
	ShipDesignSchema
} from '#lib/types/cs-proto.js';
import {
	newSides,
	newStack,
	encodeScenario,
	decodeScenario,
	validateSides,
	outcome
} from './scenario';

describe('simulation scenarios', () => {
	it('round trips race, technology, damage and independent stack plans', () => {
		const sides = newSides();
		const split = clone(BattleSimulationStackSchema, sides[0].stacks[0]);
		split.quantityDamaged = 3;
		split.damage = 12;
		split.plan!.tactic = 1;
		sides[0].stacks.push(split);
		const restored = decodeScenario(encodeScenario(sides));
		expect(restored).toEqual(sides);
		restored[0].stacks[1].design!.name = 'Changed';
		expect(sides[0].stacks[0].design!.name).toBe('Laser Destroyer');
		expect(restored[0].stacks[0].design!.name).toBe('Laser Destroyer');
	});
	it('rejects incomplete or out of range values before calling WASM', () => {
		const sides = newSides();
		sides[0].stacks[0].quantity = 1.5;
		expect(() => validateSides(sides)).toThrow(/whole numbers/);
		sides[0].stacks[0].quantity = 10;
		sides[1].techLevels!.weapons = 27;
		expect(() => validateSides(sides)).toThrow(/Tech levels/);
		expect(() => decodeScenario('{}')).toThrow(/Invalid saved/);
	});
	it('distinguishes retreating survivors from ships holding the field', () => {
		const response = create(SimulateBattleResponseSchema, {
			stacks: [
				{ playerNum: 1, survivors: 5 },
				{ playerNum: 2, survivors: 3, retreated: 3 }
			]
		});
		expect(outcome(response)).toBe('A holds the field');
		response.stacks[0].retreated = 5;
		expect(outcome(response)).toBe('No decisive outcome');
		response.stacks[1].retreated = 0;
		expect(outcome(response)).toBe('B holds the field');
	});
	it('adds starbases as a single unit', () => {
		expect(
			newStack(create(ShipDesignSchema, { hull: 'Space Station', spec: { starbase: true } }))
				.quantity
		).toBe(1);
	});
});
