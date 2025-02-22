<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { filterFleet } from '$lib/types/Filter';
	import type { AnyFleet } from '$lib/services/Universe';
	import { equal } from '$lib/types/MapObject';
	import ScannerFleet from './ScannerFleet.svelte';

	const { player, universe, commandedFleet, settings } = getGameContext();

	let fleets: AnyFleet[] = $derived(
		$universe
			.getAllFleets()
			.filter((f) => !f.orbitingPlanetNum)
			.filter((f) => equal($commandedFleet, f) || filterFleet($player, f, $settings))
	);
</script>

<!-- Fleets -->
{#each fleets as fleet}
	<ScannerFleet
		{fleet}
		color={$universe.getPlayerColor(fleet.playerNum)}
		commanded={$commandedFleet?.num === fleet.num && $commandedFleet?.playerNum === fleet.playerNum}
	/>
{/each}
