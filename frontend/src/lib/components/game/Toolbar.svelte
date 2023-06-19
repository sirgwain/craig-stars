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
		<a class="btn btn-ghost text-2xl text-primary" href={`/games/${game.id}`}
			>{game.name} - {game.year}</a
		>
	</div>
	<div class="flux-1 grow">
		<div class="dropdown dropdown-hover">
			<label for="reports" tabindex="0" class="btn">Reports</label>
		<ul
				id="reports"
				tabindex="0"
				class="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-52"
			>
				<li><a href={`/games/${game.id}/planets`}>Planets</a></li>
				<li><a href={`/games/${game.id}/fleets`}>Fleets</a></li>
				<li><a href={`/games/${game.id}/designs`}>Designs</a></li>
				<li><a href={`/games/${game.id}/messages`}>Messages</a></li>
			</ul>
		</div>

		<button on:click={() => dispatch('submit-turn')} class="last:float-right btn btn-primary"
			>Submit Turn</button
		>
	</div>
</div>
