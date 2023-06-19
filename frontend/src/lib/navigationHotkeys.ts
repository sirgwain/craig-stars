import { goto } from '$app/navigation';
import type { Page } from '@sveltejs/kit';
import hotkeys from 'hotkeys-js';
import { get, type Readable } from 'svelte/store';

export const bindNavigationHotkeys = (gameId: number, page: Readable<Page>) => {
	hotkeys('esc', () => {
		goto(`/games/${gameId}`);
	});

	hotkeys('F5', () => {
		goto(`/games/${gameId}/research`);
	});
	hotkeys('F4', () => {
		goto(`/games/${gameId}/designs`);
	});
	hotkeys('F3', () => {
		const pathname = get(page)?.url.pathname;
		switch (pathname) {
			case `/games/${gameId}`:
				goto(`/games/${gameId}/planets`);
				break;
			case `/games/${gameId}/planets`:
				goto(`/games/${gameId}/fleets`);
				break;
			default:
				goto(`/games/${gameId}`);
				break;
		}
	});
};

export const unbindNavigationHotkeys = () => {
	hotkeys.unbind('esc');
	hotkeys.unbind('F5');
	hotkeys.unbind('F4');
	hotkeys.unbind('F3');
};
