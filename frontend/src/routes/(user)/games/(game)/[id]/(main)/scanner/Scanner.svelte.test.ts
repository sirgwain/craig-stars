import { create } from '@bufbuild/protobuf';
import { zoomIdentity } from 'd3-zoom';
import { tick } from 'svelte';
import { writable } from 'svelte/store';
import { afterEach, expect, it, vi } from 'vitest';
import { render } from 'vitest-browser-svelte';
import { FleetSchema, MapObjectType, VectorSchema } from '#lib/types/cs-proto.js';
import { CommandedFleet } from '#lib/types/Fleet.js';
import { CommandedPlayer } from '#lib/types/Player.js';
import { PlayerSettings } from '#lib/types/PlayerSettings.js';
import { Universe } from '#lib/services/Universe.js';
import { onScannerContextPopup } from '#lib/components/game/tooltips/ScannerContextPopup.js';
import Scanner from './Scanner.svelte';

const state = vi.hoisted(() => ({
	context: undefined as unknown,
	canvasProps: undefined as
		{ transform?: { apply: (p: [number, number]) => [number, number] } } | undefined
}));
vi.mock('#lib/services/GameContext.js', () => ({ getGameContext: () => state.context }));
// Exercise scanner gestures and hit testing without involving canvas rendering.
vi.mock('./ScannerCanvas.svelte', () => ({
	default: (_anchor: unknown, props: typeof state.canvasProps) => {
		state.canvasProps = props;
	}
}));
vi.mock('#lib/components/game/tooltips/ScannerContextPopup.js', () => ({
	onScannerContextPopup: vi.fn()
}));

afterEach(() => {
	vi.restoreAllMocks();
	vi.clearAllMocks();
	vi.useRealTimers();
});

// A scanner commanding a fleet at (10, 10) whose selected waypoint is in deep space at (20, 20).
async function scannerWithSelectedWaypoint(addWaypoint = false) {
	const fleet = new CommandedFleet(
		create(FleetSchema, {
			mapObject: { type: MapObjectType.FLEET, playerNum: 1, num: 1, position: { x: 10, y: 10 } },
			fleetOrders: { waypoints: [{ position: { x: 10, y: 10 } }, { position: { x: 20, y: 20 } }] }
		})
	);
	const universe = new Universe();
	universe.fleets = [fleet];
	const player = new CommandedPlayer();
	player.num = 1;
	const zoomTarget = writable({ mapObject: fleet.mapObject });
	state.context = {
		game: writable({ area: { x: 100, y: 100 } }),
		player: writable(player),
		universe: writable(universe),
		settings: writable(Object.assign(new PlayerSettings(), { addWaypoint })),
		zoomTarget,
		commandedFleet: writable(fleet),
		currentSelectedWaypointIndex: writable(1),
		highlightMapObject: vi.fn(),
		mostRecentMapObject: writable(undefined),
		selectedWaypoint: writable(fleet.fleetOrders.waypoints[1])
	};
	const onSelectMapObject = vi.fn();
	const onUpdateWaypointDest = vi.fn();
	const onAddWaypoint = vi.fn(async () => true);
	const { container } = render(Scanner, {
		onAddWaypoint,
		onUpdateWaypointDest,
		onSelectWaypoint: vi.fn(),
		onSelectMapObject,
		onSetPacketDest: vi.fn(),
		onSetRouteDest: vi.fn()
	});
	const scanner = container.firstElementChild as HTMLElement;
	scanner.style.width = '200px';
	scanner.style.height = '200px';
	const root = scanner.firstElementChild as HTMLElement;
	root.style.width = '200px';
	root.style.height = '200px';
	window.dispatchEvent(new Event('resize'));
	await tick();
	zoomTarget.set({
		mapObject: { ...fleet.mapObject, position: create(VectorSchema, { x: 50, y: 50 }) }
	});
	await tick();
	const overlay = scanner.querySelector<HTMLElement>('[data-type="scanner-overlay"]')!;
	const transform = zoomIdentity.translate(-50, -50).scale(1.5);
	function pointer(type: string, x: number, y: number, pointerType = 'mouse') {
		const [layerX, layerY] = transform.apply([x * 2, y * 2]);
		const event = new PointerEvent(type, {
			bubbles: true,
			pointerType,
			button: 0,
			clientX: layerX,
			clientY: layerY
		});
		Object.defineProperties(event, { layerX: { value: layerX }, layerY: { value: layerY } });
		overlay.dispatchEvent(event);
	}
	return { pointer, onSelectMapObject, onUpdateWaypointDest, onAddWaypoint, scanner, root };
}

// The map is scaled to the scanner's shorter side, so turning a phone would otherwise leave the
// view somewhere else on the map.
it('keeps looking at the same part of the map when the scanner changes shape', async () => {
	const { scanner, root } = await scannerWithSelectedWaypoint();
	// a 200px square scanner shows the 100ly map at 2px per light year, centered on (50, 50)
	const [x, y] = state.canvasProps!.transform!.apply([100, 100]);
	expect(x).toBeCloseTo(100);
	expect(y).toBeCloseTo(100);

	// rotate to landscape
	for (const el of [scanner, root]) {
		el.style.width = '240px';
		el.style.height = '160px';
	}
	window.dispatchEvent(new Event('resize'));
	await tick();

	// now 1.6px per light year, and (50, 50) is still in the middle
	const [rx, ry] = state.canvasProps!.transform!.apply([80, 80]);
	expect(rx).toBeCloseTo(120);
	expect(ry).toBeCloseTo(80);
});

// Selecting what's under the selected waypoint can command a planet there, which would stop the
// waypoint being dragged. So a press on it waits to see whether it's a click or a drag.
it('treats a press on the selected waypoint as a click only when the pointer lifts', async () => {
	const { pointer, onSelectMapObject, onUpdateWaypointDest } = await scannerWithSelectedWaypoint();

	pointer('pointerdown', 20, 20);
	expect(onSelectMapObject).not.toHaveBeenCalled();
	pointer('pointerup', 20, 20);
	expect(onSelectMapObject).toHaveBeenCalledOnce();
	expect(onUpdateWaypointDest).not.toHaveBeenCalled();
});

it('drags the selected waypoint without selecting what is under it', async () => {
	const { pointer, onSelectMapObject, onUpdateWaypointDest } = await scannerWithSelectedWaypoint();

	pointer('pointerdown', 20, 20);
	pointer('pointermove', 30, 30);
	pointer('pointerup', 30, 30);
	expect(onUpdateWaypointDest).toHaveBeenLastCalledWith(expect.anything(), false, true);
	expect(onSelectMapObject).not.toHaveBeenCalled();
});

// Touch screens have no right button, so holding a finger still opens the context menu.
it('opens the context menu on a long press', async () => {
	const { pointer, onSelectMapObject } = await scannerWithSelectedWaypoint();
	vi.useFakeTimers({ toFake: ['setTimeout', 'clearTimeout'] });

	pointer('pointerdown', 20, 20, 'touch');
	vi.advanceTimersByTime(499);
	expect(onScannerContextPopup).not.toHaveBeenCalled();
	vi.advanceTimersByTime(1);
	expect(onScannerContextPopup).toHaveBeenCalledOnce();
	// lifting the finger afterwards isn't a click
	pointer('pointerup', 20, 20, 'touch');
	expect(onSelectMapObject).not.toHaveBeenCalled();
});

it('opens the context menu when a held finger wobbles a little', async () => {
	const { pointer } = await scannerWithSelectedWaypoint();
	vi.useFakeTimers({ toFake: ['setTimeout', 'clearTimeout'] });

	pointer('pointerdown', 20, 20, 'touch');
	// 2 light years is 6 pixels here, inside what a finger moves while held still
	pointer('pointermove', 22, 20, 'touch');
	vi.advanceTimersByTime(500);
	expect(onScannerContextPopup).toHaveBeenCalledOnce();
});

// The click sent when the finger lifts would otherwise close the menu as a click outside of it.
it('swallows the click that follows a long press', async () => {
	const { pointer } = await scannerWithSelectedWaypoint();
	vi.useFakeTimers({ toFake: ['setTimeout', 'clearTimeout'] });
	const seen = vi.fn();
	document.addEventListener('click', seen, true);

	pointer('pointerdown', 20, 20, 'touch');
	vi.advanceTimersByTime(500);
	pointer('pointerup', 20, 20, 'touch');
	document.body.dispatchEvent(new MouseEvent('click', { bubbles: true, cancelable: true }));
	expect(seen).not.toHaveBeenCalled();

	// only that one click
	document.body.dispatchEvent(new MouseEvent('click', { bubbles: true, cancelable: true }));
	expect(seen).toHaveBeenCalledOnce();
	document.removeEventListener('click', seen, true);
});

it.each([
	{ name: 'the finger moves', then: 'pointermove', x: 30, type: 'touch' },
	{ name: 'the finger lifts early', then: 'pointerup', x: 20, type: 'touch' },
	{ name: 'it is a mouse', then: 'pointermove', x: 20, type: 'mouse' }
])('does not open the context menu when $name', async ({ then, x, type }) => {
	const { pointer } = await scannerWithSelectedWaypoint();
	vi.useFakeTimers({ toFake: ['setTimeout', 'clearTimeout'] });

	pointer('pointerdown', 20, 20, type);
	pointer(then, x, x, type);
	vi.advanceTimersByTime(1000);
	expect(onScannerContextPopup).not.toHaveBeenCalled();
});

// In add waypoint mode a tap adds a waypoint. A drag still pans the map, so the waypoint is only
// added once the pointer lifts without having moved.
it('adds a waypoint on a tap in add waypoint mode', async () => {
	const { pointer, onAddWaypoint } = await scannerWithSelectedWaypoint(true);

	pointer('pointerdown', 60, 60, 'touch');
	expect(onAddWaypoint).not.toHaveBeenCalled();
	pointer('pointerup', 60, 60, 'touch');
	expect(onAddWaypoint).toHaveBeenCalledOnce();
});

// A held finger wobbles a pixel or two before it lifts, which is still a tap.
it('adds a waypoint on a tap that wobbles a little', async () => {
	const { pointer, onAddWaypoint } = await scannerWithSelectedWaypoint(true);

	pointer('pointerdown', 60, 60, 'touch');
	pointer('pointermove', 60.5, 60, 'touch');
	pointer('pointerup', 60.5, 60, 'touch');
	expect(onAddWaypoint).toHaveBeenCalledOnce();
});

it.each([
	{ name: 'the map is dragged', then: 'pointermove' },
	{ name: 'the touch is cancelled', then: 'pointercancel' }
])('does not add a waypoint when $name', async ({ then }) => {
	const { pointer, onAddWaypoint } = await scannerWithSelectedWaypoint(true);

	pointer('pointerdown', 60, 60, 'touch');
	pointer(then, 80, 80, 'touch');
	pointer('pointerup', 80, 80, 'touch');
	expect(onAddWaypoint).not.toHaveBeenCalled();
});

// Panning the map across a waypoint used to pick the waypoint up and drag it too.
it('does not drag a waypoint the pointer passes over while panning', async () => {
	const { pointer, onUpdateWaypointDest } = await scannerWithSelectedWaypoint();

	pointer('pointerdown', 60, 60, 'touch');
	pointer('pointermove', 40, 40, 'touch');
	pointer('pointermove', 20, 20, 'touch');
	pointer('pointermove', 10, 30, 'touch');
	pointer('pointerup', 10, 30, 'touch');
	expect(onUpdateWaypointDest).not.toHaveBeenCalled();
});

it('can drag a waypoint right after adding one', async () => {
	const { pointer, onAddWaypoint, onUpdateWaypointDest } = await scannerWithSelectedWaypoint(true);

	pointer('pointerdown', 60, 60, 'touch');
	pointer('pointerup', 60, 60, 'touch');
	expect(onAddWaypoint).toHaveBeenCalledOnce();
	await tick();

	pointer('pointerdown', 20, 20, 'touch');
	pointer('pointermove', 30, 30, 'touch');
	expect(onUpdateWaypointDest).toHaveBeenCalledWith(expect.anything(), false, false);
	pointer('pointerup', 30, 30, 'touch');
	expect(onUpdateWaypointDest).toHaveBeenLastCalledWith(expect.anything(), false, true);
	expect(onAddWaypoint).toHaveBeenCalledOnce();
});

it('measures with the tape from a mapobject to the pointer', async () => {
	const fleet = new CommandedFleet(
		create(FleetSchema, {
			mapObject: { type: MapObjectType.FLEET, playerNum: 1, num: 1, position: { x: 10, y: 10 } },
			fleetOrders: { waypoints: [{ position: { x: 10, y: 10 } }] }
		})
	);
	const universe = new Universe();
	universe.fleets = [fleet];
	const player = new CommandedPlayer();
	player.num = 1;
	const zoomTarget = writable({ mapObject: fleet.mapObject });
	state.context = {
		game: writable({ area: { x: 100, y: 100 } }),
		player: writable(player),
		universe: writable(universe),
		settings: writable(new PlayerSettings()),
		zoomTarget,
		commandedFleet: writable(fleet),
		currentSelectedWaypointIndex: writable(0),
		highlightMapObject: vi.fn(),
		mostRecentMapObject: writable(undefined),
		selectedWaypoint: writable(fleet.fleetOrders.waypoints[0])
	};
	const onSelectMapObject = vi.fn();
	const { container } = render(Scanner, {
		onAddWaypoint: vi.fn(),
		onUpdateWaypointDest: vi.fn(),
		onSelectWaypoint: vi.fn(),
		onSelectMapObject,
		onSetPacketDest: vi.fn(),
		onSetRouteDest: vi.fn()
	});
	const scanner = container.firstElementChild as HTMLElement;
	scanner.style.width = '200px';
	scanner.style.height = '200px';
	const root = scanner.firstElementChild as HTMLElement;
	root.style.width = '200px';
	root.style.height = '200px';
	window.dispatchEvent(new Event('resize'));
	await tick();
	zoomTarget.set({
		mapObject: { ...fleet.mapObject, position: create(VectorSchema, { x: 50, y: 50 }) }
	});
	await tick();
	const overlay = scanner.querySelector<HTMLElement>('[data-type="scanner-overlay"]')!;
	const transform = zoomIdentity.translate(-50, -50).scale(1.5);
	function pointer(type: string, x: number, y: number, init: PointerEventInit) {
		const [layerX, layerY] = transform.apply([x * 2, y * 2]);
		const event = new PointerEvent(type, { bubbles: true, cancelable: true, ...init });
		Object.defineProperties(event, { layerX: { value: layerX }, layerY: { value: layerY } });
		overlay.dispatchEvent(event);
		return event;
	}
	const status = () => scanner.querySelector('[data-type="tape-status"]')?.textContent;

	// shift + right button, starting on the fleet
	pointer('pointerdown', 11, 10, { button: 2, shiftKey: true });
	await tick();
	expect(status()).toBeUndefined();
	// a small move doesn't start the tape
	pointer('pointermove', 12, 11, { shiftKey: true });
	await tick();
	expect(status()).toBeUndefined();
	pointer('pointermove', 40, 50, { shiftKey: true });
	await tick();
	expect(status()).toContain('Deep Space');
	expect(status()).toContain('X: 40');
	expect(status()).toContain('Y: 50');
	expect(status()).toContain('50.00 light years');
	const line = scanner.querySelector('[data-type="tape"] line')!;
	const [x1, y1] = transform.apply([20, 20]);
	expect(Number(line.getAttribute('x1'))).toBeCloseTo(x1);
	expect(Number(line.getAttribute('y1'))).toBeCloseTo(y1);

	pointer('pointerup', 40, 50, { button: 2, shiftKey: true });
	await tick();
	expect(status()).toBeUndefined();
	expect(onSelectMapObject).not.toHaveBeenCalled();

	// the middle button measures too
	pointer('pointerdown', 40, 50, { button: 1 });
	pointer('pointermove', 70, 90, {});
	await tick();
	expect(status()).toContain('50.00 light years');
	pointer('pointerup', 70, 90, { button: 1 });
	await tick();
	expect(status()).toBeUndefined();
});
