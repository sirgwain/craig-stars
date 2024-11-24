<script lang="ts">
	import type {
		BattlePlanChangedProps,
		ChangeMassDriverSpeedProps,
		ClearProductionQueueProps,
		ShowCargoTransferDialogProps,
		ShowMergeFleetsDialogProps,
		ShowProductionQueueDialogProps,
		ShowSplitFleetDialogProps,
		ShowTransportTasksDialogEventProps
	} from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
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

	type Props = {
		onDeleteWaypoint: () => Promise<void>;
		onSplitAll: () => Promise<void>;
	} & ShowCargoTransferDialogProps &
		ShowSplitFleetDialogProps &
		ShowMergeFleetsDialogProps &
		ShowProductionQueueDialogProps &
		ShowTransportTasksDialogEventProps &
		ClearProductionQueueProps &
		BattlePlanChangedProps &
		ChangeMassDriverSpeedProps;

	const {
		onDeleteWaypoint,
		onSplitAll,
		onShowCargoTransferDialog,
		onShowSplitFleetDialog,
		onShowMergeFleetDialog,
		onShowProductionQueueDialog,
		onShowTransportTasksDialog,
		onBattlePlanChanged,
		onClearProductionQueue,
		onChangeMassDriverSpeed
	}: Props = $props();

	const { universe, commandedPlanet, commandedFleet, selectedWaypoint } = getGameContext();
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
			{onShowCargoTransferDialog}
		/>
		<PlanetProductionTile
			planet={$commandedPlanet}
			{onShowProductionQueueDialog}
			{onClearProductionQueue}
		/>
		<PlanetStarbaseTile
			planet={$commandedPlanet}
			starbase={$universe.getPlanetStarbase($commandedPlanet.num)}
			{onChangeMassDriverSpeed}
		/>
	</div>
{:else if $commandedFleet}
	<div class="lg:flex lg:flex-col">
		<FleetSummaryTile fleet={$commandedFleet} />
		<FleetOrbitingTile fleet={$commandedFleet} {onShowCargoTransferDialog} />
		<FleetOtherFleetsHereTile
			fleet={$commandedFleet}
			fleetsInOrbit={$universe
				.getMyFleetsByPosition($commandedFleet)
				.filter((f) => f.num !== $commandedFleet?.num)}
			{onShowSplitFleetDialog}
			{onShowCargoTransferDialog}
		/>
		<FleetCompositionTile
			fleet={$commandedFleet}
			selectedWaypoint={$selectedWaypoint}
			{onSplitAll}
			{onShowSplitFleetDialog}
			{onShowMergeFleetDialog}
			{onBattlePlanChanged}
		/>
	</div>
	<div class="lg:flex lg:flex-col">
		<FleetFuelAndCargoTile fleet={$commandedFleet} {onShowCargoTransferDialog} />
		<FleetWaypointsTile
			fleet={$commandedFleet}
			selectedWaypoint={$selectedWaypoint}
			{onDeleteWaypoint}
		/>
		<FleetWaypointTaskTile
			fleet={$commandedFleet}
			selectedWaypoint={$selectedWaypoint}
			{onShowTransportTasksDialog}
		/>
	</div>
{/if}
