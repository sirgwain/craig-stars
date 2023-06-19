<script lang="ts">
	import { commandedFleet, commandedPlanet, designs } from '$lib/services/Context';
	import type { Player } from '$lib/types/Player';
	import FleetFuelAndCargoTile from './FleetFuelAndCargoTile.svelte';
	import FleetOtherFleetsHereTile from './FleetOtherFleetsHereTile.svelte';
	import FleetSummaryTile from './FleetSummaryTile.svelte';
	import FleetWaypointsTile from './FleetWaypointsTile.svelte';
	import FleetWaypointTaskTile from './FleetWaypointTaskTile.svelte';
	import PlanetFleetsInOrbitTile from './PlanetFleetsInOrbitTile.svelte';
	import PlanetMineralsOnHandTile from './PlanetMineralsOnHandTile.svelte';
	import PlanetProductionTile from './PlanetProductionTile.svelte';
	import PlanetStarbaseTile from './PlanetStarbaseTile.svelte';
	import PlanetStatusTile from './PlanetStatusTile.svelte';
	import PlanetSummaryTile from './PlanetSummaryTile.svelte';

	export let player: Player;
</script>

{#if $commandedPlanet}
	<PlanetSummaryTile planet={$commandedPlanet} />
	<PlanetProductionTile planet={$commandedPlanet} />
	<PlanetMineralsOnHandTile planet={$commandedPlanet} />
	<PlanetStarbaseTile {player} planet={$commandedPlanet} />
	<PlanetStatusTile planet={$commandedPlanet} />
	<PlanetFleetsInOrbitTile planet={$commandedPlanet} />
{:else if $commandedFleet}
	<FleetSummaryTile fleet={$commandedFleet} designs={$designs ?? []} />
	<FleetFuelAndCargoTile fleet={$commandedFleet} />
	<FleetOtherFleetsHereTile fleet={$commandedFleet} />
	<FleetWaypointsTile fleet={$commandedFleet} />
	<FleetWaypointTaskTile fleet={$commandedFleet} />
	<!-- empty div for layout -->
	<div class="hidden md:block md:w-[14rem]" />
{:else}{/if}