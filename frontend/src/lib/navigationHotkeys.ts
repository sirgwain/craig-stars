import { goto } from '$app/navigation';
import type { Page } from '@sveltejs/kit';
import hotkeys from 'hotkeys-js';
import { get, type Readable } from 'svelte/store';

export const bindNavigationHotkeys = (gameId: number, page: Readable<Page>) => {
	hotkeys('esc', () => {
		goto(`/games/${gameId}`);
	});
	hotkeys('F2', () => {
		goto(`/games/${gameId}/techs`);
	});
	hotkeys('F5', () => {
		goto(`/games/${gameId}/research`);
	});
	hotkeys('F4', () => {
		goto(`/games/${gameId}/designer`);
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
			case `/games/${gameId}/fleets`:
				goto(`/games/${gameId}/designs`);
				break;
			case `/games/${gameId}/designs`:
				goto(`/games/${gameId}/messages`);
				break;
			case `/games/${gameId}/messages`:
				goto(`/games/${gameId}/battles`);
				break;
			default:
				goto(`/games/${gameId}`);
				break;
		}
	});
	hotkeys('F10', () => {
		goto(`/games/${gameId}/players`);
	});
};

export const unbindNavigationHotkeys = () => {
	hotkeys.unbind('esc');
	hotkeys.unbind('F5');
	hotkeys.unbind('F4');
	hotkeys.unbind('F3');
	hotkeys.unbind('F2');
	hotkeys.unbind('F10');
};
