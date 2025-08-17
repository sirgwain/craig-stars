<!--
  @component
  Salvage in the scanner
 -->
<script lang="ts">
	import type { SalvageIntel } from '$lib/types/cs-proto';
	import MapObjectScaler from './MapObjectScaler.svelte';

	type Props = {
		salvage: SalvageIntel;
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
</MapObjectScaler>

<style>
	.rotate {
		transform-box: fill-box;
		transform-origin: center;
		transform: rotate(45deg);
	}
</style>
