/**
 * Concaentate several strings with commas, adding a conjunction between the final 2 entries.
 * @param items a list of strings to concatenate
 * @param emptyResult a result to return if the string is empty, default nothing
 * @param finalConjunction A conjunction to use on the final word; defaults to the word 'and'.
 * @returns A string containing all of items' elements delimited by commas and the chosen conjunction.
 */
export function andCommaList(items: string[], emptyResult = '', finalConjunction = 'and'): string {
	const filteredItems = items.filter((i) => !!i); // remove empty strings

	if (filteredItems.length == 0) {
		return emptyResult;
	}

	if (filteredItems.length == 1) {
		return filteredItems[0];
	}

	if (filteredItems.length == 2) {
		return `${filteredItems[0]} ${finalConjunction} ${filteredItems[1]}`;
	}

	let result = filteredItems[0];

	for (let i = 1; i < filteredItems.length; i++) {
		if (i == filteredItems.length - 1) {
			result += ` ${finalConjunction} `;
		} else {
			result += ', ';
		}

		result += filteredItems[i];
	}

	return result;
}
