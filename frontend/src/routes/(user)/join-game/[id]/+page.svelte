<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import GameCard from '$lib/components/game/GameCard.svelte';
	import PlayerChooser from '$lib/components/game/newgame/PlayerChooser.svelte';
	import { GameService } from '$lib/services/GameService';
	import { Service } from '$lib/services/Service';
	import { me } from '$lib/services/Stores';
	import { humanoid } from '$lib/types/Race';
	import { RoleGuest, type GameWithPlayers } from '$lib/types/cs';
	import { onMount } from 'svelte';

	let game: GameWithPlayers | undefined = $state();
	let race = $state(humanoid());
	let name = $state($me.username);
	let valid: boolean = $state(false);
	let error = $state('');

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

			await goto(`/games/${game.id}`);
		}
	};

	$effect(() => {
		valid = !!(game && (game.openPlayerSlots ?? 0) > 0);
	});
</script>

<ItemTitle>Join Public Game</ItemTitle>
<div class="text text-error">{error}</div>

{#if game}
	<div class="flex flex-col place-items-center">
		<GameCard {game} />
	</div>

	<form
		onsubmit={(e) => {
			e.preventDefault();
			onSubmit();
		}}
	>
		{#if $me.role === RoleGuest}
			<label class="label" for="name">Name</label>
			<input name="name" bind:value={name} class="input input-bordered" />
		{/if}
		<fieldset name="players" class="form-control mt-3">
			<PlayerChooser
				raceUpdated={(updated, raceValid) => {
					race = updated;
					valid = raceValid && !!(game && (game.openPlayerSlots ?? 0) > 0);
				}}
			/>
		</fieldset>
		<button class="btn btn-primary mt-2" disabled={!valid}>Join</button>
	</form>
{/if}
