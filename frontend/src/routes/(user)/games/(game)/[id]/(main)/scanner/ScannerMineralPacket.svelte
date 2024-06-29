<!--
  @component
  MineralPacket in the scanner
 -->
<script lang="ts">
	import type { MineralPacket } from '$lib/types/MineralPacket';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import MapObjectScaler from './MapObjectScaler.svelte';

	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');

	export let mineralPacket: MineralPacket;
	export let color = '#0900FF';

	const size = 10;
</script>

<!-- ScannerMineralPacket -->
<MapObjectScaler mapObject={mineralPacket}>
	<rect x={-size / 2} y={-size / 2} width={size} height={size} stroke={color} fill="none" stroke-width={2} />
</MapObjectScaler>
