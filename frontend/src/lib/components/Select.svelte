<script lang="ts">
	import { startCase } from 'lodash-es';
	import type { HTMLSelectAttributes } from 'svelte/elements';

	type Value = {
		value: any;
		title: string;
	};

	type Props = {
		name: string;
		value: any | undefined;
		title?: string | undefined;
		titleClass?: string;
		values?: Value[];
	} & HTMLSelectAttributes;

	let {
		name,
		value = $bindable(),
		title = startCase(name),
		titleClass = 'label-text w-32 text-right',
		values = [],
		...others
	}: Props = $props();
</script>

<div class="w-full flex-grow">
	<label class="label"
		><span class={titleClass}>{title}</span>
		<select class="select input-bordered ml-2 flex-grow" name="type" bind:value {...others}>
			{#each values as value}
				<option value={value.value}>{value.title}</option>
			{/each}
		</select>
	</label>
</div>
