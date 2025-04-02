<script lang="ts">
	import { clickOutside } from '$lib/clickOutside';
	import DisclosureHeader from '$lib/components/DisclosureHeader.svelte';
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
	import { MapObjectTypeFleet, MapObjectTypePlanet } from '$lib/types/cs';
	import { commandable, equalsTarget, getMapObjectName } from '$lib/types/MapObject';
	import { distance, equal as equalPosition } from '$lib/types/Vector';
	import { slide } from 'svelte/transition';
	import MapObjectSummary from '../MapObjectSummary.svelte';
	import MapObjectSummaryCollapsed from '../MapObjectSummaryMini.svelte';
	import FleetCompositionTile from './FleetCompositionTile.svelte';
	import FleetFuelAndCargoTile from './FleetFuelAndCargoTile.svelte';
	import FleetOrbitingTile from './FleetOrbitingTile.svelte';
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

	const {
		player,
		settings,
		universe,
		commandedFleet,
		commandedMapObject,
		commandedPlanet,
		selectedMapObject,
		selectedWaypoint,
		currentSelectedWaypointIndex,
		commandMapObject
	} = getGameContext();

	type Props = {} & NextPrevMapObjectProps &
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

	let {
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

	let open = $state(false);

	let summaryMapObject = $derived.by(() => {
		if (!$selectedMapObject) {
			return $commandedMapObject;
		}

		// don't update the summaryMapObject  when adding waypoints
		if ($settings.addWaypoint) {
			return $commandedMapObject;
		}

		// if we select a fleet, put it in the summary
		if (
			$commandedMapObject &&
			$selectedMapObject &&
			$selectedMapObject.type === MapObjectTypeFleet
		) {
			return $selectedMapObject;
		}

		// if we are cycilng through our commandable fleets on a planet we own
		// make sure we select the commandable fleet (because the planet will be selected, due to the way the desktop ui works)
		if (
			$commandedMapObject &&
			equalPosition($selectedMapObject.position, $commandedMapObject.position) &&
			$commandedMapObject.type === MapObjectTypeFleet &&
			$selectedMapObject.type === MapObjectTypePlanet
		) {
			return $commandedMapObject;
		}

		// the selectedMapObject is different than the commanded map object, show it in the summary
		return $selectedMapObject;
	});

	let dist = $derived(
		$commandedMapObject && $selectedMapObject
			? distance($commandedMapObject.position, $selectedMapObject.position)
			: 0
	);

	function toggleDrawer() {
		open = !open;
		if (summaryMapObject && commandable($player.num, summaryMapObject)) {
			commandMapObject(summaryMapObject);
		}
	}
</script>

<div class="w-full md:hidden select-none z-10">
	<div class={open ? 'fixed inset-0 flex flex-col justify-end' : ''}>
		{#if open}
			<!-- Backdrop -->
			<div class="absolute inset-0 bg-black/40"></div>

			<div
				transition:slide={{ duration: 250 }}
				use:clickOutside={() => (open = false)}
				class="bg-base-200 border-t border-base-300 rounded-t-2xl h-[80vh] shadow-lg relative z-5 p-4"
			>
				<DisclosureHeader {open} onToggle={toggleDrawer}>
					<div class="flex flex-row justify-center w-full pb-1">
						{$commandedMapObject?.name}
					</div>
				</DisclosureHeader>
				<div class="overflow-y-auto max-h-[calc(80vh-3.5rem)] overflow-x-hidden">
					<div class="flex flex-col gap-1 justify-stretch">
						{#if $commandedPlanet}
							<div id="planet-summary-tile">
								<PlanetSummaryTile
									planet={$commandedPlanet}
									hideTitle={true}
									{onNextMapObject}
									{onPreviousMapObject}
								/>
							</div>
							<div id="summary">
								<MapObjectSummary
									{onShowCargoTransferDialog}
									hideTitle={true}
									hideCycleButton={true}
								/>
							</div>
							<div id="planet-production-tile">
								<PlanetProductionTile
									planet={$commandedPlanet}
									{onShowProductionQueueDialog}
									{onClearProductionQueue}
								/>
							</div>
							<div id="planet-status-tile">
								<PlanetStatusTile planet={$commandedPlanet} />
							</div>
							<div id="planet-minerals-on-hand-tile">
								<PlanetMineralsOnHandTile planet={$commandedPlanet} />
							</div>
							{#if $commandedPlanet.spec.hasStarbase}
								<div id="planet-starbase-tile">
									<PlanetStarbaseTile
										planet={$commandedPlanet}
										starbase={$universe.getMyPlanetStarbase($commandedPlanet.num)}
										{onChangeMassDriverSpeed}
									/>
								</div>
							{/if}
							<div id="planet-fleets-in-orbit-tile">
								<PlanetFleetsInOrbitTile
									planet={$commandedPlanet}
									fleetsInOrbit={$universe.getMyFleetsByPosition($commandedPlanet)}
									{onShowCargoTransferDialog}
								/>
							</div>
						{:else if $commandedFleet}
							<div id="planet-summary-tile">
								<FleetSummaryTile
									fleet={$commandedFleet}
									hideTitle={true}
									{onNextMapObject}
									{onPreviousMapObject}
									{onRenameFleet}
								/>
							</div>
							<div id="fleet-composition-tile">
								<FleetCompositionTile
									fleet={$commandedFleet}
									selectedWaypoint={$selectedWaypoint}
									{onShowSplitFleetDialog}
									{onShowMergeFleetDialog}
									{onSplitAll}
									{onBattlePlanChanged}
								/>
							</div>
							<div id="fleet-orbiting-tile">
								<FleetOrbitingTile fleet={$commandedFleet} {onShowCargoTransferDialog} />
							</div>
							<div id="fleet-fuel-and-cargo-tile">
								<FleetFuelAndCargoTile fleet={$commandedFleet} {onShowCargoTransferDialog} />
							</div>
							<div id="fleet-waypoints-tile">
								<FleetWaypointsTile
									fleet={$commandedFleet}
									selectedWaypointIndex={$currentSelectedWaypointIndex}
									{onSelectWaypoint}
									{onChangeWaypoint}
									{onDeleteWaypoint}
								/>
							</div>
							<div id="fleet-waypoint-task-tile">
								<FleetWaypointTaskTile
									fleet={$commandedFleet}
									selectedWaypointIndex={$currentSelectedWaypointIndex}
									{onShowTransportTasksDialog}
									{onChangeWaypoint}
								/>
							</div>
							<div id="fleet-other-fleets-here-tile">
								<FleetOtherFleetsHereTile
									fleet={$commandedFleet}
									cargoDestsInOrbit={$universe.getCargoDestsByPosition($commandedFleet)}
									{onShowCargoTransferDialog}
									{onShowSplitFleetDialog}
								/>
							</div>
						{/if}
					</div>
				</div>
			</div>
		{:else}
			<div
				transition:slide={{ duration: 250 }}
				class="bg-base-200 border-t border-base-300 rounded-t2xl relative p-4"
			>
				<DisclosureHeader
					{open}
					openable={commandable($player.num, summaryMapObject) ||
						equalsTarget(summaryMapObject, $selectedWaypoint)}
					onToggle={toggleDrawer}
				>
					<div class="flex flex-row w-full">
						<div class="text-sm text-left flex flex-col w-20">
							<div>
								{#if summaryMapObject?.num}
									ID: {summaryMapObject?.num}
								{/if}
							</div>
							<div>
								X: {summaryMapObject?.position.x}, Y: {summaryMapObject?.position.y}
							</div>
						</div>
						<div class="grow text-center">
							{summaryMapObject && summaryMapObject.name !== ''
								? summaryMapObject.name
								: 'Deep Space'}
						</div>
						<div class="text-sm my-auto w-20 text-right">
							{#if $commandedMapObject && dist}
								{dist.toFixed(1)} ly from {getMapObjectName($commandedMapObject)}
							{/if}
						</div>
					</div>
				</DisclosureHeader>

				<MapObjectSummaryCollapsed mapObject={summaryMapObject} />
			</div>
		{/if}
	</div>
</div>
