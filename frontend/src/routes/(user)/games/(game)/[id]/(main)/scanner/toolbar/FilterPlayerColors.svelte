<script lang="ts">
	import PlayerColors from '#lib/components/icons/PlayerColors.svelte';
	import { getGameContext } from '#lib/services/GameContext.js';

	const { settings } = getGameContext();
</script>

<a
	href="#show-player-colors"
	title="Show Player Colors"
	class:btn-primary={$settings.showPlayerColors}
	class:btn-ghost={!$settings.showPlayerColors}
	class="btn btn-ghost btn-xs size-full p-0"
	onclick={() => ($settings.showPlayerColors = !$settings.showPlayerColors)}
	><PlayerColors class="size-6 " /></a
>
