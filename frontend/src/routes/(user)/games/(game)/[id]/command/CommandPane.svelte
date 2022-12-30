<script lang="ts">
	import { commandedPlanet, commandedFleet } from '$lib/services/Context';
	import type { Game } from '$lib/types/Game';
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

	export let game: Game;
	export let player: Player;
</script>

{#if $commandedPlanet}
	<PlanetSummaryTile planet={$commandedPlanet} />
	<PlanetProductionTile planet={$commandedPlanet} />
	<PlanetMineralsOnHandTile planet={$commandedPlanet} />
	<PlanetStarbaseTile {player} planet={$commandedPlanet} />
	<PlanetStatusTile {player} planet={$commandedPlanet} />
	<PlanetFleetsInOrbitTile planet={$commandedPlanet} />
{:else if $commandedFleet}
	<FleetSummaryTile {player} fleet={$commandedFleet} />
	<FleetFuelAndCargoTile />
	<FleetOtherFleetsHereTile />
	<FleetWaypointsTile />
	<FleetWaypointTaskTile />
{:else}{/if}
