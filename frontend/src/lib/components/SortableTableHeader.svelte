<script lang="ts">
	import type { SvelteTableColumn } from '@hurtigruten/svelte-table';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { ArrowsUpDown, ArrowUp, ArrowDown } from '@steeze-ui/heroicons';

	export let column: SvelteTableColumn;
	export let isSorted: boolean = false;
	export let sortDescending: boolean = false;
</script>

{#if column.sortable ?? true}
	<th class="hover:text-accent cursor-pointer select-none"
		>{column.title}
		{#if isSorted}
			{#if sortDescending}
				<Icon src={ArrowUp} size="16" class="hover:stroke-accent inline-block" />
			{:else}
				<Icon src={ArrowDown} size="16" class="hover:stroke-accent inline-block" />
			{/if}
		{:else}
			<Icon src={ArrowsUpDown} size="16" class="hover:stroke-accent inline-block" />
		{/if}
	</th>
{:else}
	<th>{column.title}</th>
{/if}
