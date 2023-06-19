<!--
  @component
  Generates an SVG scatter plot. This component can also work if the x- or y-scale is ordinal, i.e. it has a `.bandwidth` method. See the [timeplot chart](https://layercake.graphics/example/Timeplot) for an example.
 -->
<script lang="ts">
	import { highlightedMapObject } from '$lib/services/Context';
	import type { FullGame } from '$lib/services/FullGame';
	import { MapObjectType, type MapObject } from '$lib/types/MapObject';
	import type { Planet } from '$lib/types/Planet';
	import type { ZoomTransform } from 'd3-zoom';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';

	const game = getContext<FullGame>('game');
	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');
	const scale = getContext<Writable<number>>('scale');

	export let transform: ZoomTransform;

	$: planets = $data && $data.filter((mo: MapObject) => mo.type == MapObjectType.Planet);

	function fillStyle(left: number, top: number) {
		return `top:${top}px; left: ${left}px;`;
	}
</script>

<!-- Names -->
{#each planets as planet}
	{#if $highlightedMapObject == planet}
		<div
			class="absolute w-32 text-center ml-[-4rem] mt-2 pointer-events-none z-20"
			style={fillStyle(transform.applyX($xGet(planet)), transform.applyY($yGet(planet)))}
		>
			<span class="select-none">{planet.name}</span>
			<!-- <button class="btn btn-xs btn-secondary btn-outline hover:btn-primary pointer-events-auto">blah</button> -->
		</div>
	{/if}
{/each}
