<script lang="ts">
	import DropdownButton from '$lib/components/DropdownButton.svelte';
	import OtherMapObjectsHere from '$lib/components/game/OtherMapObjectsHere.svelte';
	import MineralMini from '$lib/components/game/MineralMini.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import {
		CommandedFleet,
		emptyTransportTasks,
		WaypointTask,
		WaypointTaskTransportAction
	} from '$lib/types/Fleet';
	import { MapObjectType, owned, ownedBy, type MapObject } from '$lib/types/MapObject';
	import { Unexplored, getMineralOutput } from '$lib/types/Planet';
	import type { TransportPlan } from '$lib/types/Player';
	import { startCase } from 'lodash-es';
	import { $enum as eu } from 'ts-enum-util';
	import TransportTasksMini from '../../(plans)/transport-plans/TransportTasksMini.svelte';
	import CommandTile from './CommandTile.svelte';
	import { PencilSquare } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { createEventDispatcher } from 'svelte';
	import type { TransportTasksDialogEvent } from '../../dialogs/transport/TransportTasksDialog.svelte';
	import WarpSpeedGauge from '$lib/components/game/WarpSpeedGauge.svelte';

	const dispatch = createEventDispatcher<TransportTasksDialogEvent>();

	const { game, player, universe, selectedWaypoint, updateFleetOrders } = getGameContext();

	export let fleet: CommandedFleet;

	$: selectedWaypointTask = $selectedWaypoint?.task ?? WaypointTask.None;
	$: selectedWaypointPlanet =
		$selectedWaypoint &&
		$selectedWaypoint.targetType == MapObjectType.Planet &&
		$selectedWaypoint.targetNum
			? $universe.getPlanet($selectedWaypoint.targetNum)
			: undefined;

	const onSelectedWaypointTaskChange = (task: WaypointTask) => {
		if ($selectedWaypoint) {
			$selectedWaypoint.task = task;

			if (task != WaypointTask.Transport) {
				// if we aren't doing a transport, reset the transport tasks to blank.
				// If we don't do this, the user could pick a transport task in the future and assume it defaults to empty
				// but it will have whatever it last had
				$selectedWaypoint.transportTasks = emptyTransportTasks();
			}

			updateFleetOrders(fleet);
		}
	};

	function onPatrolRangeChanged() {
		updateFleetOrders(fleet);
	}

	async function onPatrolWarpSpeedChanged(warpSpeed: number) {
		if ($selectedWaypoint) {
			$selectedWaypoint.patrolWarpSpeed = warpSpeed;
			await updateFleetOrders(fleet);
		}
	}

	async function onPatrolWarpSpeedDragged(warpSpeed: number) {
		if ($selectedWaypoint) {
			$selectedWaypoint.patrolWarpSpeed = warpSpeed;
		}
	}

	function onLayMineFieldDurationChanged() {
		updateFleetOrders(fleet);
	}

	function onTransferToPlayerChanged() {
		updateFleetOrders(fleet);
	}

	function applyTransportPlan(plan: TransportPlan) {
		if ($selectedWaypoint) {
			$selectedWaypoint.transportTasks = plan.tasks;

			updateFleetOrders(fleet);
		}
	}

	function onTargetChanged(target: MapObject) {
		if ($selectedWaypoint) {
			$selectedWaypoint.targetName = target.name;
			$selectedWaypoint.targetType = target.type;
			$selectedWaypoint.targetNum = target.num;
			$selectedWaypoint.targetPlayerNum = target.playerNum;
			updateFleetOrders(fleet);
		}
	}
</script>

{#if $selectedWaypoint}
	<CommandTile title="Waypoint Task">
		<div class="flex justify-between">
			<div class="my-auto text-tile-item-title">Target</div>
			<div>
				<OtherMapObjectsHere
					{fleet}
					otherMapObjectsHere={$universe.getOtherMapObjectsHereByType($selectedWaypoint.position)}
					target={$selectedWaypoint}
					class="w-36"
					on:selected={(e) => onTargetChanged(e.detail)}
				/>
			</div>
		</div>
		<div class="flex justify-between my-1">
			<div class="my-auto text-tile-item-title">Task</div>
			<div>
				<select
					class="select select-outline select-secondary select-sm text-sm w-36"
					value={selectedWaypointTask}
					on:change|preventDefault={(e) =>
						onSelectedWaypointTaskChange(
							eu(WaypointTask).getValueOrDefault(e.currentTarget.value, WaypointTask.None)
						)}
				>
					{#each eu(WaypointTask).getValues() as task}
						{#if task === WaypointTask.None}
							<option value={task}>None</option>
						{:else}
							<option value={task}
								>{startCase(eu(WaypointTask).getValueOrDefault(task, 'None'))}</option
							>
						{/if}
					{/each}
				</select>
			</div>
		</div>

		{#if $selectedWaypoint?.task == WaypointTask.Transport}
			<div class="flex flex-col">
				<div>
					<TransportTasksMini transportTasks={$selectedWaypoint.transportTasks} />
				</div>
				<div class="ml-auto mt-1 flex flex-row gap-1">
					<div>
						<button
							on:click={() =>
								$selectedWaypoint &&
								dispatch('transport-tasks-dialog', { fleet, waypoint: $selectedWaypoint })}
							class="btn btn-outline btn-sm normal-case btn-secondary inline-block p-1"
							><Icon src={PencilSquare} size="16" class="hover:stroke-accent inline" /></button
						>
					</div>
					<DropdownButton
						title="Apply Plan"
						items={$player.transportPlans}
						itemTitle={(item) => item.name}
						on:selected={(e) => applyTransportPlan(e.detail)}
					/>
				</div>
			</div>
		{:else if $selectedWaypoint?.task === WaypointTask.RemoteMining}
			{#if selectedWaypointPlanet}
				<!-- if this waypoint is owned -->
				{#if selectedWaypointPlanet.reportAge == Unexplored}
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
		{:else if $selectedWaypoint?.task === WaypointTask.LayMineField}
			<select
				class="select select-outline select-secondary select-sm py-0 text-sm mt-1"
				bind:value={$selectedWaypoint.layMineFieldDuration}
				on:change|preventDefault={() => onLayMineFieldDurationChanged()}
			>
				<option value={undefined}>Indefinitely</option>
				<option value={1}>for 1 year</option>
				<option value={2}>for 2 years</option>
				<option value={3}>for 3 years</option>
				<option value={4}>for 4 years</option>
				<option value={5}>for 5 years</option>
			</select>
			<p class="text-warning">This fleet can lay {fleet.getTotalMinesLaidPerYear()} mines per year.</p>
		{:else if $selectedWaypoint?.task === WaypointTask.Patrol}
			<div class="flex justify-between my-1">
				<div class="my-auto text-tile-item-title">Intercept</div>
				<div>
					<select
						class="select select-outline select-secondary select-sm py-0 text-sm mt-1"
						bind:value={$selectedWaypoint.patrolRange}
						on:change|preventDefault={() => onPatrolRangeChanged()}
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
						on:valuechanged={(e) => onPatrolWarpSpeedChanged(e.detail)}
						on:valuedragged={(e) => onPatrolWarpSpeedDragged(e.detail)}
						bind:value={$selectedWaypoint.patrolWarpSpeed}
						warnSpeed={fleet.spec.engine.maxSafeSpeed
							? fleet.spec.engine.maxSafeSpeed + 1
							: undefined}
						warp0Text={'Automatic'}
					/>
				</span>
			</div>
		{:else if $selectedWaypoint?.task === WaypointTask.TransferFleet}
			<select
				class="select select-outline select-secondary select-sm py-0 text-sm mt-1"
				bind:value={$selectedWaypoint.transferToPlayer}
				on:change|preventDefault={() => onTransferToPlayerChanged()}
			>
				<option value={undefined}>None</option>
				{#each $game.players as otherPlayer}
					{#if otherPlayer.num != $player.num}
						<option value={otherPlayer.num}>{$universe.getPlayerPluralName(otherPlayer.num)}</option
						>
					{/if}
				{/each}
			</select>
		{:else}
			<!-- else content here -->
		{/if}
	</CommandTile>
{/if}
