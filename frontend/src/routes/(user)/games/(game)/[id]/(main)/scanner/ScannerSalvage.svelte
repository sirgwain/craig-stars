<!--
  @component
  Salvage in the scanner
 -->
<script lang="ts">
	import type { Salvage } from '$lib/types/Salvage';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';

	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');
	const scale = getContext<Writable<number>>('scale');

	export let salvage: Salvage;
	export let color = '#FFFF00';

	$: size = 8 / $scale;
</script>

<!-- ScannerSalvage -->
<rect
	x={$xGet(salvage) - size / 2}
	y={$yGet(salvage) - size / 2}
	width={size}
	height={size}
	stroke={color}
	fill="none"
	class="rotate"
/>

<style>
	.rotate {
		transform-box: fill-box;
		transform-origin: center;
		transform: rotate(45deg);
	}
</style>
