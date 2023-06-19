<!--
  @component
  Generates an SVG scatter plot. This component can also work if the x- or y-scale is ordinal, i.e. it has a `.bandwidth` method. See the [timeplot chart](https://layercake.graphics/example/Timeplot) for an example.
 -->
<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import { settings } from '$lib/services/Settings';
	import { commandedMapObject, commandedPlanet } from '$lib/services/Stores';
	import type { Fleet } from '$lib/types/Fleet';
	import { MapObjectType, type MapObject } from '$lib/types/MapObject';
	import type { Planet } from '$lib/types/Planet';
	import { PlanetViewState } from '$lib/types/PlayerSettings';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import ScannerPlanetNormal from './ScannerPlanetNormal.svelte';
	import ScannerPlanetPercent from './ScannerPlanetPercent.svelte';

	const { game, player, universe } = getGameContext();
	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');

	const commanded = (
		planet: Planet,
		commandedMapObject: MapObject | undefined,
		commandedPlanet: Planet | undefined
	): boolean => {
		if (
			commandedMapObject?.type == MapObjectType.Fleet &&
			(commandedMapObject as Fleet).orbitingPlanetNum == planet.num
		) {
			return true;
		} else if (commandedPlanet?.num === planet.num) {
			return true;
		}
		return false;
	};

	$: planets = $universe.planets;
</script>

<!-- Planets -->
{#each planets as planet}
	{#if $settings.planetViewState == PlanetViewState.Percent}
		<ScannerPlanetPercent {planet} />
	{:else}
		<ScannerPlanetNormal
			{planet}
			commanded={commanded(planet, $commandedMapObject, $commandedPlanet)}
			orbitingFleets={($universe.getMapObjectsByPosition(planet) ?? []).length > 1}
		/>
	{/if}
{/each}
