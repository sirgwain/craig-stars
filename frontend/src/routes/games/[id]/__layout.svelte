<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import ProductionQueue from '$lib/components/game/ProductionQueue.svelte';
	import Toolbar from '$lib/components/game/Toolbar.svelte';
	import { EventManager } from '$lib/EventManager';
	import { commandPlanet, game, player, selectPlanet } from '$lib/services/Context';
	import { GameService } from '$lib/services/GameService';
	import { PlayerService } from '$lib/services/PlayerService';
	import type { Game } from '$lib/types/Game';
	import type { Planet } from '$lib/types/Planet';
	import type { Player } from '$lib/types/Player';

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
			commandPlanet(homeworld);
			selectPlanet(homeworld);
		} else {
			commandPlanet($player.planets[0]);
			selectPlanet($player.planets[0]);
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
	<div class="flex flex-col h-full">
		<div class="flex-none">
			<Toolbar on:submit-turn={onSubmitTurn} />
		</div>

		<div class="flex-1">
			<div class="flex h-full">
				<slot />
			</div>
		</div>

		<div class="modal" class:modal-open={productionQueueDialogOpen}>
			<div class="modal-box max-w-full max-h-max h-full lg:max-w-[40rem] lg:max-h-[48rem]">
				<ProductionQueue
					on:ok={() => (productionQueueDialogOpen = false)}
					on:cancel={() => (productionQueueDialogOpen = false)}
				/>
			</div>
		</div>
	</div>
{/if}
