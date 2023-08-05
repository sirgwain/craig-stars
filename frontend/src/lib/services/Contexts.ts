import { getNextVisibleMessageNum } from '$lib/types/Message';
import { Player } from '$lib/types/Player';
import { PlayerSettings } from '$lib/types/PlayerSettings';
import type { ShipDesign } from '$lib/types/ShipDesign';
import { getContext, setContext } from 'svelte';
import { writable, type Readable, type Writable } from 'svelte/store';
import { FullGame } from './FullGame';
import { universe } from './Stores';
import { Universe } from './Universe';

export const playerFinderKey = Symbol();
export const designFinderKey = Symbol();
export const gameKey = Symbol();

type GameContext = {
	game: Readable<FullGame>;
	player: Readable<Player>;
	universe: Readable<Universe>;
	settings: Writable<PlayerSettings>;
	designs: Writable<ShipDesign[]>;
	messageNum: Writable<number>;
};

export const gameStore: GameContext = {
	game: writable<FullGame>(new FullGame()),
	player: writable<Player>(new Player()),
	universe: universe,
	settings: writable<PlayerSettings>(new PlayerSettings()),
	designs: writable<ShipDesign[]>([]),
	messageNum: writable<number>(0)
};

// init the game context with empty data
export const initGameContext = () => setContext<GameContext>(gameKey, gameStore);
export const getGameContext = () => getContext<GameContext>(gameKey);
export const clearGameContext = () =>
	updateGameContext(new FullGame(), new Player(), new Universe());

// update the game context after a load
export const updateGameContext = (game: FullGame, player: Player, universe: Universe) => {
	// update the settings, or use existing settings from localStorage
	const settingsStore = gameStore.settings;
	const settings = loadSettingsOrDefault(game.id, player.num);
	settingsStore.subscribe((value) => {
		value.beforeSave();
		localStorage.setItem(value.key, JSON.stringify(value));
	});
	settingsStore.update(() => settings);
	universe.settings = settings;

	const writableGameStore = gameStore.game as Writable<FullGame>;
	writableGameStore.update(() => game);

	const writablePlayerStore = gameStore.player as Writable<Player>;
	writablePlayerStore.update(() => player);

	const writableUniverseStore = gameStore.universe as Writable<Universe>;
	writableUniverseStore.update(() => universe);

	// update the designs store when the universe updates
	gameStore.designs.update(() => universe.getMyDesigns());

	// init the messageNum when we load a game or have a new turn
	gameStore.messageNum.update(() => getNextVisibleMessageNum(-1, false, player.messages, settings));
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

	// update the designs store when the universe updates
	gameStore.designs.update(() => universe.getMyDesigns());
};

export const updateDesigns = (designs: ShipDesign[]) => {
	// update the designs store when the universe updates
	gameStore.designs.update(() => designs);
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
