import { CommandedFleet, type Fleet, type Waypoint } from '$lib/types/Fleet';
import { equal, MapObjectType, None, type MapObject } from '$lib/types/MapObject';
import { CommandedPlanet } from '$lib/types/Planet';
import { emptyUser, type User } from '$lib/types/User';
import type { Vector } from '$lib/types/Vector';
import { findIndex } from 'lodash-es';
import type { ComponentType, SvelteComponent } from 'svelte';
import { derived, get, writable } from 'svelte/store';
import type { FullGame } from './FullGame';
import { rollover } from './Math';
import { TechService } from './TechService';

export type MapObjectsByPosition = {
	[k: string]: MapObject[];
};

export const me = writable<User>(emptyUser);
export const game = writable<FullGame | undefined>();
export const techs = writable<TechService>(new TechService());

export const commandedPlanet = writable<CommandedPlanet | undefined>();
export const commandedFleet = writable<CommandedFleet | undefined>();
export const selectedWaypoint = writable<Waypoint | undefined>();
export const selectedMapObject = writable<MapObject | undefined>();
export const commandedMapObject = writable<MapObject | undefined>();
export const highlightedMapObject = writable<MapObject | undefined>();

export const commandedMapObjectName = writable<string>();
export const zoomTarget = writable<MapObject | undefined>();

const currentMapObjectIndex = derived(
	[game, commandedFleet, commandedPlanet],
	([$game, $commandedFleet, $commandedPlanet]) => {
		if ($game) {
			if ($commandedPlanet) {
				return $game.universe.planets.findIndex((p) => p.num === $commandedPlanet.num);
			}
			if ($commandedFleet) {
				return $game.universe.fleets.findIndex((f) => f.num === $commandedFleet.num);
			}
		}
		return 0;
	}
);

const currentSelectedMapObjectIndex = derived(
	[game, selectedMapObject],
	([$game, $selectedMapObject]) => {
		if ($game) {
			if ($selectedMapObject) {
				const mos = $game.universe.getMapObjectsByPosition($selectedMapObject.position);
				return findIndex(mos, (mo) => equal($selectedMapObject, mo));
			}
		}
		return -1;
	}
);

export const selectNextMapObject = () => {
	const g = get(game);
	const selected = get(selectedMapObject);
	const index = get(currentSelectedMapObjectIndex);

	if (index != -1 && selected) {
		const mos = g?.universe.getMapObjectsByPosition(selected.position);
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
	const g = get(game);
	const i = get(currentMapObjectIndex);
	const mo = get(commandedMapObject);

	if (g && mo) {
		if (mo.type == MapObjectType.Planet) {
			const prevIndex = rollover(i - 1, 0, g.universe.planets.length - 1);
			const planet = g.universe.planets[prevIndex];
			commandMapObject(planet);
			zoomToMapObject(planet);
			selectMapObject(planet);
		} else if (mo.type == MapObjectType.Fleet) {
			const prevIndex = rollover(i - 1, 0, g.universe.fleets.length - 1);
			commandMapObject(g.universe.fleets[prevIndex]);
			zoomToMapObject(g.universe.fleets[prevIndex]);

			const fleet = g.universe.fleets[prevIndex];
			if (fleet.orbitingPlanetNum && fleet.orbitingPlanetNum != None) {
				const planet = g.universe.getMapObject(MapObjectType.Planet, fleet.orbitingPlanetNum);
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
	const g = get(game);
	const i = get(currentMapObjectIndex);
	const mo = get(commandedMapObject);

	if (g && mo) {
		if (mo.type == MapObjectType.Planet) {
			const nextIndex = rollover(i + 1, 0, g.universe.planets.length - 1);
			const planet = g.universe.planets[nextIndex];
			commandMapObject(planet);
			zoomToMapObject(planet);
			selectMapObject(planet);
		} else if (mo.type == MapObjectType.Fleet) {
			const nextIndex = rollover(i + 1, 0, g.universe.fleets.length - 1);
			const fleet = g.universe.fleets[nextIndex];
			commandMapObject(g.universe.fleets[nextIndex]);
			zoomToMapObject(g.universe.fleets[nextIndex]);
			if (fleet.orbitingPlanetNum && fleet.orbitingPlanetNum != None) {
				const planet = g.universe.getMapObject(MapObjectType.Planet, fleet.orbitingPlanetNum);
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
};

export const selectWaypoint = (wp: Waypoint) => {
	selectedWaypoint.update(() => wp);
};

export const commandMapObject = (mo: MapObject) => {
	commandedMapObject.update(() => mo);
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
};

export const tooltipComponent = writable<
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	{ component: typeof SvelteComponent; props: any } | undefined
>();
export const tooltipLocation = writable<Vector>({ x: 0, y: 0 });

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const showTooltip = <T>(x: number, y: number, component: ComponentType, props: T) => {
	tooltipLocation.update(() => ({
		x,
		y
	}));
	tooltipComponent.update(() => ({
		component,
		props
	}));
};
