<script lang="ts">
	import { createEvent } from '@testing-library/svelte';
	import { startCase } from 'lodash-es';
	import { createEventDispatcher } from 'svelte';
	const dispatch = createEventDispatcher();

	type Value = {
		value: any;
		title: string;
	};
	export let name: string;
	export let value: any | undefined;

	export let title: string | undefined = undefined;
	export let titleClass = 'label-text w-32 text-right';
	export let required = false;

	export let values: Value[] = [];

	$: !title && (title = startCase(name));
</script>

<div class="w-full flex-grow">
	<label class="label"
		><span class={titleClass}>{title}</span>
		<select
			class="select input-bordered ml-2 flex-grow"
			name="type"
			{required}
			bind:value
			on:change={(e) => dispatch('change', e.currentTarget.value)}
		>
			{#each values as value}
				<option value={value.value}>{value.title}</option>
			{/each}
		</select>
	</label>
</div>
