<script lang="ts">
	import { selectedWaypoint } from '$lib/services/Context';
	import { FleetService } from '$lib/services/FleetService';
	import { WaypointTask, type Fleet } from '$lib/types/Fleet';
	import { startCase } from 'lodash-es';
	import { $enum as eu } from 'ts-enum-util';
	import CommandTile from './CommandTile.svelte';

	export let fleet: Fleet;

	const onSelectedWaypointTaskChange = (task: WaypointTask) => {
		if ($selectedWaypoint) {
			$selectedWaypoint.task = task;

			FleetService.updateFleetOrders(fleet);
		}
	};
</script>

{#if $selectedWaypoint}
	<CommandTile title="Waypoint Task">
		<select
			class="select select-outline select-secondary select-sm py-0 text-sm"
			bind:value={$selectedWaypoint.task}
			on:change={(e) =>
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

		{#if $selectedWaypoint?.task != WaypointTask.None}
			<div class="flex justify-between my-1 btn-group">
				<!-- Task items -->
			</div>
		{/if}
	</CommandTile>
{/if}
