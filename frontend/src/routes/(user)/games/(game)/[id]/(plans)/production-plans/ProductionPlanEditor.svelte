<script lang="ts">
	import TextInput from '$lib/components/TextInput.svelte';
	import type { ProductionPlan } from '$lib/types/Player';
	import { ExclamationTriangle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { createEventDispatcher } from 'svelte';
	import { fade } from 'svelte/transition';

	export let plan: ProductionPlan;
</script>

<TextInput name="name" bind:value={plan.name} required />
