<script lang="ts">
	import type { Minefield } from '$lib/types/cs-proto';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';

	const { xGet, yGet, xScale, yScale } = getContext<LayerCake>('LayerCake');

	type Props = {
		minefield: Minefield;
		color?: string;
		selected?: boolean;
	};

	let { minefield, color = '#0900FF', selected = false }: Props = $props();
</script>

<circle
	cx={$xGet(minefield.mapObject)}
	cy={$yGet(minefield.mapObject)}
	r={$xScale(Math.sqrt(minefield.numMines))}
	mask="url(#mask-minefield)"
	fill={color}
	class:selected
/>

{#if selected}
	<rect
		width={$xScale(2)}
		height={$yScale(2)}
		rx={0.5}
		x={$xGet(minefield.mapObject) - $xScale(1)}
		y={$yGet(minefield.mapObject) - $yScale(1)}
		fill={color}
	/>
{/if}

<style>
	.selected {
		filter: brightness(0.6);
	}
</style>
