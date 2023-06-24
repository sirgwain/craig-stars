<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import { GameState } from '$lib/types/Game';
	import Game from './Game.svelte';
	import GameSetup from './GameSetup.svelte';
	import WaitingForPlayers from './WaitingForPlayers.svelte';
	const { game, player, universe } = getGameContext();
</script>

{#if $game.state == GameState.Setup}
	<GameSetup />
{:else if $game.state == GameState.GeneratingTurn}
	<WaitingForPlayers />
{:else if $player.submittedTurn && $game.state == GameState.WaitingForPlayers}
	<WaitingForPlayers />
{:else}
	<Game />
{/if}
