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
		onChange?: (value: number) => void;
	};

	let {
		value = $bindable(),
		step = 1,
		min = 0,
		max = 100,
		unit = '',
		onChange: onChange
	}: Props = $props();

	function increase(e: MouseEvent | PointerEvent) {
		value = clamp(
			value + step * (e.shiftKey ? 10 : 1) * (e.metaKey || e.ctrlKey ? 100 : 1),
			min,
			max
		);
		onChange?.(value);
	}

	function decrease(e: MouseEvent | PointerEvent) {
		value = clamp(
			value - step * (e.shiftKey ? 10 : 1) * (e.metaKey || e.ctrlKey ? 100 : 1),
			min,
			max
		);
		onChange?.(value);
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
			<button
				type="button"
				class="btn btn-xs"
				onclick={increase}
				data-type="spin-number-increase-button"
			>
				<Icon src={ChevronUp} size="12" class="hover:stroke-accent" />
			</button>
			<button
				type="button"
				class="btn btn-xs"
				onclick={decrease}
				data-type="spin-number-decrease-button"
			>
				<Icon src={ChevronDown} size="12" class="hover:stroke-accent" />
			</button>
		</div>
	</div>
</div>
