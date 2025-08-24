<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import GameCard from '$lib/components/game/GameCard.svelte';
	import PlayerChooser from '$lib/components/game/newgame/PlayerChooser.svelte';
	import { me } from '$lib/services/Stores';
	import { gameClient } from '$lib/services/connect';
	import { getGameWithPlayersFlat, type GameWithPlayersFlat } from '$lib/types/Game';
	import { humanoid } from '$lib/types/Race';
	import { onMount } from 'svelte';
	import { UserRole } from '$lib/types/cs-proto';

	let game: GameWithPlayersFlat | undefined = $state();
	let race = $state(humanoid());
	let name = $state($me.username);
	let valid: boolean = $derived(!!(game && game.openPlayerSlots > 0));
	let error = $state('');

	onMount(async () => {
		try {
			const resp = await gameClient.getGame({ gameId: BigInt(page.params.id) });
			if (resp.game) {
				game = getGameWithPlayersFlat(resp.game);
			}
		} catch {
			error = 'No open game found for the invite';
		}
	});

	const onSubmit = async () => {
		if (game) {
			await gameClient.joinGame({ name, gameId: game.id, race });
			await goto(`/games/${game.id}`);
		}
	};
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
		{#if $me.role === UserRole.GUEST}
			<label class="label" for="name">Name</label>
			<input name="name" bind:value={name} class="input input-bordered" />
		{/if}
		<fieldset name="players" class="form-control mt-3">
			<PlayerChooser
				raceUpdated={(updated, raceValid) => {
					race = updated;
					valid =
						raceValid &&
						!!(
							(game && game.openPlayerSlots > 0) ||
							// guests are already joined, just need to pick their race
							($me.isGuest() && game?.players.some((p) => p.userId === $me.id))
						);
				}}
			/>
		</fieldset>
		<button class="btn btn-primary mt-2" disabled={!valid}>Join</button>
	</form>
{/if}
