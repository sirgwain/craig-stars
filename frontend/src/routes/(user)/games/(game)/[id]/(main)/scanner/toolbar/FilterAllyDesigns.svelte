<script lang="ts">
	import FilterFleets from '$lib/components/icons/FilterFleets.svelte';
	import { getGameContext } from '$lib/services/GameContext';

	const { settings } = getGameContext();
</script>

<a
	href="#show-idle-fleets-only"
	title="Ally Ship Class Filter"
	class:btn-primary={$settings.filterAllyDesigns}
	class:btn-ghost={!$settings.filterAllyDesigns}
	class="btn btn-ghost fill-current btn-xs w-full h-full"
	onclick={() => ($settings.filterAllyDesigns = !$settings.filterAllyDesigns)}
	><FilterFleets primary="fill-fleet-ally" class="w-6 h-6" /></a
>
