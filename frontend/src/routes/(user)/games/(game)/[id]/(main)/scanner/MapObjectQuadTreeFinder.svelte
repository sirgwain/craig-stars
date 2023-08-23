<script lang="ts" context="module">
	export type FinderEventDetails = {
		event: PointerEvent | MouseEvent;
		position: Vector;
		found: MapObject | undefined;
	};
	export type FinderEvent = {
		pointermove: FinderEventDetails;
		pointerdown: FinderEventDetails;
		pointerup: FinderEventDetails;
		contextmenu: FinderEventDetails;
	};
</script>

<!--
  @component
  Creates an interaction layer (in HTML) using [d3-quadtree](https://github.com/d3/d3-quadtree) to find the nearest datapoint to the mouse.
  This component fires events for mouse movement/down/etc
 -->
<script lang="ts">
	import type { MapObject } from '$lib/types/MapObject';
	import type { Vector } from '$lib/types/Vector';
	import { quadtree } from 'd3-quadtree';
	import type { ZoomTransform } from 'd3-zoom';
	import type { LayerCake } from 'layercake';
	import { createEventDispatcher, getContext } from 'svelte';
	import type { Writable } from 'svelte/store';

	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');
	const scale = getContext<Writable<number>>('scale');
	const dispatch = createEventDispatcher<FinderEvent>();

	// transform to transform our mouse to world coords
	export let transform: ZoomTransform;

	/** The number of pixels to search around the mouse's location. This is the third argument passed to [`quadtree.find`](https://github.com/d3/d3-quadtree#quadtree_find) and by default a value of `undefined` means an unlimited range. */
	export let searchRadius: number;

	// find the item under
	function findItem(x: number, y: number) {
		let [x1, y1] = [x, y];

		if (transform) {
			[x1, y1] = transform.invert([x1, y1]);
		}

		const found = finder.find(x1, y1, searchRadius / $scale);
		const position = { x: Math.round(x1 / $xScale(1)), y: Math.round(y1 / $yScale(1)) };

		return { position, found };
	}

	// as the pointer moves, find the items it is under
	function onPointerMove(event: PointerEvent) {
		// this is not supported, but works for me...
		const evt = event as PointerEvent & { layerX: number; layerY: number };
		const { position, found } = findItem(evt.layerX, evt.layerY);

		dispatch('pointermove', { event, position, found });
	}

	function onPointerDown(event: PointerEvent) {
		const evt = event as PointerEvent & { layerX: number; layerY: number };
		const { position, found } = findItem(evt.layerX, evt.layerY);

		dispatch('pointerdown', { event, position, found });
	}

	// turn off dragging
	function onPointerUp(event: PointerEvent) {
		const evt = event as PointerEvent & { layerX: number; layerY: number };
		const { position, found } = findItem(evt.layerX, evt.layerY);

		dispatch('pointerup', { event, position, found });
	}

	function onContextMenu(event: MouseEvent) {
		const evt = event as PointerEvent & { layerX: number; layerY: number };
		const { position, found } = findItem(evt.layerX, evt.layerY);

		dispatch('contextmenu', { event, position, found });
	}

	$: finder = quadtree<MapObject>()
		.extent([
			[-1, -1],
			[$width + 1, $height + 1]
		])
		.x($xGet)
		.y($yGet)
		.addAll($data);
</script>

<div
	class="absolute h-full w-full z-10"
	on:contextmenu|preventDefault={onContextMenu}
	on:pointermove={onPointerMove}
	on:pointerdown={onPointerDown}
	on:pointerup={onPointerUp}
/>
