<!--
  @component
  A mysterytrader that is flying outside of a planet
 -->
<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { radiansToDegrees } from '$lib/services/Math';
	import type { MysteryTrader } from '$lib/types/MysteryTrader';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';

	const { xGet, yGet } = getContext<LayerCake>('LayerCake');
	const scale = getContext<Writable<number>>('scale');

	export let mysteryTrader: MysteryTrader;
	export let color = '#00FFFF';

	let angle = 0;

	$: size = 8 / $scale;

	// identity or default is rotated 90º, or pointing up and to the right
	const angleOffset = 225;

	$: {
		if (mysteryTrader && mysteryTrader.heading) {
			angle =
				radiansToDegrees(
					// Math.atan2(determinant(startHeading, mysterytrader.heading), dot(startHeading, mysterytrader.heading))
					Math.atan2(mysteryTrader.heading.y, mysteryTrader.heading.x)
				) + angleOffset;
		}
	}
</script>

<!-- ScannerMysteryTrader -->
<polygon
	points={`0,0 0,${size} ${size},${size}`}
	fill={color}
	transform={`translate(${$xGet(mysteryTrader)} ${$yGet(mysteryTrader)}) rotate(${angle}) translate(${-size / 2} ${
		-size / 2
	})`}
/>
