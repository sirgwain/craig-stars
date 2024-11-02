<script lang="ts">
	import { startCase } from 'lodash-es';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	export let name: string;
	export let value: number | undefined;
	export let unit: string | undefined = undefined;
	export let title: string | undefined = undefined;
	export let titleClass = 'label-text w-32 text-right';
	export let inputClass = 'input input-bordered w-full';
	export let step = 0.01;
	export let min = 0;
	export let max: number | undefined = undefined;
	export let unitLabelClass = 'w-16';
	export let required = false;
	export let disabled = false;

	$: !title && (title = startCase(name));
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
					on:change={(e) => dispatch('change', e)}
				/>
				{#if unit}
					<span class={unitLabelClass}>{unit}</span>
				{/if}
			</div>
		</div>
	</label>
</div>
