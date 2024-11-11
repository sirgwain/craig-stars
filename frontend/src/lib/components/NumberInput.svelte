<script lang="ts">
	import { run } from 'svelte/legacy';

	import { startCase } from 'lodash-es';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	type Props = {
		name: string;
		value: number | undefined;
		unit?: string | undefined;
		title?: string | undefined;
		titleClass?: string;
		inputClass?: string;
		step?: number;
		min?: number;
		max?: number | undefined;
		unitLabelClass?: string;
		required?: boolean;
		disabled?: boolean;
	};

	let {
		name,
		value = $bindable(),
		unit = undefined,
		title = $bindable(undefined),
		titleClass = 'label-text w-32 text-right',
		inputClass = 'input input-bordered w-full',
		step = 0.01,
		min = 0,
		max = undefined,
		unitLabelClass = 'w-16',
		required = false,
		disabled = false
	}: Props = $props();

	run(() => {
		!title && (title = startCase(name));
	});
</script>

<div class="w-full flex-grow">
	<label class="label"
		><span class={titleClass}>{title}</span>
		<div class="flex-grow pl-2">
			<div class="input-group">
				<input
					class={inputClass}
					type="number"
					{disabled}
					{name}
					{min}
					{max}
					{step}
					{required}
					bind:value
					onchange={(e) => dispatch('change', e)}
				/>
				{#if unit}
					<span class={unitLabelClass}>{unit}</span>
				{/if}
			</div>
		</div>
	</label>
</div>
