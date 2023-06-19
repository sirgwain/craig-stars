<script lang="ts">
	import { game } from '$lib/services/Context';
	import { GameState } from '$lib/types/Game';
	import Game from './Game.svelte';
	import GameSetup from './GameSetup.svelte';
	import WaitingForPlayers from './WaitingForPlayers.svelte';
</script>

{#if $game && $game.state == GameState.Setup}
	<GameSetup game={$game} />
{:else if $game?.state == GameState.GeneratingTurn}
	Generating turn, refresh
{:else if $game.player.submittedTurn && $game?.state == GameState.WaitingForPlayers}
	<WaitingForPlayers game={$game} />
{:else}
	<Game game={$game} />
{/if}
