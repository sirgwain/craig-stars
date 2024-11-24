import { type MapObject } from '$lib/types/MapObject';
import { User, emptyUser } from '$lib/types/User';
import type { Vector } from '$lib/types/Vector';
import type { Component, ComponentType, SvelteComponent } from 'svelte';
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
	{ component: Component<any>; props: any } | undefined
>();
export const tooltipLocation = writable<Vector>({ x: 0, y: 0 });

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const showTooltip = <T extends Partial<Record<string, any>>>(
	x: number,
	y: number,
	component: Component<T>,
	props?: T
) => {
	tooltipLocation.update(() => ({
		x,
		y
	}));
	tooltipComponent.update(() => ({
		component,
		props
	}));
};

