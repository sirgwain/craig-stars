<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import type { Fleet } from '$lib/types/cs-proto';
	import { filterFleet } from '$lib/types/Filter';
	import { equal } from '$lib/types/MapObject';
	import ScannerFleet from './ScannerFleet.svelte';

	const { player, universe, commandedFleet, settings } = getGameContext();

	let fleets: Fleet[] = $derived(
		$universe
			.getAllFleets()
			.filter((f) => !f.orbitingPlanetNum)
			.filter((f) => equal($commandedFleet, f) || filterFleet($player, f, $settings))
	);
</script>

<!-- Fleets -->
{#each fleets as fleet (fleet)}
	<ScannerFleet
		{fleet}
		color={$universe.getPlayerColor(fleet.mapObject?.playerNum)}
		commanded={$commandedFleet?.mapObject.num === fleet.mapObject?.num &&
			$commandedFleet?.mapObject.playerNum === fleet.mapObject?.playerNum}
	/>
{/each}
