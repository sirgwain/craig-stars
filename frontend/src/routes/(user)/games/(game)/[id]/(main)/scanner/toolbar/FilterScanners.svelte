<script lang="ts">
	import FleetCount from '$lib/components/icons/FleetCount.svelte';
	import IdleFleets from '$lib/components/icons/IdleFleets.svelte';
	import { getGameContext, playerFinderKey } from '$lib/services/GameContext';
	import Scanner from '$lib/components/icons/Scanner.svelte';

	const { player, settings } = getGameContext();
</script>

<a
	href="#scanner-toggle"
	title="Show Scanners"
	class:fill-accent={$settings.showScanners}
	class:fill-current={!$settings.showScanners}
	class="btn btn-ghost btn-xs w-full h-full"
	on:click|preventDefault={() => ($settings.showScanners = !$settings.showScanners)}
	><span><Scanner class="w-6 h-6" /></span></a
>
