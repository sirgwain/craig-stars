<script lang="ts">
	import PlayerColors from '$lib/components/icons/PlayerColors.svelte';
	import { getGameContext } from '$lib/services/GameContext';

	const { settings } = getGameContext();
</script>

<a
	href="#show-player-colors"
	title="Show Player Colors"
	class:btn-primary={$settings.showPlayerColors}
	class:btn-ghost={!$settings.showPlayerColors}
	class="btn btn-ghost btn-xs w-full h-full"
	onclick={() => ($settings.showPlayerColors = !$settings.showPlayerColors)}
	><PlayerColors class="w-6 h-6" /></a
>
