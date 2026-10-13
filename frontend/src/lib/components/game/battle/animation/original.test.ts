import { create } from '@bufbuild/protobuf';
import { describe, expect, it } from 'vitest';
import {
	BattleRecordTokenActionSchema,
	BattleRecordTokenActionType,
	type BattleRecordTokenAction
} from '#lib/types/cs-proto.js';
import {
	STRONG_BEAM_COLOR,
	WEAK_BEAM_COLOR,
	attackAtPhase,
	attackFlightMs,
	attackFrame,
	beamColor,
	originalStep,
	shotPoints,
	squarePoints,
	torpedoFrames,
	torpedoPosition,
	vcrFrameMs,
	vcrSpeed,
	vcrStepMs
} from './original';

// a 66px square: its center is 33px in, the beam points 22px from the center
const cell = 66;
const at = (x: number, y: number) => ({ x, y });

type Extra = {
	slot?: number;
	tokensDestroyed?: number;
	torpedoHits?: number;
	torpedoMisses?: number;
};

const shot = (
	type: BattleRecordTokenActionType,
	from: [number, number],
	to: [number, number],
	extra: Extra = {}
): BattleRecordTokenAction =>
	create(BattleRecordTokenActionSchema, {
		type,
		tokenNum: 1,
		round: 1,
		slot: 2,
		from: at(...from),
		to: at(...to),
		targetNum: 2,
		...extra
	});
const beam = (from: [number, number], to: [number, number], extra: Extra = {}) =>
	shot(BattleRecordTokenActionType.BEAM_FIRE, from, to, extra);
const torpedo = (from: [number, number], to: [number, number], extra: Extra = {}) =>
	shot(BattleRecordTokenActionType.TORPEDO_FIRE, from, to, extra);
const move = (tokenNum: number) =>
	create(BattleRecordTokenActionSchema, {
		type: BattleRecordTokenActionType.MOVE,
		tokenNum,
		round: 1,
		from: at(0, 0),
		to: at(1, 0)
	});

describe('shotPoints', () => {
	const p = squarePoints(at(4, 4), cell);
	const center = (x: number, y: number) => squarePoints(at(x, y), cell).center;

	it('puts the points a third of the square out from the center', () => {
		expect(p.center).toEqual({ x: 297, y: 297 });
		expect(p.top).toEqual({ x: 297, y: 275 });
		expect(p.bottom).toEqual({ x: 297, y: 319 });
		expect(p.left).toEqual({ x: 275, y: 297 });
		expect(p.right).toEqual({ x: 319, y: 297 });
	});

	it('shoots up and down from the left and right points', () => {
		// straight up
		expect(shotPoints(at(4, 4), at(4, 1), cell)).toEqual({
			beams: [p.right, p.left],
			torpedo: p.top,
			target: center(4, 1)
		});
		// straight down
		expect(shotPoints(at(4, 4), at(4, 6), cell)?.torpedo).toEqual(p.bottom);
		// nearly up, a long way off
		expect(shotPoints(at(4, 4), at(5, 0), cell)).toEqual({
			beams: [p.right, p.left],
			torpedo: p.top,
			target: center(5, 0)
		});
	});

	it('shoots left and right from the top and bottom points', () => {
		// straight left
		expect(shotPoints(at(4, 4), at(1, 4), cell)).toEqual({
			beams: [p.top, p.bottom],
			torpedo: p.left,
			target: center(1, 4)
		});
		// straight right
		expect(shotPoints(at(4, 4), at(8, 4), cell)?.torpedo).toEqual(p.right);
		// nearly right, a long way off
		expect(shotPoints(at(4, 4), at(8, 5), cell)).toEqual({
			beams: [p.top, p.bottom],
			torpedo: p.right,
			target: center(8, 5)
		});
	});

	it('shoots diagonally from the points facing the target, torpedoes from that corner', () => {
		// up and left
		expect(shotPoints(at(4, 4), at(2, 2), cell)).toEqual({
			beams: [p.top, p.left],
			torpedo: { x: p.left.x, y: p.top.y },
			target: center(2, 2)
		});
		// down and left
		expect(shotPoints(at(4, 4), at(2, 6), cell)).toEqual({
			beams: [p.bottom, p.left],
			torpedo: { x: p.left.x, y: p.bottom.y },
			target: center(2, 6)
		});
		// up and right
		expect(shotPoints(at(4, 4), at(6, 3), cell)).toEqual({
			beams: [p.top, p.right],
			torpedo: { x: p.right.x, y: p.top.y },
			target: center(6, 3)
		});
		// down and right: one across and two down is still a diagonal
		expect(shotPoints(at(4, 4), at(5, 6), cell)).toEqual({
			beams: [p.bottom, p.right],
			torpedo: { x: p.right.x, y: p.bottom.y },
			target: center(5, 6)
		});
	});

	it('draws nothing for a target in the same square', () => {
		expect(shotPoints(at(4, 4), at(4, 4), cell)).toBeUndefined();
	});
});

describe('beamColor', () => {
	it('is red below 200 damage per part and blue from 200', () => {
		expect(beamColor(10)).toBe(WEAK_BEAM_COLOR);
		expect(beamColor(199)).toBe(WEAK_BEAM_COLOR);
		expect(beamColor(200)).toBe(STRONG_BEAM_COLOR);
		expect(beamColor(433)).toBe(STRONG_BEAM_COLOR);
		expect(WEAK_BEAM_COLOR).toBe('#ff0000');
		expect(STRONG_BEAM_COLOR).toBe('#0000ff');
	});

	it('is red for an unknown weapon', () => {
		expect(beamColor(undefined)).toBe(WEAK_BEAM_COLOR);
	});
});

describe('torpedo flight', () => {
	it('takes 8 frames per square of the longer distance', () => {
		expect(torpedoFrames(at(4, 4), at(5, 4))).toBe(8);
		expect(torpedoFrames(at(4, 4), at(1, 4))).toBe(24);
		expect(torpedoFrames(at(4, 4), at(6, 9))).toBe(40);
		expect(torpedoFrames(at(4, 4), at(4, 4))).toBe(0);
	});

	it('shows each frame 35ms less 10ms per speed step', () => {
		expect([0.5, 1, 2, 4].map((s) => vcrFrameMs(vcrSpeed(s)))).toEqual([35, 25, 15, 5]);
		expect([0.5, 1, 2, 4].map((s) => vcrStepMs(vcrSpeed(s)))).toEqual([570, 450, 330, 210]);
	});

	it('starts with two frames at the start and stops a frame short of the target', () => {
		const from = { x: 0, y: 0 };
		const to = { x: 80, y: 0 };
		expect(torpedoPosition(0, 8, from, to)).toEqual(from);
		expect(torpedoPosition(1, 8, from, to)).toEqual(from);
		expect(torpedoPosition(2, 8, from, to)).toEqual({ x: 10, y: 0 });
		expect(torpedoPosition(7, 8, from, to)).toEqual({ x: 60, y: 0 });
	});

	it('flies one icon, cycling its four frames, then bangs on a hit', () => {
		const t = torpedo([1, 5], [4, 5], { torpedoHits: 2, torpedoMisses: 1 });
		const opts = { cell, frameMs: 25 };
		expect(attackFlightMs([t], opts)).toBe(24 * 25);
		const icons = [0, 25, 50, 75, 100].map((ms) => attackFrame([t], ms, opts).torpedo?.icon);
		expect(icons).toEqual(['torp1', 'torp2', 'torp3', 'torp4', 'torp1']);
		const flying = attackFrame([t], 300, opts);
		expect(flying.done).toBe(false);
		expect(flying.stamps).toEqual([]);
		expect(flying.torpedo?.at.y).toBe(shotPoints(t.from, t.to, cell)?.torpedo.y);

		const landed = attackFrame([t], 600, opts);
		expect(landed.done).toBe(true);
		expect(landed.torpedo).toBeUndefined();
		expect(landed.stamps.map((s) => s.icon)).toEqual(['bang2', 'bang1']);
	});

	it('has no hit bang when every torpedo missed', () => {
		const t = torpedo([1, 5], [4, 5], { torpedoHits: 0, torpedoMisses: 3 });
		const frame = attackFrame([t], Infinity, { cell, frameMs: 25 });
		expect(frame.stamps.map((s) => s.icon)).toEqual(['bang1']);
	});

	it('does not fly with reduced motion', () => {
		const t = torpedo([1, 5], [4, 5], { torpedoHits: 1 });
		const opts = { cell, frameMs: 25, still: true };
		expect(attackFlightMs([t], opts)).toBe(0);
		expect(attackFrame([t], 0, opts).done).toBe(true);
	});
});

describe('beam attacks', () => {
	it('draws two lines to the target and a bang, a big one where a ship died', () => {
		const frame = attackFrame([beam([1, 5], [6, 5], { tokensDestroyed: 1 })], 0, {
			cell,
			frameMs: 25,
			beamPower: 433
		});
		expect(frame.lines).toHaveLength(2);
		expect(frame.lines.every((l) => l.color === STRONG_BEAM_COLOR)).toBe(true);
		expect(frame.stamps.map((s) => s.icon)).toEqual(['bang1', 'bang3']);
	});

	it('draws a gatling gun to every target at once', () => {
		const shots = [beam([3, 6], [5, 5]), beam([3, 6], [5, 6]), beam([3, 6], [5, 7])];
		const frame = attackFrame(shots, 0, { cell, frameMs: 25, beamPower: 31 });
		expect(frame.lines).toHaveLength(6);
		expect(frame.lines.every((l) => l.color === WEAK_BEAM_COLOR)).toBe(true);
		expect(frame.done).toBe(true);
	});
});

describe('originalStep', () => {
	it('groups the shots of one weapon slot into one attack', () => {
		const actions = [
			move(1),
			beam([3, 6], [5, 5]),
			beam([3, 6], [5, 6]),
			beam([3, 6], [5, 7], { slot: 3 }),
			move(2),
			move(3)
		];
		expect(originalStep(actions, 0)?.toPhase).toBe(1);
		const gatling = originalStep(actions, 1);
		expect(gatling?.kind).toBe('beam');
		expect(gatling?.actions).toHaveLength(2);
		expect(gatling?.toPhase).toBe(3);
		expect(originalStep(actions, 3)?.toPhase).toBe(4);
		// moves are not batched: ships jump one at a time
		expect(originalStep(actions, 4)?.toPhase).toBe(5);
		expect(originalStep(actions, 6)).toBeUndefined();
	});

	it('shows the attack that led to a phase', () => {
		const actions = [move(1), beam([3, 6], [5, 5]), beam([3, 6], [5, 6]), move(2)];
		expect(attackAtPhase(actions, 0)).toEqual([]);
		expect(attackAtPhase(actions, 1)).toEqual([]);
		expect(attackAtPhase(actions, 2)).toEqual([actions[1]]);
		expect(attackAtPhase(actions, 3)).toEqual([actions[1], actions[2]]);
		expect(attackAtPhase(actions, 4)).toEqual([]);
	});
});
