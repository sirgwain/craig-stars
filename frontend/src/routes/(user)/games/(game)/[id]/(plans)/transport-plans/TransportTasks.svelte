<script lang="ts">
	import {
		TransportActionLoadAmount,
		TransportActionUnloadAmount,
		type WaypointTransportTasks
	} from '$lib/types/cs';
	import TransportTasks from './TransportTask.svelte';

	type Props = {
		transportTasks: WaypointTransportTasks;
	};

	let { transportTasks = $bindable() }: Props = $props();
</script>

<div class="grid grid-cols-4 mt-2">
	<!-- headers -->
	<div></div>
	<div class="text-center font-semibold col-span-2">Action</div>
	<div class="text-center font-semibold">Amount (kT)</div>
	<TransportTasks
		title="Fuel"
		textClass="text-fuel"
		bind:action={transportTasks.fuel.action}
		bind:amount={transportTasks.fuel.amount}
	/>
	<TransportTasks
		title="Ironium"
		textClass="text-ironium"
		bind:action={transportTasks.ironium.action}
		bind:amount={transportTasks.ironium.amount}
	/>
	<TransportTasks
		title="Boranium"
		textClass="text-boranium"
		bind:action={transportTasks.boranium.action}
		bind:amount={transportTasks.boranium.amount}
	/>
	<TransportTasks
		title="Germanium"
		textClass="text-germanium"
		bind:action={transportTasks.germanium.action}
		bind:amount={transportTasks.germanium.amount}
	/>
	<TransportTasks
		title="Colonists"
		textClass="text-colonists"
		bind:action={transportTasks.colonists.action}
		bind:amount={transportTasks.colonists.amount}
	/>
	{#if transportTasks.colonists.amount && (transportTasks.colonists.action === TransportActionLoadAmount || transportTasks.colonists.action === TransportActionUnloadAmount)}
		<div class="col-start-2 col-span-3 ml-2">
			<span class="italic"
				>{(transportTasks.colonists.amount * 100).toLocaleString()} colonists</span
			>
		</div>
	{/if}
</div>
