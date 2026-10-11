<script lang="ts">
	import { clickOutside } from '#lib/clickOutside.js';
	import { onScannerContextPopup } from '#lib/components/game/tooltips/ScannerContextPopup.js';
	import type { SelectWaypointProps } from '#lib/services/Events.js';
	import { getGameContext } from '#lib/services/GameContext.js';
	import { clamp } from '#lib/services/Math.js';
	import { None } from '#lib/types/Consts.js';
	import {
		MapObjectType,
		VectorSchema,
		WaypointDestSchema,
		WaypointSchema,
		type Fleet,
		type Vector,
		type Waypoint,
		type WaypointDest
	} from '#lib/types/cs-proto.js';
	import { filterFleet } from '#lib/types/Filter.js';
	import { emptyMapObject, type MapObjectLike, type Position } from '#lib/types/MapObject.js';
	import { distance, emptyVector, equal } from '#lib/types/Vector.js';
	import { clone, create } from '@bufbuild/protobuf';
	import { scaleLinear } from 'd3-scale';
	import { select } from 'd3-selection';
	import { zoom, ZoomTransform, type D3ZoomEvent, type ZoomBehavior } from 'd3-zoom';
	import hotkeys from 'hotkeys-js';
	import { onDestroy, onMount } from 'svelte';
	import { derived as derivedStore } from 'svelte/store';
	import MapObjectQuadTreeFinder, { type FinderEvent } from './MapObjectQuadTreeFinder.svelte';
	import { getVisibleMinefields } from './canvas/scene';
	import ScannerCanvas from './ScannerCanvas.svelte';

	const {
		game,
		player,
		universe,
		settings,
		commandedFleet,
		currentSelectedWaypointIndex,
		highlightMapObject,
		mostRecentMapObject,
		selectedWaypoint,
		zoomTarget
	} = getGameContext();

	type Props = {
		onAddWaypoint: (dest: WaypointDest, fastestWaypoint: boolean) => Promise<boolean>;
		onUpdateWaypointDest: (dest: WaypointDest, fastestWaypoint: boolean, done: boolean) => void;
		onSelectMapObject: (mo: MapObjectLike) => void;
		onSetPacketDest: (mo: MapObjectLike) => void;
		onSetRouteDest: (mo: MapObjectLike) => void;
	} & SelectWaypointProps;

	let {
		onSelectWaypoint,
		onAddWaypoint,
		onUpdateWaypointDest,
		onSelectMapObject,
		onSetPacketDest,
		onSetRouteDest
	}: Props = $props();

	const aspectRatio = $game.area.x / $game.area.y;
	const padding = 20; // 20 px, used in zooming

	let root: HTMLElement | undefined = $state();
	let rect: HTMLDivElement | undefined = $state();
	let clientRect = $state({ width: 100, height: 100 });

	// compute scales, derived from clientWidth/height
	let scaler = $derived({
		x: scaleLinear().range(xRange(clientRect.width, clientRect.height)).domain([0, $game.area.x]),
		y: scaleLinear().range(yRange(clientRect.width, clientRect.height)).domain([0, $game.area.y])
	});

	let transform: ZoomTransform | undefined = $state();
	let showLocator = $state(false);
	let shouldAddWaypoint = $state(false);
	let waypointHighlighted = $state(false);

	// css pixels per light year, before zooming
	let pixelsPerLightYear = $derived({ x: scaler.x(1) - scaler.x(0), y: scaler.y(1) - scaler.y(0) });

	// our map scales for .75 to 10x, but the icons for the planets and fleets shrink below 2x
	const minZoom = 0.75;
	const maxZoom = 10;
	const minObjectZoom = 2;
	let scale = $state(1.5); // default zoom

	// zoom state that changes but doesn't cause a reaction
	let zooming = false;
	let fastestWaypoint = false;

	let pointerDown = false;
	let draggingWaypoint = false;
	// the commanded fleet's waypoints from before the drag, drawn so the player can see what they're changing
	let dragOrigin: Waypoint[] | undefined = $state();
	let dragAndZoomEnabled = false;

	// set to true if we are moving a waypoint to a position rather than a target
	// this is enabled when the shift key is held
	let positionWaypoint = false;

	// In add waypoint mode a tap adds a waypoint and a drag still pans the map, so the waypoint is
	// added when the pointer lifts without having moved.
	let pendingAdd: { found: MapObjectLike | undefined; position: Position } | undefined;

	// A press on the commanded fleet's waypoint is either the start of a drag or a click. We only know
	// which once the pointer moves or lifts, so the click waits for pointer up.
	// a finger wobbles more than a mouse while it's held still
	const dragMinDistance = 3;
	const touchMinDistance = 10;
	let pressedAt: { x: number; y: number } | undefined;
	let pressedWaypoint: Waypoint | undefined;
	let pendingSelect: MapObjectLike | undefined;

	function clientPoint(event: FinderEvent['event']): { x: number; y: number } | undefined {
		if ('clientX' in event) {
			return { x: event.clientX, y: event.clientY };
		}
		const touch = event.touches[0] ?? event.changedTouches[0];
		return touch ? { x: touch.clientX, y: touch.clientY } : undefined;
	}

	// true once the pointer has moved far enough from where it went down to be a drag, not a click
	function movedSincePress(event: FinderEvent['event']): boolean {
		const point = clientPoint(event);
		if (!pressedAt || !point) {
			return true;
		}
		const mouse = event instanceof PointerEvent && event.pointerType === 'mouse';
		const minDistance = mouse ? dragMinDistance : touchMinDistance;
		return Math.hypot(point.x - pressedAt.x, point.y - pressedAt.y) >= minDistance;
	}

	function clearPress() {
		pressedAt = undefined;
		pressedWaypoint = undefined;
		pendingSelect = undefined;
		pendingAdd = undefined;
		clearTimeout(longPressTimer);
		longPressTimer = undefined;
	}

	// Holding a finger still on a mapobject opens its context menu, like a right click with a mouse.
	const longPressMs = 500;
	let longPressTimer: ReturnType<typeof setTimeout> | undefined;
	// some browsers also send a contextmenu event for a long press, ignore it if we already handled it
	let longPressed = false;

	function startLongPress(event: FinderEvent['event'], found: MapObjectLike) {
		if (!(event instanceof PointerEvent) || event.pointerType === 'mouse') {
			return;
		}
		clearTimeout(longPressTimer);
		longPressTimer = setTimeout(() => {
			longPressTimer = undefined;
			// A touch always starts a pan, so we can't go by whether the scanner is zooming. The timer
			// is cancelled if the finger actually moves.
			if (!pointerDown || draggingWaypoint) {
				return;
			}
			// this press was for the menu, not a click or a drag
			longPressed = true;
			pendingSelect = undefined;
			pressedWaypoint = undefined;
			swallowNextClick();
			onScannerContextPopup(event, found.mapObject?.position ?? emptyVector());
		}, longPressMs);
	}

	// Lifting the finger after a long press sends a click, which would close the menu we just
	// opened as a click outside of it. Stop that one click before anything else sees it.
	let swallowClick: ((e: MouseEvent) => void) | undefined;
	function swallowNextClick() {
		stopSwallowingClick();
		swallowClick = (e: MouseEvent) => {
			e.preventDefault();
			e.stopPropagation();
			stopSwallowingClick();
		};
		window.addEventListener('click', swallowClick, true);
	}

	function stopSwallowingClick() {
		if (swallowClick) {
			window.removeEventListener('click', swallowClick, true);
			swallowClick = undefined;
		}
	}

	// the measuring tape is dragged out with the middle button, or the right button while shift is held.
	// It measures from where the drag started to the pointer, snapping to mapobjects at either end
	const tapeMinDistance = 3;
	type Tape = { start: Vector; end: Vector; name: string };
	let tapeStart: Vector | undefined;
	let tapeContextMenu = false;
	let tape: Tape | undefined = $state();
	let tapeLine = $derived.by(() => {
		if (!tape || !transform) {
			return undefined;
		}
		const [x1, y1] = transform.apply([scaler.x(tape.start.x), scaler.y(tape.start.y)]);
		const [x2, y2] = transform.apply([scaler.x(tape.end.x), scaler.y(tape.end.y)]);
		return { x1, y1, x2, y2 };
	});

	// zoomBehavior is based on clientWidth/height
	let zoomBehavior: ZoomBehavior<HTMLElement, unknown> = $derived(
		zoom<HTMLElement, unknown>()
			.extent([
				[0, 0],
				[clientRect.width, clientRect.height]
			])
			.scaleExtent([minZoom, maxZoom])
			.translateExtent([
				[-20, -20],
				[clientRect.width + padding, clientRect.height + padding]
			])
			.on('zoom', handleZoom)
			.on('start', handleZoomStart)
			.on('end', handleZoomEnd)
	);

	// enable drag and zoom, but disable dblclick zoom events
	function enableDragAndZoom() {
		if (!root || dragAndZoomEnabled) {
			return;
		}
		select(root).call(zoomBehavior).on('dblclick.zoom', null);
		dragAndZoomEnabled = true;
	}

	// disable drag and zoom temporarily
	function disableDragAndZoom() {
		if (!root || !dragAndZoomEnabled) {
			return;
		}
		select(root).on('.zoom', null);
		dragAndZoomEnabled = false;
		zooming = false;
	}

	function xRange(clientWidth: number, clientHeight: number) {
		if (aspectRatio > 1 && clientHeight > clientWidth) {
			// tall skinny viewport, wide map, so fully expand on the x
			// but shrink up height
			return [0, clientWidth];
		} else if (aspectRatio > 1 && clientWidth > clientHeight) {
			// wide viewport, wide map, so fully expand on the y
			// but shrink up width
			return [0, clientHeight * aspectRatio];
		}
		return [0, Math.min(clientWidth, clientHeight)];
	}

	function yRange(clientWidth: number, clientHeight: number) {
		if (aspectRatio > 1 && clientHeight > clientWidth) {
			// tall skinny viewport, wide map, so fully expand on the x
			// but shrink up height
			return [0, clientWidth / aspectRatio];
		} else if (aspectRatio > 1 && clientWidth > clientHeight) {
			// wide viewport, wide map, so fully expand on the y
			return [0, clientHeight];
		}
		return [0, Math.min(clientWidth, clientHeight)];
	}

	// update clientWidth/height on resize
	// The map is scaled to fit the scanner, so when the scanner changes shape, like a phone being
	// rotated, we re-center on the part of the map the player was looking at.
	function handleResize() {
		const next = rect?.getBoundingClientRect() ?? { width: 100, height: 100 };
		if (next.width === clientRect.width && next.height === clientRect.height) {
			return;
		}
		const center = transform
			? create(VectorSchema, {
					x: scaler.x.invert(transform.invertX(clientRect.width / 2)),
					y: scaler.y.invert(transform.invertY(clientRect.height / 2))
				})
			: undefined;
		clientRect = next;

		// the zoom behavior is sized to the scanner, so swap in the new one
		if (dragAndZoomEnabled) {
			disableDragAndZoom();
			enableDragAndZoom();
		}
		if (center) {
			translateViewport(center);
		}
	}

	// Only the window changing size re-fits the map. The scanner also changes size when the phone's
	// command drawer opens, and the map shouldn't rescale under the player for that.
	function onWindowResize() {
		handleResize();
		// a rotating phone can fire its resize before the layout settles
		requestAnimationFrame(handleResize);
	}

	function handleKeyDown(e: KeyboardEvent) {
		// add a waypoint if we are currently commanding a fleet and we didn't just click
		// on the fleet
		shouldAddWaypoint = !!$commandedFleet && (e.shiftKey || e.metaKey);
		fastestWaypoint = !!$commandedFleet && e.metaKey;

		switch (e.key) {
			case '+':
			case '=':
				zoomViewport(clamp(scale + 1, minZoom, maxZoom));
				break;
			case '-':
			case '_':
				zoomViewport(clamp(scale - 1, minZoom, maxZoom));
				break;
		}
	}

	function handleKeyUp(e: KeyboardEvent) {
		// add a waypoint if we are currently commanding a fleet and we didn't just click
		// on the fleet
		shouldAddWaypoint = !!$commandedFleet && (e.shiftKey || e.metaKey);
		fastestWaypoint = !!$commandedFleet && e.metaKey;
	}

	function showTargetLocation() {
		showLocator = true;
		setTimeout(() => (showLocator = false), 500);
	}

	function handleZoom(e: D3ZoomEvent<HTMLElement, unknown>) {
		const zoomed = e.transform.k !== scale;
		transform = e.transform;
		scale = transform.k;
		if (zoomed) {
			// a pinch isn't a long press or a tap. A pan is left to how far the finger has moved,
			// because a finger held still always nudges the map a pixel or two
			clearTimeout(longPressTimer);
			longPressTimer = undefined;
			pendingAdd = undefined;
		}
	}

	function handleZoomStart() {
		zooming = true;
	}

	function handleZoomEnd() {
		zooming = false;
	}

	// translate/zoom the display to a point on the map
	function translateViewport(position: Position) {
		if (!root) {
			return;
		}

		select(root).call(zoomBehavior.scaleTo, scale);
		const scaled: Vector = create(VectorSchema, {
			x: scaler.x(Number(position.x ?? 0)),
			y: scaler.y(Number(position.y ?? 0))
		});
		select(root)
			.call(zoomBehavior.translateTo, scaled.x, scaled.y)
			.call(zoomBehavior.scaleTo, scale);
	}

	// zoom the viewport to a specific scale
	function zoomViewport(scaleTo: number) {
		if (!root) {
			return;
		}
		select(root).call(zoomBehavior.scaleTo, scaleTo);
		const scaled: Vector = create(VectorSchema, {
			x: scaler.x(Number($zoomTarget?.mapObject?.position?.x ?? 0)),
			y: scaler.y(Number($zoomTarget?.mapObject?.position?.y ?? 0))
		});
		select(root)
			.call(zoomBehavior.translateTo, scaled.x, scaled.y)
			.call(zoomBehavior.scaleTo, scale);
	}

	function isTapeButton(event: FinderEvent['event']): boolean {
		return (
			event instanceof PointerEvent &&
			(event.button === 1 || (event.button === 2 && event.shiftKey))
		);
	}

	function tapePosition(position: Position, mo: MapObjectLike | undefined): Vector {
		const snapped = mo?.mapObject?.position ?? position;
		return create(VectorSchema, { x: Number(snapped.x ?? 0), y: Number(snapped.y ?? 0) });
	}

	function tapeMove(position: Position, mo: MapObjectLike | undefined) {
		if (!tapeStart) {
			return;
		}
		const end = tapePosition(position, mo);
		if (
			!tape &&
			Math.abs(end.x - tapeStart.x) < tapeMinDistance &&
			Math.abs(end.y - tapeStart.y) < tapeMinDistance
		) {
			// like a click, the tape doesn't start until the pointer moves away
			return;
		}
		tape = { start: tapeStart, end, name: mo?.mapObject?.name || 'Deep Space' };
	}

	function tapeDone() {
		tapeStart = undefined;
		tape = undefined;
	}

	// turn off dragging
	function onContextMenu(e: FinderEvent) {
		const { event, found } = e;

		if (tapeContextMenu) {
			// this right click was the measuring tape
			tapeContextMenu = false;
			return;
		}

		if (longPressed) {
			// this long press already opened the menu
			longPressed = false;
			return;
		}

		if (found && event instanceof MouseEvent) {
			onScannerContextPopup(event, found.mapObject?.position ?? emptyVector());
		}
	}

	// as the pointer moves, find the items it is under
	function onPointerMove(e: FinderEvent) {
		const { event, found, position } = e;

		highlightMapObject(found);

		if (tapeStart) {
			tapeMove(position, found);
			return;
		}

		// check if we are over the commanded fleet's waypoint
		const fleetWaypoint =
			found &&
			$commandedFleet &&
			$commandedFleet.fleetOrders.waypoints
				.slice(1)
				.find((wp) =>
					equal(wp.position ?? emptyVector(), found.mapObject?.position ?? emptyVector())
				);
		waypointHighlighted = !!fleetWaypoint;
		if (waypointHighlighted) {
			if (dragAndZoomEnabled) {
				disableDragAndZoom();
			}
		} else {
			if (!draggingWaypoint && !dragAndZoomEnabled) {
				enableDragAndZoom();
			}
		}

		if ((longPressTimer || pendingAdd) && movedSincePress(event)) {
			// the finger moved, so this is a pan or a drag
			clearTimeout(longPressTimer);
			longPressTimer = undefined;
			pendingAdd = undefined;
		}

		// check if we started a waypoint drag
		// we only
		// * start dragging once
		// * if the pointer is down
		// * on the waypoint we pressed. Passing over a waypoint while panning the map doesn't pick it up
		// * once the pointer has moved, so a click on a waypoint isn't a drag
		// * if we have a commanded fleet
		const dragWaypoint = pressedWaypoint;
		if (
			!draggingWaypoint &&
			pointerDown &&
			dragWaypoint &&
			$commandedFleet &&
			movedSincePress(event)
		) {
			draggingWaypoint = true;
			pendingSelect = undefined;
			pendingAdd = undefined;
			dragOrigin = $commandedFleet.fleetOrders.waypoints.map((wp) => clone(WaypointSchema, wp));
			onSelectWaypoint?.({ fleet: $commandedFleet, waypoint: dragWaypoint });
		}

		if (draggingWaypoint && !zooming) {
			positionWaypoint = event.shiftKey;
			dragWaypointMove(position, found);
		}
	}

	async function onPointerDown(e: FinderEvent) {
		const { event, found, position } = e;

		if (isTapeButton(event)) {
			// the middle button would otherwise start scrolling the page
			event.preventDefault();
			tapeStart = tapePosition(position, found);
			tapeContextMenu = (event as PointerEvent).button === 2;
			try {
				// keep measuring if the pointer leaves the scanner
				(event.target as Element).setPointerCapture((event as PointerEvent).pointerId);
			} catch {
				// the pointer is already gone
			}
			return;
		}

		if (event instanceof MouseEvent && event.button != 0) {
			// we only care about the first button
			return;
		}

		if (
			found?.mapObject?.type === MapObjectType.FLEET &&
			!filterFleet($player, found as Fleet, $settings, $universe)
		) {
			// this object we clicked is filtered out, don't do anything
			return;
		}

		pointerDown = true;
		longPressed = false;
		pressedAt = clientPoint(event);
		const addingWaypoint = !!$commandedFleet && (shouldAddWaypoint || $settings.addWaypoint);
		if (found && !addingWaypoint) {
			startLongPress(event, found);
		}

		// check if we pressed the commanded fleet's waypoint, it might be dragged
		const fleetWaypoint =
			found &&
			$commandedFleet &&
			$commandedFleet.fleetOrders.waypoints
				.slice(1)
				.find((wp) =>
					equal(wp.position ?? emptyVector(), found.mapObject?.position ?? emptyVector())
				);
		pressedWaypoint = fleetWaypoint;
		if (pressedWaypoint) {
			// this press can drag the waypoint, so it mustn't pan the map as well
			disableDragAndZoom();
		}

		if (addingWaypoint) {
			// wait for the pointer to lift, this could be the start of a pan or a drag
			pendingAdd = { found, position };
		} else if (found) {
			if (fleetWaypoint && $selectedWaypoint != fleetWaypoint) {
				onSelectWaypoint?.({ fleet: $commandedFleet, waypoint: fleetWaypoint });
			} else if (fleetWaypoint) {
				// Selecting what's here can command a planet, which would end the drag before it starts.
				// Wait to see if this press is a drag.
				pendingSelect = found;
			} else {
				mapObjectSelected(found);
			}
		}
	}

	// turn off dragging
	function onPointerUp(e: FinderEvent) {
		const { event, found, position } = e;

		if (tapeStart) {
			tapeDone();
			return;
		}

		if (event instanceof MouseEvent && event.button != 0) {
			// we only care about the first button
			return;
		}

		// the map can pan again, whatever this press was
		enableDragAndZoom();

		if (draggingWaypoint) {
			dragWaypointDone(position, found);
		} else if (pendingAdd && !movedSincePress(event)) {
			// a tap in add waypoint mode
			addWaypoint(pendingAdd.found, pendingAdd.position);
		} else if (pendingSelect) {
			// the press on the selected waypoint was a click after all
			mapObjectSelected(pendingSelect);
		}
		if (longPressed) {
			// the click from this lift comes right after, if the browser sends one at all
			setTimeout(stopSwallowingClick, 400);
		}
		clearPress();
		draggingWaypoint = false;
		dragOrigin = undefined;
		pointerDown = false;
	}

	function onPointerCancel() {
		tapeDone();
		clearPress();
		draggingWaypoint = false;
		dragOrigin = undefined;
		pointerDown = false;
		enableDragAndZoom();
	}

	// move the selected waypoint around snapping to targets
	function dragWaypointMove(position: Position, mo: MapObjectLike | undefined) {
		if (!($selectedWaypoint && $currentSelectedWaypointIndex && $commandedFleet)) {
			return;
		}

		const dest = create(WaypointDestSchema, {
			mo: positionWaypoint || !mo?.mapObject ? emptyMapObject() : mo.mapObject,
			position: { x: Number(position.x), y: Number(position.y) }
		});
		onUpdateWaypointDest(dest, fastestWaypoint, false);
	}

	async function dragWaypointDone(position: Position, mo: MapObjectLike | undefined) {
		// reset waypoint dragging
		if ($selectedWaypoint && $commandedFleet && draggingWaypoint) {
			// don't move the waypoint to any adjacent waypoints (same check as dragWaypointMove)
			if (mo && !positionWaypoint) {
				const index = $commandedFleet.fleetOrders.waypoints.findIndex((wp) =>
					equal(wp.position ?? emptyVector(), mo.mapObject?.position ?? emptyVector())
				);
				if (
					index == $currentSelectedWaypointIndex - 1 ||
					index == $currentSelectedWaypointIndex + 1
				) {
					// Don't call updateWaypoint - let the Go code handle the deletion logic
					// Instead, still call it but the Go code will return the correct enum
				}
			}

			const dest = create(WaypointDestSchema, {
				mo: positionWaypoint || !mo?.mapObject ? emptyMapObject() : mo.mapObject,
				position: { x: Number(position.x), y: Number(position.y) }
			});
			onUpdateWaypointDest(dest, fastestWaypoint, true);
		}
	}

	// disable add waypoint mode when the user clicks outside the scanner
	function disableAddWaypointMode(event: MouseEvent) {
		// ignore clicks on the add-waypoint toolbar button
		const elem = event.target as Element;
		if (elem.closest('#add-waypoint, #delete-waypoint')) {
			// deleting a waypoint shouldn't end add waypoint mode either
			return;
		}
		if ($settings.addWaypoint) {
			$settings.addWaypoint = false;
			$settings.fastestWaypoint = false;
		}
	}

	// if the shift key is held, add a waypoint instead of selecting a mapobject
	async function addWaypoint(mo: MapObjectLike | undefined, position: Position): Promise<boolean> {
		// This is only called for a tap. We can't go by whether the scanner is zooming, a touch always
		// starts a pan and it hasn't ended yet when the finger lifts.
		if (!$commandedFleet?.fleetOrders.waypoints) {
			return false;
		}

		// for add waypoints, we always snap to planet because the "drag" and "add waypoint button" keys (shift) are the same
		const dest = create(WaypointDestSchema, {
			mo: mo?.mapObject ?? emptyMapObject(),
			position: !mo?.mapObject ? { x: Number(position.x), y: Number(position.y) } : undefined
		});
		await onAddWaypoint(dest, fastestWaypoint);
		return true;
	}
	/**
	 * When a mapobject is selected we go through a few steps.
	 * - We select it if it's a new selection
	 * - We cycle through our commandable objects at the same location if we own an object there
	 * @param mo
	 */
	function mapObjectSelected(mo: MapObjectLike) {
		if ($settings.setPacketDest) {
			onSetPacketDest(mo);
		} else if ($settings.setRouteDest) {
			onSetRouteDest(mo);
		} else {
			onSelectMapObject(mo);
		}
	}

	// handle zoom in/out
	// this behavior controls how the zoom behaves
	// below we handle zooming events by updating a transform
	onMount(() => {
		clientRect = rect?.getBoundingClientRect() ?? { width: 100, height: 100 };

		// setup zoom and translate to the zoomTarget
		translateViewport($zoomTarget?.mapObject?.position ?? create(VectorSchema));
		enableDragAndZoom();

		// setup asubscriber to draw the target X and move the viewport to a new target
		// when the zoomTarget changes
		const unsubscribeZoomTarget = zoomTarget.subscribe((target) => {
			if (target) {
				translateViewport(target.mapObject?.position ?? create(VectorSchema));
				showTargetLocation();
			}
		});

		// bind the v key to show the target with X
		hotkeys('v', 'root', showTargetLocation);

		return () => {
			unsubscribeZoomTarget();
		};
	});

	onDestroy(() => {
		clearTimeout(longPressTimer);
		stopSwallowingClick();
		hotkeys.unbind('v', 'root', showTargetLocation);
	});

	// data used by the scanner is derived from the universe/commandedFleet stores
	// when they update, our data updates
	const data = derivedStore([universe, commandedFleet, player, settings], ([u, f, p, s]) => [
		// add mapobject waypoints
		...(f?.getWaypointMapObjects(u) || []),
		...u.getAllFleets().filter((f) => f.orbitingPlanetNum === None),
		...u.mysteryTraders,
		...u.salvages,
		...u.wormholes,
		...u.mineralPackets,
		// minefields that aren't shown can't be picked
		...getVisibleMinefields(u, p, s),
		...u.allPlanets
	]);

	// all our data in LayerCake are mapObjects/waypoints. Add this custom getter to get the
	// x/y coords of a mapobject or waypoint
</script>

<svelte:window onresize={onWindowResize} onkeydown={handleKeyDown} onkeyup={handleKeyUp} />

<div
	class:cursor-grab={waypointHighlighted}
	class:cursor-cell={shouldAddWaypoint ||
		(!!$commandedFleet && $settings.addWaypoint) ||
		$settings.setPacketDest ||
		$settings.setRouteDest}
	class={`grow bg-black overflow-hidden p-[${padding}px] select-none touch-none overscroll-contain`}
	bind:this={rect}
	use:clickOutside={disableAddWaypointMode}
>
	<div class="relative size-full" bind:this={root}>
		<ScannerCanvas
			{transform}
			{pixelsPerLightYear}
			{minObjectZoom}
			locate={showLocator ? $mostRecentMapObject : undefined}
			{dragOrigin}
		/>
		{#if transform}
			<MapObjectQuadTreeFinder
				data={$data}
				{pixelsPerLightYear}
				width={clientRect.width}
				height={clientRect.height}
				contextmenu={onContextMenu}
				pointermove={onPointerMove}
				pointerdown={onPointerDown}
				pointerup={onPointerUp}
				pointercancel={onPointerCancel}
				touchmove={onPointerMove}
				searchRadius={20}
				{transform}
			/>
		{/if}
		{#if tape && tapeLine}
			<svg class="absolute inset-0 size-full z-20 pointer-events-none" data-type="tape">
				<line {...tapeLine} class="stroke-white" stroke-width="1" />
			</svg>
			<div
				class="absolute bottom-0 inset-x-0 z-20 pointer-events-none flex flex-row flex-wrap gap-x-4 px-2 py-0.5 bg-base-200 text-sm"
				data-type="tape-status"
			>
				<span class="font-semibold">{tape.name}</span>
				<span>X: {tape.end.x}</span>
				<span>Y: {tape.end.y}</span>
				<span class="ml-auto">{distance(tape.start, tape.end).toFixed(2)} light years</span>
			</div>
		{/if}
	</div>
</div>
