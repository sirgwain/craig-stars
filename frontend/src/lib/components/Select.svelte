<script lang="ts">
	import { startCase } from 'lodash-es';
	import type { HTMLSelectAttributes } from 'svelte/elements';

	type T = $$Generic;

	type Value<T> = {
		value: T;
		title: string;
	};

	type Props = {
		name: string;
		title?: string | undefined;
		tooltip?: string | undefined;
		titleClass?: string;
		stacked?: boolean;
		values?: Value<T>[];
	} & HTMLSelectAttributes;

	let {
		name,
		value = $bindable(),
		titleClass = 'cs-label-text w-32 text-right',
		stacked = false,
		values = [],
		...rest
	}: Props = $props();

	let title = $derived(rest.title ?? startCase(name));
</script>

<div class="w-full grow">
	<label class={stacked ? 'flex flex-col gap-1' : 'cs-form-label'}
		><span class={stacked ? 'cs-label-text' : titleClass}>{title}</span>
		<select class={stacked ? 'select w-full min-w-0' : 'select ml-2 grow'} bind:value {...rest}>
			{#each values as value, index (index)}
				<option value={value.value}>{value.title}</option>
			{/each}
		</select>
	</label>
</div>
