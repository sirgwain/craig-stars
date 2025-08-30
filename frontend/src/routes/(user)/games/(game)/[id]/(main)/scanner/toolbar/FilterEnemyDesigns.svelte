<script lang="ts">
	import FilterFleets from '$lib/components/icons/FilterFleets.svelte';
	import { getGameContext } from '$lib/services/GameContext';

	const { settings } = getGameContext();
</script>

<a
	href="#show-idle-fleets-only"
	title="Enemy Ship Class Filter"
	class:btn-primary={$settings.filterEnemyDesigns}
	class:btn-ghost={!$settings.filterEnemyDesigns}
	class="btn btn-ghost fill-current btn-xs w-full h-full"
	onclick={() => ($settings.filterEnemyDesigns = !$settings.filterEnemyDesigns)}
	><FilterFleets primary="fill-fleet-enemy" class="w-6 h-6" /></a
>
