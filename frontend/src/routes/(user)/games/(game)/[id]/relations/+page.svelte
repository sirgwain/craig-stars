<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import Relations from './Relations.svelte';

	const { game, player, universe } = getGameContext();

	async function update() {
		$game.updatePlayerRelations();
	}
</script>

<div class="w-full mx-auto md:max-w-2xl">
	{#if player}
		<Relations on:update-player={update} />
	{/if}
</div>
