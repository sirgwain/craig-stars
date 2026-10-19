<script lang="ts">
	import { startCase } from 'lodash-es';
	import type { HTMLInputAttributes } from 'svelte/elements';

	type Props = {
		name: string;
		value: string | undefined;
		title?: string | undefined;
		titleClass?: string;
		stacked?: boolean;
		required?: boolean;
		disabled?: boolean;
	} & HTMLInputAttributes;

	let {
		name,
		value = $bindable(),
		titleClass = 'cs-label-text w-32 text-right',
		stacked = false,
		required = false,
		disabled = false,
		...rest
	}: Props = $props();

	let title = $derived(rest.title ?? startCase(name));
</script>

<div class="w-full grow">
	<div class="cs-form-control">
		<label class={stacked ? 'flex flex-col gap-1' : 'cs-form-label'}
			><span class={stacked ? 'cs-label-text' : titleClass}>{title}</span>
			<input
				class={stacked ? 'input w-full min-w-0' : 'input ml-2 grow'}
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
