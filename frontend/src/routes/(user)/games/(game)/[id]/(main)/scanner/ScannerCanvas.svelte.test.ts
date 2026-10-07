import { afterEach, expect, it, vi } from 'vitest';
import { render } from 'vitest-browser-svelte';
import ScannerCanvas from './ScannerCanvas.svelte';

const state = vi.hoisted(() => ({ redraw: () => {} }));
vi.mock('./canvas/colors', () => ({
	readScannerColors: () => ({ background: '#000000', selection: '#ffffff' }),
	fixedScannerColors: {
		waypointLine: '#0000FF',
		waypointLineCommanded: '#00FFFF',
		yearTick: '#00FF00'
	}
}));
vi.mock('#lib/services/GameContext.js', async () => {
	const { writable } = await import('svelte/store');
	const { create } = await import('@bufbuild/protobuf');
	const { MinefieldSchema, MapObjectType } = await import('#lib/types/cs-proto.js');
	const { PlayerSettings } = await import('#lib/types/PlayerSettings.js');
	const minefields = [1, 2].map((playerNum) =>
		create(MinefieldSchema, {
			mapObject: {
				type: MapObjectType.MINEFIELD,
				num: playerNum,
				playerNum,
				position: { x: playerNum * 80, y: 80 }
			},
			numMines: 1600
		})
	);
	const context = {
		player: writable({ num: 1, isFriend: () => false }),
		universe: writable({
			allMinefields: minefields,
			planets: [],
			planetIntels: [],
			fleets: [],
			fleetIntels: [],
			mineralPackets: [],
			wormholes: [],
			mysteryTraders: [],
			salvages: [],
			getAllFleets: () => []
		}),
		settings: writable(new PlayerSettings()),
		commandedFleet: writable(undefined),
		commandedPlanet: writable(undefined),
		commandedMapObject: writable(undefined),
		selectedMapObject: writable(minefields[1]),
		highlightedMapObject: writable(undefined),
		selectedWaypoint: writable(undefined)
	};
	state.redraw = () =>
		context.settings.update((settings) => {
			settings.showNames = !settings.showNames;
			return settings;
		});
	return { getGameContext: () => context };
});

afterEach(() => vi.restoreAllMocks());

async function scanner() {
	const created = vi.spyOn(document, 'createElement');
	const { container } = render(ScannerCanvas, {
		transform: undefined,
		pixelsPerLightYear: { x: 1, y: 1 },
		minObjectZoom: 1,
		locate: undefined
	});
	// The component harness doesn't load the app's Tailwind layout. Fix the scanner size
	// so the canvas's own intrinsic height cannot resize its container in a loop.
	const bounds = container.firstElementChild as HTMLElement;
	bounds.style.width = '256px';
	bounds.style.height = '160px';
	const canvas = document.querySelector<HTMLCanvasElement>('[data-type="scanner-canvas"]')!;
	await expect.poll(() => canvas.dataset.view).toBeDefined();
	await new Promise((resolve) => requestAnimationFrame(() => requestAnimationFrame(resolve)));
	const ctx = canvas.getContext('2d')!;
	// Count the dots in each minefield's interior rather than its antialiased edges. Reading
	// back a canvas can switch Chrome's rasterizer, which slightly changes edge colors.
	const dots = () =>
		[80, 160].map((x, field) => {
			const dpr = window.devicePixelRatio;
			const data = ctx.getImageData((x - 16) * dpr, 64 * dpr, 32 * dpr, 32 * dpr).data;
			let count = 0;
			for (let i = 0; i < data.length; i += 4) {
				if (data[i + (field === 0 ? 2 : 0)] > 40 && data[i + 1] === 0) count++;
			}
			return count;
		});
	const before = dots();
	// Both the regular and darkened selected minefield must have a visible dot fill.
	expect(before[0]).toBeGreaterThan(100);
	expect(before[1]).toBeGreaterThan(100);
	// simulate the browser discarding canvas bitmaps while the phone sleeps
	const discard = () => {
		for (const result of created.mock.results) {
			const tile = result.value;
			if (tile instanceof HTMLCanvasElement && tile !== canvas) {
				tile.getContext('2d')!.reset();
			}
		}
		ctx.reset();
	};
	return { canvas, dots, before, discard };
}

it('rebuilds discarded minefield patterns on the next redraw', async () => {
	const { dots, before, discard } = await scanner();
	discard();
	state.redraw();
	await expect.poll(dots).toEqual(before);
});

it.each([
	[
		'contextrestored',
		(canvas: HTMLCanvasElement) => canvas.dispatchEvent(new Event('contextrestored'))
	],
	['visibilitychange', () => document.dispatchEvent(new Event('visibilitychange'))],
	['pageshow', () => window.dispatchEvent(new PageTransitionEvent('pageshow', { persisted: true }))]
])('redraws on %s', async (_, wake) => {
	const { canvas, dots, before, discard } = await scanner();
	discard();
	wake(canvas);
	await expect.poll(dots).toEqual(before);
});
