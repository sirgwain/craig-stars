<script lang="ts">
	import type { Game } from '$lib/types/Game';
	import type { GameContext } from '$lib/types/GameContext';
	import { createEventDispatcher, getContext } from 'svelte';

	export let game = getContext<GameContext>('game').game;

	const dispatch = createEventDispatcher();
</script>

<div class="flex py-5">
	<div class="flux-none">
		<h1 class="text-3xl mb-2">{game.name} - {game.year}</h1>
	</div>
	<div class="flux-1 grow">
		<button on:click={() => dispatch('submit-turn')} class="last:float-right button"
			>Submit Turn</button
		>
	</div>
</div>
