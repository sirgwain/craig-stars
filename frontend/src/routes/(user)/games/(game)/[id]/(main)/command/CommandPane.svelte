<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import { commandedFleet, commandedPlanet } from '$lib/services/Stores';
	import FleetCompositionTile from './FleetCompositionTile.svelte';
	import FleetFuelAndCargoTile from './FleetFuelAndCargoTile.svelte';
	import FleetOrbitingTile from './FleetOrbitingTile.svelte';
	import FleetOtherFleetsHereTile from './FleetOtherFleetsHereTile.svelte';
	import FleetSummaryTile from './FleetSummaryTile.svelte';
	import FleetWaypointTaskTile from './FleetWaypointTaskTile.svelte';
	import FleetWaypointsTile from './FleetWaypointsTile.svelte';
	import PlanetFleetsInOrbitTile from './PlanetFleetsInOrbitTile.svelte';
	import PlanetMineralsOnHandTile from './PlanetMineralsOnHandTile.svelte';
	import PlanetProductionTile from './PlanetProductionTile.svelte';
	import PlanetStarbaseTile from './PlanetStarbaseTile.svelte';
	import PlanetStatusTile from './PlanetStatusTile.svelte';
	import PlanetSummaryTile from './PlanetSummaryTile.svelte';

	const { game, player, universe } = getGameContext();
</script>

{#if $commandedPlanet}
	<div class="lg:flex lg:flex-col">
		<PlanetSummaryTile planet={$commandedPlanet} />
		<PlanetMineralsOnHandTile planet={$commandedPlanet} />
		<PlanetStatusTile planet={$commandedPlanet} />
	</div>
	<div class="lg:flex lg:flex-col">
		<PlanetFleetsInOrbitTile
			planet={$commandedPlanet}
			fleetsInOrbit={$universe.getMyFleetsByPosition($commandedPlanet)}
		/>
		<PlanetProductionTile planet={$commandedPlanet} />
		<PlanetStarbaseTile
			planet={$commandedPlanet}
			starbase={$universe.getPlanetStarbase($commandedPlanet.num)}
		/>
	</div>
{:else if $commandedFleet}
	<div class="lg:flex lg:flex-col">
		<FleetSummaryTile fleet={$commandedFleet} />
		<FleetOrbitingTile fleet={$commandedFleet} />
		<FleetOtherFleetsHereTile
			fleet={$commandedFleet}
			fleetsInOrbit={$universe
				.getMyFleetsByPosition($commandedFleet)
				.filter((f) => f.num !== $commandedFleet?.num)}
		/>
		<FleetCompositionTile
			fleet={$commandedFleet}
			on:splitAll={() => $commandedFleet && $game.splitAll($commandedFleet)}
		/>
	</div>
	<div class="lg:flex lg:flex-col">
		<FleetFuelAndCargoTile fleet={$commandedFleet} />
		<FleetWaypointsTile fleet={$commandedFleet} />
		<FleetWaypointTaskTile fleet={$commandedFleet} />
	</div>
{/if}
