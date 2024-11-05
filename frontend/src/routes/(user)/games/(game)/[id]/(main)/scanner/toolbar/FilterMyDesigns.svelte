<script lang="ts">
	import { preventDefault } from 'svelte/legacy';

	import FilterFleets from '$lib/components/icons/FilterFleets.svelte';
	import { getGameContext } from '$lib/services/GameContext';

	const { player, settings } = getGameContext();
</script>

<a
	href="#show-idle-fleets-only"
	title="Design Filter"
	class:fill-accent={$settings.filterMyDesigns}
	class:fill-current={!$settings.filterMyDesigns}
	class="btn btn-ghost btn-xs w-full h-full"
	onclick={preventDefault(() => ($settings.filterMyDesigns = !$settings.filterMyDesigns))}
	><FilterFleets class="w-6 h-6" /></a
>
