import { getContext } from 'svelte';
import { get, writable, type Writable } from 'svelte/store';

export const carouselKey = Symbol();

export type CarouselContext = {
	open: Writable<boolean>;

	onDisclosureClicked: () => void;
};

// init the game context with empty data
export const getCarouselContext = () => getContext<CarouselContext | undefined>(carouselKey);

// update the game context after a load
export function createCarouselContext(): CarouselContext {
	const open = writable(true);

	function onDisclosureClicked() {
		open.set(!get(open));
	}

	return {
		open,
		onDisclosureClicked
	};
}
