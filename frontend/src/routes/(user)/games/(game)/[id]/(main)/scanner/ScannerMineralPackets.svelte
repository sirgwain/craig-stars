<!--
  @component
  Show all mineralpackets in the universe
 -->
<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import type { MineralPacket } from '$lib/types/MineralPacket';
	import ScannerMineralPacket from './ScannerMineralPacket.svelte';

	const { player, universe } = getGameContext();

	function getColor(mineralPacket: MineralPacket) {
		if (mineralPacket.playerNum === $player.num) {
			return '#0900FF';
		}
		return $universe.getPlayerColor(mineralPacket.playerNum);
	}
</script>

<!-- MineralPackets -->
{#each $universe.mineralPackets as mineralPacket}
	<ScannerMineralPacket {mineralPacket} color={getColor(mineralPacket)} />
{/each}
