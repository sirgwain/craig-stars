import type { FocusEventHandler } from 'svelte/elements';
import { clamp } from './services/Math';

// number input blur event validator. This clamps the value between the min/max of a number input
export const validateNumberrInput: FocusEventHandler<HTMLInputElement> = (e) => {
	e.currentTarget.value = String(
		clamp(
			parseInt(e.currentTarget.value),
			parseInt(e.currentTarget.min),
			parseInt(e.currentTarget.max)
		)
	);
};
