<script lang="ts">
	import { clamp } from '$lib/services/Math';
	import { ChevronDown, ChevronUp } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	type Props = {
		value: number;
		step?: number;
		min?: number;
		max?: number;
		unit?: string;
		onchange?: (value: number) => void;
	};

	let { value = $bindable(), step = 1, min = 0, max = 100, unit = '', onchange }: Props = $props();

	function increase(e: MouseEvent | PointerEvent) {
		value = clamp(
			value + step * (e.shiftKey ? 10 : 1) * (e.metaKey || e.ctrlKey ? 100 : 1),
			min,
			max
		);
		onchange?.(value);
	}

	function decrease(e: MouseEvent | PointerEvent) {
		value = clamp(
			value - step * (e.shiftKey ? 10 : 1) * (e.metaKey || e.ctrlKey ? 100 : 1),
			min,
			max
		);
		onchange?.(value);
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
			<button type="button" class="btn btn-xs" onclick={increase}>
				<Icon src={ChevronUp} size="12" class="hover:stroke-accent" />
			</button>
			<button type="button" class="btn btn-xs" onclick={decrease}>
				<Icon src={ChevronDown} size="12" class="hover:stroke-accent" />
			</button>
		</div>
	</div>
</div>
