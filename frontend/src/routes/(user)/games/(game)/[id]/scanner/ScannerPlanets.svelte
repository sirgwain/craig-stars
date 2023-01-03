<!--
  @component
  Generates an SVG scatter plot. This component can also work if the x- or y-scale is ordinal, i.e. it has a `.bandwidth` method. See the [timeplot chart](https://layercake.graphics/example/Timeplot) for an example.
 -->
<script lang="ts">
	import {
		commandedMapObject,
		commandedPlanet,
		commandMapObject,
		getMapObjectsByPosition,
		player,
		playerColor
	} from '$lib/services/Context';
	import type { Fleet } from '$lib/types/Fleet';
	import { MapObjectType, type MapObject } from '$lib/types/MapObject';
	import { Unexplored, type Planet } from '$lib/types/Planet';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import ScannerPlanet from './ScannerPlanet.svelte';

	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');

	let planets: Planet[] = [];

	const commanded = (planet: Planet, commandedMapObject: MapObject, commandedPlanet: Planet | undefined): boolean => {
		if (
			commandedMapObject.type == MapObjectType.Fleet &&
			(commandedMapObject as Fleet).orbitingPlanetNum == planet.num
		) {
			return true;
		} else if (commandedPlanet?.num === planet.num) {
			return true;
		}
		return false;
	};

	const color = (planet: Planet): string => {
		if (planet.playerNum === $player?.num) {
			return '#00FF00';
		} else if (planet.reportAge !== Unexplored && !planet.playerNum) {
			return '#FFF';
		} else if (planet.playerNum) {
			return playerColor(planet.playerNum);
		}
		return '#555';
	};

	$: planets = $data && $data.filter((mo: MapObject) => mo.type == MapObjectType.Planet);
</script>

<!-- Planets -->
{#each planets as planet}
	<ScannerPlanet
		{planet}
		commanded={commanded(planet, $commandedMapObject, $commandedPlanet)}
		orbitingFleets={getMapObjectsByPosition(planet).length > 1}
	/>
{/each}
