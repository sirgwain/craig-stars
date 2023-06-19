<script lang="ts">
	import { page } from '$app/stores';
	import { EventManager } from '$lib/EventManager';
	import NotFound from '$lib/components/NotFound.svelte';
	import { bindNavigationHotkeys, unbindNavigationHotkeys } from '$lib/navigationHotkeys';
	import { bindQuantityModifier, unbindQuantityModifier } from '$lib/quantityModifier';
	import {
		commandedFleet,
		commandedMapObject,
		commandedPlanet,
		game,
		techs
	} from '$lib/services/Context';
	import { FullGame } from '$lib/services/FullGame';
	import type { CommandedFleet, Fleet } from '$lib/types/Fleet';
	import { GameState } from '$lib/types/Game';
	import { MapObjectType } from '$lib/types/MapObject';
	import type { CommandedPlanet, Planet } from '$lib/types/Planet';
	import { onMount } from 'svelte';
	import GameMenu from '../GameMenu.svelte';
	import ProductionQueue from '../dialogs/ProductionQueue.svelte';
	import CargoTransferDialog from '../dialogs/cargo/CargoTransferDialog.svelte';
	import MergeFleets from '../fleets/[num]/merge/MergeFleets.svelte';

	let id = parseInt($page.params.id);

	let source: Fleet | undefined;
	let dest: Fleet | Planet | undefined;

	let loadAttempted = false;

	onMount(async () => {
		if (!$game || $game.id !== id) {
			game.update(() => undefined);

			try {
				const fg = new FullGame();
				await fg.load(id);

				game.update(() => fg);
				techs.update(() => fg.techs);
			} finally {
				loadAttempted = true;
			}
		}

		if ($game?.state == GameState.WaitingForPlayers) {
			if (!$commandedMapObject || $commandedMapObject.gameId != $game.id) {
				$game.universe.commandHomeWorld();
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
		unsubscribes.push(
			EventManager.subscribeMergeFleetDialogRequestedEvent(() => showMergeFleetDialogOpen())
		);

		// if we are in an active game, bind the navigation hotkeys, i.e. F4 for research, Esc to go back
		if ($game?.state == GameState.WaitingForPlayers) {
			bindNavigationHotkeys(id, page);
		}

		return () => {
			unbindQuantityModifier();
			unbindNavigationHotkeys();
			unsubscribes.forEach((unsubscribe) => unsubscribe.apply(unsubscribe));
		};
	});

	async function onSubmitTurn() {
		$game = await $game?.submitTurn();
		if ($game?.state == GameState.WaitingForPlayers) {
			$game.universe.commandHomeWorld();
		}
	}

	let productionQueueDialogOpen: boolean;
	const showProductionQueueDialog = (planet: CommandedPlanet) => {
		productionQueueDialogOpen = !productionQueueDialogOpen;
	};

	let mergeFleetDialogOpen: boolean;
	const showMergeFleetDialogOpen = () => {
		mergeFleetDialogOpen = !mergeFleetDialogOpen;
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
			dest = $game?.universe.getMyPlanetsByPosition(src)[0]
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

	const onMergeFleetsDialogOk = async (fleet: CommandedFleet, fleetNums: number[]) => {
		$game?.merge(fleet, fleetNums);
		mergeFleetDialogOpen = false;
	};
</script>

{#if $game}
	<main class="flex flex-col mb-20 md:mb-0">
		<div class="flex-initial">
			<GameMenu game={$game} on:submit-turn={onSubmitTurn} />
		</div>
		<!-- We want our main game view to only fill the screen (minus the toolbar) -->
		<div class="grow viewport">
			<slot>Game</slot>
		</div>
	</main>
	<div class="modal" class:modal-open={productionQueueDialogOpen}>
		<div
			class="modal-box max-w-full max-h-max h-full w-full lg:max-w-[40rem] lg:max-h-[48rem] p-0 md:p-[1.25rem]"
		>
			{#if $commandedPlanet}
				<ProductionQueue
					game={$game}
					player={$game.player}
					designs={$game.player.designs}
					planet={$commandedPlanet}
					on:ok={() => (productionQueueDialogOpen = false)}
					on:cancel={() => (productionQueueDialogOpen = false)}
				/>
			{/if}
		</div>
	</div>

	<div class="modal" class:modal-open={cargoTransferDialogOpen}>
		<div class="modal-box max-w-full max-h-max h-full w-full lg:max-w-[40rem] lg:max-h-[48rem]">
			<CargoTransferDialog
				src={source}
				{dest}
				on:ok={onCargoTransferDialogOk}
				on:cancel={() => (cargoTransferDialogOpen = false)}
			/>
		</div>
	</div>

	<div class="modal" class:modal-open={mergeFleetDialogOpen}>
		<div class="modal-box max-w-full max-h-max h-full w-full md:max-w-[32rem] md:max-h-[32rem]">
			{#if $commandedFleet}
				<MergeFleets
					fleet={$commandedFleet}
					otherFleetsHere={$game.universe
						.getMyFleetsByPosition($commandedFleet)
						.filter((f) => f.num !== $commandedFleet?.num)}
					on:ok={(e) => onMergeFleetsDialogOk(e.detail.fleet, e.detail.fleetNums)}
					on:cancel={() => (mergeFleetDialogOpen = false)}
				/>
			{/if}
		</div>
	</div>
{:else if loadAttempted}
	<NotFound title="Game not found" />
{/if}

<style>
	main {
		height: 100vh; /* Fallback for browsers that do not support Custom Properties */
		height: calc(var(--vh, 1vh) * 100);
	}

	.viewport {
		max-height: calc(100vh-4rem);
		max-height: calc((var(--vh, 1vh) * 100)-4rem);
	}
</style>
