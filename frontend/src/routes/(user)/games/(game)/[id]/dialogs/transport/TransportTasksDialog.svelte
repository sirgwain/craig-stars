<script lang="ts">
	import type {
		OnCancel,
		OnOk,
		TransportTasksDialogEvent,
		ChangeWaypointTransportTasksEvent
	} from '$lib/services/Events';
	import { emptyTransportTasks, type WaypointTransportTasks } from '$lib/types/Fleet';
	import TransportTasks from '../../(plans)/transport-plans/TransportTasks.svelte';

	type Props = {
		show?: boolean;
		props: TransportTasksDialogEvent | undefined;
		onOk: OnOk<ChangeWaypointTransportTasksEvent>;
		onCancel: OnCancel;
	};

	let { show = false, props, onOk, onCancel }: Props = $props();

	let transportTasks: WaypointTransportTasks = $state(emptyTransportTasks());

	$effect(() => {
		if (!show) return;
		// use the waypoint props transport tasks when opening the dialog
		// this is a bit weird, but props.waypoint.transportTasks is also a "state" because it comes
		// from the FleetWaypointsTile, so take a snapshot of it instead of just grabbing the existing state var
		transportTasks = $state.snapshot(props?.waypoint.transportTasks) ?? emptyTransportTasks();
	});

	function ok() {
		if (!props) return;

		onOk({
			fleet: props.fleet,
			waypoint: props.waypoint,
			waypointIndex: props.waypointIndex,
			transportTasks
		});
	}

	function cancel() {
		transportTasks = emptyTransportTasks();
		onCancel();
	}
</script>

<div class="modal" class:modal-open={show}>
	<div class="modal-box max-w-full max-h-max h-full w-full lg:max-w-[40rem] lg:max-h-[32rem]">
		{#if show}
			<div class="flex flex-row justify-center px-1 w-full">
				<div class="flex flex-col grow">
					<div class="text-xl font-semibold w-full text-center">Transport Orders</div>
					<TransportTasks bind:transportTasks />
				</div>
				<div class="flex flex-col mt-7 ml-2 gap-2">
					<button
						onclick={(e) => {
							e.preventDefault();
							ok();
						}}
						type="submit"
						class="btn btn-sm normal-case btn-primary">OK</button
					>
					<button
						onclick={cancel}
						type="button"
						class="btn btn-outline btn-sm normal-case btn-secondary">Cancel</button
					>
				</div>
			</div>
		{/if}
	</div>
</div>
