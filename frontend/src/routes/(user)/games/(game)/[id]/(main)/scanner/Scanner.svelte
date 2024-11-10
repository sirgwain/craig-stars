<script lang="ts">
	import { clickOutside } from '$lib/clickOutside';
	import { onScannerContextPopup } from '$lib/components/game/tooltips/ScannerContextPopup.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { clamp } from '$lib/services/Math';
	import { None } from '$lib/types/Constants';
	import { filterFleet } from '$lib/types/Filter';
	import { type Fleet } from '$lib/types/Fleet';
	import { MapObjectType, equal as mapObjectEqual, type MapObject } from '$lib/types/MapObject';
	import { emptyVector, equal, type Vector } from '$lib/types/Vector';
	import { scaleLinear } from 'd3-scale';
	import { select } from 'd3-selection';
	import { ZoomTransform, zoom, type D3ZoomEvent, type ZoomBehavior } from 'd3-zoom';
	import hotkeys from 'hotkeys-js';
	import { Html, LayerCake, Svg } from 'layercake';
	import { onDestroy, onMount, setContext } from 'svelte';
	import { derived as derivedStore, writable } from 'svelte/store';
	import MapObjectQuadTreeFinder, {
		type FinderEventDetails
	} from './MapObjectQuadTreeFinder.svelte';
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
		commandMapObject,
		commandedFleet,
		commandedMapObject,
		commandedPlanet,
		currentSelectedWaypointIndex,
		highlightMapObject,
		mostRecentMapObject,
		selectMapObject,
		selectWaypoint,
		selectedMapObject,
		selectedWaypoint,
		zoomTarget,
		updatePlanetOrders,
		updateFleetOrders
	} = getGameContext();

	let { onDeleteWaypoint }: { onDeleteWaypoint: () => void } = $props();

	const xGetter = (mo: MapObject) => mo?.position?.x;
	const yGetter = (mo: MapObject) => mo?.position?.y;
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
	let zoomEnabled = true;
	let zooming = false;
	let showLocator = $state(false);
	let shouldAddWaypoint = $state(false);
	let fastestWaypoint = false;

	// our map scales for .75 to 10x, but the icons for the planets and fleets are 2x min
	const minZoom = 0.75;
	const maxZoom = 10;
	const minObjectZoom = 2;
	const scale = writable(3); // default 3x zoom
	const objectScale = derivedStore([scale], ([s]) => clamp(s, minObjectZoom, maxZoom));
	setContext('scale', scale);
	setContext('objectScale', objectScale);

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
		if (!root) {
			return;
		}
		select(root).call(zoomBehavior).on('dblclick.zoom', null);
		dragAndZoomEnabled = true;
	}

	// disable drag and zoom temporarily
	function disableDragAndZoom() {
		if (!root) {
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
		// console.log('handleZoom', e, transform);
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

	let pointerDown = false;
	let draggingWaypoint = false;
	let waypointHighlighted = $state(false);
	let dragAndZoomEnabled = true;

	// set to true if we are moving a waypoint to a position rather than a target
	// this is enabled when the shift key is held
	let positionWaypoint = false;

	// if we just added a waypoint, don't drag it around
	let waypointJustAdded = false;

	// turn off dragging
	function onContextMenu(e: CustomEvent<FinderEventDetails>) {
		const { event, found } = e.detail;

		if (found && event instanceof MouseEvent) {
			onScannerContextPopup(event, found.position);
		}
	}

	// as the pointer moves, find the items it is under
	function onPointerMove(e: CustomEvent<FinderEventDetails>) {
		const { event, found, position } = e.detail;

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
			selectWaypoint(fleetWaypoint);
		}
	}

	async function onPointerDown(e: CustomEvent<FinderEventDetails>) {
		const { event, found, position } = e.detail;

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
	function onPointerUp(e: CustomEvent<FinderEventDetails>) {
		const { event, found, position } = e.detail;

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

			const dest = mo ? { mo: mo } : { position: position ?? emptyVector };

			// get highest mass of the fleet ships (for stargates)
			const highestShipMass = Math.max(
				...$commandedFleet.tokens.map((t) => $universe.getMyDesign(t.designNum)?.spec.mass ?? 0)
			);

			if (
				$commandedFleet.updateWaypoint(
					$player,
					$universe,
					dest,
					$currentSelectedWaypointIndex,
					highestShipMass,
					$settings.fastestWaypoint || fastestWaypoint
				)
			) {
				// trigger reaction
				$selectedWaypoint = $selectedWaypoint;
			}
		}
	}

	async function dragWaypointDone(position: Vector, mo: MapObject | undefined) {
		// reset waypoint dragging
		if ($selectedWaypoint && $commandedFleet && draggingWaypoint) {
			const dest = mo ? { mo: mo } : { position: position ?? emptyVector };

			// get highest mass of the fleet ships (for stargates)
			const highestShipMass = Math.max(
				...$commandedFleet.tokens.map((t) => $universe.getMyDesign(t.designNum)?.spec.mass ?? 0)
			);

			if (
				$commandedFleet.updateWaypoint(
					$player,
					$universe,
					dest,
					$currentSelectedWaypointIndex,
					highestShipMass,
					$settings.fastestWaypoint || fastestWaypoint
				)
			) {
				await updateFleetOrders($commandedFleet);

				// select the new waypoint
				selectWaypoint($commandedFleet.waypoints[$currentSelectedWaypointIndex]);
				if ($selectedWaypoint && $selectedWaypoint.targetType && $selectedWaypoint.targetNum) {
					const mo = $universe.getMapObject($selectedWaypoint);

					if (mo) {
						selectMapObject(mo);
					}
				}
			} else {
				// we dragged a waypoint to the previous position, delete it
				onDeleteWaypoint();
			}
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

		const dest = mo ? { mo: mo } : { position: position ?? emptyVector };

		// get highest mass of the fleet ships (for stargates)
		const highestShipMass = Math.max(
			...$commandedFleet.tokens.map((t) => $universe.getMyDesign(t.designNum)?.spec.mass ?? 0)
		);

		const newlyAddedWaypointIndex = $commandedFleet.addWaypoint(
			$player,
			$universe,
			dest,
			$currentSelectedWaypointIndex,
			highestShipMass,
			$settings.fastestWaypoint || fastestWaypoint
		);

		if (!newlyAddedWaypointIndex) {
			return false;
		}

		waypointJustAdded = true;

		await updateFleetOrders($commandedFleet);

		// select the new waypoint
		selectWaypoint($commandedFleet.waypoints[newlyAddedWaypointIndex]);
		if ($selectedWaypoint && $selectedWaypoint.targetType && $selectedWaypoint.targetNum) {
			const mo = $universe.getMapObject($selectedWaypoint);

			if (mo) {
				selectMapObject(mo);
			}
		}

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
			if (mo.type != MapObjectType.Planet) {
				return;
			} else {
				$settings.setPacketDest = false;
				// something went wrong, can't set dest on a planet without a massdriver
				if (!$commandedPlanet?.spec.hasMassDriver) {
					return;
				}

				if (mapObjectEqual(mo, $commandedPlanet)) {
					// clear dest
					$commandedPlanet.packetTargetNum = None;
				} else {
					$commandedPlanet.packetTargetNum = mo.num;
				}

				updatePlanetOrders($commandedPlanet);
				return;
			}
		}

		if ($selectedMapObject !== mo) {
			// we selected a different object, so just select it
			selectMapObject(mo);

			// if we selected a mapobject that is a waypoint, select the waypoint as well
			if ($commandedFleet?.waypoints) {
				const fleetWaypoint = $commandedFleet.waypoints.find((wp) =>
					equal(wp.position, mo.position)
				);
				if (fleetWaypoint) {
					selectWaypoint(fleetWaypoint);
				}
			}
		} else {
			// we selected the same mapobject twice
			const myMapObjectsAtPosition = $universe.getMyMapObjectsByPosition(mo);
			if (myMapObjectsAtPosition?.length > 0) {
				let index = myMapObjectsAtPosition.findIndex((mo) =>
					mapObjectEqual(mo, $commandedMapObject)
				);
				// if our currently commanded map object is not at this location, reset the index
				if (index == -1) {
					index = 0;
				} else {
					// command the next one
					index = index >= myMapObjectsAtPosition.length - 1 ? 0 : index + 1;
				}
				const nextMapObject = myMapObjectsAtPosition[index];

				commandMapObject(nextMapObject);
			}
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
		const unsubscribe = zoomTarget.subscribe((target) => {
			if (target) {
				translateViewport(target.position);
				showTargetLocation();
			}
		});

		// bind the v key to show the target with X
		hotkeys('v', 'root', showTargetLocation);

		return unsubscribe;
	});

	onDestroy(() => {
		hotkeys.unbind('v', 'root', showTargetLocation);
	});

	// enable/disable zoom on update
	$effect(() => {
		if ($settings.addWaypoint && zoomEnabled) {
			disableDragAndZoom();
		} else if (!$settings.addWaypoint && !zoomEnabled) {
			enableDragAndZoom();
		}
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
		x={xGetter}
		y={yGetter}
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
					on:contextmenu={onContextMenu}
					on:pointermove={onPointerMove}
					on:pointerdown={onPointerDown}
					on:pointerup={onPointerUp}
					on:touchmove={onPointerMove}
					searchRadius={20}
					{transform}
				/>
			{/if}
		</Html>
	</LayerCake>
</div>
