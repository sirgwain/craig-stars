<script lang="ts">
	import { getGameContext } from '#lib/services/GameContext.js';
	import PlanetViewStateHab from './PlanetViewStateHab.svelte';
	import PlanetViewStateMineralConc from './PlanetViewStateMineralConc.svelte';
	import PlanetViewStatesNormal from './PlanetViewStateNormal.svelte';
	import PlanetViewStatePopulation from './PlanetViewStatePopulation.svelte';
	import PlanetViewStateSurfaceMinerals from './PlanetViewStateSurfaceMinerals.svelte';
	import FilterAllyScanners from './FilterAllyScanners.svelte';
	import FilterFleetCounts from './FilterFleetCounts.svelte';
	import FilterIdleFleets from './FilterIdleFleets.svelte';
	import FilterMinefields from './FilterMinefields.svelte';
	import FilterScanners from './FilterScanners.svelte';
	import FilterMyDesigns from './FilterMyDesigns.svelte';
	import FilterAllyDesigns from './FilterAllyDesigns.svelte';
	import FilterEnemyDesigns from './FilterEnemyDesigns.svelte';
	import FilterNotes from './FilterNotes.svelte';
	import FilterPlanetNames from './FilterPlanetNames.svelte';
	import FilterPlayerColors from './FilterPlayerColors.svelte';
	import ScannerPercent from './ScannerPercent.svelte';

	const { player } = getGameContext();

	let shareMap = $derived($player.relations.some((r) => r.shareMap));

	type Props = {
		// true if the toolbar is already showing the planet views
		hideViews?: boolean;
	};
	let { hideViews = false }: Props = $props();
</script>

<!-- everything is in the same order as the toolbar, in rows narrow enough for the menu -->
<div class="flex flex-col bg-base-200">
	{#if !hideViews}
		<div class="menu-title">Planet View</div>
		<div class="flex flex-row justify-center w-full">
			<div class="size-10">
				<PlanetViewStatesNormal />
			</div>
			<div class="size-10">
				<PlanetViewStateSurfaceMinerals />
			</div>
			<div class="size-10">
				<PlanetViewStateMineralConc />
			</div>
			<div class="size-10">
				<PlanetViewStateHab />
			</div>
			<div class="size-10">
				<PlanetViewStatePopulation />
			</div>
		</div>
	{/if}
	<div class="menu-title">Scanners</div>
	<div class="flex flex-row justify-center items-center w-full">
		<div class="size-10">
			<FilterScanners />
		</div>
		{#if shareMap}
			<!-- optionally turn off ally scanners if we are sharing maps -->
			<div class="size-10">
				<FilterAllyScanners />
			</div>
		{/if}
		<div class="px-1">
			<ScannerPercent />
		</div>
	</div>
	<div class="menu-title">Filter</div>
	<div class="flex flex-row justify-center w-full">
		<div class="size-10">
			<FilterMinefields />
		</div>
		<div class="size-10">
			<FilterPlanetNames />
		</div>
		<div class="size-10">
			<FilterNotes />
		</div>
		<div class="size-10">
			<FilterFleetCounts />
		</div>
		<div class="size-10">
			<FilterIdleFleets />
		</div>
	</div>
	<div class="flex flex-row justify-center w-full">
		<FilterMyDesigns />
		{#if shareMap}
			<FilterAllyDesigns />
		{/if}
		<FilterEnemyDesigns />
		<div class="size-10">
			<FilterPlayerColors />
		</div>
	</div>
</div>
