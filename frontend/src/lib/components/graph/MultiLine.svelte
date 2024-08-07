<!--
  @component
  Generates an SVG multi-series line chart. It expects your data to be an array of objects, each with a `values` key that is an array of data objects.
 -->
<script lang="ts">
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';

	export let colorIsCode = true;

	const { data, xGet, yGet, zGet, xScale, yScale, width, height } =
		getContext<LayerCake>('LayerCake');

	$: path = (values: any) => {
		return (
			'M' +
			values
				.map((d: any) => {
					return $xGet(d) + ',' + $yGet(d);
				})
				.join('L')
		);
	};
</script>

<g class="line-group">
	{#each $data as group}
		<path
			stroke={colorIsCode ? $zGet(group) : undefined}
			class="path-line {colorIsCode ?? $zGet(group)}"
			d={path(group.values)}
		/>
	{/each}
</g>

<style>
	.path-line {
		fill: none;
		stroke-linejoin: round;
		stroke-linecap: round;
		stroke-width: 3px;
	}
</style>
