<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { GameState } from '$lib/types/Game';
	import Game from './Game.svelte';
	import GameSetup from './GameSetup.svelte';
	import WaitingForPlayers from './WaitingForPlayers.svelte';
	const { game, player, universe } = getGameContext();
</script>

{#if $game.state == GameState.Setup}
	<GameSetup />
{:else if $game.state == GameState.GeneratingTurn || $game.state == GameState.GeneratingTurnError || $game.state == GameState.GeneratingUniverse}
	<WaitingForPlayers />
{:else if $player.submittedTurn && $game.state == GameState.WaitingForPlayers}
	<WaitingForPlayers />
{:else}
	<Game />
{/if}
