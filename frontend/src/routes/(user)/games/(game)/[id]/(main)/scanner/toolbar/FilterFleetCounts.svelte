<script lang="ts">
	import FleetCount from '#lib/components/icons/FleetCount.svelte';
	import { getGameContext } from '#lib/services/GameContext.js';

	const { settings } = getGameContext();
</script>

<a
	href="#show-fleet-token-count"
	title="Show Fleet Counts"
	class:btn-primary={$settings.showFleetTokenCounts}
	class:btn-ghost={!$settings.showFleetTokenCounts}
	class="btn btn-ghost fill-current btn-xs size-full p-0"
	onclick={() => ($settings.showFleetTokenCounts = !$settings.showFleetTokenCounts)}
	><FleetCount class="size-6 " /></a
>
