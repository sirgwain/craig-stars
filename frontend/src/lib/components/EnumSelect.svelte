<script lang="ts">
	import { QuestionMarkCircle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { startCase } from 'lodash-es';
	import type { HTMLSelectAttributes } from 'svelte/elements';

	// enums are strings or numbers
	type T = $$Generic;

	type Props<T> = {
		name: string;
		options: T[];
		title?: string | undefined;
		tooltip?: string | undefined;
		titleClass?: string;
		typeTitle?: (type: T) => string;
		showEmpty?: boolean;
	} & HTMLSelectAttributes;

	let {
		name,
		options,
		value = $bindable(),
		title = startCase(name),
		tooltip,
		titleClass = 'label-text w-32 text-right',
		typeTitle = (type: T) => startCase(`${type}`),
		showEmpty = false,
		...rest
	}: Props<T> = $props();
</script>

<div class="w-full flex-grow">
	<label class="label"
		><span class={titleClass}>{title}</span>
		<select class="select input-bordered ml-2 flex-grow" {name} bind:value {...rest}>
			{#each options as type (type)}
				{#if showEmpty || `${type}` !== ''}
					<option value={type}>{typeTitle(type)}</option>
				{/if}
			{/each}
		</select>
		{#if tooltip}
			<div class="tooltip tooltip-left mx-2" data-tip={tooltip}>
				<Icon src={QuestionMarkCircle} size="16" class=" cursor-help inline-block" />
			</div>
		{/if}
	</label>
</div>
