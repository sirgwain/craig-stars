<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import Relations from './Relations.svelte';

	const { player, updatePlayerRelations } = getGameContext();
</script>

<div class="w-full mx-auto md:max-w-2xl">
	{#if player}
		<Relations on:update-player={updatePlayerRelations} />
	{/if}
</div>
