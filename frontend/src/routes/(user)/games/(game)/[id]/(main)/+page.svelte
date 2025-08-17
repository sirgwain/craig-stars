<script lang="ts">
	import { GameState } from '$lib/types/cs-proto';
	import { getGameContext } from '$lib/services/GameContext';
	import Game from './Game.svelte';
	import GameSetup from './GameSetup.svelte';
	import WaitingForPlayers from './WaitingForPlayers.svelte';
	const { game, player, fullyLoaded } = getGameContext();
</script>

{#if $game.state == GameState.SETUP}
	<GameSetup />
{:else if $game.state == GameState.GENERATING_TURN || $game.state == GameState.GENERATING_TURN_ERROR || $game.state == GameState.GENERATING_UNIVERSE}
	<WaitingForPlayers />
{:else if $player.submittedTurn && $game.state == GameState.WAITING_FOR_PLAYERS}
	<WaitingForPlayers />
{:else if $game.state == GameState.WAITING_FOR_PLAYERS && $fullyLoaded}
	<Game />
{/if}
