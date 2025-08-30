<script lang="ts">
	import type { Planet } from '$lib/types/cs-proto';
	import { getGameContext } from '$lib/services/GameContext';
	import { ReportAgeUnexplored } from '$lib/types/Consts';
	import MapObjectScaler from './MapObjectScaler.svelte';
	import ScannerFleetCount from './ScannerPlanetFleetCount.svelte';
	import ScannerPlanetNormal from './ScannerPlanetNormal.svelte';
	import { getDisplayColor } from '$lib/utils/colorUtils';

	const { universe, player, settings } = getGameContext();

	type Props = {
		planet: Planet;
	};

	let { planet }: Props = $props();

	// area of cirlce is a = πr^2, so r = √(a/π)
	const fullyHabitableRadius = 15;
	const fullyHabitableArea = Math.PI * fullyHabitableRadius * fullyHabitableRadius;
	const minRadius = 3;
	const minArea = Math.PI * minRadius * minRadius;

	let planetProps = $derived.by(() => {
		// green for us, gray for unexplored, white for explored
		let color = '#555';
		let strokeWidth = 0;
		let strokeColor = '#888';
		let radius = minRadius;
		let flagColor = color;

		if (planet.mapObject?.reportAge !== ReportAgeUnexplored) {
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

			if (planet.mapObject?.playerNum) {
				flagColor = getDisplayColor(planet.mapObject.playerNum, $player, $universe, $settings);
			}
		}

		return {
			radius,
			flagColor,
			// setup the properties of our planet circle
			circleProps: {
				r: radius,
				fill: color,
				stroke: strokeColor,
				'stroke-width': strokeWidth
			}
		};
	});
</script>

{#if planet.mapObject?.reportAge !== ReportAgeUnexplored}
	<MapObjectScaler mapObject={planet}>
		<circle cx={0} cy={0} {...planetProps.circleProps} />
		{#if planet.mapObject?.playerNum}
			<!-- draw the flag  -->
			<rect
				width="12"
				height="10"
				x={0}
				y={-fullyHabitableRadius * 2}
				fill={planetProps.flagColor}
			/>
			<path
				d={`M${0}, ${0}L${0}, ${-fullyHabitableRadius * 2}`}
				stroke={planetProps.flagColor}
				stroke-width={2}
			/>
		{/if}
	</MapObjectScaler>

	<ScannerFleetCount {planet} yOffset={planetProps.radius - 5} />
{:else}
	<ScannerPlanetNormal {planet} />
{/if}
