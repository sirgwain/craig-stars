<script lang="ts">
	import SpinnerNumber from '$lib/components/SpinnerNumber.svelte';
	import { type Snippet } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	type Props = {
		value: number;
		step?: number;
		min?: number;
		max?: number;
		unit?: string;
		begin?: Snippet;
		end?: Snippet;
		onChange?: (value: number) => void;
	} & Omit<HTMLAttributes<HTMLDivElement>, 'onchange'>;

	let {
		value = $bindable(),
		step = 1,
		min = 0,
		max = 100,
		unit = '',
		begin,
		end,
		onChange: onChange,
		...rest
	}: Props = $props();
</script>

<div class="flex flex-row gap-1" {...rest}>
	<div class="my-auto align-middle">
		{@render begin?.()}
	</div>
	<SpinnerNumber bind:value {onChange} {min} {max} {step} {unit} />
	<div class="my-auto align-middle">
		{@render end?.()}
	</div>
</div>
