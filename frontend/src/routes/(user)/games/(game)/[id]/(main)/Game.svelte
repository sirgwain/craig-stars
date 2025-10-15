<script lang="ts">
	import { addError } from '$lib/services/Errors';
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
	import { MapObjectType, type WaypointDest } from '$lib/types/cs-proto';
	import {
		commandable,
		equal as mapObjectEqual,
		ownedBy,
		type MapObjectLike
	} from '$lib/types/MapObject';
	import { None } from '$lib/types/Consts';
	import { emptyVector, equal } from '$lib/types/Vector';
	import type { ConnectError } from '@connectrpc/connect';
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

		if (wp.mapObjectTarget?.targetType && wp.mapObjectTarget.targetNum) {
			const mo = $universe.getMapObject(wp.mapObjectTarget);
			if (mo) {
				selectMapObject(mo);
			}
		}
	}

	async function onChangeWaypoint(e: ChangeWaypointEvent) {
		e.fleet.fleetOrders.waypoints[e.waypointIndex] = e.waypoint;
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

		if (absoluteSize(e.transferAmount) > 0) {
			try {
				await transferCargo(e.src, e.dest, e.transferAmount);
			} catch (err) {
				addError(err as ConnectError);
			}
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
		try {
			await split(e.src, e.dest, e.srcTokens, e.destTokens, e.transferAmount);
		} catch (err) {
			addError(err as ConnectError);
			throw err;
		}

		// close the dialog
		showSplitFleetDialog = false;
	}

	function selectSearchResult(mo: MapObjectLike | undefined) {
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
	function onSelectMapObject(mo: MapObjectLike) {
		if (
			!$selectedMapObject ||
			!equal(
				$selectedMapObject.mapObject?.position ?? emptyVector(),
				mo.mapObject?.position ?? emptyVector()
			)
		) {
			// nothing selected, or nothing selected at this location yet, select this object
			selectMapObject(mo);
			return;
		}

		// get all the mapobjects here we want to cycle, starting with commandable map objects
		const commandables = $universe.getCommandableMapObjectsByPosition(
			mo.mapObject?.position ?? emptyVector()
		);
		const selectables = $universe
			.getMapObjectsByPosition(mo.mapObject?.position ?? emptyVector())
			.filter((mo) => !commandable($player.num, mo));
		let commandedIndex = commandables.findIndex((mo) => mapObjectEqual($commandedMapObject, mo));
		let selectedIndex = selectables.findIndex((mo) => mapObjectEqual($selectedMapObject, mo));

		if (commandables.length === 0) {
			return;
		}

		if (commandedIndex < commandables.length - 1) {
			// we either havne't commanded anything yet (commandedIndex=-1) or there is a commandable object to cycle to
			// if we are at the end of the commanded list, this will skip
			commandMapObject(commandables[commandedIndex + 1]);
		} else if (selectedIndex === -1) {
			if (selectables.length > 0) {
				// nothing selected, selecting first selectable
				selectMapObject(selectables[0]);
			} else {
				// nothing selected, nothing selectable, selecting first commandable
				selectMapObject(commandables[0]);
				commandMapObject(commandables[0]);
			}
		} else if (selectedIndex !== -1 && selectedIndex < selectables.length - 1) {
			// Cycle to the next selectable object
			selectMapObject(selectables[selectedIndex + 1]);
		} else if (selectables.length > 0) {
			// If at the end, wrap around to the first selectable object
			if (commandables.length > 0) {
				selectMapObject(commandables[0]);
				commandMapObject(commandables[0]);
			} else {
				selectMapObject(selectables[0]);
			}
		}
	}

	function onSetPacketDest(mo: MapObjectLike) {
		if (mo.mapObject?.type !== MapObjectType.PLANET) {
			return;
		} else {
			$settings.setPacketDest = false;
			// something went wrong, can't set dest on a planet without a massdriver
			if (!$commandedPlanet?.spec.planetStarbaseSpec?.hasMassDriver) {
				return;
			}

			if (mapObjectEqual(mo, $commandedPlanet)) {
				// clear dest
				$commandedPlanet.planetOrders.packetTargetNum = None;
			} else {
				$commandedPlanet.planetOrders.packetTargetNum = mo.mapObject.num;
			}

			updatePlanetOrders($commandedPlanet);
		}
	}

	function onSetRouteDest(mo: MapObjectLike) {
		if (!$commandedPlanet) {
			return;
		}
		if (mo.mapObject?.type !== MapObjectType.PLANET) {
			return;
		} else {
			$settings.setRouteDest = false;

			if (mapObjectEqual(mo, $commandedPlanet)) {
				// clear dest
				$commandedPlanet.planetOrders.routeTargetNum = None;
				$commandedPlanet.planetOrders.routeTargetPlayerNum = None;
				$commandedPlanet.planetOrders.routeTargetType = MapObjectType.UNSPECIFIED;
			} else {
				$commandedPlanet.planetOrders.routeTargetNum = mo.mapObject.num;
				$commandedPlanet.planetOrders.routeTargetPlayerNum = mo.mapObject.playerNum;
				$commandedPlanet.planetOrders.routeTargetType = mo.mapObject.type;
			}

			updatePlanetOrders($commandedPlanet);
		}
	}
</script>

<!-- for small mobile displays we put the scanner on top and the command pane below it-->
<div class="flex flex-col h-full md:flex-row select-none" data-type="game-view" data-id={$game.id}>
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
				{onSetRouteDest}
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

<div class="select-none">
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
</div>
