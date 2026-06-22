<script lang="ts">
	import {
		WaypointTaskTransportAction,
		type WaypointTransportTask,
		type WaypointTransportTasks,
		WaypointTransportTaskSchema
	} from '$lib/types/cs-proto';
	import { create } from '@bufbuild/protobuf';
	import TransportTasks from './TransportTask.svelte';

	type Props = {
		transportTasks: WaypointTransportTasks;
	};

	let { transportTasks = $bindable() }: Props = $props();

	// Helper to produce a WaypointTransportTask-like object with defaults.
	// Note: The generated Message type requires internal fields at runtime,
	// but for UI state we can keep a plain object with the same shape,
	// then assign it back to transportTasks (which is also Message-typed with optional fields).
	function makeTask(task: WaypointTransportTask | undefined = undefined): WaypointTransportTask {
		return create(WaypointTransportTaskSchema, {
			// amount defaults to 0, action defaults to empty string
			amount: task?.amount ?? 0,
			action: task?.action ?? WaypointTaskTransportAction.UNSPECIFIED
		});
	}

	// Svelte 5 runes local state per resource
	let fuelTask: WaypointTransportTask = $state(makeTask(transportTasks.fuel));
	let ironiumTask: WaypointTransportTask = $state(makeTask(transportTasks.ironium));
	let boraniumTask: WaypointTransportTask = $state(makeTask(transportTasks.boranium));
	let germaniumTask: WaypointTransportTask = $state(makeTask(transportTasks.germanium));
	let colonistsTask: WaypointTransportTask = $state(makeTask(transportTasks.colonists));

	// Mirror local state back to parent using runes-compatible effect
	$effect(() => {
		transportTasks.fuel = fuelTask;
		transportTasks.ironium = ironiumTask;
		transportTasks.boranium = boraniumTask;
		transportTasks.germanium = germaniumTask;
		transportTasks.colonists = colonistsTask;
	});
</script>

<div class="grid grid-cols-4 mt-2">
	<!-- headers -->
	<div></div>
	<div class="text-center font-semibold col-span-2">Action</div>
	<div class="text-center font-semibold">Amount (kT)</div>

	<TransportTasks
		title="Fuel"
		textClass="text-fuel"
		bind:action={fuelTask.action}
		bind:amount={fuelTask.amount}
	/>
	<TransportTasks
		title="Ironium"
		textClass="text-ironium"
		bind:action={ironiumTask.action}
		bind:amount={ironiumTask.amount}
	/>
	<TransportTasks
		title="Boranium"
		textClass="text-boranium"
		bind:action={boraniumTask.action}
		bind:amount={boraniumTask.amount}
	/>
	<TransportTasks
		title="Germanium"
		textClass="text-germanium"
		bind:action={germaniumTask.action}
		bind:amount={germaniumTask.amount}
	/>
	<TransportTasks
		title="Colonists"
		textClass="text-colonists"
		bind:action={colonistsTask.action}
		bind:amount={colonistsTask.amount}
	/>
	{#if colonistsTask.amount && (colonistsTask.action === WaypointTaskTransportAction.LOAD_AMOUNT || colonistsTask.action === WaypointTaskTransportAction.UNLOAD_AMOUNT)}
		<div class="col-start-2 col-span-3 ml-2">
			<span class="italic">{(colonistsTask.amount * 100).toLocaleString()} colonists</span>
		</div>
	{/if}
</div>
