<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { GameState } from '$lib/types/Game';
	import { onDestroy, onMount } from 'svelte';
	import PlayerStatus from './PlayerStatus.svelte';

	const { game, universe, loadStatus, startPollingStatus, stopPollingStatus } = getGameContext();

	$: settingUp = $game.state === GameState.Setup;

	onMount(async () => {
		await loadStatus();
		startPollingStatus();
	});

	onDestroy(() => stopPollingStatus());
</script>

<div class:grid-cols-2={settingUp} class:grid-cols-3={!settingUp} class="grid px-2">
	<div>Player</div>
	{#if !settingUp}
		<div class="font-semibold text-xl">Race</div>
	{/if}
	<div class="font-semibold text-xl">Status</div>
	{#each $game.players as playerStatus}
		<div class="flex flex-row">
			<div class="w-4 my-auto">
				{playerStatus.num}
			</div>
			<div
				class="h-4 w-4 my-auto border border-secondary mx-2"
				style={`background-color: ${playerStatus.color}`}
			/>
			<div class="my-auto">
				{playerStatus.name}
			</div>
		</div>
		{#if !settingUp}
			<div class="my-auto">
				{$universe.getPlayerIntel(playerStatus.num)?.racePluralName ?? 'unknown'}
			</div>
		{/if}
		<PlayerStatus {playerStatus} />
	{/each}
</div>
