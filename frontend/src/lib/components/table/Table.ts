export interface TableColumn<T> {
	key: keyof Partial<T>;
	title: string;
	sortable?: boolean;
	filterable?: boolean;
	hidden?: boolean;
	sortBy?: (a: T, b: T) => number;
	filterBy?: (value: string, row: T) => boolean;
}

// If a column-level sortBy is defined for the provided key, it will be used.
export function defaultSortBy<T extends Partial<Record<K, unknown>>, K extends keyof T>(
	a: T,
	b: T,
	key: K,
	sortDescending: boolean,
	columns?: TableColumn<T>[]
): number {
	const col = columns?.find((c) => c.key === (key as unknown as keyof Partial<T>));
	if (col?.sortable !== false && typeof col?.sortBy === 'function') {
		return sortDescending ? col.sortBy(b as T, a as T) : col.sortBy(a as T, b as T);
	}

	let [aField, bField] = [a[key], b[key]];
	if (sortDescending) [bField, aField] = [aField, bField];

	if (typeof aField === 'number' && typeof bField === 'number') return aField - bField;
	if (typeof aField === 'bigint' && typeof bField === 'bigint') return Number(aField - bField);
	if (typeof aField === 'boolean' && typeof bField === 'boolean')
		return aField === bField ? 0 : aField ? -1 : 1;

	return `${aField ?? ''}`.localeCompare(`${bField ?? ''}`);
}
