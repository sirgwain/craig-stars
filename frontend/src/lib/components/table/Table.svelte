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
	export function defaultSortBy<T extends Partial<Record<K, any>>, K extends keyof T>(
		a: T,
		b: T,
		key: K,
		sortDescending: boolean
	): number {
		let [aField, bField] = [a[key], b[key]];
		if (sortDescending) [bField, aField] = [aField, bField];
		if (typeof aField === 'number') return aField - bField;
		if (typeof aField === 'boolean') return aField ? -1 : 1;
		return aField?.localeCompare(bField);
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
		cell?: Snippet<[{ row: T; column: TableColumn<C>; cell: any }]>;
		empty?: Snippet;
	};

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

	let lastSortedKey: keyof C | '' = $state('');
	let sortDescending = $state(false);

	/**
	 * sort rows by a column key
	 * @param key the column key to sort by
	 * @param override true to force sort by descending
	 */
	function sortRowsBy(key: keyof C, override = false): void {
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
		rows = [...rows].sort((a, b) => defaultSortBy(a, b, key, sortDescending));
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

	let filteredRows = $derived(
		(() => {
			if (externalSortAndFilter) {
				// rows come filtered and sorted, return them as is
				return rows;
			}
			if (lastSortedKey) {
				sortRowsBy(lastSortedKey, true);
			}
			return filterRowsBy(filterBy, rows);
		})()
	);

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
							{@render head?.({
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
		{#each filteredRows as row}
			<tr class={`${assignedClasses.tr}`}>
				{#each columns as column}
					{#if !column.hidden}
						<td class={assignedClasses.td}>
							{#if cell}
								{@render cell?.({ row, column, cell: row[column.key] })}
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
