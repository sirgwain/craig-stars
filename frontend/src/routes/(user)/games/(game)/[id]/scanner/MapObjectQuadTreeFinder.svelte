<!--
  @component
  Creates an interaction layer (in HTML) using [d3-quadtree](https://github.com/d3/d3-quadtree) to find the nearest datapoint to the mouse. This component creates a slot that exposes variables `x`, `y`, `found` (the found datapoint), `visible` (a Boolean whether any data was found) and `e` (the event object).

  The quadtree searches across both the x and y dimensions at the same time. But if you want to only search across one, set the `x` and `y` props to the same value. For example, the [shared tooltip component](https://layercake.graphics/components/SharedTooltip.html.svelte) sets `y='x'` since it's nicer behavior to only pick up on the nearest x-value.
 -->
<script lang="ts">
	import { highlightMapObject } from '$lib/services/Context';
	import type { MapObject } from '$lib/types/MapObject';
	import type { Planet } from '$lib/types/Planet';
	import type { Vector } from '$lib/types/Vector';
	import { quadtree } from 'd3-quadtree';
	import type { LayerCake } from 'layercake';
	import type { ZoomTransform } from 'node_modules.nosync/@types/d3-zoom';
	import { createEventDispatcher, getContext } from 'svelte';

	const { data, xGet, yGet, xScale, yScale, xReverse, yReverse, width, height } =
		getContext<LayerCake>('LayerCake');
	const dispatch = createEventDispatcher();

	let found: MapObject | undefined;
	let position: Vector = { x: 0, y: 0 };
	let e = {};

	/** The dimension to search across when moving the mouse left and right. */
	export let x: string = 'x';

	/** The dimension to search across when moving the mouse up and down. */
	export let y: string = 'y';

	export let transform: ZoomTransform;

	/** The number of pixels to search around the mouse's location. This is the third argument passed to [`quadtree.find`](https://github.com/d3/d3-quadtree#quadtree_find) and by default a value of `undefined` means an unlimited range. */
	export let searchRadius: number | undefined = undefined;

	/** @type {Array} [dataset] – The dataset to work off of—defaults to $data if left unset. You can pass override the default here in here in case you don't want to use the main data or it's in a strange format. */
	export let dataset: [] | undefined = undefined;

	$: xGetter = x === 'x' ? $xGet : $yGet;
	$: yGetter = y === 'y' ? $yGet : $xGet;

	function findItem(evt: any) {
		e = evt;

		const xLayerKey = `layer${x.toUpperCase()}`;
		const yLayerKey = `layer${y.toUpperCase()}`;

		let [x1, y1] = [evt[xLayerKey] as number, evt[yLayerKey] as number];

		if (transform) {
			[x1, y1] = transform.invert([x1, y1]);
		}

		found = finder.find(
			x1,
			y1,
			transform && searchRadius ? transform.scale(searchRadius).k : searchRadius
		);
		position = { x: Math.round(x1 / $xScale(1)), y: Math.round(y1 / $yScale(1)) };

		// console.log(
		// 	'x, y',
		// 	[evt[xLayerKey] as number, evt[yLayerKey] as number],
		// 	'x1, y1',
		// 	[x1, y1],
		// 	'found',
		// 	found?.name,
		// 	'found?.position',
		// 	`${found?.position.x}, ${found?.position.y}`,
		// 	'position',
		// 	`${position.x}, ${position.y}`,
		// );

		highlightMapObject(found as Planet);
	}

	function selectItem(evt: MouseEvent) {
		if (found) {
			if (evt.shiftKey) {
				dispatch('add-waypoint', { mo: found });
			} else {
				dispatch('mapobject-selected', found);
			}
		} else {
			if (evt.shiftKey) {
				dispatch('add-waypoint', { position });
			}
		}
	}

	$: finder = quadtree<MapObject>()
		.extent([
			[-1, -1],
			[$width + 1, $height + 1]
		])
		.x(xGetter)
		.y(yGetter)
		.addAll(dataset || $data);
</script>

<div class="absolute h-full w-full z-10" on:mousemove={findItem} on:mousedown={selectItem} />
<slot x={xGetter(found) || 0} y={yGetter(found) || 0} {found} {e} />
