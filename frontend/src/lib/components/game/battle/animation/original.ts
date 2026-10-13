// The "Original" battle style: the look of the original Stars! battle replay (VCR), following
// how the original draws an attack.
//
// A beam is two thin lines from two points on the shooter's square to the center of the target,
// red for weak beams and blue for strong ones, with a small bang on the target. A torpedo is a
// single small icon that flies from the shooter to the target, a few pixels per frame, followed
// by a bigger bang if it hit. When the attack is over, every target gets a bang: a big one if a
// ship was destroyed. Ships jump between squares. There is no sound, shake or particles.
//
// Everything here is pure (no canvas) so it can be tested; OriginalRenderer draws the result.

import { BattleRecordTokenActionType, type BattleRecordTokenAction } from '#lib/types/cs-proto.js';

export type Point = { x: number; y: number };

export type VcrIcon = 'bang1' | 'bang2' | 'bang3' | 'torp1' | 'torp2' | 'torp3' | 'torp4';
export const VCR_ICONS: VcrIcon[] = ['bang1', 'bang2', 'bang3', 'torp1', 'torp2', 'torp3', 'torp4'];
// the icons are 32x32 and drawn centered on a point
export const VCR_ICON_SIZE = 32;
const TORPEDO_ICONS: VcrIcon[] = ['torp1', 'torp2', 'torp3', 'torp4'];

// beams of this much damage per part or more are drawn blue, weaker ones red
export const STRONG_BEAM_POWER = 200;
export const WEAK_BEAM_COLOR = '#ff0000';
export const STRONG_BEAM_COLOR = '#0000ff';

export function beamColor(power: number | undefined): string {
	return (power ?? 0) >= STRONG_BEAM_POWER ? STRONG_BEAM_COLOR : WEAK_BEAM_COLOR;
}

// the four points a third of the square out from a square's center
export type SquarePoints = { center: Point; top: Point; bottom: Point; left: Point; right: Point };

export function squarePoints(square: { x?: number; y?: number } | undefined, cell: number) {
	const center = {
		x: ((square?.x ?? 0) + 0.5) * cell,
		y: ((square?.y ?? 0) + 0.5) * cell
	};
	const r = cell / 3;
	return {
		center,
		top: { x: center.x, y: center.y - r },
		bottom: { x: center.x, y: center.y + r },
		left: { x: center.x - r, y: center.y },
		right: { x: center.x + r, y: center.y }
	} satisfies SquarePoints;
}

export type ShotPoints = {
	// the two beam lines start here and meet at the target's center
	beams: [Point, Point];
	// a torpedo starts here, on the side facing the target
	torpedo: Point;
	target: Point;
};

// Where a shot starts on the shooter's square. dx and dy are shooter minus target, in squares.
// A target straight up or down (or nearly, a long way off) is shot from the left and right
// points, one straight left or right from the top and bottom points, and a diagonal one from
// the top or bottom point and the left or right point on its side. A torpedo starts at the
// point (or corner) facing the target. undefined when they share a square: nothing is drawn.
export function shotPoints(
	from: { x?: number; y?: number } | undefined,
	to: { x?: number; y?: number } | undefined,
	cell: number
): ShotPoints | undefined {
	const dx = (from?.x ?? 0) - (to?.x ?? 0);
	const dy = (from?.y ?? 0) - (to?.y ?? 0);
	if (dx === 0 && dy === 0) {
		return undefined;
	}
	const p = squarePoints(from, cell);
	const target = squarePoints(to, cell).center;
	if (dx === 0 || (Math.abs(dx) === 1 && Math.abs(dy) > 2)) {
		// up or down
		return { beams: [p.right, p.left], torpedo: dy > 0 ? p.top : p.bottom, target };
	}
	if (dy === 0 || (Math.abs(dy) === 1 && Math.abs(dx) > 2)) {
		// left or right
		return { beams: [p.top, p.bottom], torpedo: dx > 0 ? p.left : p.right, target };
	}
	// diagonal: the top or bottom point toward the target, and the side point toward it
	const vertical = dy > 0 ? p.top : p.bottom;
	const side = dx > 0 ? p.left : p.right;
	return { beams: [vertical, side], torpedo: { x: side.x, y: vertical.y }, target };
}

// a torpedo flight takes 8 frames per square of distance (the longer of across and down)
export function torpedoFrames(
	from: { x?: number; y?: number } | undefined,
	to: { x?: number; y?: number } | undefined
): number {
	const dx = (from?.x ?? 0) - (to?.x ?? 0);
	const dy = (from?.y ?? 0) - (to?.y ?? 0);
	return Math.max(Math.abs(dx), Math.abs(dy)) * 8;
}

// Where the torpedo is drawn in a frame of its flight. Like the original, the first two frames
// are drawn at the start and the last one a frame short of the target, where the hit bang goes.
export function torpedoPosition(frame: number, frames: number, from: Point, to: Point): Point {
	const k = frames > 0 ? Math.max(0, frame - 1) / frames : 1;
	return {
		x: Math.round(from.x + (to.x - from.x) * k),
		y: Math.round(from.y + (to.y - from.y) * k)
	};
}

// The original had five playback speeds (0-4); ours are 0.5x, 1x, 2x and 4x, mapped to its
// first four (its fastest skipped the torpedo flights).
export function vcrSpeed(speed: number): number {
	if (speed <= 0.5) return 0;
	if (speed <= 1) return 1;
	if (speed <= 2) return 2;
	return 3;
}

// milliseconds between steps of the playback at a speed (0-4)
export function vcrStepMs(vcrSpeed: number): number {
	return 570 - 120 * vcrSpeed;
}

// milliseconds a torpedo frame is shown at a speed (0-4)
export function vcrFrameMs(vcrSpeed: number): number {
	return Math.max(1, 35 - 10 * vcrSpeed);
}

const isFire = (a: BattleRecordTokenAction | undefined) =>
	a?.type === BattleRecordTokenActionType.BEAM_FIRE ||
	a?.type === BattleRecordTokenActionType.TORPEDO_FIRE;

export type StepKind = 'move' | 'beam' | 'torpedo' | 'flee' | 'other';

export function actionKind(action: BattleRecordTokenAction | undefined): StepKind {
	switch (action?.type) {
		case BattleRecordTokenActionType.MOVE:
			return 'move';
		case BattleRecordTokenActionType.BEAM_FIRE:
			return 'beam';
		case BattleRecordTokenActionType.TORPEDO_FIRE:
			return 'torpedo';
		case BattleRecordTokenActionType.RAN_AWAY:
			return 'flee';
		default:
			return 'other';
	}
}

// A step of the playback covers battle.actions[fromPhase .. toPhase - 1]; phase n is the board
// after battle.actions[n - 1].
export type VcrStep = {
	kind: StepKind;
	fromPhase: number;
	toPhase: number;
	actions: BattleRecordTokenAction[];
};

// A step of the original playback: one phase, except that the shots one weapon slot fires in a
// round (a gatling gun hitting everything in range, or damage carried on to the next target)
// are one attack and drawn together.
export function originalStep(
	actions: BattleRecordTokenAction[],
	phase: number
): VcrStep | undefined {
	if (phase < 0 || phase >= actions.length) {
		return undefined;
	}
	const first = actions[phase];
	const stepActions = [first];
	if (isFire(first)) {
		for (let i = phase + 1; i < actions.length; i++) {
			const a = actions[i];
			if (
				a.type !== first.type ||
				a.tokenNum !== first.tokenNum ||
				a.slot !== first.slot ||
				a.round !== first.round
			) {
				break;
			}
			stepActions.push(a);
		}
	}
	return {
		kind: actionKind(first),
		fromPhase: phase,
		toPhase: phase + stepActions.length,
		actions: stepActions
	};
}

// the shots of the attack shown on the board at a phase (the board shows the attack that led to
// it), up to that phase
export function attackAtPhase(
	actions: BattleRecordTokenAction[],
	phase: number
): BattleRecordTokenAction[] {
	const last = actions[phase - 1];
	if (!isFire(last)) {
		return [];
	}
	let start = phase - 1;
	while (start > 0) {
		const a = actions[start - 1];
		if (
			a.type !== last.type ||
			a.tokenNum !== last.tokenNum ||
			a.slot !== last.slot ||
			a.round !== last.round
		) {
			break;
		}
		start--;
	}
	return actions.slice(start, phase);
}

export type VcrLine = { from: Point; to: Point; color: string };
export type VcrStamp = { icon: VcrIcon; at: Point };

// what the board shows of an attack at a moment
export type VcrFrame = {
	lines: VcrLine[];
	// icons, drawn in order (later ones over earlier ones)
	stamps: VcrStamp[];
	// the torpedo in flight
	torpedo?: VcrStamp;
	// every torpedo has arrived
	done: boolean;
};

export type VcrAttackOptions = {
	cell: number;
	// ms per torpedo frame
	frameMs: number;
	// damage per part of the firing beam weapon, for its color
	beamPower?: number;
	// no torpedo flights (reduced motion)
	still?: boolean;
};

// how long the torpedo flights of an attack take, in ms
export function attackFlightMs(shots: BattleRecordTokenAction[], opts: VcrAttackOptions): number {
	if (opts.still) return 0;
	return shots
		.filter((s) => s.type === BattleRecordTokenActionType.TORPEDO_FIRE)
		.reduce((ms, s) => ms + torpedoFrames(s.from, s.to) * opts.frameMs, 0);
}

// The attack drawn at t ms after it started (Infinity: the finished attack, as the paused board
// shows it). Each target in turn: beam lines and a bang, or a torpedo flight and a hit bang
// unless every torpedo missed. Then a bang on every target, a big one where a ship died.
export function attackFrame(
	shots: BattleRecordTokenAction[],
	t: number,
	opts: VcrAttackOptions
): VcrFrame {
	const frame: VcrFrame = { lines: [], stamps: [], done: false };
	let cursor = 0;
	for (const shot of shots) {
		const points = shotPoints(shot.from, shot.to, opts.cell);
		if (!points) {
			continue;
		}
		if (shot.type === BattleRecordTokenActionType.BEAM_FIRE) {
			const color = beamColor(opts.beamPower);
			frame.lines.push(
				{ from: points.beams[0], to: points.target, color },
				{ from: points.beams[1], to: points.target, color }
			);
			frame.stamps.push({ icon: 'bang1', at: points.target });
		} else if (shot.type === BattleRecordTokenActionType.TORPEDO_FIRE) {
			const frames = opts.still ? 0 : torpedoFrames(shot.from, shot.to);
			const flight = frames * opts.frameMs;
			if (t < cursor + flight) {
				const i = Math.floor((t - cursor) / opts.frameMs);
				frame.torpedo = {
					icon: TORPEDO_ICONS[i & 3],
					at: torpedoPosition(i, frames, points.torpedo, points.target)
				};
				return frame;
			}
			cursor += flight;
			if ((shot.torpedoHits ?? 0) > 0) {
				frame.stamps.push({ icon: 'bang2', at: points.target });
			}
		}
	}
	for (const shot of shots) {
		frame.stamps.push({
			icon: (shot.tokensDestroyed ?? 0) > 0 ? 'bang3' : 'bang1',
			at: squarePoints(shot.to, opts.cell).center
		});
	}
	frame.done = true;
	return frame;
}
