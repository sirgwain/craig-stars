<script lang="ts">
	import { PlanetViewState } from '$lib/types/PlayerSettings';

	import SelectedMapObject from '$lib/components/icons/SelectedMapObject.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { equal, type MapObjectLike } from '$lib/types/MapObject';
	import MapObjectScaler from './MapObjectScaler.svelte';
	import type { Fleet } from '$lib/types/cs-proto';
	import { MapObjectType } from '$lib/types/cs-proto';

	const { selectedMapObject, commandedMapObject, settings } = getGameContext();

	const commanded = (
		selectedMapObject: MapObjectLike | undefined,
		commandedMapObject: MapObjectLike | undefined
	): boolean => {
		if (
			equal(selectedMapObject, commandedMapObject) ||
			(commandedMapObject?.mapObject?.type === MapObjectType.FLEET &&
				selectedMapObject?.mapObject?.type === MapObjectType.PLANET &&
				(commandedMapObject as Fleet).orbitingPlanetNum == selectedMapObject.mapObject?.num)
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
