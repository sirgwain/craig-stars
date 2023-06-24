<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import { GameState } from '$lib/types/Game';
	import type { PlayerIntel } from '$lib/types/Player';
	import { onDestroy, onMount } from 'svelte';
	import PlayerStatus from './PlayerStatus.svelte';

	const { game, player, universe } = getGameContext();

	$: settingUp = $game.state === GameState.Setup;

	function getPlayerIntel(num: number): PlayerIntel | undefined {
		if (num == $player.num) {
			return $player;
		}
		if (num > 0 && num <= $universe.players.length) {
			return $universe.players[num];
		}
	}

	onMount(async () => {
		await $game.loadPlayersStatus();
		$game.startPollingPlayersStatus();
	});

	onDestroy(() => $game.stopPollingPlayersStatus());
</script>

<div class:grid-cols-2={settingUp} class:grid-cols-3={!settingUp} class="grid px-2">
	<div>Player</div>
	{#if !settingUp}
		<div class="font-semibold text-xl">Race</div>
	{/if}
	<div class="font-semibold text-xl">Status</div>
	{#each $game.playersStatus as playerStatus}
		<div class="flex flex-row">
			<div class="w-4">
				{playerStatus.num}
			</div>
			<div
				class="h-4 w-4 my-auto border border-secondary mx-2"
				style={`background-color: ${playerStatus.color}`}
			/>
			{playerStatus.name}
		</div>
		{#if !settingUp}
			<div>{getPlayerIntel(playerStatus.num)?.racePluralName ?? 'unknown'}</div>
		{/if}
		<div>
			<PlayerStatus {playerStatus} />
		</div>
	{/each}
</div>
