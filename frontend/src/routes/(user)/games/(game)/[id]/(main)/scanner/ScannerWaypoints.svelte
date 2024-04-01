<script lang="ts">

	import { getGameContext } from '$lib/services/GameContext';
	import ScannerWaypointLine from './ScannerWaypointLine.svelte';

	const { universe, commandedFleet, selectedWaypoint } = getGameContext();
</script>

{#each $universe.fleets as fleet}
	{#if fleet.waypoints && fleet.waypoints.length > 1}
		<ScannerWaypointLine
			{fleet}
			selectedWaypoint={$selectedWaypoint}
			commanded={fleet.num === $commandedFleet?.num}
		/>
	{/if}
{/each}
