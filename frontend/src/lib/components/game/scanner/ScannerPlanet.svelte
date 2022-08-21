<!--
  @component
  Generates an SVG scatter plot. This component can also work if the x- or y-scale is ordinal, i.e. it has a `.bandwidth` method. See the [timeplot chart](https://layercake.graphics/example/Timeplot) for an example.
 -->
<script lang="ts">
	import {
	commandedMapObject,mapObjectsByPosition,
	myMapObjectsByPosition,player
	} from '$lib/services/Context';
	import { positionKey } from '$lib/types/MapObject';
	import { Unexplored,type Planet } from '$lib/types/Planet';
	import { getContext } from 'svelte';

	const { data, xGet, yGet, xScale, yScale, width, height } = getContext('LayerCake');

	export let planet: Planet;

	let commanded = false;
	let props = {};
	let ringProps: any | undefined = undefined;

	$: {
		if ($mapObjectsByPosition && $myMapObjectsByPosition && $player) {
			// green for us, gray for unexplored, white for explored
			let color = '#555';
			let strokeWidth = 0;
			let strokeColor = '#555';

			const key = positionKey(planet);
			const mapObjectsAtPosition = $mapObjectsByPosition[key];
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

			if (planet.playerNum === $player.num) {
				color = '#00FF00';
				strokeWidth = $xScale(2);
			} else if ((planet as Planet)?.reportAge !== Unexplored) {
				color = '#FFF';
			}

			// if anything is orbiting our planet, put a ring on it
			if (mapObjectsAtPosition.length > 1) {
				ringProps = {
					cx: $xGet(planet),
					cy: $yGet(planet),
					stroke: '#fff',
					'stroke-width': $xScale(1),
					r: $xScale(1) * (commanded ? 10 : 5),
					'fill-opacity': 0
				};
			} else {
				ringProps = undefined;
			}

			// setup the properties of our planet circle
			props = {
				r: $xScale(1) * (commanded ? 7 : 3),
				fill: color,
				stroke: strokeColor,
				'stroke-width': strokeWidth
			};
		}
	}
</script>

{#if ringProps}
	<circle {...ringProps} />
{/if}
<circle cx={$xGet(planet)} cy={$yGet(planet)} {...props} />
