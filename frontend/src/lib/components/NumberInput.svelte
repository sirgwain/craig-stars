<script lang="ts">
	import { startCase } from 'lodash-es';
	import type { HTMLInputAttributes } from 'svelte/elements';

	type Props = {
		name: string;
		value: number | undefined;
		unit?: string | undefined;
		title?: string | undefined;
		titleClass?: string;
		unitLabelClass?: string;
	} & HTMLInputAttributes;

	let {
		name,
		value = $bindable(),
		unit = undefined,
		titleClass = 'cs-label-text w-32 text-right',
		unitLabelClass = 'w-16',
		...rest
	}: Props = $props();

	let title = $derived(rest.title ?? startCase(name));
</script>

<div class="w-full grow">
	<label class="cs-form-label"
		><span class={titleClass}>{title}</span>
		<div class="grow pl-2">
			<div class="flex items-center">
				<input type="number" class="input w-full" {name} bind:value {...rest} />
				{#if unit}
					<span class={unitLabelClass}>{unit}</span>
				{/if}
			</div>
		</div>
	</label>
</div>
