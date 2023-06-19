<script lang="ts">
	import { commandedFleet, commandedPlanet } from '$lib/services/Context';
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

	export let game: FullGame;
</script>

{#if $commandedPlanet}
	<PlanetSummaryTile planet={$commandedPlanet} />
	<PlanetProductionTile planet={$commandedPlanet} />
	<PlanetMineralsOnHandTile planet={$commandedPlanet} />
	<PlanetStarbaseTile starbase={game.universe.getPlanetStarbase($commandedPlanet.num)} />
	<PlanetStatusTile planet={$commandedPlanet} />
	<PlanetFleetsInOrbitTile
		planet={$commandedPlanet}
		fleetsInOrbit={game.universe.getMyFleetsByPosition($commandedPlanet)}
	/>
{:else if $commandedFleet}
	<FleetSummaryTile fleet={$commandedFleet} player={game.player} />
	<FleetFuelAndCargoTile {game} fleet={$commandedFleet} />
	<FleetCompositionTile
		fleet={$commandedFleet}
		player={game.player}
		on:splitAll={() => $commandedFleet && game.splitAll($commandedFleet)}
	/>
	<FleetOtherFleetsHereTile
		fleet={$commandedFleet}
		fleetsInOrbit={game.universe
			.getMyFleetsByPosition($commandedFleet)
			.filter((f) => f.num !== $commandedFleet?.num)}
	/>
	<FleetWaypointsTile {game} fleet={$commandedFleet} />
	<FleetWaypointTaskTile {game} player={game.player} fleet={$commandedFleet} />
	<!-- empty div for layout -->
	<div class="hidden md:block md:w-[14rem]" />
{:else}{/if}
