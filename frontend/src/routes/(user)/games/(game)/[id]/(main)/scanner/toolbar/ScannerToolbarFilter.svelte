<script lang="ts">
	import FleetCount from '$lib/components/icons/FleetCount.svelte';
	import IdleFleets from '$lib/components/icons/IdleFleets.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import Scanner from '$lib/components/icons/Scanner.svelte';

	const { settings } = getGameContext();
</script>

<ul {...$$restProps}>
	<li class="h-10 w-10">
		<a
			href="#scanner-toggle"
			class:fill-accent={$settings.showScanners}
			class:fill-current={!$settings.showScanners}
			class="btn btn-ghost btn-xs w-full h-full"
			on:click|preventDefault={() => ($settings.showScanners = !$settings.showScanners)}
			><span><Scanner class="w-6 h-6" /></span></a
		>
	</li>
	<li class="h-10 w-10">
		<a
			href="#show-fleet-token-count"
			title="Show Fleet Counts"
			class:fill-accent={$settings.showFleetTokenCounts}
			class:fill-current={!$settings.showFleetTokenCounts}
			class="btn btn-ghost btn-xs w-full h-full"
			on:click|preventDefault={() =>
				($settings.showFleetTokenCounts = !$settings.showFleetTokenCounts)}
			><FleetCount class="w-6 h-6" /></a
		>
	</li>
	<li class="h-10 w-10">
		<a
			href="#show-idle-fleets-only"
			title="Idle Fleets Filter"
			class:fill-accent={$settings.showIdleFleetsOnly}
			class:fill-current={!$settings.showIdleFleetsOnly}
			class="btn btn-ghost btn-xs w-full h-full"
			on:click|preventDefault={() => ($settings.showIdleFleetsOnly = !$settings.showIdleFleetsOnly)}
			><IdleFleets class="w-6 h-6" /></a
		>
	</li>
</ul>
