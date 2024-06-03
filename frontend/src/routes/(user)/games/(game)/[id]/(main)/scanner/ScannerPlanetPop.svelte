<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import type { Planet } from '$lib/types/Planet';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';
	import ScannerPlanetNormal from './ScannerPlanetNormal.svelte';
	import ScannerFleetCount from './ScannerPlanetFleetCount.svelte';
	import MapObjectScaler from './MapObjectScaler.svelte';

	const { game, player, universe, settings } = getGameContext();
	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');

	export let planet: Planet;

	let props = {};

	const fullyPopulatedRadius = 18;
	const fullyPopulatedArea = Math.PI * fullyPopulatedRadius * fullyPopulatedRadius;
	const minRadius = 2;
	const minArea = Math.PI * minRadius * minRadius;
	const population = planet.spec.population ?? 0;
	$: radius = 0;

	$: {
		// green for us, gray for unexplored, white for explored
		let color = '#555';
		let strokeColor = '#555';
		let strokeWidth = 1;

		if (population > 0) {
			radius = Math.sqrt(
				Math.max((population / 1_300_000) * fullyPopulatedArea, minArea) / Math.PI
			);

			strokeWidth = (population / 1_300_000) * strokeWidth;

			if (planet.playerNum) {
				color = $universe.getPlayerColor(planet.playerNum) ?? '#FF0000';
			}

			// setup the properties of our planet circle
			props = {
				r: radius,
				fill: color,
				stroke: strokeColor,
				'stroke-width': strokeWidth
			};
		}
	}
</script>

{#if population > 0}
	<MapObjectScaler mapObject={planet}>
		<circle cx={0} cy={0} {...props} />
	</MapObjectScaler>
	<ScannerFleetCount {planet} yOffset={radius-3} />
{:else}
	<ScannerPlanetNormal {planet} />
{/if}
