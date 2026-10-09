<script lang="ts">
	import { Icon } from '@steeze-ui/svelte-icon';
	import { ArrowDown, ArrowUp } from '@steeze-ui/heroicons';
	import SortableTableHeader from './SortableTableHeader.svelte';
	import type { TableColumn } from './Table';

	// a sortable header for a column of minerals that can also sort by ironium, boranium or germanium
	type T = $$Generic;
	type Props = {
		column: TableColumn<T>;
		sortKey: string;
		sortDescending?: boolean;
		onSorted?: (key: string, descending: boolean) => void;
	};

	let { column, sortKey, sortDescending = false, onSorted }: Props = $props();

	const minerals = [
		{ type: 'ironium', label: 'I', title: 'Ironium', class: 'text-ironium' },
		{ type: 'boranium', label: 'B', title: 'Boranium', class: 'text-boranium' },
		{ type: 'germanium', label: 'G', title: 'Germanium', class: 'text-germanium' }
	];

	function sort(key: string) {
		onSorted?.(key, sortKey === key ? !sortDescending : false);
	}
</script>

<div class="flex flex-col items-center">
	<SortableTableHeader
		{column}
		isSorted={sortKey === column.key}
		{sortDescending}
		onSorted={() => sort(String(column.key))}
	/>
	<div class="flex flex-row gap-2">
		{#each minerals as mineral (mineral.type)}
			{@const key = `${String(column.key)}.${mineral.type}`}
			<button
				type="button"
				title={`Sort by ${mineral.title}`}
				aria-label={`Sort ${column.title} by ${mineral.title}`}
				class="{mineral.class} hover:text-accent cursor-pointer select-none"
				onclick={() => sort(key)}
			>
				{mineral.label}{#if sortKey === key}<Icon
						src={sortDescending ? ArrowUp : ArrowDown}
						size="12"
						class="inline-block"
					/>{/if}
			</button>
		{/each}
	</div>
</div>
