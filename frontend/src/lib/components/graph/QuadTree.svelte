<!--
  @component
  Creates an interaction layer (in HTML) using [d3-quadtree](https://github.com/d3/d3-quadtree) to find the nearest datapoint to the mouse. This component creates a slot that exposes variables `x`, `y`, `found` (the found datapoint), `visible` (a Boolean whether any data was found) and `e` (the event object).

  The quadtree searches across both the x and y dimensions at the same time. But if you want to only search across one, set the `x` and `y` props to the same value. For example, the [shared tooltip component](https://layercake.graphics/components/SharedTooltip.html.svelte) sets `y='x'` since it's nicer behavior to only pick up on the nearest x-value.
 -->
<script lang="ts">
	import { getContext } from 'svelte';
	import { quadtree } from 'd3-quadtree';
	import type { LayerCake } from 'layercake';
	const { data, xGet, yGet, zGet, xScale, yScale, width, height, config } =
		getContext<LayerCake>('LayerCake');

	let visible = $state(false);
	let found = $state({});
	let e = $state({});

	interface Props {
		x?: String;
		y?: String;
		/** @type {String} [searchRadius] – The number of pixels to search around the mouse's location. This is the third argument passed to [`quadtree.find`](https://github.com/d3/d3-quadtree#quadtree_find) and by default a value of `undefined` means an unlimited range. */
		searchRadius?: number | undefined;
		/** @type {Array} [dataset] – The dataset to work off of—defaults to $data if left unset. You can pass override the default here in here in case you don't want to use the main data or it's in a strange format. */
		dataset?: any;
		children?: import('svelte').Snippet<[any]>;
	}

	let {
		x = 'x',
		y = 'y',
		searchRadius = undefined,
		dataset = undefined,
		children
	}: Props = $props();

	let xGetter = $derived(x === 'x' ? $xGet : $yGet);
	let yGetter = $derived(y === 'y' ? $yGet : $xGet);

	function findItem(evt: any) {
		e = evt;

		const xLayerKey = `layer${x.toUpperCase()}`;
		const yLayerKey = `layer${y.toUpperCase()}`;

		found = finder.find(evt[xLayerKey], evt[yLayerKey], searchRadius) || {};
		visible = Object.keys(found).length > 0;
	}

	let finder = $derived(
		quadtree()
			.extent([
				[-1, -1],
				[$width + 1, $height + 1]
			])
			.x(xGetter)
			.y(yGetter)
			.addAll(dataset || $data)
	);
</script>

<div
	class="bg"
	onmousemove={findItem}
	onmouseout={() => (visible = false)}
	onblur={() => (visible = false)}
></div>
{@render children?.({ x: xGetter(found) || 0, y: yGetter(found) || 0, found, visible, e })}

<style>
	.bg {
		position: absolute;
		top: 0;
		right: 0;
		bottom: 0;
		left: 0;
	}
</style>
