<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { ownedBy, type MapObject } from '$lib/types/MapObject';
	import hotkeys from 'hotkeys-js';
	import { onMount } from 'svelte';
	import CargoTranfserDialog, {
		type CargoTransferDialogEventDetails
	} from '../dialogs/cargo/CargoTranfserDialog.svelte';
	import MergeFleetsDialog, {
		type MergeFleetsDialogEventDetails
	} from '../dialogs/merge/MergeFleetsDialog.svelte';
	import ProductionQueueDialog from '../dialogs/production/ProductionQueueDialog.svelte';
	import SplitFleetDialog, {
		type SplitFleetDialogEventDetails
	} from '../dialogs/split/SplitFleetDialog.svelte';
	import type { TransportTasksDialogEventDetails } from '../dialogs/transport/TransportTasksDialog.svelte';
	import TransportTasksDialog from '../dialogs/transport/TransportTasksDialog.svelte';
	import SearchDialog from '../search/SearchDialog.svelte';
	import HighlightedMapObjectStats from './HighlightedMapObjectStats.svelte';
	import MapObjectSummary from './MapObjectSummary.svelte';
	import CommandPane from './command/CommandPane.svelte';
	import CommandPaneCarousel from './command/CommandPaneCarousel.svelte';
	import Scanner from './scanner/Scanner.svelte';
	import ScannerToolbar from './scanner/ScannerToolbar.svelte';

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
		updateFleetOrders
	} = getGameContext();

	let showProductionQueueDialog = false;
	let showCargoTransferDialog = false;
	let showMergeFleetsDialog = false;
	let showSplitFleetDialog = false;
	let showTransportTasksDialog = false;
	let showSearchDialog = false;
	let cargoTransferDetails: CargoTransferDialogEventDetails | undefined = undefined;
	let mergeFleetsDialogEventDetails: MergeFleetsDialogEventDetails | undefined = undefined;
	let splitFleetDialogEventDetails: SplitFleetDialogEventDetails | undefined = undefined;
	let transportTasksDialogEventDetails: TransportTasksDialogEventDetails | undefined = undefined;

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

		return () => {
			hotkeys.unbind('n', 'root');
			hotkeys.unbind('p', 'root');
			hotkeys.unbind('q', 'root');
			hotkeys.unbind('Delete', 'root');
			hotkeys.unbind('Backspace', 'root');
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

	function selectSearchResult(mo: MapObject | undefined) {
		if (mo) {
			if (ownedBy(mo, $player.num)) {
				commandMapObject(mo);
			}
			selectMapObject(mo);
			zoomToMapObject(mo);
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
				on:change-production={(e) => (showProductionQueueDialog = true)}
				on:cargo-transfer-dialog={(e) => {
					showCargoTransferDialog = true;
					cargoTransferDetails = e?.detail;
				}}
				on:merge-fleets-dialog={(e) => {
					showMergeFleetsDialog = true;
					mergeFleetsDialogEventDetails = e.detail;
				}}
				on:split-fleet-dialog={(e) => {
					showSplitFleetDialog = true;
					splitFleetDialogEventDetails = e.detail;
				}}
				on:transport-tasks-dialog={(e) => {
					showTransportTasksDialog = true;
					transportTasksDialogEventDetails = e.detail;
				}}
				on:delete-waypoint={onDeleteWaypoint}
			/>
		</div>
		<div class="hidden lg:block lg:p-1 mx-2">
			<MapObjectSummary
				on:cargo-transfer-dialog={(e) => {
					showCargoTransferDialog = true;
					cargoTransferDetails = e?.detail;
				}}
			/>
		</div>
	</div>

	<div class="flex flex-col grow">
		<div class="flex flex-col grow border-gray-700 border-2 shadow-sm">
			<ScannerToolbar on:show-search={() => (showSearchDialog = true)} />
			<Scanner />
		</div>
		<div class="hidden md:block">
			<HighlightedMapObjectStats />
		</div>
		<div class="hidden md:block md:w-full lg:hidden mb-2">
			<MapObjectSummary
				on:cargo-transfer-dialog={(e) => {
					showCargoTransferDialog = true;
					cargoTransferDetails = e?.detail;
				}}
			/>
		</div>
	</div>

	<!-- for phone displays, use a carousel -->
	<div class="flex flex-col flex-0">
		<CommandPaneCarousel
			on:change-production={(e) => (showProductionQueueDialog = true)}
			on:cargo-transfer-dialog={(e) => {
				showCargoTransferDialog = true;
				cargoTransferDetails = e?.detail;
			}}
			on:merge-fleets-dialog={(e) => {
				showMergeFleetsDialog = true;
				mergeFleetsDialogEventDetails = e.detail;
			}}
			on:split-fleet-dialog={(e) => {
				showSplitFleetDialog = true;
				splitFleetDialogEventDetails = e.detail;
			}}
			on:transport-tasks-dialog={(e) => {
				showTransportTasksDialog = true;
				transportTasksDialogEventDetails = e.detail;
			}}
			on:delete-waypoint={onDeleteWaypoint}
		/>
	</div>
</div>

<!-- dialog modals -->
<ProductionQueueDialog bind:show={showProductionQueueDialog} />
<CargoTranfserDialog bind:show={showCargoTransferDialog} bind:props={cargoTransferDetails} />
<MergeFleetsDialog bind:show={showMergeFleetsDialog} bind:props={mergeFleetsDialogEventDetails} />
<SplitFleetDialog bind:show={showSplitFleetDialog} bind:props={splitFleetDialogEventDetails} />
<TransportTasksDialog
	bind:show={showTransportTasksDialog}
	bind:props={transportTasksDialogEventDetails}
/>
<SearchDialog bind:show={showSearchDialog} on:select-result={(e) => selectSearchResult(e.detail)} />
