<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import {
		commandMapObject,
		commandedFleet,
		commandedMapObject,
		commandedPlanet,
		currentSelectedWaypointIndex,
		highlightMapObject,
		selectMapObject,
		selectWaypoint,
		selectedMapObject,
		selectedWaypoint,
		zoomTarget
	} from '$lib/services/Stores';
	import { WaypointTask, type Waypoint } from '$lib/types/Fleet';
	import {
		MapObjectType,
		None,
		equal as mapObjectEqual,
		type MapObject,
		owned
	} from '$lib/types/MapObject';
	import { emptyVector, equal, type Vector } from '$lib/types/Vector';
	import type { ScaleLinear } from 'd3-scale';
	import { scaleLinear } from 'd3-scale';
	import { select } from 'd3-selection';
	import { ZoomTransform, zoom, type D3ZoomEvent, type ZoomBehavior } from 'd3-zoom';
	import { Html, LayerCake, Svg } from 'layercake';
	import { setContext } from 'svelte';
	import { writable } from 'svelte/store';
	import MapObjectQuadTreeFinder, {
		type FinderEventDetails
	} from './MapObjectQuadTreeFinder.svelte';
	import ScannerFleets from './ScannerFleets.svelte';
	import ScannerMineFieldPattern from './ScannerMineFieldPattern.svelte';
	import ScannerMineFields from './ScannerMineFields.svelte';
	import ScannerMineralPackets from './ScannerMineralPackets.svelte';
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
	import SelectedWaypoint from './SelectedWaypoint.svelte';
	import type { Planet } from '$lib/types/Planet';

	const { game, player, universe, settings } = getGameContext();

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

	const scale = writable($game.area.y / 400); // tiny games are at 1x starting zoom, the rest zoom in based on universe size
	const clampedScale = writable($scale);
	$: $clampedScale = Math.min(3, $scale); // don't let the scale used for scanner objects go more than 1/2th size
	// $: console.log('scale ', $scale, ' clampedScale', $clampedScale);

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
				.scaleExtent([0.75, 10])
				.translateExtent([
					[-20, -20],
					[clientWidth + padding, clientHeight + padding]
				])
				.on('zoom', handleZoom);

			enableDragAndZoom();
		}
	}

	$: setPacketDest = $settings.setPacketDest;

	// enable drag and zoom, but disable dblclick zoom events
	const enableDragAndZoom = () => select(root).call(zoomBehavior).on('dblclick.zoom', null);
	// disable drag and zoom temporarily
	const disableDragAndZoom = () => select(root).on('.zoom', null);

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

	function handleZoom(e: D3ZoomEvent<HTMLElement, any>) {
		transform = e.transform;
		$scale = transform.k;
		// console.log('handleZoom', e, transform);
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

	let pointerDown = false;
	let dragging = false;
	let waypointHighlighted = false;
	let dragAndZoomEnabled = true;

	// set to true if we are moving a waypoint to a position rather than a target
	// this is enabled when the shift key is held
	let positionWaypoint = false;

	// if we just added a waypoint, don't drag it around
	let waypointJustAdded = false;

	// as the pointer moves, find the items it is under
	function onPointerMove(e: CustomEvent<FinderEventDetails>) {
		const { event, found, position } = e.detail;

		highlightMapObject(found);

		// check if we are over the commanded fleet's waypoint
		const fleetWaypoint =
			found &&
			$commandedFleet &&
			$commandedFleet.waypoints.slice(1).find((wp) => equal(wp.position, found.position));
		waypointHighlighted = !!fleetWaypoint;
		if (waypointHighlighted) {
			if (dragAndZoomEnabled) {
				dragAndZoomEnabled = false;
				disableDragAndZoom();
			}
		} else {
			if (!dragging && !dragAndZoomEnabled) {
				dragAndZoomEnabled = true;
				enableDragAndZoom();
			}
		}

		// check if we started a waypoint drag
		// we only
		// * start dragging once
		// * if the pointer is down
		// * if we are over a mapobject
		// * if we have a commanded fleet
		if (!waypointJustAdded && !dragging && pointerDown && fleetWaypoint) {
			dragging = true;
			selectWaypoint(fleetWaypoint);
		}

		if (dragging) {
			positionWaypoint = event.shiftKey;
			dragWaypointMove(position, found);
		}
	}

	function onPointerDown(e: CustomEvent<FinderEventDetails>) {
		const { event, found, position } = e.detail;
		pointerDown = true;

		if (found) {
			if (event.shiftKey || $settings.addWaypoint) {
				addWaypoint(found, position);
			} else {
				mapObjectSelected(found);
			}
		} else {
			if (event.shiftKey || $settings.addWaypoint) {
				addWaypoint(found, position);
			}
		}
	}

	// turn off dragging
	function onPointerUp(e: CustomEvent<FinderEventDetails>) {
		const { event, found, position } = e.detail;

		if (dragging) {
			if (!dragAndZoomEnabled) {
				dragAndZoomEnabled = true;
				enableDragAndZoom();
			}

			dragWaypointDone(position, found);
		}
		dragging = false;
		pointerDown = false;
		waypointJustAdded = false;
	}

	// move the selected waypoint around snapping to targets
	function dragWaypointMove(position: Vector, target: MapObject | undefined) {
		if ($selectedWaypoint && $currentSelectedWaypointIndex && $commandedFleet) {
			// don't move the waypoint to any adjacent waypoints
			if (target) {
				const index = $commandedFleet.waypoints.findIndex((wp) =>
					equal(wp.position, target.position)
				);
				if (
					index == $currentSelectedWaypointIndex - 1 ||
					index == $currentSelectedWaypointIndex + 1
				) {
					return;
				}
			}

			if (positionWaypoint || !target) {
				$selectedWaypoint.position = position;
				$selectedWaypoint.targetType = MapObjectType.None;
				$selectedWaypoint.targetNum = None;
				$selectedWaypoint.targetPlayerNum = None;
				$selectedWaypoint.targetName = '';
			} else if (target) {
				$selectedWaypoint.position = target.position;
				$selectedWaypoint.targetType = target.type;
				$selectedWaypoint.targetNum = target.num;
				$selectedWaypoint.targetPlayerNum = target.playerNum;
				$selectedWaypoint.targetName = target.name;
			}
		}
	}

	function dragWaypointDone(position: Vector, mo: MapObject | undefined) {
		// reset waypoint dragging
		if ($selectedWaypoint && $commandedFleet && dragging) {
			$game.updateFleetOrders($commandedFleet);
		}
	}
	// if the shift key is held, add a waypoint instead of selecting a mapobject
	async function addWaypoint(mo: MapObject | undefined, position: Vector) {
		if (!$commandedFleet?.waypoints) {
			return;
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
			return;
		}

		const warpSpeed = $selectedWaypoint?.warpSpeed ?? $commandedFleet.spec?.engine?.idealSpeed ?? 5;
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

			// if this is a colonizer and the target is a habitable planet
			if (
				$commandedFleet.spec.colonizer &&
				$commandedFleet.cargo.colonists &&
				mo.type === MapObjectType.Planet &&
				!owned(mo) &&
				((mo as Planet).spec.terraformedHabitability ?? 0) > 0
			) {
				wp.task = WaypointTask.Colonize;
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

		await $game.updateFleetOrders($commandedFleet);

		// select the new waypoint
		selectWaypoint($commandedFleet.waypoints[$commandedFleet.waypoints.length - 1]);
		if ($selectedWaypoint && $selectedWaypoint.targetType && $selectedWaypoint.targetNum) {
			const mo = $universe.getMapObject($selectedWaypoint);

			if (mo) {
				selectMapObject(mo);
			}
		}
	}
	/**
	 * When a mapobject is selected we go through a few steps.
	 * - We select it if it's a new selection
	 * - We cycle through our commandable objects at the same location if we own an object there
	 * @param mo
	 */
	function mapObjectSelected(mo: MapObject) {
		if (setPacketDest) {
			if (mo.type != MapObjectType.Planet) {
				return;
			} else {
				$settings.setPacketDest = false;
				// something went wrong, can't set dest on a planet without a massdriver
				if (!$commandedPlanet?.spec.hasMassDriver) {
					return;
				}
				$commandedPlanet.packetTargetNum = mo.num;
				$game.updatePlanetOrders($commandedPlanet);
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

	setContext('scale', clampedScale);
</script>

<svelte:window on:resize={handleResize} />

<div
	class:cursor-grab={waypointHighlighted}
	class={`grow bg-black overflow-hidden p-[${padding}px]`}
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
			<g preserveAspectRatio="true" transform={transform?.toString()}>
				<ScannerScanners />
				<ScannerMineFieldPattern />
				<ScannerMineFields />
				<ScannerPacketDests />
				<ScannerWaypoints />
				<SelectedWaypoint />
				<ScannerPlanets />
				<ScannerMineralPackets />
				<ScannerWormholes />
				<ScannerFleets />
				<ScannerWarpLine />
				<ScannerWormholeLinks />
				<ScannerSalvages />
				<SelectedMapObject />
			</g>
		</Svg>
		<Html>
			{#if transform}
				<ScannerNames {transform} />

				<MapObjectQuadTreeFinder
					on:pointermove={onPointerMove}
					on:pointerdown={onPointerDown}
					on:pointerup={onPointerUp}
					searchRadius={20}
					{transform}
				/>
			{/if}
		</Html>
	</LayerCake>
</div>
