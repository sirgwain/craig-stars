import type { Battle } from '#lib/types/Battle.js';
import type { BattleRecordTokenAction } from '#lib/types/cs-proto.js';
import {
	actionKind,
	attackAtPhase,
	attackFlightMs,
	attackFrame,
	originalStep,
	vcrFrameMs,
	vcrSpeed,
	vcrStepMs,
	type VcrAttackOptions,
	type VcrFrame,
	type VcrStep
} from './original';
import { OriginalRenderer } from './OriginalRenderer';

export type PlaybackHost = {
	getPhase(): number;
	// commit a new phase to the board
	setPhase(phase: number): void;
	// damage per part of the weapon in a token's slot, when its design is known
	weaponPower?(tokenNum: number, slot: number): number | undefined;
};

type ActiveStep = {
	step: VcrStep;
	// seconds on the playback clock
	start: number;
	// seconds from the start until the step's phase is committed (after a torpedo flight)
	impactAt: number;
	// seconds from the start until the next step
	advanceAt: number;
	impactDone: boolean;
};

export const SPEEDS = [0.5, 1, 2, 4];
const playbackKey = 'cs.battle.playback';

// remembered per browser: whether attacks are animated, whether the damage they do is marked on
// their targets, and how fast the battle plays
export type PlaybackPrefs = { animations: boolean; damage: boolean; speed: number };
export const defaultPlaybackPrefs: PlaybackPrefs = { animations: true, damage: true, speed: 1 };

export function loadPlaybackPrefs(): PlaybackPrefs {
	try {
		const json = localStorage.getItem(playbackKey);
		if (json) {
			const prefs = JSON.parse(json) as Partial<PlaybackPrefs>;
			return {
				animations:
					typeof prefs.animations === 'boolean'
						? prefs.animations
						: defaultPlaybackPrefs.animations,
				damage: typeof prefs.damage === 'boolean' ? prefs.damage : defaultPlaybackPrefs.damage,
				speed:
					prefs.speed && SPEEDS.includes(prefs.speed) ? prefs.speed : defaultPlaybackPrefs.speed
			};
		}
	} catch {
		// storage unavailable, use defaults
	}
	return { ...defaultPlaybackPrefs };
}

export function savePlaybackPrefs(prefs: PlaybackPrefs) {
	try {
		localStorage.setItem(playbackKey, JSON.stringify(prefs));
	} catch {
		// storage unavailable, ignore
	}
}

// Plays a battle back in the style of the original game's battle replay: one step per phase (the
// shots of one weapon slot together) at the original's beat for the speed, drawing each attack on
// a canvas over the board. It waits at the start until it is played.
export class BattlePlayback {
	// torpedoes fly to their targets; when off, an attack is drawn right away
	animations = $state(true);
	// the damage of an attack is marked on its targets
	damage = $state(true);
	playing = $state(false);
	speed = $state(1);
	// a step is being played (playing, or a single step)
	animating = $state(false);
	// a torpedo is on its way: the board still shows the phase before the attack
	inFlight = $state(false);
	reducedMotion = $state(false);

	private renderer = new OriginalRenderer(() => this.draw());
	// seconds
	private clock = 0;
	private frozen = false;
	private current: ActiveStep | undefined;
	private raf = 0;
	private lastTs = 0;
	private ctx: CanvasRenderingContext2D | undefined;
	private width = 0;
	private height = 0;
	private dpr = 1;
	private cell = 68.6;
	private motionQuery: MediaQueryList | undefined;

	constructor(
		private battle: Battle,
		private host: PlaybackHost
	) {
		const prefs = loadPlaybackPrefs();
		this.animations = prefs.animations;
		this.damage = prefs.damage;
		this.speed = prefs.speed;
		this.renderer.load();
		if (typeof window !== 'undefined' && 'matchMedia' in window) {
			this.motionQuery = window.matchMedia('(prefers-reduced-motion: reduce)');
			this.reducedMotion = this.motionQuery.matches;
			this.motionQuery.addEventListener('change', this.onMotionChange);
		}
	}

	private onMotionChange = (e: MediaQueryListEvent) => {
		this.reducedMotion = e.matches;
	};

	destroy() {
		cancelAnimationFrame(this.raf);
		this.raf = 0;
		this.motionQuery?.removeEventListener('change', this.onMotionChange);
	}

	setBattle(battle: Battle) {
		if (battle === this.battle) {
			return;
		}
		// a universe refresh re-creates the same battle; keep playing in that case
		const same =
			battle.num === this.battle.num &&
			battle.totalPhases === this.battle.totalPhases &&
			battle.tokens.length === this.battle.tokens.length;
		this.battle = battle;
		if (!same) {
			this.stopAll();
		}
	}

	// attach the canvas the attacks are drawn on, and its size in css pixels
	setCanvas(canvas: HTMLCanvasElement | undefined, width: number, height: number) {
		this.ctx = canvas?.getContext('2d') ?? undefined;
		this.dpr = typeof window !== 'undefined' ? window.devicePixelRatio || 1 : 1;
		this.width = width;
		this.height = height;
		this.cell = width / 10;
		if (canvas) {
			canvas.width = Math.round(width * this.dpr);
			canvas.height = Math.round(height * this.dpr);
		}
		this.draw();
	}

	// ---- settings ----

	setSpeed(speed: number) {
		this.speed = speed;
		this.savePrefs();
	}

	setAnimations(animations: boolean) {
		if (animations === this.animations) {
			return;
		}
		// a step in progress was timed for the old setting
		this.finishCurrent();
		this.animations = animations;
		this.savePrefs();
		this.draw();
	}

	setDamage(damage: boolean) {
		this.damage = damage;
		this.savePrefs();
	}

	private savePrefs() {
		savePlaybackPrefs({ animations: this.animations, damage: this.damage, speed: this.speed });
	}

	// ---- transport ----

	togglePlay() {
		if (this.playing) {
			this.pause();
		} else {
			this.play();
		}
	}

	play() {
		if (this.host.getPhase() >= this.battle.totalPhases && !this.current) {
			// replay from the start
			this.host.setPhase(0);
		}
		this.playing = true;
		this.frozen = false;
		this.ensureLoop();
	}

	pause() {
		this.playing = false;
		// freeze mid-flight; resuming continues where it stopped
		this.frozen = !!this.current;
	}

	// play exactly one step, then stop
	stepForward() {
		this.playing = false;
		this.frozen = false;
		this.finishCurrent();
		this.startStep(this.clock);
		this.ensureLoop();
	}

	// go back exactly one phase from the committed phase
	stepBack() {
		this.playing = false;
		this.cancelCurrent();
		this.host.setPhase(Math.max(0, this.host.getPhase() - 1));
		this.draw();
	}

	// jump to a phase; keeps playing if we were playing
	seek(phase: number) {
		this.cancelCurrent();
		this.host.setPhase(Math.max(0, Math.min(this.battle.totalPhases, phase)));
		this.draw();
		if (this.playing) {
			this.ensureLoop();
		}
	}

	begin() {
		this.playing = false;
		this.seek(0);
	}

	end() {
		this.playing = false;
		this.seek(this.battle.totalPhases);
	}

	// go to the next attack and play it
	nextAttack() {
		const phase = this.current ? this.current.step.toPhase : this.host.getPhase();
		const index = this.battle.actions.findIndex((a, i) => {
			const kind = actionKind(a);
			return i >= phase && (kind === 'beam' || kind === 'torpedo');
		});
		if (index === -1) {
			this.end();
			return;
		}
		const wasPlaying = this.playing;
		this.seek(index);
		if (!wasPlaying) {
			this.stepForward();
		}
	}

	private stopAll() {
		this.playing = false;
		this.cancelCurrent();
		this.draw();
	}

	// ---- steps ----

	private cancelCurrent() {
		this.current = undefined;
		this.frozen = false;
		this.animating = false;
		this.inFlight = false;
	}

	// jump to the end of the current step
	private finishCurrent() {
		const current = this.current;
		if (current) {
			if (!current.impactDone) {
				this.host.setPhase(current.step.toPhase);
			}
			this.cancelCurrent();
		}
	}

	// The original moves on at a fixed beat per step. A torpedo flight takes as long as it takes
	// and its bang stays up for a moment after it lands.
	private startStep(at: number) {
		const step = originalStep(this.battle.actions, this.host.getPhase());
		if (!step) {
			this.playing = false;
			this.animating = false;
			return;
		}
		const interval = vcrStepMs(vcrSpeed(this.speed)) / 1000;
		let impactAt = 0;
		let advanceAt = interval;
		if (step.kind === 'torpedo') {
			impactAt = attackFlightMs(step.actions, this.vcrOptions(step.actions)) / 1000;
			advanceAt = Math.max(interval, impactAt + interval / 3);
		}
		this.current = { step, start: at, impactAt, advanceAt, impactDone: false };
		this.animating = true;
		this.inFlight = impactAt > 0;
	}

	private vcrOptions(shots: BattleRecordTokenAction[]): VcrAttackOptions {
		const first = shots.at(0);
		return {
			cell: this.cell,
			frameMs: vcrFrameMs(vcrSpeed(this.speed)),
			beamPower: first ? this.host.weaponPower?.(first.tokenNum, first.slot) : undefined,
			still: this.reducedMotion || !this.animations
		};
	}

	// what the board shows: the attack being played, or the one that led to the phase on the board
	private frame(): VcrFrame | undefined {
		const current = this.current;
		if (current) {
			const t = (this.clock - current.start) * 1000;
			if (t < 0 || (current.step.kind !== 'beam' && current.step.kind !== 'torpedo')) {
				return undefined;
			}
			return attackFrame(current.step.actions, t, this.vcrOptions(current.step.actions));
		}
		const shots = attackAtPhase(this.battle.actions, this.host.getPhase());
		return shots.length ? attackFrame(shots, Infinity, this.vcrOptions(shots)) : undefined;
	}

	// ---- frame loop ----

	private ensureLoop() {
		if (!this.raf && typeof window !== 'undefined') {
			this.lastTs = 0;
			this.raf = requestAnimationFrame(this.tick);
		}
	}

	private tick = (ts: number) => {
		this.raf = 0;
		const dt = this.lastTs ? Math.min(0.1, (ts - this.lastTs) / 1000) : 0;
		this.lastTs = ts;
		if (!this.frozen) {
			this.clock += dt;
			this.process();
		}
		this.draw();
		if (!this.frozen && (this.playing || this.current)) {
			this.raf = requestAnimationFrame(this.tick);
		}
	};

	private process() {
		const current = this.current;
		if (current) {
			const t = this.clock - current.start;
			if (!current.impactDone && t >= current.impactAt) {
				current.impactDone = true;
				this.inFlight = false;
				this.host.setPhase(current.step.toPhase);
			}
			if (t >= current.advanceAt) {
				const end = current.start + current.advanceAt;
				this.cancelCurrent();
				if (this.playing) {
					this.startStep(end);
					this.process();
				}
			}
		} else if (this.playing) {
			this.startStep(this.clock);
		}
	}

	private draw() {
		const ctx = this.ctx;
		if (!ctx) return;
		ctx.setTransform(this.dpr, 0, 0, this.dpr, 0, 0);
		ctx.clearRect(0, 0, this.width, this.height);
		const frame = this.frame();
		if (frame) {
			this.renderer.render(ctx, frame);
		}
	}
}
