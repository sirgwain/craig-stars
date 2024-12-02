<!--
  @component
  Creates an interaction layer (in HTML) using [d3-quadtree](https://github.com/d3/d3-quadtree) to find the nearest datapoint to the mouse. This component creates a slot that exposes variables `x`, `y`, `found` (the found datapoint), `visible` (a Boolean whether any data was found) and `e` (the event object).

  The quadtree searches across both the x and y dimensions at the same time. But if you want to only search across one, set the `x` and `y` props to the same value. For example, the [shared tooltip component](https://layercake.graphics/components/SharedTooltip.html.svelte) sets `y='x'` since it's nicer behavior to only pick up on the nearest x-value.
 -->
<script lang="ts">
	import { quadtree } from 'd3-quadtree';
	import type { LayerCake } from 'layercake';
	import { getContext, type Snippet } from 'svelte';
	const { data, xGet, yGet, width, height } = getContext<LayerCake>('LayerCake');

	let visible = $state(false);
	let found: [number, number] | undefined = $state();

	type Props = {
		x?: string;
		y?: string;
		/** @type {String} [searchRadius] – The number of pixels to search around the mouse's location. This is the third argument passed to [`quadtree.find`](https://github.com/d3/d3-quadtree#quadtree_find) and by default a value of `undefined` means an unlimited range. */
		searchRadius?: number | undefined;
		/** @type {Array} [dataset] – The dataset to work off of—defaults to $data if left unset. You can pass override the default here in here in case you don't want to use the main data or it's in a strange format. */
		dataset?: unknown[];
		children?: Snippet<
			[{ x: number; y: number; found: [number, number] | undefined; visible: boolean }]
		>;
	};

	let {
		x = 'x',
		y = 'y',
		searchRadius = undefined,
		dataset = undefined,
		children
	}: Props = $props();

	let xGetter = $derived(x === 'x' ? $xGet : $yGet);
	let yGetter = $derived(y === 'y' ? $yGet : $xGet);

	function findItem(event: MouseEvent | PointerEvent) {
		const evt = event as PointerEvent & { layerX: number; layerY: number };

		found = finder.find(evt.layerX, evt.layerY, searchRadius);
		visible = found !== undefined;
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

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
	class="bg"
	onmousemove={findItem}
	onmouseout={() => (visible = false)}
	onblur={() => (visible = false)}
></div>
{@render children?.({ x: xGetter(found) || 0, y: yGetter(found) || 0, found, visible })}

<style>
	.bg {
		position: absolute;
		top: 0;
		right: 0;
		bottom: 0;
		left: 0;
	}
</style>
