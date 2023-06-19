<!--
  @component
  Show all mineralpackets in the universe
 -->
<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import { selectedMapObject } from '$lib/services/Stores';
	import { MapObjectType } from '$lib/types/MapObject';
	import type { MineralPacket } from '$lib/types/MineralPacket';
	import ScannerMineralPacket from './ScannerMineralPacket.svelte';

	const { player, universe } = getGameContext();

	function getColor(mineralPacket: MineralPacket) {
		if (mineralPacket.playerNum === $player.num) {
			return '#0900FF';
		}
		return $universe.getPlayerColor(mineralPacket.playerNum);
	}

	$: mineralpackets = $universe.mineralPackets;

	$: selectedMineralPacket =
		$selectedMapObject && $selectedMapObject.type === MapObjectType.MineralPacket
			? ($selectedMapObject as MineralPacket)
			: undefined;
</script>

<!-- MineralPackets -->
{#each mineralpackets as mineralPacket}
	<ScannerMineralPacket {mineralPacket} color={getColor(mineralPacket)} />
{/each}
