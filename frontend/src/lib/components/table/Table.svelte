<script lang="ts" module>
	export interface TableColumn<T> {
		key: string;
		title: string;
		sortable?: boolean;
		filterable?: boolean;
		hidden?: boolean;
		sortBy?: (a: T, b: T) => number;
		filterBy?: (value: string, row: T) => boolean;
	}
</script>

<script lang="ts">
	type T = $$Generic<Record>;
	type TableClasses = Partial<
		Record<'table' | 'thead' | 'headtr' | 'th' | 'tbody' | 'tr' | 'td', string>
	>;

	const defaultClasses: TableClasses = {
		table: '',
		headtr: '',
		thead: '',
		tbody: '',
		tr: '',
		th: '',
		td: ''
	};

	interface Props {
		classes?: TableClasses;
		columns?: TableColumn<T>[];
		rows?: T[];
		filterBy?: string;
		externalSortAndFilter?: boolean;
		head?: import('svelte').Snippet<[any]>;
		cell?: import('svelte').Snippet<[any]>;
		empty?: import('svelte').Snippet;
	}

	let {
		classes = defaultClasses,
		columns = [],
		rows = $bindable([]),
		filterBy = '',
		externalSortAndFilter = false,
		head,
		cell,
		empty
	}: Props = $props();

	let lastSortedKey = $state('');
	let sortDescending = $state(false);

	/**
	 * sort rows by a column key
	 * @param key the column key to sort by
	 * @param override true to force sort by descending
	 */
	function sortRowsBy(key: string, override = false): void {
		const columnData = columns.find((column) => column.key === key);
		if (!columnData || columnData.sortable === false) {
			return;
		}

		sortDescending = getSortingOrder(key, override);
		lastSortedKey = key;

		// call column sortBy
		if (columnData.sortBy) {
			const sortBy = columnData.sortBy;
			rows = [...rows].sort((a, b) => {
				[a, b] = sortDescending ? [a, b] : [b, a];
				return sortBy(a, b);
			});
			return;
		}

		// sort by content by default
		rows = [...rows].sort((a, b) => {
			[a, b] = [a[key], b[key]];
			if (sortDescending) [b, a] = [a, b];
			if (typeof a === 'number') return a - b;
			if (typeof a === 'boolean') return a ? -1 : 1;
			return a?.localeCompare(b);
		});
	}

	function getSortingOrder(key: any, override = false): boolean {
		if (override) return sortDescending;
		if (lastSortedKey === key) return !sortDescending;
		return false;
	}

	function filterRowsBy(value: string, rows: T[]) {
		const numColumns = columns.length;
		return rows.filter((row, rowIndex) => {
			for (let colIndex = 0; colIndex < numColumns; colIndex++) {
				const col = columns[colIndex];
				if (col.filterable === false) {
					continue;
				}

				if (col.filterBy && col.filterBy(value, row)) {
					// this column contains the text we're looking for, return the row
					return row;
				} else {
					if (`${row[col.key]}`.toLowerCase().indexOf(value) != -1) {
						// this cell contains the text we are looking fork include the row
						return row;
					}
				}
			}
		});
	}

	let filteredRows = $derived((() => {
		if (externalSortAndFilter) {
			// rows come filtered and sorted, return them as is
			return rows;
		}
		if (lastSortedKey) {
			sortRowsBy(lastSortedKey, true);
		}
		return filterRowsBy(filterBy, rows);
	})());

	let assignedClasses = $derived({ ...defaultClasses, ...classes });
</script>

<table class={assignedClasses.table} style="border-spacing: 0">
	<thead class={assignedClasses.thead}>
		<tr class={assignedClasses.headtr}>
			{#each columns as column}
				{#if !column.hidden}
					<th
						scope="col"
						class={assignedClasses.th}
						onclick={() => !externalSortAndFilter && sortRowsBy(column.key)}
					>
						{#if head}
							{@render head?.({ column, isSorted: lastSortedKey === column.key, sortDescending, sortable: column.sortable !== false, })}
						{:else}
							<span>{column.title}</span>
						{/if}
					</th>
				{/if}
			{/each}
		</tr>
	</thead>
	<tbody class={assignedClasses.tbody}>
		{#each filteredRows as row}
			<tr class={`${assignedClasses.tr}`}>
				{#each columns as column}
					{#if !column.hidden}
						<td class={assignedClasses.td}>
							{#if cell}
								{@render cell?.({ row, column, cell: row[column.key], })}
							{:else}
								<span>{row[column.key]}</span>
							{/if}
						</td>
					{/if}
				{/each}
			</tr>
		{:else}
			{@render empty?.()}
		{/each}
	</tbody>
</table>
