<script lang="ts">
	import { game,player } from '$lib/services/Context';
	import FleetFuelAndCargoTile from './FleetFuelAndCargoTile.svelte';
	import FleetOtherFleetsHereTile from './FleetOtherFleetsHereTile.svelte';
	import FleetSummaryTile from './FleetSummaryTile.svelte';
	import FleetWaypointTaskTile from './FleetWaypointTaskTile.svelte';
	import PlanetFleetsInOrbitTile from './PlanetFleetsInOrbitTile.svelte';
	import PlanetMineralsOnHandTile from './PlanetMineralsOnHandTile.svelte';
	import PlanetProductionTile from './PlanetProductionTile.svelte';
	import PlanetStarbaseTile from './PlanetStarbaseTile.svelte';
	import PlanetStatusTile from './PlanetStatusTile.svelte';
	import PlanetSummaryTile from './PlanetSummaryTile.svelte';
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
	<FleetWaypointTaskTile />
{/if}
