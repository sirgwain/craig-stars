<script lang="ts">
	import { clickOutside } from '$lib/clickOutside';
	import { onScannerContextPopup } from '$lib/components/game/tooltips/ScannerContextPopup.svelte';
	import type { SelectWaypointProps } from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import { clamp } from '$lib/services/Math';
	import { None } from '$lib/types/Consts';
	import {
		MapObjectType,
		VectorSchema,
		WaypointDestSchema,
		type Fleet,
		type Vector,
		type WaypointDest
	} from '$lib/types/cs-proto';
	import { filterFleet } from '$lib/types/Filter';
	import { emptyMapObject, type MapObjectLike, type Position } from '$lib/types/MapObject';
	import { emptyVector, equal } from '$lib/types/Vector';
	import { create } from '@bufbuild/protobuf';
	import { scaleLinear } from 'd3-scale';
	import { select } from 'd3-selection';
	import { zoom, ZoomTransform, type D3ZoomEvent, type ZoomBehavior } from 'd3-zoom';
	import hotkeys from 'hotkeys-js';
	import { Html, LayerCake, Svg } from 'layercake';
	import { onDestroy, onMount } from 'svelte';
	import { derived as derivedStore, writable } from 'svelte/store';
	import MapObjectQuadTreeFinder, { type FinderEvent } from './MapObjectQuadTreeFinder.svelte';
	import { setScannerContext } from './Scanner';
	import ScannerFleets from './ScannerFleets.svelte';
	import ScannerMapObjectLocation from './ScannerMapObjectLocation.svelte';
	import ScannerMinefieldPattern from './ScannerMinefieldPattern.svelte';
	import ScannerMinefields from './ScannerMinefields.svelte';
	import ScannerMineralPackets from './ScannerMineralPackets.svelte';
	import ScannerMysteryTraders from './ScannerMysteryTraders.svelte';
	import ScannerNames from './ScannerNames.svelte';
	import ScannerPacketDests from './ScannerPacketDests.svelte';
	import ScannerPlanets from './ScannerPlanets.svelte';
	import ScannerRouteDests from './ScannerRouteDests.svelte';
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

	// our map scales for .75 to 10x, but the icons for the planets and fleets are 2x min
	const minZoom = 0.75;
	const maxZoom = 10;
	const minObjectZoom = 2;
	const scale = writable(1.5); // default 3x zoom
	const objectScale = derivedStore([scale], ([s]) => clamp(s, minObjectZoom, maxZoom));
	setScannerContext({ scale, objectScale });

	// zoom state that changes but doesn't cause a reaction
	let zooming = false;
	let fastestWaypoint = false;

	let pointerDown = false;
	let draggingWaypoint = false;
	let dragAndZoomEnabled = false;

	// set to true if we are moving a waypoint to a position rather than a target
	// this is enabled when the shift key is held
	let positionWaypoint = false;

	// if we just added a waypoint, don't drag it around
	let waypointJustAdded = false;

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
	function handleResize() {
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

	function handleZoom(e: D3ZoomEvent<HTMLElement, unknown>) {
		transform = e.transform;
		$scale = transform.k;
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

		select(root).call(zoomBehavior.scaleTo, $scale);
		const scaled: Vector = create(VectorSchema, {
			x: scaler.x(Number(position.x ?? 0)),
			y: scaler.y(Number(position.y ?? 0))
		});
		select(root)
			.call(zoomBehavior.translateTo, scaled.x, scaled.y)
			.call(zoomBehavior.scaleTo, $scale);
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
			.call(zoomBehavior.scaleTo, $scale);
	}

	// turn off dragging
	function onContextMenu(e: FinderEvent) {
		const { event, found } = e;

		if (found && event instanceof MouseEvent) {
			onScannerContextPopup(event, found.mapObject?.position ?? emptyVector());
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

		// check if we started a waypoint drag
		// we only
		// * start dragging once
		// * if the pointer is down
		// * if we are over a mapobject
		// * if we have a commanded fleet
		if (!waypointJustAdded && !draggingWaypoint && pointerDown && fleetWaypoint) {
			draggingWaypoint = true;
			onSelectWaypoint?.({ fleet: $commandedFleet, waypoint: fleetWaypoint });
		}
	}

	async function onPointerDown(e: FinderEvent) {
		const { event, found, position } = e;

		if (event instanceof MouseEvent && event.button != 0) {
			// we only care about the first button
			return;
		}

		if (
			found?.mapObject?.type === MapObjectType.FLEET &&
			!filterFleet($player, found as Fleet, $settings)
		) {
			// this object we clicked is filtered out, don't do anything
			return;
		}

		pointerDown = true;

		if (found) {
			if ((shouldAddWaypoint || $settings.addWaypoint) && (await addWaypoint(found, position))) {
				// ignore
			} else {
				// check if we are clicked the commanded fleet's waypoint
				const fleetWaypoint =
					found &&
					$commandedFleet &&
					$commandedFleet.fleetOrders.waypoints
						.slice(1)
						.find((wp) =>
							equal(wp.position ?? emptyVector(), found.mapObject?.position ?? emptyVector())
						);

				if (fleetWaypoint && $selectedWaypoint != fleetWaypoint) {
					onSelectWaypoint?.({ fleet: $commandedFleet, waypoint: fleetWaypoint });
				} else {
					mapObjectSelected(found);
				}
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
		if (elem.id == 'add-waypoint' || elem.parentElement?.id == 'add-waypoint') {
			return;
		}
		if ($settings.addWaypoint) {
			$settings.addWaypoint = false;
			$settings.fastestWaypoint = false;
		}
	}

	// if the shift key is held, add a waypoint instead of selecting a mapobject
	async function addWaypoint(mo: MapObjectLike | undefined, position: Position): Promise<boolean> {
		if (zooming) {
			return false;
		}
		if (!$commandedFleet?.fleetOrders.waypoints) {
			return false;
		}

		// for add waypoints, we always snap to planet because the "drag" and "add waypoint button" keys (shift) are the same
		const dest = create(WaypointDestSchema, {
			mo: mo?.mapObject ?? emptyMapObject(),
			position: !mo?.mapObject ? { x: Number(position.x), y: Number(position.y) } : undefined
		});
		waypointJustAdded = await onAddWaypoint(dest, fastestWaypoint);
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

		// if settings change, enable/disable drag and zoom depending on if the addWaypoint
		// mode is enabled or disabled
		const unsubscribeSettings = settings.subscribe((settings) => {
			if (settings.addWaypoint) {
				disableDragAndZoom();
			} else {
				enableDragAndZoom();
			}
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
		...u.getAllFleets().filter((f) => f.orbitingPlanetNum === None),
		...u.mysteryTraders,
		...u.salvages,
		...u.wormholes,
		...u.mineralPackets,
		...u.allMinefields,
		...u.allPlanets
	]);

	// all our data in LayerCake are mapObjects/waypoints. Add this custom getter to get the
	// x/y coords of a mapobject or waypoint
	const xGet = (mo: { position: Position | undefined }) => mo.position?.x ?? 0;
	const yGet = (mo: { position: Position | undefined }) => mo.position?.y ?? 0;
</script>

<svelte:window onresize={handleResize} onkeydown={handleKeyDown} onkeyup={handleKeyUp} />

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
	<LayerCake
		data={$data}
		x={xGet}
		y={yGet}
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
				<ScannerMinefieldPattern />
				<ScannerMinefields />
				<ScannerPacketDests />
				<ScannerRouteDests />
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
