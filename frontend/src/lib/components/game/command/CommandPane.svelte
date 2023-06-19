<script lang="ts">
	import MineralsOnHandTile from './MineralsOnHandTile.svelte';
	import PlanetProductionTile from './PlanetProductionTile.svelte';
	import PlanetStatusTile from './PlanetStatusTile.svelte';
	import { game, player } from '$lib/services/Context';
	import PlanetSummaryTile from './PlanetSummaryTile.svelte';
	import MessagesPane from '../MessagesPane.svelte';
</script>

{#if $game && $player.planets?.length}
	<div class="flex flex-col h-full max-w-[14rem] lg:max-w-[33rem]">
		<div class="flex-1">
			<div class="flex flex-row">
				<div class="flex grow flex-wrap gap-2">
					<PlanetSummaryTile />
					<MineralsOnHandTile />
					<PlanetStatusTile />
					<PlanetProductionTile />
				</div>
			</div>
		</div>

		<div class="flex-none w-full"><MessagesPane /></div>
		<div class="flex-none w-full">summary</div>
	</div>
{/if}
