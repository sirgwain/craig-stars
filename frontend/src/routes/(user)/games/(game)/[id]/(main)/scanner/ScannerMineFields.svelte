<!--
  @component
  Show all minefields in the universe
 -->
<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import { selectedMapObject } from '$lib/services/Stores';
	import { MapObjectType, equal } from '$lib/types/MapObject';
	import type { MineField } from '$lib/types/MineField';
	import ScannerMineField from './ScannerMineField.svelte';

	const { player, universe } = getGameContext();

	function getColor(mineField: MineField) {
		if (mineField.playerNum === $player.num) {
			return '#0900FF';
		}
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
