<!--
  @component
  Show all mineralpackets in the universe
 -->
<script lang="ts">
	import type { FullGame } from '$lib/services/FullGame';
	import { MapObjectType, type MapObject, equal } from '$lib/types/MapObject';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import ScannerMineralPacket from './ScannerMineralPacket.svelte';
	import { selectedMapObject } from '$lib/services/Context';
	import type { MineralPacket } from '$lib/types/MineralPacket';
	import { min } from 'date-fns';

	const game = getContext<FullGame>('game');
	const { data } = getContext<LayerCake>('LayerCake');

	function getColor(mineralPacket: MineralPacket) {
		if (mineralPacket.playerNum === game.player.num) {
			return '#0900FF';
		}
		return game.getPlayerColor(mineralPacket.playerNum);
	}

	$: mineralpackets =
		$data && $data.filter((mo: MapObject) => mo.type == MapObjectType.MineralPacket);
	$: selectedMineralPacket =
		$selectedMapObject && $selectedMapObject.type === MapObjectType.MineralPacket
			? ($selectedMapObject as MineralPacket)
			: undefined;
</script>

<!-- MineralPackets -->
{#each mineralpackets as mineralPacket}
	{#if mineralPacket !== selectedMineralPacket}
		<ScannerMineralPacket
			{mineralPacket}
			color={getColor(mineralPacket)}
			selected={equal($selectedMapObject, mineralPacket)}
		/>
	{/if}
{/each}
{#if selectedMineralPacket}
	<ScannerMineralPacket
		mineralPacket={selectedMineralPacket}
		color={getColor(selectedMineralPacket)}
		selected={true}
	/>
{/if}
