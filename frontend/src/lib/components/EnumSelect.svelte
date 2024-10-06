<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { startCase } from 'lodash-es';
	import { $enum as eu } from 'ts-enum-util';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { QuestionMarkCircle } from '@steeze-ui/heroicons';

	const dispatch = createEventDispatcher();

	export let name: string;
	export let value: string | undefined;
	export let title: string | undefined = undefined;
	export let tooltip: string | undefined = undefined;
	export let enumType: any;
	export let titleClass = 'label-text w-32 text-right';
	export let required = false;
	export let typeTitle = (type: any) => startCase(type);
	export let showEmpty = false;

	$: !title && (title = startCase(name));
</script>

<div class="w-full flex-grow">
	<label class="label"
		><span class={titleClass}>{title}</span>
		<select
			class="select input-bordered ml-2 flex-grow"
			{name}
			{required}
			bind:value
			on:change={(e) => dispatch('change', e)}
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
