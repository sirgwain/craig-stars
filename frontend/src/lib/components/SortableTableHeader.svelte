<script lang="ts">
	import type { SvelteTableColumn } from '@hurtigruten/svelte-table';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { ArrowsUpDown, ArrowUp, ArrowDown } from '@steeze-ui/heroicons';
	import { createEventDispatcher } from 'svelte';
	const dispatch = createEventDispatcher();

	export let column: SvelteTableColumn;
	export let isSorted: boolean = false;
	export let sortDescending: boolean = false;
</script>

<div class="h-full">
	{#if column.sortable ?? true}
		<button
			class="hover:text-accent cursor-pointer select-none"
			on:click={() =>
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
