<script lang="ts">
	import { page } from '$app/stores';
	import { GameService } from '$lib/services/GameService';
	import type { Game } from '$lib/types/Game';

	import { goto } from '$app/navigation';
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import GameCard from '$lib/components/game/GameCard.svelte';
	import { Service } from '$lib/services/Service';
	import { me } from '$lib/services/Stores';
	import { humanoid } from '$lib/types/Race';
	import { UserRole } from '$lib/types/User';
	import { onMount } from 'svelte';
	import PlayerChooser from '../../../../lib/components/game/newgame/PlayerChooser.svelte';

	let game: Game | undefined;
	let race = Object.assign({}, humanoid());
	let name = $me.username;

	onMount(async () => {
		try {
			let id = parseInt($page.params.id);
			game = await GameService.loadGame(id);
		} catch {
			error = 'No open game found for the invite';
		}
	});

	const onSubmit = async () => {
		if (game) {
			const data = JSON.stringify({ race, name });

			const response = await fetch(`/api/games/${game.id}/join`, {
				method: 'POST',
				headers: {
					accept: 'application/json'
				},
				body: data
			});

			if (!response.ok) {
				await Service.throwError(response);
			}
			goto(`/games/${game.id}`);
		}
	};

	let error = '';
	$: valid = game && game.openPlayerSlots > 0;
</script>

<ItemTitle>Join Public Game</ItemTitle>
<div class="text text-error">{error}</div>

{#if game}
	<div class="flex flex-col place-items-center">
		<GameCard {game} />
	</div>

	<form on:submit|preventDefault={onSubmit}>
		{#if $me.role == UserRole.guest}
			<label class="label" for="name">Name</label>
			<input name="name" bind:value={name} class="input input-bordered" />
		{/if}
		<fieldset name="players" class="form-control mt-3">
			<PlayerChooser bind:race bind:valid />
		</fieldset>
		<button class="btn btn-primary mt-2" disabled={!valid}>Join</button>
	</form>
{/if}
