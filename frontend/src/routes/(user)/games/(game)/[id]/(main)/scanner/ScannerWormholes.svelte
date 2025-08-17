<!--
  @component
  Show all wormholes in the universe
 -->
<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import ScannerWormhole from './ScannerWormhole.svelte';

	const { universe } = getGameContext();
</script>

<!-- Wormholes -->
{#each $universe.wormholeIntels as wormhole (wormhole.mapObject?.num)}
	<ScannerWormhole {wormhole} />
{/each}
