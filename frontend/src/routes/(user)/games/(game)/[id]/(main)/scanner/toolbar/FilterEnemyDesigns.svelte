<script lang="ts">
	import FilterFleets from '$lib/components/icons/FilterFleets.svelte';
	import { getGameContext } from '$lib/services/GameContext';

	const { player, settings } = getGameContext();
</script>

<a
	href="#show-idle-fleets-only"
	title="Enemy Ship Class Filter"
	class:fill-accent={$settings.filterEnemyDesigns}
	class:fill-current={!$settings.filterEnemyDesigns}
	class="btn btn-ghost btn-xs w-full h-full"
	on:click|preventDefault={() => ($settings.filterEnemyDesigns = !$settings.filterEnemyDesigns)}
	><FilterFleets primary="fill-fleet-enemy" class="w-6 h-6" /></a
>
