import type { Fleet } from '$lib/types/Fleet';
import type { Game } from '$lib/types/Game';
import { MapObjectType, ownedBy, type MapObject } from '$lib/types/MapObject';
import type { Planet } from '$lib/types/Planet';
import type { Player } from '$lib/types/Player';
import { derived, writable } from 'svelte/store';

export type MapObjectsByPosition = {
	[k: string]: MapObject[];
};

export const game = writable<Game>();
export const player = writable<Player>();

export const commandedPlanet = writable<Planet | undefined>();
export const commandedFleet = writable<Fleet | undefined>();
export const selectedMapObject = writable<MapObject>();
export const commandedMapObject = writable<MapObject>();
export const highlightedMapObject = writable<MapObject | undefined>();

export const mapObjectsByPosition = derived(player, ($player) => {
	const dict: MapObjectsByPosition = {};
	const addtoDict = (mo: MapObject) => {
		const key = `${mo.position.x},${mo.position.y}`;
		if (!dict[key]) {
			dict[key] = [];
		}
		dict[key].push(mo);
	};

	$player.planetIntels?.forEach(addtoDict);
	$player.fleetIntels?.forEach(addtoDict);
	$player.fleets?.forEach(addtoDict);

	return dict;
});

export const myMapObjectsByPosition = derived(player, ($player) => {
	const dict: MapObjectsByPosition = {};
	const addtoDict = (mo: MapObject) => {
		if (!ownedBy(mo, $player.num)) {
			return;
		}
		const key = `${mo.position.x},${mo.position.y}`;
		if (!dict[key]) {
			dict[key] = [];
		}
		dict[key].push(mo);
	};

	$player.planetIntels?.forEach(addtoDict);
	$player.fleetIntels?.forEach(addtoDict);
	$player.fleets?.forEach(addtoDict);

	return dict;
});

export const selectMapObject = (mo: MapObject) => {
	selectedMapObject.update(() => mo);
};

export const commandMapObject = (mo: MapObject) => {
	// console.log(`Commanded ${mo.type}:${mo.name}`);
	commandedMapObject.update(() => mo);
	if (mo.type == MapObjectType.Planet) {
		commandedPlanet.update(() => mo as Planet);
		commandedFleet.update(() => undefined);
	} else if (mo.type == MapObjectType.Fleet) {
		commandedPlanet.update(() => undefined);
		commandedFleet.update(() => mo as Fleet);
	}
};

export const highlightMapObject = (mo: MapObject | undefined) => {
	highlightedMapObject.update(() => mo);
};
