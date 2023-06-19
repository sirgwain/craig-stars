<script lang="ts">
	import { page } from '$app/stores';
	import GameMenu from '$lib/components/game/GameMenu.svelte';
	import ProductionQueue from '$lib/components/game/ProductionQueue.svelte';
	import { EventManager } from '$lib/EventManager';
	import { commandMapObject, game, player, selectMapObject } from '$lib/services/Context';
	import { GameService } from '$lib/services/GameService';
	import { PlayerService } from '$lib/services/PlayerService';
	import type { Planet } from '$lib/types/Planet';
	import { onMount } from 'svelte';

	let id = parseInt($page.params.id);
	let playerService: PlayerService;
	let gameService: GameService = new GameService();

	onMount(async () => {
		// load the game on mount
		const result = await gameService.loadGame(id);
		game.update((store) => (store = result.game));
		player.update((store) => (store = result.player));

		playerService = new PlayerService($player);
	});

	// all other components will use this context
	$: if ($game && $player) {
		// setGameContext(game, player);
		const homeworld = $player.planets.find((p) => p.homeworld);
		if (homeworld) {
			commandMapObject(homeworld);
			selectMapObject(homeworld);
		} else {
			commandMapObject($player.planets[0]);
			selectMapObject($player.planets[0]);
		}
	}

	async function onSubmitTurn() {
		const result = await playerService.submitTurn();
		if (result !== undefined) {
			game.update((store) => (store = result.game));
			player.update((store) => (store = result.player));
		}
	}

	let productionQueueDialogOpen: boolean;
	const showProductionQueueDialog = (planet?: Planet | undefined) => {
		productionQueueDialogOpen = !productionQueueDialogOpen;
	};

	EventManager.productionQueueDialogRequestedEvent = (planet) => showProductionQueueDialog(planet);
</script>

{#if $game && $player}
	<main class="flex flex-col h-screen">
		<div class="flex-initial">
			<GameMenu on:submit-turn={onSubmitTurn} />
		</div>
		<div class="p-2 flex-1">
			<slot>Game</slot>
		</div>
	</main>
	<div class="modal" class:modal-open={productionQueueDialogOpen}>
		<div class="modal-box max-w-full max-h-max h-full lg:max-w-[40rem] lg:max-h-[48rem]">
			<ProductionQueue
				on:ok={() => (productionQueueDialogOpen = false)}
				on:cancel={() => (productionQueueDialogOpen = false)}
			/>
		</div>
	</div>
{/if}
