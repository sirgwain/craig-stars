<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import type { Planet } from '$lib/types/Planet';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import ScannerPlanetNormal from './ScannerPlanetNormal.svelte';
	import { min } from 'date-fns';
	import type { Writable } from 'svelte/store';

	const { game, player, universe, settings } = getGameContext();
	const scale = getContext<Writable<number>>('scale');
	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');

	export let planet: Planet;
	$: population = planet.spec.population ?? 0;

	let props = {};

	const fullyPopulatedRadius = 10;
	const minRadius = 2;

	$: {
		// green for us, gray for unexplored, white for explored
		let color = '#555';
		let radius = 3;
		const population = planet.spec.population ?? 0;

		if (population > 0) {
			if (population > 0) {
				radius = Math.max((population / 1_000_000) * fullyPopulatedRadius, minRadius);
			}

			if (planet.playerNum) {
				color = $universe.getPlayerColor(planet.playerNum) ?? '#FF0000';
			}

			// setup the properties of our planet circle
			props = {
				r: radius / $scale,
				fill: color
			};
		}
	}
</script>

{#if population > 0}
	<circle cx={$xGet(planet)} cy={$yGet(planet)} {...props} />
{:else}
	<ScannerPlanetNormal {planet} />
{/if}
