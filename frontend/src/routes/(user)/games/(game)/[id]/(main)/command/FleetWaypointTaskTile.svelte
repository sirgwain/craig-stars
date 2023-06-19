<script lang="ts">
	import DropdownButton from '$lib/components/DropdownButton.svelte';
	import MineralMini from '$lib/components/game/MineralMini.svelte';
	import { selectedWaypoint } from '$lib/services/Context';
	import { FleetService } from '$lib/services/FleetService';
	import { CommandedFleet, WaypointTask } from '$lib/types/Fleet';
	import { MapObjectType, owned, ownedBy } from '$lib/types/MapObject';
	import { Unexplored, getMineralOutput } from '$lib/types/Planet';
	import type { Player, TransportPlan } from '$lib/types/Player';
	import { startCase } from 'lodash-es';
	import { $enum as eu } from 'ts-enum-util';
	import TransportTasksMini from '../../(plans)/transport-plans/TransportTasksMini.svelte';
	import CommandTile from './CommandTile.svelte';
	import type { FullGame } from '$lib/services/FullGame';

	export let fleet: CommandedFleet;
	export let player: Player;
	export let game: FullGame;

	$: selectedWaypointTask = $selectedWaypoint?.task ?? WaypointTask.None;
	$: selectedWaypointPlanet =
		$selectedWaypoint &&
		$selectedWaypoint.targetType == MapObjectType.Planet &&
		$selectedWaypoint.targetNum
			? player.getPlanetIntel($selectedWaypoint.targetNum)
			: undefined;

	const onSelectedWaypointTaskChange = (task: WaypointTask) => {
		if ($selectedWaypoint) {
			$selectedWaypoint.task = task;

			FleetService.updateFleetOrders(fleet);
		}
	};

	function onLayMineFieldDurationChanged() {
		FleetService.updateFleetOrders(fleet);
	}

	function applyTransportPlan(plan: TransportPlan) {
		if ($selectedWaypoint) {
			$selectedWaypoint.transportTasks = plan.tasks;

			FleetService.updateFleetOrders(fleet);
		}
	}
</script>

{#if $selectedWaypoint}
	<CommandTile title="Waypoint Task">
		<select
			class="select select-outline select-secondary select-sm py-0 text-sm"
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
					<option value={task}>{startCase(eu(WaypointTask).getValueOrDefault(task, 'None'))}</option
					>
				{/if}
			{/each}
		</select>

		{#if $selectedWaypoint?.task == WaypointTask.Transport}
			<div class="flex flex-col">
				<div>
					<TransportTasksMini transportTasks={$selectedWaypoint.transportTasks} />
				</div>
				<div class="ml-auto mt-1">
					<DropdownButton
						title="Apply Plan"
						items={player.transportPlans}
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
				{:else if owned(selectedWaypointPlanet) && !(player.race.spec?.canRemoteMineOwnPlanets && ownedBy(selectedWaypointPlanet, player.num))}
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
							game.rules.remoteMiningMineOutput
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
		{:else}
			<!-- else content here -->
		{/if}
	</CommandTile>
{/if}
