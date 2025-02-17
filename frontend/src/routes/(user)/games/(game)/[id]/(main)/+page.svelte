<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import {
		GameStateGeneratingTurn,
		GameStateGeneratingTurnError,
		GameStateGeneratingUniverse,
		GameStateSetup,
		GameStateWaitingForPlayers
	} from '$lib/types/cs';
	import Game from './Game.svelte';
	import GameSetup from './GameSetup.svelte';
	import WaitingForPlayers from './WaitingForPlayers.svelte';
	const { game, player, fullyLoaded } = getGameContext();
</script>

{#if $game.state == GameStateSetup}
	<GameSetup />
{:else if $game.state == GameStateGeneratingTurn || $game.state == GameStateGeneratingTurnError || $game.state == GameStateGeneratingUniverse}
	<WaitingForPlayers />
{:else if $player.submittedTurn && $game.state == GameStateWaitingForPlayers}
	<WaitingForPlayers />
{:else if $game.state == GameStateWaitingForPlayers && $fullyLoaded}
	<Game />
{/if}
