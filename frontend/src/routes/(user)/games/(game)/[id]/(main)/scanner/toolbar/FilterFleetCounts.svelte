<script lang="ts">
	import FleetCount from '$lib/components/icons/FleetCount.svelte';
	import { getGameContext } from '$lib/services/GameContext';

	const { settings } = getGameContext();
</script>

<a
	href="#show-fleet-token-count"
	title="Show Fleet Counts"
	class:fill-accent={$settings.showFleetTokenCounts}
	class:fill-current={!$settings.showFleetTokenCounts}
	class="btn btn-ghost btn-xs w-full h-full"
	onclick={() => ($settings.showFleetTokenCounts = !$settings.showFleetTokenCounts)}
	><FleetCount class="w-6 h-6" /></a
>
