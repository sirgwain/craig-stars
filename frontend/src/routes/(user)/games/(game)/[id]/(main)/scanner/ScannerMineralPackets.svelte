<!--
  @component
  Show all mineralpackets in the universe
 -->
<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import type { AnyMineralPacket } from '$lib/services/Universe';
	import ScannerMineralPacket from './ScannerMineralPacket.svelte';

	const { player, universe } = getGameContext();

	function getColor(mineralPacket: AnyMineralPacket) {
		if (mineralPacket.playerNum === $player.num) {
			return '#0900FF';
		}
		return $universe.getPlayerColor(mineralPacket.playerNum);
	}
</script>

<!-- MineralPackets -->
{#each $universe.allMineralPackets as mineralPacket}
	<ScannerMineralPacket {mineralPacket} color={getColor(mineralPacket)} />
{/each}
