<script lang="ts">
	import TextInput from '$lib/components/TextInput.svelte';
	import type { TransportPlan } from '$lib/types/Player';
	import { ExclamationTriangle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { fade } from 'svelte/transition';
	import TransportPlanTask from './TransportPlanTask.svelte';

	export let plan: TransportPlan;
	export let error: string = '';
</script>

{#if error !== ''}
	<div
		class="alert alert-error shadow-lg w-1/2 mx-auto"
		in:fade
		out:fade={{ delay: 5000 }}
		on:introend={(e) => (error = '')}
	>
		<div>
			<Icon src={ExclamationTriangle} size="24" class="hover:stroke-accent" />
			<span>{error}</span>
		</div>
	</div>
{/if}

<TextInput name="name" bind:value={plan.name} required />
<div class="grid grid-cols-3">
	<TransportPlanTask
		title="Fuel"
		textClass="text-fuel"
		bind:action={plan.tasks.fuel.action}
		bind:amount={plan.tasks.fuel.amount}
	/>
	<TransportPlanTask
		title="Ironium"
		textClass="text-ironium"
		bind:action={plan.tasks.ironium.action}
		bind:amount={plan.tasks.ironium.amount}
	/>
	<TransportPlanTask
		title="Boranium"
		textClass="text-boranium"
		bind:action={plan.tasks.boranium.action}
		bind:amount={plan.tasks.boranium.amount}
	/>
	<TransportPlanTask
		title="Germanium"
		textClass="text-germanium"
		bind:action={plan.tasks.germanium.action}
		bind:amount={plan.tasks.germanium.amount}
	/>
	<TransportPlanTask
		title="Colonists"
		textClass="text-colonists"
		bind:action={plan.tasks.colonists.action}
		bind:amount={plan.tasks.colonists.amount}
	/>
</div>
