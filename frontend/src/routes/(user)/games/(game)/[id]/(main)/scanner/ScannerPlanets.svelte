<script lang="ts">
	import { MapObjectType } from '$lib/types/cs-proto';
	import type { Fleet } from '$lib/types/cs-proto';
	import type { Planet } from '$lib/types/cs-proto';
	import { getGameContext } from '$lib/services/GameContext';
	import type { MapObjectLike } from '$lib/types/MapObject';
	import type { CommandedPlanet } from '$lib/types/Planet';
	import { PlanetViewState } from '$lib/types/PlayerSettings';
	import ScannerPlanetMineralConcentration from './ScannerPlanetMineralConcentration.svelte';
	import ScannerPlanetNormal from './ScannerPlanetNormal.svelte';
	import ScannerPlanetPercent from './ScannerPlanetPercent.svelte';
	import ScannerPlanetPop from './ScannerPlanetPop.svelte';
	import ScannerPlanetSurfaceMinerals from './ScannerPlanetSurfaceMinerals.svelte';

	const { universe, settings, commandedMapObject, commandedPlanet } = getGameContext();

	const commanded = (
		planet: Planet,
		commandedMapObject: MapObjectLike | undefined,
		commandedPlanet: CommandedPlanet | undefined
	): boolean => {
		if (
			commandedMapObject?.mapObject?.type === MapObjectType.FLEET &&
			(commandedMapObject as Fleet).orbitingPlanetNum === planet.mapObject?.num
		) {
			return true;
		} else if (commandedPlanet?.mapObject.num === planet.mapObject?.num) {
			return true;
		}
		return false;
	};
</script>

<!-- Planets -->
{#each $universe.planetIntels as planet (planet.mapObject?.num)}
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
