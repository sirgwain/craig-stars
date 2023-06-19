<!--
  @component
  Generates an SVG scatter plot. This component can also work if the x- or y-scale is ordinal, i.e. it has a `.bandwidth` method. See the [timeplot chart](https://layercake.graphics/example/Timeplot) for an example.
 -->
<script lang="ts">
	import SelectedMapObject from '$lib/components/icons/SelectedMapObject.svelte';
	import { selectedMapObject } from '$lib/services/Context';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';

	const { xGet, yGet, xScale, yScale } = getContext<LayerCake>('LayerCake');
</script>

<SelectedMapObject
	x={$xGet($selectedMapObject) - $xScale(7.5)}
	y={$yGet($selectedMapObject) + $yScale(9)}
	width={$xScale(15)}
	height={$yScale(15)}
/>
