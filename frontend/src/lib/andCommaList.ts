/**
 * Concatenate several non-empty strings with commas, adding a conjunction between the final 2 entries.
 * Empty strings are omitted inside the end result.
 * @param items an array of strings to concatenate
 * @param emptyResult a result to return if the array is undefined, empty or contains no non-blank strings.
 * @param finalConjunction a conjunction to use on the final word
 * @returns A string containing all the provided elements delimited by commas and the chosen conjunction.
 */
export function andCommaList(
	items: string[] | undefined,
	emptyResult = '',
	finalConjunction = 'and'
): string {
	if (!items) {
		return emptyResult;
	}
	const filteredItems = items.filter((i) => !!i); // remove empty strings

	if (filteredItems.length == 0) {
		return emptyResult;
	}

	let result = ''
	for (const [i, val] of filteredItems.entries()) {
		result += val;
		if (i == filteredItems.length - 2) {
			// add final conjunction for 2nd to last
			result += ` ${finalConjunction} `;
		} else if (i != filteredItems.length - 1) {
			// add commas after all but the last item
			result += ', ';
		}

	}

	return result;
}
