import { Player } from '$lib/types/Player';
import { getContext, setContext } from 'svelte';
import { writable, type Readable, type Writable } from 'svelte/store';
import { FullGame } from './FullGame';
import { Universe } from './Universe';

export const playerFinderKey = Symbol();
export const designFinderKey = Symbol();
export const gameKey = Symbol();

type GameContext = {
	game: Readable<FullGame>;
	player: Readable<Player>;
	universe: Readable<Universe>;
};

export const gameStore: GameContext = {
	game: writable<FullGame>(new FullGame()),
	player: writable<Player>(new Player()),
	universe: writable<Universe>(new Universe())
};

// init the game context with empty data
export const initGameContext = () => setContext<GameContext>(gameKey, gameStore);
export const getGameContext = () => getContext<GameContext>(gameKey);

// update the game context after a load
export const updateGameContext = (game: FullGame, player: Player, universe: Universe) => {
	const writableGameStore = gameStore.game as Writable<FullGame>;
	writableGameStore.update(() => game);

	const writablePlayerStore = gameStore.player as Writable<Player>;
	writablePlayerStore.update(() => player);

	const writableUniverseStore = gameStore.universe as Writable<Universe>;
	writableUniverseStore.update(() => universe);
};

export const updateGame = (game: FullGame) => {
	const writableGameStore = gameStore.game as Writable<FullGame>;
	writableGameStore.update(() => game);
};

export const updatePlayer = (player: Player) => {
	const writablePlayerStore = gameStore.player as Writable<Player>;
	writablePlayerStore.update(() => player);
};

export const updateUniverse = (universe: Universe) => {
	const writableUniverseStore = gameStore.universe as Writable<Universe>;
	writableUniverseStore.update(() => universe);
};
