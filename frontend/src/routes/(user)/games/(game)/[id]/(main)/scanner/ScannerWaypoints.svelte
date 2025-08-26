<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { filterFleet } from '$lib/types/Filter';
	import type { Fleet } from '$lib/types/cs-proto';
	import { equal } from '$lib/types/MapObject';
	import ScannerWaypointLine from './ScannerWaypointLine.svelte';

	const { player, universe, settings, commandedFleet, selectedWaypoint } = getGameContext();

	let fleets = $derived(
		$universe.fleets.filter(
			(f: Fleet) => equal($commandedFleet, f) || filterFleet($player, f, $settings)
		)
	);
</script>

{#each fleets as fleet (fleet.mapObject?.num)}
	{#if fleet.fleetOrders?.waypoints && fleet.fleetOrders.waypoints.length > 1 && fleet.mapObject?.num !== $commandedFleet?.mapObject.num}
		<ScannerWaypointLine {fleet} selectedWaypoint={$selectedWaypoint} />
	{/if}
{/each}
{#if $commandedFleet && $commandedFleet.fleetOrders.waypoints.length > 1}
	<ScannerWaypointLine
		fleet={$commandedFleet}
		selectedWaypoint={$selectedWaypoint}
		commanded={true}
	/>
{/if}
