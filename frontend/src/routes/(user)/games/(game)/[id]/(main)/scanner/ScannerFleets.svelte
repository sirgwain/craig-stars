<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { filterFleet } from '$lib/types/Filter';
	import { type Fleet } from '$lib/types/cs';
	import { equal } from '$lib/types/MapObject';
	import ScannerFleet from './ScannerFleet.svelte';

	const { player, universe, commandedFleet, settings } = getGameContext();

	let fleets: Fleet[] = $derived(
		$universe.fleets
			.filter((f: Fleet) => !f.orbitingPlanetNum)
			.filter((f: Fleet) => equal($commandedFleet, f) || filterFleet($player, f, $settings))
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
