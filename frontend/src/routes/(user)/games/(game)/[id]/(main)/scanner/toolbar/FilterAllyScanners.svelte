<script lang="ts">
	import FleetCount from '$lib/components/icons/FleetCount.svelte';
	import IdleFleets from '$lib/components/icons/IdleFleets.svelte';
	import { getGameContext, playerFinderKey } from '$lib/services/GameContext';
	import Scanner from '$lib/components/icons/Scanner.svelte';

	const { player, settings } = getGameContext();
</script>

<a
	href="#ally-scanner-toggle"
	title="Show Ally Scanners"
	class:fill-ally-selected={$settings.showAllyScanners}
	class:fill-ally-current={!$settings.showAllyScanners}
	class="btn btn-ghost btn-xs w-full h-full"
	on:click|preventDefault={() => ($settings.showAllyScanners = !$settings.showAllyScanners)}
	><span><Scanner class="w-6 h-6" /></span></a
>
