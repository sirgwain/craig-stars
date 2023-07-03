<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import { None } from '$lib/types/MapObject';
	import { Unexplored, type Planet } from '$lib/types/Planet';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';

	const { game, player, universe, settings } = getGameContext();
	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');

	export let planet: Planet;

	let props = {};
	let flagColor = '#555';

	$: planetX = $xGet(planet);
	$: planetY = $yGet(planet);

	const fullyHabitableRadius = 7;

	$: {
		// green for us, gray for unexplored, white for explored
		let color = '#555';
		let strokeWidth = 0;
		let strokeColor = '#888';
		let radius = 3;
		let minRadius = 2;

		if (planet.reportAge !== Unexplored) {
			strokeWidth = 1;
			let habitability = planet.spec?.habitability ?? 0;
			let habitabilityTerraformed = planet.spec?.terraformedHabitability ?? 0;
			if (habitability > 0) {
				color = '#00FF00';
				radius = Math.max((habitability / 100.0) * fullyHabitableRadius, minRadius);
				strokeWidth = (habitability / 100.0) * strokeWidth;
			} else {
				if (habitabilityTerraformed > 0) {
					color = '#FFFF00';
					radius = Math.max((habitabilityTerraformed / 100.0) * fullyHabitableRadius, minRadius);
					strokeWidth = (habitabilityTerraformed / 100.0) * strokeWidth;
				} else {
					color = '#FF0000';
					radius = Math.max((-habitability / 45.0) * fullyHabitableRadius, minRadius);
					strokeWidth = (-habitability / 45.0) * strokeWidth;
				}
			}

			if (planet.playerNum) {
				flagColor = $universe.getPlayerColor(planet.playerNum) ?? '#FF0000';
			}
		}

		// setup the properties of our planet circle
		props = {
			r: radius,
			fill: color,
			stroke: strokeColor,
			'stroke-width': strokeWidth
		};
	}
</script>

<circle cx={planetX} cy={planetY} {...props} />

{#if planet.reportAge !== Unexplored && planet.playerNum != None}
	<!-- draw the flag  -->
	<rect width="8" height="6" x={planetX} y={planetY - fullyHabitableRadius * 2} fill={flagColor} />
	<path
		d={`M${planetX}, ${planetY}L${planetX}, ${planetY - fullyHabitableRadius * 2}`}
		stroke={flagColor}
	/>
{/if}
