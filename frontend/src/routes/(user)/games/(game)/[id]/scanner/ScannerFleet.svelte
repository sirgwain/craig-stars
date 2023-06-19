<!--
  @component
  Generates an SVG scatter plot. This component can also work if the x- or y-scale is ordinal, i.e. it has a `.bandwidth` method. See the [timeplot chart](https://layercake.graphics/example/Timeplot) for an example.
 -->
<script lang="ts">
	import {
		commandedMapObject,
		getMyMapObjectsByPosition,
		player
	} from '$lib/services/Context';
	import { radiansToDegrees } from '$lib/services/Math';
	import type { Fleet } from '$lib/types/Fleet';
	import { normalized } from '$lib/types/Vector';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';

	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');

	export let fleet: Fleet;

	let commanded = false;
	let color = '#FF0000';
	let angle = 0;

	let size = 8;
	$: size = $xScale(8);

	// the drawn svg triangle points up and to the left
	const startHeading = normalized({ x: -1, y: 1 });

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

	$: {
		if ($player) {
			if (
				getMyMapObjectsByPosition(fleet).find(
					(m) => m.type == $commandedMapObject.type && m.id == $commandedMapObject.id
				)
			) {
				commanded = true;
			} else {
				commanded = false;
			}

			if (fleet.playerNum === $player.num) {
				color = commanded ? '#FFFF00' : $player.color;
			}
		}
	}
</script>

{#if !fleet.orbitingPlanetNum}
	<polygon
		id="fleet"
		points={`0,0 0,${size} ${size},${size}`}
		xlink:href="#fleet"
		fill={color}
		transform={`translate(${$xGet(fleet)} ${$yGet(fleet)}) rotate(${angle}) translate(${
			-size / 2
		} ${-size / 2})`}
	/>
{/if}
