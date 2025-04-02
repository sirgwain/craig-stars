<script lang="ts">
	import type {
		BattlePlanChangedEvent,
		CargoTransferDialogEvent,
		ChangeMassDriverSpeedEvent,
		ChangeWaypointEvent,
		ChangeWaypointTransportTasksEvent,
		ClearProductionQueueEvent,
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
	import { commandable, equal as mapObjectEqual, ownedBy } from '$lib/types/MapObject';
	import { equal } from '$lib/types/Vector';
	import hotkeys from 'hotkeys-js';
	import { onMount } from 'svelte';
	import CargoTranfserDialog from '../dialogs/cargo/CargoTransferDialog.svelte';
	import MergeFleetsDialog from '../dialogs/merge/MergeFleetsDialog.svelte';
	import ProductionQueueDialog from '../dialogs/production/ProductionQueueDialog.svelte';
	import SplitFleetDialog from '../dialogs/split/SplitFleetDialog.svelte';
	import TransportTasksDialog from '../dialogs/transport/TransportTasksDialog.svelte';
	import SearchDialog from '../search/SearchDialog.svelte';
	import CommandPane from './command/CommandPane.svelte';
	import CommandPaneDrawer from './command/CommandPaneDrawer.svelte';
	import MapObjectStatsBar from './MapObjectStatsBar.svelte';
	import MapObjectSummary from './MapObjectSummary.svelte';
	import Scanner from './scanner/Scanner.svelte';
	import ScannerToolbar from './scanner/ScannerToolbar.svelte';

	const {
		game,
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
		selectNextMapObject,
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

	async function onClearProductionQueue(e: ClearProductionQueueEvent) {
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

	// onSelectMapObject cycles through commanding MapObjects at a location, and then
	// selecting non-commandable mapobjects
	function onSelectMapObject(mo: MapObject) {
		// nothing selected, select this
		if (!$selectedMapObject) {
			console.log('nothing selected, selecting ', mo.name);
			selectMapObject(mo);
			return;
		}

		// selected something at a new location, select it
		if (!equal($selectedMapObject.position, mo.position)) {
			console.log('new location selected, selecting ', mo.name);
			selectMapObject(mo);
			return;
		}

		// get all the mapobjects here we want to cycle, starting with commandable map objects
		const commandableObjects = $universe.getCommandableMapObjectsByPosition(mo.position);
		const selectableObjects = $universe
			.getMapObjectsByPosition(mo.position)
			.filter((mo) => !commandable($player.num, mo));
		let commandedIndex = commandableObjects.findIndex((mo) =>
			mapObjectEqual(mo, $commandedMapObject)
		);
		let selectedIndex = selectableObjects.findIndex((mo) => mapObjectEqual(mo, $selectedMapObject));

		if (commandedIndex === -1 && commandableObjects.length > 0) {
			// we haven't commanded anything here yet, command the first one
			console.log('nothing commanded, commanding ', commandableObjects[0].name);
			commandMapObject(commandableObjects[0]);
		} else if (commandedIndex !== -1 && commandedIndex < commandableObjects.length - 1) {
			// command the next map object in the cycle
			console.log('commanding next ', commandableObjects[commandedIndex + 1].name);
			commandMapObject(commandableObjects[commandedIndex + 1]);
		} else if (selectedIndex === -1) {
			if (selectableObjects.length > 0) {
				console.log('nothing selected, selecting first ', selectableObjects[0].name);
				selectMapObject(selectableObjects[0]);
			} else {
				console.log(
					'nothing selected, nothing selectable selecting first commandable ',
					commandableObjects[0].name
				);
				selectMapObject(commandableObjects[0]);
				commandMapObject(commandableObjects[0]);
			}
		} else if (selectedIndex !== -1 && selectedIndex < selectableObjects.length - 1) {
			// Cycle to the next selectable object
			console.log('cycling selection to ', selectableObjects[selectedIndex + 1].name);
			selectMapObject(selectableObjects[selectedIndex + 1]);
		} else if (selectableObjects.length > 0) {
			// If at the end, wrap around to the first selectable object
			console.log('looping back to first selectable ', selectableObjects[0].name);
			if (commandableObjects.length > 0) {
				selectMapObject(commandableObjects[0]);
				commandMapObject(commandableObjects[0]);
			} else {
				selectMapObject(selectableObjects[0]);
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
<div class="flex flex-col h-full md:flex-row" data-type="game-view" data-id={$game.id}>
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
				{onClearProductionQueue}
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
			<ScannerToolbar
				onShowSearch={() => (showSearchDialog = true)}
				onCycleMapObject={() => selectNextMapObject()}
				{onNextMapObject}
				{onPreviousMapObject}
			/>
			<Scanner
				{onSelectWaypoint}
				{onAddWaypoint}
				{onUpdateWaypointDest}
				{onSelectMapObject}
				{onSetPacketDest}
			/>
		</div>
		<div class="hidden md:block">
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

	<!-- for phone displays, use a drawer -->
	<div class="flex flex-col">
		<CommandPaneDrawer
			{onNextMapObject}
			{onPreviousMapObject}
			{onRenameFleet}
			{onSelectWaypoint}
			{onChangeWaypoint}
			{onDeleteWaypoint}
			{onSplitAll}
			{onBattlePlanChanged}
			{onChangeMassDriverSpeed}
			{onClearProductionQueue}
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
