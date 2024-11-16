<script lang="ts">
	import { QuestionMarkCircle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { startCase } from 'lodash-es';
	import type { HTMLSelectAttributes } from 'svelte/elements';
	import { $enum as eu } from 'ts-enum-util';

	type Props = {
		name: string;
		title?: string | undefined;
		tooltip?: string | undefined;
		enumType: any;
		titleClass?: string;
		typeTitle?: (type: any) => string;
		showEmpty?: boolean;
	} & HTMLSelectAttributes;

	let {
		name,
		value = $bindable(),
		title = startCase(name),
		tooltip,
		enumType,
		titleClass = 'label-text w-32 text-right',
		typeTitle = (type: any) => startCase(type),
		showEmpty = false,
		...others
	}: Props = $props();
</script>

<div class="w-full flex-grow">
	<label class="label"
		><span class={titleClass}>{title}</span>
		<select class="select input-bordered ml-2 flex-grow" {name} bind:value {...others}>
			{#each eu(enumType).getValues() as type}
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
