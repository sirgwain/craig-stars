<script lang="ts">
	import { goto } from '$app/navigation';
	import GameStatus from '../GameStatus.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import { me } from '$lib/services/Stores';

	const { game } = getGameContext();

	const onLeave = async () => {
		const response = await fetch(`/api/games/${$game.id}/leave`, {
			method: 'post',
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			goto(`/games`);
		} else {
			const resolvedResponse = await response?.json();
			error = resolvedResponse.error;
			console.error(error);
		}
	};

	const onSubmit = async () => {
		const response = await fetch(`/api/games/${$game.id}/generate-universe`, {
			method: 'post',
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			// force an update so the game reloads
			await $game.loadPlayersStatus();
			goto(`/games/${$game.id}`);
		} else {
			const resolvedResponse = await response?.json();
			error = resolvedResponse.error;
			console.error(error);
		}
	};
	let error = '';
</script>

<GameStatus title="Waiting for players to join" game={$game}>
	{#if $me?.id == $game.hostId}
		<form class="mt-2" on:submit|preventDefault={onSubmit}>
			<button class="btn btn-primary">Generate Universe</button>
		</form>
	{:else}
		<form class="mt-2" on:submit|preventDefault={onLeave}>
			<button class="btn btn-primary">Leave Game</button>
		</form>
	{/if}
</GameStatus>
