<script lang="ts">
	import CommandPane from './command/CommandPane.svelte';
	import HighlightedMapObjectStats from './HighlightedMapObjectStats.svelte';
	import MapObjectSummary from './MapObjectSummary.svelte';
	import MessagesPane from './MessagesPane.svelte';
	import Scanner from './scanner/Scanner.svelte';
</script>

<div class="flex flex-row h-full w-full">
	<div class="hidden sm:block flex-initial w-[16rem] lg:w-[30rem] overflow-y-auto">
		<div class="flex flex-col">
			<div class="flex-1">
				<div class="flex flex-row flex-wrap gap-3.5 justify-between">
					<CommandPane />
				</div>
			</div>
			<div class="hidden mt-3 lg:block flex-initial"><MessagesPane /></div>
			<div class="hidden mt-3 lg:block flex-initial"><MapObjectSummary /></div>
		</div>
	</div>

	<div class="flex-1 flex flex-col h-full ml-3">
		<div class="flex-1 border-gray-700 border-2 shadow-sm min-h-[50%]">
			<Scanner />
		</div>
		<HighlightedMapObjectStats />
		<div class="lg:hidden mt-3 block flex-initial"><MapObjectSummary /></div>
		<div class="lg:hidden mt-3 block flex-initial"><MessagesPane /></div>
		<div class="md:hidden mt-3 block flex-initial"><CommandPane /></div>
	</div>
</div>
