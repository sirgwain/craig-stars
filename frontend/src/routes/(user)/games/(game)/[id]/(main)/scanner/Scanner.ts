import type { MapObjectLike } from '#lib/types/MapObject.js';
import type { CommandedPlayer } from '#lib/types/Player.js';
import { find } from 'lodash-es';

// for a list of orbiting fleets, return whether there are enemies, friends, both or neither
export function getEnemiesAndFriends(
	orbitingFleets: MapObjectLike[],
	player: CommandedPlayer
): { enemies: boolean; friends: boolean } {
	const playerNums = new Set<number>(orbitingFleets.map((f) => f.mapObject?.playerNum ?? 0));
	if (playerNums.size == 1) {
		if (orbitingFleets[0].mapObject?.playerNum === player.num) {
			return { enemies: false, friends: false };
		} else if (player.isEnemy(orbitingFleets[0].mapObject?.playerNum ?? 0)) {
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
