<script lang="ts">
	import type { AnyMineralPacket } from '$lib/services/Universe';
	import MapObjectScaler from './MapObjectScaler.svelte';

	type Props = {
		mineralPacket: AnyMineralPacket;
		color?: string;
	};

	let { mineralPacket, color = '#0900FF' }: Props = $props();

	const size = 10;
</script>

<MapObjectScaler mapObject={mineralPacket}>
	<rect
		x={-size / 2}
		y={-size / 2}
		width={size}
		height={size}
		stroke={color}
		fill="none"
		stroke-width={2}
	/>
</MapObjectScaler>
