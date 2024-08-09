<!--
  @component
  Generates an SVG scatter plot. This component can also work if the x- or y-scale is ordinal, i.e. it has a `.bandwidth` method. See the [timeplot chart](https://layercake.graphics/example/Timeplot) for an example.
 -->
<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import type { ZoomTransform } from 'd3-zoom';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import type { Readable } from 'svelte/store';

	const { universe, selectedMapObject, highlightedMapObject } = getGameContext();
	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');
	const objectScale = getContext<Readable<number>>('objectScale');

	export let transform: ZoomTransform;

	function fillStyle(left: number, top: number) {
		return `top:${top}px; left: ${left}px;`;
	}

</script>

<!-- Names -->
{#each $universe.planets as planet}
	{#if $highlightedMapObject == planet || $selectedMapObject == planet || $objectScale >= 5}
		<div
			class="absolute w-32 text-center ml-[-4rem] mt-2 pointer-events-none z-20 text-white"
			style={fillStyle(transform.applyX($xGet(planet)), transform.applyY($yGet(planet)))}
		>
			<span class="select-none">{planet.name}</span>
		</div>
	{/if}
{/each}
