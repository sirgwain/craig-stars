<script lang="ts">
	import type { WaypointTransportTasks } from '#lib/types/cs-proto.js';
	import type {
		ChangeWaypointTransportTasksEvent,
		OnCancel,
		OnOk,
		TransportTasksDialogEvent
	} from '#lib/services/Events.js';
	import { emptyTransportTasks } from '#lib/types/Fleet.js';
	import TransportTasks from '../../(plans)/transport-plans/TransportTasks.svelte';
	import { untrack } from 'svelte';

	type Props = {
		show?: boolean;
		props: TransportTasksDialogEvent | undefined;
		onOk: OnOk<ChangeWaypointTransportTasksEvent>;
		onCancel: OnCancel;
	};

	let { show = false, props, onOk, onCancel }: Props = $props();

	let transportTasks: WaypointTransportTasks = $state(emptyTransportTasks());

	$effect.pre(() => {
		if (!show || !props) return;
		const waypoint = props.waypoint;
		// Start each edit with a copy before rendering the inputs. Changes to the source orders
		// while the dialog is open must not overwrite the player's draft.
		transportTasks = untrack(
			() => $state.snapshot(waypoint.transportTasks) ?? emptyTransportTasks()
		);
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

<div data-type="dialog" data-id="Transport Orders" class="modal" class:modal-open={show}>
	<div class="modal-box max-w-full max-h-max size-full lg:max-w-160 lg:max-h-128">
		{#if show}
			<div class="flex flex-col gap-2 px-1 w-full">
				<div class="text-xl font-semibold w-full text-center">Transport Orders</div>
				<div class="flex justify-end gap-2">
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
				<TransportTasks bind:transportTasks />
			</div>
		{/if}
	</div>
</div>
