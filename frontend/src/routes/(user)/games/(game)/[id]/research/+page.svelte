<script lang="ts">
	import { game } from '$lib/services/Context';
	import Research from './Research.svelte';

	// make sure player updates are reactive
	$: player = $game?.player;

	async function update() {
		player = await $game?.updatePlayer();
	}
</script>

<div class="w-full mx-auto md:max-w-2xl">
	{#if player}
		<Research {player} on:update-player={update} />
	{/if}
</div>
