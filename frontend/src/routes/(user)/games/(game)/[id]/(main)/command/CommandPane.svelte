<script lang="ts">
	import type {
		BattlePlanChangedProps,
		ChangeMassDriverSpeedProps,
		ChangeWaypointProps,
		ClearProductionQueueProps,
		DeleteWaypointProps,
		NextPrevMapObjectProps,
		RenameFleetProps,
		SelectWaypointProps,
		ShowCargoTransferDialogProps,
		ShowMergeFleetsDialogProps,
		ShowProductionQueueDialogProps,
		ShowSplitFleetDialogProps,
		ShowTransportTasksDialogEventProps,
		SplitAllProps
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

	type Props = NextPrevMapObjectProps &
		RenameFleetProps &
		SplitAllProps &
		ShowCargoTransferDialogProps &
		ShowSplitFleetDialogProps &
		ShowMergeFleetsDialogProps &
		ShowProductionQueueDialogProps &
		ShowTransportTasksDialogEventProps &
		ClearProductionQueueProps &
		BattlePlanChangedProps &
		SelectWaypointProps &
		ChangeWaypointProps &
		DeleteWaypointProps &
		ChangeMassDriverSpeedProps;

	const {
		onNextMapObject,
		onPreviousMapObject,
		onRenameFleet,
		onSplitAll,
		onShowCargoTransferDialog,
		onShowSplitFleetDialog,
		onShowMergeFleetDialog,
		onShowProductionQueueDialog,
		onShowTransportTasksDialog,
		onBattlePlanChanged,
		onClearProductionQueue,
		onChangeMassDriverSpeed,
		onSelectWaypoint,
		onChangeWaypoint,
		onDeleteWaypoint
	}: Props = $props();

	const {
		universe,
		commandedPlanet,
		commandedFleet,
		selectedWaypoint,
		currentSelectedWaypointIndex
	} = getGameContext();
</script>

{#if $commandedPlanet}
	<div class="lg:flex lg:flex-col">
		<PlanetSummaryTile planet={$commandedPlanet} {onNextMapObject} {onPreviousMapObject} />
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
			starbase={$universe.getMyPlanetStarbase($commandedPlanet.num)}
			{onChangeMassDriverSpeed}
		/>
	</div>
{:else if $commandedFleet && $selectedWaypoint}
	<div class="lg:flex lg:flex-col">
		<FleetSummaryTile
			fleet={$commandedFleet}
			{onNextMapObject}
			{onPreviousMapObject}
			{onRenameFleet}
		/>
		<FleetOrbitingTile fleet={$commandedFleet} {onShowCargoTransferDialog} />
		<FleetOtherFleetsHereTile
			fleet={$commandedFleet}
			fleetsInOrbit={$universe.getFleetsByPosition($commandedFleet)}
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
			selectedWaypointIndex={$currentSelectedWaypointIndex}
			{onSelectWaypoint}
			{onChangeWaypoint}
			{onDeleteWaypoint}
		/>
		<FleetWaypointTaskTile
			fleet={$commandedFleet}
			selectedWaypointIndex={$currentSelectedWaypointIndex}
			{onShowTransportTasksDialog}
			{onChangeWaypoint}
		/>
	</div>
{/if}
