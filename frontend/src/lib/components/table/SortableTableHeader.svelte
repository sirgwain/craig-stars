<script lang="ts">
	import { Icon } from '@steeze-ui/svelte-icon';
	import { ArrowsUpDown, ArrowUp, ArrowDown } from '@steeze-ui/heroicons';
	import { createEventDispatcher } from 'svelte';
	import type { TableColumn } from './Table.svelte';
	const dispatch = createEventDispatcher();

	type T = $$Generic;
	type Props = {
		column: TableColumn<T>;
		isSorted?: boolean;
		sortDescending?: boolean;
	};

	let { column, isSorted = false, sortDescending = false }: Props = $props();
</script>

<div class="h-full">
	{#if column.sortable ?? true}
		<button
			class="hover:text-accent cursor-pointer select-none"
			onclick={() =>
				dispatch('sorted', { sortDescending: isSorted ? !sortDescending : false, column })}
		>
			{column.title}
			{#if isSorted}
				{#if sortDescending}
					<Icon src={ArrowUp} size="16" class="hover:stroke-accent inline-block" />
				{:else}
					<Icon src={ArrowDown} size="16" class="hover:stroke-accent inline-block" />
				{/if}
			{:else}
				<Icon src={ArrowsUpDown} size="16" class="hover:stroke-accent inline-block" />
			{/if}
		</button>
	{:else}
		{column.title}
	{/if}
</div>
