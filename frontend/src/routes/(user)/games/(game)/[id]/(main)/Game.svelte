<script lang="ts">
	import CommandPane from './command/CommandPane.svelte';
	import CommandPaneCarousel from './command/CommandPaneCarousel.svelte';
	import Dialogs from './Dialogs.svelte';
	import HighlightedMapObjectStats from './HighlightedMapObjectStats.svelte';
	import MapObjectSummary from './MapObjectSummary.svelte';
	import Scanner from './scanner/Scanner.svelte';
	import ScannerToolbar from './scanner/ScannerToolbar.svelte';
</script>

<!-- for small mobile displays we put the scanner on top and the command pane below it-->
<div class="flex flex-col h-full md:flex-row">
	<!-- for medium+ displays, command pane goes on the left -->
	<div
		class="hidden md:flex md:flex-col md:flex-none justify-between md:w-[15.5rem] lg:w-[30rem] overflow-y-auto md:max-h-[calc(100vh-4rem)]"
	>
		<div class="flex flex-row flex-wrap gap-2 justify-center">
			<CommandPane />
		</div>
		<div class="hidden lg:block lg:p-1">
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
		<CommandPaneCarousel />
	</div>
</div>
<!-- dialog modals -->
<Dialogs />
