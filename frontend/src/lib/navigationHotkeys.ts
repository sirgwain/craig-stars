import { goto } from '$app/navigation';
import hotkeys from 'hotkeys-js';

export const bindNavigationHotkeys = (gameId: number) => {
	hotkeys('esc', () => {
		goto(`/games/${gameId}`);
	});

	hotkeys('F5', () => {
		goto(`/games/${gameId}/research`);
	});
	hotkeys('F4', () => {
		goto(`/games/${gameId}/designs`);
	});
};

export const unbindNavigationHotkeys = () => {
	hotkeys.unbind('esc');
	hotkeys.unbind('F5');
	hotkeys.unbind('F4');
};
