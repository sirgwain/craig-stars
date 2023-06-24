<script lang="ts">
	import { page } from '$app/stores';
	import { GameService } from '$lib/services/GameService';
	import type { Game } from '$lib/types/Game';

	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import PlayerChooser from '../../../../lib/components/game/newgame/PlayerChooser.svelte';
	import GameCard from '$lib/components/game/GameCard.svelte';
	import ColorInput from '$lib/components/ColorInput.svelte';
	import ItemTitle from '$lib/components/ItemTitle.svelte';

	let game: Game | undefined;
	let raceId: number;
	let color: string = '#0000FF';

	onMount(async () => {
		try {
			let id = parseInt($page.params.id);
			game = await GameService.loadGame(id);
		} catch {
			const games = await GameService.loadGameByHash($page.params.id);
			if (games.length == 1) {
				game = games[0];
			} else {
				error = 'No open game found for the invite';
			}
		}
	});

	const onSubmit = async () => {
		if (game) {
			const data = JSON.stringify({ raceId, color });

			const response = await fetch(`/api/games/${game.id}/join`, {
				method: 'post',
				headers: {
					accept: 'application/json'
				},
				body: data
			});

			if (response.ok) {
				goto(`/games/${game.id}`);
			} else {
				const resolvedResponse = await response?.json();
				error = resolvedResponse.error;
				console.error(error);
			}
		}
	};

	let error = '';
</script>

<ItemTitle>Join Game</ItemTitle>
<div class="text text-error">{error}</div>

{#if game}
	<div class="flex flex-col place-items-center">
		<GameCard {game} />
	</div>

	<form on:submit|preventDefault={onSubmit}>
		<fieldset name="players" class="form-control mt-3">
			<PlayerChooser bind:raceId />

			<ColorInput bind:value={color} name="color" />
			<div class="text-right text-md italic">
				Note: color may be changed during universe generation
			</div>
		</fieldset>
		<button class="btn btn-primary">Join</button>
	</form>
{/if}
