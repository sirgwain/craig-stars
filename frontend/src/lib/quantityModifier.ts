import { clamp } from './services/Math';

// get the quantity modifier for multiplying adding/removing items
export const quantityModifier = (e: MouseEvent, min = 0, max = Number.MAX_SAFE_INTEGER) =>
	clamp((e.shiftKey ? 10 : 1) * (e.metaKey || e.ctrlKey ? 100 : 1), min, max);
