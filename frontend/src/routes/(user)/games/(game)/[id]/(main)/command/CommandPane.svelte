<script lang="ts">
	import { commandedFleet, commandedPlanet, designs } from '$lib/services/Context';
	import type { Fleet } from '$lib/types/Fleet';
	import type { CommandedPlanet } from '$lib/types/Planet';
	import type { PlayerResponse, PlayerMapObjects } from '$lib/types/Player';
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

	export let mapObjects: PlayerMapObjects;

	const getPlanetStarbase = (planet: CommandedPlanet) =>
		mapObjects.starbases.find((sb) => sb.planetNum == planet.num);
</script>

{#if $commandedPlanet}
	<PlanetSummaryTile planet={$commandedPlanet} />
	<PlanetProductionTile planet={$commandedPlanet} />
	<PlanetMineralsOnHandTile planet={$commandedPlanet} />
	<PlanetStarbaseTile starbase={getPlanetStarbase($commandedPlanet)} />
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
