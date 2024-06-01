<!--
  @component
  A fleet that is flying outside of a planet
 -->
<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { radiansToDegrees } from '$lib/services/Math';
	import type { Fleet } from '$lib/types/Fleet';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';
	import MapObjectScaler from './MapObjectScaler.svelte';

	const { settings } = getGameContext();
	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');
	const scale = getContext<Writable<number>>('scale');

	export let fleet: Fleet;
	export let commanded = false;
	export let color = '#0000FF';
	export let commandedColor = '#FFFF00';

	let angle = 0;

	const size = 8;

	// identity or default is rotated 90º, or pointing up and to the right
	const angleOffset = 225;

	$: {
		if (fleet && fleet.heading) {
			angle =
				radiansToDegrees(
					// Math.atan2(determinant(startHeading, fleet.heading), dot(startHeading, fleet.heading))
					Math.atan2(fleet.heading.y, fleet.heading.x)
				) + angleOffset;
		}
	}

	$: tokenCount = fleet.tokens ? fleet.tokens.reduce((count, t) => count + t.quantity, 0) : 0;
</script>

<!-- ScannerFleet -->
<MapObjectScaler mapObject={fleet}>
	<polygon
		points={`0,0 0,${size} ${size},${size}`}
		fill={commanded ? commandedColor : color}
		transform={`rotate(${angle}) translate(${-size / 2} ${-size / 2})`}
	/>
	{#if $settings.showFleetTokenCounts}
		<!-- position the text below the fleet -->
		<text transform={`translate(0 ${size * 2.5})`} text-anchor="middle" class="fill-base-content"
			>{tokenCount}</text
		>
	{/if}
</MapObjectScaler>
