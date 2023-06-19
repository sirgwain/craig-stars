<script lang="ts">
	import MineralsOnHandTile from './MineralsOnHandTile.svelte';
	import PlanetProductionTile from './PlanetProductionTile.svelte';
	import PlanetStatusTile from './PlanetStatusTile.svelte';
	import { game, player } from '$lib/services/Context';
	import PlanetSummaryTile from './PlanetSummaryTile.svelte';
	import MessagesPane from '../MessagesPane.svelte';
	import MapObjectSummary from '../MapObjectSummary.svelte';
</script>

{#if $game && $player.planets?.length}
	<div class="flex flex-col h-full">
		<div class="flex-1">
			<div class="flex flex-col h-full">
				<div class="flex-1">
					<div class="flex grow flex-wrap gap-3.5">
						<PlanetSummaryTile />
						<MineralsOnHandTile />
						<PlanetStatusTile />
						<PlanetProductionTile />
					</div>
				</div>
				<div class="flex-none"><MessagesPane /></div>
				<div class="flex-none"><MapObjectSummary /></div>
			</div>
		</div>
	</div>
{/if}
