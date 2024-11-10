<script lang="ts">
	import type {
		CargoTransferDialogEvent,
		MergeFleetsDialogEvent,
		MergeFleetsEvent,
		SplitFleetDialogEvent,
		SplitFleetEvent,
		TransferCargoEvent,
		TransportTasksDialogEvent,
		TransportTasksUpdateEvent
	} from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import { ownedBy, type MapObject } from '$lib/types/MapObject';
	import { newSalvage } from '$lib/types/Salvage';
	import hotkeys from 'hotkeys-js';
	import { onMount } from 'svelte';
	import CargoTranfserDialog from '../dialogs/cargo/CargoTransferDialog.svelte';
	import MergeFleetsDialog from '../dialogs/merge/MergeFleetsDialog.svelte';
	import ProductionQueueDialog from '../dialogs/production/ProductionQueueDialog.svelte';
	import SplitFleetDialog from '../dialogs/split/SplitFleetDialog.svelte';
	import TransportTasksDialog from '../dialogs/transport/TransportTasksDialog.svelte';
	import SearchDialog from '../search/SearchDialog.svelte';
	import MapObjectStatsBar from './MapObjectStatsBar.svelte';
	import MapObjectSummary from './MapObjectSummary.svelte';
	import CommandPane from './command/CommandPane.svelte';
	import CommandPaneCarousel from './command/CommandPaneCarousel.svelte';
	import Scanner from './scanner/Scanner.svelte';
	import ScannerToolbar from './scanner/ScannerToolbar.svelte';
	import type { CommandedPlanet } from '$lib/types/Planet';

	const {
		game,
		universe,
		player,
		commandedPlanet,
		commandedFleet,
		selectedWaypoint,
		currentSelectedWaypointIndex,
		commandMapObject,
		zoomToMapObject,
		nextMapObject,
		previousMapObject,
		selectWaypoint,
		selectMapObject,
		updateFleetOrders,
		updatePlanetOrders,
		transferCargo,
		split,
		splitAll,
		merge
	} = getGameContext();

	let carouselOpen = $state(true);
	let showProductionQueueDialog = $state(false);
	let showCargoTransferDialog = $state(false);
	let showMergeFleetsDialog = $state(false);
	let showSplitFleetDialog = $state(false);
	let showTransportTasksDialog = $state(false);
	let showSearchDialog = $state(false);
	let cargoTransferDialogEvent: CargoTransferDialogEvent | undefined = $state(undefined);
	let mergeFleetsDialogEvent: MergeFleetsDialogEvent | undefined = $state(undefined);
	let splitFleetDialogEvent: SplitFleetDialogEvent | undefined = $state(undefined);
	let transportTasksDialogEvent: TransportTasksDialogEvent | undefined = $state(undefined);

	onMount(() => {
		hotkeys('n', 'root', () => {
			nextMapObject();
		});
		hotkeys('p', 'root', () => {
			previousMapObject();
		});
		hotkeys('q', 'root', () => {
			if ($commandedPlanet) {
				showProductionQueueDialog = true;
			}
		});
		hotkeys('⌘+k', 'root', () => {
			showSearchDialog = true;
		});
		hotkeys('Delete', 'root', () => {
			onDeleteWaypoint();
		});
		hotkeys('Backspace', 'root', () => {
			onDeleteWaypoint();
		});

		// on the game view, prevent scrolling
		document.documentElement.classList.add('overscroll-none');

		return () => {
			hotkeys.unbind('n', 'root');
			hotkeys.unbind('p', 'root');
			hotkeys.unbind('q', 'root');
			hotkeys.unbind('Delete', 'root');
			hotkeys.unbind('Backspace', 'root');

			// go back to normal
			document.documentElement.classList.remove('overscroll-none');
		};
	});

	async function onDeleteWaypoint() {
		const selectedWaypointIndex = $currentSelectedWaypointIndex;
		if (selectedWaypoint && $commandedFleet && selectedWaypointIndex > 0) {
			$commandedFleet.waypoints = $commandedFleet.waypoints.filter((wp) => wp != $selectedWaypoint);

			// select the previous waypoint
			const wp = $commandedFleet.waypoints[selectedWaypointIndex - 1];
			selectWaypoint(wp);

			const mo = $universe.getMapObject(wp);
			if (mo) {
				selectMapObject(mo);
			}

			await updateFleetOrders($commandedFleet);
		}
	}

	const onUpdateTransportTasks = async (e: TransportTasksUpdateEvent) => {
		e.waypoint.transportTasks = e.transportTasks;
		await updateFleetOrders(e.fleet);

		// close the dialog
		showTransportTasksDialog = false;
	};

	async function onSplitAll() {
		if (!$commandedFleet) {
			return;
		}
		splitAll($commandedFleet);
	}

	async function onTransferCargo(e: TransferCargoEvent) {
		// close the dialog
		showCargoTransferDialog = false;

		if (e && e.transferAmount.absoluteSize() > 0) {
			if (!e.dest) {
				e.dest = newSalvage();
			}
			await transferCargo(e.src, e.dest, e.transferAmount);
		}
	}

	async function onUpdatePlanetOrders(planet: CommandedPlanet) {
		updatePlanetOrders(planet);
	}

	async function onNextPlanet(updateOrders: boolean) {
		if (!$commandedPlanet) {
			return;
		}
		if (updateOrders) {
			await updatePlanetOrders($commandedPlanet);
		}

		nextMapObject();
	}

	async function onPrevPlanet(updateOrders: boolean) {
		if (!$commandedPlanet) {
			return;
		}
		if (updateOrders) {
			await updatePlanetOrders($commandedPlanet);
		}

		previousMapObject();
	}

	async function onMergeFleets(e: MergeFleetsEvent) {
		await merge(e.fleet, e.fleetNums);
		// close the dialog
		showMergeFleetsDialog = false;
	}

	async function onSplitFleet(e: SplitFleetEvent) {
		await split(e.src, e.dest, e.srcTokens, e.destTokens, e.transferAmount);

		// close the dialog
		showSplitFleetDialog = false;
	}

	function selectSearchResult(mo: MapObject | undefined) {
		if (mo) {
			if (ownedBy(mo, $player.num)) {
				commandMapObject(mo);
			}
			selectMapObject(mo);
			zoomToMapObject(mo);
		}
		showSearchDialog = false;
	}
</script>

<!-- for small mobile displays we put the scanner on top and the command pane below it-->
<div class="flex flex-col h-full md:flex-row">
	<!-- for medium+ displays, command pane goes on the left -->
	<div
		class="hidden overflow-x-hidden md:flex md:flex-col md:flex-none justify-between md:w-[15.5rem] lg:w-[30rem] overflow-y-auto md:max-h-[calc(100dvh-4rem)]"
	>
		<div class="flex flex-row flex-wrap gap-2 justify-center">
			<CommandPane
				{onDeleteWaypoint}
				{onSplitAll}
				onShowProductionQueueDialog={(e) => (showProductionQueueDialog = true)}
				onShowCargoTransferDialog={(e) => {
					showCargoTransferDialog = true;
					cargoTransferDialogEvent = e;
				}}
				onShowMergeFleetDialog={(e) => {
					showMergeFleetsDialog = true;
					mergeFleetsDialogEvent = e;
				}}
				onShowSplitFleetDialog={(e) => {
					showSplitFleetDialog = true;
					splitFleetDialogEvent = e;
				}}
				onShowTransportTasksDialog={(e) => {
					showTransportTasksDialog = true;
					transportTasksDialogEvent = e;
				}}
			/>
		</div>
		<div class="hidden lg:block lg:p-1 mx-2">
			<MapObjectSummary
				onShowCargoTransferDialog={(e) => {
					showCargoTransferDialog = true;
					cargoTransferDialogEvent = e;
				}}
			/>
		</div>
	</div>

	<div class="flex flex-col grow">
		<div class="flex flex-col grow border-gray-700 border-2 shadow-sm">
			<ScannerToolbar on:show-search={() => (showSearchDialog = true)} />
			<Scanner {onDeleteWaypoint} />
		</div>
		<div class:hidden={!carouselOpen}>
			<MapObjectStatsBar />
		</div>
		<div class="hidden md:block md:w-full lg:hidden mb-2">
			<MapObjectSummary
				onShowCargoTransferDialog={(e) => {
					showCargoTransferDialog = true;
					cargoTransferDialogEvent = e;
				}}
			/>
		</div>
	</div>

	<!-- for phone displays, use a carousel -->
	<div class="flex flex-col flex-0">
		<CommandPaneCarousel
			bind:isOpen={carouselOpen}
			{onDeleteWaypoint}
			{onSplitAll}
			onShowProductionQueueDialog={(e) => (showProductionQueueDialog = true)}
			onShowCargoTransferDialog={(e) => {
				showCargoTransferDialog = true;
				cargoTransferDialogEvent = e;
			}}
			onShowMergeFleetDialog={(e) => {
				showMergeFleetsDialog = true;
				mergeFleetsDialogEvent = e;
			}}
			onShowSplitFleetDialog={(e) => {
				showSplitFleetDialog = true;
				splitFleetDialogEvent = e;
			}}
			onShowTransportTasksDialog={(e) => {
				showTransportTasksDialog = true;
				transportTasksDialogEvent = e;
			}}
		/>
	</div>
</div>

<!-- dialog modals -->
<ProductionQueueDialog
	show={showProductionQueueDialog}
	onNext={() => onNextPlanet(true)}
	onPrev={() => onPrevPlanet(true)}
	onOk={(planet) => {
		showProductionQueueDialog = false;
		updatePlanetOrders(planet);
	}}
	onCancel={() => (showProductionQueueDialog = false)}
/>
<CargoTranfserDialog
	show={showCargoTransferDialog}
	props={cargoTransferDialogEvent}
	onOk={onTransferCargo}
	onCancel={() => (showCargoTransferDialog = false)}
/>
<MergeFleetsDialog
	show={showMergeFleetsDialog}
	props={mergeFleetsDialogEvent}
	onOk={onMergeFleets}
	onCancel={() => (showMergeFleetsDialog = false)}
/>
<SplitFleetDialog
	show={showSplitFleetDialog}
	props={splitFleetDialogEvent}
	onOk={onSplitFleet}
	onCancel={() => (showSplitFleetDialog = false)}
/>
<TransportTasksDialog
	props={transportTasksDialogEvent}
	onOk={onUpdateTransportTasks}
	onCancel={() => (showTransportTasksDialog = false)}
/>
<SearchDialog
	show={showSearchDialog}
	onOk={(e) => selectSearchResult(e)}
	onCancel={() => (showSearchDialog = false)}
/>
