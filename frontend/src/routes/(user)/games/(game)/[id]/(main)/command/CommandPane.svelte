<script lang="ts">
	import { commandedFleet, commandedPlanet } from '$lib/services/Stores';
	import type { FullGame } from '$lib/services/FullGame';
	import FleetCompositionTile from './FleetCompositionTile.svelte';
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
	import { getGameContext } from '$lib/services/Contexts';
	import FleetOrbitingTile from './FleetOrbitingTile.svelte';

	const { game, player, universe } = getGameContext();
</script>

{#if $commandedPlanet}
	<PlanetSummaryTile planet={$commandedPlanet} />
	<PlanetProductionTile planet={$commandedPlanet} />
	<PlanetMineralsOnHandTile planet={$commandedPlanet} />
	<PlanetStarbaseTile
		planet={$commandedPlanet}
		starbase={$universe.getPlanetStarbase($commandedPlanet.num)}
	/>
	<PlanetStatusTile planet={$commandedPlanet} />
	<PlanetFleetsInOrbitTile
		planet={$commandedPlanet}
		fleetsInOrbit={$universe.getMyFleetsByPosition($commandedPlanet)}
	/>
{:else if $commandedFleet}
	<FleetSummaryTile fleet={$commandedFleet} />
	<FleetOrbitingTile fleet={$commandedFleet} />
	<FleetFuelAndCargoTile fleet={$commandedFleet} />
	<FleetWaypointsTile fleet={$commandedFleet} />
	<FleetCompositionTile
		fleet={$commandedFleet}
		on:splitAll={() => $commandedFleet && $game.splitAll($commandedFleet)}
	/>
	<FleetWaypointTaskTile fleet={$commandedFleet} />
	<FleetOtherFleetsHereTile
		fleet={$commandedFleet}
		fleetsInOrbit={$universe
			.getMyFleetsByPosition($commandedFleet)
			.filter((f) => f.num !== $commandedFleet?.num)}
	/>
	<!-- empty div for layout -->
	<div class="hidden md:block md:w-[14rem]" />
{/if}
