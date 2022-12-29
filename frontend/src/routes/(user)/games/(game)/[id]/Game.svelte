<script lang="ts">
	import { settings } from '$lib/services/Context';
	import type { Game } from '$lib/types/Game';
	import type { Player } from '$lib/types/Player';
	import CommandPane from './command/CommandPane.svelte';
	import HighlightedMapObjectStats from './HighlightedMapObjectStats.svelte';
	import MapObjectSummary from './MapObjectSummary.svelte';
	import MessagesPane from './MessagesPane.svelte';
	import Scanner from './scanner/Scanner.svelte';
	import ScannerToolbar from './scanner/ScannerToolbar.svelte';

	export let game: Game;
	export let player: Player;
</script>

<div class="flex flex-row h-full w-full">
	<div class="hidden sm:block flex-initial w-[16rem] lg:w-[30rem] overflow-y-auto">
		<div class="flex flex-col">
			<div class="flex-1">
				<div class="flex flex-row flex-wrap gap-3.5 justify-between">
					<CommandPane />
				</div>
			</div>
			<div class="hidden mt-3 lg:block flex-initial"><MapObjectSummary /></div>
		</div>
	</div>

	<div class="flex-1 flex flex-col grow ml-3">
		<div class="flex flex-col border-gray-700 border-2 shadow-sm grow min-h-[68%]">
			<ScannerToolbar {game} {player} />
			<Scanner />
		</div>
		<HighlightedMapObjectStats />
		<div class="lg:hidden mt-3 block flex-initial"><MapObjectSummary /></div>
		<div class="md:hidden mt-3 block flex-initial"><CommandPane /></div>
	</div>
</div>
