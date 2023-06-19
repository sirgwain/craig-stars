<!--
  @component
  Generates an SVG scatter plot. This component can also work if the x- or y-scale is ordinal, i.e. it has a `.bandwidth` method. See the [timeplot chart](https://layercake.graphics/example/Timeplot) for an example.
 -->
<script lang="ts">
	import { player, playerColor } from '$lib/services/Context';
	import { Unexplored, type Planet } from '$lib/types/Planet';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';

	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');

	export let planet: Planet;
	export let commanded = false;
	export let orbitingFleets = false;

	let props = {};
	let ringProps: any | undefined = undefined;

	$: {
		if ($player) {
			// green for us, gray for unexplored, white for explored
			let color = '#555';
			let strokeWidth = 0;
			let strokeColor = '#555';

			if (planet.playerNum === $player.num) {
				color = '#00FF00';
				strokeWidth = 2;
			} else if (planet.reportAge !== Unexplored && !planet.playerNum) {
				color = '#FFF';
			} else if (planet.playerNum) {
				color = playerColor(planet.playerNum);
			}

			// if anything is orbiting our planet, put a ring on it
			if (orbitingFleets) {
				ringProps = {
					cx: $xGet(planet),
					cy: $yGet(planet),
					stroke: '#fff',
					'stroke-width': 1,
					r: 1 * (commanded ? 10 : 5),
					'fill-opacity': 0
				};
			} else {
				ringProps = undefined;
			}

			// setup the properties of our planet circle
			props = {
				r: 1 * (commanded ? 7 : 3),
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
