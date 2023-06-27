<script lang="ts">
	import { goto } from '$app/navigation';
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import GameCard from '$lib/components/game/GameCard.svelte';
	import GameSettingsEditor from '$lib/components/game/newgame/GameSettingsEditor.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import { GameService } from '$lib/services/GameService';
	import { Service } from '$lib/services/Service';
	import { me } from '$lib/services/Stores';
	import type { GameSettings } from '$lib/types/Game';
	import { CheckBadge, XMark } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { onDestroy, onMount } from 'svelte';

	const { game } = getGameContext();

	let settings: GameSettings = {
		name: $game.name,
		public: $game.public,
		quickStartTurns: $game.quickStartTurns,
		size: $game.size,
		area: $game.area,
		density: $game.density,
		playerPositions: $game.playerPositions,
		randomEvents: $game.randomEvents,
		computerPlayersFormAlliances: $game.computerPlayersFormAlliances,
		publicPlayerScores: $game.publicPlayerScores,
		startMode: $game.startMode,
		year: $game.year,
		victoryConditions: $game.victoryConditions
	};

	const onLeave = async () => {
		const response = await fetch(`/api/games/${$game.id}/leave`, {
			method: 'POST',
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

	const onUpdateGame = async () => {
		const result = await GameService.updateSettings($game.id, settings);
		$game = Object.assign($game, result);
	};

	const onGenerateUniverse = async () => {
		const response = await fetch(`/api/games/${$game.id}/generate-universe`, {
			method: 'POST',
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			await Service.throwError(response);
		}
		// force an update so the game reloads
		await $game.loadStatus();
		goto(`/games/${$game.id}`);
	};
	let error = '';

	onMount(async () => {
		await $game.loadStatus();
		$game.startPollingStatus();
	});

	onDestroy(() => $game.stopPollingStatus());

	$: isHost = $me?.id === $game.hostId;
</script>

<div class="w-full mx-auto md:max-w-2xl">
	<ItemTitle>{$game.name}</ItemTitle>

	<form class="mt-2">
		<div class="flex flex-col justify-center gap-2 place-items-center">
			<div>
				{#if isHost}
					<GameSettingsEditor bind:settings />
				{:else}
					<GameCard game={$game} />
				{/if}
			</div>
			<div class="w-full bg-base-200 shadow rounded-sm border-2 border-base-300 pt-2 m-1 ">
				<div class="grid grid-cols-2 gap-x-5 px-2" class:grid-cols-3={isHost}>
					<div class="text-center border-b border-b-secondary mb-1">Player</div>
					<div class="text-center border-b border-b-secondary mb-1 font-semibold text-xl">Status</div>
					{#if isHost}
						<div class="border-b border-b-secondary mb-1" ></div>
					{/if}

					{#each $game.players as playerStatus}
						<div class="flex flex-row h-6">
							<div class="w-4">
								{playerStatus.num}
							</div>
							<div
								class="h-4 w-4 my-auto border border-secondary mx-2"
								style={`background-color: ${playerStatus.color}`}
							/>
							{playerStatus.name}
						</div>
						<div>
							{#if playerStatus.ready}
								<div class="flex flex-row">
									<div class="w-20">Ready</div>
									<Icon src={CheckBadge} size="24" class="stroke-success" />
								</div>
							{:else}
								<div class="flex flex-row">
									<div class="w-20">Open</div>
									<Icon src={XMark} size="24" class="stroke-error" />
								</div>
							{/if}
						</div>
						{#if isHost}
							<div class="flex justify-center">
								<button type="button" class="btn btn-outline btn-sm my-1 normal-case">Kick</button>
							</div>
						{/if}
					{/each}
				</div>
			</div>
		</div>
		{#if isHost}
			<button class="btn btn-primary" on:click|preventDefault={onUpdateGame}>Update Game</button>
			<button class="btn btn-secondary" on:click|preventDefault={onGenerateUniverse}
				>Generate Universe</button
			>
		{:else}
			<button class="btn btn-secondary" on:click|preventDefault={onLeave}>Leave Game</button>
		{/if}
	</form>
</div>
