<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { None } from '$lib/types/MapObject';
	import { Unexplored, type Planet } from '$lib/types/Planet';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import ScannerFleetCount from './ScannerPlanetFleetCount.svelte';
	import ScannerPlanetNormal from './ScannerPlanetNormal.svelte';
	import type { Writable } from 'svelte/store';
	import { clamp } from '$lib/services/Math';
	import MapObjectScaler from './MapObjectScaler.svelte';

	const { game, player, universe, settings } = getGameContext();
	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');
	const scale = getContext<Writable<number>>('scale');

	export let planet: Planet;

	let props = {};
	let flagColor = '#555';

	$: planetX = $xGet(planet);
	$: planetY = $yGet(planet);

	// area of cirlce is a = πr^2, so r = √(a/π)
	const fullyHabitableRadius = 15;
	const fullyHabitableArea = Math.PI * fullyHabitableRadius * fullyHabitableRadius;
	const minRadius = 3;
	const minArea = Math.PI * minRadius * minRadius;
	let radius = minRadius;

	$: {
		// green for us, gray for unexplored, white for explored
		let color = '#555';
		let strokeWidth = 0;
		let strokeColor = '#888';
		radius = minRadius;

		if (planet.reportAge !== Unexplored) {
			strokeWidth = 1;
			let habitability = planet.spec?.habitability ?? 0;
			let habitabilityTerraformed = planet.spec?.terraformedHabitability ?? 0;
			if (habitability > 0) {
				color = '#00FF00';
				radius = Math.sqrt(
					Math.max((habitability / 100.0) * fullyHabitableArea, minArea) / Math.PI
				);
				strokeWidth = (habitability / 100.0) * strokeWidth;
			} else {
				if (habitabilityTerraformed > 0) {
					color = '#FFFF00';
					radius = Math.sqrt(
						Math.max((habitabilityTerraformed / 100.0) * fullyHabitableArea, minArea) / Math.PI
					);
					strokeWidth = (habitabilityTerraformed / 100.0) * strokeWidth;
				} else {
					color = '#FF0000';
					radius = Math.sqrt(
						Math.max((-habitability / 45.0) * fullyHabitableArea, minArea) / Math.PI
					);
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

{#if planet.reportAge !== Unexplored}
	<MapObjectScaler mapObject={planet}>
		<circle cx={0} cy={0} {...props} />
		{#if planet.playerNum != None}
			<!-- draw the flag  -->
			<rect width="12" height="10" x={0} y={-fullyHabitableRadius * 2} fill={flagColor} />
			<path
				d={`M${0}, ${0}L${0}, ${-fullyHabitableRadius * 2}`}
				stroke={flagColor}
				stroke-width={2}
			/>
		{/if}
	</MapObjectScaler>

	<ScannerFleetCount {planet} yOffset={radius-5} />
{:else}
	<ScannerPlanetNormal {planet} />
{/if}
