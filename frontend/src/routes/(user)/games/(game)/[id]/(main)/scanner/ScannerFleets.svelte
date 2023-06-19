<!--
  @component
  Generates an SVG scatter plot. This component can also work if the x- or y-scale is ordinal, i.e. it has a `.bandwidth` method. See the [timeplot chart](https://layercake.graphics/example/Timeplot) for an example.
 -->
<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import { commandedFleet } from '$lib/services/Stores';
	import type { Fleet } from '$lib/types/Fleet';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import ScannerFleet from './ScannerFleet.svelte';

	const { data } = getContext<LayerCake>('LayerCake');
	const { game, universe } = getGameContext();

	let fleets: Fleet[] = [];

	$: fleets = $universe.fleets.filter((f: Fleet) => !f.orbitingPlanetNum);
</script>

<!-- Fleets -->
{#each fleets as fleet}
	<ScannerFleet
		{fleet}
		color={$universe.getPlayerColor(fleet.playerNum)}
		commanded={$commandedFleet?.num === fleet.num && $commandedFleet?.playerNum === fleet.playerNum}
	/>
{/each}
