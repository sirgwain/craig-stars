import { clone, create, fromJsonString, toJsonString } from '@bufbuild/protobuf';
import {
	BattleAttackWho,
	BattlePlanSchema,
	BattleTactic,
	BattleTarget,
	BattleSimulationSideSchema,
	BattleSimulationStackSchema,
	PlayerSchema,
	ShipDesignSchema,
	SimulateBattleRequestSchema,
	TechLevelSchema,
	type BattleSimulationSide,
	type BattleSimulationStack,
	type ShipDesign,
	type SimulateBattleResponse
} from '#lib/types/cs-proto.js';
import { CommandedPlayer } from '#lib/types/Player.js';
import { humanoid } from '#lib/types/Race.js';

export const techFields = [
	'energy',
	'weapons',
	'propulsion',
	'construction',
	'electronics',
	'biotechnology'
] as const;
export function newStack(design: ShipDesign): BattleSimulationStack {
	return create(BattleSimulationStackSchema, {
		design: clone(ShipDesignSchema, design),
		quantity: design.spec?.starbase ? 1 : 10,
		plan: {
			name: 'Default',
			tactic: BattleTactic.MAXIMIZE_DAMAGE,
			primaryTarget: BattleTarget.ANY,
			secondaryTarget: BattleTarget.ANY,
			attackWho: BattleAttackWho.EVERYONE
		}
	});
}
export function sampleDesign(): ShipDesign {
	return create(ShipDesignSchema, {
		name: 'Laser Destroyer',
		hull: 'Destroyer',
		slots: [
			{ hullSlotIndex: 1, hullComponent: 'Long Hump 6', quantity: 1 },
			{ hullSlotIndex: 2, hullComponent: 'Laser', quantity: 1 },
			{ hullSlotIndex: 3, hullComponent: 'Laser', quantity: 1 }
		]
	});
}
export function newSides(): BattleSimulationSide[] {
	return [0, 1].map(() =>
		create(BattleSimulationSideSchema, {
			race: humanoid(),
			techLevels: create(TechLevelSchema, Object.fromEntries(techFields.map((f) => [f, 10]))),
			stacks: [newStack(sampleDesign())]
		})
	);
}
export function owner(side: BattleSimulationSide): CommandedPlayer {
	return new CommandedPlayer(
		create(PlayerSchema, {
			num: 1,
			race: side.race,
			techLevels: side.techLevels,
			acquiredTechs: side.acquiredTechs
		})
	);
}
export function encodeScenario(sides: BattleSimulationSide[]): string {
	return toJsonString(SimulateBattleRequestSchema, create(SimulateBattleRequestSchema, { sides }));
}
export function decodeScenario(json: string): BattleSimulationSide[] {
	const sides = fromJsonString(SimulateBattleRequestSchema, json).sides;
	if (
		sides.length !== 2 ||
		sides.some(
			(s) =>
				!s.race ||
				!s.techLevels ||
				s.stacks.length > 64 ||
				s.stacks.some((t) => !t.design || !t.plan)
		)
	)
		throw new Error('Invalid saved scenario');
	return sides;
}
export function validateSides(sides: BattleSimulationSide[]): void {
	for (const [i, side] of sides.entries()) {
		if (!side.stacks.length || side.stacks.length > 64)
			throw new Error(`Side ${i + 1} needs 1 to 64 stacks.`);
		for (const field of techFields)
			if (
				!Number.isInteger(side.techLevels?.[field]) ||
				side.techLevels![field] < 0 ||
				side.techLevels![field] > 26
			)
				throw new Error('Tech levels must be whole numbers from 0 to 26.');
		for (const stack of side.stacks) {
			if (!Number.isInteger(stack.quantity) || stack.quantity < 1 || stack.quantity > 10000)
				throw new Error('Stack quantities must be whole numbers from 1 to 10,000.');
			if (
				!Number.isInteger(stack.quantityDamaged) ||
				stack.quantityDamaged < 0 ||
				stack.quantityDamaged > stack.quantity ||
				!Number.isFinite(stack.damage) ||
				stack.damage < 0
			)
				throw new Error('Check the number of damaged ships and damage per ship.');
		}
	}
}
export type SimulationRun = {
	response: SimulateBattleResponse;
	sides: BattleSimulationSide[];
	seed: bigint;
};
export function outcome(response: SimulateBattleResponse): string {
	const active = [1, 2].map((p) =>
		response.stacks
			.filter((s) => s.playerNum === p)
			.reduce((n, s) => n + s.survivors - s.retreated, 0)
	);
	if (active[0] && !active[1]) return 'A holds the field';
	if (active[1] && !active[0]) return 'B holds the field';
	return 'No decisive outcome';
}
export function copyStack(stack: BattleSimulationStack): BattleSimulationStack {
	return clone(BattleSimulationStackSchema, stack);
}
export const defaultPlan = () => create(BattlePlanSchema, newStack(sampleDesign()).plan);
