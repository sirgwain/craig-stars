<script lang="ts">
	import type { Game } from '$lib/types/Game';
	import type { GameContext } from '$lib/types/GameContext';
	import { createEventDispatcher, getContext } from 'svelte';

	export let game = getContext<GameContext>('game').game;

	const updateTitle = () => (document.title = `${game.name} - ${game.year}`);
	$: game && updateTitle();
	const dispatch = createEventDispatcher();
</script>

<div class="flex py-5">
	<div class="flux-none">
		<div class="prose">
			<h2>{game.name} - {game.year}</h2>
		</div>
	</div>
	<div class="flux-1 grow">
		<button on:click={() => dispatch('submit-turn')} class="last:float-right btn btn-primary"
			>Submit Turn</button
		>
	</div>
</div>
