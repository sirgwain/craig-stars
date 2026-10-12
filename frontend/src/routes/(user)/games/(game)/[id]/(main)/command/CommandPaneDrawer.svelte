<script lang="ts">
	import { clickOutside } from '#lib/clickOutside.js';
	import DisclosureHeader from '#lib/components/DisclosureHeader.svelte';
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
	} from '#lib/services/Events.js';
	import { getGameContext } from '#lib/services/GameContext.js';
	import { MapObjectType } from '#lib/types/cs-proto.js';
	import { commandable, getMapObjectName } from '#lib/types/MapObject.js';
	import { distance, emptyVector } from '#lib/types/Vector.js';
	import { PencilSquare } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
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
	import TagBadge from '#lib/components/game/TagBadge.svelte';
	import PlanetStatusTile from './PlanetStatusTile.svelte';
	import PlanetSummaryTile from './PlanetSummaryTile.svelte';
	import { pickSummaryMapObject, summaryKey, type LastPicked } from './CommandPaneSummary';

	const {
		player,
		settings,
		universe,
		readOnly,
		editNote,
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

	// track whether the player last picked something to command or something to select. Refreshes
	// of the same mapobject don't count
	let lastPicked: LastPicked = $state('selected');
	let commandedKey = '';
	let selectedKey = '';
	$effect.pre(() => {
		const key = summaryKey($selectedMapObject);
		if (key !== selectedKey) {
			selectedKey = key;
			lastPicked = 'selected';
		}
	});
	$effect.pre(() => {
		const key = summaryKey($commandedMapObject);
		if (key !== commandedKey) {
			commandedKey = key;
			lastPicked = 'commanded';
		}
	});

	let summaryMapObject = $derived(
		pickSummaryMapObject($selectedMapObject, $commandedMapObject, $settings.addWaypoint, lastPicked)
	);

	// notes are written about planets and the player's own fleets
	let noteTarget = $derived(
		summaryMapObject && $universe.canHaveNote(summaryMapObject) ? summaryMapObject : undefined
	);
	let note = $derived($universe.getNote(noteTarget));

	let dist = $derived(
		$commandedMapObject && $selectedMapObject
			? distance(
					$commandedMapObject.mapObject?.position ?? emptyVector(),
					$selectedMapObject.mapObject?.position ?? emptyVector()
				)
			: 0
	);

	// true if the summary is showing the selected waypoint of the commanded fleet, i.e. the planet a
	// fleet was just sent back to. Planets are only keyed by num, so don't compare the owner.
	let summaryIsSelectedWaypoint = $derived.by(() => {
		const target = $commandedFleet ? $selectedWaypoint?.mapObjectTarget : undefined;
		const mo = summaryMapObject?.mapObject;
		return (
			!!target &&
			!!mo &&
			mo.type === target.targetType &&
			mo.num === target.targetNum &&
			(mo.type === MapObjectType.PLANET || mo.playerNum === target.targetPlayerNum)
		);
	});

	// The drawer opens on whatever is commanded, even when the summary is showing something else the
	// player selected. A summary of something commandable gets commanded first, unless it's where the
	// commanded fleet's selected waypoint is. Picking a waypoint doesn't mean leaving the fleet.
	function toggleDrawer() {
		if (
			!open &&
			summaryMapObject &&
			commandable($player.num, summaryMapObject) &&
			!summaryIsSelectedWaypoint
		) {
			commandMapObject(summaryMapObject);
		}
		open = !open;
	}
</script>

<div
	class="w-full desk:hidden phone-landscape:hidden select-none z-10"
	data-type="command-drawer"
	data-open={open}
>
	<div class={open ? 'fixed inset-0 flex flex-col justify-end' : ''}>
		{#if open}
			<!-- Backdrop -->
			<button
				type="button"
				aria-label="Close command pane backdrop"
				class="absolute inset-0 bg-black/40"
				onclick={() => (open = false)}
			></button>

			<div
				transition:slide={{ duration: 250 }}
				use:clickOutside={(event) => {
					// Dialogs open over the drawer; their controls and backdrops aren't outside taps.
					if (
						event.target instanceof Element &&
						event.target.closest('.modal-open, dialog[open]')
					) {
						return;
					}
					open = false;
				}}
				class="bg-base-200 border-t border-base-300 rounded-t-2xl h-[80dvh] shadow-lg relative z-5 p-4"
			>
				<DisclosureHeader {open} onToggle={toggleDrawer}>
					<div class="w-full min-w-0 truncate text-center pb-1">
						{$commandedMapObject?.mapObject?.name}
					</div>
				</DisclosureHeader>
				<div class="overflow-y-auto max-h-[calc(80dvh-3.5rem)] overflow-x-hidden">
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
							{#if $commandedPlanet.spec.planetStarbaseSpec?.hasStarbase}
								<div id="planet-starbase-tile">
									<PlanetStarbaseTile
										planet={$commandedPlanet}
										starbase={$universe.getMyPlanetStarbase($commandedPlanet.mapObject.num)}
										{onChangeMassDriverSpeed}
									/>
								</div>
							{/if}
							<div id="planet-fleets-in-orbit-tile">
								<PlanetFleetsInOrbitTile
									planet={$commandedPlanet}
									fleetsInOrbit={$universe.getMyFleetsByPosition(
										$commandedPlanet.mapObject.position ?? emptyVector()
									)}
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
							<div id="fleet-fuel-and-cargo-tile">
								<FleetFuelAndCargoTile fleet={$commandedFleet} {onShowCargoTransferDialog} />
							</div>
							<div id="fleet-orbiting-tile">
								<FleetOrbitingTile fleet={$commandedFleet} {onShowCargoTransferDialog} />
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
								<FleetWaypointTaskTile {onShowTransportTasksDialog} {onChangeWaypoint} />
							</div>
							<div id="fleet-other-fleets-here-tile">
								<FleetOtherFleetsHereTile
									fleet={$commandedFleet}
									cargoDestsInOrbit={$universe.getCargoDestsByPosition(
										$commandedFleet.mapObject.position ?? emptyVector()
									)}
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
					openable={!!$commandedMapObject || commandable($player.num, summaryMapObject)}
					onToggle={toggleDrawer}
				>
					<div class="flex flex-row w-full">
						<div class="text-sm text-left flex flex-col w-20">
							<div>
								{#if summaryMapObject?.mapObject?.num}
									ID: {summaryMapObject.mapObject.num}
								{/if}
							</div>
							<div>
								X: {summaryMapObject?.mapObject?.position?.x ?? 0}, Y: {summaryMapObject?.mapObject
									?.position?.y ?? 0}
							</div>
						</div>
						<div class="grow text-center">
							{summaryMapObject && summaryMapObject.mapObject?.name !== ''
								? summaryMapObject.mapObject?.name
								: 'Deep Space'}
						</div>
						<div class="text-sm my-auto w-20 text-right">
							{#if $commandedMapObject && dist}
								{dist.toFixed(1)} ly from {getMapObjectName($commandedMapObject)}
							{/if}
						</div>
					</div>
				</DisclosureHeader>

				<MapObjectSummaryCollapsed mapObject={summaryMapObject} {onShowCargoTransferDialog} />
				{#if noteTarget && (note || !$readOnly)}
					<button
						type="button"
						class="flex flex-row items-center gap-1 mt-2 w-full overflow-hidden text-sm text-left"
						aria-label={note ? 'Edit notes and tags' : 'Add notes and tags'}
						data-type="map-object-note-button"
						onclick={() => editNote(noteTarget)}
					>
						<Icon src={PencilSquare} class={`size-4 flex-none ${note ? 'stroke-accent' : ''}`} />
						{#each note?.tags ?? [] as tag (tag)}
							<TagBadge {tag} class="flex-none" />
						{/each}
						<span class="truncate italic" data-type="map-object-note"
							>{note ? note.text : 'Add notes or tags'}</span
						>
					</button>
				{/if}
			</div>
		{/if}
	</div>
</div>
