<!--
  @component
  Generates an SVG multi-series line chart. It expects your data to be an array of objects, each with a `values` key that is an array of data objects.
 -->
<script lang="ts">
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';

	interface Props {
		zIsColorCode?: boolean;
	}

	let { zIsColorCode = false }: Props = $props();

	const { data, xGet, yGet, zGet, xScale, yScale, width, height } =
		getContext<LayerCake>('LayerCake');

	let path = $derived((values: any) => {
		return (
			'M' +
			values
				.map((d: any) => {
					return $xGet(d) + ',' + $yGet(d);
				})
				.join('L')
		);
	});
</script>

<g class="line-group">
	{#each $data as group}
		<path
			stroke={zIsColorCode ? $zGet(group) : undefined}
			class="path-line {zIsColorCode ? undefined : $zGet(group)}"
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
