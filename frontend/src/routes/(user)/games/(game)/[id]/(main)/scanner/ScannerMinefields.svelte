<!--
  @component
  Show all minefields in the universe
 -->
<script lang="ts">
	import { MapObjectType } from '$lib/types/cs-proto';
	import { getGameContext } from '$lib/services/GameContext';
	import type { AnyMinefield } from '$lib/services/Universe';
	import { equal, key } from '$lib/types/MapObject';
	import ScannerMinefield from './ScannerMinefield.svelte';

	const { universe, selectedMapObject } = getGameContext();

	function getColor(minefield: AnyMinefield) {
		return $universe.getPlayerColor(minefield.mapObject?.playerNum);
	}

	let minefields = $derived($universe.allMinefields);
	let selectedMinefield = $derived(
		$selectedMapObject && $selectedMapObject.mapObject?.type === MapObjectType.MINEFIELD
			? ($selectedMapObject as AnyMinefield)
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
