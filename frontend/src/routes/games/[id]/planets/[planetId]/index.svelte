<script lang="ts">
	import { page } from '$app/stores';
	import ProductionQueue from '$lib/components/game/ProductionQueue.svelte';
	import { setGameContext } from '$lib/services/Context';
	import { GameService } from '$lib/services/GameService';
	import { PlayerService } from '$lib/services/PlayerService';
	import type { Game } from '$lib/types/Game';
	import type { Planet } from '$lib/types/Planet';
	import type { Player } from '$lib/types/Player';
	import { onMount } from 'svelte';

	let gameId = parseInt($page.params.id);
	let planetId = parseInt($page.params.gameId);

	let player: Player;
	let game: Game;
	let playerService: PlayerService;
	let gameService: GameService = new GameService();
	let planet: Planet;

	onMount(async () => {
		// load the game on mount
		({ game, player } = await gameService.loadGame(gameId));
		playerService = new PlayerService(player);
	});

	// all other components will use this context
	$: if (game && player) {
		setGameContext(game, player);
		const homeworld = player.planets.find((p) => p.homeworld);
		if (homeworld) {
			planet = homeworld;
		}
	}
</script>

{#if game}
	<ProductionQueue {planet} />
{/if}
