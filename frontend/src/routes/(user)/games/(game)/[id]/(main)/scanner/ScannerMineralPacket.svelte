<!--
  @component
  MineralPacket in the scanner
 -->
<script lang="ts">
	import type { MineralPacket } from '$lib/types/MineralPacket';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';

	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');
	const scale = getContext<Writable<number>>('scale');

	export let mineralPacket: MineralPacket;
	export let color = '#0900FF';

	$: size = 8 / $scale;
</script>

<!-- ScannerMineralPacket -->
<rect
	x={$xGet({ position: mineralPacket.position }) - size / 2}
	y={$yGet({ position: mineralPacket.position }) - size / 2}
	width={size}
	height={size}
	stroke={color}
	fill="none"
/>
