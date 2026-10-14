<script lang="ts">
	import { ArrowsPointingOut } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { addError } from '#lib/services/Errors.js';
	import type {
		BattlePlanChangedEvent,
		CargoTransferDialogEvent,
		ChangeMassDriverSpeedEvent,
		ChangeWaypointEvent,
		ChangeWaypointTransportTasksEvent,
		ClearProductionQueueEvent,
		MergeFleetsDialogEvent,
		MergeFleetsEvent,
		RenameFleetEvent,
		SelectWaypointEvent,
		SplitFleetDialogEvent,
		SplitFleetEvent,
		TransferCargoEvent,
		TransportTasksDialogEvent
	} from '#lib/services/Events.js';
	import { getGameContext } from '#lib/services/GameContext.js';
	import { absoluteSize } from '#lib/types/CargoTransferRequest.js';
	import { MapObjectType, type WaypointDest } from '#lib/types/cs-proto.js';
	import {
		commandable,
		deepSpaceMapObject,
		equal as mapObjectEqual,
		ownedBy,
		type MapObjectLike
	} from '#lib/types/MapObject.js';
	import { MaxNameLength, None } from '#lib/types/Consts.js';
	import { emptyVector, equal } from '#lib/types/Vector.js';
	import type { ConnectError } from '@connectrpc/connect';
	import hotkeys from 'hotkeys-js';
	import { onMount } from 'svelte';
	import CargoTranfserDialog from '../dialogs/cargo/CargoTransferDialog.svelte';
	import MergeFleetsDialog from '../dialogs/merge/MergeFleetsDialog.svelte';
	import MailDialog from '../dialogs/mail/MailDialog.svelte';
	import NoteDialog from '../dialogs/note/NoteDialog.svelte';
	import ProductionQueueDialog from '../dialogs/production/ProductionQueueDialog.svelte';
	import SplitFleetDialog from '../dialogs/split/SplitFleetDialog.svelte';
	import TransportTasksDialog from '../dialogs/transport/TransportTasksDialog.svelte';
	import SearchDialog from '../search/SearchDialog.svelte';
	import CommandPane from './command/CommandPane.svelte';
	import CommandPaneDrawer from './command/CommandPaneDrawer.svelte';
	import MapObjectStatsBar from './MapObjectStatsBar.svelte';
	import MapObjectSummary from './MapObjectSummary.svelte';
	import MapObjectSummaryDialog from './MapObjectSummaryDialog.svelte';
	import Scanner from './scanner/Scanner.svelte';
	import MessagesPane from './MessagesPane.svelte';
	import ScannerToolbar from './scanner/ScannerToolbar.svelte';

	const {
		game,
		settings,
		universe,
		player,
		selectedMapObject,
		commandedMapObject,
		commandedPlanet,
		commandedFleet,
		commandMapObject,
		currentSelectedWaypointIndex,
		zoomToMapObject,
		nextMapObject,
		previousMapObject,
		selectWaypoint,
		selectMapObject,
		selectNextMapObject,
		addWaypoint,
		updateWaypoint,
		deleteWaypoint,
		renameFleet,
		updateFleetOrders,
		updatePlanetOrders,
		transferCargo,
		split,
		splitAll,
		merge
	} = getGameContext();

	// the fastest warp the speed hotkeys will set
	const maxWarpSpeed = 10;

	// like the original, the toolbar spans the whole window when the window is tall enough that the
	// command pane can spare the room. Otherwise it sits over the scanner.
	// The single column command pane is taller, so it needs a taller window
	const toolbarFullWidthMinHeight = 960;
	const toolbarFullWidthMinHeightSingleColumn = 1100;
	let innerWidth = $state(0);
	let innerHeight = $state(0);
	let toolbarHeight = $state(0);
	let toolbarFullWidth = $derived(
		innerWidth >= 768 &&
			innerHeight >=
				(innerWidth >= 1024 ? toolbarFullWidthMinHeight : toolbarFullWidthMinHeightSingleColumn)
	);

	let showProductionQueueDialog = $state(false);
	let showCargoTransferDialog = $state(false);
	let showMergeFleetsDialog = $state(false);
	let showSplitFleetDialog = $state(false);
	let showTransportTasksDialog = $state(false);
	let showSearchDialog = $state(false);
	let showSummaryDialog = $state(false);
	let cargoTransferDialogEvent: CargoTransferDialogEvent | undefined = $state(undefined);
	let mergeFleetsDialogEvent: MergeFleetsDialogEvent | undefined = $state(undefined);
	let splitFleetDialogEvent: SplitFleetDialogEvent | undefined = $state(undefined);
	let transportTasksDialogEvent: TransportTasksDialogEvent | undefined = $state(undefined);

	onMount(() => {
		hotkeys('n', 'root', () => {
			nextMapObject();
		});
		hotkeys('p', 'root', () => {
			previousMapObject();
		});
		hotkeys('q', 'root', () => {
			if ($commandedPlanet) {
				showProductionQueueDialog = true;
			}
		});
		hotkeys('⌘+k', 'root', (event) => {
			event.preventDefault();
			showSearchDialog = true;
		});
		hotkeys('Delete', 'root', () => {
			onDeleteWaypoint();
		});
		hotkeys('Backspace', 'root', () => {
			onDeleteWaypoint();
		});
		hotkeys('1, 2, 3, 4, 5, 6', 'root', (event) => {
			$settings.planetViewState = Number(event.code.slice(-1)) - 1;
		});
		hotkeys('7', 'root', () => {
			$settings.showScanners = !$settings.showScanners;
		});
		hotkeys('8', 'root', () => {
			settings.update((s) => {
				s.toggleMinefields();
				return s;
			});
		});
		hotkeys('0', 'root', () => {
			$settings.showNames = !$settings.showNames;
		});
		hotkeys('shift+0', 'root', () => {
			$settings.showFleetTokenCounts = !$settings.showFleetTokenCounts;
		});
		hotkeys(',', 'root', () => changeWarpSpeed(-1));
		hotkeys('.', 'root', () => changeWarpSpeed(1));
		hotkeys('[', 'root', () => selectAdjacentEnemyFleet(-1));
		hotkeys(']', 'root', () => selectAdjacentEnemyFleet(1));
		hotkeys('f', 'root', selectFirstMapObject);
		hotkeys('shift+n', 'root', () => {
			if (!selectAdjacentStarbase(1)) {
				nextMapObject();
			}
		});
		hotkeys('shift+p', 'root', () => {
			if (!selectAdjacentStarbase(-1)) {
				previousMapObject();
			}
		});
		hotkeys('r', 'root', (event) => {
			if ($commandedFleet) {
				// don't type the r into the prompt
				event.preventDefault();
				promptRenameFleet();
			}
		});

		// on the game view, prevent scrolling
		document.documentElement.classList.add('overscroll-none');

		return () => {
			hotkeys.unbind('n', 'root');
			hotkeys.unbind('p', 'root');
			hotkeys.unbind('q', 'root');
			hotkeys.unbind('⌘+k', 'root');
			hotkeys.unbind('Delete', 'root');
			hotkeys.unbind('Backspace', 'root');
			hotkeys.unbind('1, 2, 3, 4, 5, 6, 7, 8, 0, shift+0', 'root');
			hotkeys.unbind(',', 'root');
			hotkeys.unbind('., [, ], f, shift+n, shift+p, r', 'root');

			// go back to normal
			document.documentElement.classList.remove('overscroll-none');
		};
	});

	function commandAndZoom(mo: MapObjectLike) {
		commandMapObject(mo);
		selectMapObject(mo);
		zoomToMapObject(mo);
	}

	// speed up or slow down the selected waypoint. With no waypoint selected, this changes the next one
	function changeWarpSpeed(change: number) {
		const fleet = $commandedFleet;
		const waypointIndex = Math.max(1, $currentSelectedWaypointIndex);
		if (!fleet || waypointIndex >= fleet.fleetOrders.waypoints.length) {
			return;
		}
		const waypoint = fleet.fleetOrders.waypoints[waypointIndex];
		const warpSpeed = waypoint.warpSpeed + change;
		if (warpSpeed < 0 || warpSpeed > maxWarpSpeed) {
			return;
		}
		waypoint.warpSpeed = warpSpeed;
		onChangeWaypoint({ fleet, waypoint, waypointIndex });
	}

	// cycle through the fleets we know about that belong to other players
	function selectAdjacentEnemyFleet(offset: number) {
		const fleets = $universe.fleetIntels.filter((f) => !ownedBy(f, $player.num));
		if (fleets.length === 0) {
			return;
		}
		const index = fleets.findIndex((f) => mapObjectEqual($selectedMapObject, f));
		const next = index === -1 ? (offset > 0 ? 0 : -1) : index + offset;
		const fleet = fleets[(next + fleets.length) % fleets.length];
		selectMapObject(fleet);
		zoomToMapObject(fleet);
	}

	// go to our first planet, or our first fleet if we are already there
	function selectFirstMapObject() {
		const byNum = (a: MapObjectLike, b: MapObjectLike) =>
			(a.mapObject?.num ?? 0) - (b.mapObject?.num ?? 0);
		const first = [
			...$universe.getMyPlanets('num', false).sort(byNum).slice(0, 1),
			...$universe.getMyFleets('num', false).sort(byNum).slice(0, 1)
		].find((mo) => !mapObjectEqual($commandedMapObject, mo));
		if (first) {
			commandAndZoom(first);
		}
	}

	// from a planet, go to the next or previous of our planets with a starbase that can build ships
	function selectAdjacentStarbase(offset: number): boolean {
		const num = $commandedPlanet?.mapObject.num;
		if (num === undefined) {
			return false;
		}
		const planets = $universe
			.getMyPlanets('num', false)
			.filter(
				(p) =>
					p.spec?.planetStarbaseSpec?.hasStarbase && p.spec.planetStarbaseSpec.dockCapacity !== 0
			)
			.sort((a, b) => (a.mapObject?.num ?? 0) - (b.mapObject?.num ?? 0));
		const planet =
			offset > 0
				? (planets.find((p) => (p.mapObject?.num ?? 0) > num) ?? planets[0])
				: (planets.findLast((p) => (p.mapObject?.num ?? 0) < num) ?? planets[planets.length - 1]);
		if (planets.length > 0) {
			commandAndZoom(planet);
		}
		return true;
	}

	function promptRenameFleet() {
		const fleet = $commandedFleet;
		if (!fleet) {
			return;
		}
		const name =
			prompt('Enter fleet name', fleet.baseName) ||
			$universe.getMyDesign(fleet.tokens[0]?.designNum)?.name;
		if (name) {
			renameFleet(fleet, name.slice(0, MaxNameLength));
		}
	}

	function onNextMapObject() {
		nextMapObject();
	}

	function onPreviousMapObject() {
		previousMapObject();
	}

	async function onAddWaypoint(dest: WaypointDest, fastestWaypoint: boolean): Promise<boolean> {
		return addWaypoint(dest, fastestWaypoint);
	}

	async function onUpdateWaypointDest(dest: WaypointDest, fastestWaypoint: boolean, done: boolean) {
		updateWaypoint(dest, fastestWaypoint, done);
	}

	async function onDeleteWaypoint() {
		deleteWaypoint();
	}

	function onSelectWaypoint(e: SelectWaypointEvent) {
		const wp = e.waypoint;
		selectWaypoint(wp);

		if (wp.mapObjectTarget?.targetType && wp.mapObjectTarget.targetNum) {
			const mo = $universe.getMapObject(wp.mapObjectTarget);
			if (mo) {
				selectMapObject(mo);
			}
		} else {
			// a waypoint in deep space selects that point in space
			selectMapObject(deepSpaceMapObject(wp.position));
		}
	}

	async function onChangeWaypoint(e: ChangeWaypointEvent) {
		e.fleet.fleetOrders.waypoints[e.waypointIndex] = e.waypoint;
		updateFleetOrders(e.fleet);
	}

	async function onUpdateTransportTasks(e: ChangeWaypointTransportTasksEvent) {
		// update the transport tasks for this waypoint and update it
		e.waypoint.transportTasks = e.transportTasks;
		await onChangeWaypoint(e);

		// close the dialog
		showTransportTasksDialog = false;
	}

	async function onRenameFleet(e: RenameFleetEvent) {
		renameFleet(e.fleet, e.name);
	}

	async function onBattlePlanChanged(e: BattlePlanChangedEvent) {
		updateFleetOrders(e.fleet);
	}

	async function onChangeMassDriverSpeed(e: ChangeMassDriverSpeedEvent) {
		updatePlanetOrders(e.planet);
	}

	async function onClearProductionQueue(e: ClearProductionQueueEvent) {
		updatePlanetOrders(e.planet);
	}

	async function onSplitAll() {
		if (!$commandedFleet) {
			return;
		}
		splitAll($commandedFleet);
	}

	async function onTransferCargo(e: TransferCargoEvent) {
		// close the dialog
		showCargoTransferDialog = false;

		if (absoluteSize(e.transferAmount) > 0) {
			try {
				await transferCargo(e.src, e.dest, e.transferAmount);
			} catch (err) {
				addError(err as ConnectError);
			}
		}
	}

	async function onNextPlanet(updateOrders: boolean) {
		if (!$commandedPlanet) {
			return;
		}
		if (updateOrders) {
			await updatePlanetOrders($commandedPlanet);
		}

		nextMapObject();
	}

	async function onPrevPlanet(updateOrders: boolean) {
		if (!$commandedPlanet) {
			return;
		}
		if (updateOrders) {
			await updatePlanetOrders($commandedPlanet);
		}

		previousMapObject();
	}

	async function onMergeFleets(e: MergeFleetsEvent) {
		await merge(e.fleet, e.fleetNums);
		// close the dialog
		showMergeFleetsDialog = false;
	}

	async function onSplitFleet(e: SplitFleetEvent) {
		try {
			await split(e.src, e.dest, e.srcTokens, e.destTokens, e.transferAmount);
		} catch (err) {
			addError(err as ConnectError);
			throw err;
		}

		// close the dialog
		showSplitFleetDialog = false;
	}

	function selectSearchResult(mo: MapObjectLike | undefined) {
		if (mo) {
			if (ownedBy(mo, $player.num)) {
				commandMapObject(mo);
			}
			selectMapObject(mo);
			zoomToMapObject(mo);
		}
		showSearchDialog = false;
	}

	// onSelectMapObject cycles through commanding MapObjects at a location, and then
	// selecting non-commandable mapobjects
	function onSelectMapObject(mo: MapObjectLike) {
		const here = (other: MapObjectLike | undefined) =>
			!!other &&
			equal(other.mapObject?.position ?? emptyVector(), mo.mapObject?.position ?? emptyVector());
		// Commanding a fleet with a route moves the selection to the end of its route. We're still
		// cycling through what's here if that's the only reason the selection is somewhere else.
		const route = $commandedFleet?.fleetOrders.waypoints ?? [];
		const selectionOnRouteEnd =
			route.length > 1 &&
			!!$selectedMapObject &&
			equal(
				$selectedMapObject.mapObject?.position ?? emptyVector(),
				route[route.length - 1].position ?? emptyVector()
			);
		if (!here($selectedMapObject) && !(here($commandedMapObject) && selectionOnRouteEnd)) {
			// nothing selected, or nothing selected at this location yet, select this object
			selectMapObject(mo);
			return;
		}

		// get all the mapobjects here we want to cycle, starting with commandable map objects
		const commandables = $universe.getCommandableMapObjectsByPosition(
			mo.mapObject?.position ?? emptyVector()
		);
		const selectables = $universe
			.getMapObjectsByPosition(mo.mapObject?.position ?? emptyVector())
			.filter((mo) => !commandable($player.num, mo));
		let commandedIndex = commandables.findIndex((mo) => mapObjectEqual($commandedMapObject, mo));
		let selectedIndex = selectables.findIndex((mo) => mapObjectEqual($selectedMapObject, mo));

		if (commandables.length === 0) {
			// nothing of ours is here, i.e. an enemy planet, cycle through what is
			if (selectables.length > 0) {
				selectMapObject(selectables[(selectedIndex + 1) % selectables.length]);
			}
			return;
		}

		if (commandedIndex < commandables.length - 1) {
			// we either havne't commanded anything yet (commandedIndex=-1) or there is a commandable object to cycle to
			// if we are at the end of the commanded list, this will skip
			commandMapObject(commandables[commandedIndex + 1]);
		} else if (selectedIndex === -1) {
			if (selectables.length > 0) {
				// nothing selected, selecting first selectable
				selectMapObject(selectables[0]);
			} else {
				// nothing selected, nothing selectable, selecting first commandable
				selectMapObject(commandables[0]);
				commandMapObject(commandables[0]);
			}
		} else if (selectedIndex !== -1 && selectedIndex < selectables.length - 1) {
			// Cycle to the next selectable object
			selectMapObject(selectables[selectedIndex + 1]);
		} else if (selectables.length > 0) {
			// If at the end, wrap around to the first selectable object
			if (commandables.length > 0) {
				selectMapObject(commandables[0]);
				commandMapObject(commandables[0]);
			} else {
				selectMapObject(selectables[0]);
			}
		}
	}

	function onSetPacketDest(mo: MapObjectLike) {
		if (mo.mapObject?.type !== MapObjectType.PLANET) {
			return;
		} else {
			$settings.setPacketDest = false;
			// something went wrong, can't set dest on a planet without a massdriver
			if (!$commandedPlanet?.spec.planetStarbaseSpec?.hasMassDriver) {
				return;
			}

			if (mapObjectEqual(mo, $commandedPlanet)) {
				// clear dest
				$commandedPlanet.planetOrders.packetTargetNum = None;
			} else {
				$commandedPlanet.planetOrders.packetTargetNum = mo.mapObject.num;
			}

			updatePlanetOrders($commandedPlanet);
		}
	}

	function onSetRouteDest(mo: MapObjectLike) {
		if (!$commandedPlanet) {
			return;
		}
		if (mo.mapObject?.type !== MapObjectType.PLANET) {
			return;
		} else {
			$settings.setRouteDest = false;

			if (mapObjectEqual(mo, $commandedPlanet)) {
				// clear dest
				$commandedPlanet.planetOrders.routeTargetNum = None;
				$commandedPlanet.planetOrders.routeTargetPlayerNum = None;
				$commandedPlanet.planetOrders.routeTargetType = MapObjectType.UNSPECIFIED;
			} else {
				$commandedPlanet.planetOrders.routeTargetNum = mo.mapObject.num;
				$commandedPlanet.planetOrders.routeTargetPlayerNum = mo.mapObject.playerNum;
				$commandedPlanet.planetOrders.routeTargetType = mo.mapObject.type;
			}

			updatePlanetOrders($commandedPlanet);
		}
	}
</script>

<svelte:window bind:innerWidth bind:innerHeight />

{#snippet toolbar()}
	<ScannerToolbar
		onShowSearch={() => (showSearchDialog = true)}
		onCycleMapObject={() => selectNextMapObject()}
		{onNextMapObject}
		{onPreviousMapObject}
	/>
{/snippet}

<div class="flex flex-col h-full select-none" data-type="game-view" data-id={$game.id}>
	{#if toolbarFullWidth}
		<div class="flex-none mb-1" bind:clientHeight={toolbarHeight}>
			{@render toolbar()}
		</div>
	{/if}
	<!-- for small mobile displays we put the scanner on top and the command pane below it-->
	<div class="flex flex-col grow min-h-0 desk:flex-row phone-landscape:flex-row">
		<!-- for medium+ displays, and a phone on its side, command pane goes on the left -->
		<div
			class="hidden overflow-x-hidden desk:flex desk:flex-col desk:flex-none justify-between desk:w-62 desk:lg:w-120 overflow-y-auto desk:max-h-[calc(100dvh-4rem)] phone-landscape:flex phone-landscape:flex-col phone-landscape:flex-none phone-landscape:w-62 phone-landscape:max-h-[calc(100dvh-4rem)]"
			style:max-height={toolbarFullWidth
				? `calc(100dvh - 4.25rem - ${toolbarHeight}px)`
				: undefined}
		>
			<div class="flex flex-row flex-wrap gap-2 justify-center">
				<CommandPane
					{onNextMapObject}
					{onPreviousMapObject}
					{onRenameFleet}
					{onSelectWaypoint}
					{onChangeWaypoint}
					{onDeleteWaypoint}
					{onSplitAll}
					{onBattlePlanChanged}
					{onChangeMassDriverSpeed}
					{onClearProductionQueue}
					onShowProductionQueueDialog={() => (showProductionQueueDialog = true)}
					onShowCargoTransferDialog={(e) => {
						showCargoTransferDialog = true;
						cargoTransferDialogEvent = e;
					}}
					onShowMergeFleetDialog={(e) => {
						showMergeFleetsDialog = true;
						mergeFleetsDialogEvent = e;
					}}
					onShowSplitFleetDialog={(e) => {
						showSplitFleetDialog = true;
						splitFleetDialogEvent = e;
					}}
					onShowTransportTasksDialog={(e) => {
						showTransportTasksDialog = true;
						transportTasksDialogEvent = e;
					}}
				/>
			</div>
			<div class="hidden desk:lg:block lg:p-1 mx-2">
				<MapObjectSummary
					onShowCargoTransferDialog={(e) => {
						showCargoTransferDialog = true;
						cargoTransferDialogEvent = e;
					}}
				/>
			</div>
		</div>

		<div class="flex flex-col grow min-w-0">
			<div class="flex flex-col grow border-gray-700 border-2 shadow-xs">
				{#if !toolbarFullWidth}
					{@render toolbar()}
				{/if}
				<MessagesPane bind:showMessages={$settings.showMessagePane} messages={$player.messages} />
				<Scanner
					{onSelectWaypoint}
					{onAddWaypoint}
					{onUpdateWaypointDest}
					{onSelectMapObject}
					{onSetPacketDest}
					{onSetRouteDest}
				/>
			</div>
			<div class="hidden desk:block">
				<MapObjectStatsBar />
			</div>
			<!-- a phone on its side has no room for the summary, so it pops up from a button -->
			<button
				type="button"
				class="hidden phone-landscape:flex flex-row items-center gap-2 w-full text-left pl-2 pr-1 py-1 bg-base-200 border-t border-base-300"
				aria-label="show summary"
				data-type="show-map-object-summary-button"
				onclick={() => (showSummaryDialog = true)}
			>
				<div class="grow min-w-0 overflow-hidden"><MapObjectStatsBar large /></div>
				<span class="btn btn-xs btn-outline btn-secondary flex-none">
					<Icon src={ArrowsPointingOut} size="16" class="hover:stroke-accent" />
				</span>
			</button>
			<div class="hidden desk:block desk:w-full desk:lg:hidden mb-2">
				<MapObjectSummary
					onShowCargoTransferDialog={(e) => {
						showCargoTransferDialog = true;
						cargoTransferDialogEvent = e;
					}}
				/>
			</div>
		</div>

		<!-- for upright phone displays, use a drawer -->
		<div class="flex flex-col">
			<CommandPaneDrawer
				{onNextMapObject}
				{onPreviousMapObject}
				{onRenameFleet}
				{onSelectWaypoint}
				{onChangeWaypoint}
				{onDeleteWaypoint}
				{onSplitAll}
				{onBattlePlanChanged}
				{onChangeMassDriverSpeed}
				{onClearProductionQueue}
				onShowProductionQueueDialog={() => (showProductionQueueDialog = true)}
				onShowCargoTransferDialog={(e) => {
					showCargoTransferDialog = true;
					cargoTransferDialogEvent = e;
				}}
				onShowMergeFleetDialog={(e) => {
					showMergeFleetsDialog = true;
					mergeFleetsDialogEvent = e;
				}}
				onShowSplitFleetDialog={(e) => {
					showSplitFleetDialog = true;
					splitFleetDialogEvent = e;
				}}
				onShowTransportTasksDialog={(e) => {
					showTransportTasksDialog = true;
					transportTasksDialogEvent = e;
				}}
			/>
		</div>
	</div>
</div>

<div class="select-none">
	<!-- dialog modals -->
	<MapObjectSummaryDialog
		show={showSummaryDialog}
		onCancel={() => (showSummaryDialog = false)}
		onShowCargoTransferDialog={(e) => {
			showSummaryDialog = false;
			showCargoTransferDialog = true;
			cargoTransferDialogEvent = e;
		}}
	/>
	<ProductionQueueDialog
		show={showProductionQueueDialog}
		onNext={() => onNextPlanet(true)}
		onPrev={() => onPrevPlanet(true)}
		onOk={(planet) => {
			showProductionQueueDialog = false;
			updatePlanetOrders(planet);
		}}
		onCancel={() => (showProductionQueueDialog = false)}
	/>
	<CargoTranfserDialog
		show={showCargoTransferDialog}
		props={cargoTransferDialogEvent}
		onOk={onTransferCargo}
		onCancel={() => (showCargoTransferDialog = false)}
	/>
	<MergeFleetsDialog
		show={showMergeFleetsDialog}
		props={mergeFleetsDialogEvent}
		onOk={onMergeFleets}
		onCancel={() => (showMergeFleetsDialog = false)}
	/>
	<SplitFleetDialog
		show={showSplitFleetDialog}
		props={splitFleetDialogEvent}
		onOk={onSplitFleet}
		onCancel={() => (showSplitFleetDialog = false)}
	/>
	<TransportTasksDialog
		show={showTransportTasksDialog}
		props={transportTasksDialogEvent}
		onOk={onUpdateTransportTasks}
		onCancel={() => (showTransportTasksDialog = false)}
	/>
	<SearchDialog
		show={showSearchDialog}
		onOk={(e) => selectSearchResult(e)}
		onCancel={() => (showSearchDialog = false)}
	/>
	<NoteDialog />
	<MailDialog />
</div>
