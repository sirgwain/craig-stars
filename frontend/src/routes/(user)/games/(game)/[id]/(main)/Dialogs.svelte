<script lang="ts">
	import { EventManager } from '$lib/EventManager';
	import Tooltip from '$lib/components/game/tooltips/Tooltip.svelte';
	import { commandedFleet, commandedPlanet } from '$lib/services/Context';
	import type { FullGame } from '$lib/services/FullGame';
	import type { CommandedFleet, Fleet } from '$lib/types/Fleet';
	import { MapObjectType } from '$lib/types/MapObject';
	import type { CommandedPlanet, Planet } from '$lib/types/Planet';
	import { onMount } from 'svelte';
	import ProductionQueue from '../dialogs/ProductionQueue.svelte';
	import CargoTransferDialog from '../dialogs/cargo/CargoTransferDialog.svelte';
	import MergeFleets from '../fleets/[num]/merge/MergeFleets.svelte';

	export let game: FullGame;

	let source: CommandedFleet | undefined;
	let dest: Fleet | Planet | undefined;

	onMount(async () => {
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

		return () => {
			unsubscribes.forEach((unsubscribe) => unsubscribe.apply(unsubscribe));
		};
	});

	let productionQueueDialogOpen: boolean;
	const showProductionQueueDialog = (planet: CommandedPlanet) => {
		productionQueueDialogOpen = !productionQueueDialogOpen;
	};

	let mergeFleetDialogOpen: boolean;
	const showMergeFleetDialogOpen = () => {
		mergeFleetDialogOpen = !mergeFleetDialogOpen;
	};

	let cargoTransferDialogOpen: boolean;
	const showCargoTransferDialog = (src: CommandedFleet, target?: Fleet | Planet): void => {
		if (src.spec?.cargoCapacity === 0 && target?.type != MapObjectType.Fleet) {
			// can't transfer cargo with no cargo capcity
			// we can only transfer fuel to another fleet, so don't show the dialog at all in this case
			return;
		}

		if (!target) {
			// no explicit target checked, see if this fleet is orbiting a planet, otherwise it's a jettison
			dest = game.universe.getMyPlanetsByPosition(src)[0];
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
		game.merge(fleet, fleetNums);
		mergeFleetDialogOpen = false;
	};
</script>

<div class="modal" class:modal-open={productionQueueDialogOpen}>
	<div
		class="modal-box max-w-full max-h-max h-full w-full lg:max-w-[40rem] lg:max-h-[48rem] p-0 md:p-[1.25rem]"
	>
		{#if $commandedPlanet}
			<ProductionQueue
				{game}
				player={game.player}
				designs={game.player.designs}
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
				otherFleetsHere={game.universe
					.getMyFleetsByPosition($commandedFleet)
					.filter((f) => f.num !== $commandedFleet?.num)}
				on:ok={(e) => onMergeFleetsDialogOk(e.detail.fleet, e.detail.fleetNums)}
				on:cancel={() => (mergeFleetDialogOpen = false)}
			/>
		{/if}
	</div>
</div>

<Tooltip />
