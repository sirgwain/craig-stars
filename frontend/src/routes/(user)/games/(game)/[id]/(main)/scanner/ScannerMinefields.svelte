<!--
  @component
  Show all minefields in the universe
 -->
<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { MapObjectType, type Minefield } from '$lib/types/cs-proto';
	import { equal, key } from '$lib/types/MapObject';
	import { getDisplayColor } from '$lib/utils/colorUtils';
	import ScannerMinefield from './ScannerMinefield.svelte';

	const { player, universe, selectedMapObject, settings } = getGameContext();

	function getColor(minefield: Minefield) {
		return getDisplayColor(minefield.mapObject?.playerNum, $player, $universe, $settings);
	}

	let minefields = $derived($universe.allMinefields);
	let selectedMinefield = $derived(
		$selectedMapObject && $selectedMapObject.mapObject?.type === MapObjectType.MINEFIELD
			? ($selectedMapObject as Minefield)
			: undefined
	);
</script>

<!-- Minefields -->
{#each minefields as minefield (key(minefield))}
	{#if minefield !== selectedMinefield}
		<ScannerMinefield
			{minefield}
			color={getColor(minefield)}
			selected={equal($selectedMapObject, minefield)}
		/>
	{/if}
{/each}
{#if selectedMinefield}
	<ScannerMinefield
		minefield={selectedMinefield}
		color={getColor(selectedMinefield)}
		selected={true}
	/>
{/if}
