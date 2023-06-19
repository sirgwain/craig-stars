import type { Game } from '$lib/types/Game';
import type { GameContext } from '$lib/types/GameContext';
import type { Planet } from '$lib/types/Planet';
import type { Player } from '$lib/types/Player';
import { getContext, setContext } from 'svelte';

export function getGameContext(): GameContext {
	return getContext<GameContext>('game');
}

export function setGameContext(game: Game, player: Player) {
	setContext<GameContext>('game', { game, player });
}

export function getCommandedPlanet(): Planet {
	return getContext<Planet>('commandedPlanet');
}

export function setCommandedPlanet(planet: Planet) {
	setContext<Planet>('commandedPlanet', planet);
}
