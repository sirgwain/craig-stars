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
	import FilterMyDesigns from './FilterMyDesigns.svelte';
	import FilterAllyDesigns from './FilterAllyDesigns.svelte';
	import FilterEnemyDesigns from './FilterEnemyDesigns.svelte';
	import { clamp } from 'lodash-es';

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
	<div class="flex flex-row justify-center w-full">
		<div class="h-10 w-10">
			<FilterFleetCounts />
		</div>
		<div class="h-10 w-10">
			<FilterIdleFleets />
		</div>
		<div class="h-10 w-10">
			<FilterMyDesigns />
		</div>
		{#if $player.relations.filter((r) => r.shareMap).length > 0}
			<div class="h-10 w-10">
				<FilterAllyDesigns />
			</div>
		{/if}
		<div class="h-10 w-10">
			<FilterEnemyDesigns />
		</div>
	</div>
	<div class="menu-title">Scanners</div>
	<div class="flex flex-row justify-center w-full">
		<div class="h-10 w-10">
			<FilterScanners />
		</div>

		{#if $player.relations.filter((r) => r.shareMap).length > 0}
			<!-- optionally turn off ally scanners if we are sharing maps -->
			<div class="h-10 w-10">
				<FilterAllyScanners />
			</div>
		{/if}

		<div class="px-1 my-auto">
			<input
				class="input input-sm input-bordered w-16 pr-0 pl-1"
				type="number"
				min={0}
				max={100}
				step={10}
				value={$settings.scannerPercent}
				on:change={(e) => {
					const val = parseInt(e.currentTarget.value);
					if (val) {
						$settings.scannerPercent = clamp(val, 0, 100);
						e.currentTarget.value = `${$settings.scannerPercent}`;
					} else {
						$settings.scannerPercent = 0;
						e.currentTarget.value = '0';
					}
				}}
			/>
		</div>
	</div>
</div>
