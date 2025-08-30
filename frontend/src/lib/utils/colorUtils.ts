import type { CommandedPlayer } from '$lib/types/Player';
import type { PlayerSettings } from '$lib/types/PlayerSettings';
import type { Universe } from '$lib/services/Universe';

/**
 * Gets the display color for a map object based on player relationship and settings
 * @param playerNum - The player number of the object
 * @param currentPlayer - The current player (to check relationships)
 * @param universe - The universe service (for getting player colors)
 * @param settings - The player settings (for showPlayerColors flag)
 * @returns The color string to use for display
 */
export function getDisplayColor(
	playerNum: number | undefined,
	currentPlayer: CommandedPlayer,
	universe: Universe,
	settings: PlayerSettings
): string {
	if (!playerNum) {
		return '#999999'; // gray for no owner
	}

	// If showPlayerColors is enabled, use the actual player colors
	if (settings.showPlayerColors) {
		return universe.getPlayerColor(playerNum);
	}

	// Otherwise use relationship-based colors
	if (playerNum === currentPlayer.num) {
		return '#0000FF'; // blue for myself
	} else if (currentPlayer.isFriend(playerNum)) {
		return '#FFFF00'; // yellow for allies
	} else {
		return '#FF0000'; // red for enemies and neutrals
	}
}
// returns "black" if the color is light, otherwise "white"
export function getStrokeColor(color: string) {
	// strip #
	color = color.replace('#', '');

	// expand shorthand like #0F0 → #00FF00
	if (color.length === 3) {
		color = color
			.split('')
			.map((c) => c + c)
			.join('');
	}

	const r = parseInt(color.substr(0, 2), 16);
	const g = parseInt(color.substr(2, 2), 16);
	const b = parseInt(color.substr(4, 2), 16);

	// relative luminance formula
	const luminance = (0.299 * r + 0.587 * g + 0.114 * b) / 255;

	return luminance > 0.5 ? '#333' : '#bbb'; // dark gray for light fills, light gray for dark fills
}
