<!--
  @component
  Generates an SVG scatter plot. This component can also work if the x- or y-scale is ordinal, i.e. it has a `.bandwidth` method. See the [timeplot chart](https://layercake.graphics/example/Timeplot) for an example.
 -->
<script lang="ts">
	import { commandedFleet } from '$lib/services/Context';
	import type { FullGame } from '$lib/services/FullGame';
	import type { Fleet } from '$lib/types/Fleet';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import ScannerFleet from './ScannerFleet.svelte';

	const game = getContext<FullGame>('game');
	const { data } = getContext<LayerCake>('LayerCake');

	let fleets: Fleet[] = [];

	$: fleets = [...game.universe.fleets, ...game.universe.fleetIntels].filter(
		(f: Fleet) => !f.orbitingPlanetNum
	);
</script>

<!-- Fleets -->
{#each fleets as fleet}
	<ScannerFleet
		{fleet}
		color={game.getPlayerColor(fleet.playerNum)}
		commanded={$commandedFleet?.num === fleet.num && $commandedFleet?.playerNum === fleet.playerNum}
	/>
{/each}
