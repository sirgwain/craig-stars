<script lang="ts" module>
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

	interface Props {
		value: number;
		step?: number;
		min?: number;
		max?: number;
		unit?: string;
	}

	let { value = $bindable(), step = 1, min = 0, max = 100, unit = '' }: Props = $props();

	function increase(e) {
		value = clamp(
			value + step * (e.shiftKey ? 10 : 1) * (e.metaKey || e.ctrlKey ? 100 : 1),
			min,
			max
		);
		dispatch('change', value);
	}

	function decrease(e) {
		value = clamp(
			value - step * (e.shiftKey ? 10 : 1) * (e.metaKey || e.ctrlKey ? 100 : 1),
			min,
			max
		);
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
			<button type="button" class="btn btn-xs" onclick={(e) => increase(e)}>
				<Icon src={ChevronUp} size="12" class="hover:stroke-accent" />
			</button>
			<button type="button" class="btn btn-xs" onclick={(e) => decrease(e)}>
				<Icon src={ChevronDown} size="12" class="hover:stroke-accent" />
			</button>
		</div>
	</div>
</div>
