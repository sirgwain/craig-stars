import type { Game } from '$lib/types/Game';
import type { GameContext } from '$lib/types/GameContext';
import type { Planet } from '$lib/types/Planet';
import type { Player } from '$lib/types/Player';
import { getContext, setContext } from 'svelte';
import { writable } from 'svelte/store';

export const commandedPlanet = writable<Planet>();

export function getGameContext(): GameContext {
	return getContext<GameContext>('game');
}

export function setGameContext(game: Game, player: Player) {
	setContext<GameContext>('game', { game, player });
}
