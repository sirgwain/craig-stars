<script lang="ts">
	import type { Game } from '$lib/types/Game';
	import type { Player, PlayerMapObjects } from '$lib/types/Player';
	import type { ShipDesign } from '$lib/types/ShipDesign';
	import CommandPane from './command/CommandPane.svelte';
	import HighlightedMapObjectStats from './HighlightedMapObjectStats.svelte';
	import MapObjectSummary from './MapObjectSummary.svelte';
	import Scanner from './scanner/Scanner.svelte';
	import ScannerToolbar from './scanner/ScannerToolbar.svelte';

	export let game: Game;
	export let player: Player;
	export let designs: ShipDesign[];
	export let mapObjects: PlayerMapObjects;
</script>

<!-- for small mobile displays we put the scanner on top and the command pane below it-->
<div class="flex flex-col h-full md:flex-row">
	<!-- for medium+ displays, command pane goes on the left -->
	<div
		class="hidden md:flex md:flex-col justify-between md:w-[14.5rem] lg:w-[29rem] max-h-full overflow-y-auto"
	>
		<div class="flex flex-row flex-wrap gap-2 justify-center place-items-stretch">
			<CommandPane {mapObjects} />
		</div>
		<div class="hidden lg:block lg:p-1 mb-2">
			<MapObjectSummary {player} {designs} />
		</div>
	</div>

	<div class="flex flex-col grow min-h-[515px] md:h-full">
		<div class="flex flex-col grow border-gray-700 border-2 shadow-sm">
			<ScannerToolbar {game} {player} />
			<Scanner {game} {player} {mapObjects} />
		</div>
		<div>
			<HighlightedMapObjectStats />
		</div>
		<div class="hidden md:block md:w-full lg:hidden mb-2">
			<MapObjectSummary {player} {designs} />
		</div>
	</div>

	<div class="carousel md:hidden">
		<div class="carousel-item">
			<div class="w-screen">
				<MapObjectSummary {player} {designs} />
			</div>
		</div>
		<div class="carousel-item">
			<CommandPane {mapObjects} />
		</div>
	</div>
</div>
