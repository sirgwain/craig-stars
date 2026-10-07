<script lang="ts">
	import { startCase } from 'lodash-es';
	import type { HTMLInputAttributes } from 'svelte/elements';

	type Props = {
		name: string;
		value: string | undefined;
		title?: string | undefined;
		titleClass?: string;
		required?: boolean;
		disabled?: boolean;
	} & HTMLInputAttributes;

	let {
		name,
		value = $bindable(),
		titleClass = 'cs-label-text w-32 text-right',
		required = false,
		disabled = false,
		...rest
	}: Props = $props();

	let title = $derived(rest.title ?? startCase(name));
</script>

<div class="w-full grow">
	<div class="cs-form-control">
		<label class="cs-form-label"
			><span class={titleClass}>{title}</span>
			<input
				class="input ml-2 grow"
				type="text"
				{name}
				{required}
				{disabled}
				bind:value
				{...rest}
			/>
		</label>
	</div>
</div>
