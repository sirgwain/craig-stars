<script lang="ts">
	import FleetCount from '$lib/components/icons/FleetCount.svelte';
	import IdleFleets from '$lib/components/icons/IdleFleets.svelte';
	import { getGameContext, playerFinderKey } from '$lib/services/GameContext';
	import Scanner from '$lib/components/icons/Scanner.svelte';

	const { player, settings } = getGameContext();
</script>

<a
	href="#show-fleet-token-count"
	title="Show Fleet Counts"
	class:fill-accent={$settings.showFleetTokenCounts}
	class:fill-current={!$settings.showFleetTokenCounts}
	class="btn btn-ghost btn-xs w-full h-full"
	on:click|preventDefault={() => ($settings.showFleetTokenCounts = !$settings.showFleetTokenCounts)}
	><FleetCount class="w-6 h-6" /></a
>
