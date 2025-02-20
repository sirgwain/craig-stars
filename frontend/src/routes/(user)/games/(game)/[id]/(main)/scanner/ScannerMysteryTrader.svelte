<!--
  @component
  A mysterytrader that is flying outside of a planet
 -->
<script lang="ts">
	import { radiansToDegrees } from '$lib/services/Math';
	import type { MysteryTraderIntel } from '$lib/types/cs';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import { getScannerContext } from './Scanner';

	const { xGet, yGet } = getContext<LayerCake>('LayerCake');
	const { scale } = getScannerContext();

	// identity or default is rotated 90º, or pointing up and to the right
	const angleOffset = 225;

	type Props = {
		mysteryTrader: MysteryTraderIntel;
	};

	let { mysteryTrader }: Props = $props();

	let angle = $derived.by(() => {
		if (!mysteryTrader || !mysteryTrader.heading) {
			return 0;
		}
		return (
			radiansToDegrees(
				// Math.atan2(determinant(startHeading, mysterytrader.heading), dot(startHeading, mysterytrader.heading))
				Math.atan2(mysteryTrader.heading.y, mysteryTrader.heading.x)
			) + angleOffset
		);
	});

	let size = $derived(8 / $scale);
</script>

<!-- ScannerMysteryTrader -->
<polygon
	class="fill-mystery-trader"
	points={`0,0 0,${size} ${size},${size}`}
	transform={`translate(${$xGet(mysteryTrader)} ${$yGet(mysteryTrader)}) rotate(${angle}) translate(${-size / 2} ${
		-size / 2
	})`}
/>
