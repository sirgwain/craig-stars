<script lang="ts">
	import {
		game,
		commandedFleet,
		commandedMapObjectName,
		selectedWaypoint
	} from '$lib/services/Context';
	import { FleetService } from '$lib/services/FleetService';
	import { WaypointTask } from '$lib/types/Fleet';
	import { $enum as eu } from 'ts-enum-util';
	import CommandTile from './CommandTile.svelte';

	let selectedWaypointTask: WaypointTask = WaypointTask.None;

	const fleetService = new FleetService();

	commandedMapObjectName.subscribe(
		() => $selectedWaypoint && (selectedWaypointTask = $selectedWaypoint?.task ?? WaypointTask.None)
	);

	const onSelectedWaypointTaskChange = (task: WaypointTask) => {
		if ($game && $commandedFleet && $selectedWaypoint) {
			selectedWaypointTask = task;
			$selectedWaypoint.task = task;

			fleetService.updateFleetOrders($commandedFleet);
		}
	};
</script>

{#if $commandedFleet && $selectedWaypoint}
	<CommandTile title="Waypoint Task">
		<select
			class="select select-bordered"
			bind:value={selectedWaypointTask}
			on:change={(e) =>
				onSelectedWaypointTaskChange(
					eu(WaypointTask).getValueOrDefault(e.currentTarget.value, WaypointTask.None)
				)}
		>
			{#each eu(WaypointTask).getValues() as task}
				{#if task === WaypointTask.None}
					<option value={task}>None</option>
				{:else}
					<option value={task}>{eu(WaypointTask).getValueOrDefault(task, 'None')}</option>
				{/if}
			{/each}
		</select>

		{#if selectedWaypointTask != WaypointTask.None}
			<div class="flex justify-between my-1 btn-group">
				<!-- Task items -->
			</div>
		{/if}
	</CommandTile>
{/if}
