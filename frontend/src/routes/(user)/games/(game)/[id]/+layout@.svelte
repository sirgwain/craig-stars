<script lang="ts">
	import { page } from '$app/stores';
	import CargoTransferDialog from './dialogs/cargo/CargoTransferDialog.svelte';
	import ProductionQueueDialog from './dialogs/ProductionQueueDialog.svelte';
	import GameMenu from './GameMenu.svelte';
	import { EventManager } from '$lib/EventManager';
	import { bindQuantityModifier, unbindQuantityModifier } from '$lib/quantityModifier';
	import {
		commandMapObject,
		game,
		myMapObjectsByPosition,
		player,
		selectMapObject
	} from '$lib/services/Context';
	import { GameService } from '$lib/services/GameService';
	import { PlayerService } from '$lib/services/PlayerService';
	import type { Fleet } from '$lib/types/Fleet';
	import { MapObjectType, positionKey } from '$lib/types/MapObject';
	import type { Planet } from '$lib/types/Planet';
	import { onMount } from 'svelte';

	let id = parseInt($page.params.id);
	let playerService: PlayerService;
	let gameService: GameService = new GameService();

	let source: Fleet | undefined;
	let dest: Fleet | Planet | undefined;

	onMount(async () => {
		game.update(() => undefined);
		player.update(() => undefined);

		// load the game on mount
		const result = await gameService.loadGame(id);
		game.update(() => result.game);
		player.update(() => result.player);

		playerService = new PlayerService(result.player);

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

		return () => {
			unbindQuantityModifier();
			unsubscribes.forEach((unsubscribe) => unsubscribe.apply(unsubscribe));
		};
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

	let cargoTransferDialogOpen: boolean;
	const showCargoTransferDialog = (src: Fleet, target?: Fleet | Planet): void => {
		if (!$myMapObjectsByPosition) {
			return;
		}

		if (src.spec?.cargoCapacity === 0 && target?.type != MapObjectType.Fleet) {
			// can't transfer cargo with no cargo capcity
			// we can only transfer fuel to another fleet, so don't show the dialog at all in this case
			return;
		}

		if (!target) {
			// no explicit target checked, see if this fleet is orbiting a planet, otherwise it's a jettison
			const key = positionKey(src);
			const myMapObjectsAtPosition = $myMapObjectsByPosition[key];
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
			<GameMenu on:submit-turn={onSubmitTurn} />
		</div>
		<div class="p-2 flex-1 overflow-y-auto">
			<slot>Game</slot>
		</div>
	</main>
	<div class="modal" class:modal-open={productionQueueDialogOpen}>
		<div class="modal-box max-w-full max-h-max h-full lg:max-w-[40rem] lg:max-h-[48rem]">
			<ProductionQueueDialog
				on:ok={() => (productionQueueDialogOpen = false)}
				on:cancel={() => (productionQueueDialogOpen = false)}
			/>
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
{/if}
