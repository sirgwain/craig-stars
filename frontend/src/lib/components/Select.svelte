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
		values?: Value<T>[];
	} & HTMLSelectAttributes;

	let {
		name,
		value = $bindable(),
		titleClass = 'label-text w-32 text-right',
		values = [],
		...rest
	}: Props = $props();

	let title = $derived(rest.title ?? startCase(name));
</script>

<div class="w-full flex-grow">
	<label class="label"
		><span class={titleClass}>{title}</span>
		<select class="select input-bordered ml-2 flex-grow" bind:value {...rest}>
			{#each values as value, index (index)}
				<option value={value.value}>{value.title}</option>
			{/each}
		</select>
	</label>
</div>
