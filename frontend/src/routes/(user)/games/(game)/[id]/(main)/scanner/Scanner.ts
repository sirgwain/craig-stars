import { type Fleet } from '$lib/types/Fleet';
import type { MapObject } from '$lib/types/MapObject';
import type { Player } from '$lib/types/Player';
import { find } from 'lodash-es';

export type ViewportCoords = {
	x: number;
	y: number;
	visible: boolean;
};

// return the viewport coords and whether they are in the viewport
export function getViewportCoords(
	x: number,
	y: number,
	width: number,
	height: number,
	padding = 0
): ViewportCoords {
	return {
		x,
		y,
		visible: x > -padding && y > -padding && x < width + padding && y < height + padding
	};
}

// for a list of orbiting fleets, return whether there are enemies, friends, both or neither
export function getEnemiesAndFriends(
	orbitingFleets: MapObject[],
	player: Player
): { enemies: boolean; friends: boolean } {
	const playerNums = new Set<number>(orbitingFleets.map((f) => f.playerNum));
	if (playerNums.size == 1) {
		if (orbitingFleets[0].playerNum === player.num) {
			return { enemies: false, friends: false };
		} else if (player.isEnemy(orbitingFleets[0].playerNum)) {
			return { enemies: true, friends: false };
		} else {
			return { enemies: false, friends: true };
		}
	} else {
		const enemies = find(playerNums, (n: number) => player.isEnemy(n));
		const friends = find(playerNums, (n: number) => player.isFriendOrNeutral(n));

		if (friends && !enemies) {
			return { enemies: false, friends: true };
		} else if (!friends && enemies) {
			return { enemies: true, friends: false };
		} else {
			return { enemies: true, friends: true };
		}
	}
}
