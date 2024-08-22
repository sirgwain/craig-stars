<script lang="ts">
	import { carouselKey, createCarouselContext } from '$lib/services/CarouselContext';
	import { getGameContext } from '$lib/services/GameContext';
	import { equal, getMapObjectName } from '$lib/types/MapObject';
	import { ChevronDown, ChevronUp } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { afterUpdate, createEventDispatcher, onDestroy, setContext } from 'svelte';
	import type { MouseEventHandler, UIEventHandler } from 'svelte/elements';
	import type { CargoTransferDialogEvent } from '../../dialogs/cargo/CargoTranfserDialog.svelte';
	import type { MergeFleetsDialogEvent } from '../../dialogs/merge/MergeFleetsDialog.svelte';
	import type { ProductionQueueDialogEvent } from '../../dialogs/production/ProductionQueueDialog.svelte';
	import type { SplitFleetDialogEvent } from '../../dialogs/split/SplitFleetDialog.svelte';
	import type { TransportTasksDialogEvent } from '../../dialogs/transport/TransportTasksDialog.svelte';
	import FleetSummary from '../FleetSummary.svelte';
	import MapObjectSummary from '../MapObjectSummary.svelte';
	import FleetCompositionTile from './FleetCompositionTile.svelte';
	import FleetFuelAndCargoTile from './FleetFuelAndCargoTile.svelte';
	import FleetOrbitingTile from './FleetOrbitingTile.svelte';
	import FleetOtherFleetsHereTile from './FleetOtherFleetsHereTile.svelte';
	import FleetWaypointTaskTile from './FleetWaypointTaskTile.svelte';
	import FleetWaypointsTile, { type DeleteWaypointEvent } from './FleetWaypointsTile.svelte';
	import PlanetFleetsInOrbitTile from './PlanetFleetsInOrbitTile.svelte';
	import PlanetMineralsOnHandTile from './PlanetMineralsOnHandTile.svelte';
	import PlanetProductionTile from './PlanetProductionTile.svelte';
	import PlanetStarbaseTile from './PlanetStarbaseTile.svelte';
	import PlanetStatusTile from './PlanetStatusTile.svelte';

	const dispatch = createEventDispatcher<
		SplitFleetDialogEvent &
			MergeFleetsDialogEvent &
			CargoTransferDialogEvent &
			ProductionQueueDialogEvent &
			TransportTasksDialogEvent &
			DeleteWaypointEvent
	>();
	const {
		universe,
		commandedFleet,
		commandedPlanet,
		commandedMapObjectKey,
		selectedMapObject,
		selectedWaypoint,
		currentSelectedWaypointIndex,
		splitAll
	} = getGameContext();

	// setup the carouselContext for our child CommandTiles
	const carouselContext = createCarouselContext();
	setContext(carouselKey, carouselContext);

	const { open } = carouselContext;

	export let isOpen: boolean = $open;
	$: isOpen = $open;

	let carousel: HTMLDivElement | undefined;
	let activeNav = '#summary';
	let activeWaypointIndex: number | undefined;

	const onNavClicked: MouseEventHandler<HTMLAnchorElement> = (e) => {
		const a = e.currentTarget;
		const href = a.getAttribute('href') ?? '';
		if (href != '') {
			activeNav = href;
			scrollTo(href);
		}
	};

	// onScroll with a touch device, update the active nav
	const onScroll: UIEventHandler<HTMLDivElement> = (e) => {
		if (!e.currentTarget || !carousel) {
			return;
		}

		if ($commandedPlanet) {
			switch (e.currentTarget.scrollLeft) {
				case carousel.querySelector<HTMLDivElement>('#summary')!.offsetLeft:
					activeNav = '#summary';
					break;
				case carousel.querySelector<HTMLDivElement>('#planet-status-tile')!.offsetLeft:
					activeNav = '#planet-status-tile';
					break;
				case carousel.querySelector<HTMLDivElement>('#planet-minerals-on-hand-tile')!.offsetLeft:
					activeNav = '#planet-minerals-on-hand-tile';
					break;
				case carousel.querySelector<HTMLDivElement>('#planet-production-tile')!.offsetLeft:
					activeNav = '#planet-production-tile';
					break;
				case carousel.querySelector<HTMLDivElement>('#planet-starbase-tile')!.offsetLeft:
					activeNav = '#planet-starbase-tile';
					break;
				case carousel.querySelector<HTMLDivElement>('#planet-fleets-in-orbit-tile')!.offsetLeft:
					activeNav = '#planet-fleets-in-orbit-tile';
					break;
			}
		} else if ($commandedFleet) {
			switch (e.currentTarget.scrollLeft) {
				case carousel.querySelector<HTMLDivElement>('#summary')!.offsetLeft:
					activeNav = '#summary';
					break;
				case carousel.querySelector<HTMLDivElement>('#fleet-composition-tile')!.offsetLeft:
					activeNav = '#fleet-composition-tile';
					break;
				case carousel.querySelector<HTMLDivElement>('#fleet-orbiting-tile')!.offsetLeft:
					activeNav = '#fleet-orbiting-tile';
					break;
				case carousel.querySelector<HTMLDivElement>('#fleet-fuel-and-cargo-tile')!.offsetLeft:
					activeNav = '#fleet-fuel-and-cargo-tile';
					break;
				case carousel.querySelector<HTMLDivElement>('#fleet-waypoints-tile')!.offsetLeft:
					activeNav = '#fleet-waypoints-tile';
					break;
				case carousel.querySelector<HTMLDivElement>('#fleet-waypoint-task-tile')!.offsetLeft:
					activeNav = '#fleet-waypoint-task-tile';
					break;
				case carousel.querySelector<HTMLDivElement>('#fleet-other-fleets-here-tile')!.offsetLeft:
					activeNav = '#fleet-other-fleets-here-tile';
					break;
			}
		}
	};

	// scroll the view to an href like "#summary"
	function scrollTo(href: string) {
		if (!carousel) {
			return;
		}
		const target = carousel.querySelector<HTMLDivElement>(href)!;
		if (target && href != '') {
			carousel.scrollTo({ left: target.offsetLeft });
		}
	}

	// anytime the commandedMapObject is updated, open the command pane
	const unsuscribeCommandedMapObjectKey = commandedMapObjectKey.subscribe((mo) => {
		activeNav = '#summary';
	});

	// anytime the selectedMapObject is updated, show the summary
	const unsuscribeSelectedMapObject = selectedMapObject.subscribe((mo) => {
		if (!$selectedWaypoint || !equal(mo, $universe.getMapObject($selectedWaypoint))) {
			activeNav = '#summary';
		}
	});

	const unsubscribecurrentSelectedWaypointIndex = currentSelectedWaypointIndex.subscribe((i) => {
		// if there is a new currentlySelectedWaypointIndex and it's not the first one, show the waypoints tile
		if (i != -1 && i !== 0 && activeWaypointIndex !== i) {
			activeNav = '#fleet-waypoints-tile';
			activeWaypointIndex = i;
		}
	});

	// wait until the DOM is updated before calling scrollTo
	afterUpdate(() => {
		scrollTo(activeNav);
	});
	// unsubscribe on destroy to keep things tidy
	onDestroy(() => {
		unsuscribeCommandedMapObjectKey();
		unsuscribeSelectedMapObject();
		unsubscribecurrentSelectedWaypointIndex();
	});
</script>

<div
	class:hidden={!$open}
	class="carousel w-full md:hidden"
	bind:this={carousel}
	on:scrollend={onScroll}
>
	<div id="summary" class="carousel-item w-full">
		{#if $commandedFleet && equal($selectedMapObject, $universe.getPlanet($commandedFleet.orbitingPlanetNum))}
			<div class="w-full card bg-base-200 shadow rounded-sm border-2 border-base-300">
				<div class="card-body p-2 gap-0">
					<div class="flex flex-row items-center">
						<button class="w-full" on:click={carouselContext?.onDisclosureClicked}>
							<div class="flex-1 text-center text-lg font-semibold text-secondary">
								{$commandedFleet?.name ?? ''}
							</div>
						</button>
					</div>
					<FleetSummary
						fleet={$commandedFleet}
						on:cargo-transfer-dialog={(e) => dispatch('cargo-transfer-dialog', e?.detail)}
					/>
				</div>
			</div>
		{:else}
			<MapObjectSummary />
		{/if}
	</div>
	{#if $commandedPlanet}
		<div id="planet-status-tile" class="carousel-item w-full">
			<PlanetStatusTile planet={$commandedPlanet} />
		</div>
		<div id="planet-minerals-on-hand-tile" class="carousel-item w-full">
			<PlanetMineralsOnHandTile planet={$commandedPlanet} />
		</div>
		<div id="planet-production-tile" class="carousel-item w-full">
			<PlanetProductionTile
				planet={$commandedPlanet}
				on:change-production={(e) => dispatch('change-production', e?.detail)}
			/>
		</div>
		{#if $commandedPlanet.spec.hasStarbase}
			<div id="planet-starbase-tile" class="carousel-item w-full">
				<PlanetStarbaseTile
					planet={$commandedPlanet}
					starbase={$universe.getPlanetStarbase($commandedPlanet.num)}
				/>
			</div>
		{/if}
		<div id="planet-fleets-in-orbit-tile" class="carousel-item w-full">
			<PlanetFleetsInOrbitTile
				planet={$commandedPlanet}
				fleetsInOrbit={$universe.getMyFleetsByPosition($commandedPlanet)}
				on:cargo-transfer-dialog={(e) => dispatch('cargo-transfer-dialog', e?.detail)}
			/>
		</div>
	{:else if $commandedFleet}
		<div id="fleet-composition-tile" class="carousel-item w-full">
			<FleetCompositionTile
				fleet={$commandedFleet}
				on:split-all={() => $commandedFleet && splitAll($commandedFleet)}
				on:split-fleet-dialog={(e) => dispatch('split-fleet-dialog', e?.detail)}
				on:merge-fleets-dialog={(e) => dispatch('merge-fleets-dialog', e?.detail)}
			/>
		</div>
		<div id="fleet-orbiting-tile" class="carousel-item w-full">
			<FleetOrbitingTile
				fleet={$commandedFleet}
				on:cargo-transfer-dialog={(e) => dispatch('cargo-transfer-dialog', e?.detail)}
			/>
		</div>
		<div id="fleet-fuel-and-cargo-tile" class="carousel-item w-full">
			<FleetFuelAndCargoTile
				fleet={$commandedFleet}
				on:cargo-transfer-dialog={(e) => dispatch('cargo-transfer-dialog', e?.detail)}
			/>
		</div>
		<div id="fleet-waypoints-tile" class="carousel-item w-full">
			<FleetWaypointsTile
				fleet={$commandedFleet}
				on:delete-waypoint={(e) => dispatch('delete-waypoint')}
			/>
		</div>
		<div id="fleet-waypoint-task-tile" class="carousel-item w-full">
			<FleetWaypointTaskTile
				fleet={$commandedFleet}
				on:transport-tasks-dialog={(e) => dispatch('transport-tasks-dialog', e?.detail)}
			/>
		</div>
		<div id="fleet-other-fleets-here-tile" class="carousel-item w-full">
			<FleetOtherFleetsHereTile
				fleet={$commandedFleet}
				fleetsInOrbit={$universe
					.getMyFleetsByPosition($commandedFleet)
					.filter((f) => f.num !== $commandedFleet?.num)}
				on:cargo-transfer-dialog={(e) => dispatch('cargo-transfer-dialog', e?.detail)}
			/>
		</div>
	{/if}
</div>
<!-- Bottom menu -->
<div class="flex justify-center w-full md:hidden py-2 gap-2" class:hidden={!$open}>
	<a
		on:click|preventDefault={onNavClicked}
		href="#summary"
		class:btn-accent={activeNav == '#summary'}
		class="btn btn-xs">S</a
	>
	{#if $commandedPlanet}
		<a
			on:click|preventDefault={onNavClicked}
			href="#planet-status-tile"
			class:btn-accent={activeNav == '#planet-status-tile'}
			class="btn btn-xs">I</a
		>
		<a
			on:click|preventDefault={onNavClicked}
			href="#planet-minerals-on-hand-tile"
			class:btn-accent={activeNav == '#planet-minerals-on-hand-tile'}
			class="btn btn-xs">M</a
		>
		<a
			on:click|preventDefault={onNavClicked}
			href="#planet-production-tile"
			class:btn-accent={activeNav == '#planet-production-tile'}
			class="btn btn-xs">P</a
		>
		{#if $commandedPlanet.spec.hasStarbase}
			<a
				on:click|preventDefault={onNavClicked}
				href="#planet-starbase-tile"
				class:btn-accent={activeNav == '#planet-starbase-tile'}
				class="btn btn-xs">B</a
			>
		{/if}
		<a
			on:click|preventDefault={onNavClicked}
			href="#planet-fleets-in-orbit-tile"
			class:btn-accent={activeNav == '#planet-fleets-in-orbit-tile'}
			class="btn btn-xs">O</a
		>
	{:else if $commandedFleet}
		<a
			on:click|preventDefault={onNavClicked}
			href="#fleet-composition-tile"
			class:btn-accent={activeNav == '#fleet-composition-tile'}
			class="btn btn-xs">F</a
		>
		<a
			on:click|preventDefault={onNavClicked}
			href="#fleet-orbiting-tile"
			class:btn-accent={activeNav == '#fleet-orbiting-tile'}
			class="btn btn-xs">O</a
		>
		<a
			on:click|preventDefault={onNavClicked}
			href="#fleet-fuel-and-cargo-tile"
			class:btn-accent={activeNav == '#fleet-fuel-and-cargo-tile'}
			class="btn btn-xs">C</a
		>
		<a
			on:click|preventDefault={onNavClicked}
			href="#fleet-waypoints-tile"
			class:btn-accent={activeNav == '#fleet-waypoints-tile'}
			class="btn btn-xs">W</a
		>
		<a
			on:click|preventDefault={onNavClicked}
			href="#fleet-waypoint-task-tile"
			class:btn-accent={activeNav == '#fleet-waypoint-task-tile'}
			class="btn btn-xs">T</a
		>
		<a
			on:click|preventDefault={onNavClicked}
			href="#fleet-other-fleets-here-tile"
			class:btn-accent={activeNav == '#fleet-other-fleets-here-tile'}
			class="btn btn-xs">H</a
		>
	{/if}
</div>
<button class:hidden={$open} class="w-full p-2" on:click={carouselContext.onDisclosureClicked}>
	<div class="flex flex-row items-center">
		<div class="flex-1 text-center text-lg font-semibold text-secondary">
			{#if $commandedFleet && equal($selectedMapObject, $universe.getPlanet($commandedFleet.orbitingPlanetNum))}
				{getMapObjectName($commandedFleet)}
			{:else}
				{getMapObjectName($selectedMapObject) ?? 'Command Pane'}
			{/if}
		</div>
		{#if $open}
			<Icon src={ChevronUp} size="16" class="hover:stroke-accent" />
		{:else}
			<Icon src={ChevronDown} size="16" class="hover:stroke-accent" />
		{/if}
	</div>
</button>
