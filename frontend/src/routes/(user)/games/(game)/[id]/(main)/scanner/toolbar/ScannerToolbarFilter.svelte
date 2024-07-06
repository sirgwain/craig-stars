<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import FilterAllyScanners from './FilterAllyScanners.svelte';
	import FilterFleetCounts from './FilterFleetCounts.svelte';
	import FilterIdleFleets from './FilterIdleFleets.svelte';
	import FilterScanners from './FilterScanners.svelte';

	const { player, settings } = getGameContext();
</script>

<ul {...$$restProps}>
	<li class="h-10 w-10">
		<FilterScanners />
	</li>

	{#if $player.relations.filter((r) => r.shareMap).length > 0}
		<!-- optionally turn off ally scanners if we are sharing maps -->
		<li class="h-10 w-10">
			<FilterAllyScanners />
		</li>
	{/if}

	<li class="h-10 w-10">
		<FilterFleetCounts />
	</li>
	<li class="h-10 w-10">
		<FilterIdleFleets />
	</li>
</ul>
