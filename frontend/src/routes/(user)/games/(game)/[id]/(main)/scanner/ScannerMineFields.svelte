<!--
  @component
  Show all minefields in the universe
 -->
<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import type { AnyMineField } from '$lib/services/Universe';
	import { equal } from '$lib/types/MapObject';
	import { MapObjectTypeMineField } from '$lib/types/cs';
	import ScannerMineField from './ScannerMineField.svelte';

	const { universe, selectedMapObject } = getGameContext();

	function getColor(mineField: AnyMineField) {
		return $universe.getPlayerColor(mineField.playerNum);
	}

	let minefields = $derived($universe.allMineFields);
	let selectedMineField = $derived(
		$selectedMapObject && $selectedMapObject.type === MapObjectTypeMineField
			? ($selectedMapObject as AnyMineField)
			: undefined
	);
</script>

<!-- MineFields -->
{#each minefields as mineField}
	{#if mineField !== selectedMineField}
		<ScannerMineField
			{mineField}
			color={getColor(mineField)}
			selected={equal($selectedMapObject, mineField)}
		/>
	{/if}
{/each}
{#if selectedMineField}
	<ScannerMineField
		mineField={selectedMineField}
		color={getColor(selectedMineField)}
		selected={true}
	/>
{/if}
