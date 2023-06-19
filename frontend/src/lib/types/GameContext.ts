import type { Game } from './Game';
import type { Player } from './Player';

export type GameContext = {
	game: Game;
	player: Player;
};
