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
	import { ownedBy } from '$lib/types/MapObject';
	import { minus, normalized } from '$lib/types/Vector';

	const { settings, player } = getGameContext();
	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');
	const scale = getContext<Writable<number>>('scale');

	export let fleet: Fleet;
	export let commanded = false;
	export let color = '#0000FF';
	export let commandedColor = '#FFFF00';

	const size = 8;

	// identity or default is rotated 90º, or pointing up and to the right
	const angleOffset = 225;

	function getAngle(fleet: Fleet): number {
		return radiansToDegrees(Math.atan2(fleet.heading.y, fleet.heading.x)) + angleOffset;
	}

	function getTokenCount(fleet: Fleet): number {
		return fleet.tokens ? fleet.tokens.reduce((count, t) => count + t.quantity, 0) : 0;
	}
</script>

<!-- ScannerFleet -->
<MapObjectScaler mapObject={fleet}>
	<polygon
		points={`0,0 0,${size} ${size},${size}`}
		fill={commanded ? commandedColor : color}
		transform={`rotate(${getAngle(fleet)}) translate(${-size / 2} ${-size / 2})`}
	/>
	{#if $settings.showFleetTokenCounts}
		<!-- position the text below the fleet -->
		<text transform={`translate(0 ${size * 2.5})`} text-anchor="middle" class="fill-white"
			>{getTokenCount(fleet)}</text
		>
	{/if}
</MapObjectScaler>
