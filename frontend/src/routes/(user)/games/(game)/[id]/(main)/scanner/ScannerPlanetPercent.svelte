<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import type { FullGame } from '$lib/services/FullGame';
	import { settings } from '$lib/services/Settings';
	import { None } from '$lib/types/MapObject';
	import { Unexplored, type Planet } from '$lib/types/Planet';
	import { PlanetViewState } from '$lib/types/PlayerSettings';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';

	const { game, player, universe } = getGameContext();
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

		if (planet.reportAge !== Unexplored) {
			strokeWidth = 1;
			let habitability = planet.spec?.habitability ?? 0;
			let habitabilityTerraformed = planet.spec?.habitabilityTerraformed ?? 0;
			if (habitability > 0) {
				color = '#00FF00';
				radius = (habitability / 100.0) * fullyHabitableRadius;
			} else {
				if (habitabilityTerraformed > 0) {
					color = '#00FFFF';
					radius = (habitabilityTerraformed / 100.0) * fullyHabitableRadius;
				} else {
					color = '#FF0000';
					radius = (-habitability / 45.0) * fullyHabitableRadius;
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
