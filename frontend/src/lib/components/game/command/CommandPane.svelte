<script lang="ts">
	import MineralsOnHandTile from './MineralsOnHandTile.svelte';
	import PlanetProductionTile from './PlanetProductionTile.svelte';
	import PlanetStatusTile from './PlanetStatusTile.svelte';
	import { game, player } from '$lib/services/Context';
	import PlanetSummaryTile from './PlanetSummaryTile.svelte';
	import FleetSummaryTile from './FleetSummaryTile.svelte';
</script>

{#if $game && $player}
		<PlanetSummaryTile />
		<MineralsOnHandTile />
		<PlanetStatusTile />
		<PlanetProductionTile />
		<FleetSummaryTile />
{/if}
