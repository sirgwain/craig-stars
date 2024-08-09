<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import PlanetViewStateHab from './PlanetViewStateHab.svelte';
	import PlanetViewStateMineralConc from './PlanetViewStateMineralConc.svelte';
	import PlanetViewStatesNormal from './PlanetViewStateNormal.svelte';
	import PlanetViewStatePopulation from './PlanetViewStatePopulation.svelte';
	import PlanetViewStateSurfaceMinerals from './PlanetViewStateSurfaceMinerals.svelte';
	import FilterAllyScanners from './FilterAllyScanners.svelte';
	import FilterFleetCounts from './FilterFleetCounts.svelte';
	import FilterIdleFleets from './FilterIdleFleets.svelte';
	import FilterScanners from './FilterScanners.svelte';

	const { player, settings } = getGameContext();
</script>

<div class="flex flex-col bottom-0 top-0 bg-base-200">
    <div class="menu-title">Planet View</div>
	<div class="flex flex-row w-full">
		<div class="h-10 w-10">
			<PlanetViewStatesNormal />
		</div>
		<div class="h-10 w-10">
			<PlanetViewStateSurfaceMinerals />
		</div>
		<div class="h-10 w-10">
			<PlanetViewStateMineralConc />
		</div>
		<div class="h-10 w-10">
			<PlanetViewStateHab />
		</div>
		<div class="h-10 w-10">
			<PlanetViewStatePopulation />
		</div>
	</div>
    <div class="menu-title">Filter</div>
	<div class="flex flex-row w-full">
		<div class="h-10 w-10">
			<FilterScanners />
		</div>

		{#if $player.relations.filter((r) => r.shareMap).length > 0}
			<!-- optionally turn off ally scanners if we are sharing maps -->
			<div class="h-10 w-10">
				<FilterAllyScanners />
			</div>
		{/if}

		<div class="h-10 w-10">
			<FilterFleetCounts />
		</div>
		<div class="h-10 w-10">
			<FilterIdleFleets />
		</div>
	</div>
</div>
