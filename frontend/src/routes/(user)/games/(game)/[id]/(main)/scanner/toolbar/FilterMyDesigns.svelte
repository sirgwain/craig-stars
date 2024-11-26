<script lang="ts">
	import FilterFleets from '$lib/components/icons/FilterFleets.svelte';
	import { getGameContext } from '$lib/services/GameContext';

	const { settings } = getGameContext();
</script>

<a
	href="#show-idle-fleets-only"
	title="Design Filter"
	class:fill-accent={$settings.filterMyDesigns}
	class:fill-current={!$settings.filterMyDesigns}
	class="btn btn-ghost btn-xs w-full h-full"
	onclick={() => ($settings.filterMyDesigns = !$settings.filterMyDesigns)}
	><FilterFleets class="w-6 h-6" /></a
>
