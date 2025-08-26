<!--
  @component
  A mysterytrader that is flying outside of a planet
 -->
<script lang="ts">
	import { radiansToDegrees } from '$lib/services/Math';
	import type { MysteryTrader } from '$lib/types/cs-proto';
	import MapObjectScaler from './MapObjectScaler.svelte';

	// identity or default is rotated 90º, or pointing up and to the right
	const angleOffset = 225;

	type Props = {
		mysteryTrader: MysteryTrader;
	};

	let { mysteryTrader }: Props = $props();

	function getAngle(mysteryTrader: MysteryTrader): number {
		return (
			radiansToDegrees(Math.atan2(mysteryTrader.heading?.y ?? 0, mysteryTrader.heading?.x ?? 0)) +
			angleOffset
		);
	}

	let size = 8;
</script>

<!-- ScannerMysteryTrader -->
<MapObjectScaler mapObject={mysteryTrader}>
	<polygon
		class="fill-mystery-trader"
		points={`0,0 0,${size} ${size},${size}`}
		transform={`rotate(${getAngle(mysteryTrader)}) translate(${-size / 2} ${-size / 2})`}
	/>
</MapObjectScaler>
