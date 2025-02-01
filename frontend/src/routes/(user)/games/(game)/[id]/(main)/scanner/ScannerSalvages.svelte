<!--
  @component
  Show all salvages in the universe
 -->
<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import ScannerSalvage from './ScannerSalvage.svelte';

	const { universe } = getGameContext();
</script>

<!-- Salvages -->
{#each $universe.salvages as salvage}
	<ScannerSalvage {salvage} />
{/each}
