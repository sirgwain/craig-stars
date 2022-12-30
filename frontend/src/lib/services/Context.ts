import type { Fleet, Waypoint } from '$lib/types/Fleet';
import type { Game } from '$lib/types/Game';
import { MapObjectType, ownedBy, positionKey, type MapObject } from '$lib/types/MapObject';
import type { Planet } from '$lib/types/Planet';
import { findMyPlanet, type Player } from '$lib/types/Player';
import { PlayerSettings } from '$lib/types/PlayerSettings';
import { emptyUser, type User } from '$lib/types/User';
import { derived, get, writable } from 'svelte/store';
import { rollover } from './Math';

export type MapObjectsByPosition = {
	[k: string]: MapObject[];
};

export const me = writable<User>(emptyUser);
export const game = writable<Game | undefined>();
export const player = writable<Player | undefined>();
export const settings = writable<PlayerSettings>(loadSettingsOrDefault());

export const commandedPlanet = writable<Planet | undefined>();
export const commandedFleet = writable<Fleet | undefined>();
export const selectedWaypoint = writable<Waypoint | undefined>();
export const selectedMapObject = writable<MapObject>();
export const commandedMapObject = writable<MapObject>();
export const highlightedMapObject = writable<MapObject | undefined>();
export const commandedMapObjectName = writable<string>();
export const zoomTarget = writable<MapObject | undefined>();

export const currentMapObjectIndex = derived(
	[player, commandedFleet, commandedPlanet],
	([$player, $commandedFleet, $commandedPlanet]) => {
		if ($player && $commandedPlanet) {
			return $player.planets.findIndex((p) => p === $commandedPlanet);
		}
		if ($player && $commandedFleet) {
			return $player.fleets.findIndex((f) => f === $commandedFleet);
		}
		return 0;
	}
);

// command the previous mapObject for this type, i.e. the previous planet or fleet
export const previousMapObject = () => {
	const p = get(player);
	if (!p) {
		return;
	}
	const i = get(currentMapObjectIndex);
	const mo = get(commandedMapObject);

	if (mo.type == MapObjectType.Planet) {
		const prevIndex = rollover(i - 1, 0, p.planets.length - 1);
		commandMapObject(p.planets[prevIndex]);
		zoomToMapObject(p.planets[prevIndex]);
	} else if (mo.type == MapObjectType.Fleet) {
		const prevIndex = rollover(i - 1, 0, p.fleets.length - 1);
		commandMapObject(p.fleets[prevIndex]);
		zoomToMapObject(p.fleets[prevIndex]);
	}
};

// command the next mapObject for this type, i.e. the next planet or fleet
export const nextMapObject = () => {
	const p = get(player);
	if (!p) {
		return;
	}
	const i = get(currentMapObjectIndex);
	const mo = get(commandedMapObject);
	if (mo.type == MapObjectType.Planet) {
		const nextIndex = rollover(i + 1, 0, p.planets.length - 1);
		commandMapObject(p.planets[nextIndex]);
		zoomToMapObject(p.planets[nextIndex]);
	} else if (mo.type == MapObjectType.Fleet) {
		const nextIndex = rollover(i + 1, 0, p.fleets.length - 1);
		commandMapObject(p.fleets[nextIndex]);
		zoomToMapObject(p.fleets[nextIndex]);
	}
};

export const mapObjectsByPosition = derived(player, ($player) => {
	if (!$player) return undefined;

	const dict: MapObjectsByPosition = {};
	const addtoDict = (mo: MapObject) => {
		const key = positionKey(mo);
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
	if (!$player) return undefined;

	const dict: MapObjectsByPosition = {};
	const addtoDict = (mo: MapObject) => {
		if (!ownedBy(mo, $player.num)) {
			return;
		}
		const key = positionKey(mo);
		if (!dict[key]) {
			dict[key] = [];
		}
		dict[key].push(mo);
	};

	$player.planetIntels?.forEach(addtoDict);
	$player.fleetIntels?.forEach(addtoDict);
	$player.fleets?.forEach(addtoDict);
	$player.planets?.forEach(addtoDict);

	return dict;
});

// get a mapobject by type, number, and optionally player num
export const getMapObject = (
	player: Player,
	type: MapObjectType,
	num: number,
	playerNum?: number
): MapObject | undefined => {
	let mo: MapObject;
	switch (type) {
		case MapObjectType.Planet:
			mo = player.planetIntels[num - 1];
			if (mo.playerNum == player.num) {
				return findMyPlanet(player, mo as Planet);
			}
			return mo;
		case MapObjectType.Fleet:
			if (playerNum == player.num) {
				return player.fleets.find((f) => f.num == num);
			}
			return player.fleetIntels?.find((f) => f.num == num && f.playerNum == playerNum);
		case MapObjectType.Wormhole:
			break;
		case MapObjectType.MineField:
			break;
		case MapObjectType.MysteryTrader:
			break;
		case MapObjectType.Salvage:
			break;
		case MapObjectType.MineralPacket:
			break;
		case MapObjectType.PositionWaypoint:
			break;
	}
};

export const selectMapObject = (mo: MapObject) => {
	selectedMapObject.update(() => mo);
};

export const selectWaypoint = (wp: Waypoint) => {
	selectedWaypoint.update(() => wp);
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
		selectedWaypoint.update(() => (mo as Fleet).waypoints[0]);
	}

	commandedMapObjectName.update(() => mo.name);
};

export const highlightMapObject = (mo: MapObject | undefined) => {
	highlightedMapObject.update(() => mo);
};

export const zoomToMapObject = (mo: MapObject) => {
	zoomTarget.update(() => mo);
};

settings.subscribe((value) => {
	localStorage.setItem('playerSettings', JSON.stringify(value));
});

function loadSettingsOrDefault(): PlayerSettings {
	const json = localStorage.getItem('playerSettings');
	if (json) {
		const settings = JSON.parse(json) as PlayerSettings;
		if (settings) {
			return settings;
		}
	}

	return new PlayerSettings();
}

export const playerName = (num: number) => {
	const p = get(player);
	if (p && num > 0 && num <= p.playerIntels.length) {
		const intel = p.playerIntels[num - 1];
		return intel.racePluralName ?? intel.name;
	}
};
