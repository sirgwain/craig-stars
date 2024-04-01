<!--
  @component
  Generates an SVG scatter plot. This component can also work if the x- or y-scale is ordinal, i.e. it has a `.bandwidth` method. See the [timeplot chart](https://layercake.graphics/example/Timeplot) for an example.
 -->
<script lang="ts">
	import SelectedMapObject from '$lib/components/icons/SelectedMapObject.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';

	const { selectedMapObject } = getGameContext();
	const { xGet, yGet, xScale, yScale } = getContext<LayerCake>('LayerCake');
	const scale = getContext<Writable<number>>('scale');

	$: size = 15 / $scale;
</script>

{#if $selectedMapObject}
	<SelectedMapObject
		x={$xGet($selectedMapObject) - 7.5 / $scale}
		y={$yGet($selectedMapObject) + 9 / $scale}
		width={size}
		height={size}
	/>
{/if}
