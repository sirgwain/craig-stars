<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { me } from '$lib/services/Stores';
	import { onDestroy, onMount } from 'svelte';
	import GameStatus from '../GameStatus.svelte';
	import { playerClient } from '$lib/services/connect';

	const {
		game,
		player,
		forceGenerateTurn,
		loadStatus,
		startPollingStatus,
		stopPollingStatus,
		updateGame
	} = getGameContext();

	async function onForceGenerate() {
		if (
			confirm(
				'Some players have not submitted their turns, are you sure you want to generate a new turn?'
			)
		) {
			await forceGenerateTurn();
		}
	}

	async function onUnsubmitTurn() {
		await playerClient.unsubmitTurn({ gameId: $game.id });
		$player.submittedTurn = false;
		// trigger reactivity on the layout so our hotkeys are wired up again
		// (not super happy with this workaround...)
		updateGame($game);
	}

	// poll for game status when this view is shown
	onMount(async () => {
		await loadStatus();
		startPollingStatus();
	});
	onDestroy(stopPollingStatus);
</script>

<GameStatus title="Waiting for players to play" game={$game.toGameWithPlayers()}>
	<form>
		<div class="gap-2 mt-2">
			{#if $me.id == $game.hostId}
				<button onclick={onForceGenerate} type="button" class="btn btn-primary"
					>Force Generate Turn</button
				>
			{/if}
			<button onclick={onUnsubmitTurn} type="button" class="btn btn-secondary">Unsubmit Turn</button
			>
		</div>
	</form>
</GameStatus>
