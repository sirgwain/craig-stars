<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { createEventDispatcher } from 'svelte';
	import type { CargoTransferDialogEvent } from '../../dialogs/cargo/CargoTranfserDialog.svelte';
	import type { MergeFleetsDialogEvent } from '../../dialogs/merge/MergeFleetsDialog.svelte';
	import type { ProductionQueueDialogEvent } from '../../dialogs/production/ProductionQueueDialog.svelte';
	import type { SplitFleetDialogEvent } from '../../dialogs/split/SplitFleetDialog.svelte';
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
	import type { TransportTasksDialogEvent } from '../../dialogs/transport/TransportTasksDialog.svelte';

	const dispatch = createEventDispatcher<
		SplitFleetDialogEvent &
			MergeFleetsDialogEvent &
			CargoTransferDialogEvent &
			ProductionQueueDialogEvent &
			TransportTasksDialogEvent
	>();

	const { universe, commandedPlanet, commandedFleet, splitAll } = getGameContext();

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
			on:split-fleet-dialog={(e) => dispatch('split-fleet-dialog', e.detail)}
			on:cargo-transfer-dialog={(e) => dispatch('cargo-transfer-dialog', e.detail)}
		/>
		<FleetCompositionTile
			fleet={$commandedFleet}
			on:split-all={() => $commandedFleet && splitAll($commandedFleet)}
			on:split-fleet-dialog={(e) => dispatch('split-fleet-dialog', e.detail)}
			on:merge-fleets-dialog={(e) => dispatch('merge-fleets-dialog', e.detail)}
		/>
	</div>
	<div class="lg:flex lg:flex-col">
		<FleetFuelAndCargoTile
			fleet={$commandedFleet}
			on:cargo-transfer-dialog={(e) => dispatch('cargo-transfer-dialog', e.detail)}
		/>
		<FleetWaypointsTile fleet={$commandedFleet} />
		<FleetWaypointTaskTile
			fleet={$commandedFleet}
			on:transport-tasks-dialog={(e) => dispatch('transport-tasks-dialog', e.detail)}
		/>
	</div>
{/if}
