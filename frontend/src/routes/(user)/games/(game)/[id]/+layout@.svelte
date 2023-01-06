<script lang="ts">
	import { page } from '$app/stores';
	import Menu from '$lib/components/Menu.svelte';
	import NotFound from '$lib/components/NotFound.svelte';
	import { EventManager } from '$lib/EventManager';
	import { bindNavigationHotkeys, unbindNavigationHotkeys } from '$lib/navigationHotkeys';
	import { bindQuantityModifier, unbindQuantityModifier } from '$lib/quantityModifier';
	import {
		commandMapObject,
		game,
		commandedPlanet,
		getMyMapObjectsByPosition,
		player,
		selectMapObject,
		zoomToMapObject,
		me
	} from '$lib/services/Context';
	import { GameService } from '$lib/services/GameService';
	import { PlayerService } from '$lib/services/PlayerService';
	import type { Fleet } from '$lib/types/Fleet';
	import { GameState } from '$lib/types/Game';
	import { MapObjectType } from '$lib/types/MapObject';
	import type { CommandedPlanet, Planet } from '$lib/types/Planet';
	import { onMount } from 'svelte';
	import CargoTransferDialog from './dialogs/cargo/CargoTransferDialog.svelte';
	import ProductionQueueDialog from './dialogs/ProductionQueueDialog.svelte';
	import GameMenu from './GameMenu.svelte';

	let id = parseInt($page.params.id);
	let gameService: GameService = new GameService();

	let source: Fleet | undefined;
	let dest: Fleet | Planet | undefined;

	let loadAttempted = false;

	onMount(async () => {
		if ($game?.id !== id || !$game || !$player) {
			game.update(() => undefined);
			player.update(() => undefined);

			try {
				// load the game on mount
				const result = await gameService.loadGame(id);
				game.update(() => result.game);
				player.update(() => result.player);
			} finally {
				loadAttempted = true;
			}
		}

		// setup the quantityModifier
		bindQuantityModifier();

		// subscribe to events

		const unsubscribes: (() => void)[] = [];
		unsubscribes.push(
			EventManager.subscribeProductionQueueDialogRequestedEvent((planet) =>
				showProductionQueueDialog(planet)
			)
		);
		unsubscribes.push(
			EventManager.subscribeCargoTransferDialogRequestedEvent((src, target) =>
				showCargoTransferDialog(src, target)
			)
		);

		// if we are in an active game, bind the navigation hotkeys, i.e. F4 for research, Esc to go back
		if ($game?.state == GameState.WaitingForPlayers) {
			bindNavigationHotkeys(id);
		}

		return () => {
			unbindQuantityModifier();
			unbindNavigationHotkeys();
			unsubscribes.forEach((unsubscribe) => unsubscribe.apply(unsubscribe));
		};
	});

	// all other components will use this context
	$: if ($game && $player) {
		if ($game.state == GameState.WaitingForPlayers) {
			// setGameContext(game, player);
			const homeworld = $player.planets.find((p) => p.homeworld);
			if (homeworld) {
				commandMapObject(homeworld);
				selectMapObject(homeworld);
				zoomToMapObject(homeworld);
			} else {
				commandMapObject($player.planets[0]);
				selectMapObject($player.planets[0]);
				zoomToMapObject($player.planets[0]);
			}
		}
	}

	async function onSubmitTurn() {
		if ($player) {
			const result = await PlayerService.submitTurn($player);
			if (result !== undefined) {
				game.update((store) => (store = result.game));
				player.update((store) => (store = result.player));
			}
		}
	}

	let productionQueueDialogOpen: boolean;
	const showProductionQueueDialog = (planet: CommandedPlanet) => {
		productionQueueDialogOpen = !productionQueueDialogOpen;
	};

	let cargoTransferDialogOpen: boolean;
	const showCargoTransferDialog = (src: Fleet, target?: Fleet | Planet): void => {
		if (src.spec?.cargoCapacity === 0 && target?.type != MapObjectType.Fleet) {
			// can't transfer cargo with no cargo capcity
			// we can only transfer fuel to another fleet, so don't show the dialog at all in this case
			return;
		}

		if (!target) {
			// no explicit target checked, see if this fleet is orbiting a planet, otherwise it's a jettison
			const myMapObjectsAtPosition = getMyMapObjectsByPosition(src);
			dest = myMapObjectsAtPosition.find((mo) => mo.type == MapObjectType.Planet) as Planet;
		} else {
			dest = target;
		}

		source = src;
		cargoTransferDialogOpen = !cargoTransferDialogOpen;
	};

	const onCargoTransferDialogOk = () => {
		// let any subscribed components know we transferred cargo
		source && EventManager.publishCargoTransferredEvent(source);
		dest && EventManager.publishCargoTransferredEvent(dest);

		// close the dialog
		cargoTransferDialogOpen = false;
	};
</script>

{#if $game && $player}
	<main class="flex flex-col h-screen">
		<div class="flex-initial">
			<GameMenu game={$game} on:submit-turn={onSubmitTurn} />
		</div>
		<div class="p-2 flex-1 overflow-y-auto">
			<slot>Game</slot>
		</div>
	</main>
	<div class="modal" class:modal-open={productionQueueDialogOpen}>
		<div
			class="modal-box max-w-full max-h-max h-full w-full lg:max-w-[40rem] lg:max-h-[48rem] p-0 md:p-[1.25rem]"
		>
			{#if $commandedPlanet}
				<ProductionQueueDialog
					player={$player}
					planet={$commandedPlanet}
					on:ok={() => (productionQueueDialogOpen = false)}
					on:cancel={() => (productionQueueDialogOpen = false)}
				/>
			{/if}
		</div>
	</div>

	<div class="modal" class:modal-open={cargoTransferDialogOpen}>
		<div class="modal-box max-w-full max-h-max h-full lg:max-w-[40rem] lg:max-h-[48rem]">
			<CargoTransferDialog
				src={source}
				{dest}
				on:ok={onCargoTransferDialogOk}
				on:cancel={() => (cargoTransferDialogOpen = false)}
			/>
		</div>
	</div>
{:else if loadAttempted}
	<NotFound title="Game not found" />
{/if}
