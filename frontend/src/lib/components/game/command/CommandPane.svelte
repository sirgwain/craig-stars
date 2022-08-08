<script lang="ts">
	import PlanetMineralsOnHandTile from './PlanetMineralsOnHandTile.svelte';
	import PlanetProductionTile from './PlanetProductionTile.svelte';
	import PlanetStatusTile from './PlanetStatusTile.svelte';
	import { game, player } from '$lib/services/Context';
	import PlanetStarbaseTile from './PlanetStarbaseTile.svelte';
	import PlanetSummaryTile from './PlanetSummaryTile.svelte';
	import PlanetFleetsInOrbitTile from './PlanetFleetsInOrbitTile.svelte';
	import FleetSummaryTile from './FleetSummaryTile.svelte';
	import FleetFuelAndCargoTile from './FleetFuelAndCargoTile.svelte';
import FleetOtherFleetsHereTile from './FleetOtherFleetsHereTile.svelte';
</script>

{#if $game && $player}
	<PlanetSummaryTile />
	<PlanetProductionTile />
	<PlanetMineralsOnHandTile />
	<PlanetStarbaseTile />
	<PlanetStatusTile />
	<PlanetFleetsInOrbitTile />
	<FleetSummaryTile />
	<FleetFuelAndCargoTile />
	<FleetOtherFleetsHereTile />
{/if}
