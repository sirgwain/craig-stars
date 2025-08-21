<!--
  @component
  Show all salvages in the universe
 -->
<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { key } from '$lib/types/MapObject';
	import ScannerSalvage from './ScannerSalvage.svelte';

	const { universe } = getGameContext();
</script>

<!-- Salvages -->
{#each $universe.salvages as salvage (key(salvage))}
	<ScannerSalvage {salvage} />
{/each}
