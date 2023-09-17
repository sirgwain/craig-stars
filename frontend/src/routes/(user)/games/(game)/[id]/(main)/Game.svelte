<script lang="ts">
	import { commandedPlanet } from '$lib/services/Stores';
	import hotkeys from 'hotkeys-js';
	import { onMount } from 'svelte';
	import CargoTranfserDialog, {
		type CargoTransferDialogEventDetails
	} from '../dialogs/cargo/CargoTranfserDialog.svelte';
	import MergeFleetsDialog, {
		type MergeFleetsDialogEventDetails
	} from '../dialogs/merge/MergeFleetsDialog.svelte';
	import ProductionQueueDialog from '../dialogs/production/ProductionQueueDialog.svelte';
	import HighlightedMapObjectStats from './HighlightedMapObjectStats.svelte';
	import MapObjectSummary from './MapObjectSummary.svelte';
	import CommandPane from './command/CommandPane.svelte';
	import CommandPaneCarousel from './command/CommandPaneCarousel.svelte';
	import Scanner from './scanner/Scanner.svelte';
	import ScannerToolbar from './scanner/ScannerToolbar.svelte';
	import type { TransportTasksDialogEventDetails } from '../dialogs/transport/TransportTasksDialog.svelte';
	import TransportTasksDialog from '../dialogs/transport/TransportTasksDialog.svelte';

	let showProductionQueueDialog = false;
	let showCargoTransferDialog = false;
	let showMergeFleetsDialog = false;
	let showTransportTasksDialog = false;
	let cargoTransferDetails: CargoTransferDialogEventDetails | undefined = undefined;
	let mergeFleetsDialogEventDetails: MergeFleetsDialogEventDetails | undefined = undefined;
	let transportTasksDialogEventDetails: TransportTasksDialogEventDetails | undefined = undefined;

	onMount(() => {
		hotkeys('q', 'root', () => {
			if ($commandedPlanet) {
				showProductionQueueDialog = true;
			}
		});
		return () => {
			hotkeys.unbind('q', 'root');
		};
	});
</script>

<!-- for small mobile displays we put the scanner on top and the command pane below it-->
<div class="flex flex-col h-full md:flex-row">
	<!-- for medium+ displays, command pane goes on the left -->
	<div
		class="hidden md:flex md:flex-col md:flex-none justify-between md:w-[15.5rem] lg:w-[30rem] overflow-y-auto md:max-h-[calc(100dvh-4rem)]"
	>
		<div class="flex flex-row flex-wrap gap-2 justify-center">
			<CommandPane
				on:change-production={(e) => (showProductionQueueDialog = true)}
				on:cargo-transfer-dialog={(e) => {
					showCargoTransferDialog = true;
					cargoTransferDetails = e.detail;
				}}
				on:merge-fleets-dialog={(e) => {
					showMergeFleetsDialog = true;
					mergeFleetsDialogEventDetails = e.detail;
				}}
				on:transport-tasks-dialog={(e) => {
					showTransportTasksDialog = true;
					transportTasksDialogEventDetails = e.detail;
				}}
			/>
		</div>
		<div class="hidden lg:block lg:p-1 mx-2">
			<MapObjectSummary />
		</div>
	</div>

	<div class="flex flex-col grow">
		<div class="flex flex-col grow border-gray-700 border-2 shadow-sm">
			<ScannerToolbar />
			<Scanner />
		</div>
		<div>
			<HighlightedMapObjectStats />
		</div>
		<div class="hidden md:block md:w-full lg:hidden mb-2">
			<MapObjectSummary />
		</div>
	</div>

	<!-- for phone displays, use a carousel -->
	<div class="flex flex-col flex-0">
		<CommandPaneCarousel
			on:change-production={(e) => (showProductionQueueDialog = true)}
			on:cargo-transfer-dialog={(e) => {
				showCargoTransferDialog = true;
				cargoTransferDetails = e.detail;
			}}
			on:merge-fleets-dialog={(e) => {
				showMergeFleetsDialog = true;
				mergeFleetsDialogEventDetails = e.detail;
			}}
			on:transport-tasks-dialog={(e) => {
				showTransportTasksDialog = true;
				transportTasksDialogEventDetails = e.detail;
			}}
		/>
	</div>
</div>

<!-- dialog modals -->
<ProductionQueueDialog bind:show={showProductionQueueDialog} />
<CargoTranfserDialog bind:show={showCargoTransferDialog} bind:props={cargoTransferDetails} />
<MergeFleetsDialog bind:show={showMergeFleetsDialog} bind:props={mergeFleetsDialogEventDetails} />
<TransportTasksDialog
	bind:show={showTransportTasksDialog}
	bind:props={transportTasksDialogEventDetails}
/>
