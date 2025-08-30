<script lang="ts">
	import FilterFleets from '$lib/components/icons/FilterFleets.svelte';
	import { getGameContext } from '$lib/services/GameContext';

	const { settings } = getGameContext();
</script>

<a
	href="#show-idle-fleets-only"
	title="Design Filter"
	class:btn-primary={$settings.filterMyDesigns}
	class:btn-ghost={!$settings.filterMyDesigns}
	class="btn btn-ghost fill-current btn-xs w-full h-full"
	onclick={() => ($settings.filterMyDesigns = !$settings.filterMyDesigns)}
	><FilterFleets class="w-6 h-6" /></a
>
