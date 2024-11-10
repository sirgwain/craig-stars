<script lang="ts">
	import { clickOutside } from '$lib/clickOutside';
	import { onScannerContextPopup } from '$lib/components/game/tooltips/ScannerContextPopup.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { clamp } from '$lib/services/Math';
	import { None } from '$lib/types/Constants';
	import { filterFleet } from '$lib/types/Filter';
	import { type Fleet, type Waypoint, type WaypointDest } from '$lib/types/Fleet';
	import { MapObjectType, type MapObject } from '$lib/types/MapObject';
	import { emptyVector, equal, type Vector } from '$lib/types/Vector';
	import { scaleLinear } from 'd3-scale';
	import { select } from 'd3-selection';
	import { ZoomTransform, zoom, type D3ZoomEvent, type ZoomBehavior } from 'd3-zoom';
	import hotkeys from 'hotkeys-js';
	import { Html, LayerCake, Svg } from 'layercake';
	import { onDestroy, onMount, setContext } from 'svelte';
	import { derived as derivedStore, writable } from 'svelte/store';
	import MapObjectQuadTreeFinder, { type FinderEvent } from './MapObjectQuadTreeFinder.svelte';
	import ScannerFleets from './ScannerFleets.svelte';
	import ScannerMapObjectLocation from './ScannerMapObjectLocation.svelte';
	import ScannerMineFieldPattern from './ScannerMineFieldPattern.svelte';
	import ScannerMineFields from './ScannerMineFields.svelte';
	import ScannerMineralPackets from './ScannerMineralPackets.svelte';
	import ScannerMysteryTraders from './ScannerMysteryTraders.svelte';
	import ScannerNames from './ScannerNames.svelte';
	import ScannerPacketDests from './ScannerPacketDests.svelte';
	import ScannerPlanets from './ScannerPlanets.svelte';
	import ScannerSalvages from './ScannerSalvages.svelte';
	import ScannerScanners from './ScannerScanners.svelte';
	import ScannerWarpLine from './ScannerWarpLine.svelte';
	import ScannerWaypoints from './ScannerWaypoints.svelte';
	import ScannerWormholeLinks from './ScannerWormholeLinks.svelte';
	import ScannerWormholes from './ScannerWormholes.svelte';
	import SelectedMapObject from './SelectedMapObject.svelte';

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
		onSelectWaypoint: (wp: Waypoint) => void;
		onAddWaypoint: (dest: WaypointDest, fastestWaypoint: boolean) => Promise<boolean>;
		onUpdateWaypoint: (dest: WaypointDest, fastestWaypoint: boolean, done: boolean) => void;
		onSelectMapObject: (mo: MapObject) => void;
		onSetPacketDest: (mo: MapObject) => void;
	};

	let {
		onSelectWaypoint,
		onAddWaypoint,
		onUpdateWaypoint,
		onSelectMapObject,
		onSetPacketDest
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

	// our map scales for .75 to 10x, but the icons for the planets and fleets are 2x min
	const minZoom = 0.75;
	const maxZoom = 10;
	const minObjectZoom = 2;
	const scale = writable(3); // default 3x zoom
	const objectScale = derivedStore([scale], ([s]) => clamp(s, minObjectZoom, maxZoom));
	setContext('scale', scale);
	setContext('objectScale', objectScale);

	// zoom state that changes but doesn't cause a reaction
	let zooming = false;
	let fastestWaypoint = false;

	let pointerDown = false;
	let draggingWaypoint = false;
	let dragAndZoomEnabled = true;

	// set to true if we are moving a waypoint to a position rather than a target
	// this is enabled when the shift key is held
	let positionWaypoint = false;

	// if we just added a waypoint, don't drag it around
	let waypointJustAdded = false;

	// zoomBehavior is based on clientWidth/height
	let zoomBehavior: ZoomBehavior<HTMLElement, any> = $derived(
		zoom<HTMLElement, any>()
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
	function handleResize(event: UIEvent & { currentTarget: EventTarget & Window }) {
		clientRect = rect?.getBoundingClientRect() ?? { width: 100, height: 100 };
	}

	function handleKeyDown(e: KeyboardEvent) {
		// add a waypoint if we are currently commanding a fleet and we didn't just click
		// on the fleet
		shouldAddWaypoint = !!$commandedFleet && (e.shiftKey || e.metaKey);
		fastestWaypoint = !!$commandedFleet && e.metaKey;

		switch (e.key) {
			case '+':
			case '=':
				zoomViewport(clamp($scale + 1, minZoom, maxZoom));
				break;
			case '-':
			case '_':
				zoomViewport(clamp($scale - 1, minZoom, maxZoom));
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

	function handleZoom(e: D3ZoomEvent<HTMLElement, any>) {
		transform = e.transform;
		$scale = transform.k;
	}

	function handleZoomStart(e: D3ZoomEvent<HTMLElement, any>) {
		zooming = true;
	}

	function handleZoomEnd(e: D3ZoomEvent<HTMLElement, any>) {
		zooming = false;
	}

	// translate/zoom the display to a point on the map
	function translateViewport(position: Vector) {
		if (!root) {
			return;
		}

		select(root).call(zoomBehavior.scaleTo, $scale);
		const scaled: Vector = {
			x: scaler.x(position.x),
			y: scaler.y(position.y)
		};
		select(root)
			.call(zoomBehavior.translateTo, scaled.x, scaled.y)
			.call(zoomBehavior.scaleTo, $scale);
	}

	// zoom the viewport to a specific scale
	function zoomViewport(scaleTo: number) {
		if (!root || !zoomBehavior) {
			return;
		}
		select(root).call(zoomBehavior.scaleTo, scaleTo);
	}

	// turn off dragging
	function onContextMenu(e: FinderEvent) {
		const { event, found } = e;

		if (found && event instanceof MouseEvent) {
			onScannerContextPopup(event, found.position);
		}
	}

	// as the pointer moves, find the items it is under
	function onPointerMove(e: FinderEvent) {
		const { event, found, position } = e;

		highlightMapObject(found);

		if (draggingWaypoint && !zooming) {
			positionWaypoint = event.shiftKey;
			dragWaypointMove(position, found);
		}

		// check if we are over the commanded fleet's waypoint
		const fleetWaypoint =
			found &&
			$commandedFleet &&
			$commandedFleet.waypoints.slice(1).find((wp) => equal(wp.position, found.position));
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

		// check if we started a waypoint drag
		// we only
		// * start dragging once
		// * if the pointer is down
		// * if we are over a mapobject
		// * if we have a commanded fleet
		if (!waypointJustAdded && !draggingWaypoint && pointerDown && fleetWaypoint) {
			draggingWaypoint = true;
			onSelectWaypoint(fleetWaypoint);
		}
	}

	async function onPointerDown(e: FinderEvent) {
		const { event, found, position } = e;

		if (event instanceof MouseEvent && event.button != 0) {
			// we only care about the first button
			return;
		}

		if (found?.type == MapObjectType.Fleet && !filterFleet($player, found as Fleet, $settings)) {
			// this object we clicked is filtered out, don't do anything
			return;
		}

		pointerDown = true;

		if (found) {
			if ((shouldAddWaypoint || $settings.addWaypoint) && (await addWaypoint(found, position))) {
			} else {
				mapObjectSelected(found);
			}
		} else {
			if (shouldAddWaypoint || $settings.addWaypoint) {
				addWaypoint(found, position);
			}
		}
	}

	// turn off dragging
	function onPointerUp(e: FinderEvent) {
		const { event, found, position } = e;

		if (event instanceof MouseEvent && event.button != 0) {
			// we only care about the first button
			return;
		}

		if (draggingWaypoint) {
			if (!dragAndZoomEnabled) {
				enableDragAndZoom();
			}

			dragWaypointDone(position, found);
		}
		draggingWaypoint = false;
		pointerDown = false;
		waypointJustAdded = false;
	}

	// move the selected waypoint around snapping to targets
	function dragWaypointMove(position: Vector, mo: MapObject | undefined) {
		if ($selectedWaypoint && $currentSelectedWaypointIndex && $commandedFleet) {
			// don't move the waypoint to any adjacent waypoints
			if (mo && !positionWaypoint) {
				const index = $commandedFleet.waypoints.findIndex((wp) => equal(wp.position, mo.position));
				if (
					index == $currentSelectedWaypointIndex - 1 ||
					index == $currentSelectedWaypointIndex + 1
				) {
					return;
				}
			}

			const dest = mo && !positionWaypoint ? { mo: mo } : { position: position ?? emptyVector };
			onUpdateWaypoint(dest, fastestWaypoint, false);
		}
	}

	async function dragWaypointDone(position: Vector, mo: MapObject | undefined) {
		// reset waypoint dragging
		if ($selectedWaypoint && $commandedFleet && draggingWaypoint) {
			const dest = mo && !positionWaypoint ? { mo: mo } : { position: position ?? emptyVector };
			onUpdateWaypoint(dest, fastestWaypoint, true);
		}
	}

	// disable add waypoint mode when the user clicks outside the scanner
	function disableAddWaypointMode(event: MouseEvent) {
		// ignore clicks on the add-waypoint toolbar button
		const elem = event.target as Element;
		if (elem?.id == 'add-waypoint' || elem?.parentElement?.id == 'add-waypoint') {
			return;
		}
		if ($settings.addWaypoint) {
			$settings.addWaypoint = false;
			$settings.fastestWaypoint = false;
		}
	}

	// if the shift key is held, add a waypoint instead of selecting a mapobject
	async function addWaypoint(mo: MapObject | undefined, position: Vector): Promise<boolean> {
		if (zooming) {
			return false;
		}
		if (!$commandedFleet?.waypoints) {
			return false;
		}

		// for add waypoints, we always snap to planet because the "drag" and "add waypoint button" keys (shift) are the same
		const dest = mo ? { mo: mo } : { position: position ?? emptyVector };
		waypointJustAdded = await onAddWaypoint(dest, fastestWaypoint);
		return true;
	}
	/**
	 * When a mapobject is selected we go through a few steps.
	 * - We select it if it's a new selection
	 * - We cycle through our commandable objects at the same location if we own an object there
	 * @param mo
	 */
	function mapObjectSelected(mo: MapObject) {
		if ($settings.setPacketDest) {
			onSetPacketDest(mo);
		} else {
			onSelectMapObject(mo);
		}
	}

	// handle zoom in/out
	// this behavior controls how the zoom behaves
	// below we handle zooming events by updating a transform
	onMount(() => {
		clientRect = rect?.getBoundingClientRect() ?? { width: 100, height: 100 };
		if (!$zoomTarget) {
			return;
		}

		// setup zoom and translate to the zoomTarget
		translateViewport($zoomTarget.position);
		enableDragAndZoom();

		// setup asubscriber to draw the target X and move the viewport to a new target
		// when the zoomTarget changes
		const unsubscribeZoomTarget = zoomTarget.subscribe((target) => {
			if (target) {
				translateViewport(target.position);
				showTargetLocation();
			}
		});

		// if settings change, enable/disable drag and zoom depending on if the addWaypoint
		// mode is enabled or disabled
		const unsubscribeSettings = settings.subscribe((settings) => {
			settings.addWaypoint ? disableDragAndZoom() : enableDragAndZoom();
		});

		// bind the v key to show the target with X
		hotkeys('v', 'root', showTargetLocation);

		return () => {
			unsubscribeZoomTarget();
			unsubscribeSettings();
		};
	});

	onDestroy(() => {
		hotkeys.unbind('v', 'root', showTargetLocation);
	});

	// data used by the scanner is derived from the universe/commandedFleet stores
	// when they update, our data updates
	const data = derivedStore([universe, commandedFleet], ([u, f]) => [
		// add mapobject waypoints
		...(f?.getWaypointMapObjects(u) || []),
		...u.fleets.filter((f) => f.orbitingPlanetNum === None || f.orbitingPlanetNum === undefined),
		...u.mysteryTraders,
		...u.mineralPackets,
		...u.salvages,
		...u.wormholes,
		...u.mineFields,
		...u.planets
	]);
</script>

<svelte:window onresize={handleResize} onkeydown={handleKeyDown} onkeyup={handleKeyUp} />

<div
	class:cursor-grab={waypointHighlighted}
	class:cursor-cell={shouldAddWaypoint ||
		(!!$commandedFleet && $settings.addWaypoint) ||
		$settings.setPacketDest}
	class={`grow bg-black overflow-hidden p-[${padding}px] select-none`}
	bind:this={rect}
	use:clickOutside={disableAddWaypointMode}
>
	<LayerCake
		data={$data}
		x={(mo: MapObject | undefined) => mo?.position?.x}
		y={(mo: MapObject | undefined) => mo?.position?.y}
		xDomain={[0, $game.area.x]}
		yDomain={[0, $game.area.y]}
		xRange={xRange(clientRect.width, clientRect.height)}
		yRange={yRange(clientRect.width, clientRect.height)}
		yReverse={true}
		bind:element={root}
	>
		<Svg>
			<g transform={transform?.toString()}>
				<ScannerScanners />
				<ScannerMineFieldPattern />
				<ScannerMineFields />
				<ScannerPacketDests />
				<ScannerWaypoints />
				<ScannerPlanets />
				<ScannerMineralPackets />
				<ScannerWormholes />
				<ScannerFleets />
				<ScannerMysteryTraders />
				<ScannerWarpLine />
				<ScannerWormholeLinks />
				<ScannerSalvages />
				<SelectedMapObject />
				{#if showLocator}
					<ScannerMapObjectLocation show={$mostRecentMapObject} />
				{/if}
			</g>
		</Svg>
		<Html>
			{#if transform}
				<ScannerNames {transform} />

				<MapObjectQuadTreeFinder
					contextmenu={onContextMenu}
					pointermove={onPointerMove}
					pointerdown={onPointerDown}
					pointerup={onPointerUp}
					touchmove={onPointerMove}
					searchRadius={20}
					{transform}
				/>
			{/if}
		</Html>
	</LayerCake>
</div>
