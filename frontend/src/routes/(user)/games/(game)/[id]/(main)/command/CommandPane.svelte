<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import { commandedFleet, commandedPlanet } from '$lib/services/Stores';
	import { createEventDispatcher } from 'svelte';
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

	const dispatch = createEventDispatcher();
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
			on:cargo-transfer-dialog={(e) => dispatch('cargo-transfer-dialog', e.detail)}
		/>
		<PlanetProductionTile
			planet={$commandedPlanet}
			on:change-production={(e) => dispatch('change-production', e.detail)}
		/>
		<PlanetStarbaseTile
			planet={$commandedPlanet}
			starbase={$universe.getPlanetStarbase($commandedPlanet.num)}
		/>
	</div>
{:else if $commandedFleet}
	<div class="lg:flex lg:flex-col">
		<FleetSummaryTile fleet={$commandedFleet} />
		<FleetOrbitingTile
			fleet={$commandedFleet}
			on:cargo-transfer-dialog={(e) => dispatch('cargo-transfer-dialog', e.detail)}
		/>
		<FleetOtherFleetsHereTile
			fleet={$commandedFleet}
			fleetsInOrbit={$universe
				.getMyFleetsByPosition($commandedFleet)
				.filter((f) => f.num !== $commandedFleet?.num)}
			on:cargo-transfer-dialog={(e) => dispatch('cargo-transfer-dialog', e.detail)}
		/>
		<FleetCompositionTile
			fleet={$commandedFleet}
			on:split-all={() => $commandedFleet && $game.splitAll($commandedFleet)}
			on:merge-fleets-dialog={(e) => dispatch('merge-fleets-dialog', e.detail)}
		/>
	</div>
	<div class="lg:flex lg:flex-col">
		<FleetFuelAndCargoTile
			fleet={$commandedFleet}
			on:cargo-transfer-dialog={(e) => dispatch('cargo-transfer-dialog', e.detail)}
		/>
		<FleetWaypointsTile fleet={$commandedFleet} />
		<FleetWaypointTaskTile fleet={$commandedFleet} />
	</div>
{/if}
