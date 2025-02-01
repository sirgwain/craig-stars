<script lang="ts">
	import TextInput from '$lib/components/TextInput.svelte';
	import type { TransportPlan } from '$lib/types/Player';
	import TransportTasks from './TransportTasks.svelte';

	type Props = {
		plan: TransportPlan;
	};

	let { plan = $bindable() }: Props = $props();
</script>

<TextInput name="name" bind:value={plan.name} required />
<TransportTasks bind:transportTasks={plan.tasks} />
