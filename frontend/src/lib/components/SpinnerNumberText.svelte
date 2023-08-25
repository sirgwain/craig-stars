<script lang="ts">
	import SpinnerNumber, { type SpinnerNumberEvent } from '$lib/components/SpinnerNumber.svelte';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher<SpinnerNumberEvent>();

	export let value: number;
	export let step = 1;
	export let min = 0;
	export let max = 100;
	export let unit = '';
</script>

<div class="flex flex-row gap-1" {...$$restProps}>
	<div class="my-auto align-middle">
		<slot name="begin" />
	</div>
	<SpinnerNumber
		bind:value
		on:change={(e) => dispatch('change', e.detail)}
		{min}
		{max}
		{step}
		{unit}
	/>
	<div class="my-auto align-middle">
		<slot name="end" />
	</div>
</div>
