import type { Fleet, Waypoint } from '$lib/types/Fleet';
import type { Game } from '$lib/types/Game';
import { MapObjectType, None, ownedBy, positionKey, type MapObject } from '$lib/types/MapObject';
import { CommandedPlanet, type Planet } from '$lib/types/Planet';
import type { Player } from '$lib/types/Player';
import { emptyUser, type User } from '$lib/types/User';
import type { Vector } from '$lib/types/Vector';
import { derived, get, writable } from 'svelte/store';
import { rollover } from './Math';

export type MapObjectsByPosition = {
	[k: string]: MapObject[];
};

export const me = writable<User>(emptyUser);
export const game = writable<Game | undefined>();
export const player = writable<Player | undefined>();

export const commandedPlanet = writable<CommandedPlanet | undefined>();
export const commandedFleet = writable<Fleet | undefined>();
export const selectedWaypoint = writable<Waypoint | undefined>();
export const selectedMapObject = writable<MapObject>();
export const commandedMapObject = writable<MapObject>();
export const highlightedMapObject = writable<MapObject | undefined>();
export const commandedMapObjectName = writable<string>();
export const zoomTarget = writable<MapObject | undefined>();

const mapObjectsByPosition = derived(player, ($player) => {
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

const myMapObjectsByPosition = derived(player, ($player) => {
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

const currentMapObjectIndex = derived(
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
		const planet = p.planets[prevIndex];
		commandMapObject(planet);
		zoomToMapObject(planet);
		selectMapObject(planet);
	} else if (mo.type == MapObjectType.Fleet) {
		const prevIndex = rollover(i - 1, 0, p.fleets.length - 1);
		commandMapObject(p.fleets[prevIndex]);
		zoomToMapObject(p.fleets[prevIndex]);

		const fleet = p.fleets[prevIndex];
		if (fleet.orbitingPlanetNum && fleet.orbitingPlanetNum != None) {
			const planet = getMapObject(MapObjectType.Planet, fleet.orbitingPlanetNum);
			if (planet) {
				selectMapObject(planet);
			}
		} else {
			selectMapObject(fleet);
		}
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
		const planet = p.planets[nextIndex];
		commandMapObject(planet);
		zoomToMapObject(planet);
		selectMapObject(planet);
	} else if (mo.type == MapObjectType.Fleet) {
		const nextIndex = rollover(i + 1, 0, p.fleets.length - 1);
		const fleet = p.fleets[nextIndex];
		commandMapObject(p.fleets[nextIndex]);
		zoomToMapObject(p.fleets[nextIndex]);
		if (fleet.orbitingPlanetNum && fleet.orbitingPlanetNum != None) {
			const planet = getMapObject(MapObjectType.Planet, fleet.orbitingPlanetNum);
			if (planet) {
				selectMapObject(planet);
			}
		} else {
			selectMapObject(fleet);
		}
	}
};

// get a mapobject by type, number, and optionally player num
export const getMapObject = (
	type: MapObjectType,
	num: number,
	playerNum?: number
): MapObject | undefined => {
	const p = get(player);
	if (p) {
		let mo: MapObject;
		switch (type) {
			case MapObjectType.Planet:
				mo = p.planetIntels[num - 1];
				if (mo.playerNum == p.num) {
					return findMyPlanetByNum(mo.num);
				}
				return mo;
			case MapObjectType.Fleet:
				if (playerNum == p.num) {
					return p?.fleets.find((f) => f.num == num);
				}
				return p?.fleetIntels?.find((f) => f.num == num && f.playerNum == playerNum);
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
	}
};

export const getMapObjectsByPosition = (position: MapObject | Vector): MapObject[] => {
	const mos = get(mapObjectsByPosition);
	if (mos) {
		return mos[positionKey(position)] ?? [];
	}
	return [];
};

export const getMyMapObjectsByPosition = (position: MapObject | Vector): MapObject[] => {
	const mos = get(myMapObjectsByPosition);
	if (mos) {
		return mos[positionKey(position)] ?? [];
	}
	return [];
};

export const findMyPlanetByNum = (num: number): Planet | undefined => {
	const p = get(player);
	return p?.planets?.find((p) => p.num == num);
};

export const findIntelMapObject = (mo: MapObject): MapObject | undefined => {
	const p = get(player);
	if (mo.type === MapObjectType.Planet) {
		return p?.planetIntels?.find((planet) => planet.num == mo.num) ?? mo;
	} else if (mo.type === MapObjectType.Fleet) {
		return p?.fleetIntels?.find((fleet) => fleet.num == mo.num) ?? mo;
	}
	return mo;
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
		commandedPlanet.update(() => Object.assign(new CommandedPlanet(), mo));
		commandedFleet.update(() => undefined);
	} else if (mo.type == MapObjectType.Fleet) {
		commandedPlanet.update(() => undefined);
		commandedFleet.update(() => mo as Fleet);
		selectedWaypoint.update(() => {
			const fleet = mo as Fleet;
			if (fleet?.waypoints && fleet.waypoints.length) {
				return fleet.waypoints[0];
			}
			return undefined;
		});
	}

	commandedMapObjectName.update(() => mo.name);
};

export const highlightMapObject = (mo: MapObject | undefined) => {
	highlightedMapObject.update(() => mo);
};

export const zoomToMapObject = (mo: MapObject) => {
	zoomTarget.update(() => mo);
};

export const playerName = (playerNum: number | undefined) => {
	const p = get(player);
	if (p && playerNum && playerNum > 0 && playerNum <= p.playerIntels.length) {
		const intel = p.playerIntels[playerNum - 1];
		return intel.racePluralName ?? intel.name;
	}
};

export const playerColor = (playerNum: number | undefined): string => {
	const p = get(player);
	if (p && playerNum && playerNum > 0 && playerNum <= p.playerIntels.length) {
		const intel = p.playerIntels[playerNum - 1];
		return intel.color ?? '#FF0000';
	}
	return '#FF0000';
};
