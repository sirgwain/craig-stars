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
import Scanner from './Scanner.svelte';

const state = vi.hoisted(() => ({ context: undefined as unknown }));
vi.mock('#lib/services/GameContext.js', () => ({ getGameContext: () => state.context }));
// Exercise scanner gestures and hit testing without involving canvas rendering.
vi.mock('./ScannerCanvas.svelte', () => ({ default: () => {} }));
vi.mock('#lib/components/game/tooltips/ScannerContextPopup.js', () => ({
	onScannerContextPopup: vi.fn()
}));

afterEach(() => {
	vi.restoreAllMocks();
});

it.each(['pointerup', 'pointercancel'])(
	'does not re-enable the just-added guard after %s',
	async (release) => {
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
		const settings = writable(Object.assign(new PlayerSettings(), { addWaypoint: true }));
		const zoomTarget = writable({ mapObject: fleet.mapObject });
		state.context = {
			game: writable({ area: { x: 100, y: 100 } }),
			player: writable(player),
			universe: writable(universe),
			settings,
			zoomTarget,
			commandedFleet: writable(fleet),
			currentSelectedWaypointIndex: writable(1),
			highlightMapObject: vi.fn(),
			mostRecentMapObject: writable(undefined),
			selectedWaypoint: writable(fleet.fleetOrders.waypoints[1])
		};
		let finishAdd!: (added: boolean) => void;
		const pendingAdd = new Promise<boolean>((resolve) => {
			finishAdd = resolve;
		});
		const onAddWaypoint = vi.fn(() => pendingAdd);
		const onUpdateWaypointDest = vi.fn();
		const { container } = render(Scanner, {
			onAddWaypoint,
			onUpdateWaypointDest,
			onSelectWaypoint: vi.fn(),
			onSelectMapObject: vi.fn(),
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
		function pointer(type: string, x: number, y: number) {
			const [layerX, layerY] = transform.apply([x * 2, y * 2]);
			const event = new PointerEvent(type, { bubbles: true, pointerType: 'touch', button: 0 });
			Object.defineProperties(event, { layerX: { value: layerX }, layerY: { value: layerY } });
			overlay.dispatchEvent(event);
		}

		pointer('pointerdown', 20, 20);
		expect(onAddWaypoint).toHaveBeenCalledOnce();
		pointer('pointermove', 20, 20);
		pointer('pointermove', 30, 30);
		expect(onUpdateWaypointDest).not.toHaveBeenCalled();
		pointer(release, 20, 20);
		// Completing the old addition after pointer-up must not affect the next gesture.
		finishAdd(true);
		await pendingAdd;
		await tick();
		settings.update((s) => {
			s.addWaypoint = false;
			return s;
		});
		await tick();
		pointer('pointerdown', 20, 20);
		pointer('pointermove', 20, 20);
		pointer('pointermove', 30, 30);
		expect(onUpdateWaypointDest).toHaveBeenCalledWith(expect.anything(), false, false);
		pointer('pointerup', 30, 30);
		expect(onUpdateWaypointDest).toHaveBeenLastCalledWith(expect.anything(), false, true);
	}
);
