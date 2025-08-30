<!--
  @component
  Show all mineralpackets in the universe
 -->
<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import type { MineralPacket } from '$lib/types/cs-proto';
	import { key } from '$lib/types/MapObject';
	import { getDisplayColor } from '$lib/utils/colorUtils';
	import ScannerMineralPacket from './ScannerMineralPacket.svelte';

	const { player, universe, settings } = getGameContext();

	function getColor(mineralPacket: MineralPacket) {
		if (mineralPacket.mapObject?.playerNum === $player.num) {
			return '#0900FF';
		}
		return getDisplayColor(mineralPacket.mapObject?.playerNum, $player, $universe, $settings);
	}
</script>

<!-- MineralPackets -->
{#each $universe.mineralPackets as mineralPacket (key(mineralPacket))}
	<ScannerMineralPacket {mineralPacket} color={getColor(mineralPacket)} />
{/each}
