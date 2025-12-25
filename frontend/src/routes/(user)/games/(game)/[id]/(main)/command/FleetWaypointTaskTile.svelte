<script lang="ts">
	import MineralMini from '$lib/components/game/MineralMini.svelte';
	import OtherMapObjectsHere from '$lib/components/game/OtherMapObjectsHere.svelte';
	import WarpSpeedGauge from '$lib/components/game/WarpSpeedGauge.svelte';
	import type {
		ChangeWaypointProps,
		ShowTransportTasksDialogEventProps
	} from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import { ReportAgeUnexplored } from '$lib/types/Consts';
	import {
		MapObjectTargetSchema,
		MapObjectType,
		WaypointTask,
		type TransportPlan
	} from '$lib/types/cs-proto';
	import { enumToString } from '$lib/types/Enums';
	import { emptyTransportTasks, WaypointTasks } from '$lib/types/Fleet';
	import { owned, ownedBy, type MapObjectLike } from '$lib/types/MapObject';
	import { getMineralOutput } from '$lib/types/Planet';
	import { emptyVector } from '$lib/types/Vector';
	import { create } from '@bufbuild/protobuf';
	import { PencilSquare } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import TransportTasksMini from '../../(plans)/transport-plans/TransportTasksMini.svelte';
	import CommandTile from './CommandTile.svelte';

	const { game, player, universe, commandedFleet, selectedWaypoint, currentSelectedWaypointIndex } =
		getGameContext();

	type Props = ChangeWaypointProps & ShowTransportTasksDialogEventProps;

	let { onShowTransportTasksDialog, onChangeWaypoint }: Props = $props();

	let selectedWaypointTask = $derived($selectedWaypoint?.task || WaypointTask.UNSPECIFIED);
	let selectedWaypointPlanet = $derived(
		$selectedWaypoint &&
			$selectedWaypoint.mapObjectTarget?.targetType === MapObjectType.PLANET &&
			$selectedWaypoint.mapObjectTarget.targetNum
			? $universe.getPlanet($selectedWaypoint.mapObjectTarget.targetNum)
			: undefined
	);

	const onSelectedWaypointTaskChange = (task: WaypointTask) => {
		if (!$selectedWaypoint || !$commandedFleet) {
			return;
		}
		$selectedWaypoint.task = task;

		if (task != WaypointTask.TRANSPORT) {
			// if we aren't doing a transport, reset the transport tasks to blank.
			// If we don't do this, the user could pick a transport task in the future and assume it defaults to empty
			// but it will have whatever it last had
			$selectedWaypoint.transportTasks = emptyTransportTasks();
		}

		onChangeWaypoint?.({
			fleet: $commandedFleet,
			waypoint: $selectedWaypoint,
			waypointIndex: $currentSelectedWaypointIndex
		});
	};

	function onPatrolRangeChanged(value: number) {
		if (!$selectedWaypoint || !$commandedFleet) {
			return;
		}
		$selectedWaypoint.patrolRange = value;
		onChangeWaypoint?.({
			fleet: $commandedFleet,
			waypoint: $selectedWaypoint,
			waypointIndex: $currentSelectedWaypointIndex
		});
	}

	function onPatrolWarpSpeedChanged(warpSpeed: number) {
		if (!$selectedWaypoint || !$commandedFleet) {
			return;
		}
		$selectedWaypoint.patrolWarpSpeed = warpSpeed;
		onChangeWaypoint?.({
			fleet: $commandedFleet,
			waypoint: $selectedWaypoint,
			waypointIndex: $currentSelectedWaypointIndex
		});
	}

	function onPatrolWarpSpeedDragged(warpSpeed: number) {
		if (!$selectedWaypoint) {
			return;
		}
		$selectedWaypoint.patrolWarpSpeed = warpSpeed;
	}

	function onLayMinefieldDurationChanged(value: number | undefined) {
		if (!$selectedWaypoint || !$commandedFleet) {
			return;
		}
		$selectedWaypoint.layMinefieldDuration = value ?? $selectedWaypoint.layMinefieldDuration;
		onChangeWaypoint?.({
			fleet: $commandedFleet,
			waypoint: $selectedWaypoint,
			waypointIndex: $currentSelectedWaypointIndex
		});
	}

	function onTransferToPlayerChanged(playerNum: number) {
		if (!$selectedWaypoint || !$commandedFleet) {
			return;
		}
		$selectedWaypoint.transferToPlayer = playerNum;
		onChangeWaypoint?.({
			fleet: $commandedFleet,
			waypoint: $selectedWaypoint,
			waypointIndex: $currentSelectedWaypointIndex
		});
	}

	function applyTransportPlan(plan: TransportPlan | undefined) {
		if (!plan || !$selectedWaypoint || !$commandedFleet) {
			return;
		}

		$selectedWaypoint.transportTasks = plan.tasks;
		onChangeWaypoint?.({
			fleet: $commandedFleet,
			waypoint: $selectedWaypoint,
			waypointIndex: $currentSelectedWaypointIndex
		});
	}

	function onTargetChanged(target: Partial<MapObjectLike>) {
		if (!$selectedWaypoint || !$commandedFleet) {
			return;
		}
		$selectedWaypoint.mapObjectTarget = create(MapObjectTargetSchema, {
			targetName: target.mapObject?.name ?? '',
			targetType: target.mapObject?.type ?? MapObjectType.UNSPECIFIED,
			targetNum: target.mapObject?.num ?? 0,
			targetPlayerNum: target.mapObject?.playerNum ?? 0
		});

		onChangeWaypoint?.({
			fleet: $commandedFleet,
			waypoint: $selectedWaypoint,
			waypointIndex: $currentSelectedWaypointIndex
		});
	}
</script>

{#if $commandedFleet && $selectedWaypoint}
	<CommandTile title="Waypoint Task">
		<div class="flex justify-between">
			<div class="my-auto text-tile-item-title">Target</div>
			<div>
				<OtherMapObjectsHere
					fleet={$commandedFleet}
					otherMapObjectsHere={$universe.getOtherMapObjectsHereByType(
						$selectedWaypoint.position ?? emptyVector()
					)}
					target={$selectedWaypoint.mapObjectTarget}
					position={$selectedWaypoint.position}
					class="w-36"
					onSelected={onTargetChanged}
				/>
			</div>
		</div>
		<div class="flex justify-between my-1">
			<div class="my-auto text-tile-item-title">Task</div>
			<div>
				<select
					data-type="select-waypoint-task"
					data-id="waypoint-task"
					class="select select-outline select-secondary select-sm text-sm w-36"
					value={selectedWaypointTask}
					onchange={(e) => {
						onSelectedWaypointTaskChange(parseInt(e.currentTarget.value));
					}}
				>
					{#each WaypointTasks as task (task)}
						{#if task === WaypointTask.UNSPECIFIED}
							<option value={task}>None</option>
						{:else}
							<option value={task}>{enumToString(WaypointTask, task)}</option>
						{/if}
					{/each}
				</select>
			</div>
		</div>

		{#if $selectedWaypoint.task === WaypointTask.TRANSPORT}
			<div class="flex flex-col">
				<div>
					<TransportTasksMini transportTasks={$selectedWaypoint.transportTasks} />
				</div>
				<div class="ml-auto mt-1 flex flex-row gap-1">
					<div>
						<button
							type="button"
							onclick={() =>
								onShowTransportTasksDialog?.({
									fleet: $commandedFleet,
									waypoint: $selectedWaypoint,
									waypointIndex: $currentSelectedWaypointIndex
								})}
							class="btn btn-outline btn-sm normal-case btn-secondary inline-block p-1"
							><Icon src={PencilSquare} size="16" class="hover:stroke-accent inline" /></button
						>
					</div>
					<select
						class="select select-outline select-sm select-secondary w-12 sm:w-full text-secondary"
						onchange={(e) => {
							applyTransportPlan(
								$player.playerPlans.transportPlans.find(
									(p) => p.num == parseInt(e.currentTarget.value)
								)
							);
							e.currentTarget.value = '0';
						}}
					>
						<option value={0}>Apply Plan</option>
						{#each $player.playerPlans.transportPlans as plan (plan.num)}
							<option value={plan.num}>{plan.name}</option>
						{/each}
					</select>
				</div>
			</div>
		{:else if $selectedWaypoint.task === WaypointTask.REMOTE_MINING}
			{#if selectedWaypointPlanet}
				<!-- if this waypoint is owned -->
				{#if selectedWaypointPlanet.mapObject?.reportAge === ReportAgeUnexplored}
					<span class="text-warning"
						>Warning: This planet is unexplored. We have no way of knowing if we can mine it.</span
					>
				{:else if owned(selectedWaypointPlanet) && !($player.race.spec.canRemoteMineOwnPlanets && ownedBy(selectedWaypointPlanet, $player.num))}
					<span class="text-error">Note: You can only remote mine unoccupied planets.</span>
				{:else if !$commandedFleet.spec.shipDesignSpec?.miningRate}
					<span class="text-error"
						>Warning: This fleet contains no ships with remote mining modules.</span
					>
				{:else}
					Mining Rate per Year:
					<MineralMini
						mineral={getMineralOutput(
							selectedWaypointPlanet,
							$commandedFleet.spec.shipDesignSpec?.miningRate ?? 0,
							$game.rules.remoteMiningMineOutput
						)}
						showUnits={true}
					/>
				{/if}
			{:else}
				<span class="text-error">Warning: Can only remote mine planets.</span>
			{/if}
		{:else if $selectedWaypoint.task === WaypointTask.LAY_MINEFIELD}
			<select
				class="select select-outline select-secondary select-sm py-0 text-sm mt-1"
				value={$selectedWaypoint.layMinefieldDuration}
				onchange={(e) => onLayMinefieldDurationChanged(parseInt(e.currentTarget.value))}
			>
				<option value={0}>Indefinitely</option>
				<option value={1}>for 1 year</option>
				<option value={2}>for 2 years</option>
				<option value={3}>for 3 years</option>
				<option value={4}>for 4 years</option>
				<option value={5}>for 5 years</option>
			</select>
			<p class="text-warning">
				This fleet can lay {$commandedFleet.getTotalMinesLaidPerYear()} mines per year.
			</p>
		{:else if $selectedWaypoint.task === WaypointTask.PATROL}
			<div class="flex justify-between my-1">
				<div class="my-auto text-tile-item-title">Intercept</div>
				<div>
					<select
						class="select select-outline select-secondary select-sm py-0 text-sm mt-1"
						value={$selectedWaypoint.patrolRange}
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
						<option value={0}>any enemy</option>
					</select>
				</div>
			</div>
			<div class="flex mt-1">
				<span class="text-tile-item-title">Warp Factor</span>
				<span class="flex-1 ml-1">
					<WarpSpeedGauge
						onValueChanged={(value) => onPatrolWarpSpeedChanged(value)}
						onValueDragged={(value) => onPatrolWarpSpeedDragged(value)}
						value={$selectedWaypoint.patrolWarpSpeed}
						warnSpeed={$commandedFleet.spec.shipDesignSpec?.engine?.maxSafeSpeed
							? $commandedFleet.spec.shipDesignSpec.engine.maxSafeSpeed + 1
							: undefined}
						warp0Text="Automatic"
					/>
				</span>
			</div>
		{:else if $selectedWaypoint.task === WaypointTask.TRANSFER_FLEET}
			<select
				class="select select-outline select-secondary select-sm py-0 text-sm mt-1"
				value={$selectedWaypoint.transferToPlayer}
				onchange={(e) => onTransferToPlayerChanged(parseInt(e.currentTarget.value))}
			>
				<option value={undefined}>None</option>
				{#each $game.players as otherPlayer (otherPlayer.num)}
					{#if otherPlayer.num !== $player.num}
						<option value={otherPlayer.num}>{$universe.getPlayerPluralName(otherPlayer.num)}</option
						>
					{/if}
				{/each}
			</select>
		{/if}
	</CommandTile>
{/if}
