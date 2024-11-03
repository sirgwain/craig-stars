<script lang="ts">
	import { preventDefault } from 'svelte/legacy';

	import FilterFleets from '$lib/components/icons/FilterFleets.svelte';
	import { getGameContext } from '$lib/services/GameContext';

	const { player, settings } = getGameContext();
</script>

<a
	href="#show-idle-fleets-only"
	title="Ally Ship Class Filter"
	class:fill-accent={$settings.filterAllyDesigns}
	class:fill-current={!$settings.filterAllyDesigns}
	class="btn btn-ghost btn-xs w-full h-full"
	onclick={preventDefault(() => ($settings.filterAllyDesigns = !$settings.filterAllyDesigns))}
	><FilterFleets primary="fill-fleet-ally" class="w-6 h-6" /></a
>
