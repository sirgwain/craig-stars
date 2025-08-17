<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import ScannerMysteryTrader from './ScannerMysteryTrader.svelte';

	const { universe } = getGameContext();
</script>

<!-- MysteryTraders -->
{#each $universe.mysteryTraderIntels as mysteryTrader (mysteryTrader.mapObject?.num)}
	<ScannerMysteryTrader {mysteryTrader} />
{/each}
