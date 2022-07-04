import type { Game } from '$lib/types/Game';
import type { Planet } from '$lib/types/Planet';
import type { Player } from '$lib/types/Player';
import { writable } from 'svelte/store';

export const commandedPlanet = writable<Planet>();
export const game = writable<Game>();
export const player = writable<Player>();
