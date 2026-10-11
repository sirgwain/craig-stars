<script lang="ts">
	import Scanner from '#lib/components/icons/Scanner.svelte';
	import { getGameContext } from '#lib/services/GameContext.js';

	const { settings } = getGameContext();
</script>

<a
	href="#scanner-toggle"
	title="Show Scanners"
	class:fill-accent={$settings.showScanners}
	class:fill-current={!$settings.showScanners}
	class="btn btn-ghost btn-xs size-full p-0"
	onclick={() => ($settings.showScanners = !$settings.showScanners)}
	><span><Scanner class="size-6 " /></span></a
>
