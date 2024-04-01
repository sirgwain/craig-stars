<!--
  @component
  Show all minefields in the universe
 -->
<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { MapObjectType, equal } from '$lib/types/MapObject';
	import type { MineField } from '$lib/types/MineField';
	import ScannerMineField from './ScannerMineField.svelte';

	const { universe, selectedMapObject } = getGameContext();

	function getColor(mineField: MineField) {
		return $universe.getPlayerColor(mineField.playerNum);
	}

	$: minefields = $universe.mineFields;
	$: selectedMineField =
		$selectedMapObject && $selectedMapObject.type === MapObjectType.MineField
			? ($selectedMapObject as MineField)
			: undefined;
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
