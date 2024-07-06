<script lang="ts">
	import FleetCount from '$lib/components/icons/FleetCount.svelte';
	import IdleFleets from '$lib/components/icons/IdleFleets.svelte';
	import { getGameContext, playerFinderKey } from '$lib/services/GameContext';
	import Scanner from '$lib/components/icons/Scanner.svelte';

	const { player, settings } = getGameContext();
</script>

<ul {...$$restProps}>
	<li class="h-10 w-10">
		<a
			href="#scanner-toggle"
			title="Show Scanners"
			class:fill-accent={$settings.showScanners}
			class:fill-current={!$settings.showScanners}
			class="btn btn-ghost btn-xs w-full h-full"
			on:click|preventDefault={() => ($settings.showScanners = !$settings.showScanners)}
			><span><Scanner class="w-6 h-6" /></span></a
		>
	</li>

	{#if $player.relations.filter((r) => r.shareMap).length > 0}
		<!-- optionally turn off ally scanners if we are sharing maps -->
		<li class="h-10 w-10">
			<a
				href="#ally-scanner-toggle"
				title="Show Ally Scanners"
				class:fill-ally-selected={$settings.showAllyScanners}
				class:fill-ally-current={!$settings.showAllyScanners}
				class="btn btn-ghost btn-xs w-full h-full"
				on:click|preventDefault={() => ($settings.showAllyScanners = !$settings.showAllyScanners)}
				><span><Scanner class="w-6 h-6" /></span></a
			>
		</li>
	{/if}

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
