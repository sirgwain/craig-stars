import { type MapObject } from '$lib/types/MapObject';
import { User, emptyUser } from '$lib/types/User';
import type { Vector } from '$lib/types/Vector';
import type { ComponentType, SvelteComponent } from 'svelte';
import { writable } from 'svelte/store';
import { TechService } from './TechService';

export type MapObjectsByPosition = {
	[k: string]: MapObject[];
};

export const me = writable<User>(emptyUser);
export const techs = writable<TechService>(new TechService());
export const loadingModalText = writable<string | undefined>(undefined);

export const setLoadingModalText = (text: string) => {
	loadingModalText.update(() => text);
};

export const clearLoadingModalText = () => {
	loadingModalText.update(() => undefined);
};

export const tooltipComponent = writable<
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	{ component: typeof SvelteComponent; props: any } | undefined
>();
export const tooltipLocation = writable<Vector>({ x: 0, y: 0 });

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const showTooltip = <T>(x: number, y: number, component: ComponentType, props?: T) => {
	tooltipLocation.update(() => ({
		x,
		y
	}));
	tooltipComponent.update(() => ({
		component,
		props
	}));
};

export const popupComponent = writable<
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	{ component: typeof SvelteComponent; props: any } | undefined
>();
export const popupLocation = writable<Vector>({ x: 0, y: 0 });

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const showPopup = <T>(x: number, y: number, component: ComponentType, props?: T) => {
	popupLocation.update(() => ({
		x,
		y
	}));
	popupComponent.update(() => ({
		component,
		props
	}));
};
