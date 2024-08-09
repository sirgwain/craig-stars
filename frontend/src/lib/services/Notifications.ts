import { get, writable } from 'svelte/store';

// simple store to display notifications on save
export const notification = writable<string>();
export function notify(value: string) {
	const n = get(notification);
	if (n != value) {
		notification.update(() => value);
	}
}
