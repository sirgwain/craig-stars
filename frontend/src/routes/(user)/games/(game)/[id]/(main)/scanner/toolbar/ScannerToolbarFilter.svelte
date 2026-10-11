<script lang="ts">
	import { getGameContext } from '#lib/services/GameContext.js';
	import type { HTMLAttributes } from 'svelte/elements';
	import FilterAllyDesigns from './FilterAllyDesigns.svelte';
	import FilterAllyScanners from './FilterAllyScanners.svelte';
	import FilterEnemyDesigns from './FilterEnemyDesigns.svelte';
	import FilterFleetCounts from './FilterFleetCounts.svelte';
	import FilterIdleFleets from './FilterIdleFleets.svelte';
	import FilterMinefields from './FilterMinefields.svelte';
	import FilterMyDesigns from './FilterMyDesigns.svelte';
	import FilterPlanetNames from './FilterPlanetNames.svelte';
	import FilterPlayerColors from './FilterPlayerColors.svelte';
	import FilterScanners from './FilterScanners.svelte';
	import ScannerPercent from './ScannerPercent.svelte';

	let { ...props }: HTMLAttributes<HTMLDivElement> = $props();

	const { player } = getGameContext();

	let shareMap = $derived($player.relations.some((r) => r.shareMap));
</script>

<!-- the filters are in the same order as the original's toolbar, with ours beside their closest relatives -->
<div {...props}>
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
	<div class="size-10">
		<FilterMinefields />
	</div>
	<div class="size-10">
		<FilterPlanetNames />
	</div>
	<div class="size-10">
		<FilterFleetCounts />
	</div>
	<div class="size-10">
		<FilterIdleFleets />
	</div>
	<FilterMyDesigns />
	{#if shareMap}
		<FilterAllyDesigns />
	{/if}
	<FilterEnemyDesigns />
	<div class="size-10">
		<FilterPlayerColors />
	</div>
</div>
