
// get the quantity modifier for multiplying adding/removing items
export const quantityModifier = (e: MouseEvent) =>
	(e.shiftKey ? 10 : 1) * (e.metaKey || e.ctrlKey ? 100 : 1);
