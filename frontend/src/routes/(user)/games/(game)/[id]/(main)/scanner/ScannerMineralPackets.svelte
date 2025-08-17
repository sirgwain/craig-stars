<!--
  @component
  Show all mineralpackets in the universe
 -->
<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import type { AnyMineralPacket } from '$lib/services/Universe';
	import { key } from '$lib/types/MapObject';
	import ScannerMineralPacket from './ScannerMineralPacket.svelte';

	const { player, universe } = getGameContext();

	function getColor(mineralPacket: AnyMineralPacket) {
		if (mineralPacket.mapObject?.playerNum === $player.num) {
			return '#0900FF';
		}
		return $universe.getPlayerColor(mineralPacket.mapObject?.playerNum);
	}
</script>

<!-- MineralPackets -->
{#each $universe.allMineralPackets as mineralPacket (key(mineralPacket))}
	<ScannerMineralPacket {mineralPacket} color={getColor(mineralPacket)} />
{/each}
