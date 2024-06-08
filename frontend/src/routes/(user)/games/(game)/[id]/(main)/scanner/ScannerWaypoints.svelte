<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { type Fleet, idleFleetsFilter } from '$lib/types/Fleet';
	import { equal } from '$lib/types/MapObject';
	import ScannerWaypointLine from './ScannerWaypointLine.svelte';

	const { universe, settings, commandedFleet, selectedWaypoint } = getGameContext();

	$: fleets = $universe.fleets.filter(
		(f: Fleet) => equal(f, $commandedFleet) || idleFleetsFilter(f, $settings.showIdleFleetsOnly)
	);
</script>

{#each fleets as fleet}
	{#if fleet.waypoints && fleet.waypoints.length > 1}
		<ScannerWaypointLine
			{fleet}
			selectedWaypoint={$selectedWaypoint}
			commanded={fleet.num === $commandedFleet?.num}
		/>
	{/if}
{/each}
