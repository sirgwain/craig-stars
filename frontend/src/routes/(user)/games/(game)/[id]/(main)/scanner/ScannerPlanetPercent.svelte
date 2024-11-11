<script lang="ts">
	import { run } from 'svelte/legacy';

	import { getGameContext } from '$lib/services/GameContext';
	import { None, Unexplored } from '$lib/types/Constants';
	import { type Planet } from '$lib/types/Planet';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import MapObjectScaler from './MapObjectScaler.svelte';
	import ScannerFleetCount from './ScannerPlanetFleetCount.svelte';
	import ScannerPlanetNormal from './ScannerPlanetNormal.svelte';
	import { getScannerContext } from './Scanner';

	const { game, player, universe, settings } = getGameContext();
	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');
	const { scale } = getScannerContext();

	type Props = {
		planet: Planet;
	};

	let { planet }: Props = $props();

	let props = $state({});
	let flagColor = $state('#555');

	let planetX = $derived($xGet(planet));
	let planetY = $derived($yGet(planet));

	// area of cirlce is a = πr^2, so r = √(a/π)
	const fullyHabitableRadius = 15;
	const fullyHabitableArea = Math.PI * fullyHabitableRadius * fullyHabitableRadius;
	const minRadius = 3;
	const minArea = Math.PI * minRadius * minRadius;
	let radius = $state(minRadius);

	run(() => {
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
	});
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

	<ScannerFleetCount {planet} yOffset={radius - 5} />
{:else}
	<ScannerPlanetNormal {planet} />
{/if}
