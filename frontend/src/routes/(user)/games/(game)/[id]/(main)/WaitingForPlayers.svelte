<script lang="ts">
	import { goto } from '$app/navigation';
	import { getGameContext } from '$lib/services/Contexts';
	import { me } from '$lib/services/Stores';
	import GameStatus from '../GameStatus.svelte';

	const { game } = getGameContext();

	const onSubmit = async () => {
		const response = await fetch(`/api/games/${$game.id}/generate-turn`, {
			method: 'post',
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			goto('/');
		} else {
			const resolvedResponse = await response?.json();
			error = resolvedResponse.error;
			console.error(error);
		}
	};
	let error = '';
</script>

<GameStatus title="Waiting for players to play" game={$game}>
	{#if $me?.id == $game.hostId}
		<form on:submit|preventDefault={onSubmit}>
			<button class="btn btn-primary">Force Generate Turn</button>
		</form>
	{/if}
</GameStatus>
