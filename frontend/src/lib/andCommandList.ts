import type { FocusEventHandler } from 'svelte/elements';
import { clamp } from './services/Math';

// for a list of items, return a comma separated list with an and on the final word
// ex:
// [a, b, c] will return `a, b and c`.
// [a, b] will return `a and b`.
// [a] will return `a`.
// empty strings in the list are discarded
export function andCommaList(items: string[], emptyResult = ''): string {
	let result = '';
	const filteredItems = items.filter((i) => !!i);

	if (filteredItems.length == 0) {
		return emptyResult;
	}

	if (filteredItems.length == 1) {
		return filteredItems[0];
	}

	if (filteredItems.length == 2) {
		return `${filteredItems[0]} and ${filteredItems[1]}`;
	}

	for (let i = 0; i < filteredItems.length; i++) {
		if (i == filteredItems.length - 1) {
			result += ' and ';
		} else if (i > 0) {
			result += ', ';
		}

		result += filteredItems[i];
	}

	return result;
}
