<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import type { Planet } from '$lib/types/Planet';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';
	import ScannerPlanetNormal from './ScannerPlanetNormal.svelte';

	const { game, player, universe, settings } = getGameContext();
	const scale = getContext<Writable<number>>('scale');
	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');

	export let planet: Planet;

	let props = {};

	const fullyPopulatedRadius = 10;
	const minRadius = 3;

	const population = planet.spec.population ?? 0;

	$: {
		// green for us, gray for unexplored, white for explored
		let color = '#555';

		if (population > 0) {
			let radius = Math.max((population / 1_000_000) * fullyPopulatedRadius, minRadius);

			if (planet.playerNum) {
				color = $universe.getPlayerColor(planet.playerNum) ?? '#FF0000';
			}

			// setup the properties of our planet circle
			props = {
				r: $xScale(radius),
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
