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
	cx={$xGet({ position: mineField.position })}
	cy={$yGet({ position: mineField.position })}
	r={$xScale(mineField.spec.radius)}
	mask="url(#mask-minefield)"
	fill={color}
	class:selected
/>

<style>
	.selected {
		filter: brightness(0.6);
	}
</style>
