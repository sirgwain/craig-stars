<!--
  @component
  Generates an SVG scatter plot. This component can also work if the x- or y-scale is ordinal, i.e. it has a `.bandwidth` method. See the [timeplot chart](https://layercake.graphics/example/Timeplot) for an example.
 -->
<script lang="ts">
	import SelectedMapObject from '$lib/components/icons/SelectedMapObject.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import type { Fleet } from '$lib/types/Fleet';
	import { MapObjectType, equal, type MapObject } from '$lib/types/MapObject';
	import type { Planet } from '$lib/types/Planet';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import MapObjectScaler from './MapObjectScaler.svelte';

	const { selectedMapObject, commandedMapObject } = getGameContext();
	const { xGet, yGet, xScale, yScale } = getContext<LayerCake>('LayerCake');

	const commanded = (
		selectedMapObject: MapObject | undefined,
		commandedMapObject: MapObject | undefined
	): boolean => {
		if (
			equal(selectedMapObject, commandedMapObject) ||
			(commandedMapObject?.type == MapObjectType.Fleet &&
				selectedMapObject?.type == MapObjectType.Planet &&
				(commandedMapObject as Fleet).orbitingPlanetNum == selectedMapObject.num)
		) {
			return true;
		}
		return false;
	};

	$: size = commanded($selectedMapObject, $commandedMapObject) ? 15 : 10;
</script>

{#if $selectedMapObject}
	<MapObjectScaler mapObject={$selectedMapObject}>
		<SelectedMapObject x={-size / 2} y={size * 0.5} width={size} height={size} />
	</MapObjectScaler>
{/if}
