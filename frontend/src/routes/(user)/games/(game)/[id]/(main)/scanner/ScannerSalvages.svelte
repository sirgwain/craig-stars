<!--
  @component
  Show all salvages in the universe
 -->
<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import type { Salvage } from '$lib/types/Salvage';
	import ScannerSalvage from './ScannerSalvage.svelte';

	const { player, universe } = getGameContext();
</script>

<!-- Salvages -->
{#each $universe.salvages as salvage}
	<ScannerSalvage {salvage} />
{/each}
