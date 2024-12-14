<script lang="ts">
	import { PlanetViewState } from '$lib/types/PlayerSettings';

	import SelectedMapObject from '$lib/components/icons/SelectedMapObject.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import type { Fleet } from '$lib/types/Fleet';
	import { MapObjectType, equal, type MapObject } from '$lib/types/MapObject';
	import MapObjectScaler from './MapObjectScaler.svelte';

	const { selectedMapObject, commandedMapObject, settings } = getGameContext();

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

	let size = $derived.by(() => {
		switch ($settings.planetViewState) {
			case PlanetViewState.Normal:
			case PlanetViewState.SurfaceMinerals:
			case PlanetViewState.MineralConcentration:
				return commanded($selectedMapObject, $commandedMapObject) ? 15 : 10;
			case PlanetViewState.Percent:
			case PlanetViewState.Population:
			case PlanetViewState.None:
				return 21;
			default:
				return 10;
		}
	});
</script>

{#if $selectedMapObject}
	<MapObjectScaler mapObject={$selectedMapObject}>
		<SelectedMapObject x={-size / 2} y={size * 0.5} width={size} height={size} />
	</MapObjectScaler>
{/if}
