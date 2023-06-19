<!--
  @component
  Generates an SVG scatter plot. This component can also work if the x- or y-scale is ordinal, i.e. it has a `.bandwidth` method. See the [timeplot chart](https://layercake.graphics/example/Timeplot) for an example.
 -->
<script lang="ts">
	import {
		commandedMapObject,
		mapObjectsByPosition,
		myMapObjectsByPosition,
		player
	} from '$lib/services/Context';
	import { positionKey } from '$lib/types/MapObject';
	import type { Fleet } from '$lib/types/Fleet';
	import { getContext } from 'svelte';
	import { determinant, dot, normalized } from '$lib/types/Vector';
	import { radiansToDegrees } from '$lib/services/Math';

	const { data, xGet, yGet, xScale, yScale, width, height } = getContext('LayerCake');

	export let fleet: Fleet;

	let commanded = false;
	let color = '#FF0000';
	let angle = 0;

	let size = 5;
	$: size = $xScale(5);

	// the drawn svg triangle points up and to the left
	const startHeading = normalized({ x: -1, y: 1 });

	// identity or default is rotated 90º, or pointing up and to the right
	const angleOffset = 45;

	$: {
		if (fleet && fleet.heading) {
			angle =
				radiansToDegrees(
					// Math.atan2(determinant(startHeading, fleet.heading), dot(startHeading, fleet.heading))
					Math.atan2(fleet.position.y + fleet.heading.y, fleet.position.x + fleet.heading.x)
				) + angleOffset;
		}
	}

	$: {
		if ($mapObjectsByPosition && $myMapObjectsByPosition && $player) {
			const key = positionKey(fleet);
			const myMapObjectsAtPosition = $myMapObjectsByPosition[key];
			if (
				myMapObjectsAtPosition?.find(
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

{#if !fleet.orbiting}
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
