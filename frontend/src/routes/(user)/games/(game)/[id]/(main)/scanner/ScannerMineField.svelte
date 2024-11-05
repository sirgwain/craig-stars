<!--
  @component
  Minefield in the scanner
 -->
<script lang="ts">
	import type { MineField } from '$lib/types/MineField';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';

	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');

	export let mineField: MineField;
	export let color = '#0900FF';
	export let selected = false;
</script>

<!-- ScannerMineField -->
<circle
	cx={$xGet(mineField)}
	cy={$yGet(mineField)}
	r={$xScale(mineField.spec.radius)}
	mask="url(#mask-minefield)"
	fill={color}
	class:selected
/>

{#if selected}
	<rect
		width={$xScale(2)}
		height={$yScale(2)}
		rx={0.5}
		x={$xGet(mineField) - $xScale(1)}
		y={$yGet(mineField) - $yScale(1)}
		fill={color}
	/>
{/if}

<style>
	.selected {
		filter: brightness(0.6);
	}
</style>
