<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import GameCard from '$lib/components/game/GameCard.svelte';
	import { me } from '$lib/services/Stores';
	import { gameClient } from '$lib/services/connect';
	import { getGameWithPlayersFlat, type GameWithPlayersFlat } from '$lib/types/Game';
	import { humanoid } from '$lib/types/Race';
	import { onMount } from 'svelte';
	import PlayerChooser from '../../../../lib/components/game/newgame/PlayerChooser.svelte';
	import { UserRole } from '$lib/types/cs-proto';

	let game: GameWithPlayersFlat | undefined = $state();
	let race = $state(humanoid());
	let name = $state($me.username);

	onMount(async () => {
		const resp = await gameClient.getGameByInviteHash({ hash: page.params.hash });
		if (resp.game) {
			game = getGameWithPlayersFlat(resp.game);
		}
	});

	const onSubmit = async () => {
		if (game) {
			await gameClient.joinGame({ gameId: game.id, name, race });
			goto(`/games/${game.id}`);
		}
	};

	let valid = $derived(!!(game && (game.openPlayerSlots ?? 0) > 0));
</script>

<ItemTitle>Join Private Game</ItemTitle>

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
					valid = raceValid && !!(game && (game.openPlayerSlots ?? 0) > 0);
				}}
			/>
		</fieldset>
		<button class="btn btn-primary mt-2" disabled={!valid}>Join</button>
	</form>
{/if}
