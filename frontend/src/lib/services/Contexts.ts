import { Player } from '$lib/types/Player';
import { PlayerSettings } from '$lib/types/PlayerSettings';
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
	settings: Writable<PlayerSettings>;
};

export const gameStore: GameContext = {
	game: writable<FullGame>(new FullGame()),
	player: writable<Player>(new Player()),
	universe: writable<Universe>(new Universe()),
	settings: writable<PlayerSettings>(new PlayerSettings())
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

	// update the settings, or use existing settings from localStorage
	const settingsStore = gameStore.settings;
	const settings = loadSettingsOrDefault(game.id, player.num);
	settingsStore.subscribe((value) => {
		value.beforeSave();
		localStorage.setItem(value.key, JSON.stringify(value));
	});
	settingsStore.update(() => settings);
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

function loadSettingsOrDefault(gameId: number, playerNum: number): PlayerSettings {
	const key = PlayerSettings.key(gameId, playerNum);

	const json = localStorage.getItem(key);
	if (json) {
		const settingsJSON = JSON.parse(json) as PlayerSettings;
		if (settingsJSON) {
			// create a new object
			const settings = new PlayerSettings(gameId, playerNum);
			Object.assign(settings, settingsJSON);
			settings.afterLoad();
			return settings;
		}
	}

	return new PlayerSettings(gameId, playerNum);
}
