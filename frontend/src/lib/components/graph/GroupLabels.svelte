<!--
  @component
  Generates HTML text labels for a nested data structure. It places the label near the y-value of the highest x-valued data point. This is useful for labeling the final point in a multi-series line chart, for example. It expects your data to be an array of objects where each has `values` field that is an array of data objects. It uses the `z` field accessor to pull the text label.
 -->
<script lang="ts">
	import { getContext } from 'svelte';
	import { max } from 'd3-array';
	import { LayerCake } from 'layercake';

	const { data, x, y, xScale, yScale, xRange, yRange, z } = getContext<LayerCake>('LayerCake');

	/* --------------------------------------------
	 * Title case the first letter
	 */
	const cap = (val: string) => val.replace(/^\w/, (d) => d.toUpperCase());

	/* --------------------------------------------
	 * Put the label on the highest value
	 */
	let left = $derived((values: number[]) => $xScale(max(values, $x)) / Math.max(...$xRange));
	let top = $derived((values: number[]) => $yScale(max(values, $y)) / Math.max(...$yRange));
</script>

{#each $data as group (group)}
	<div
		class="label"
		style="
        top:{top(group.values) * 100}%;
        left:{left(group.values) * 100}%;
      "
	>
		{cap($z(group))}
	</div>
{/each}

<style>
	.label {
		position: absolute;
		transform: translate(-100%, -100%) translateY(1px);
		font-size: 13px;
	}
</style>
