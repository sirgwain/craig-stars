<!--
  @component
  Generates an SVG scatter plot. This component can also work if the x- or y-scale is ordinal, i.e. it has a `.bandwidth` method. See the [timeplot chart](https://layercake.graphics/example/Timeplot) for an example.
 -->
<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import { commandedMapObject, commandedPlanet } from '$lib/services/Stores';
	import type { Fleet } from '$lib/types/Fleet';
	import { MapObjectType, type MapObject } from '$lib/types/MapObject';
	import type { Planet } from '$lib/types/Planet';
	import { PlanetViewState } from '$lib/types/PlayerSettings';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import ScannerPlanetNormal from './ScannerPlanetNormal.svelte';
	import ScannerPlanetPercent from './ScannerPlanetPercent.svelte';
	import ScannerPlanetPop from './ScannerPlanetPop.svelte';
	import ScannerPlanetMineralConcentration from './ScannerPlanetMineralConcentration.svelte';
	import ScannerPlanetSurfaceMinerals from './ScannerPlanetSurfaceMinerals.svelte';

	const { game, player, universe, settings } = getGameContext();
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
	{:else if $settings.planetViewState == PlanetViewState.Population}
		<ScannerPlanetPop {planet} />
	{:else if $settings.planetViewState == PlanetViewState.MineralConcentration}
		<ScannerPlanetMineralConcentration {planet} />
	{:else if $settings.planetViewState == PlanetViewState.SurfaceMinerals}
		<ScannerPlanetSurfaceMinerals {planet} />
	{:else}
		<ScannerPlanetNormal
			{planet}
			commanded={commanded(planet, $commandedMapObject, $commandedPlanet)}
		/>
	{/if}
{/each}
