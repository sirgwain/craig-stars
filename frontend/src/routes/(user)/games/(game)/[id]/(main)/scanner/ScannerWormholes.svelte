<!--
  @component
  Show all wormholes in the universe
 -->
<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import ScannerWormhole from './ScannerWormhole.svelte';

	const { universe } = getGameContext();

	let wormholes = $derived($universe.wormholes);

</script>

<!-- Wormholes -->
{#each wormholes as wormhole}
	<ScannerWormhole {wormhole} />
{/each}
