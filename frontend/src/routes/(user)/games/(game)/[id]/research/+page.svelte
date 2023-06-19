<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import Research from './Research.svelte';

	const { game, player, universe } = getGameContext();

	async function update() {
		$game.updatePlayerOrders();
	}
</script>

<div class="w-full mx-auto md:max-w-2xl">
	{#if player}
		<Research player={$player} on:update-player={update} />
	{/if}
</div>
