<script lang="ts">
	import { goto } from '$app/navigation';
	import { getGameContext, playerFinderKey } from '$lib/services/Contexts';
	import { GameService } from '$lib/services/GameService';
	import { PlayerService } from '$lib/services/PlayerService';
	import { me } from '$lib/services/Stores';
	import GameStatus from '../GameStatus.svelte';

	const { game, player } = getGameContext();

	async function onForceGenerate() {
		if (
			confirm(
				'Some players have not submitted their turns, are you sure you want to generate a new turn?'
			)
		) {
			await $game.forceGenerateTurn();
		}
	}

	async function onUnsubmitTurn() {
		await PlayerService.unsubmitTurn($game.id);
		$player.submittedTurn = false;
	}
	let error = '';
</script>

<GameStatus title="Waiting for players to play" game={$game}>
	<form>
		<div class="gap-2 mt-2">
			{#if $me?.id == $game.hostId}
				<button on:click={onForceGenerate} type="button" class="btn btn-primary"
					>Force Generate Turn</button
				>
			{/if}
			<button on:click={onUnsubmitTurn} type="button" class="btn btn-secondary"
				>Unsubmit Turn</button
			>
		</div>
	</form>
</GameStatus>
