<script lang="ts">
	import { run } from 'svelte/legacy';

	import { createEventDispatcher } from 'svelte';
	import { startCase } from 'lodash-es';
	import { $enum as eu } from 'ts-enum-util';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { QuestionMarkCircle } from '@steeze-ui/heroicons';

	const dispatch = createEventDispatcher();

	type Props = {
		name: string;
		value: string | undefined;
		title?: string | undefined;
		tooltip?: string | undefined;
		enumType: any;
		titleClass?: string;
		required?: boolean;
		typeTitle?: any;
		showEmpty?: boolean;
	};

	let {
		name,
		value = $bindable(),
		title = $bindable(undefined),
		tooltip = undefined,
		enumType,
		titleClass = 'label-text w-32 text-right',
		required = false,
		typeTitle = (type: any) => startCase(type),
		showEmpty = false
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
			{name}
			{required}
			bind:value
			onchange={(e) => dispatch('change', e)}
		>
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
