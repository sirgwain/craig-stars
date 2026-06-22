import type { Component } from 'svelte';
import { writable } from 'svelte/store';

export type PopupPropsBase = {
	onClose?: () => void;
};

export type PopupComponentStore<T extends PopupPropsBase = PopupPropsBase> = {
	component: Component<T>;
	props?: Omit<T, 'onClose'>;
};

export type PopupComponentBase = {
	component: Component<PopupPropsBase>;
	props?: Record<string, unknown>;
};

export const popupComponent = writable<PopupComponentBase | undefined>();
export const popupLocation = writable<{ x: number; y: number }>({ x: 0, y: 0 });

export function showPopup<T extends PopupPropsBase>(
	x: number,
	y: number,
	component: Component<T>,
	props?: Omit<T, 'onClose'>
) {
	popupLocation.set({ x, y });

	popupComponent.set({
		component: component as Component<PopupPropsBase>,
		props: props as Record<string, unknown>
	});
}
