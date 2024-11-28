<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import type { Planet } from '$lib/types/Planet';
	import MapObjectScaler from './MapObjectScaler.svelte';
	import ScannerFleetCount from './ScannerPlanetFleetCount.svelte';
	import ScannerPlanetNormal from './ScannerPlanetNormal.svelte';

	const { universe } = getGameContext();

	type Props = {
		planet: Planet;
	};

	let { planet }: Props = $props();

	const fullyPopulatedRadius = 18;
	const fullyPopulatedArea = Math.PI * fullyPopulatedRadius * fullyPopulatedRadius;
	const minRadius = 2;
	const minArea = Math.PI * minRadius * minRadius;

	let planetProps = $derived.by(() => {
		const population = planet.spec.population ?? 0;
		if (population <= 0) {
			return {
				radius: 0,
				circleProps: {}
			};
		}

		// green for us, gray for unexplored, white for explored
		let color = '#555';
		let strokeColor = '#555';

		let radius = Math.sqrt(
			Math.max((population / 1_300_000) * fullyPopulatedArea, minArea) / Math.PI
		);
		let strokeWidth = population / 1_300_000;

		if (planet.playerNum) {
			color = $universe.getPlayerColor(planet.playerNum) ?? '#FF0000';
		}

		// setup the properties of our planet circle
		return {
			radius,
			circleProps: {
				r: radius,
				fill: color,
				stroke: strokeColor,
				'stroke-width': strokeWidth
			}
		};
	});
</script>

{#if planetProps.radius}
	<MapObjectScaler mapObject={planet}>
		<circle cx={0} cy={0} {...planetProps.circleProps} />
	</MapObjectScaler>
	<ScannerFleetCount {planet} yOffset={planetProps.radius - 3} />
{:else}
	<ScannerPlanetNormal {planet} />
{/if}
