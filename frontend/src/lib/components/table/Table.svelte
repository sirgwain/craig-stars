<script lang="ts" module>
	export interface TableColumn<T> {
		key: keyof Partial<T>;
		title: string;
		sortable?: boolean;
		filterable?: boolean;
		hidden?: boolean;
		sortBy?: (a: T, b: T) => number;
		filterBy?: (value: string, row: T) => boolean;
	}

	// generic sortBy function
	// If a column-level sortBy is defined for the provided key, it will be used.
	export function defaultSortBy<T extends Partial<Record<K, unknown>>, K extends keyof T>(
		a: T,
		b: T,
		key: K,
		sortDescending: boolean,
		columns?: TableColumn<T>[]
	): number {
		// Prefer a column-specific sortBy when available
		const col = columns?.find((c) => c.key === (key as unknown as keyof Partial<T>));
		if (col?.sortable !== false && typeof col?.sortBy === 'function') {
			// Apply descending by swapping args to keep column sortBy simple
			if (sortDescending) {
				return col.sortBy(b as T, a as T);
			}
			return col.sortBy(a as T, b as T);
		}

		// Fallback: default comparison by field value
		let [aField, bField] = [a[key], b[key]];
		if (sortDescending) [bField, aField] = [aField, bField];

		if (typeof aField === 'number' && typeof bField === 'number')
			return (aField as number) - (bField as number);
		if (typeof aField === 'bigint' && typeof bField === 'bigint')
			return Number((aField as bigint) - (bField as bigint));
		if (typeof aField === 'boolean' && typeof bField === 'boolean')
			return aField === bField ? 0 : aField ? -1 : 1;

		// String compare fallback
		return `${aField ?? ''}`.localeCompare(`${bField ?? ''}`);
	}
</script>

<script lang="ts">
	import type { Snippet } from 'svelte';

	type T = $$Generic<Partial<Record>>;
	type C = $$Generic<T>;
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

	type Props = {
		classes?: TableClasses;
		columns?: TableColumn<C>[];
		rows?: T[];
		filterBy?: string;
		externalSortAndFilter?: boolean;
		head?: Snippet<
			[{ isSorted: boolean; sortDescending: boolean; sortable: boolean; column: TableColumn<T> }]
		>;
		cell?: Snippet<[{ row: T; column: TableColumn<C>; cell: unknown }]>;
		empty?: Snippet;
	};

	let {
		classes = defaultClasses,
		columns = [],
		rows = [],
		filterBy = '',
		externalSortAndFilter = false,
		head,
		cell,
		empty
	}: Props = $props();

	let lastSortedKey: keyof C | '' = $state('');
	let sortDescending = $state(false);

	/**
	 * sort rows by a column key
	 * @param key the column key to sort by
	 */
	function sortRowsBy(key: keyof C): T[] {
		const columnData = columns.find((column) => column.key === key);
		if (!columnData || columnData.sortable === false) {
			return rows;
		}

		// call column sortBy
		if (columnData.sortBy) {
			const sortBy = columnData.sortBy;
			return [...rows].sort((a, b) => {
				[a, b] = sortDescending ? [a, b] : [b, a];
				return sortBy(a, b);
			});
		}

		// sort by content by default
		return [...rows].sort((a, b) => defaultSortBy(a, b, key, sortDescending));
	}

	function getSortingOrder(key: keyof C): boolean {
		if (lastSortedKey === key) return !sortDescending;
		return false;
	}

	function filterRowsBy(value: string, rows: T[]) {
		const numColumns = columns.length;
		return rows.filter((row) => {
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

	let filteredRows = $derived(
		(() => {
			if (externalSortAndFilter) {
				// rows come filtered and sorted, return them as is
				return rows;
			}
			if (lastSortedKey) {
				return sortRowsBy(lastSortedKey);
			}
			return filterRowsBy(filterBy, rows);
		})()
	);

	let assignedClasses = $derived({ ...defaultClasses, ...classes });
</script>

<table class={assignedClasses.table} style="border-spacing: 0">
	<thead class={assignedClasses.thead}>
		<tr class={assignedClasses.headtr}>
			{#each columns as column (column.key)}
				{#if !column.hidden}
					<th
						scope="col"
						class={assignedClasses.th}
						onclick={() => {
							sortDescending = getSortingOrder(lastSortedKey);
							if (!externalSortAndFilter) {
								lastSortedKey = column.key;
							}
						}}
					>
						{#if head}
							{@render head({
								column,
								isSorted: lastSortedKey === column.key,
								sortDescending,
								sortable: column.sortable !== false
							})}
						{:else}
							<span>{column.title}</span>
						{/if}
					</th>
				{/if}
			{/each}
		</tr>
	</thead>
	<tbody class={assignedClasses.tbody}>
		{#each filteredRows as row (row)}
			<tr class={`${assignedClasses.tr}`}>
				{#each columns as column (column.key)}
					{#if !column.hidden}
						<td class={assignedClasses.td}>
							{#if cell}
								{@render cell({ row, column, cell: row[column.key] })}
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
