<script lang="ts" module>
	import type { MapObjectLike } from '$lib/types/MapObject';

	/**
	 * FinderEvents are pointer/touch/mouse events that target a MapObject in the scanner
	 */
	export type FinderEvent = {
		event: PointerEvent | MouseEvent | TouchEvent;
		position: { x: number; y: number };
		found: MapObjectLike | undefined;
	};
</script>

<!--
  @component
  Creates an interaction layer (in HTML) using [d3-quadtree](https://github.com/d3/d3-quadtree) to find the nearest datapoint to the mouse.
  This component fires events for mouse movement/down/etc
 -->
<script lang="ts">
	import { quadtree } from 'd3-quadtree';
	import type { ZoomTransform } from 'd3-zoom';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import { getScannerContext } from './Scanner';

	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');
	const { scale } = getScannerContext();

	type Props = {
		// transform to transform our mouse to world coords
		transform: ZoomTransform;
		/** The number of pixels to search around the mouse's location. This is the third argument passed to [`quadtree.find`](https://github.com/d3/d3-quadtree#quadtree_find) and by default a value of `undefined` means an unlimited range. */
		searchRadius: number;

		pointermove: (e: FinderEvent) => void;
		pointerdown: (e: FinderEvent) => void;
		pointerup: (e: FinderEvent) => void;
		touchmove: (e: FinderEvent) => void;
		touchstart?: (e: FinderEvent) => void;
		touchend?: (e: FinderEvent) => void;
		contextmenu: (e: FinderEvent) => void;
	};

	let {
		transform,
		searchRadius,
		pointermove,
		pointerdown,
		pointerup,
		touchmove,
		touchstart,
		touchend,
		contextmenu
	}: Props = $props();

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

	function onPointerDown(event: PointerEvent) {
		const evt = event as PointerEvent & { layerX: number; layerY: number };
		const { position, found } = findItem(evt.layerX, evt.layerY);

		pointerdown({ event, position, found });
	}

	// as the pointer moves, find the items it is under
	function onPointerMove(event: PointerEvent) {
		// this is not supported, but works for me...
		const evt = event as PointerEvent & { layerX: number; layerY: number };
		const { position, found } = findItem(evt.layerX, evt.layerY);

		pointermove({ event, position, found });
	}

	// turn off dragging
	function onPointerUp(event: PointerEvent) {
		const evt = event as PointerEvent & { layerX: number; layerY: number };
		const { position, found } = findItem(evt.layerX, evt.layerY);

		pointerup({ event, position, found });
	}

	function onContextMenu(event: MouseEvent) {
		event.preventDefault();
		const evt = event as PointerEvent & { layerX: number; layerY: number };
		const { position, found } = findItem(evt.layerX, evt.layerY);

		contextmenu({ event, position, found });
	}

	function onTouchStart(event: TouchEvent) {
		if (event.target instanceof Element) {
			const bcr = event.target.getBoundingClientRect();
			const x = event.targetTouches[0].clientX - bcr.x;
			const y = event.targetTouches[0].clientY - bcr.y;
			const { position, found } = findItem(x, y);

			touchstart?.({ event, position, found });
		}
	}

	function onTouchMove(event: TouchEvent) {
		event.preventDefault();
		if (event.target instanceof Element) {
			const bcr = event.target.getBoundingClientRect();
			const x = event.targetTouches[0].clientX - bcr.x;
			const y = event.targetTouches[0].clientY - bcr.y;
			const { position, found } = findItem(x, y);

			touchmove({ event, position, found });
		}
	}

	function onTouchEnd(event: TouchEvent) {
		if (event.target instanceof Element) {
			const bcr = event.target.getBoundingClientRect();
			const x = event.changedTouches[0].clientX - bcr.x;
			const y = event.changedTouches[0].clientY - bcr.y;
			const { position, found } = findItem(x, y);

			touchend?.({ event, position, found });
		}
	}

	let finder = $derived(
		quadtree<MapObjectLike>()
			.extent([
				[-1, -1],
				[$width + 1, $height + 1]
			])
			.x((d) => $xGet(d.mapObject))
			.y((d) => $yGet(d.mapObject))
			.addAll($data)
	);
</script>

<div
	class="absolute h-full w-full z-10"
	role="link"
	tabindex="-1"
	ontouchstart={onTouchStart}
	ontouchmove={onTouchMove}
	ontouchend={onTouchEnd}
	oncontextmenu={onContextMenu}
	onpointerdown={onPointerDown}
	onpointermove={onPointerMove}
	onpointerup={onPointerUp}
></div>
