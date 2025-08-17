<script lang="ts">
	import { enumToString } from '$lib/types/Enums';
	import { QuestionMarkCircle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { startCase } from 'lodash-es';
	import type { HTMLSelectAttributes } from 'svelte/elements';
	import { $enum as eu } from 'ts-enum-util';
	import type { StringKeyOf } from 'ts-enum-util/dist/types/types';

	// enums are strings or numbers
	type T = $$Generic<Record<StringKeyOf<E>, number | string>>;
	type TT = T[Extract<keyof T, string>];

	type Props<T> = {
		name: string;
		title?: string | undefined;
		tooltip?: string | undefined;
		enumType: T;
		titleClass?: string;
		typeTitle?: (type: TT) => string;
		typeFilter?: (type: TT) => boolean;
		showEmpty?: boolean;
	} & HTMLSelectAttributes;

	let {
		name,
		value = $bindable(),
		title = startCase(name),
		tooltip,
		enumType,
		titleClass = 'label-text w-32 text-right',
		typeTitle = (type: TT) => startCase(enumToString(enumType, type)),
		typeFilter = (_: TT) => true,
		showEmpty = false,
		...rest
	}: Props<T> = $props();
</script>

<div class="w-full flex-grow">
	<label class="label"
		><span class={titleClass}>{title}</span>
		<select class="select input-bordered ml-2 flex-grow" {name} bind:value {...rest}>
			{#each eu(enumType).getValues() as type (type)}
				{#if typeFilter(type) && (showEmpty || type !== 0)}
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
