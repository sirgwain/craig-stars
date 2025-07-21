<script lang="ts">
	import MineralMini from '$lib/components/game/MineralMini.svelte';
	import OtherMapObjectsHere from '$lib/components/game/OtherMapObjectsHere.svelte';
	import WarpSpeedGauge from '$lib/components/game/WarpSpeedGauge.svelte';
	import type {
		ChangeWaypointProps,
		ShowTransportTasksDialogEventProps
	} from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import type { MapObject, TransportPlan, WaypointTask } from '$lib/types/cs';
	import {
		MapObjectTypePlanet,
		ReportAgeUnexplored,
		WaypointTaskLayMinefield,
		WaypointTaskNone,
		WaypointTaskPatrol,
		WaypointTaskRemoteMining,
		WaypointTaskTransferFleet,
		WaypointTaskTransport
	} from '$lib/types/cs';
	import { CommandedFleet, emptyTransportTasks, WaypointTasks } from '$lib/types/Fleet';
	import { owned, ownedBy } from '$lib/types/MapObject';
	import { getMineralOutput } from '$lib/types/Planet';
	import { PencilSquare } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { startCase } from 'lodash-es';
	import TransportTasksMini from '../../(plans)/transport-plans/TransportTasksMini.svelte';
	import CommandTile from './CommandTile.svelte';

	const { game, player, universe } = getGameContext();

	type Props = {
		fleet: CommandedFleet;
		selectedWaypointIndex: number;
	} & ChangeWaypointProps &
		ShowTransportTasksDialogEventProps;

	let {
		fleet: propFleet,
		selectedWaypointIndex,
		onShowTransportTasksDialog,
		onChangeWaypoint
	}: Props = $props();

	// local state for the ui components
	let fleet = $state(propFleet);
	let waypoint = $state(propFleet.waypoints[selectedWaypointIndex]);

	$effect(() => {
		// update state when the props change
		fleet = propFleet;
		waypoint = propFleet.waypoints[selectedWaypointIndex];
	});

	let selectedWaypointTask = $derived(waypoint.task ?? WaypointTaskNone);
	let selectedWaypointPlanet = $derived(
		waypoint.targetType == MapObjectTypePlanet && waypoint.targetNum
			? $universe.getPlanet(waypoint.targetNum)
			: undefined
	);

	const onSelectedWaypointTaskChange = (task: WaypointTask) => {
		waypoint.task = task;

		if (task != WaypointTaskTransport) {
			// if we aren't doing a transport, reset the transport tasks to blank.
			// If we don't do this, the user could pick a transport task in the future and assume it defaults to empty
			// but it will have whatever it last had
			waypoint.transportTasks = emptyTransportTasks();
		}

		onChangeWaypoint?.({ fleet, waypoint, waypointIndex: selectedWaypointIndex });
	};

	function onPatrolRangeChanged(value: number) {
		waypoint.patrolRange = value;
		onChangeWaypoint?.({ fleet, waypoint, waypointIndex: selectedWaypointIndex });
	}

	function onPatrolWarpSpeedChanged(warpSpeed: number) {
		waypoint.patrolWarpSpeed = warpSpeed;
		onChangeWaypoint?.({ fleet, waypoint, waypointIndex: selectedWaypointIndex });
	}

	function onPatrolWarpSpeedDragged(warpSpeed: number) {
		waypoint.patrolWarpSpeed = warpSpeed;
	}

	function onLayMinefieldDurationChanged(value: number | undefined) {
		waypoint.layMinefieldDuration = value;
		onChangeWaypoint?.({ fleet, waypoint, waypointIndex: selectedWaypointIndex });
	}

	function onTransferToPlayerChanged(playerNum: number) {
		waypoint.transferToPlayer = playerNum;
		onChangeWaypoint?.({ fleet, waypoint, waypointIndex: selectedWaypointIndex });
	}

	function applyTransportPlan(plan: TransportPlan | undefined) {
		if (!plan) return;

		waypoint.transportTasks = plan.tasks;
		onChangeWaypoint?.({ fleet, waypoint, waypointIndex: selectedWaypointIndex });
	}

	function onTargetChanged(target: Partial<MapObject>) {
		waypoint.targetName = target.name;
		waypoint.targetType = target.type;
		waypoint.targetNum = target.num;
		waypoint.targetPlayerNum = target.playerNum;
		onChangeWaypoint?.({ fleet, waypoint, waypointIndex: selectedWaypointIndex });
	}
</script>

<CommandTile title="Waypoint Task">
	<div class="flex justify-between">
		<div class="my-auto text-tile-item-title">Target</div>
		<div>
			<OtherMapObjectsHere
				{fleet}
				otherMapObjectsHere={$universe.getOtherMapObjectsHereByType(waypoint.position)}
				target={waypoint}
				position={waypoint.position}
				class="w-36"
				onSelected={onTargetChanged}
			/>
		</div>
	</div>
	<div class="flex justify-between my-1">
		<div class="my-auto text-tile-item-title">Task</div>
		<div>
			<select
				class="select select-outline select-secondary select-sm text-sm w-36"
				value={selectedWaypointTask}
				onchange={(e) => {
					onSelectedWaypointTaskChange(e.currentTarget.value ?? WaypointTaskNone);
				}}
			>
				{#each WaypointTasks as task (task)}
					{#if task === WaypointTaskNone}
						<option value={task}>None</option>
					{:else}
						<option value={task}>{startCase(task ?? 'None')}</option>
					{/if}
				{/each}
			</select>
		</div>
	</div>

	{#if waypoint.task === WaypointTaskTransport}
		<div class="flex flex-col">
			<div>
				<TransportTasksMini transportTasks={waypoint.transportTasks} />
			</div>
			<div class="ml-auto mt-1 flex flex-row gap-1">
				<div>
					<button
						type="button"
						onclick={() =>
							onShowTransportTasksDialog?.({
								fleet,
								waypoint,
								waypointIndex: selectedWaypointIndex
							})}
						class="btn btn-outline btn-sm normal-case btn-secondary inline-block p-1"
						><Icon src={PencilSquare} size="16" class="hover:stroke-accent inline" /></button
					>
				</div>
				<select
					class="select select-outline select-sm select-secondary w-12 sm:w-full text-secondary"
					onchange={(e) => {
						applyTransportPlan(
							$player.transportPlans.find((p) => p.num == parseInt(e.currentTarget.value))
						);
						e.currentTarget.value = '0';
					}}
				>
					<option value={0}>Apply Plan</option>
					{#each $player.transportPlans as plan (plan.num)}
						<option value={plan.num}>{plan.name}</option>
					{/each}
				</select>
			</div>
		</div>
	{:else if waypoint.task === WaypointTaskRemoteMining}
		{#if selectedWaypointPlanet}
			<!-- if this waypoint is owned -->
			{#if 'reportAge' in selectedWaypointPlanet && selectedWaypointPlanet.reportAge === ReportAgeUnexplored}
				<span class="text-warning"
					>Warning: This planet is unexplored. We have no way of knowing if we can mine it.</span
				>
			{:else if owned(selectedWaypointPlanet) && !($player.race.spec?.canRemoteMineOwnPlanets && ownedBy(selectedWaypointPlanet, $player.num))}
				<span class="text-error">Note: You can only remote mine unoccupied planets.</span>
			{:else if !fleet.spec.miningRate}
				<span class="text-error"
					>Warning: This fleet contains no ships with remote mining modules.</span
				>
			{:else}
				Mining Rate per Year:
				<MineralMini
					mineral={getMineralOutput(
						selectedWaypointPlanet,
						fleet.spec.miningRate ?? 0,
						$game.rules.remoteMiningMineOutput
					)}
					showUnits={true}
				/>
			{/if}
		{:else}
			<span class="text-error">Warning: Can only remote mine planets.</span>
		{/if}
	{:else if waypoint.task === WaypointTaskLayMinefield}
		<select
			class="select select-outline select-secondary select-sm py-0 text-sm mt-1"
			value={waypoint.layMinefieldDuration}
			onchange={(e) => onLayMinefieldDurationChanged(parseInt(e.currentTarget.value))}
		>
			<option value={undefined}>Indefinitely</option>
			<option value={1}>for 1 year</option>
			<option value={2}>for 2 years</option>
			<option value={3}>for 3 years</option>
			<option value={4}>for 4 years</option>
			<option value={5}>for 5 years</option>
		</select>
		<p class="text-warning">
			This fleet can lay {fleet.getTotalMinesLaidPerYear()} mines per year.
		</p>
	{:else if waypoint.task === WaypointTaskPatrol}
		<div class="flex justify-between my-1">
			<div class="my-auto text-tile-item-title">Intercept</div>
			<div>
				<select
					class="select select-outline select-secondary select-sm py-0 text-sm mt-1"
					value={waypoint.patrolRange}
					onchange={(e) => onPatrolRangeChanged(parseInt(e.currentTarget.value))}
				>
					<option value={50}>within 50 l.y.</option>
					<option value={100}>within 100 l.y.</option>
					<option value={150}>within 150 l.y.</option>
					<option value={200}>within 200 l.y.</option>
					<option value={250}>within 250 l.y.</option>
					<option value={300}>within 300 l.y.</option>
					<option value={350}>within 350 l.y.</option>
					<option value={450}>within 450 l.y.</option>
					<option value={550}>within 550 l.y.</option>
					<option value={undefined}>any enemy</option>
				</select>
			</div>
		</div>
		<div class="flex mt-1">
			<span class="text-tile-item-title">Warp Factor</span>
			<span class="flex-1 ml-1">
				<WarpSpeedGauge
					onValueChanged={(value) => onPatrolWarpSpeedChanged(value)}
					onValueDragged={(value) => onPatrolWarpSpeedDragged(value)}
					value={waypoint.patrolWarpSpeed}
					warnSpeed={fleet.spec.engine?.maxSafeSpeed
						? fleet.spec.engine.maxSafeSpeed + 1
						: undefined}
					warp0Text="Automatic"
				/>
			</span>
		</div>
	{:else if waypoint.task === WaypointTaskTransferFleet}
		<select
			class="select select-outline select-secondary select-sm py-0 text-sm mt-1"
			value={waypoint.transferToPlayer}
			onchange={(e) => onTransferToPlayerChanged(parseInt(e.currentTarget.value))}
		>
			<option value={undefined}>None</option>
			{#each $game.players as otherPlayer (otherPlayer.num)}
				{#if otherPlayer.num !== $player.num}
					<option value={otherPlayer.num}>{$universe.getPlayerPluralName(otherPlayer.num)}</option>
				{/if}
			{/each}
		</select>
	{/if}
</CommandTile>
