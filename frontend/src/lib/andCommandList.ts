/**
 * Concaentate several strings with commas, adding an "and" on the final word.
 * @param items a list of strings to concatenate
 * @param emptyResult a result to return if the string is empty
 * @returns The concaetation of all items
*/
export function andCommaList(items: string[], emptyResult = ''): string {
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

	return filteredItems.reduce((prevVal, currVal: string, i: number) => {
		if (i == filteredItems.length - 1) {
			prevVal += ' and ';
		} else if (i > 0) {
			prevVal += ', ';
		}

		return prevVal + currVal;
	});

}
