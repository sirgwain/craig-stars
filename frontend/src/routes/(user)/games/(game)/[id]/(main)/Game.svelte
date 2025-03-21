<script lang="ts">
	import type {
		BattlePlanChangedEvent,
		CargoTransferDialogEvent,
		ChangeMassDriverSpeedEvent,
		ChangeWaypointEvent,
		ChangeWaypointTransportTasksEvent,
		MergeFleetsDialogEvent,
		MergeFleetsEvent,
		RenameFleetEvent,
		SelectWaypointEvent,
		SplitFleetDialogEvent,
		SplitFleetEvent,
		TransferCargoEvent,
		TransportTasksDialogEvent
	} from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import { absoluteSize } from '$lib/types/CargoTransferRequest.svelte';
	import { MapObjectTypePlanet, None, type MapObject } from '$lib/types/cs';
	import { type WaypointDest } from '$lib/types/Fleet';
	import { equal as mapObjectEqual, ownedBy } from '$lib/types/MapObject';
	import { newSalvage } from '$lib/types/Salvage';
	import { equal } from '$lib/types/Vector';
	import hotkeys from 'hotkeys-js';
	import { onMount } from 'svelte';
	import CargoTransferDialog from '../dialogs/cargo/CargoTransferDialog.svelte';
	import MergeFleetsDialog from '../dialogs/merge/MergeFleetsDialog.svelte';
	import ProductionQueueDialog from '../dialogs/production/ProductionQueueDialog.svelte';
	import SplitFleetDialog from '../dialogs/split/SplitFleetDialog.svelte';
	import TransportTasksDialog from '../dialogs/transport/TransportTasksDialog.svelte';
	import SearchDialog from '../search/SearchDialog.svelte';
	import CommandPane from './command/CommandPane.svelte';
	import CommandPaneCarousel from './command/CommandPaneCarousel.svelte';
	import MapObjectStatsBar from './MapObjectStatsBar.svelte';
	import MapObjectSummary from './MapObjectSummary.svelte';
	import Scanner from './scanner/Scanner.svelte';
	import ScannerToolbar from './scanner/ScannerToolbar.svelte';

	const {
		settings,
		universe,
		player,
		selectedMapObject,
		commandedMapObject,
		commandedPlanet,
		commandedFleet,
		commandMapObject,
		zoomToMapObject,
		nextMapObject,
		previousMapObject,
		selectWaypoint,
		selectMapObject,
		addWaypoint,
		updateWaypoint,
		deleteWaypoint,
		renameFleet,
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

	function onNextMapObject() {
		nextMapObject();
	}

	function onPreviousMapObject() {
		previousMapObject();
	}

	async function onAddWaypoint(dest: WaypointDest, fastestWaypoint: boolean): Promise<boolean> {
		return addWaypoint(dest, fastestWaypoint);
	}

	async function onUpdateWaypointDest(dest: WaypointDest, fastestWaypoint: boolean, done: boolean) {
		updateWaypoint(dest, fastestWaypoint, done);
	}

	async function onDeleteWaypoint() {
		deleteWaypoint();
	}

	function onSelectWaypoint(e: SelectWaypointEvent) {
		const wp = e.waypoint;
		selectWaypoint(wp);

		if (wp.targetType && wp.targetNum) {
			const mo = $universe.getMapObject(wp);
			if (mo) {
				selectMapObject(mo);
			}
		}
	}

	async function onChangeWaypoint(e: ChangeWaypointEvent) {
		e.fleet.waypoints[e.waypointIndex] = e.waypoint;
		updateFleetOrders(e.fleet);
	}

	async function onUpdateTransportTasks(e: ChangeWaypointTransportTasksEvent) {
		// update the transport tasks for this waypoint and update it
		e.waypoint.transportTasks = e.transportTasks;
		await onChangeWaypoint(e);

		// close the dialog
		showTransportTasksDialog = false;
	}

	async function onRenameFleet(e: RenameFleetEvent) {
		renameFleet(e.fleet, e.name);
	}

	async function onBattlePlanChanged(e: BattlePlanChangedEvent) {
		updateFleetOrders(e.fleet);
	}

	async function onChangeMassDriverSpeed(e: ChangeMassDriverSpeedEvent) {
		updatePlanetOrders(e.planet);
	}

	async function onSplitAll() {
		if (!$commandedFleet) {
			return;
		}
		splitAll($commandedFleet);
	}

	async function onTransferCargo(e: TransferCargoEvent) {
		// close the dialog
		showCargoTransferDialog = false;

		if (e && absoluteSize(e.transferAmount) > 0) {
			if (!e.dest) {
				e.dest = newSalvage();
			}
			await transferCargo(e.src, e.dest, e.transferAmount);
		}
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

	function onSelectMapObject(mo: MapObject) {
		if ($selectedMapObject !== mo) {
			// we selected a different object, so just select it
			selectMapObject(mo);

			// if we selected a mapobject that is a waypoint, select the waypoint as well
			if ($commandedFleet?.waypoints) {
				const fleetWaypoint = $commandedFleet.waypoints.find((wp) =>
					equal(wp.position, mo.position)
				);
				if (fleetWaypoint) {
					selectWaypoint(fleetWaypoint);
				}
			}
		} else {
			// we selected the same mapobject twice
			const myMapObjectsAtPosition = $universe.getMyMapObjectsByPosition(mo);
			if (myMapObjectsAtPosition?.length > 0) {
				let index = myMapObjectsAtPosition.findIndex((mo) =>
					mapObjectEqual(mo, $commandedMapObject)
				);
				// if our currently commanded map object is not at this location, reset the index
				if (index == -1) {
					index = 0;
				} else {
					// command the next one
					index = index >= myMapObjectsAtPosition.length - 1 ? 0 : index + 1;
				}
				const nextMapObject = myMapObjectsAtPosition[index];

				commandMapObject(nextMapObject);
			}
		}
	}

	function onSetPacketDest(mo: MapObject) {
		if (mo.type != MapObjectTypePlanet) {
			return;
		} else {
			$settings.setPacketDest = false;
			// something went wrong, can't set dest on a planet without a massdriver
			if (!$commandedPlanet?.spec.hasMassDriver) {
				return;
			}

			if (mapObjectEqual(mo, $commandedPlanet)) {
				// clear dest
				$commandedPlanet.packetTargetNum = None;
			} else {
				$commandedPlanet.packetTargetNum = mo.num;
			}

			updatePlanetOrders($commandedPlanet);
		}
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
				{onNextMapObject}
				{onPreviousMapObject}
				{onRenameFleet}
				{onSelectWaypoint}
				{onChangeWaypoint}
				{onDeleteWaypoint}
				{onSplitAll}
				{onBattlePlanChanged}
				{onChangeMassDriverSpeed}
				onShowProductionQueueDialog={() => (showProductionQueueDialog = true)}
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
			<ScannerToolbar onShowSearch={() => (showSearchDialog = true)} />
			<Scanner
				{onSelectWaypoint}
				{onAddWaypoint}
				{onUpdateWaypointDest}
				{onSelectMapObject}
				{onSetPacketDest}
			/>
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
			{onSelectWaypoint}
			{onChangeWaypoint}
			{onDeleteWaypoint}
			{onSplitAll}
			{onBattlePlanChanged}
			{onChangeMassDriverSpeed}
			onShowProductionQueueDialog={() => (showProductionQueueDialog = true)}
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
<CargoTransferDialog
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
	show={showTransportTasksDialog}
	props={transportTasksDialogEvent}
	onOk={onUpdateTransportTasks}
	onCancel={() => (showTransportTasksDialog = false)}
/>
<SearchDialog
	show={showSearchDialog}
	onOk={(e) => selectSearchResult(e)}
	onCancel={() => (showSearchDialog = false)}
/>
