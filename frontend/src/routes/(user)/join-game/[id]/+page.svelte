<script lang="ts">
	import { page } from '$app/stores';
	import { GameService } from '$lib/services/GameService';
	import type { Game, NewGamePlayer } from '$lib/types/Game';

	import type { PlayerService } from '$lib/services/PlayerService';
	import { onMount } from 'svelte';
	import PlayerChooser from './PlayerChooser.svelte';
	import { goto } from '$app/navigation';

	let id = parseInt($page.params.id);

	let gameService: GameService = new GameService();
	let game: Game;
	let raceId: number;

	onMount(async () => {
		game = await gameService.loadOpenGame(id);
	});

	const onSubmit = async () => {
		const data = JSON.stringify({ raceId });

		const response = await fetch(`/api/games/open/${id}`, {
			method: 'post',
			headers: {
				accept: 'application/json'
			},
			body: data
		});

		if (response.ok) {
			goto(`/games/${id}`);
		} else {
			const resolvedResponse = await response?.json();
			error = resolvedResponse.error;
			console.error(error);
		}
	};

	let error = '';
</script>

<h2 class="font-semibold text-xl">Join Game</h2>
<div class="text text-error">{error}</div>

{#if game}
	<div class="">
		<div class="flex">
			<div class="font-semibold w-[8rem]">Name:</div>
			<div class="text-left ml-2">{game.name}</div>
		</div>
		<div class="flex">
			<div class="font-semibold w-[8rem]">Size:</div>
			<div class="text-left ml-2">{game.size}</div>
		</div>
		<div class="flex">
			<div class="font-semibold w-[8rem]">Density:</div>
			<div class="text-left ml-2">{game.density}</div>
		</div>
		<div class="flex">
			<div class="font-semibold w-[8rem]">Player Distance:</div>
			<div class="text-left ml-2">{game.playerPositions}</div>
		</div>
	</div>

	<form on:submit|preventDefault={onSubmit}>
		<fieldset name="players" class="form-control mt-3">
			<PlayerChooser bind:raceId />
		</fieldset>
		<button class="btn btn-primary">Join</button>
	</form>
{/if}
