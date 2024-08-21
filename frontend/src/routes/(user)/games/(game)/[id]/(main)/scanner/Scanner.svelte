<script lang="ts">
	import { clickOutside } from '$lib/clickOutside';
	import { onScannerContextPopup } from '$lib/components/game/tooltips/ScannerContextPopup.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { clamp } from '$lib/services/Math';
	import { totalCargo } from '$lib/types/Cargo';
	import { WaypointTask, type Fleet, type Waypoint } from '$lib/types/Fleet';
	import {
		MapObjectType,
		None,
		StargateWarpSpeed,
		equal as mapObjectEqual,
		owned,
		type MapObject
	} from '$lib/types/MapObject';
	import type { Planet } from '$lib/types/Planet';
	import { distance, emptyVector, equal, type Vector } from '$lib/types/Vector';
	import type { ScaleLinear } from 'd3-scale';
	import { scaleLinear } from 'd3-scale';
	import { select } from 'd3-selection';
	import { ZoomTransform, zoom, type D3ZoomEvent, type ZoomBehavior } from 'd3-zoom';
	import hotkeys from 'hotkeys-js';
	import { Html, LayerCake, Svg } from 'layercake';
	import { onDestroy, onMount, setContext } from 'svelte';
	import { derived, writable } from 'svelte/store';
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
	import { filterFleet } from '$lib/types/Filter';

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

	const xGetter = (mo: MapObject) => mo?.position?.x;
	const yGetter = (mo: MapObject) => mo?.position?.y;

	let clientWidth = 100;
	let clientHeight = 100;
	let aspectRatio = 1;
	let transform: ZoomTransform;
	let zoomBehavior: ZoomBehavior<HTMLElement, any>;
	let root: HTMLElement;
	let padding = 20; // 20 px, used in zooming
	let scaleX: ScaleLinear<number, number, never>;
	let scaleY: ScaleLinear<number, number, never>;
	let zoomEnabled = true;
	let zooming = false;
	let showLocator = false;
	let shouldAddWaypoint = false;

	// our map scales for .75 to 10x, but the icons for the planets and fleets are 2x min
	const minZoom = 0.75;
	const maxZoom = 10;
	const minObjectZoom = 2;
	const scale = writable(3); // default 3x zoom
	const objectScale = derived([scale], ([s]) => clamp(s, minObjectZoom, maxZoom));
	setContext('scale', scale);
	setContext('objectScale', objectScale);

	// $: console.log('scale ', $scale);

	const unsubscribe = zoomTarget.subscribe(() => showTargetLocation());

	onMount(() => {
		hotkeys('v', 'root', showTargetLocation);
	});

	onDestroy(() => {
		hotkeys.unbind('v', 'root', showTargetLocation);
		unsubscribe();
	});

	// handle zoom in/out
	// this behavior controls how the zoom behaves
	// below we handle zooming events by updating a transform
	$: {
		if (root) {
			handleResize();

			zoomBehavior = zoom<HTMLElement, any>()
				.extent([
					[0, 0],
					[clientWidth, clientHeight]
				])
				.scaleExtent([minZoom, maxZoom])
				.translateExtent([
					[-20, -20],
					[clientWidth + padding, clientHeight + padding]
				])
				.on('zoom', handleZoom)
				.on('start', handleZoomStart)
				.on('end', handleZoomEnd);

			enableDragAndZoom();
		}
	}

	$: {
		if ($settings.addWaypoint && zoomEnabled) {
			disableDragAndZoom();
		} else if (!$settings.addWaypoint && !zoomEnabled) {
			enableDragAndZoom();
		}
	}

	// enable drag and zoom, but disable dblclick zoom events
	function enableDragAndZoom() {
		select(root).call(zoomBehavior).on('dblclick.zoom', null);
		dragAndZoomEnabled = true;
	}

	// disable drag and zoom temporarily
	function disableDragAndZoom() {
		select(root).on('.zoom', null);
		dragAndZoomEnabled = false;
		zooming = false;
	}

	const xRange = () => {
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
	};
	const yRange = () => {
		if (aspectRatio > 1 && clientHeight > clientWidth) {
			// tall skinny viewport, wide map, so fully expand on the x
			// but shrink up height
			return [0, clientWidth / aspectRatio];
		} else if (aspectRatio > 1 && clientWidth > clientHeight) {
			// wide viewport, wide map, so fully expand on the y
			return [0, clientHeight];
		}
		return [0, Math.min(clientWidth, clientHeight)];
	};

	function handleResize() {
		clientWidth = root?.clientWidth ?? 100;
		clientHeight = root?.clientHeight ?? 100;
		aspectRatio = $game.area.x / $game.area.y;

		// compute scales
		scaleX = scaleLinear().range(xRange()).domain([0, $game.area.x]);
		scaleY = scaleLinear().range(yRange()).domain([0, $game.area.y]);
	}

	function handleKeyDown(e: KeyboardEvent) {
		// add a waypoint if we are currently commanding a fleet and we didn't just click
		// on the fleet
		shouldAddWaypoint = !!$commandedFleet && e.shiftKey;

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
		shouldAddWaypoint = !!$commandedFleet && e.shiftKey;
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

	// zoom to the commanded map object every time it changes
	$: if (root && $zoomTarget) {
		translateViewport($zoomTarget.position);
	}

	// zoom the display to a point on the map
	function translateViewport(position: Vector, scaleTo?: number) {
		if (root) {
			select(root).call(zoomBehavior.scaleTo, $scale);
			const scaled: Vector = {
				x: scaleX(position.x),
				y: scaleY(position.y)
			};
			let localScale = $scale;
			if (scaleTo) {
				localScale = scaleTo;
			}
			select(root)
				.call(zoomBehavior.translateTo, scaled.x, scaled.y)
				.call(zoomBehavior.scaleTo, localScale);
		}
	}

	// zoom the viewport to a specific scale
	function zoomViewport(scaleTo: number) {
		if (root) {
			select(root).call(zoomBehavior.scaleTo, scaleTo);
		}
	}

	let pointerDown = false;
	let draggingWaypoint = false;
	let waypointHighlighted = false;
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
			return
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

			let warpSpeed = $selectedWaypoint?.warpSpeed;

			// update the ideal speed
			let waypointIndex = $currentSelectedWaypointIndex;

			if (waypointIndex > 0) {
				const previousWaypoint = $commandedFleet.waypoints[waypointIndex - 1];
				const dist = distance(mo?.position ?? position, previousWaypoint.position);

				warpSpeed = $commandedFleet.getMinimalWarp(dist, previousWaypoint.warpSpeed);
			}

			if (positionWaypoint || !mo) {
				$selectedWaypoint.position = position;
				$selectedWaypoint.warpSpeed = warpSpeed;
				$selectedWaypoint.targetType = MapObjectType.None;
				$selectedWaypoint.targetNum = None;
				$selectedWaypoint.targetPlayerNum = None;
				$selectedWaypoint.targetName = '';
			} else if (mo) {
				$selectedWaypoint.position = mo.position;
				$selectedWaypoint.warpSpeed = warpSpeed;
				$selectedWaypoint.targetType = mo.type;
				$selectedWaypoint.targetNum = mo.num;
				$selectedWaypoint.targetPlayerNum = mo.playerNum;
				$selectedWaypoint.targetName = mo.name;
			}
		}
	}

	function dragWaypointDone(position: Vector, mo: MapObject | undefined) {
		// reset waypoint dragging
		if ($selectedWaypoint && $commandedFleet && draggingWaypoint) {
			let waypointIndex = $currentSelectedWaypointIndex;

			if (waypointIndex > 0) {
				const previousWaypoint = $commandedFleet.waypoints[waypointIndex - 1];
				let warpSpeed = $selectedWaypoint?.warpSpeed
					? $selectedWaypoint?.warpSpeed
					: $commandedFleet.spec?.engine?.idealSpeed ?? 5;
				const dist = distance(mo?.position ?? position, previousWaypoint.position);

				warpSpeed = $commandedFleet.getMinimalWarp(dist, warpSpeed);
				$selectedWaypoint.warpSpeed = warpSpeed;
			}

			updateFleetOrders($commandedFleet);
		}
	}

	// disable add waypoint mode when the user clicks outside the
	// scanner
	function disableAddWaypointMode(event: MouseEvent) {
		// ignore clicks on the add-waypoint toolbar button
		const elem = event.target as Element;
		if (elem?.id == 'add-waypoint' || elem?.parentElement?.id == 'add-waypoint') {
			return;
		}
		if ($settings.addWaypoint) {
			$settings.addWaypoint = false;
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

		let waypointIndex = $currentSelectedWaypointIndex;
		if (waypointIndex == -1) {
			waypointIndex = 0;
		}
		const currentWaypoint = $commandedFleet.waypoints[waypointIndex];
		let nextWaypoint = currentWaypoint;
		if (waypointIndex < $commandedFleet.waypoints.length - 1) {
			nextWaypoint = $commandedFleet.waypoints[waypointIndex + 1];
		}

		position = mo ? mo.position : position ?? emptyVector;
		if (equal(position, currentWaypoint.position) || equal(position, nextWaypoint.position)) {
			// don't duplicate waypoints
			return false;
		}

		const dist = distance(mo?.position ?? position, currentWaypoint.position);

		const colonizing =
			$commandedFleet.spec.colonizer &&
			$commandedFleet.cargo.colonists &&
			mo &&
			mo.type === MapObjectType.Planet &&
			!owned(mo) &&
			((mo as Planet).spec.terraformedHabitability ?? 0) > 0;

		// use a stargate automatically if it's safe and in range
		const orbiting = $universe.getPlanet($commandedFleet.orbitingPlanetNum);
		const targetPlanet = mo?.type == MapObjectType.Planet ? (mo as Planet) : undefined;
		let stargate = false;
		if (orbiting && targetPlanet) {
			const destSafeHullMass = targetPlanet.spec.safeHullMass ?? 0;
			const destSafeRange = targetPlanet.spec.safeRange ?? 0;
			const sourceSafeHullMass = orbiting.spec.safeHullMass ?? 0;
			const sourceSafeRange = orbiting.spec.safeRange ?? 0;
			const destStargateSafe =
				(totalCargo($commandedFleet.cargo) == 0 || $player.race.spec?.canGateCargo) &&
				owned(targetPlanet) &&
				destSafeRange >= dist &&
				Math.max(
					...$commandedFleet.tokens.map((t) => $universe.getMyDesign(t.designNum)?.spec.mass ?? 0)
				) < destSafeHullMass;
			const sourceStargateSafe =
				(totalCargo($commandedFleet.cargo) == 0 || $player.race.spec?.canGateCargo) &&
				owned(orbiting) &&
				sourceSafeRange >= dist &&
				Math.max(
					...$commandedFleet.tokens.map((t) => $universe.getMyDesign(t.designNum)?.spec.mass ?? 0)
				) < sourceSafeHullMass;
			stargate = !!destStargateSafe && !!sourceStargateSafe;
		}

		// if the last waypoint warp is higher than the default engine speed, use it
		// otherwise use the default engine speed
		const engineIdealSpeed = $commandedFleet.spec?.engine?.idealSpeed ?? 5;
		let warpSpeed =
			$selectedWaypoint?.warpSpeed &&
			$selectedWaypoint.warpSpeed != StargateWarpSpeed &&
			$selectedWaypoint?.warpSpeed > engineIdealSpeed
				? $selectedWaypoint?.warpSpeed
				: engineIdealSpeed;

		// if colonizing, we want the max possible warp
		if (colonizing) {
			warpSpeed = $commandedFleet.getMaxWarp(
				dist,
				$universe,
				$player.race.spec?.fuelEfficiencyOffset ?? 0
			);
		} else {
			// use the minimal warp based on our ideal speed
			warpSpeed = $commandedFleet.getMinimalWarp(dist, warpSpeed);
		}

		// use a stargate if it's safe
		if (stargate) {
			warpSpeed = StargateWarpSpeed;
		}
		const task = $selectedWaypoint?.task ?? WaypointTask.None;
		const transportTasks = $selectedWaypoint?.transportTasks ?? {
			fuel: {},
			ironium: {},
			boranium: {},
			germanium: {},
			colonists: {}
		};

		if (!mo) {
			$commandedFleet.waypoints.splice(waypointIndex + 1, 0, {
				position: position,
				warpSpeed: warpSpeed,
				task: task,
				transportTasks: transportTasks
			});
		} else {
			const wp: Waypoint = {
				position: mo.position,
				targetName: mo.name,
				targetPlayerNum: mo.playerNum,
				targetNum: mo.num,
				targetType: mo.type,
				warpSpeed: warpSpeed,
				task: task,
				transportTasks: transportTasks
			};
			$commandedFleet.waypoints.splice(waypointIndex + 1, 0, wp);

			const remoteMining =
				$commandedFleet.spec.miningRate &&
				$commandedFleet.spec.miningRate > 0 &&
				mo.type === MapObjectType.Planet &&
				(!owned(mo) || $player.race.spec?.canRemoteMineOwnPlanets);

			// if this is a colonizer and the target is a habitable planet
			if (colonizing) {
				wp.task = WaypointTask.Colonize;
				wp.transportTasks = {
					fuel: {},
					ironium: {},
					boranium: {},
					germanium: {},
					colonists: {}
				};
			} else if (remoteMining) {
				wp.task = WaypointTask.RemoteMining;
				wp.transportTasks = {
					fuel: {},
					ironium: {},
					boranium: {},
					germanium: {},
					colonists: {}
				};
			}
		}
		waypointJustAdded = true;

		await updateFleetOrders($commandedFleet);

		// select the new waypoint
		selectWaypoint($commandedFleet.waypoints[waypointIndex + 1]);
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

	let data: MapObject[] = [];
	$: {
		const waypoints: MapObject[] = [];
		if ($commandedFleet?.waypoints) {
			waypoints.push(
				...$commandedFleet.waypoints.map((wp) => {
					const mo = $universe.getMapObject(wp);
					if (mo) {
						return mo;
					} else {
						return {
							position: wp.position,
							type: wp.targetType ?? MapObjectType.PositionWaypoint,
							name: wp.targetName ?? '',
							num: wp.targetNum ?? 0,
							playerNum: wp.targetPlayerNum ?? 0
						} as MapObject;
					}
				})
			);
		}
		data = [
			...waypoints,
			...$universe.fleets.filter(
				(f) => f.orbitingPlanetNum === None || f.orbitingPlanetNum === undefined
			),
			...$universe.mysteryTraders,
			...$universe.mineralPackets,
			...$universe.salvages,
			...$universe.wormholes,
			...$universe.mineFields,
			...$universe.planets
		];
	}
</script>

<svelte:window on:resize={handleResize} on:keydown={handleKeyDown} on:keyup={handleKeyUp} />

<div
	class:cursor-grab={waypointHighlighted}
	class:cursor-cell={shouldAddWaypoint ||
		(!!$commandedFleet && $settings.addWaypoint) ||
		$settings.setPacketDest}
	class={`grow bg-black overflow-hidden p-[${padding}px] select-none`}
	use:clickOutside={disableAddWaypointMode}
>
	<LayerCake
		{data}
		x={xGetter}
		y={yGetter}
		xDomain={[0, $game.area.x]}
		yDomain={[0, $game.area.y]}
		{xRange}
		{yRange}
		yReverse={true}
		bind:element={root}
	>
		<!-- <Svg viewBox={`0 0 ${game.area.x} ${game.area.y}`}> -->
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
