<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import type { Fleet, PlanetIntel } from '$lib/types/cs';
	import { MapObjectTypeFleet, type MapObject } from '$lib/types/cs';
	import type { CommandedPlanet } from '$lib/types/Planet';
	import { PlanetViewState } from '$lib/types/PlayerSettings';
	import ScannerPlanetMineralConcentration from './ScannerPlanetMineralConcentration.svelte';
	import ScannerPlanetNormal from './ScannerPlanetNormal.svelte';
	import ScannerPlanetPercent from './ScannerPlanetPercent.svelte';
	import ScannerPlanetPop from './ScannerPlanetPop.svelte';
	import ScannerPlanetSurfaceMinerals from './ScannerPlanetSurfaceMinerals.svelte';

	const { universe, settings, commandedMapObject, commandedPlanet } = getGameContext();

	const commanded = (
		planet: PlanetIntel,
		commandedMapObject: MapObject | undefined,
		commandedPlanet: CommandedPlanet | undefined
	): boolean => {
		if (
			commandedMapObject?.type == MapObjectTypeFleet &&
			(commandedMapObject as Fleet).orbitingPlanetNum == planet.num
		) {
			return true;
		} else if (commandedPlanet?.num === planet.num) {
			return true;
		}
		return false;
	};
</script>

<!-- Planets -->
{#each $universe.planets as planet (planet.num)}
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
