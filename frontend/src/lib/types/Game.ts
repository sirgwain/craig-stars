import type { PlayerStatus } from '$lib/types/cs-proto';
import { GameSchema, type Game, type GameWithPlayers } from '$lib/types/cs-proto';
import { create } from '@bufbuild/protobuf';

export type GameWithPlayersFlat = Game & { players: PlayerStatus[] };

export function getGameWithPlayersFlat(gwp: GameWithPlayers): GameWithPlayersFlat {
	return { ...(gwp.game ?? create(GameSchema)), players: gwp.players ?? [] };
}
