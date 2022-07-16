<!--
  @component
  Generates an SVG scatter plot. This component can also work if the x- or y-scale is ordinal, i.e. it has a `.bandwidth` method. See the [timeplot chart](https://layercake.graphics/example/Timeplot) for an example.
 -->
<script lang="ts">
	import { selectedMapObject } from '$lib/services/Context';
	import { getContext } from 'svelte';

	const { xGet, yGet } = getContext('LayerCake');
</script>

<svg
	x={$xGet($selectedMapObject) - 10}
	y={$yGet($selectedMapObject)}
	fill="none"
	stroke="currentColor"
	stroke-width="2"
>
	<path stroke-linecap="round" stroke-linejoin="round" d="M5 15l5-5 5 5" />
</svg>
