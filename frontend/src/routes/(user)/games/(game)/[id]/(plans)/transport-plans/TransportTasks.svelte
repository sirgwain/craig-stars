<script lang="ts">
	import {
		WaypointTaskTransportAction,
		type WaypointTransportTasks,
		WaypointTransportTaskSchema
	} from '#lib/types/cs-proto.js';
	import { create } from '@bufbuild/protobuf';
	import TransportTask from './TransportTask.svelte';

	type Props = {
		transportTasks: WaypointTransportTasks;
	};

	let { transportTasks = $bindable() }: Props = $props();

	const resources = [
		{ key: 'fuel', title: 'Fuel', textClass: 'text-fuel' },
		{ key: 'ironium', title: 'Ironium', textClass: 'text-ironium' },
		{ key: 'boranium', title: 'Boranium', textClass: 'text-boranium' },
		{ key: 'germanium', title: 'Germanium', textClass: 'text-germanium' },
		{ key: 'colonists', title: 'Colonists', textClass: 'text-colonists' }
	] as const;
</script>

<div class="grid grid-cols-4 mt-2">
	<!-- headers -->
	<div></div>
	<div class="text-center font-semibold col-span-2">Action</div>
	<div class="text-center font-semibold">Amount</div>

	{#each resources as { key, title, textClass } (key)}
		<TransportTask
			{title}
			{textClass}
			bind:action={
				() => transportTasks[key]?.action ?? WaypointTaskTransportAction.UNSPECIFIED,
				(action) => {
					const task = (transportTasks[key] ??= create(WaypointTransportTaskSchema));
					task.action = action ?? WaypointTaskTransportAction.UNSPECIFIED;
				}
			}
			bind:amount={
				() => transportTasks[key]?.amount ?? 0,
				(amount) => {
					const task = (transportTasks[key] ??= create(WaypointTransportTaskSchema));
					task.amount = amount ?? 0;
				}
			}
		/>
	{/each}
</div>
