<script lang="ts">
	import {
		commandedFleet,
		commandedMapObject,
		commandMapObject,
		game,
		myMapObjectsByPosition,
		player,
		selectedMapObject,
		selectedWaypoint,
		selectMapObject,
		selectWaypoint,
		zoomTarget
	} from '$lib/services/Context';
	import { FleetService } from '$lib/services/FleetService';
	import { clamp } from '$lib/services/Math';
	import { MapObjectType, ownedBy, positionKey, type MapObject } from '$lib/types/MapObject';
	import type { Planet } from '$lib/types/Planet';
	import { findIntelMapObject, findMyPlanet } from '$lib/types/Player';
	import { emptyVector, equal, type Vector } from '$lib/types/Vector';
	import { select } from 'd3-selection';
	import { zoom, zoomIdentity, ZoomTransform, type D3ZoomEvent, type ZoomBehavior } from 'd3-zoom';
	import { Html, LayerCake, Svg } from 'layercake';
	import { merge } from 'lodash-es';
	import MapObjectQuadTreeFinder from './MapObjectQuadTreeFinder.svelte';
	import ScannerFleets from './ScannerFleets.svelte';
	import ScannerPlanets from './ScannerPlanets.svelte';
	import ScannerScanners from './ScannerScanners.svelte';
	import ScannerWaypoints from './ScannerWaypoints.svelte';
	import SelectedMapObject from './SelectedMapObject.svelte';

	const xGetter = (mo: MapObject) => mo?.position?.x;
	const yGetter = (mo: MapObject) => mo?.position?.y;

	const fleetService = new FleetService();

	let clientWidth = 100;
	let clientHeight = 100;
	let aspectRatio = 1;
	let transform: ZoomTransform;
	let zoomBehavior: ZoomBehavior<HTMLElement, any>;
	let root: HTMLElement;
	let commandedMapObjectIndex = 0;
	let scale = 2;
	let padding = 20; // 20 px, used in zooming

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
				.scaleExtent([1, 5])
				.translateExtent([
					[-20, -20],
					[clientWidth, clientHeight]
				])
				.on('zoom', handleZoom);
		}
	}

	function handleResize() {
		clientWidth = root?.clientWidth ?? 100;
		clientHeight = root?.clientHeight ?? 100;
		if (clientWidth > clientHeight) {
			aspectRatio = clientHeight / clientWidth;
		} else {
			aspectRatio = clientWidth / clientHeight;
		}
	}

	function handleZoom(e: D3ZoomEvent<HTMLElement, any>) {
		transform = e.transform;
		scale = transform.k;
		// console.log('initialZoom', initialZoom, transform);
	}

	// attach the zoom behavior to the root element
	$: if (root) {
		// disable double click on zoom, we use this to cycle commanded mapobjects
		select(root).call(zoomBehavior).on('dblclick.zoom', null);
	}

	// zoom to the commanded map object every time it changes
	$: if (root && $zoomTarget) {
		zoomToPosition($zoomTarget.position);
	}

	// zoom the display to a point on the map
	function zoomToPosition(position: Vector) {
		if (root) {
			select(root).call(
				zoomBehavior.transform,
				zoomIdentity
					.translate(clientWidth / 2 + padding, clientHeight / 2 + padding) // translate to center
					.scale(scale)
					.translate(
						// translate to selected mapobject
						clamp(-position.x / 2 - padding, -(clientWidth / 2 + padding) * scale, 0),
						clamp(-position.y / 2 + padding, -(clientHeight / 2 + padding) * scale, 0)
					)
			);
		}
	}

	// if the shift key is held, add a waypoint instead of selecting a mapobject
	async function addWaypoint(options: { mo?: MapObject; position?: Vector }) {
		if (!$myMapObjectsByPosition || !$player || !$commandedFleet) {
			return;
		}

		let waypointIndex = $commandedFleet.waypoints.findIndex((wp) => $selectedWaypoint == wp);
		if (waypointIndex == -1) {
			waypointIndex = 0;
		}
		const currentWaypoint = $commandedFleet.waypoints[waypointIndex];
		let nextWaypoint = currentWaypoint;
		if (waypointIndex < $commandedFleet.waypoints.length - 1) {
			nextWaypoint = $commandedFleet.waypoints[waypointIndex + 1];
		}

		const position = options.mo ? options.mo.position : options.position ?? emptyVector;
		if (equal(position, currentWaypoint.position) || equal(position, nextWaypoint.position)) {
			// don't duplicate waypoints
			return;
		}

		if (options.mo) {
			const mo = options.mo;
			$commandedFleet.waypoints.splice(waypointIndex + 1, 0, {
				position: mo.position,
				warpFactor: $commandedFleet.spec?.idealSpeed ?? 5,
				targetName: mo.name,
				targetPlayerNum: mo.playerNum,
				targetNum: mo.num,
				targetType: mo.type
			});
		} else if (options.position) {
			$commandedFleet.waypoints.splice(waypointIndex + 1, 0, {
				position: options.position,
				warpFactor: $commandedFleet.spec?.idealSpeed ?? 5
			});
		}

		// select the new waypoint
		selectWaypoint($commandedFleet.waypoints[$commandedFleet.waypoints.length - 1]);

		// save the commanded fleet
		const fleet = await fleetService.updateFleetOrders($commandedFleet);

		// update the player fleet
		merge($commandedFleet, fleet);

		commandedFleet.update(() => $commandedFleet);
	}
	/**
	 * When a mapobject is selected we go through a few steps.
	 * - We select it if it's a new selection
	 * - We cycle through our commandable objects at the same location if we own an object there
	 * @param mo
	 */
	function mapobjectSelected(mo: MapObject) {
		if (!$myMapObjectsByPosition || !$player) {
			return;
		}

		let myMapObject = mo;
		if (ownedBy(mo, $player.num) && mo.type === MapObjectType.Planet) {
			myMapObject = findMyPlanet($player, mo as Planet) as MapObject;
		}

		const commandedIntelObject = findIntelMapObject($player, $commandedMapObject);

		if ($selectedMapObject !== myMapObject) {
			// we selected a different object, so just select it
			selectMapObject(myMapObject);

			// if we selected a mapobject that is a waypoint, select the waypoint as well
			if ($commandedFleet) {
				const selectedWaypoint = $commandedFleet.waypoints.find((wp) =>
					equal(wp.position, myMapObject.position)
				);
				if (selectedWaypoint) {
					selectWaypoint(selectedWaypoint);
				}
			}
		} else {
			// we selected the same mapobject twice
			const myMapObjectsAtPosition = $myMapObjectsByPosition[positionKey(mo)];
			if ($player && myMapObjectsAtPosition?.length > 0) {
				// if our currently commanded map object is not at this location, reset the index
				if (!myMapObjectsAtPosition.find((mo) => mo == commandedIntelObject)) {
					commandedMapObjectIndex = 0;
				} else {
					// command the next one
					commandedMapObjectIndex =
						commandedMapObjectIndex >= myMapObjectsAtPosition.length - 1
							? 0
							: commandedMapObjectIndex + 1;
				}
				const nextMapObject = myMapObjectsAtPosition[commandedMapObjectIndex];
				if (nextMapObject.type === MapObjectType.Planet) {
					commandMapObject(findMyPlanet($player, nextMapObject as Planet) as MapObject);
					commandedMapObjectIndex = 0;
				} else {
					commandMapObject(nextMapObject);
				}
			}
		}
	}

	let data: MapObject[] = [];
	$: {
		if ($player) {
			const waypoints: MapObject[] = [];
			if ($commandedFleet) {
				waypoints.push(
					...$commandedFleet.waypoints.map(
						(wp) =>
							({
								position: wp.position,
								type: wp.targetType ?? MapObjectType.PositionWaypoint,
								name: wp.targetName ?? '',
								num: wp.targetNum ?? 0,
								playerNum: wp.targetPlayerNum ?? 0
							} as MapObject)
					)
				);
			}
			data = [
				...waypoints,
				...$player.fleets.filter((f) => !f.orbitingPlanetNum),
				...($player.fleetIntels?.filter((f) => !f.orbitingPlanetNum) ?? []),
				...$player.planetIntels
			];
		}
	}
</script>

<svelte:window on:resize={handleResize} />

<div class={`grow bg-black overflow-hidden p-[${padding}px]`}>
	{#if $game && $player}
		<LayerCake
			{data}
			x={xGetter}
			y={yGetter}
			xDomain={[0, $game.area.x]}
			yDomain={[0, $game.area.y]}
			xRange={[0, clientWidth * aspectRatio]}
			yRange={[0, clientHeight * aspectRatio]}
			yReverse={true}
			bind:element={root}
		>
			<!-- <Svg viewBox={`0 0 ${$game.area.x} ${$game.area.y}`}> -->
			<Svg>
				<g transform={transform?.toString()}>
					<ScannerScanners />
					<ScannerWaypoints />
					<ScannerPlanets />
					<ScannerFleets />
					<SelectedMapObject />
				</g>
			</Svg>
			<Html>
				{#if transform}
					<MapObjectQuadTreeFinder
						on:mapobject-selected={(mo) => mapobjectSelected(mo.detail)}
						on:add-waypoint={(mo) => addWaypoint(mo.detail)}
						searchRadius={20}
						let:x
						let:y
						let:found
						{transform}
					>
						<!-- <div
							class="border-b-2 border-info absolute rounded-b-box"
							style="top:{transform.applyY(y - 10)}px;left:{transform.applyX(
								x - 10
							)}px;width:{transform.scale(20).k}px;height:{transform.scale(20).k}px;display: {found
								? 'block'
								: 'none'};"
						/> -->
					</MapObjectQuadTreeFinder>
				{/if}
			</Html>
		</LayerCake>
	{/if}
</div>
