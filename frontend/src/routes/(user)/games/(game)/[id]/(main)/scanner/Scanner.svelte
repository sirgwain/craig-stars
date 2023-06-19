<script lang="ts">
	import {
		commandedFleet,
		commandedMapObject,
		commandMapObject,
		selectedMapObject,
		selectedWaypoint,
		selectMapObject,
		selectWaypoint,
		zoomTarget
	} from '$lib/services/Context';
	import { FleetService } from '$lib/services/FleetService';
	import type { FullGame } from '$lib/services/FullGame';
	import { MapObjectType, ownedBy, type MapObject } from '$lib/types/MapObject';
	import { emptyVector, equal, type Vector } from '$lib/types/Vector';
	import type { ScaleLinear } from 'd3-scale';
	import { scaleLinear } from 'd3-scale';
	import { select } from 'd3-selection';
	import { zoom, ZoomTransform, type D3ZoomEvent, type ZoomBehavior } from 'd3-zoom';
	import { Html, LayerCake, Svg } from 'layercake';
	import { setContext } from 'svelte';
	import MapObjectQuadTreeFinder from './MapObjectQuadTreeFinder.svelte';
	import ScannerFleets from './ScannerFleets.svelte';
	import ScannerPlanets from './ScannerPlanets.svelte';
	import ScannerScanners from './ScannerScanners.svelte';
	import ScannerWaypoints from './ScannerWaypoints.svelte';
	import SelectedMapObject from './SelectedMapObject.svelte';

	export let game: FullGame;

	setContext('game', game);

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
	let scale = 3;
	let padding = 20; // 20 px, used in zooming
	let scaleX: ScaleLinear<number, number, never>;
	let scaleY: ScaleLinear<number, number, never>;

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
				.scaleExtent([0.75, 8])
				.translateExtent([
					[-20, -20],
					[clientWidth + padding, clientHeight + padding]
				])
				.on('zoom', handleZoom);

			// disable double click on zoom, we use this to cycle commanded mapobjects
			select(root).call(zoomBehavior).on('dblclick.zoom', null);
			// jump to 0,0 at our scale
			// select(root).call(zoomBehavior.scaleTo, scale).call(zoomBehavior.translateTo, 0, 0);
		}
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
		aspectRatio = game.area.x / game.area.y;

		// compute scales
		scaleX = scaleLinear().range(xRange()).domain([0, game.area.x]);
		scaleY = scaleLinear().range(yRange()).domain([0, game.area.y]);
	}

	function handleZoom(e: D3ZoomEvent<HTMLElement, any>) {
		transform = e.transform;
		scale = transform.k;
		// console.log('handleZoom', e, transform);
	}

	// zoom to the commanded map object every time it changes
	$: if (root && $zoomTarget) {
		translateViewport($zoomTarget.position);
	}

	// zoom the display to a point on the map
	function translateViewport(position: Vector, scaleTo?: number) {
		if (root) {
			select(root).call(zoomBehavior.scaleTo, scale);
			const scaled: Vector = {
				x: scaleX(position.x),
				y: scaleY(position.y)
			};
			let localScale = scale;
			if (scaleTo) {
				localScale = scaleTo;
			}
			select(root)
				.call(zoomBehavior.translateTo, scaled.x, scaled.y)
				.call(zoomBehavior.scaleTo, localScale);
		}
	}

	// if the shift key is held, add a waypoint instead of selecting a mapobject
	async function addWaypoint(options: { mo?: MapObject; position?: Vector }) {
		if (!$commandedFleet?.waypoints) {
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

		await game.updateFleetOrders($commandedFleet);
		// select the new waypoint
		selectWaypoint($commandedFleet.waypoints[$commandedFleet.waypoints.length - 1]);
		if ($selectedWaypoint && $selectedWaypoint.targetType && $selectedWaypoint.targetNum) {
			const mo = game.universe.getMapObject(
				$selectedWaypoint.targetType,
				$selectedWaypoint.targetNum,
				$selectedWaypoint.targetPlayerNum
			);

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
	function mapobjectSelected(mo: MapObject) {
		let myMapObject = mo;
		if (ownedBy(mo, game.player.num) && mo.type === MapObjectType.Planet) {
			myMapObject = game.getPlanet(mo.num) ?? mo;
		}

		const commandedIntelObject = $commandedMapObject
			? game.universe.getMapObject(
					$commandedMapObject.type,
					$commandedMapObject.num,
					$commandedMapObject.playerNum
			  )
			: undefined;

		if ($selectedMapObject !== myMapObject) {
			// we selected a different object, so just select it
			selectMapObject(myMapObject);

			// if we selected a mapobject that is a waypoint, select the waypoint as well
			if ($commandedFleet?.waypoints) {
				const selectedWaypoint = $commandedFleet.waypoints.find((wp) =>
					equal(wp.position, myMapObject.position)
				);
				if (selectedWaypoint) {
					selectWaypoint(selectedWaypoint);
				}
			}
		} else {
			// we selected the same mapobject twice
			const myMapObjectsAtPosition = game.universe.getMyMapObjectsByPosition(mo);
			if (myMapObjectsAtPosition?.length > 0) {
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
					commandMapObject(game.getPlanet(nextMapObject.num) as MapObject);
					commandedMapObjectIndex = 0;
				} else {
					commandMapObject(nextMapObject);
				}
			}
		}
	}

	let data: MapObject[] = [];
	$: {
		const waypoints: MapObject[] = [];
		if ($commandedFleet?.waypoints) {
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
			...game.universe.fleets.filter((f) => !f.orbitingPlanetNum),
			...(game.player.fleetIntels?.filter((f) => !f.orbitingPlanetNum) ?? []),
			...game.player.planetIntels
		];
	}
</script>

<svelte:window on:resize={handleResize} />

<div class={`grow bg-black overflow-hidden p-[${padding}px]`}>
	<LayerCake
		{data}
		x={xGetter}
		y={yGetter}
		xDomain={[0, game.area.x]}
		yDomain={[0, game.area.y]}
		{xRange}
		{yRange}
		yReverse={true}
		bind:element={root}
	>
		<!-- <Svg viewBox={`0 0 ${game.area.x} ${game.area.y}`}> -->
		<Svg>
			<g preserveAspectRatio="true" transform={transform?.toString()}>
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
				/>
			{/if}
		</Html>
	</LayerCake>
</div>
