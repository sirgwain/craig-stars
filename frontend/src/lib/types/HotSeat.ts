import type { PlayerStatus } from '#lib/types/cs-proto.js';

// the players in a game controlled by a user
export function getUserPlayers(players: PlayerStatus[], userId: bigint): PlayerStatus[] {
	return players.filter((p) => !p.aiControlled && p.userId === userId);
}

// a hot seat game is one where a user controls more than one player
export function isHotSeat(players: PlayerStatus[], userId: bigint): boolean {
	return getUserPlayers(players, userId).length > 1;
}

// the first of a user's players that still needs to submit a turn, or their first player
export function getFirstHotSeatPlayer(
	players: PlayerStatus[],
	userId: bigint
): PlayerStatus | undefined {
	const userPlayers = getUserPlayers(players, userId);
	return userPlayers.find((p) => !p.submittedTurn) ?? userPlayers[0];
}
