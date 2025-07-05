<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import type { ZoomTransform } from 'd3-zoom';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import { getScannerContext } from './Scanner';

	const { universe, selectedMapObject, highlightedMapObject } = getGameContext();
	const { xGet, yGet } = getContext<LayerCake>('LayerCake');
	const { objectScale } = getScannerContext();

	type Props = {
		transform: ZoomTransform;
	};

	let { transform }: Props = $props();

	function fillStyle(left: number, top: number) {
		return `top:${top}px; left: ${left}px;`;
	}
</script>

<!-- Names -->
{#each $universe.planetIntels as planet (planet.num)}
	{#if $highlightedMapObject == planet || $selectedMapObject == planet || $objectScale >= 5}
		<div
			class="absolute w-32 text-center ml-[-4rem] mt-2 pointer-events-none z-10 text-white"
			style={fillStyle(transform.applyX($xGet(planet)), transform.applyY($yGet(planet)))}
		>
			<span class="select-none">{planet.name}</span>
		</div>
	{/if}
{/each}
