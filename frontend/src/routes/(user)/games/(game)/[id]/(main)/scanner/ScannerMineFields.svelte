<!--
  @component
  Show all minefields in the universe
 -->
<script lang="ts">
	import type { FullGame } from '$lib/services/FullGame';
	import { MapObjectType, type MapObject, equal } from '$lib/types/MapObject';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import ScannerMineField from './ScannerMineField.svelte';
	import { selectedMapObject } from '$lib/services/Stores';
	import type { MineField } from '$lib/types/MineField';
	import { min } from 'date-fns';

	const game = getContext<FullGame>('game');
	const { data } = getContext<LayerCake>('LayerCake');

	function getColor(mineField: MineField) {
		if (mineField.playerNum === game.player.num) {
			return '#0900FF';
		}
		return game.getPlayerColor(mineField.playerNum);
	}

	$: minefields = $data && $data.filter((mo: MapObject) => mo.type == MapObjectType.MineField);
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
