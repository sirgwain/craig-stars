<!--
  @component
  A fleet that is flying outside of a planet
 -->
<script lang="ts">
	import { game } from '$lib/services/Stores';
	import { radiansToDegrees } from '$lib/services/Math';
	import type { Fleet } from '$lib/types/Fleet';
	import { None } from '$lib/types/MapObject';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';

	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');
	const scale = getContext<Writable<number>>('scale');

	export let fleet: Fleet;
	export let commanded = false;
	export let color = '#0000FF';
	export let commandedColor = '#FFFF00';

	let angle = 0;

	$: size = 8 / $scale;

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
</script>

<!-- ScannerFleet -->
<polygon
	points={`0,0 0,${size} ${size},${size}`}
	fill={commanded ? commandedColor : color}
	transform={`translate(${$xGet(fleet)} ${$yGet(fleet)}) rotate(${angle}) translate(${-size / 2} ${
		-size / 2
	})`}
/>
