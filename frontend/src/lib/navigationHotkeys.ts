import { goto } from '$app/navigation';
import type { Page } from '@sveltejs/kit';
import hotkeys from 'hotkeys-js';
import { get, type Readable } from 'svelte/store';

export const bindNavigationHotkeys = (gameId: number, page: Readable<Page>) => {
	hotkeys('esc', 'root', () => {
		goto(`/games/${gameId}`);
	});
	hotkeys('F2', 'root', (event) => {
		event.preventDefault();
		goto(`/games/${gameId}/techs`);
	});
	hotkeys('F3', 'root', (event) => {
		event.preventDefault();
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
	hotkeys('F4', 'root', (event) => {
		event.preventDefault();
		goto(`/games/${gameId}/designer`);
	});
	hotkeys('F5', 'root', (event) => {
		event.preventDefault();
		goto(`/games/${gameId}/research`);
	});
	hotkeys('F7', 'root', (event) => {
		event.preventDefault();
		goto(`/games/${gameId}/relations`);
	});
	hotkeys('F8', 'root', (event) => {
		event.preventDefault();
		goto(`/games/${gameId}/race`);
	});
	hotkeys('F10', 'root', (event) => {
		event.preventDefault();
		goto(`/games/${gameId}/players`);
	});
};

export const unbindNavigationHotkeys = () => {
	hotkeys.unbind('esc', 'root');
	hotkeys.unbind('F2', 'root');
	hotkeys.unbind('F3', 'root');
	hotkeys.unbind('F4', 'root');
	hotkeys.unbind('F5', 'root');
	hotkeys.unbind('F7', 'root');
	hotkeys.unbind('F8', 'root');
	hotkeys.unbind('F10', 'root');
};
