<script lang="ts">
	import { page } from '$app/stores';
	import type { Game } from '$lib/types/Game';
	import { goto } from '$app/navigation';
	import { me } from '$lib/services/Stores';
	import { onMount } from 'svelte';
	import { GameService } from '$lib/services/GameService';
	import type { PlayerResponse } from '$lib/types/Player';
	import { getGameContext } from '$lib/services/Contexts';

	const { game } = getGameContext();
	let id = parseInt($page.params.id);

	let playerStatuses: PlayerResponse[] = [];

	onMount(async () => {
		playerStatuses = await GameService.loadPlayerStatuses(id);
	});

	const onSubmit = async () => {
		const response = await fetch(`/api/games/${id}/generate`, {
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

<h2 class="font-semibold text-xl">Waiting for players</h2>
<div class="text text-error">{error}</div>

{#if game}
	<div class="">
		<div class="flex">
			<div class="font-semibold w-[8rem]">Name:</div>
			<div class="text-left ml-2">{$game.name}</div>
		</div>
		<div class="flex">
			<div class="font-semibold w-[8rem]">Size:</div>
			<div class="text-left ml-2">{$game.size}</div>
		</div>
		<div class="flex">
			<div class="font-semibold w-[8rem]">Density:</div>
			<div class="text-left ml-2">{$game.density}</div>
		</div>
		<div class="flex">
			<div class="font-semibold w-[8rem]">Player Distance:</div>
			<div class="text-left ml-2">{$game.playerPositions}</div>
		</div>
	</div>

	<ul>
		{#each playerStatuses as status}
			<li>Player {status.num} ({status.name}) - Submitted Turn: {status.submittedTurn ?? false}</li>
		{/each}
	</ul>
	{#if $me?.id == $game.hostId}
		<form on:submit|preventDefault={onSubmit}>
			<button class="btn btn-primary">Force Generate Turn</button>
		</form>
	{/if}
{/if}
