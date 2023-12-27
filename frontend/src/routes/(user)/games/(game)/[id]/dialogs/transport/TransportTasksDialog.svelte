<script lang="ts" context="module">
	import type { CommandedFleet, Waypoint, WaypointTransportTasks } from '$lib/types/Fleet';

	export type TransportTasksDialogEventDetails = {
		fleet: CommandedFleet;
		waypoint: Waypoint;
	};
	export type TransportTasksUpdateEventDetails = {
		fleet: CommandedFleet;
		waypoint: Waypoint;
		transportTasks: WaypointTransportTasks;
	};
	export type TransportTasksDialogEvent = {
		'transport-tasks-dialog': TransportTasksDialogEventDetails;
	};
</script>

<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import TransportTasks from '../../(plans)/transport-plans/TransportTasks.svelte';

	const { game, player, universe } = getGameContext();

	export let show = false;
	export let props: TransportTasksDialogEventDetails | undefined;

	$: transportTasks = props?.waypoint.transportTasks;

	const onUpdateTransportTasks = async () => {
		if (props && transportTasks) {
			props.waypoint.transportTasks = transportTasks;
			await $game.updateFleetOrders(props.fleet);
		}

		// close the dialog
		show = false;
	};
</script>

<div class="modal" class:modal-open={show}>
	<div class="modal-box max-w-full max-h-max h-full w-full lg:max-w-[40rem] lg:max-h-[32rem]">
		{#if props && transportTasks}
			<div class="flex flex-row justify-center px-1 w-full">
				<div class="flex flex-col grow">
					<div class="text-xl font-semibold w-full text-center">Transport Orders</div>
					<TransportTasks bind:transportTasks />
				</div>
				<div class="flex flex-col mt-7 ml-2 gap-2">
					<button
						on:click|preventDefault={() => onUpdateTransportTasks()}
						type="submit"
						class="btn btn-sm normal-case btn-primary">OK</button
					>
					<button
						on:click={() => (show = false)}
						class="btn btn-outline btn-sm normal-case btn-secondary">Cancel</button
					>
				</div>
			</div>
		{/if}
	</div>
</div>
