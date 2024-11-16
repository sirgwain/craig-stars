<!--
  @component
  Salvage in the scanner
 -->
<script lang="ts">
	import type { Salvage } from '$lib/types/Salvage';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import MapObjectScaler from './MapObjectScaler.svelte';

	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');

	type Props = {
		salvage: Salvage;
		color?: string;
	};

	let { salvage, color = '#FFFF00' }: Props = $props();

	const size = 10;
</script>

<!-- ScannerSalvage -->
<MapObjectScaler mapObject={salvage}>
	<rect
		x={-size / 2}
		y={-size / 2}
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
</MapObjectScaler>
