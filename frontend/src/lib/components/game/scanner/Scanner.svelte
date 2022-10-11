<script lang="ts">
	import {
		commandedMapObject,
		commandedPlanet,
		commandMapObject,
		game,
		myMapObjectsByPosition,
		player,
		selectedMapObject,
		selectMapObject
	} from '$lib/services/Context';
	import { NotOrbitingPlanet } from '$lib/types/Fleet';
	import { MapObjectType, ownedBy, positionKey, type MapObject } from '$lib/types/MapObject';
	import type { Planet } from '$lib/types/Planet';
	import { findIntelMapObject, findMyPlanet } from '$lib/types/Player';
	import { select } from 'd3-selection';
	import { zoom, zoomIdentity, ZoomTransform, type D3ZoomEvent, type ZoomBehavior } from 'd3-zoom';
	import { Html, LayerCake, Svg } from 'layercake';
	import MapObjectQuadTreeFinder from './MapObjectQuadTreeFinder.svelte';
	import ScannerFleets from './ScannerFleets.svelte';
	import ScannerPlanets from './ScannerPlanets.svelte';
	import ScannerScanners from './ScannerScanners.svelte';
	import ScannerWaypoints from './ScannerWaypoints.svelte';
	import SelectedMapObject from './SelectedMapObject.svelte';

	const xGetter = (mo: MapObject) => mo?.position?.x;
	const yGetter = (mo: MapObject) => mo?.position?.y;

	let clientWidth = 100;
	let clientHeight = 100;
	let aspectRatio = 1;
	let transform: ZoomTransform;
	let zoomBehavior: ZoomBehavior<HTMLElement, any>;
	let root: HTMLElement;
	let commandedMapObjectIndex = 0;

	// handle zoom in/out
	// this behavior controls how the zoom behaves
	// below we handle zooming events by updating a transform
	$: {
		if (root) {
			handleResize();
			zoomBehavior = zoom<HTMLElement, any>()
				.extent([
					[0, 0],
					[clientWidth * aspectRatio, clientHeight * aspectRatio]
				])
				.scaleExtent([1, 5])
				.translateExtent([
					[0, 0],
					[clientWidth * aspectRatio, clientHeight * aspectRatio]
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
		// console.log(transform);
	}

	// attach the zoom behavior to the root element
	$: if (root && $game?.area) {
		select(root).call(zoomBehavior).on('dblclick.zoom', null);
		// if ($commandedPlanet) {
		// set initial zoom
		select(root).call(zoomBehavior.transform, zoomIdentity);

		// select(root).call(
		// 	zoomBehavior.transform,
		// 	zoomIdentity
		// 		.scale(2)
		// 		.translate(
		// 			clamp(-($game.area.x - $commandedPlanet.position.x) / 2, -root.clientWidth, 0),
		// 			clamp(($game.area.y - $commandedPlanet.position.y) / 2, -$game.area.y, 0)
		// 		)
		// );
		// }
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

		const me = $player;

		let myMapObject = mo;
		if (ownedBy(mo, $player.num) && mo.type === MapObjectType.Planet) {
			myMapObject = findMyPlanet($player, mo as Planet) as MapObject;
		}

		const commandedIntelObject = findIntelMapObject($player, $commandedMapObject);

		if ($selectedMapObject !== myMapObject) {
			// we selected a different object, so just select it
			selectMapObject(myMapObject);
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
			data = [
				...$player.fleets.filter((f) => f.orbitingPlanetNum == NotOrbitingPlanet),
				...($player.fleetIntels?.filter((f) => f.orbitingPlanetNum == NotOrbitingPlanet) ?? []),
				...$player.planetIntels
			];
		}
	}
</script>

<svelte:window on:resize={handleResize} />

<div class="flex-1 h-full bg-black overflow-hidden p-5">
	{#if $game && $player}
		<LayerCake
			{data}
			x={xGetter}
			y={yGetter}
			xDomain={[0, $game.area.x]}
			yDomain={[0, $game.area.y]}
			xRange={[0, clientWidth * aspectRatio]}
			yRange={[clientHeight * aspectRatio, 0]}
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
