<script lang="ts">
	import SpinnerNumber, { type SpinnerNumberEvent } from '$lib/components/SpinnerNumber.svelte';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher<SpinnerNumberEvent>();

	interface Props {
		value: number;
		step?: number;
		min?: number;
		max?: number;
		unit?: string;
		begin?: import('svelte').Snippet;
		end?: import('svelte').Snippet;
		[key: string]: any
	}

	let {
		value = $bindable(),
		step = 1,
		min = 0,
		max = 100,
		unit = '',
		begin,
		end,
		...rest
	}: Props = $props();
</script>

<div class="flex flex-row gap-1" {...rest}>
	<div class="my-auto align-middle">
		{@render begin?.()}
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
		{@render end?.()}
	</div>
</div>
