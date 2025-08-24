<script lang="ts">
	import TextInput from '$lib/components/TextInput.svelte';
	import { WaypointTransportTasksSchema, type TransportPlan } from '$lib/types/cs-proto';
	import { create } from '@bufbuild/protobuf';
	import TransportTasks from './TransportTasks.svelte';

	type Props = {
		plan: TransportPlan;
	};

	let { plan = $bindable() }: Props = $props();

	// Local runes state for reactive bindings
	let name: string = $state(plan.name);
	let tasks = $state(plan.tasks ?? create(WaypointTransportTasksSchema)); // may be undefined; child handles defaults internally

	// Sync local state back to plan
	$effect(() => {
		plan.name = name;
		plan.tasks = tasks;
	});
</script>

<TextInput name="name" bind:value={name} required />
<TransportTasks bind:transportTasks={tasks} />
