<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { createEventDispatcher } from 'svelte';
	import FleetSummary from '../FleetSummary.svelte';
	import MapObjectSummary from '../MapObjectSummary.svelte';
	import PlanetSummary from '../PlanetSummary.svelte';
	import FleetCompositionTile from './FleetCompositionTile.svelte';
	import FleetFuelAndCargoTile from './FleetFuelAndCargoTile.svelte';
	import FleetOrbitingTile from './FleetOrbitingTile.svelte';
	import FleetOtherFleetsHereTile from './FleetOtherFleetsHereTile.svelte';
	import FleetWaypointTaskTile from './FleetWaypointTaskTile.svelte';
	import FleetWaypointsTile from './FleetWaypointsTile.svelte';
	import PlanetFleetsInOrbitTile from './PlanetFleetsInOrbitTile.svelte';
	import PlanetMineralsOnHandTile from './PlanetMineralsOnHandTile.svelte';
	import PlanetProductionTile from './PlanetProductionTile.svelte';
	import PlanetStarbaseTile from './PlanetStarbaseTile.svelte';
	import PlanetStatusTile from './PlanetStatusTile.svelte';
	import type { CargoTransferDialogEvent } from '../../dialogs/cargo/CargoTranfserDialog.svelte';
	import type { MergeFleetsDialogEvent } from '../../dialogs/merge/MergeFleetsDialog.svelte';
	import type { ProductionQueueDialogEvent } from '../../dialogs/production/ProductionQueueDialog.svelte';
	import type { SplitFleetDialogEvent } from '../../dialogs/split/SplitFleetDialog.svelte';
	import type { TransportTasksDialogEvent } from '../../dialogs/transport/TransportTasksDialog.svelte';

	const dispatch = createEventDispatcher<
		SplitFleetDialogEvent &
			MergeFleetsDialogEvent &
			CargoTransferDialogEvent &
			ProductionQueueDialogEvent &
			TransportTasksDialogEvent
	>();
	const { game, universe, commandedFleet, commandedPlanet, splitAll } = getGameContext();
</script>

{#if $commandedPlanet}
	<div class="carousel w-full md:hidden">
		<div id="planet-summary-tile" class="carousel-item w-full">
			<div class="w-full card bg-base-200 shadow rounded-sm border-2 border-base-300">
				<div class="card-body p-2 gap-0">
					<div class="flex flex-row items-center">
						<div class="flex-1 text-center text-lg font-semibold text-secondary">
							{$commandedPlanet?.name ?? ''}
						</div>
					</div>
					<PlanetSummary planet={$commandedPlanet} />
				</div>
			</div>
		</div>
		<div id="mapobject-summary" class="carousel-item w-full">
			<MapObjectSummary />
		</div>
		<div id="planet-production-tile" class="carousel-item w-full">
			<PlanetProductionTile
				planet={$commandedPlanet}
				on:change-production={(e) => dispatch('change-production', e.detail)}
			/>
		</div>
		<div id="planet-minerals-on-hand-tile" class="carousel-item w-full">
			<PlanetMineralsOnHandTile planet={$commandedPlanet} />
		</div>
		<div id="planet-starbase-tile" class="carousel-item w-full">
			<PlanetStarbaseTile
				planet={$commandedPlanet}
				starbase={$universe.getPlanetStarbase($commandedPlanet.num)}
			/>
		</div>
		<div id="planet-status-tile" class="carousel-item w-full">
			<PlanetStatusTile planet={$commandedPlanet} />
		</div>
		<div id="planet-fleets-in-orbit-tile" class="carousel-item w-full">
			<PlanetFleetsInOrbitTile
				planet={$commandedPlanet}
				fleetsInOrbit={$universe.getMyFleetsByPosition($commandedPlanet)}
				on:cargo-transfer-dialog={(e) => dispatch('cargo-transfer-dialog', e.detail)}
			/>
		</div>
	</div>
	<div class="flex justify-center w-full md:hidden py-2 gap-2">
		<a href="#planet-summary-tile" class="btn btn-xs">1</a>
		<a href="#mapobject-summary" class="btn btn-xs">2</a>
		<a href="#planet-production-tile" class="btn btn-xs">3</a>
		<a href="#planet-minerals-on-hand-tile" class="btn btn-xs">4</a>
		<a href="#planet-starbase-tile" class="btn btn-xs">5</a>
		<a href="#planet-status-tile" class="btn btn-xs">6</a>
		<a href="#planet-fleets-in-orbit-tile" class="btn btn-xs">7</a>
	</div>
{:else if $commandedFleet}
	<div class="carousel w-full md:hidden">
		<div id="fleet-summary-tile" class="carousel-item w-full">
			<div class="w-full card bg-base-200 shadow rounded-sm border-2 border-base-300">
				<div class="card-body p-2 gap-0">
					<div class="flex flex-row items-center">
						<div class="flex-1 text-center text-lg font-semibold text-secondary">
							{$commandedFleet?.name ?? ''}
						</div>
					</div>
					<FleetSummary fleet={$commandedFleet} />
				</div>
			</div>
		</div>
		<div id="mapobject-summary" class="carousel-item w-full">
			<MapObjectSummary />
		</div>
		<div id="fleet-orbiting-tile" class="carousel-item w-full">
			<FleetOrbitingTile
				fleet={$commandedFleet}
				on:cargo-transfer-dialog={(e) => dispatch('cargo-transfer-dialog', e.detail)}
			/>
		</div>
		<div id="fleet-fuel-and-cargo-tile" class="carousel-item w-full">
			<FleetFuelAndCargoTile
				fleet={$commandedFleet}
				on:cargo-transfer-dialog={(e) => dispatch('cargo-transfer-dialog', e.detail)}
			/>
		</div>
		<div id="fleet-waypoints-tile" class="carousel-item w-full">
			<FleetWaypointsTile fleet={$commandedFleet} />
		</div>
		<div id="fleet-composition-tile" class="carousel-item w-full">
			<FleetCompositionTile
				fleet={$commandedFleet}
				on:split-all={() => $commandedFleet && splitAll($commandedFleet)}
				on:split-fleet-dialog={(e) => dispatch('split-fleet-dialog', e.detail)}
				on:merge-fleets-dialog={(e) => dispatch('merge-fleets-dialog', e.detail)}
			/>
		</div>
		<div id="fleet-waypoint-task-tile" class="carousel-item w-full">
			<FleetWaypointTaskTile
				fleet={$commandedFleet}
				on:transport-tasks-dialog={(e) => dispatch('transport-tasks-dialog', e.detail)}
			/>
		</div>
		<div id="fleet-other-fleets-here-tile" class="carousel-item w-full">
			<FleetOtherFleetsHereTile
				fleet={$commandedFleet}
				fleetsInOrbit={$universe
					.getMyFleetsByPosition($commandedFleet)
					.filter((f) => f.num !== $commandedFleet?.num)}
				on:cargo-transfer-dialog={(e) => dispatch('cargo-transfer-dialog', e.detail)}
			/>
		</div>
	</div>
	<div class="flex justify-center w-full md:hidden py-2 gap-2">
		<a href="#fleet-summary-tile" class="btn btn-xs">1</a>
		<a href="#mapobject-summary" class="btn btn-xs">2</a>
		<a href="#fleet-orbiting-tile" class="btn btn-xs">3</a>
		<a href="#fleet-fuel-and-cargo-tile" class="btn btn-xs">4</a>
		<a href="#fleet-waypoints-tile" class="btn btn-xs">5</a>
		<a href="#fleet-composition-tile" class="btn btn-xs">6</a>
		<a href="#fleet-waypoint-task-tile" class="btn btn-xs">7</a>
		<a href="#fleet-other-fleets-here-tile" class="btn btn-xs">8</a>
	</div>
{/if}
