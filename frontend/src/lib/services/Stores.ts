import { goto } from '$app/navigation';
import { CommandedFleet, type Fleet, type Waypoint } from '$lib/types/Fleet';
import { equal, MapObjectType, None, ownedBy, type MapObject } from '$lib/types/MapObject';
import { getMapObjectTypeForMessageType, MessageType, type Message } from '$lib/types/Message';
import { CommandedPlanet } from '$lib/types/Planet';
import type { Player } from '$lib/types/Player';
import { emptyUser, User } from '$lib/types/User';
import type { Vector } from '$lib/types/Vector';
import { findIndex, kebabCase } from 'lodash-es';
import type { ComponentType, SvelteComponent } from 'svelte';
import { derived, get, writable } from 'svelte/store';
import type { FullGame } from './FullGame';
import { rollover } from './Math';
import { TechService } from './TechService';
import { Universe } from './Universe';

export type MapObjectsByPosition = {
	[k: string]: MapObject[];
};

export const me = writable<User>(emptyUser);
export const techs = writable<TechService>(new TechService());
export const loadingModalText = writable<string | undefined>(undefined);

export const commandedPlanet = writable<CommandedPlanet | undefined>();
export const commandedFleet = writable<CommandedFleet | undefined>();
export const commandedMapObject = writable<MapObject | undefined>();
export const commandedMapObjectPeers = writable<MapObject[]>([]);
export const selectedMapObject = writable<MapObject | undefined>();
export const selectedMapObjectPeers = writable<MapObject[]>([]);
export const selectedWaypoint = writable<Waypoint | undefined>();
export const highlightedMapObject = writable<MapObject | undefined>();
export const highlightedMapObjectPeers = writable<MapObject[]>([]);
export const mostRecentMapObject = writable<MapObject | undefined>();

export const commandedMapObjectName = writable<string>();
export const zoomTarget = writable<MapObject | undefined>();

export const universe = writable<Universe>(new Universe());

const currentCommandedMapObjectIndex = derived(
	[universe, commandedFleet, commandedPlanet],
	([$universe, $commandedFleet, $commandedPlanet]) => {
		if ($commandedPlanet) {
			return $universe.getMyPlanets().findIndex((p) => p.num === $commandedPlanet.num);
		}
		if ($commandedFleet) {
			return $universe.getMyFleets().findIndex((f) => f.num === $commandedFleet.num);
		}
		return 0;
	}
);

const currentSelectedMapObjectIndex = derived(
	[universe, selectedMapObject],
	([$universe, $selectedMapObject]) => {
		if ($selectedMapObject) {
			const mos = $universe.getMapObjectsByPosition($selectedMapObject.position);
			return findIndex(mos, (mo) => equal($selectedMapObject, mo));
		}
		return -1;
	}
);

// derived store of all commandableMapObjects at the position of the current commandedMapObjects
const commandableMapObjectsAtCommandedMapObjectPosition = derived(
	[universe, commandedMapObject],
	([$universe, $commandedMapObject]) => {
		if ($commandedMapObject) {
			return $universe.getMyMapObjectsByPosition($commandedMapObject);
		}
	}
);

// derived store of the current index
const currentCommandedMapObjectPositionIndex = derived(
	[commandedMapObject, commandableMapObjectsAtCommandedMapObjectPosition],
	([$commandedMapObject, $commandableMapObjectsAtCommandedMapObjectPosition]) => {
		if ($commandedMapObject && $commandableMapObjectsAtCommandedMapObjectPosition) {
			return findIndex($commandableMapObjectsAtCommandedMapObjectPosition, (mo) =>
				equal($commandedMapObject, mo)
			);
		}
		return -1;
	}
);

export const nextCommandableMapObjectAtPosition = () => {
	const index = get(currentCommandedMapObjectPositionIndex);
	const commandable = get(commandableMapObjectsAtCommandedMapObjectPosition);

	if (commandable && commandable?.length > 0) {
		if (index + 1 > commandable.length) {
			commandMapObject(commandable[0]);
		} else {
			commandMapObject(commandable[index + 1]);
		}
	}
};

export const currentSelectedWaypointIndex = derived(
	[selectedWaypoint, commandedFleet],
	([$selectedWaypoint, $commandedFleet]) => {
		if ($selectedWaypoint && $commandedFleet) {
			return findIndex($commandedFleet.waypoints, (wp) => wp === $selectedWaypoint);
		}
		return -1;
	}
);

export const selectNextMapObject = () => {
	const u = get(universe);
	const selected = get(selectedMapObject);
	const index = get(currentSelectedMapObjectIndex);

	if (index != -1 && selected) {
		const mos = u.getMapObjectsByPosition(selected.position);
		if (mos) {
			if (index >= mos.length - 1) {
				selectMapObject(mos[0]);
			} else {
				selectMapObject(mos[index + 1]);
			}
		}
	}
};

// command the previous mapObject for this type, i.e. the previous planet or fleet
export const previousMapObject = () => {
	const u = get(universe);
	const i = get(currentCommandedMapObjectIndex);
	const mo = get(commandedMapObject);

	if (mo) {
		if (mo.type == MapObjectType.Planet) {
			const planets = u.getMyPlanets();
			const prevIndex = rollover(i - 1, 0, planets.length - 1);
			const planet = planets[prevIndex];
			commandMapObject(planet);
			zoomToMapObject(planet);
			selectMapObject(planet);
		} else if (mo.type == MapObjectType.Fleet) {
			const fleets = u.getMyFleets();
			const prevIndex = rollover(i - 1, 0, fleets.length - 1);
			commandMapObject(fleets[prevIndex]);
			zoomToMapObject(fleets[prevIndex]);

			const fleet = fleets[prevIndex];
			if (fleet.orbitingPlanetNum && fleet.orbitingPlanetNum != None) {
				const planet = u.getMapObject({
					targetType: MapObjectType.Planet,
					targetNum: fleet.orbitingPlanetNum
				});
				if (planet) {
					selectMapObject(planet);
				}
			} else {
				selectMapObject(fleet);
			}
		}
	}
};

// command the next mapObject for this type, i.e. the next planet or fleet
export const nextMapObject = () => {
	const u = get(universe);
	const i = get(currentCommandedMapObjectIndex);
	const mo = get(commandedMapObject);

	if (mo) {
		if (mo.type == MapObjectType.Planet) {
			const planets = u.getMyPlanets();
			const nextIndex = rollover(i + 1, 0, u.getMyPlanets().length - 1);
			const planet = planets[nextIndex];
			commandMapObject(planet);
			zoomToMapObject(planet);
			selectMapObject(planet);
		} else if (mo.type == MapObjectType.Fleet) {
			const fleets = u.getMyFleets();

			const nextIndex = rollover(i + 1, 0, fleets.length - 1);
			const fleet = fleets[nextIndex];
			commandMapObject(fleets[nextIndex]);
			zoomToMapObject(fleets[nextIndex]);
			if (fleet.orbitingPlanetNum && fleet.orbitingPlanetNum != None) {
				const planet = u.getMapObject({
					targetType: MapObjectType.Planet,
					targetNum: fleet.orbitingPlanetNum
				});
				if (planet) {
					selectMapObject(planet);
				}
			} else {
				selectMapObject(fleet);
			}
		}
	}
};

export const selectMapObject = (mo: MapObject) => {
	selectedMapObject.update(() => mo);
	mostRecentMapObject.update(() => mo);
};

export const selectWaypoint = (wp: Waypoint) => {
	selectedWaypoint.update(() => wp);
};

export const commandMapObject = (mo: MapObject) => {
	commandedMapObject.update(() => mo);
	mostRecentMapObject.update(() => mo);
	if (mo.type == MapObjectType.Planet) {
		commandedPlanet.update(() => Object.assign(new CommandedPlanet(), mo));
		commandedFleet.update(() => undefined);
	} else if (mo.type == MapObjectType.Fleet) {
		commandedPlanet.update(() => undefined);
		commandedFleet.update(() => Object.assign(new CommandedFleet(), mo));
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
	mostRecentMapObject.update(() => mo);
};

export const setLoadingModalText = (text: string) => {
	loadingModalText.update(() => text);
};

export const clearLoadingModalText = () => {
	loadingModalText.update(() => undefined);
};

export const tooltipComponent = writable<
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	{ component: typeof SvelteComponent; props: any } | undefined
>();
export const tooltipLocation = writable<Vector>({ x: 0, y: 0 });

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const showTooltip = <T>(x: number, y: number, component: ComponentType, props?: T) => {
	tooltipLocation.update(() => ({
		x,
		y
	}));
	tooltipComponent.update(() => ({
		component,
		props
	}));
};

export const popupComponent = writable<
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	{ component: typeof SvelteComponent; props: any } | undefined
>();
export const popupLocation = writable<Vector>({ x: 0, y: 0 });

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const showPopup = <T>(x: number, y: number, component: ComponentType, props?: T) => {
	popupLocation.update(() => ({
		x,
		y
	}));
	popupComponent.update(() => ({
		component,
		props
	}));
};
