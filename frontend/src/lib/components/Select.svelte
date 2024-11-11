<script lang="ts">
	import { run } from 'svelte/legacy';

	import { createEvent } from '@testing-library/svelte';
	import { startCase } from 'lodash-es';
	import { createEventDispatcher } from 'svelte';
	const dispatch = createEventDispatcher();

	type Value = {
		value: any;
		title: string;
	};

	type Props = {
		name: string;
		value: any | undefined;
		title?: string | undefined;
		titleClass?: string;
		required?: boolean;
		values?: Value[];
	};

	let {
		name,
		value = $bindable(),
		title = $bindable(undefined),
		titleClass = 'label-text w-32 text-right',
		required = false,
		values = []
	}: Props = $props();

	run(() => {
		!title && (title = startCase(name));
	});
</script>

<div class="w-full flex-grow">
	<label class="label"
		><span class={titleClass}>{title}</span>
		<select
			class="select input-bordered ml-2 flex-grow"
			name="type"
			{required}
			bind:value
			onchange={(e) => dispatch('change', e.currentTarget.value)}
		>
			{#each values as value}
				<option value={value.value}>{value.title}</option>
			{/each}
		</select>
	</label>
</div>
