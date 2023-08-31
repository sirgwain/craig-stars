<script lang="ts" context="module">
	export type SpinnerNumberEvent = {
		change: number;
	};
</script>

<script lang="ts">
	import { clamp } from '$lib/services/Math';
	import { ChevronDown, ChevronUp } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher<SpinnerNumberEvent>();

	export let value: number;
	export let step = 1;
	export let min = 0;
	export let max = 100;
	export let unit = '';

	function increase() {
		value = clamp(value + step, min, max);
		dispatch('change', value);
	}

	function decrease() {
		value = clamp(value - step, min, max);
		dispatch('change', value);
	}
</script>

<div class="inline-block">
	<div class="flex flex-row">
		<div class="my-auto text-primary font-semibold mr-1">
			{value}
		</div>
		<div class="my-auto mr-1">
			{unit}
		</div>
		<div class="flex flex-col">
			<button class="btn btn-xs" on:click={increase}>
				<Icon src={ChevronUp} size="12" class="hover:stroke-accent" />
			</button>
			<button class="btn btn-xs" on:click={decrease}>
				<Icon src={ChevronDown} size="12" class="hover:stroke-accent" />
			</button>
		</div>
	</div>
</div>
