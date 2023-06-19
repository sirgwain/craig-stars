import { goto } from '$app/navigation';
import hotkeys from 'hotkeys-js';


export const bindNavigationHotkeys = (gameId: number) => {
	hotkeys('esc', () => {
		goto(`/games/${gameId}`);
	});

	hotkeys('F5', () => {
		goto(`/games/${gameId}/research`);
	});

};

export const unbindNavigationHotkeys = () => {
	hotkeys.unbind('esc');
	hotkeys.unbind('F5');
};
