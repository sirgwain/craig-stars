import { writable } from 'svelte/store';

// set when the server has newer changes for the player than the ones we loaded,
// i.e. the player made more recent changes on another device
export const staleData = writable(false);
