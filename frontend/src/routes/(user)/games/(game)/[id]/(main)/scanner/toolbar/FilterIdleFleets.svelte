<script lang="ts">
	import FleetCount from '$lib/components/icons/FleetCount.svelte';
	import IdleFleets from '$lib/components/icons/IdleFleets.svelte';
	import { getGameContext, playerFinderKey } from '$lib/services/GameContext';
	import Scanner from '$lib/components/icons/Scanner.svelte';

	const { player, settings } = getGameContext();
</script>

<a
	href="#show-idle-fleets-only"
	title="Idle Fleets Filter"
	class:fill-accent={$settings.showIdleFleetsOnly}
	class:fill-current={!$settings.showIdleFleetsOnly}
	class="btn btn-ghost btn-xs w-full h-full"
	on:click|preventDefault={() => ($settings.showIdleFleetsOnly = !$settings.showIdleFleetsOnly)}
	><IdleFleets class="w-6 h-6" /></a
>
