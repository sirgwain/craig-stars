<script lang="ts">
	import type {
		OnCancel,
		OnOk,
		TransportTasksDialogEvent,
		TransportTasksUpdateEvent
	} from '$lib/services/Events';
	import type { WaypointTransportTasks } from '$lib/types/Fleet';
	import TransportTasks from '../../(plans)/transport-plans/TransportTasks.svelte';

	interface Props {
		show?: boolean;
		props: TransportTasksDialogEvent | undefined;
		onOk: OnOk<TransportTasksUpdateEvent>;
		onCancel: OnCancel;
	}

	let { show = false, props, onOk, onCancel }: Props = $props();

	let transportTasks: WaypointTransportTasks | undefined = $state();

	$effect(() => {
		transportTasks = props?.waypoint.transportTasks;
	});

	function ok() {
		if (!transportTasks || !props) {
			return;
		}
		onOk({ fleet: props.fleet, waypoint: props.waypoint, transportTasks });
	}
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
						onclick={(e) => {
							e.preventDefault();
							ok();
						}}
						type="submit"
						class="btn btn-sm normal-case btn-primary">OK</button
					>
					<button onclick={onCancel} class="btn btn-outline btn-sm normal-case btn-secondary"
						>Cancel</button
					>
				</div>
			</div>
		{/if}
	</div>
</div>
