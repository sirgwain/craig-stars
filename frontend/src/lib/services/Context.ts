import type { Game } from '$lib/types/Game';
import type { Planet } from '$lib/types/Planet';
import type { Player } from '$lib/types/Player';
import { writable } from 'svelte/store';

export const commandedPlanet = writable<Planet>();
export const selectedPlanet = writable<Planet>();
export const game = writable<Game>();
export const player = writable<Player>();

export const selectPlanet = (planet: Planet) => {
	selectedPlanet.update((p) => (p = planet));
};

export const commandPlanet = (planet: Planet) => {
	commandedPlanet.update((p) => (p = planet));
};
