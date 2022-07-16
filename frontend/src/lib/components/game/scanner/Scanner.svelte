<script lang="ts">
	import {
		commandedMapObject,
		commandedPlanet,
		game,
		player,
		selectedMapObject,
		selectMapObject,
		myMapObjectsByPosition,
		commandMapObject
	} from '$lib/services/Context';
	import { MapObjectType, ownedBy, positionKey, type MapObject } from '$lib/types/MapObject';
	import { select } from 'd3-selection';
	import { zoom, zoomIdentity, ZoomTransform, type D3ZoomEvent } from 'd3-zoom';
	import { Html, LayerCake, Svg } from 'layercake';
	import MapObjectQuadTreeFinder from './MapObjectQuadTreeFinder.svelte';
	import ScannerPlanets from './ScannerPlanets.svelte';
	import SelectedMapObject from './SelectedMapObject.svelte';
	import { clamp } from '$lib/services/Math';
	import { findIntelMapObject, findMyPlanet } from '$lib/types/Player';
	import type { Planet } from '$lib/types/Planet';
	import { findIndex } from 'lodash-es';

	const xGetter = (mo: MapObject) => mo?.position?.x;
	const yGetter = (mo: MapObject) => mo?.position?.y;

	let commandedMapObjectIndex = 0;

	// handle zoom in/out
	// this behavior controls how the zoom behaves
	// below we handle zooming events by updating a transform
	$: zoomBehavior = zoom<HTMLElement, any>()
		.extent([
			[0, 0],
			[root?.clientWidth ?? 100, root?.clientHeight ?? 100]
		])
		.scaleExtent([1, 5])
		.translateExtent([
			[0, 0],
			[root?.clientWidth ?? 100, root?.clientHeight ?? 100]
		])
		.on('zoom', handleZoom);

	function handleZoom(e: D3ZoomEvent<HTMLElement, any>) {
		transform = e.transform;
		// console.log(transform);
	}

	// attach the zoom behavior to the root element
	$: if (root && $game.area && $commandedMapObject.position) {
		select(root).call(zoomBehavior).on('dblclick.zoom', null);
		if ($commandedPlanet) {
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
		}
	}

	/**
	 * When a mapobject is selected we go through a few steps.
	 * - We select it if it's a new selection
	 * - We cycle through our commandable objects at the same location if we own an object there
	 * @param mo
	 */
	function mapobjectSelected(mo: MapObject) {
		const myMapObjectsAtPosition = $myMapObjectsByPosition[positionKey(mo)];
		let myMapObject = mo;
		if (ownedBy(mo, $player.num) && mo.type === MapObjectType.Planet) {
			myMapObject = findMyPlanet($player, mo as Planet) as MapObject;
		}

		const commandedIntelObject = findIntelMapObject($player, $commandedMapObject);

		if ($selectedMapObject !== myMapObject) {
			// we selected a different object, so just select it
			selectMapObject(myMapObject);
		} else if (
			ownedBy(mo, $player.num) &&
			myMapObjectsAtPosition?.length > 1 && // there are fleets orbiting this planet
			myMapObjectsAtPosition.find((m) => m === mo) && // the object we selected is in the things orbiting the planet
			myMapObjectsAtPosition.find((m) => m === commandedIntelObject) // the commanded object is also the planet or orbiting the planet
		) {
			// we selected one of our mapobjects and it's not currently commanded
			commandedMapObjectIndex =
				commandedMapObjectIndex >= myMapObjectsAtPosition.length - 1
					? 0
					: commandedMapObjectIndex + 1;
			// command the next object at this position
			const nextMapObject = myMapObjectsAtPosition[commandedMapObjectIndex];
			if (nextMapObject.type === MapObjectType.Planet) {
				commandMapObject(findMyPlanet($player, nextMapObject as Planet) as MapObject);
				commandedMapObjectIndex = 0;
			} else {
				commandMapObject(nextMapObject);
			}
		} else if ($selectedMapObject == myMapObject && commandedIntelObject != $selectedMapObject) {
			if (mo.type === MapObjectType.Planet) {
				commandMapObject(findMyPlanet($player, mo as Planet) as MapObject);
				commandedMapObjectIndex = 0;
			} else {
				commandMapObject(mo);
			}
		}
	}

	let transform: ZoomTransform;
	let root: HTMLElement;
</script>

<!-- <svelte:window on:resize={handleResize} /> -->

<div class="flex-1 h-full bg-black overflow-hidden p-5">
	{#if $game}
		<LayerCake
			data={$player.planetIntels}
			x={xGetter}
			y={yGetter}
			xDomain={[0, $game.area.x]}
			yDomain={[0, $game.area.y]}
			xRange={[0, root?.clientWidth ?? 100]}
			yRange={[root?.clientHeight ?? 100, 0]}
			bind:element={root}
		>
			<Svg>
				<g transform={transform?.toString()}>
					<ScannerPlanets />
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
						<div
							class="border-b-2 border-info absolute rounded-b-box"
							style="top:{transform.applyY(y - 10)}px;left:{transform.applyX(
								x - 10
							)}px;width:{transform.scale(20).k}px;height:{transform.scale(20).k}px;display: {found
								? 'block'
								: 'none'};"
						/>
					</MapObjectQuadTreeFinder>
				{/if}
			</Html>
		</LayerCake>
	{/if}
</div>
