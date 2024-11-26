<script lang="ts">
	import Scanner from '$lib/components/icons/Scanner.svelte';
	import { getGameContext } from '$lib/services/GameContext';

	const { player, settings } = getGameContext();
</script>

<a
	href="#scanner-toggle"
	title="Show Scanners"
	class:fill-accent={$settings.showScanners}
	class:fill-current={!$settings.showScanners}
	class="btn btn-ghost btn-xs w-full h-full"
	onclick={() => ($settings.showScanners = !$settings.showScanners)}
	><span><Scanner class="w-6 h-6" /></span></a
>
