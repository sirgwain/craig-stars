<script lang="ts">
	import { commandedFleet, selectedWaypoint } from '$lib/services/Stores';

	import { getGameContext } from '$lib/services/Contexts';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';
	import ScannerWaypointLine from './ScannerWaypointLine.svelte';

	const { game, player, universe } = getGameContext();
	const scale = getContext<Writable<number>>('scale');
	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');
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
