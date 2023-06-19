<!--
  @component
  Generates an SVG scatter plot. This component can also work if the x- or y-scale is ordinal, i.e. it has a `.bandwidth` method. See the [timeplot chart](https://layercake.graphics/example/Timeplot) for an example.
 -->
<script lang="ts">
	import { radiansToDegrees } from '$lib/services/Math';
	import type { Fleet } from '$lib/types/Fleet';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';

	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');

	export let fleet: Fleet;
	export let commanded = false;
	export let color = '#0000FF';
	export let commandedColor = '#FFFF00';

	let angle = 0;

	let size = 8;
	// $: size = $xScale(8);

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

{#if !fleet.orbitingPlanetNum}
	<!-- ScannerFleet -->
	<polygon
		id="fleet"
		points={`0,0 0,${size} ${size},${size}`}
		xlink:href="#fleet"
		fill={commanded ? commandedColor : color}
		transform={`translate(${$xGet(fleet)} ${$yGet(fleet)}) rotate(${angle}) translate(${
			-size / 2
		} ${-size / 2})`}
	/>
{/if}
