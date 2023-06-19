<!--
  @component
  Generates an SVG scatter plot. This component can also work if the x- or y-scale is ordinal, i.e. it has a `.bandwidth` method. See the [timeplot chart](https://layercake.graphics/example/Timeplot) for an example.
 -->
<script lang="ts">
	import { commandedMapObject, commandedPlanet } from '$lib/services/Stores';
	import type { Fleet } from '$lib/types/Fleet';
	import { MapObjectType, type MapObject } from '$lib/types/MapObject';
	import type { Planet } from '$lib/types/Planet';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import ScannerPlanetNormal from './ScannerPlanetNormal.svelte';
	import type { FullGame } from '$lib/services/FullGame';
	import { settings } from '$lib/services/Settings';
	import { PlanetViewState } from '$lib/types/PlayerSettings';
	import ScannerPlanetPercent from './ScannerPlanetPercent.svelte';

	const game = getContext<FullGame>('game');
	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');

	let planets: Planet[] = [];

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

	$: planets = $data && $data.filter((mo: MapObject) => mo.type == MapObjectType.Planet);
</script>

<!-- Planets -->
{#each planets as planet}
	{#if $settings.planetViewState == PlanetViewState.Percent}
		<ScannerPlanetPercent {planet} />
	{:else}
		<ScannerPlanetNormal
			{planet}
			commanded={commanded(planet, $commandedMapObject, $commandedPlanet)}
			orbitingFleets={(game.universe.getMapObjectsByPosition(planet) ?? []).length > 1}
		/>
	{/if}
{/each}
