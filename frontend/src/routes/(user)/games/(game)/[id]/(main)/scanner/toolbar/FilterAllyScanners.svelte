<script lang="ts">
	import Scanner from '$lib/components/icons/Scanner.svelte';
	import { getGameContext } from '$lib/services/GameContext';

	const { settings } = getGameContext();
</script>

<a
	href="#ally-scanner-toggle"
	title="Show Ally Scanners"
	class:fill-ally-selected={$settings.showAllyScanners}
	class:fill-ally-current={!$settings.showAllyScanners}
	class="btn btn-ghost btn-xs w-full h-full"
	onclick={() => ($settings.showAllyScanners = !$settings.showAllyScanners)}
	><span><Scanner class="w-6 h-6" /></span></a
>
