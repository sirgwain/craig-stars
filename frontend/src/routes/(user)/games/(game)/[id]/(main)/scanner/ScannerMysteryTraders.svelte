<!--
  @component
  Generates an SVG scatter plot. This component can also work if the x- or y-scale is ordinal, i.e. it has a `.bandwidth` method. See the [timeplot chart](https://layercake.graphics/example/Timeplot) for an example.
 -->
<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import ScannerMysteryTrader from './ScannerMysteryTrader.svelte';

	const { universe } = getGameContext();
</script>

<!-- MysteryTraders -->
{#each $universe.mysteryTraders as mysteryTrader}
	<ScannerMysteryTrader mysteryTrader={mysteryTrader} />
{/each}
